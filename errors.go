package hl7

import "fmt"

// ParseError represents an error that occurs during HL7 message parsing.
type ParseError struct {
	Line    int
	Field   string
	Message string
}

// Error implements the error interface.
func (e *ParseError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("parse error at line %d, field %s: %s", e.Line, e.Field, e.Message)
	}
	return fmt.Sprintf("parse error at field %s: %s", e.Field, e.Message)
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
