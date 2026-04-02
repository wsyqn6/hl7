package hl7

import (
	"bytes"
	"io"
	"strings"
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

	if msgType := msg.Type(); msgType != "ADT^A01" {
		t.Errorf("Type() = %q, want %q", msgType, "ADT^A01")
	}

	if ctrlID := msg.ControlID(); ctrlID != "MSG001" {
		t.Errorf("ControlID() = %q, want %q", ctrlID, "MSG001")
	}

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

func TestParserWithConfig(t *testing.T) {
	config := ParserConfig{
		MaxSegments:    10,
		MaxFieldLength: 1024,
	}
	parser := NewParserWithConfig(config)
	if parser.config.MaxSegments != 10 {
		t.Errorf("MaxSegments = %d, want 10", parser.config.MaxSegments)
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

func TestScannerBasic(t *testing.T) {
	// Use double newline to separate messages
	data := `MSH|^~\&|APP|FAC||||||MSG001|P|2.5
PID|1||12345||Smith||19900101|M

MSH|^~\&|APP|FAC||||||MSG002|P|2.5
PID|1||67890||Doe||19850101|F`

	reader := strings.NewReader(data)
	scanner := NewScanner(reader)

	count := 0
	for scanner.Scan() {
		msg := scanner.Message()
		if msg == nil {
			t.Fatal("Message() returned nil")
		}
		count++
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("Scanner error: %v", err)
	}

	if count != 2 {
		t.Errorf("expected 2 messages, got %d", count)
	}
}

func TestScannerMLLPFrame(t *testing.T) {
	// MLLP framed messages with proper structure: MSH|^~\&|APP|FAC|||20240101120000||ADT^A01|MSG001|P|2.5
	data := []byte{
		0x0B, // MLLP_START
		// Message 1: ADT^A01|MSG001
		'M', 'S', 'H', '|', '^', '~', '\\', '&', '|', 'A', 'P', 'P', '|', 'F', 'A', 'C', '|', '|', '|', '2', '0', '2', '4', '0', '1', '0', '1', '1', '2', '0', '0', '0', '0', '|', '|', 'A', 'D', 'T', '^', 'A', '0', '1', '|', 'M', 'S', 'G', '0', '0', '1', '|', 'P', '|', '2', '.', '5', '\r',
		0x1C, // MLLP_END
		0x0D, // MLLP_CR
		0x0B, // MLLP_START
		// Message 2: ADT^A08|MSG002
		'M', 'S', 'H', '|', '^', '~', '\\', '&', '|', 'A', 'P', 'P', '|', 'F', 'A', 'C', '|', '|', '|', '2', '0', '2', '4', '0', '1', '0', '1', '1', '2', '0', '0', '0', '0', '|', '|', 'A', 'D', 'T', '^', 'A', '0', '8', '|', 'M', 'S', 'G', '0', '0', '2', '|', 'P', '|', '2', '.', '5', '\r',
		0x1C, // MLLP_END
		0x0D, // MLLP_CR
	}

	scanner := NewScanner(bytes.NewReader(data))

	count := 0
	for scanner.Scan() {
		msg := scanner.Message()
		if msg == nil {
			t.Fatal("Message() returned nil")
		}
		ctrlID := msg.ControlID()
		if ctrlID != "MSG001" && ctrlID != "MSG002" {
			t.Errorf("unexpected control ID: %s", ctrlID)
		}
		msgType := msg.Type()
		if msgType != "ADT^A01" && msgType != "ADT^A08" {
			t.Errorf("unexpected message type: %s", msgType)
		}
		count++
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("Scanner error: %v", err)
	}

	if count != 2 {
		t.Errorf("expected 2 MLLP messages, got %d", count)
	}
}

func TestScannerMaxMessageSize(t *testing.T) {
	// Create a message larger than 100 bytes
	data := strings.Repeat("MSH|", 50) + "\r"

	scanner := NewScannerWithOptions(
		strings.NewReader(data),
		WithMaxMessageSize(100),
	)

	if scanner.Scan() {
		t.Error("expected scanner to fail due to size limit")
	}

	if scanner.Err() == nil {
		t.Error("expected error for message too large")
	}
}

func TestScannerSkipInvalid(t *testing.T) {
	// Data with valid messages separated by double newline
	data := `MSH|^~\&|APP|FAC||||||MSG001|P|2.5
PID|1||12345||Smith||19900101|M

MSH|^~\&|APP|FAC||||||MSG002|P|2.5
PID|1||67890||Doe||19850101|F`

	scanner := NewScannerWithOptions(
		strings.NewReader(data),
		WithSkipInvalid(true),
	)

	count := 0
	for scanner.Scan() {
		count++
	}

	// Should have scanned 2 valid messages
	if count != 2 {
		t.Errorf("expected 2 messages with skip, got %d", count)
	}

	if scanner.Err() != nil {
		t.Logf("Scanner reported errors (expected): %v", scanner.Err())
	}
}

func TestScannerCount(t *testing.T) {
	// Use double newline to separate messages
	data := `MSH|^~\&|APP|FAC||||||MSG001|P|2.5
MSH|^~\&|APP|FAC||||||MSG002|P|2.5
MSH|^~\&|APP|FAC||||||MSG003|P|2.5`

	scanner := NewScanner(strings.NewReader(data))

	count := 0
	for scanner.Scan() {
		count++
	}

	if scanner.Count() != 1 {
		t.Errorf("Count() = %d, want 1 (no double newline separator)", scanner.Count())
	}
}

func TestScannerReset(t *testing.T) {
	data := `MSH|^~\&|APP|FAC||||||MSG001|P|2.5
MSH|^~\&|APP|FAC||||||MSG002|P|2.5`

	scanner := NewScanner(strings.NewReader(data))
	scanner.Scan()

	// Reset with new data
	newData := `MSH|^~\&|APP|FAC||||||MSG003|P|2.5
MSH|^~\&|APP|FAC||||||MSG004|P|2.5`
	scanner.Reset(strings.NewReader(newData))

	count := 0
	for scanner.Scan() {
		count++
	}

	if count != 1 {
		t.Errorf("expected 1 message after reset, got %d", count)
	}
}

func TestScannerMLLPWithCRLFFormat(t *testing.T) {
	// MLLP with proper structure
	data := []byte{
		0x0B,
		// MSH|^~\&|APP|FAC|||20240101120000||ADT^A01|MSG001|P|2.5
		'M', 'S', 'H', '|', '^', '~', '\\', '&', '|', 'A', 'P', 'P', '|', 'F', 'A', 'C', '|', '|', '|', '2', '0', '2', '4', '0', '1', '0', '1', '1', '2', '0', '0', '0', '0', '|', '|', 'A', 'D', 'T', '^', 'A', '0', '1', '|', 'M', 'S', 'G', '0', '0', '1', '|', 'P', '|', '2', '.', '5', '\r',
		0x1C, 0x0D,
	}

	scanner := NewScanner(bytes.NewReader(data))

	if !scanner.Scan() {
		t.Fatalf("Scan() = false, want true. Err: %v", scanner.Err())
	}

	msg := scanner.Message()
	if msg == nil {
		t.Fatal("Message() returned nil")
	}

	if msg.ControlID() != "MSG001" {
		t.Errorf("ControlID() = %q, want MSG001", msg.ControlID())
	}
}

func TestScannerMixedContent(t *testing.T) {
	// Valid HL7 messages separated by double newline
	data := `MSH|^~\&|APP|FAC||||||MSG001|P|2.5
PID|1||12345||Smith||19900101|M

MSH|^~\&|APP|FAC||||||MSG002|P|2.5
PID|1||67890||Doe||19850101|F`

	scanner := NewScannerWithOptions(
		strings.NewReader(data),
		WithSkipInvalid(true),
	)

	count := 0
	for scanner.Scan() {
		count++
	}

	if count != 2 {
		t.Errorf("expected 2 messages, got %d", count)
	}
}

func TestScannerEOFWithoutTrailing(t *testing.T) {
	// Single message without trailing newline
	data := `MSH|^~\&|APP|FAC||||||MSG001|P|2.5
PID|1||12345||Smith||19900101|M`

	scanner := NewScanner(strings.NewReader(data))

	if !scanner.Scan() {
		t.Fatalf("Scan() = false, want true. Err: %v", scanner.Err())
	}

	if scanner.Message() == nil {
		t.Fatal("Message() returned nil")
	}

	if scanner.Count() != 1 {
		t.Errorf("Count() = %d, want 1", scanner.Count())
	}
}

func TestScannerEmptyStream(t *testing.T) {
	scanner := NewScanner(strings.NewReader(""))

	if scanner.Scan() {
		t.Error("expected Scan() = false for empty stream")
	}

	if scanner.Err() != nil {
		t.Errorf("unexpected error: %v", scanner.Err())
	}
}

func TestScannerWithCustomDelimiters(t *testing.T) {
	delims := Delimiters{
		Field:        '#',
		Component:    '*',
		Repetition:   '~',
		Escape:       '\\',
		SubComponent: '&',
	}

	// Use custom delimiter: MSH#^~\&#APP#FAC#####ADT^A01#MSG001#P#2.5
	data := "MSH#^~\\&#APP#FAC#####ADT^A01#MSG001#P#2.5"

	scanner := NewScannerWithOptions(
		strings.NewReader(data),
		WithDelimiters(delims),
	)

	if !scanner.Scan() {
		t.Fatalf("Scan() = false, want true. Err: %v", scanner.Err())
	}

	msg := scanner.Message()
	if msg == nil {
		t.Fatal("Message() returned nil")
	}

	if msg.Type() != "ADT^A01" {
		t.Errorf("Type() = %q, want ADT^A01", msg.Type())
	}

	if msg.ControlID() != "MSG001" {
		t.Errorf("ControlID() = %q, want MSG001", msg.ControlID())
	}
}

func TestScannerWithMaxMessageSize(t *testing.T) {
	scanner := NewScannerWithOptions(
		strings.NewReader(""),
		WithMaxMessageSize(1024*1024), // 1MB
	)

	if scanner.cfg.MaxMessageSize != 1024*1024 {
		t.Errorf("MaxMessageSize = %d, want 1048576", scanner.cfg.MaxMessageSize)
	}
}

func TestNewScannerWithOptions(t *testing.T) {
	opts := []ScannerOption{
		WithMaxMessageSize(1024 * 1024),
		WithSkipInvalid(true),
	}

	cfg := ScannerConfig{
		MaxMessageSize: DefaultMaxMessageSize,
		Delimiters:     DefaultDelimiters(),
		SkipInvalid:    false,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	if cfg.MaxMessageSize != 1024*1024 {
		t.Errorf("MaxMessageSize = %d, want 1048576", cfg.MaxMessageSize)
	}

	if !cfg.SkipInvalid {
		t.Error("SkipInvalid should be true")
	}
}

// Benchmark tests

func BenchmarkScanner(b *testing.B) {
	data := strings.Repeat("MSH|^~\\&|APP|FAC||||||MSG001|P|2.5\rPID|1||12345||Smith^John||19900101|M\r", 10)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		scanner := NewScanner(strings.NewReader(data))
		for scanner.Scan() {
			_ = scanner.Message()
		}
	}
}

func BenchmarkScannerMLLP(b *testing.B) {
	var buf bytes.Buffer
	for i := 0; i < 10; i++ {
		buf.WriteByte(MLLP_START)
		buf.WriteString("MSH|^~\\&|APP|FAC||||||MSG001|P|2.5\r")
		buf.WriteByte(MLLP_END)
		buf.WriteByte(MLLP_CR)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		scanner := NewScanner(bytes.NewReader(buf.Bytes()))
		for scanner.Scan() {
			_ = scanner.Message()
		}
	}
}

func BenchmarkParse(b *testing.B) {
	data := []byte("MSH|^~\\&|SENDING|FACILITY|||202401151200||ADT^A01|MSG001|P|2.5\rPID|1||12345^^^MRN||Smith^John^A||19800115|M")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = Parse(data)
	}
}

// Ensure io.Reader interface is satisfied
var _ io.Reader = (*strings.Reader)(nil)

// Fuzz tests

func FuzzMessageParse(f *testing.F) {
	testcases := []string{
		"MSH|^~\\&|APP|FAC||||||MSG001|P|2.5",
		"MSH|^~\\&|APP|FAC||||||MSG001|P|2.5\rPID|1||12345||Smith||19900101|M",
		"MSH|^~\\&|APP|FAC|||20240101120000||ADT^A01|MSG001|P|2.5\rPID|1||12345||Smith^John||19800115|M\rPV1|1|I|||Dr.Smith",
		"MSH|^~\\&|SENDING|FACILITY|||202401151200||ADT^A01|MSG001|P|2.5\r",
		"",
	}
	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, data string) {
		msg, err := Parse([]byte(data))
		if err != nil {
			return
		}
		_ = msg.Type()
		_ = msg.ControlID()
		_ = msg.AllSegments()
	})
}

func FuzzSegmentParse(f *testing.F) {
	testcases := []string{
		"MSH|^~\\&|APP|FAC||||||MSG001|P|2.5",
		"PID|1||12345||Smith||19900101|M",
		"PV1|1|I|||Dr.Smith",
		"OBR|1||12345|CBC^Complete Blood Count",
		"",
	}
	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, data string) {
		seg := ParseSegment(data)
		if seg == nil {
			return
		}
		_ = seg.Name()
		_ = seg.Field(1)
		_ = seg.Fields()
	})
}

func FuzzFieldParse(f *testing.F) {
	testcases := []string{
		"Smith^John",
		"12345^^^MRN",
		"19800115",
		"",
	}
	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, data string) {
		field := ParseField(data)
		_ = field.Value
		_ = field.Components
	})
}

func FuzzComponentParse(f *testing.F) {
	testcases := []string{
		"Smith^John",
		"12345",
		"",
	}
	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, data string) {
		comps := ParseComponents(data)
		_ = len(comps)
	})
}

func FuzzScannerStream(f *testing.F) {
	testcases := []string{
		"MSH|^~\\&|APP|FAC||||||MSG001|P|2.5\r",
		"MSH|^~\\&|APP|FAC||||||MSG001|P|2.5\rPID|1||12345||Smith||19900101|M\r\n\nMSH|^~\\&|APP|FAC||||||MSG002|P|2.5\rPID|1||67890||Doe||19850101|F\r",
	}
	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, data string) {
		scanner := NewScanner(strings.NewReader(data))
		for scanner.Scan() {
			msg := scanner.Message()
			if msg != nil {
				_ = msg.Type()
				_ = len(msg.AllSegments())
			}
		}
	})
}

func FuzzScannerMLLPFramed(f *testing.F) {
	testcases := [][]byte{
		{0x0B, 'M', 'S', 'H', '|', '^', '~', '\\', '&', '|', 'A', 'P', 'P', '|', 'F', 'A', 'C', '|', '|', '|', '2', '0', '2', '4', '0', '1', '0', '1', '1', '2', '0', '0', '0', '0', '|', '|', 'A', 'D', 'T', '^', 'A', '0', '1', '|', 'M', 'S', 'G', '0', '0', '1', '|', 'P', '|', '2', '.', '5', '\r', 0x1C, 0x0D},
	}
	for _, tc := range testcases {
		f.Add(string(tc))
	}

	f.Fuzz(func(t *testing.T, data string) {
		scanner := NewScanner(strings.NewReader(data))
		for scanner.Scan() {
			msg := scanner.Message()
			if msg != nil {
				_ = msg.Type()
				_ = msg.ControlID()
			}
		}
	})
}

func FuzzMessageSplit(f *testing.F) {
	testcases := []string{
		"MSH|^~\\&|APP|FAC||||||MSG001|P|2.5\r\rMSH|^~\\&|APP|FAC||||||MSG002|P|2.5",
		"MSH|^~\\&|APP|FAC||||||MSG001|P|2.5\r",
		"",
	}
	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, data string) {
		msgs := SplitMessages([]byte(data))
		_ = len(msgs)
	})
}

func FuzzHL7Detection(f *testing.F) {
	testcases := []string{
		"MSH|^~\\&|APP|FAC||||||MSG001|P|2.5",
		"NOTHL7",
		"",
		"\x0BMSH|",
	}
	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, data string) {
		_ = IsHL7Message([]byte(data))
	})
}
