package hl7

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

const (
	MLLP_START = 0x0B
	MLLP_END   = 0x1C
	MLLP_CR    = 0x0D
)

type Handler func(ctx context.Context, msg *Message) (*Message, error)

type Server struct {
	addr           string
	handler        Handler
	readTimeout    time.Duration
	writeTimeout   time.Duration
	maxMessageSize int
	tlsConfig      *tls.Config
	listener       net.Listener
	wg             sync.WaitGroup
}

type ServerOption func(*Server)

func WithReadTimeout(d time.Duration) ServerOption {
	return func(s *Server) {
		s.readTimeout = d
	}
}

func WithWriteTimeout(d time.Duration) ServerOption {
	return func(s *Server) {
		s.writeTimeout = d
	}
}

func WithServerMaxMessageSize(size int) ServerOption {
	return func(s *Server) {
		s.maxMessageSize = size
	}
}

func WithTLS(config *tls.Config) ServerOption {
	return func(s *Server) {
		s.tlsConfig = config
	}
}

func WithInsecureTLS() ServerOption {
	return func(s *Server) {
		s.tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}
}

func NewServer(addr string, handler Handler, opts ...ServerOption) *Server {
	s := &Server{
		addr:           addr,
		handler:        handler,
		readTimeout:    60 * time.Second,
		writeTimeout:   30 * time.Second,
		maxMessageSize: 1024 * 1024,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Server) ListenAndServe() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.listener = ln
	return s.Serve(ln)
}

func (s *Server) ListenAndServeTLS(certFile, keyFile string) error {
	if s.tlsConfig == nil {
		s.tlsConfig = &tls.Config{}
	}

	var err error
	s.tlsConfig.Certificates = make([]tls.Certificate, 1)
	s.tlsConfig.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return fmt.Errorf("failed to load TLS certificate: %w", err)
	}

	config := s.tlsConfig
	if config.MinVersion == 0 {
		config.MinVersion = tls.VersionTLS12
	}

	ln, err := tls.Listen("tcp", s.addr, config)
	if err != nil {
		return err
	}
	s.listener = ln
	return s.Serve(ln)
}

func (s *Server) Serve(ln net.Listener) error {
	s.listener = ln
	for {
		conn, err := ln.Accept()
		if err != nil {
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
		if s.readTimeout > 0 {
			conn.SetReadDeadline(time.Now().Add(s.readTimeout))
		}

		data, err := readMLLPMessage(scanner, s.maxMessageSize)
		if err != nil {
			if err == io.EOF {
				return
			}
			nak, _ := Generate(nil, Reject("Failed to read message"))
			if nakData, err := encoder.Encode(nak); err == nil {
				conn.Write(nakData)
			}
			return
		}

		msg, err := parser.Parse(data)
		if err != nil {
			nak, _ := Generate(nil, Reject(fmt.Sprintf("Parse error: %v", err)))
			if nakData, err := encoder.Encode(nak); err == nil {
				conn.Write(nakData)
			}
			continue
		}

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

func readMLLPMessage(r *bufio.Reader, maxSize int) ([]byte, error) {
	for {
		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		if b == MLLP_START {
			break
		}
	}

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

		if len(msg) >= 2 {
			if msg[len(msg)-2] == MLLP_END && msg[len(msg)-1] == MLLP_CR {
				return msg[:len(msg)-2], nil
			}
		}
	}
}

type Client struct {
	conn    net.Conn
	parser  *Parser
	encoder *Encoder
	timeout time.Duration
	tls     bool
}

type ClientOption func(*Client)

func WithDialTimeout(d time.Duration) ClientOption {
	return func(c *Client) {
		c.timeout = d
	}
}

func WithTLSClient(tls bool) ClientOption {
	return func(c *Client) {
		c.tls = tls
	}
}

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

func DialTLS(addr string, config *tls.Config, opts ...ClientOption) (*Client, error) {
	c := &Client{
		parser:  NewParser(),
		encoder: NewEncoder().WithMLLPFraming(true),
		timeout: 10 * time.Second,
	}
	for _, opt := range opts {
		opt(c)
	}

	dialer := &net.Dialer{Timeout: c.timeout}
	conn, err := tls.DialWithDialer(dialer, "tcp", addr, config)
	if err != nil {
		return nil, err
	}
	c.conn = conn
	return c, nil
}

func DialAndSend(ctx context.Context, addr string, msg *Message, opts ...ClientOption) (*Message, error) {
	client, err := Dial(addr, opts...)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	return client.Send(ctx, msg)
}

func DialAndSendTLS(ctx context.Context, addr string, msg *Message, config *tls.Config, opts ...ClientOption) (*Message, error) {
	client, err := DialTLS(addr, config, opts...)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	return client.Send(ctx, msg)
}

func (c *Client) Send(ctx context.Context, msg *Message) (*Message, error) {
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

	if deadline, ok := ctx.Deadline(); ok {
		c.conn.SetReadDeadline(deadline)
	}

	respData, err := readMLLPMessage(bufio.NewReader(c.conn), 1024*1024)
	if err != nil {
		return nil, fmt.Errorf("read error: %w", err)
	}

	resp, err := c.parser.Parse(respData)
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	return resp, nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
