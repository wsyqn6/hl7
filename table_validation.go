package hl7

import (
	"fmt"
	"time"
)

type TableRule struct {
	location string
	tableID  string
	required bool
}

func Table(location, tableID string) Rule {
	return &TableRule{
		location: location,
		tableID:  tableID,
		required: false,
	}
}

func RequiredTable(location, tableID string) Rule {
	return &TableRule{
		location: location,
		tableID:  tableID,
		required: true,
	}
}

func (r *TableRule) Validate(msg *Message) *ValidationError {
	data, err := msg.Get(r.location)
	if err != nil {
		if r.required {
			return &ValidationError{
				Location: r.location,
				Message:  fmt.Sprintf("required field not found: %s", r.location),
			}
		}
		return nil
	}

	if data == "" {
		if r.required {
			return &ValidationError{
				Location: r.location,
				Message:  fmt.Sprintf("field is required: %s", r.location),
			}
		}
		return nil
	}

	table, ok := HL7Tables[r.tableID]
	if !ok {
		return nil
	}

	if _, valid := table.Values[data]; !valid {
		return &ValidationError{
			Location: r.location,
			Message:  fmt.Sprintf("invalid value %q for table %s; valid values: %v", data, r.tableID, r.getValidValues(table)),
		}
	}

	return nil
}

func (r *TableRule) getValidValues(table HL7Table) []string {
	values := make([]string, 0, len(table.Values))
	for v := range table.Values {
		values = append(values, v)
	}
	return values
}

type VersionRule struct {
	requiredVersions []string
}

func SupportedVersion(versions ...string) Rule {
	return &VersionRule{
		requiredVersions: versions,
	}
}

func (r *VersionRule) Validate(msg *Message) *ValidationError {
	msh, ok := msg.Segment("MSH")
	if !ok {
		return &ValidationError{
			Location: "MSH.12",
			Message:  "MSH segment not found",
		}
	}

	version := msh.Field(12)
	if version == "" {
		return &ValidationError{
			Location: "MSH.12",
			Message:  "version ID is required",
		}
	}

	for _, v := range r.requiredVersions {
		if version == v {
			return nil
		}
	}

	return &ValidationError{
		Location: "MSH.12",
		Message:  fmt.Sprintf("unsupported version %q; supported versions: %v", version, r.requiredVersions),
	}
}

type MessageTypeRule struct {
	messageTypes []string
}

func SupportedMessageType(types ...string) Rule {
	return &MessageTypeRule{
		messageTypes: types,
	}
}

func (r *MessageTypeRule) Validate(msg *Message) *ValidationError {
	msh, ok := msg.Segment("MSH")
	if !ok {
		return &ValidationError{
			Location: "MSH.9",
			Message:  "MSH segment not found",
		}
	}

	msgType := msh.Field(9)
	if msgType == "" {
		return &ValidationError{
			Location: "MSH.9",
			Message:  "message type is required",
		}
	}

	for _, t := range r.messageTypes {
		if msgType == t {
			return nil
		}
	}

	return &ValidationError{
		Location: "MSH.9",
		Message:  fmt.Sprintf("unsupported message type %q; supported: %v", msgType, r.messageTypes),
	}
}

type MessageStructureRule struct {
	expectedStructure string
}

func ExpectedStructure(messageType string) Rule {
	return &MessageStructureRule{
		expectedStructure: messageType,
	}
}

func (r *MessageStructureRule) Validate(msg *Message) *ValidationError {
	ms, ok := LookupMessageStructure(r.expectedStructure)
	if !ok {
		return nil
	}

	var errors []*ValidationError

	for _, req := range ms.Segments {
		segments := msg.Segments(req.Name)
		count := len(segments)

		if count == 0 && req.IsRequired {
			errors = append(errors, &ValidationError{
				Location: req.Name,
				Message:  fmt.Sprintf("required segment %s not found in message structure %s", req.Name, r.expectedStructure),
			})
			continue
		}

		if req.MaxOccurrence > 0 && count > req.MaxOccurrence {
			errors = append(errors, &ValidationError{
				Location: req.Name,
				Message:  fmt.Sprintf("segment %s occurs %d times, maximum allowed is %d", req.Name, count, req.MaxOccurrence),
			})
		}
	}

	if len(errors) > 0 {
		return errors[0]
	}
	return nil
}

type CompositeRule struct {
	rules []Rule
}

func AllOf(rules ...Rule) Rule {
	return &CompositeRule{rules: rules}
}

func (r *CompositeRule) Validate(msg *Message) *ValidationError {
	for _, rule := range r.rules {
		if err := rule.Validate(msg); err != nil {
			return err
		}
	}
	return nil
}

type AnyOfRule struct {
	rules []Rule
}

func AnyOf(rules ...Rule) Rule {
	return &AnyOfRule{rules: rules}
}

func (r *AnyOfRule) Validate(msg *Message) *ValidationError {
	if len(r.rules) == 0 {
		return nil
	}

	for _, rule := range r.rules {
		if err := rule.Validate(msg); err == nil {
			return nil
		}
	}

	return &ValidationError{
		Location: "",
		Message:  "none of the alternative rules matched",
	}
}

type NotRule struct {
	rule Rule
}

func Not(rule Rule) Rule {
	return &NotRule{rule: rule}
}

func (r *NotRule) Validate(msg *Message) *ValidationError {
	if err := r.rule.Validate(msg); err == nil {
		return &ValidationError{
			Location: "",
			Message:  "rule should not have matched but did",
		}
	}
	return nil
}

type WhenRule struct {
	condition func(*Message) bool
	rule      Rule
}

func When(condition func(*Message) bool, rule Rule) Rule {
	return &WhenRule{
		condition: condition,
		rule:      rule,
	}
}

func (r *WhenRule) Validate(msg *Message) *ValidationError {
	if r.condition(msg) {
		return r.rule.Validate(msg)
	}
	return nil
}

type RangeRule struct {
	location   string
	min        float64
	max        float64
	allowEmpty bool
}

func Range(location string, min, max float64) Rule {
	return &RangeRule{
		location: location,
		min:      min,
		max:      max,
	}
}

func (r *RangeRule) Validate(msg *Message) *ValidationError {
	data, err := msg.Get(r.location)
	if err != nil {
		if !r.allowEmpty {
			return &ValidationError{
				Location: r.location,
				Message:  err.Error(),
			}
		}
		return nil
	}

	if data == "" {
		return nil
	}

	var value float64
	_, err = fmt.Sscanf(data, "%f", &value)
	if err != nil {
		return &ValidationError{
			Location: r.location,
			Message:  fmt.Sprintf("invalid numeric value: %s", data),
		}
	}

	if value < r.min || value > r.max {
		return &ValidationError{
			Location: r.location,
			Message:  fmt.Sprintf("value %f is outside range [%f, %f]", value, r.min, r.max),
		}
	}

	return nil
}

type DateFormatRule struct {
	location string
	format   string
}

func DateFormat(location, format string) Rule {
	return &DateFormatRule{
		location: location,
		format:   format,
	}
}

func (r *DateFormatRule) Validate(msg *Message) *ValidationError {
	data, err := msg.Get(r.location)
	if err != nil {
		return nil
	}

	if data == "" {
		return nil
	}

	_, err = time.Parse(r.format, data)
	if err != nil {
		return &ValidationError{
			Location: r.location,
			Message:  fmt.Sprintf("invalid date format; expected %s, got %s", r.format, data),
		}
	}

	return nil
}

func ValidateWithSchema(msg *Message, messageType string) []*ValidationError {
	var errors []*ValidationError

	ms, ok := LookupMessageStructure(messageType)
	if !ok {
		return errors
	}

	msh, ok := msg.Segment("MSH")
	if ok {
		version := msh.Field(12)
		if version != "" && ms.Version != "" && version != ms.Version {
			errors = append(errors, &ValidationError{
				Location: "MSH.12",
				Message:  fmt.Sprintf("version mismatch: expected %s, got %s", ms.Version, version),
			})
		}
	}

	for _, req := range ms.Segments {
		segments := msg.Segments(req.Name)
		count := len(segments)

		if count == 0 && req.IsRequired {
			errors = append(errors, &ValidationError{
				Location: req.Name,
				Message:  fmt.Sprintf("required segment %s not found", req.Name),
			})
			continue
		}

		if req.MaxOccurrence > 0 && count > req.MaxOccurrence {
			errors = append(errors, &ValidationError{
				Location: req.Name,
				Message:  fmt.Sprintf("too many occurrences of %s: %d (max: %d)", req.Name, count, req.MaxOccurrence),
			})
		}
	}

	return errors
}
