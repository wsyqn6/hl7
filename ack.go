package hl7

import (
	"fmt"
	"strings"
	"time"
)

// ACK codes (HL7 v2.x Standard)
const (
	AA    = "AA" // Application Accept
	AE    = "AE" // Application Error
	AR    = "AR" // Application Reject
	CA    = "CA" // Commit Accept (for Enhanced Acknowledgment Mode)
	ACKCE = "CE" // Commit Error (for Enhanced Acknowledgment Mode)
	CR    = "CR" // Commit Reject (for Enhanced Acknowledgment Mode)
)

// ACK code descriptions
var ackDescriptions = map[string]string{
	AA:    "Application Accept",
	AE:    "Application Error",
	AR:    "Application Reject",
	CA:    "Commit Accept",
	ACKCE: "Commit Error",
	CR:    "Commit Reject",
}

// ACKOption configures an ACK message.
type ACKOption func(*ackConfig)

type ackConfig struct {
	code       string
	text       string
	errorCode  string
	location   string
	errSegment *errConfig
}

type errConfig struct {
	code        string
	severity    string
	diagnostics string
	location    string
	expression  string
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

// CommitAccept creates an option for a commit accept (Enhanced Acknowledgment).
func CommitAccept() ACKOption {
	return func(c *ackConfig) {
		c.code = CA
	}
}

// CommitError creates an option for a commit error (Enhanced Acknowledgment).
func CommitError(text string) ACKOption {
	return func(c *ackConfig) {
		c.code = ACKCE
		c.text = text
	}
}

// CommitReject creates an option for a commit reject (Enhanced Acknowledgment).
func CommitReject(text string) ACKOption {
	return func(c *ackConfig) {
		c.code = CR
		c.text = text
	}
}

// WithCode sets a custom ACK code.
func WithCode(code string) ACKOption {
	return func(c *ackConfig) {
		c.code = code
	}
}

// WithText sets the acknowledgment text.
func WithText(text string) ACKOption {
	return func(c *ackConfig) {
		c.text = text
	}
}

// WithErrorCode sets the error code (HL7 Table 0357 or custom).
func WithErrorCode(code string) ACKOption {
	return func(c *ackConfig) {
		c.errorCode = code
	}
}

// WithErrorLocation sets the error location (e.g., "PID.5").
func WithErrorLocation(location string) ACKOption {
	return func(c *ackConfig) {
		c.location = location
	}
}

// WithERR adds an ERR segment to the ACK message.
// severity: E (Error), W (Warning), I (Information)
func WithERR(code, severity, diagnostics, location string) ACKOption {
	return func(c *ackConfig) {
		c.errSegment = &errConfig{
			code:        code,
			severity:    severity,
			diagnostics: diagnostics,
			location:    location,
		}
	}
}

// WithValidationError adds an ERR segment from a ValidationError.
func WithValidationError(err *ValidationError) ACKOption {
	return func(c *ackConfig) {
		if err != nil {
			c.errSegment = &errConfig{
				code:        "207", // Application internal error
				severity:    "E",
				diagnostics: err.Error(),
				location:    err.Location,
			}
		}
	}
}

// Generate generates an ACK message in response to the given message.
func Generate(original *Message, options ...ACKOption) (*Message, error) {
	if original == nil {
		return nil, fmt.Errorf("original message is nil")
	}

	cfg := &ackConfig{
		code: AA,
	}
	for _, opt := range options {
		opt(cfg)
	}

	ack := NewMessage()

	msh := NewSegment("MSH")
	msh.SetField(1, "|")
	msh.SetField(2, "^~\\&")
	msh.SetField(3, "")
	msh.SetField(4, "")

	if origMSH, ok := original.Segment("MSH"); ok {
		msh.SetField(5, origMSH.Field(3))
		msh.SetField(6, origMSH.Field(4))
	}

	msh.SetField(7, time.Now().Format("20060102150405"))
	msh.SetField(8, "")
	msh.SetField(9, "ACK")
	msh.SetField(10, generateControlID())
	msh.SetField(11, "P")

	if origMSH, ok := original.Segment("MSH"); ok {
		msh.SetField(12, origMSH.Field(12))
	}

	ack.AddSegment(msh)

	msa := NewSegment("MSA")
	msa.SetField(1, cfg.code)

	if origMSH, ok := original.Segment("MSH"); ok {
		msa.SetField(2, origMSH.Field(10))
	}

	text := cfg.text
	if text == "" {
		if desc, ok := ackDescriptions[cfg.code]; ok {
			text = desc
		}
	}
	msa.SetField(3, text)
	msa.SetField(4, "")
	msa.SetField(5, "")
	msa.SetField(6, cfg.errorCode)
	msa.SetField(7, cfg.location)

	ack.AddSegment(msa)

	if cfg.errSegment != nil {
		errSeg := NewSegment("ERR")
		errSeg.SetField(1, cfg.errSegment.location)
		errSeg.SetField(2, cfg.errSegment.code)
		errSeg.SetField(3, "HL70000")
		errSeg.SetField(4, cfg.errSegment.severity)
		errSeg.SetField(5, cfg.errSegment.diagnostics)
		errSeg.SetField(6, "")
		errSeg.SetField(7, "")
		errSeg.SetField(8, "")
		errSeg.SetField(9, cfg.errSegment.expression)
		ack.AddSegment(errSeg)
	}

	return ack, nil
}

// GenerateACK is an alias for Generate.
func GenerateACK(original *Message, options ...ACKOption) (*Message, error) {
	return Generate(original, options...)
}

// GenerateACKWithValidation generates an ACK with multiple validation errors.
func GenerateACKWithValidation(original *Message, errors []*ValidationError) (*Message, error) {
	if original == nil {
		return nil, fmt.Errorf("original message is nil")
	}

	if len(errors) == 0 {
		return Generate(original, Accept())
	}

	cfg := &ackConfig{
		code: AR,
		text: fmt.Sprintf("%d validation error(s)", len(errors)),
	}

	ack := NewMessage()

	msh := NewSegment("MSH")
	msh.SetField(1, "|")
	msh.SetField(2, "^~\\&")
	msh.SetField(3, "")
	msh.SetField(4, "")

	if origMSH, ok := original.Segment("MSH"); ok {
		msh.SetField(5, origMSH.Field(3))
		msh.SetField(6, origMSH.Field(4))
	}

	msh.SetField(7, time.Now().Format("20060102150405"))
	msh.SetField(8, "")
	msh.SetField(9, "ACK")
	msh.SetField(10, generateControlID())
	msh.SetField(11, "P")

	if origMSH, ok := original.Segment("MSH"); ok {
		msh.SetField(12, origMSH.Field(12))
	}

	ack.AddSegment(msh)

	msa := NewSegment("MSA")
	msa.SetField(1, cfg.code)

	if origMSH, ok := original.Segment("MSH"); ok {
		msa.SetField(2, origMSH.Field(10))
	}

	msa.SetField(3, cfg.text)
	ack.AddSegment(msa)

	for i, verr := range errors {
		errSeg := NewSegment("ERR")
		errSeg.SetField(1, verr.Location)
		errSeg.SetField(2, "207")
		errSeg.SetField(3, "HL70000")
		errSeg.SetField(4, "E")
		errSeg.SetField(5, verr.Message)
		errSeg.SetField(6, "")
		errSeg.SetField(7, fmt.Sprintf("%d", i+1))
		errSeg.SetField(8, "")
		errSeg.SetField(9, "")
		ack.AddSegment(errSeg)
	}

	return ack, nil
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
	return strings.HasPrefix(msgType, "ACK")
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
	return code == AA || code == CA, nil
}

// IsRejected checks if an ACK message indicates rejection.
func IsRejected(msg *Message) (bool, error) {
	code, err := GetACKCode(msg)
	if err != nil {
		return false, err
	}
	return code == AR || code == CR, nil
}

// HasError checks if an ACK message indicates an error condition.
func HasError(msg *Message) (bool, error) {
	code, err := GetACKCode(msg)
	if err != nil {
		return false, err
	}
	return code == AE || code == ACKCE, nil
}

// GetErrors returns all ERR segments from an ACK message.
func GetErrors(msg *Message) []ERRInfo {
	var errors []ERRInfo
	if msg == nil {
		return errors
	}

	for _, seg := range msg.Segments("ERR") {
		errors = append(errors, ERRInfo{
			ErrorLocation:       seg.Field(1),
			HL7ErrorCode:        seg.Field(2),
			HL7ErrorCodeSys:     seg.Field(3),
			Severity:            seg.Field(4),
			Diagnostics:         seg.Field(5),
			UserMessage:         seg.Field(6),
			HelpLocation:        seg.Field(7),
			VerboseHelpLocation: seg.Field(8),
			Expression:          seg.Field(9),
		})
	}
	return errors
}

// ERRInfo contains information from an ERR segment.
type ERRInfo struct {
	ErrorLocation       string
	HL7ErrorCode        string
	HL7ErrorCodeSys     string
	Severity            string
	Diagnostics         string
	UserMessage         string
	HelpLocation        string
	VerboseHelpLocation string
	Expression          string
}

// String returns a formatted error string.
func (e ERRInfo) String() string {
	if e.Diagnostics != "" {
		return fmt.Sprintf("%s: %s", e.ErrorLocation, e.Diagnostics)
	}
	if e.UserMessage != "" {
		return fmt.Sprintf("%s: %s", e.ErrorLocation, e.UserMessage)
	}
	return e.ErrorLocation
}
