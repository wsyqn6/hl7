package hl7

// Message represents an HL7 v2.x message containing multiple segments.
type Message struct {
	segments []Segment
	delims   Delimiters
}

// Delimiters contains the HL7 delimiter characters.
type Delimiters struct {
	Field        rune
	Component    rune
	Repetition   rune
	Escape       rune
	SubComponent rune
}

// DefaultDelimiters returns the default HL7 delimiters.
func DefaultDelimiters() Delimiters {
	return Delimiters{
		Field:        '|',
		Component:    '^',
		Repetition:   '~',
		Escape:       '\\',
		SubComponent: '&',
	}
}

// NewMessage creates a new empty Message.
func NewMessage() *Message {
	return &Message{
		segments: make([]Segment, 0),
		delims:   DefaultDelimiters(),
	}
}

// Type returns the message type (e.g., "ADT^A01").
func (m *Message) Type() string {
	if len(m.segments) == 0 {
		return ""
	}
	// MSH segment is always first, field 9 is message type
	if m.segments[0].Name() == "MSH" {
		if field := m.segments[0].Field(9); field != "" {
			return field
		}
	}
	return ""
}

// ControlID returns the message control ID.
func (m *Message) ControlID() string {
	if len(m.segments) == 0 {
		return ""
	}
	// MSH segment is always first, field 10 is control ID
	if m.segments[0].Name() == "MSH" {
		if field := m.segments[0].Field(10); field != "" {
			return field
		}
	}
	return ""
}

// Segment returns the first segment with the given name.
func (m *Message) Segment(name string) (Segment, bool) {
	for _, seg := range m.segments {
		if seg.Name() == name {
			return seg, true
		}
	}
	return Segment{}, false
}

// Segments returns all segments with the given name.
func (m *Message) Segments(name string) []Segment {
	var result []Segment
	for _, seg := range m.segments {
		if seg.Name() == name {
			result = append(result, seg)
		}
	}
	return result
}

// AllSegments returns all segments in the message.
func (m *Message) AllSegments() []Segment {
	return m.segments
}

// AddSegment adds a segment to the message.
func (m *Message) AddSegment(seg Segment) {
	m.segments = append(m.segments, seg)
}

// SetSegment updates or adds a segment to the message.
// If a segment with the same name exists, it updates the first occurrence.
// Otherwise, it adds the segment to the end.
func (m *Message) SetSegment(seg Segment) {
	for i, s := range m.segments {
		if s.Name() == seg.Name() {
			m.segments[i] = seg
			return
		}
	}
	m.segments = append(m.segments, seg)
}
