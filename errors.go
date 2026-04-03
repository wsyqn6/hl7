package hl7

import "fmt"

// ParseError represents an error that occurs during HL7 message parsing.
type ParseError struct {
	Field   string
	Message string
}

// Error implements the error interface.
func (e *ParseError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("parse error at field %s: %s", e.Field, e.Message)
	}
	return fmt.Sprintf("parse error: %s", e.Message)
}

// ValidationError represents a validation error for an HL7 message.
type ValidationError struct {
	Location string
	Message  string
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error at %s: %s", e.Location, e.Message)
}

// FieldError represents an error related to a specific field.
type FieldError struct {
	Segment string
	Field   int
	Err     error
	Value   string
}

// Error implements the error interface.
func (e *FieldError) Error() string {
	return fmt.Sprintf("error in %s.%d: %v (value=%q)", e.Segment, e.Field, e.Err, e.Value)
}

// Unwrap returns the underlying error.
func (e *FieldError) Unwrap() error {
	return e.Err
}

// ScannerError represents an error during message scanning.
type ScannerError struct {
	Message   string
	BytesRead int
	Cause     error
}

// Error implements the error interface.
func (e *ScannerError) Error() string {
	if e.BytesRead > 0 {
		return fmt.Sprintf("scanner error at byte %d: %s", e.BytesRead, e.Message)
	}
	return fmt.Sprintf("scanner error: %s", e.Message)
}

// Unwrap returns the underlying error.
func (e *ScannerError) Unwrap() error {
	return e.Cause
}

// Common scanner errors
var (
	ErrMessageTooLarge = &ScannerError{Message: "message exceeds maximum size limit"}
	ErrInvalidMLLP     = &ScannerError{Message: "invalid MLLP frame"}
)
