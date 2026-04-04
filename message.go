package hl7

import (
	"fmt"
	"strconv"
	"strings"
)

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

// Location represents a parsed HL7 location string.
type Location struct {
	Segment      string
	SegmentIndex int // 0 means first/none
	Field        int
	Component    int
	SubComponent int
	Repetition   int // For repeated fields within a segment
}

// ParseLocation parses a location string into a Location struct.
// Supports formats:
//   - "PID" -> segment only
//   - "PID.3" -> segment.field
//   - "PID.3.1" -> segment.field.component
//   - "PID.3.1.2" -> segment.field.component.subcomponent
//   - "PID[1]" -> first PID segment
//   - "PID[2].3" -> second PID segment, field 3
//   - "PID.3[0]" -> first repetition of field 3
//   - "PID.3[1].1" -> second repetition of field 3, component 1
func ParseLocation(location string) (*Location, error) {
	loc := &Location{}

	if location == "" {
		return nil, fmt.Errorf("empty location")
	}

	parts := strings.Split(location, ".")
	if len(parts) == 0 {
		return nil, fmt.Errorf("invalid location: %s", location)
	}

	segPart := parts[0]

	idxStart := strings.Index(segPart, "[")
	idxEnd := strings.Index(segPart, "]")

	if idxStart >= 0 && idxEnd > idxStart {
		loc.Segment = segPart[:idxStart]
		idxStr := segPart[idxStart+1 : idxEnd]
		loc.SegmentIndex, _ = strconv.Atoi(idxStr)
	} else {
		loc.Segment = segPart
	}

	if loc.Segment == "" {
		return nil, fmt.Errorf("empty segment name")
	}

	if len(parts) > 1 {
		fieldStr := parts[1]

		repStart := strings.Index(fieldStr, "[")
		repEnd := strings.Index(fieldStr, "]")

		var actualFieldStr string
		if repStart >= 0 && repEnd > repStart {
			actualFieldStr = fieldStr[:repStart]
			repStr := fieldStr[repStart+1 : repEnd]
			loc.Repetition, _ = strconv.Atoi(repStr)
		} else {
			actualFieldStr = fieldStr
		}

		loc.Field, _ = strconv.Atoi(actualFieldStr)
		if loc.Field < 1 {
			return nil, fmt.Errorf("invalid field index: %d", loc.Field)
		}
	}

	if len(parts) > 2 {
		compStr := parts[2]
		loc.Component, _ = strconv.Atoi(compStr)
		if loc.Component < 1 {
			return nil, fmt.Errorf("invalid component index: %d", loc.Component)
		}
	}

	if len(parts) > 3 {
		subStr := parts[3]
		loc.SubComponent, _ = strconv.Atoi(subStr)
		if loc.SubComponent < 1 {
			return nil, fmt.Errorf("invalid subcomponent index: %d", loc.SubComponent)
		}
	}

	return loc, nil
}

// Get retrieves a field value by location string.
// Location format: "SEGMENT.FIELD" or "SEGMENT.FIELD.COMPONENT" or "SEGMENT.FIELD.COMPONENT.SUBCOMPONENT"
// For repeated segments, use: "SEGMENT[INDEX].FIELD"
func (m *Message) Get(location string) (string, error) {
	loc, err := ParseLocation(location)
	if err != nil {
		return "", err
	}

	var seg Segment
	var ok bool

	if loc.SegmentIndex > 0 {
		segs := m.Segments(loc.Segment)
		if loc.SegmentIndex-1 < len(segs) {
			seg = segs[loc.SegmentIndex-1]
			ok = true
		}
	} else {
		seg, ok = m.Segment(loc.Segment)
	}

	if !ok {
		return "", fmt.Errorf("segment not found: %s", loc.Segment)
	}

	if loc.Field == 0 {
		return seg.Name(), nil
	}

	fieldValue := seg.Field(loc.Field)
	if fieldValue == "" {
		return "", nil
	}

	if loc.Repetition > 0 {
		repetitions := SplitField(fieldValue, '~')
		if loc.Repetition-1 < len(repetitions) {
			fieldValue = repetitions[loc.Repetition-1]
		} else {
			return "", nil
		}
	}

	if loc.Component == 0 {
		return fieldValue, nil
	}

	components := SplitField(fieldValue, '^')
	if loc.Component > len(components) {
		return "", nil
	}

	compValue := components[loc.Component-1]

	if loc.SubComponent == 0 {
		return compValue, nil
	}

	subComponents := SplitField(compValue, '&')
	if loc.SubComponent > len(subComponents) {
		return "", nil
	}

	return subComponents[loc.SubComponent-1], nil
}

// MustGet retrieves a field value by location string, panics on error.
func (m *Message) MustGet(location string) string {
	val, err := m.Get(location)
	if err != nil {
		panic(err)
	}
	return val
}

// Iterate returns a channel that yields all segments in the message.
func (m *Message) Iterate() <-chan SegmentIterator {
	ch := make(chan SegmentIterator)
	go func() {
		defer close(ch)
		for _, seg := range m.segments {
			ch <- SegmentIterator{
				Name:   seg.Name(),
				Index:  1,
				Fields: seg.Fields(),
				seg:    seg,
			}
		}
	}()
	return ch
}

// SegmentIterator provides iteration over segment fields.
type SegmentIterator struct {
	Name   string
	Index  int
	Fields []string
	seg    Segment
}

// Next advances to the next field in the segment.
func (i *SegmentIterator) Next() bool {
	if i.Index < len(i.Fields) {
		i.Index++
		return true
	}
	return false
}

// Value returns the value of the current field.
func (i *SegmentIterator) Value() string {
	if i.Index >= 1 && i.Index <= len(i.Fields) {
		return i.Fields[i.Index-1]
	}
	return ""
}

// ValueAt returns the value at the specified field index.
func (i *SegmentIterator) ValueAt(index int) string {
	if index >= 1 && index <= len(i.Fields) {
		return i.Fields[index-1]
	}
	return ""
}

// Segment returns the underlying segment.
func (i *SegmentIterator) Segment() Segment {
	return i.seg
}

// Count returns the number of fields in the segment.
func (i *SegmentIterator) Count() int {
	return len(i.Fields)
}

// GetAllRepetitions returns all repetitions of a field value.
func (m *Message) GetAllRepetitions(location string) ([]string, error) {
	loc, err := ParseLocation(location)
	if err != nil {
		return nil, err
	}

	var seg Segment
	var ok bool

	if loc.SegmentIndex > 0 {
		segs := m.Segments(loc.Segment)
		if loc.SegmentIndex-1 < len(segs) {
			seg = segs[loc.SegmentIndex-1]
			ok = true
		}
	} else {
		seg, ok = m.Segment(loc.Segment)
	}

	if !ok {
		return nil, fmt.Errorf("segment not found: %s", loc.Segment)
	}

	if loc.Field == 0 {
		return nil, fmt.Errorf("field index required")
	}

	fieldValue := seg.Field(loc.Field)
	if fieldValue == "" {
		return []string{}, nil
	}

	return SplitField(fieldValue, '~'), nil
}

// CountSegment returns the count of segments with the given name.
func (m *Message) CountSegment(name string) int {
	return len(m.Segments(name))
}

// HasSegment checks if a segment exists in the message.
func (m *Message) HasSegment(name string) bool {
	_, ok := m.Segment(name)
	return ok
}

// GetNthSegment returns the nth segment (1-based index) with the given name.
func (m *Message) GetNthSegment(name string, n int) (Segment, bool) {
	segs := m.Segments(name)
	if n > 0 && n <= len(segs) {
		return segs[n-1], true
	}
	return Segment{}, false
}

// MessageStats contains statistics about a message.
type MessageStats struct {
	SegmentCount int
	SegmentTypes map[string]int
	TotalFields  int
	EmptyFields  int
	HasMSH       bool
	HasPID       bool
	HasPV1       bool
	MessageType  string
	Version      string
}

// Stats returns statistics about the message.
func (m *Message) Stats() MessageStats {
	stats := MessageStats{
		SegmentTypes: make(map[string]int),
	}

	for _, seg := range m.segments {
		stats.SegmentCount++
		stats.SegmentTypes[seg.Name()]++
		stats.TotalFields += seg.FieldCount()

		for _, f := range seg.Fields() {
			if f == "" {
				stats.EmptyFields++
			}
		}

		switch seg.Name() {
		case "MSH":
			stats.HasMSH = true
			stats.MessageType = seg.Field(9)
			stats.Version = seg.Field(12)
		case "PID":
			stats.HasPID = true
		case "PV1":
			stats.HasPV1 = true
		}
	}

	return stats
}

// Summary returns a human-readable summary of the message.
func (m *Message) Summary() string {
	stats := m.Stats()
	var b strings.Builder
	b.WriteString(fmt.Sprintf("HL7 Message (%s)\n", stats.MessageType))
	b.WriteString(fmt.Sprintf("  Version: %s\n", stats.Version))
	b.WriteString(fmt.Sprintf("  Segments: %d\n", stats.SegmentCount))
	b.WriteString(fmt.Sprintf("  Total Fields: %d (Empty: %d)\n", stats.TotalFields, stats.EmptyFields))
	b.WriteString("  Segment Types:\n")
	for name, count := range stats.SegmentTypes {
		b.WriteString(fmt.Sprintf("    %s: %d\n", name, count))
	}
	return b.String()
}
