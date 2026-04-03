package hl7

import (
	"fmt"
	"regexp"
)

// Validator validates HL7 messages against rules.
type Validator struct {
	rules []Rule
}

// Rule represents a validation rule.
type Rule interface {
	Validate(msg *Message) *ValidationError
}

// NewValidator creates a new Validator.
func NewValidator(rules ...Rule) *Validator {
	return &Validator{
		rules: rules,
	}
}

// AddRule adds a validation rule.
func (v *Validator) AddRule(rule Rule) {
	v.rules = append(v.rules, rule)
}

// Validate validates a message and returns all validation errors.
func (v *Validator) Validate(msg *Message) []*ValidationError {
	var errors []*ValidationError
	for _, rule := range v.rules {
		if err := rule.Validate(msg); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

// Required creates a rule that checks if a field is required.
func Required(location string) Rule {
	return &requiredRule{location: location}
}

type requiredRule struct {
	location string
}

func (r *requiredRule) Validate(msg *Message) *ValidationError {
	data, err := msg.Get(r.location)
	if err != nil || data == "" {
		return &ValidationError{
			Location: r.location,
			Message:  "field is required",
		}
	}
	return nil
}

// Value creates a rule that checks if a field has a specific value.
func Value(location string, expected string) Rule {
	return &valueRule{location: location, expected: expected}
}

type valueRule struct {
	location string
	expected string
}

func (r *valueRule) Validate(msg *Message) *ValidationError {
	data, err := msg.Get(r.location)
	if err != nil {
		return &ValidationError{
			Location: r.location,
			Message:  err.Error(),
		}
	}
	if data != r.expected {
		return &ValidationError{
			Location: r.location,
			Message:  fmt.Sprintf("expected %q, got %q", r.expected, data),
		}
	}
	return nil
}

// Pattern creates a rule that checks if a field matches a regex pattern.
func Pattern(location string, pattern string) Rule {
	return &patternRule{location: location, pattern: regexp.MustCompile(pattern)}
}

type patternRule struct {
	location string
	pattern  *regexp.Regexp
}

func (r *patternRule) Validate(msg *Message) *ValidationError {
	data, err := msg.Get(r.location)
	if err != nil {
		return nil // Skip if location doesn't exist
	}
	if data == "" {
		return nil // Skip empty values (use Required if needed)
	}
	if !r.pattern.MatchString(data) {
		return &ValidationError{
			Location: r.location,
			Message:  fmt.Sprintf("does not match pattern %s", r.pattern.String()),
		}
	}
	return nil
}

// OneOf creates a rule that checks if a field value is one of the allowed values.
func OneOf(location string, allowed ...string) Rule {
	return &oneOfRule{location: location, allowed: allowed}
}

type oneOfRule struct {
	location string
	allowed  []string
}

func (r *oneOfRule) Validate(msg *Message) *ValidationError {
	data, err := msg.Get(r.location)
	if err != nil {
		return nil // Skip if location doesn't exist
	}
	if data == "" {
		return nil
	}
	for _, a := range r.allowed {
		if data == a {
			return nil
		}
	}
	return &ValidationError{
		Location: r.location,
		Message:  fmt.Sprintf("must be one of %v, got %q", r.allowed, data),
	}
}

// MinLength creates a rule that checks the minimum length of a field.
func MinLength(location string, min int) Rule {
	return &minLengthRule{location: location, min: min}
}

type minLengthRule struct {
	location string
	min      int
}

func (r *minLengthRule) Validate(msg *Message) *ValidationError {
	data, err := msg.Get(r.location)
	if err != nil {
		return nil
	}
	if len(data) < r.min {
		return &ValidationError{
			Location: r.location,
			Message:  fmt.Sprintf("minimum length is %d, got %d", r.min, len(data)),
		}
	}
	return nil
}

// MaxLength creates a rule that checks the maximum length of a field.
func MaxLength(location string, max int) Rule {
	return &maxLengthRule{location: location, max: max}
}

type maxLengthRule struct {
	location string
	max      int
}

func (r *maxLengthRule) Validate(msg *Message) *ValidationError {
	data, err := msg.Get(r.location)
	if err != nil {
		return nil
	}
	if len(data) > r.max {
		return &ValidationError{
			Location: r.location,
			Message:  fmt.Sprintf("maximum length is %d, got %d", r.max, len(data)),
		}
	}
	return nil
}

// Length creates a rule that checks the length range of a field.
func Length(location string, min, max int) Rule {
	return &lengthRule{location: location, min: min, max: max}
}

type lengthRule struct {
	location string
	min      int
	max      int
}

func (r *lengthRule) Validate(msg *Message) *ValidationError {
	data, err := msg.Get(r.location)
	if err != nil {
		return nil
	}
	if len(data) < r.min || len(data) > r.max {
		return &ValidationError{
			Location: r.location,
			Message:  fmt.Sprintf("length must be between %d and %d, got %d", r.min, r.max, len(data)),
		}
	}
	return nil
}

// Custom creates a custom validation rule.
func Custom(location string, fn func(value string) error) Rule {
	return &customRule{location: location, fn: fn}
}

type customRule struct {
	location string
	fn       func(value string) error
}

func (r *customRule) Validate(msg *Message) *ValidationError {
	data, err := msg.Get(r.location)
	if err != nil {
		return nil
	}
	if err := r.fn(data); err != nil {
		return &ValidationError{
			Location: r.location,
			Message:  err.Error(),
		}
	}
	return nil
}
