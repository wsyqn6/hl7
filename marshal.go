package hl7

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// Marshaler is the interface implemented by types that can marshal themselves into HL7.
type Marshaler interface {
	MarshalHL7() ([]byte, error)
}

// Unmarshaler is the interface implemented by types that can unmarshal HL7 data.
type Unmarshaler interface {
	UnmarshalHL7(data []byte) error
}

// Unmarshal parses HL7 data and stores the result in the value pointed to by v.
func Unmarshal(data []byte, v interface{}) error {
	msg, err := Parse(data)
	if err != nil {
		return err
	}
	return UnmarshalMessage(msg, v)
}

// UnmarshalMessage unmarshals a Message into the value pointed to by v.
func UnmarshalMessage(msg *Message, v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return fmt.Errorf("unmarshal requires a non-nil pointer")
	}

	return unmarshalValue(msg, val.Elem(), nil)
}

// unmarshalValue recursively unmarshals a Message into a reflect.Value.
func unmarshalValue(msg *Message, val reflect.Value, tag *hl7Tag) error {
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			val.Set(reflect.New(val.Type().Elem()))
		}
		return unmarshalValue(msg, val.Elem(), tag)
	}

	switch val.Kind() {
	case reflect.Struct:
		return unmarshalStruct(msg, val)
	case reflect.String:
		return unmarshalString(val, tag, msg)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return unmarshalInt(val, tag, msg)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return unmarshalUint(val, tag, msg)
	case reflect.Float32, reflect.Float64:
		return unmarshalFloat(val, tag, msg)
	case reflect.Bool:
		return unmarshalBool(val, tag, msg)
	case reflect.Slice:
		return unmarshalSlice(msg, val, tag)
	default:
		// Check for Unmarshaler interface
		if u, ok := val.Interface().(Unmarshaler); ok {
			if tag != nil {
				data := getFieldData(msg, tag)
				return u.UnmarshalHL7([]byte(data))
			}
		}
		return fmt.Errorf("unsupported type: %s", val.Type())
	}
}

// unmarshalStruct unmarshals into a struct.
func unmarshalStruct(msg *Message, val reflect.Value) error {
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Get hl7 tag
		tagStr := field.Tag.Get("hl7")
		tag := parseHL7Tag(tagStr)

		// No tag or embedded struct without tag - recurse into struct
		if tag == nil || tag.segment == "" {
			if field.Type.Kind() == reflect.Struct {
				if err := unmarshalStruct(msg, fieldVal); err != nil {
					return err
				}
			}
			continue
		}

		// Tag with only segment name on a struct - treat as nested struct
		// e.g., hl7:"MSH" on a struct field
		if tag.field == 0 && field.Type.Kind() == reflect.Struct {
			// Recurse into nested struct with the segment context
			if err := unmarshalStruct(msg, fieldVal); err != nil {
				return err
			}
			continue
		}

		if err := unmarshalValue(msg, fieldVal, tag); err != nil {
			return fmt.Errorf("field %s: %w", field.Name, err)
		}
	}
	return nil
}

// hl7Tag represents a parsed hl7 struct tag.
type hl7Tag struct {
	segment string
	field   int
	comp    int
	sub     int
}

// parseHL7Tag parses an hl7 struct tag.
// Supports formats: "segment:field.comp.sub", "segment.field.comp.sub", "PID.3.1"
func parseHL7Tag(tag string) *hl7Tag {
	if tag == "" || tag == "-" {
		return nil
	}

	t := &hl7Tag{}
	// Format: "segment:PID" or "segment:PID.3.1.2" or "PID.3.1"
	parts := splitTag(tag)
	if len(parts) > 0 {
		t.segment = parts[0]
	}
	if len(parts) > 1 {
		t.field, _ = strconv.Atoi(parts[1])
	}
	if len(parts) > 2 {
		t.comp, _ = strconv.Atoi(parts[2])
	}
	if len(parts) > 3 {
		t.sub, _ = strconv.Atoi(parts[3])
	}
	return t
}

// splitTag splits a tag by colon and dot.
func splitTag(tag string) []string {
	var result []string
	current := ""
	for _, r := range tag {
		if r == ':' || r == '.' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(r)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

// getFieldData retrieves field data from a message based on tag.
func getFieldData(msg *Message, tag *hl7Tag) string {
	if tag.segment == "" {
		return ""
	}

	seg, ok := msg.Segment(tag.segment)
	if !ok {
		return ""
	}

	if tag.sub > 0 {
		return seg.SubComponent(tag.field, tag.comp, tag.sub)
	}
	if tag.comp > 0 {
		return seg.Component(tag.field, tag.comp)
	}
	return seg.Field(tag.field)
}

// unmarshalString unmarshals into a string field.
func unmarshalString(val reflect.Value, tag *hl7Tag, msg *Message) error {
	if tag == nil {
		return nil
	}
	data := getFieldData(msg, tag)
	val.SetString(data)
	return nil
}

// unmarshalInt unmarshals into an int field.
func unmarshalInt(val reflect.Value, tag *hl7Tag, msg *Message) error {
	if tag == nil {
		return nil
	}
	data := getFieldData(msg, tag)
	if data == "" {
		return nil
	}
	n, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return fmt.Errorf("parsing int: %w", err)
	}
	val.SetInt(n)
	return nil
}

// unmarshalUint unmarshals into a uint field.
func unmarshalUint(val reflect.Value, tag *hl7Tag, msg *Message) error {
	if tag == nil {
		return nil
	}
	data := getFieldData(msg, tag)
	if data == "" {
		return nil
	}
	n, err := strconv.ParseUint(data, 10, 64)
	if err != nil {
		return fmt.Errorf("parsing uint: %w", err)
	}
	val.SetUint(n)
	return nil
}

// unmarshalFloat unmarshals into a float field.
func unmarshalFloat(val reflect.Value, tag *hl7Tag, msg *Message) error {
	if tag == nil {
		return nil
	}
	data := getFieldData(msg, tag)
	if data == "" {
		return nil
	}
	f, err := strconv.ParseFloat(data, 64)
	if err != nil {
		return fmt.Errorf("parsing float: %w", err)
	}
	val.SetFloat(f)
	return nil
}

// unmarshalBool unmarshals into a bool field.
func unmarshalBool(val reflect.Value, tag *hl7Tag, msg *Message) error {
	if tag == nil {
		return nil
	}
	data := getFieldData(msg, tag)
	if data == "" {
		return nil
	}
	b, err := strconv.ParseBool(data)
	if err != nil {
		return fmt.Errorf("parsing bool: %w", err)
	}
	val.SetBool(b)
	return nil
}

// unmarshalSlice unmarshals into a slice field.
func unmarshalSlice(msg *Message, val reflect.Value, tag *hl7Tag) error {
	if tag == nil {
		return nil
	}
	data := getFieldData(msg, tag)
	if data == "" {
		return nil
	}

	// Split by repetition separator
	repetitions := SplitField(data, '~')
	sliceType := val.Type()
	slice := reflect.MakeSlice(sliceType, len(repetitions), len(repetitions))

	for i, rep := range repetitions {
		elemVal := slice.Index(i)
		// Create a temporary message for this repetition
		tempSeg := NewSegment(tag.segment)
		tempSeg.SetField(tag.field, rep)
		tempMsg := NewMessage()
		tempMsg.AddSegment(tempSeg)

		if err := unmarshalValue(tempMsg, elemVal, &hl7Tag{
			segment: tag.segment,
			field:   tag.field,
			comp:    tag.comp,
			sub:     tag.sub,
		}); err != nil {
			return err
		}
	}

	val.Set(slice)
	return nil
}

// Marshal marshals a value to HL7 format.
func Marshal(v interface{}) ([]byte, error) {
	msg := NewMessage()
	if err := MarshalMessage(v, msg); err != nil {
		return nil, err
	}
	return Encode(msg)
}

// MarshalMessage marshals a value into an HL7 Message.
func MarshalMessage(v interface{}, msg *Message) error {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return fmt.Errorf("marshal requires a struct")
	}

	return marshalStruct(val, msg)
}

// marshalStruct marshals a struct into a Message.
func marshalStruct(val reflect.Value, msg *Message) error {
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		if !field.IsExported() {
			continue
		}

		tagStr := field.Tag.Get("hl7")
		tag := parseHL7Tag(tagStr)

		// No tag - check for embedded struct
		if tag == nil || tag.segment == "" {
			if field.Type.Kind() == reflect.Struct {
				if err := marshalStruct(fieldVal, msg); err != nil {
					return err
				}
			}
			continue
		}

		// Tag with only segment name on a struct - treat as nested struct
		// e.g., hl7:"MSH" on a struct field
		if tag.field == 0 && field.Type.Kind() == reflect.Struct {
			if err := marshalStruct(fieldVal, msg); err != nil {
				return err
			}
			continue
		}

		// Get or create segment
		seg, ok := msg.Segment(tag.segment)
		if !ok {
			seg = NewSegment(tag.segment)
		}

		// Check if field is a struct with nested hl7 tags (component structure)
		if field.Type.Kind() == reflect.Struct {
			// Marshal nested struct to components
			if err := marshalNestedStruct(seg, tag.field, fieldVal); err != nil {
				return fmt.Errorf("field %s: %w", field.Name, err)
			}
		} else {
			// Marshal field value
			value, err := marshalValue(fieldVal)
			if err != nil {
				return fmt.Errorf("field %s: %w", field.Name, err)
			}

			if tag.comp > 0 {
				// Component
				currentField := seg.Field(tag.field)
				components := SplitField(currentField, '^')
				for len(components) < tag.comp {
					components = append(components, "")
				}
				components[tag.comp-1] = value
				seg.SetField(tag.field, joinComponents(components, '^'))
			} else {
				seg.SetField(tag.field, value)
			}
		}

		// Update segment in message
		msg.SetSegment(seg)
	}
	return nil
}

// marshalNestedStruct marshals a nested struct into components of a field.
func marshalNestedStruct(seg Segment, fieldIndex int, val reflect.Value) error {
	typ := val.Type()
	var components []string

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		if !field.IsExported() {
			continue
		}

		// Get component index from tag
		tagStr := field.Tag.Get("hl7")
		if tagStr == "" || tagStr == "-" {
			continue
		}

		// Parse the tag - expect format like "1", "2", etc.
		var compIndex int
		_, err := fmt.Sscanf(tagStr, "%d", &compIndex)
		if err != nil || compIndex < 1 {
			continue
		}

		// Marshal the component value
		value, err := marshalValue(fieldVal)
		if err != nil {
			return fmt.Errorf("component %d: %w", compIndex, err)
		}

		// Ensure we have enough components
		for len(components) < compIndex {
			components = append(components, "")
		}
		components[compIndex-1] = value
	}

	// Set the field
	seg.SetField(fieldIndex, joinComponents(components, '^'))
	return nil
}

// marshalValue marshals a reflect.Value to a string.
func marshalValue(val reflect.Value) (string, error) {
	// Check for Marshaler interface
	if m, ok := val.Interface().(Marshaler); ok {
		data, err := m.MarshalHL7()
		if err != nil {
			return "", err
		}
		return string(data), nil
	}

	switch val.Kind() {
	case reflect.String:
		return val.String(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(val.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(val.Uint(), 10), nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(val.Float(), 'f', -1, 64), nil
	case reflect.Bool:
		return strconv.FormatBool(val.Bool()), nil
	case reflect.Struct:
		// Check for time.Time
		if t, ok := val.Interface().(time.Time); ok {
			return t.Format("20060102150405"), nil
		}
		return "", fmt.Errorf("unsupported struct type: %s", val.Type())
	default:
		return "", fmt.Errorf("unsupported type: %s", val.Type())
	}
}

// joinComponents joins components with the given separator.
func joinComponents(components []string, sep rune) string {
	var result string
	for i, comp := range components {
		if i > 0 {
			result += string(sep)
		}
		result += comp
	}
	return result
}
