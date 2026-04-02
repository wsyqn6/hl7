package hl7

import (
	"bufio"
	"bytes"
	"io"
)

// Default values for Scanner configuration
const (
	DefaultMaxMessageSize = 64 * 1024 // 64KB
	DefaultMaxSegments    = 1000
	DefaultMaxFieldLength = 32 * 1024 // 32KB
)

// Parser parses HL7 messages.
type Parser struct {
	delims Delimiters
	config ParserConfig
}

// ParserConfig holds configuration for DoS protection.
type ParserConfig struct {
	MaxSegments    int // Maximum number of segments allowed
	MaxFieldLength int // Maximum length of a field
}

// NewParser creates a new Parser with default settings.
func NewParser() *Parser {
	return &Parser{
		delims: DefaultDelimiters(),
		config: ParserConfig{
			MaxSegments:    DefaultMaxSegments,
			MaxFieldLength: DefaultMaxFieldLength,
		},
	}
}

// NewParserWithDelimiters creates a new Parser with custom delimiters.
func NewParserWithDelimiters(delims Delimiters) *Parser {
	return &Parser{
		delims: delims,
		config: ParserConfig{
			MaxSegments:    DefaultMaxSegments,
			MaxFieldLength: DefaultMaxFieldLength,
		},
	}
}

// NewParserWithConfig creates a new Parser with custom configuration.
func NewParserWithConfig(config ParserConfig) *Parser {
	if config.MaxSegments == 0 {
		config.MaxSegments = DefaultMaxSegments
	}
	if config.MaxFieldLength == 0 {
		config.MaxFieldLength = DefaultMaxFieldLength
	}
	return &Parser{
		delims: DefaultDelimiters(),
		config: config,
	}
}

// Parse parses raw HL7 data into a Message.
func (p *Parser) Parse(data []byte) (*Message, error) {
	if len(data) == 0 {
		return nil, &ParseError{Message: "empty data"}
	}

	msg := &Message{
		segments: make([]Segment, 0),
		delims:   p.delims,
	}

	data = bytes.ReplaceAll(data, []byte{'\r', '\n'}, []byte{'\r'})
	data = bytes.ReplaceAll(data, []byte{'\n'}, []byte{'\r'})
	data = bytes.TrimSuffix(data, []byte{'\r'})

	segments := bytes.Split(data, []byte{'\r'})
	if len(segments) > p.config.MaxSegments {
		return nil, &ParseError{Message: "too many segments"}
	}

	for _, segData := range segments {
		if len(segData) == 0 {
			continue
		}
		segData = bytes.TrimSpace(segData)
		if len(segData) == 0 {
			continue
		}
		seg, err := p.parseSegment(segData)
		if err != nil {
			return nil, err
		}
		msg.segments = append(msg.segments, seg)
	}

	return msg, nil
}

// parseSegment parses a single segment.
func (p *Parser) parseSegment(data []byte) (Segment, error) {
	fields := bytes.Split(data, []byte{byte(p.delims.Field)})

	name := string(fields[0])
	if name == "" {
		return Segment{}, &ParseError{Message: "empty segment name"}
	}

	if len(fields) > p.config.MaxFieldLength {
		return Segment{}, &ParseError{Field: name, Message: "too many fields"}
	}

	seg := NewSegment(name)

	if name == "MSH" {
		seg.SetField(1, string(p.delims.Field))
		for i, field := range fields[1:] {
			seg.SetField(i+2, string(field))
		}
	} else {
		for i, field := range fields[1:] {
			seg.SetField(i+1, string(field))
		}
	}

	return seg, nil
}

// Parse parses raw HL7 data into a Message (package-level function).
func Parse(data []byte) (*Message, error) {
	p := NewParser()
	return p.Parse(data)
}

// ParseString parses an HL7 message from a string.
func ParseString(data string) (*Message, error) {
	return Parse([]byte(data))
}

// ScannerConfig holds configuration for the Scanner.
type ScannerConfig struct {
	MaxMessageSize int        // Maximum message size in bytes
	Delimiters     Delimiters // Custom delimiters
	SkipInvalid    bool       // Skip invalid messages and continue scanning
}

// ScannerOption configures a Scanner.
type ScannerOption func(*ScannerConfig)

// WithMaxMessageSize sets the maximum message size for the Scanner.
// Default is 64KB.
func WithMaxMessageSize(size int) ScannerOption {
	return func(cfg *ScannerConfig) {
		if size > 0 {
			cfg.MaxMessageSize = size
		}
	}
}

// WithSkipInvalid configures the scanner to skip invalid messages.
// When true, the scanner will continue scanning after encountering
// a parsing error and report the error via Err() at the end.
func WithSkipInvalid(skip bool) ScannerOption {
	return func(cfg *ScannerConfig) {
		cfg.SkipInvalid = skip
	}
}

// WithDelimiters sets custom delimiters for parsing.
func WithDelimiters(d Delimiters) ScannerOption {
	return func(cfg *ScannerConfig) {
		cfg.Delimiters = d
	}
}

// Scanner scans HL7 messages from a stream.
type Scanner struct {
	cfg     ScannerConfig
	reader  *bufio.Reader
	parser  *Parser
	message *Message
	err     error
	count   int // Count of messages scanned
}

// NewScanner creates a new Scanner with default configuration.
func NewScanner(reader io.Reader) *Scanner {
	return NewScannerWithOptions(reader)
}

// NewScannerWithOptions creates a new Scanner with custom options.
func NewScannerWithOptions(reader io.Reader, opts ...ScannerOption) *Scanner {
	cfg := ScannerConfig{
		MaxMessageSize: DefaultMaxMessageSize,
		Delimiters:     DefaultDelimiters(),
		SkipInvalid:    false,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	return &Scanner{
		cfg:    cfg,
		reader: bufio.NewReaderSize(reader, 4096),
		parser: NewParserWithDelimiters(cfg.Delimiters),
	}
}

// Scan advances the scanner to the next message.
// Returns true if a message was scanned, false at EOF or on error.
func (s *Scanner) Scan() bool {
	if s.err != nil && !s.cfg.SkipInvalid {
		return false
	}

	if s.cfg.SkipInvalid && s.err != nil {
		s.err = nil
	}

	var lineBuf bytes.Buffer
	emptyLines := 0
	inMLLP := false

	for {
		b, err := s.reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				if lineBuf.Len() > 0 {
					s.parseMessage(lineBuf.Bytes())
					return s.err == nil || s.cfg.SkipInvalid
				}
				return false
			}
			s.err = &ScannerError{Message: "read error", BytesRead: lineBuf.Len(), Cause: err}
			return false
		}

		// Check message size limit
		if lineBuf.Len() >= s.cfg.MaxMessageSize {
			s.err = &ScannerError{Message: "message exceeds maximum size", BytesRead: lineBuf.Len()}
			return false
		}

		// Handle MLLP start
		if !inMLLP && b == MLLP_START {
			inMLLP = true
			continue
		}

		if inMLLP {
			if b == MLLP_END {
				// Read one more byte (should be CR)
				cr, err := s.reader.ReadByte()
				if err == nil && cr == MLLP_CR {
					s.parseMessage(lineBuf.Bytes())
					lineBuf.Reset()
					inMLLP = false
					return s.err == nil || s.cfg.SkipInvalid
				}
				// Not valid MLLP, include the bytes
				lineBuf.WriteByte(MLLP_END)
				if err == nil {
					lineBuf.WriteByte(cr)
				}
				continue
			}
			lineBuf.WriteByte(b)
			continue
		}

		// Handle regular line endings
		if b == '\r' {
			// Check for CR+LF
			peek, err := s.reader.Peek(1)
			if err == nil && len(peek) > 0 && peek[0] == '\n' {
				// Consume the LF
				s.reader.Discard(1)
			}

			if lineBuf.Len() == 0 {
				emptyLines++
				if emptyLines >= 1 {
					// Empty line after content - message boundary
					if s.count > 0 || lineBuf.Len() > 0 {
						// We have a previous message
					}
				}
				continue
			}

			if lineBuf.Len() > 0 {
				// End of line
				s.parseMessage(lineBuf.Bytes())
				lineBuf.Reset()
				return s.err == nil || s.cfg.SkipInvalid
			}
			continue
		}

		if b == '\n' {
			if lineBuf.Len() == 0 {
				emptyLines++
				continue
			}

			if lineBuf.Len() > 0 {
				// Check for MSH segment in buffer
				if bytes.HasPrefix(lineBuf.Bytes(), []byte("MSH")) || s.count > 0 {
					// Check for double newline (empty line = message boundary)
					// Peek at next character
					peek, err := s.reader.Peek(1)
					if err == nil && len(peek) > 0 {
						if peek[0] == '\n' || peek[0] == '\r' {
							// This is a message boundary
							s.parseMessage(lineBuf.Bytes())
							lineBuf.Reset()
							return s.err == nil || s.cfg.SkipInvalid
						}
					}
				}
				lineBuf.WriteByte(b)
			}
			continue
		}

		// Regular character
		lineBuf.WriteByte(b)
		emptyLines = 0
	}
}

// parseMessage parses the current line buffer as an HL7 message.
func (s *Scanner) parseMessage(data []byte) {
	// Build full message from lines
	var fullMsg bytes.Buffer
	for _, line := range bytes.Split(data, []byte{'\n'}) {
		line = bytes.TrimSpace(line)
		if len(line) > 0 {
			fullMsg.Write(line)
			fullMsg.WriteByte('\r')
		}
	}

	msgData := fullMsg.Bytes()
	if len(msgData) > 0 {
		msgData = bytes.TrimSuffix(msgData, []byte{'\r'})
	}

	if len(msgData) == 0 {
		return
	}

	msg, err := s.parser.Parse(msgData)
	if err != nil {
		s.err = err
		return
	}

	// Verify it's an HL7 message (has MSH segment)
	if len(msg.AllSegments()) == 0 {
		if s.cfg.SkipInvalid {
			return
		}
		s.err = &ParseError{Message: "no segments found"}
		return
	}

	s.message = msg
	s.count++
}

// Message returns the current message.
func (s *Scanner) Message() *Message {
	return s.message
}

// Err returns the first error encountered by the scanner.
func (s *Scanner) Err() error {
	return s.err
}

// Count returns the number of messages successfully scanned.
func (s *Scanner) Count() int {
	return s.count
}

// Reset resets the scanner to scan from the beginning.
func (s *Scanner) Reset(reader io.Reader) {
	s.reader = bufio.NewReaderSize(reader, 4096)
	s.message = nil
	s.err = nil
	s.count = 0
}

// SplitMessages splits multiple HL7 messages from a byte slice.
func SplitMessages(data []byte) [][]byte {
	// HL7 messages are separated by double carriage returns
	parts := bytes.Split(data, []byte{'\r', '\r'})
	var messages [][]byte
	for _, part := range parts {
		trimmed := bytes.TrimSpace(part)
		if len(trimmed) > 0 {
			messages = append(messages, trimmed)
		}
	}
	return messages
}

// IsHL7Message checks if the data appears to be an HL7 message.
func IsHL7Message(data []byte) bool {
	if len(data) < 4 {
		return false
	}
	if data[0] == MLLP_START {
		return bytes.Contains(data, []byte("MSH|"))
	}
	return bytes.HasPrefix(data, []byte("MSH|"))
}

// StripMLLP removes MLLP framing from data.
func StripMLLP(data []byte) []byte {
	data = bytes.TrimPrefix(data, []byte{MLLP_START})
	data = bytes.TrimSuffix(data, []byte{MLLP_END, MLLP_CR})
	return bytes.TrimSpace(data)
}
