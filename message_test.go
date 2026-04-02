package hl7

import (
	"testing"
)

func TestNewMessage(t *testing.T) {
	msg := NewMessage()
	if msg == nil {
		t.Fatal("NewMessage returned nil")
	}
	if len(msg.AllSegments()) != 0 {
		t.Errorf("expected 0 segments, got %d", len(msg.AllSegments()))
	}
}

func TestDefaultDelimiters(t *testing.T) {
	delims := DefaultDelimiters()
	if delims.Field != '|' {
		t.Errorf("expected field separator '|', got %c", delims.Field)
	}
	if delims.Component != '^' {
		t.Errorf("expected component separator '^', got %c", delims.Component)
	}
}

func TestMessageAddSegment(t *testing.T) {
	msg := NewMessage()
	seg := NewSegment("MSH")
	seg.SetField(1, "|")
	seg.SetField(2, "^~\\&")
	msg.AddSegment(seg)

	if len(msg.AllSegments()) != 1 {
		t.Errorf("expected 1 segment, got %d", len(msg.AllSegments()))
	}
}

func TestSegmentField(t *testing.T) {
	seg := NewSegment("PID")
	seg.SetField(1, "1")
	seg.SetField(3, "12345")

	if seg.Field(1) != "1" {
		t.Errorf("expected field 1 to be '1', got %s", seg.Field(1))
	}
	if seg.Field(3) != "12345" {
		t.Errorf("expected field 3 to be '12345', got %s", seg.Field(3))
	}
	if seg.Field(99) != "" {
		t.Errorf("expected field 99 to be empty, got %s", seg.Field(99))
	}
}

func TestSegmentComponent(t *testing.T) {
	seg := NewSegment("PID")
	seg.SetField(5, "Smith^John^A")

	if comp := seg.Component(5, 1); comp != "Smith" {
		t.Errorf("expected component 1 to be 'Smith', got %s", comp)
	}
	if comp := seg.Component(5, 2); comp != "John" {
		t.Errorf("expected component 2 to be 'John', got %s", comp)
	}
	if comp := seg.Component(5, 3); comp != "A" {
		t.Errorf("expected component 3 to be 'A', got %s", comp)
	}
}

func TestSplitField(t *testing.T) {
	tests := []struct {
		input     string
		separator rune
		expected  []string
	}{
		{"A^B^C", '^', []string{"A", "B", "C"}},
		{"A|B|C", '|', []string{"A", "B", "C"}},
		{"single", '^', []string{"single"}},
		{"", '^', nil},
	}

	for _, tt := range tests {
		result := SplitField(tt.input, tt.separator)
		if len(result) != len(tt.expected) {
			t.Errorf("SplitField(%q, %c) returned %d fields, expected %d", tt.input, tt.separator, len(result), len(tt.expected))
			continue
		}
		for i, v := range result {
			if v != tt.expected[i] {
				t.Errorf("SplitField(%q, %c)[%d] = %q, expected %q", tt.input, tt.separator, i, v, tt.expected[i])
			}
		}
	}
}
