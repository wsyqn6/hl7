package hl7

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

const (
	// MLLP framing characters
	// SOH (Start of Heading)
	MLLP_START = 0x0B
	// FS (File Separator)
	MLLP_END = 0x1C
	// CR (Carriage Return)
	MLLP_CR = 0x0D
)

// Handler is a function that handles an HL7 message and returns a response.
type Handler func(ctx context.Context, msg *Message) (*Message, error)

// Server is an MLLP server.
type Server struct {
	addr           string
	handler        Handler
	readTimeout    time.Duration
	writeTimeout   time.Duration
	maxMessageSize int
	listener       net.Listener
	wg             sync.WaitGroup
}

// ServerOption configures an MLLP server.
type ServerOption func(*Server)

// WithReadTimeout sets the read timeout.
func WithReadTimeout(d time.Duration) ServerOption {
	return func(s *Server) {
		s.readTimeout = d
	}
}

// WithWriteTimeout sets the write timeout.
func WithWriteTimeout(d time.Duration) ServerOption {
	return func(s *Server) {
		s.writeTimeout = d
	}
}

// WithMaxMessageSize sets the maximum message size.
func WithMaxMessageSize(size int) ServerOption {
	return func(s *Server) {
		s.maxMessageSize = size
	}
}

// NewServer creates a new MLLP server.
func NewServer(addr string, handler Handler, opts ...ServerOption) *Server {
	s := &Server{
		addr:           addr,
		handler:        handler,
		readTimeout:    60 * time.Second,
		writeTimeout:   30 * time.Second,
		maxMessageSize: 1024 * 1024, // 1MB default
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// ListenAndServe starts the server and listens for connections.
func (s *Server) ListenAndServe() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.listener = ln
	return s.Serve(ln)
}

// Serve accepts connections from the listener and handles them.
func (s *Server) Serve(ln net.Listener) error {
	s.listener = ln
	for {
		conn, err := ln.Accept()
		if err != nil {
			// Check if listener was closed
			select {
			case <-s.done():
				return nil
			default:
				return err
			}
		}
		s.wg.Add(1)
		go s.handleConn(conn)
	}
}

// Stop gracefully stops the server.
func (s *Server) Stop() error {
	if s.listener != nil {
		s.listener.Close()
	}
	s.wg.Wait()
	return nil
}

func (s *Server) done() <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(ch)
	}()
	return ch
}

func (s *Server) handleConn(conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close()

	scanner := bufio.NewReader(conn)
	parser := NewParser()
	encoder := NewEncoder().WithMLLPFraming(true)

	for {
		// Read message with MLLP framing
		if s.readTimeout > 0 {
			conn.SetReadDeadline(time.Now().Add(s.readTimeout))
		}

		data, err := readMLLPMessage(scanner, s.maxMessageSize)
		if err != nil {
			if err == io.EOF {
				return
			}
			// Send NAK on error
			nak, _ := Generate(nil, Reject("Failed to read message"))
			if nakData, err := encoder.Encode(nak); err == nil {
				conn.Write(nakData)
			}
			return
		}

		// Parse message
		msg, err := parser.Parse(data)
		if err != nil {
			nak, _ := Generate(nil, Reject(fmt.Sprintf("Parse error: %v", err)))
			if nakData, err := encoder.Encode(nak); err == nil {
				conn.Write(nakData)
			}
			continue
		}

		// Handle message
		ctx := context.Background()
		if s.readTimeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, s.readTimeout)
			defer cancel()
		}

		response, err := s.handler(ctx, msg)
		if err != nil {
			response, _ = Generate(msg, Error(err.Error()))
		}
		if response == nil {
			response, _ = Generate(msg, Accept())
		}

		// Send response
		if s.writeTimeout > 0 {
			conn.SetWriteDeadline(time.Now().Add(s.writeTimeout))
		}

		respData, err := encoder.Encode(response)
		if err != nil {
			return
		}
		conn.Write(respData)
	}
}

// readMLLPMessage reads an MLLP-framed message.
func readMLLPMessage(r *bufio.Reader, maxSize int) ([]byte, error) {
	// Read until start byte
	for {
		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		if b == MLLP_START {
			break
		}
	}

	// Read until end sequence
	var msg []byte
	for {
		if len(msg) >= maxSize {
			return nil, fmt.Errorf("message too large")
		}

		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}

		msg = append(msg, b)

		// Check for end sequence (FS + CR)
		if len(msg) >= 2 {
			if msg[len(msg)-2] == MLLP_END && msg[len(msg)-1] == MLLP_CR {
				// Remove end sequence
				return msg[:len(msg)-2], nil
			}
		}
	}
}

// Client is an MLLP client.
type Client struct {
	conn    net.Conn
	parser  *Parser
	encoder *Encoder
	timeout time.Duration
}

// ClientOption configures an MLLP client.
type ClientOption func(*Client)

// WithDialTimeout sets the dial timeout.
func WithDialTimeout(d time.Duration) ClientOption {
	return func(c *Client) {
		c.timeout = d
	}
}

// Dial creates a new MLLP client connected to the given address.
func Dial(addr string, opts ...ClientOption) (*Client, error) {
	c := &Client{
		parser:  NewParser(),
		encoder: NewEncoder().WithMLLPFraming(true),
		timeout: 10 * time.Second,
	}
	for _, opt := range opts {
		opt(c)
	}

	conn, err := net.DialTimeout("tcp", addr, c.timeout)
	if err != nil {
		return nil, err
	}
	c.conn = conn
	return c, nil
}

// Send sends a message and returns the response.
func (c *Client) Send(ctx context.Context, msg *Message) (*Message, error) {
	// Encode and send message
	data, err := c.encoder.Encode(msg)
	if err != nil {
		return nil, fmt.Errorf("encode error: %w", err)
	}

	if deadline, ok := ctx.Deadline(); ok {
		c.conn.SetWriteDeadline(deadline)
	}
	_, err = c.conn.Write(data)
	if err != nil {
		return nil, fmt.Errorf("write error: %w", err)
	}

	// Read response
	if deadline, ok := ctx.Deadline(); ok {
		c.conn.SetReadDeadline(deadline)
	}

	respData, err := readMLLPMessage(bufio.NewReader(c.conn), 1024*1024)
	if err != nil {
		return nil, fmt.Errorf("read error: %w", err)
	}

	// Parse response
	resp, err := c.parser.Parse(respData)
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	return resp, nil
}

// Close closes the client connection.
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// DialAndSend is a convenience function to dial and send a message.
func DialAndSend(ctx context.Context, addr string, msg *Message, opts ...ClientOption) (*Message, error) {
	client, err := Dial(addr, opts...)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	return client.Send(ctx, msg)
}
