package hl7

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

// Parser parses HL7 messages.
type Parser struct {
	delims Delimiters
}

// NewParser creates a new Parser with default delimiters.
func NewParser() *Parser {
	return &Parser{
		delims: DefaultDelimiters(),
	}
}

// NewParserWithDelimiters creates a new Parser with custom delimiters.
func NewParserWithDelimiters(delims Delimiters) *Parser {
	return &Parser{
		delims: delims,
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

	// Normalize line endings: replace \r\n with \r, then split by \r or \n
	data = bytes.ReplaceAll(data, []byte{'\r', '\n'}, []byte{'\r'})
	data = bytes.ReplaceAll(data, []byte{'\n'}, []byte{'\r'})
	data = bytes.TrimSuffix(data, []byte{'\r'})

	// Split by segment terminator
	segments := bytes.Split(data, []byte{'\r'})
	for _, segData := range segments {
		if len(segData) == 0 {
			continue
		}
		// Trim any remaining whitespace (including line feeds)
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
	// Split by field separator
	fields := bytes.Split(data, []byte{byte(p.delims.Field)})

	// Get segment name from first field
	name := string(fields[0])
	if name == "" {
		return Segment{}, &ParseError{Message: "empty segment name"}
	}

	seg := NewSegment(name)

	// For MSH segment, the field separator is field 1
	// MSH fields start at index 2 because index 0 is "MSH" and index 1 is empty (the field separator)
	if name == "MSH" {
		seg.SetField(1, string(p.delims.Field))
		// Start from second field for MSH (field 2 in HL7, index 1 in array)
		for i, field := range fields[1:] {
			seg.SetField(i+2, string(field))
		}
	} else {
		// For other segments, fields start at index 1 (index 0 is segment name)
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

// Scanner scans HL7 messages from a stream.
type Scanner struct {
	reader  *bufio.Reader
	parser  *Parser
	message *Message
	err     error
}

// NewScanner creates a new Scanner.
func NewScanner(reader io.Reader) *Scanner {
	return &Scanner{
		reader: bufio.NewReader(reader),
		parser: NewParser(),
	}
}

// Scan advances the scanner to the next message.
// Returns true if a message was scanned, false at EOF or error.
func (s *Scanner) Scan() bool {
	var buf bytes.Buffer

	for {
		ch, err := s.reader.ReadByte()
		if err != nil {
			if err == io.EOF && buf.Len() > 0 {
				// Final message without trailing CR
				s.message, s.err = s.parser.Parse(buf.Bytes())
				return s.err == nil
			}
			if err != io.EOF {
				s.err = err
			}
			return false
		}

		buf.WriteByte(ch)

		// Check for segment terminator (CR)
		if ch == '\r' {
			// Check if this is the end of the message (empty line)
			nextCh, err := s.reader.Peek(1)
			if err != nil || len(nextCh) == 0 || nextCh[0] == '\r' || nextCh[0] == '\n' {
				// End of message
				s.message, s.err = s.parser.Parse(buf.Bytes())
				return s.err == nil
			}
		}
	}
}

// Message returns the current message.
func (s *Scanner) Message() *Message {
	return s.message
}

// Err returns the first error encountered by the scanner.
func (s *Scanner) Err() error {
	return s.err
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
	// Check for MLLP start byte
	if data[0] == 0x0B {
		// MLLP framed message
		return bytes.Contains(data, []byte("MSH|"))
	}
	// Check for MSH segment
	return bytes.HasPrefix(data, []byte("MSH|")) || bytes.HasPrefix(data, []byte("MSH|"))
}

// StripMLLP removes MLLP framing from data.
func StripMLLP(data []byte) []byte {
	// Remove start byte (0x0B) and end bytes (0x1C, 0x0D)
	data = bytes.TrimPrefix(data, []byte{0x0B})
	data = bytes.TrimSuffix(data, []byte{0x1C, 0x0D})
	return bytes.TrimSpace(data)
}

// Ensure imports are used
var _ = strings.NewReader
