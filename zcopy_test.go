package hl7

import (
	"bytes"
	"testing"
)

func TestBytesToString(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want string
	}{
		{"empty", []byte{}, ""},
		{"nil", nil, ""},
		{"simple", []byte("hello"), "hello"},
		{"with spaces", []byte("hello world"), "hello world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := bytesToString(tt.data)
			if got != tt.want {
				t.Errorf("bytesToString() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestStringToBytes(t *testing.T) {
	tests := []struct {
		name string
		s    string
	}{
		{"empty", ""},
		{"simple", "hello"},
		{"with spaces", "hello world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stringToBytes(tt.s)
			if !bytes.Equal(got, []byte(tt.s)) {
				t.Errorf("stringToBytes() = %v, want %v", got, []byte(tt.s))
			}
		})
	}
}

func TestBytesBuffer(t *testing.T) {
	buf := newBytesBuffer(64)

	buf.WriteString("hello")
	buf.WriteByte(' ')
	buf.WriteString("world")

	if got := buf.String(); got != "hello world" {
		t.Errorf("String() = %q, want %q", got, "hello world")
	}

	if !bytes.Equal(buf.Bytes(), []byte("hello world")) {
		t.Errorf("Bytes() = %v, want %v", buf.Bytes(), []byte("hello world"))
	}

	buf.Reset()
	if buf.String() != "" {
		t.Errorf("After Reset, String() = %q, want empty", buf.String())
	}
}

func TestZeroCopyParser(t *testing.T) {
	parser := newZeroCopyParser()

	tests := []struct {
		name    string
		input   string
		wantSeg int
		wantErr bool
	}{
		{
			name:    "simple message",
			input:   "MSH|^~\\&|SEND|RECV|||20240115||ADT^A01|1|P|2.5\rPID|1||12345||DOE^JOHN",
			wantSeg: 2,
			wantErr: false,
		},
		{
			name:    "empty",
			input:   "",
			wantSeg: 0,
			wantErr: true,
		},
		{
			name:    "MSH only",
			input:   "MSH|^~\\&|SEND|RECV",
			wantSeg: 1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := parser.Parse([]byte(tt.input))
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(msg.segments) != tt.wantSeg {
				t.Errorf("got %d segments, want %d", len(msg.segments), tt.wantSeg)
			}
			PutMessage(msg)
		})
	}
}

func TestSplitBytes(t *testing.T) {
	tests := []struct {
		name   string
		data   []byte
		sep    byte
		expect []string
	}{
		{"simple", []byte("a|b|c"), '|', []string{"a", "b", "c"}},
		{"empty", []byte(""), '|', nil},
		{"no sep", []byte("abc"), '|', []string{"abc"}},
		{"trailing sep", []byte("a|b|"), '|', []string{"a", "b", ""}},
		{"leading sep", []byte("|a|b"), '|', []string{"", "a", "b"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitBytes(tt.data, tt.sep)
			if len(got) != len(tt.expect) {
				t.Errorf("splitBytes() count = %d, want %d", len(got), len(tt.expect))
				return
			}
			for i, s := range got {
				if string(s) != tt.expect[i] {
					t.Errorf("splitBytes()[%d] = %q, want %q", i, string(s), tt.expect[i])
				}
			}
		})
	}
}

func TestNormalizeLineEndings(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  []byte
	}{
		{"windows", []byte("line1\r\nline2\r\nline3"), []byte("line1\rline2\rline3")},
		{"unix", []byte("line1\nline2\nline3"), []byte("line1\rline2\rline3")},
		{"mac", []byte("line1\rline2\rline3"), []byte("line1\rline2\rline3")},
		{"mixed", []byte("line1\r\nline2\nline3\rline4"), []byte("line1\rline2\rline3\rline4")},
		{"trailing crlf", []byte("line1\r\n"), []byte("line1")},
		{"trailing lf", []byte("line1\n"), []byte("line1")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeLineEndings(tt.input)
			if !bytes.Equal(got, tt.want) {
				t.Errorf("NormalizeLineEndings() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestBuildSegmentString(t *testing.T) {
	delims := DefaultDelimiters()

	tests := []struct {
		name    string
		segName string
		fields  []string
		want    string
	}{
		{"MSH", "MSH", []string{"^~\\&", "SEND", "RECV"}, "MSH|^~\\&|SEND|RECV"},
		{"PID", "PID", []string{"1", "12345", "DOE^JOHN"}, "PID|1|12345|DOE^JOHN"},
		{"empty fields", "PID", []string{}, "PID"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildSegmentString(tt.segName, tt.fields, delims)
			if got != tt.want {
				t.Errorf("BuildSegmentString() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBuildSegmentBytes(t *testing.T) {
	delims := DefaultDelimiters()

	seg := BuildSegmentBytes("MSH", []string{"^~\\&", "SEND"}, delims)
	want := []byte("MSH|^~\\&|SEND")

	if !bytes.Equal(seg, want) {
		t.Errorf("BuildSegmentBytes() = %v, want %v", string(seg), string(want))
	}
}

func BenchmarkZeroCopyParser(b *testing.B) {
	parser := newZeroCopyParser()
	data := []byte("MSH|^~\\&|SENDING|FACILITY|||202401151200||ADT^A01|MSG001|P|2.5\r" +
		"EVN|A01|202401151200|||\r" +
		"PID|1||12345^^^MRN^MR^N||Smith^John^A||19800115|M|||123 Main St^^New York^NY^10001||(555)123-4567|||\r" +
		"PV1|1|I|||Bed 1^Room 100^^Facility 1|||Dr. Johnson^John^M^MD|||ED|||||||||||")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		msg, _ := parser.Parse(data)
		PutMessage(msg)
	}
}

func BenchmarkStandardParser(b *testing.B) {
	parser := NewParser()
	data := []byte("MSH|^~\\&|SENDING|FACILITY|||202401151200||ADT^A01|MSG001|P|2.5\r" +
		"EVN|A01|202401151200|||\r" +
		"PID|1||12345^^^MRN^MR^N||Smith^John^A||19800115|M|||123 Main St^^New York^NY^10001||(555)123-4567|||\r" +
		"PV1|1|I|||Bed 1^Room 100^^Facility 1|||Dr. Johnson^John^M^MD|||ED|||||||||||")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		parser.Parse(data)
	}
}

func BenchmarkBytesToString(b *testing.B) {
	data := []byte("MSH|^~\\&|SENDING|FACILITY|||202401151200||ADT^A01|MSG001|P|2.5")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = bytesToString(data)
	}
}

func BenchmarkStringConversion(b *testing.B) {
	data := []byte("MSH|^~\\&|SENDING|FACILITY|||202401151200||ADT^A01|MSG001|P|2.5")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = string(data)
	}
}
