package hl7

// Segment represents an HL7 segment (e.g., MSH, PID, OBX).
type Segment struct {
	name   string
	fields []string
}

// NewSegment creates a new Segment with the given name.
func NewSegment(name string) Segment {
	return Segment{
		name:   name,
		fields: make([]string, 0),
	}
}

// Name returns the segment name.
func (s Segment) Name() string {
	return s.name
}

// Field returns the field at the given index (1-based).
// For MSH segment, field 1 is the field separator.
func (s Segment) Field(index int) string {
	if index < 1 || index > len(s.fields) {
		return ""
	}
	return s.fields[index-1]
}

// FieldCount returns the number of fields in the segment.
func (s Segment) FieldCount() int {
	return len(s.fields)
}

// SetField sets the field at the given index (1-based).
// Creates empty fields if necessary.
func (s *Segment) SetField(index int, value string) {
	for len(s.fields) < index {
		s.fields = append(s.fields, "")
	}
	s.fields[index-1] = value
}

// Fields returns all fields in the segment.
func (s Segment) Fields() []string {
	return s.fields
}

// Component returns the component at the given index within a field.
func (s Segment) Component(fieldIndex, componentIndex int) string {
	field := s.Field(fieldIndex)
	if field == "" {
		return ""
	}
	return ParseComponent(field, componentIndex)
}

// SubComponent returns the subcomponent at the given indices.
func (s Segment) SubComponent(fieldIndex, componentIndex, subComponentIndex int) string {
	field := s.Field(fieldIndex)
	if field == "" {
		return ""
	}
	components := SplitField(field, '^')
	if componentIndex < 1 || componentIndex > len(components) {
		return ""
	}
	subComponents := SplitField(components[componentIndex-1], '&')
	if subComponentIndex < 1 || subComponentIndex > len(subComponents) {
		return ""
	}
	return subComponents[subComponentIndex-1]
}

// ParseComponent extracts a component from a field value.
func ParseComponent(field string, index int) string {
	components := SplitField(field, '^')
	if index < 1 || index > len(components) {
		return ""
	}
	return components[index-1]
}

// SplitField splits a field by the given separator.
func SplitField(value string, separator rune) []string {
	if value == "" {
		return nil
	}
	var result []string
	current := ""
	for _, r := range value {
		if r == separator {
			result = append(result, current)
			current = ""
		} else {
			current += string(r)
		}
	}
	result = append(result, current)
	return result
}

// ParseSegment parses a segment string into a Segment.
func ParseSegment(data string) *Segment {
	if data == "" {
		return nil
	}
	delims := DefaultDelimiters()
	fields := SplitField(data, delims.Field)
	if len(fields) == 0 {
		return nil
	}
	seg := NewSegment(fields[0])
	if fields[0] == "MSH" {
		seg.SetField(1, string(delims.Field))
		for i, field := range fields[1:] {
			seg.SetField(i+2, field)
		}
	} else {
		for i, field := range fields[1:] {
			seg.SetField(i+1, field)
		}
	}
	return &seg
}

// Field represents a parsed HL7 field.
type Field struct {
	Value      string
	Components []string
}

// ParseField parses a field string into a Field struct with components.
func ParseField(data string) Field {
	components := SplitField(data, '^')
	value := data
	if len(components) > 0 {
		value = components[0]
	}
	return Field{
		Value:      value,
		Components: components,
	}
}

// ParseComponents parses a field string into its components.
func ParseComponents(data string) []string {
	return SplitField(data, '^')
}
