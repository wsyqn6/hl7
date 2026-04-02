package hl7

import (
	"testing"
)

func TestParseSimpleMessage(t *testing.T) {
	data := []byte("MSH|^~\\&|SENDING|FACILITY|||202401151200||ADT^A01|MSG001|P|2.5\rPID|1||12345^^^MRN||Smith^John^A||19800115|M")

	msg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if msg == nil {
		t.Fatal("Parse() returned nil message")
	}

	// Check message type
	if msgType := msg.Type(); msgType != "ADT^A01" {
		t.Errorf("Type() = %q, want %q", msgType, "ADT^A01")
	}

	// Check control ID
	if ctrlID := msg.ControlID(); ctrlID != "MSG001" {
		t.Errorf("ControlID() = %q, want %q", ctrlID, "MSG001")
	}

	// Check segments
	if segs := msg.AllSegments(); len(segs) != 2 {
		t.Errorf("expected 2 segments, got %d", len(segs))
	}
}

func TestParseEmptyData(t *testing.T) {
	_, err := Parse([]byte{})
	if err == nil {
		t.Error("Parse() expected error for empty data")
	}
}

func TestParseString(t *testing.T) {
	data := "MSH|^~\\&|APP|FAC||||||MSG001|P|2.5"
	msg, err := ParseString(data)
	if err != nil {
		t.Fatalf("ParseString() error = %v", err)
	}
	if msg == nil {
		t.Fatal("ParseString() returned nil")
	}
}

func TestParserWithCustomDelimiters(t *testing.T) {
	delims := Delimiters{
		Field:        '#',
		Component:    '*',
		Repetition:   '~',
		Escape:       '\\',
		SubComponent: '&',
	}
	parser := NewParserWithDelimiters(delims)
	data := []byte("MSH#^~\\&#APP#FAC#||MSG001#P#2.5")

	msg, err := parser.Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if msg == nil {
		t.Fatal("Parse() returned nil")
	}
}

func TestSplitMessages(t *testing.T) {
	data := []byte("MSH|1\rPID|1\r\rMSH|2\rPID|2")
	msgs := SplitMessages(data)
	if len(msgs) != 2 {
		t.Errorf("expected 2 messages, got %d", len(msgs))
	}
}

func TestIsHL7Message(t *testing.T) {
	tests := []struct {
		data     []byte
		expected bool
	}{
		{[]byte("MSH|^~\\&|APP|FAC||||||MSG001|P|2.5"), true},
		{[]byte{0x0B, 'M', 'S', 'H', '|'}, true},
		{[]byte("NOTHL7"), false},
		{[]byte{}, false},
		{[]byte("MSH|"), true},
	}

	for _, tt := range tests {
		result := IsHL7Message(tt.data)
		if result != tt.expected {
			t.Errorf("IsHL7Message(%q) = %v, want %v", tt.data, result, tt.expected)
		}
	}
}

func TestStripMLLP(t *testing.T) {
	data := []byte{0x0B, 'M', 'S', 'H', '|', 'T', 'E', 'S', 'T', 0x1C, 0x0D}
	result := StripMLLP(data)
	expected := []byte("MSH|TEST")
	if string(result) != string(expected) {
		t.Errorf("StripMLLP() = %q, want %q", result, expected)
	}
}
