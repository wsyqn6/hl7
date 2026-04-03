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

// Get retrieves a field value by location string.
// Location format: "SEGMENT.FIELD" or "SEGMENT.FIELD.COMPONENT" or "SEGMENT.FIELD.COMPONENT.SUBCOMPONENT"
// For repeated segments, use: "SEGMENT[INDEX].FIELD"
func (m *Message) Get(location string) (string, error) {
	parts := SplitField(location, '.')
	if len(parts) == 0 {
		return "", fmt.Errorf("invalid location: %s", location)
	}

	segmentName := parts[0]

	// Check for repeated segment index
	var segIndex int
	var segNamePart string
	segNamePart = segmentName
	if len(segmentName) > 0 {
		for i := 0; i < len(segmentName); i++ {
			if segmentName[i] == '[' {
				segNamePart = segmentName[:i]
				break
			}
		}
	}

	// Parse index if present
	if idxStart := strings.Index(segmentName, "["); idxStart >= 0 {
		segNamePart = segmentName[:idxStart]
		if idxEnd := strings.Index(segmentName, "]"); idxEnd > idxStart {
			idxStr := segmentName[idxStart+1 : idxEnd]
			segIndex, _ = strconv.Atoi(idxStr)
		}
	}

	var seg Segment
	var ok bool

	if segIndex > 0 {
		segs := m.Segments(segNamePart)
		if segIndex-1 < len(segs) {
			seg = segs[segIndex-1]
			ok = true
		}
	} else {
		seg, ok = m.Segment(segNamePart)
	}

	if !ok {
		return "", fmt.Errorf("segment not found: %s", segNamePart)
	}

	if len(parts) == 1 {
		return seg.Name(), nil
	}

	// Field index
	fieldStr := parts[1]
	fieldIdx := 0
	for _, c := range fieldStr {
		if c < '0' || c > '9' {
			return "", fmt.Errorf("invalid field index: %s", fieldStr)
		}
		fieldIdx = fieldIdx*10 + int(c-'0')
	}
	if fieldIdx < 1 {
		return "", fmt.Errorf("invalid field index: %d", fieldIdx)
	}

	fieldValue := seg.Field(fieldIdx)
	if len(parts) == 2 {
		return fieldValue, nil
	}

	// Component index
	compStr := parts[2]
	compIdx := 0
	for _, c := range compStr {
		if c < '0' || c > '9' {
			return "", fmt.Errorf("invalid component index: %s", compStr)
		}
		compIdx = compIdx*10 + int(c-'0')
	}
	if compIdx < 1 {
		return "", fmt.Errorf("invalid component index: %d", compIdx)
	}

	compValue := seg.Component(fieldIdx, compIdx)
	if len(parts) == 3 {
		return compValue, nil
	}

	// Subcomponent index
	subStr := parts[3]
	subIdx := 0
	for _, c := range subStr {
		if c < '0' || c > '9' {
			return "", fmt.Errorf("invalid subcomponent index: %s", subStr)
		}
		subIdx = subIdx*10 + int(c-'0')
	}
	if subIdx < 1 {
		return "", fmt.Errorf("invalid subcomponent index: %d", subIdx)
	}

	return seg.SubComponent(fieldIdx, compIdx, subIdx), nil
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
