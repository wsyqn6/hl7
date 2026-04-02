package hl7

import (
	"fmt"
	"time"
)

// ACK codes
const (
	// AA - Application Accept
	AA = "AA"
	// AE - Application Error
	AE = "AE"
	// AR - Application Reject
	AR = "AR"
)

// ACKOption is a function that configures an ACK message.
type ACKOption func(*ackConfig)

type ackConfig struct {
	code      string
	text      string
	errorCode string
	location  string
}

// Accept creates an option for accepting a message.
func Accept() ACKOption {
	return func(c *ackConfig) {
		c.code = AA
	}
}

// Error creates an option for an error response.
func Error(text string) ACKOption {
	return func(c *ackConfig) {
		c.code = AE
		c.text = text
	}
}

// Reject creates an option for rejecting a message.
func Reject(text string) ACKOption {
	return func(c *ackConfig) {
		c.code = AR
		c.text = text
	}
}

// WithText sets the acknowledgment text.
func WithText(text string) ACKOption {
	return func(c *ackConfig) {
		c.text = text
	}
}

// WithErrorCode sets the error code.
func WithErrorCode(code string) ACKOption {
	return func(c *ackConfig) {
		c.errorCode = code
	}
}

// WithErrorLocation sets the error location.
func WithErrorLocation(location string) ACKOption {
	return func(c *ackConfig) {
		c.location = location
	}
}

// Generate generates an ACK message in response to the given message.
func Generate(original *Message, options ...ACKOption) (*Message, error) {
	if original == nil {
		return nil, fmt.Errorf("original message is nil")
	}

	cfg := &ackConfig{
		code: AA, // Default to accept
	}
	for _, opt := range options {
		opt(cfg)
	}

	ack := NewMessage()

	// Create MSH segment for ACK
	msh := NewSegment("MSH")

	// MSH.1 - Field separator
	msh.SetField(1, "|")
	// MSH.2 - Encoding characters
	msh.SetField(2, "^~\\&")
	// MSH.3 - Sending application (empty for ACK)
	msh.SetField(3, "")
	// MSH.4 - Sending facility (empty for ACK)
	msh.SetField(4, "")
	// MSH.5 - Receiving application (from original MSH.3)
	if origMSH, ok := original.Segment("MSH"); ok {
		msh.SetField(5, origMSH.Field(3))
		// MSH.6 - Receiving facility (from original MSH.4)
		msh.SetField(6, origMSH.Field(4))
	}
	// MSH.7 - Date/Time of message
	msh.SetField(7, time.Now().Format("20060102150405"))
	// MSH.8 - Security (empty)
	msh.SetField(8, "")
	// MSH.9 - Message type (ACK)
	msh.SetField(9, "ACK")
	// MSH.10 - Message control ID
	msh.SetField(10, generateControlID())
	// MSH.11 - Processing ID
	msh.SetField(11, "P")
	// MSH.12 - Version ID
	if origMSH, ok := original.Segment("MSH"); ok {
		msh.SetField(12, origMSH.Field(12))
	}

	ack.AddSegment(msh)

	// Create MSA segment
	msa := NewSegment("MSA")
	// MSA.1 - Acknowledgment code
	msa.SetField(1, cfg.code)
	// MSA.2 - Message control ID (from original message)
	if origMSH, ok := original.Segment("MSH"); ok {
		msa.SetField(2, origMSH.Field(10))
	}
	// MSA.3 - Text message (optional)
	msa.SetField(3, cfg.text)
	// MSA.4 - Expected sequence number (empty)
	msa.SetField(4, "")
	// MSA.5 - Delayed acknowledgment type (empty)
	msa.SetField(5, "")
	// MSA.6 - Error code
	msa.SetField(6, cfg.errorCode)
	// MSA.7 - Error location
	msa.SetField(7, cfg.location)

	ack.AddSegment(msa)

	return ack, nil
}

// GenerateACK is an alias for Generate.
func GenerateACK(original *Message, options ...ACKOption) (*Message, error) {
	return Generate(original, options...)
}

// generateControlID generates a unique control ID.
func generateControlID() string {
	return fmt.Sprintf("ACK%d", time.Now().UnixNano())
}

// IsACK checks if a message is an ACK message.
func IsACK(msg *Message) bool {
	if msg == nil {
		return false
	}
	seg, ok := msg.Segment("MSH")
	if !ok {
		return false
	}
	msgType := seg.Field(9)
	return msgType == "ACK"
}

// GetACKCode returns the ACK code from an ACK message.
func GetACKCode(msg *Message) (string, error) {
	if !IsACK(msg) {
		return "", fmt.Errorf("not an ACK message")
	}
	msa, ok := msg.Segment("MSA")
	if !ok {
		return "", fmt.Errorf("MSA segment not found")
	}
	return msa.Field(1), nil
}

// GetACKText returns the text message from an ACK message.
func GetACKText(msg *Message) (string, error) {
	if !IsACK(msg) {
		return "", fmt.Errorf("not an ACK message")
	}
	msa, ok := msg.Segment("MSA")
	if !ok {
		return "", fmt.Errorf("MSA segment not found")
	}
	return msa.Field(3), nil
}

// IsAccepted checks if an ACK message indicates acceptance.
func IsAccepted(msg *Message) (bool, error) {
	code, err := GetACKCode(msg)
	if err != nil {
		return false, err
	}
	return code == AA, nil
}
