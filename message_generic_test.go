package hl7

import (
	"testing"
)

func TestMessageGet(t *testing.T) {
	data := []byte("MSH|^~\\&|APP|FAC|||20240101120000||ADT^A01|MSG001|P|2.5\r" +
		"PID|1||12345^^^MRN||Smith^John||19900101|M\r" +
		"PV1|1|I|ICU^Room101^Bed1")

	msg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	tests := []struct {
		location string
		want     string
		wantErr  bool
	}{
		{"MSH.3", "APP", false},
		{"MSH.4", "FAC", false},
		{"MSH.9", "ADT^A01", false},
		{"MSH.10", "MSG001", false},
		{"MSH.12", "2.5", false},
		{"MSH.1", "|", false},
		{"PID.3", "12345^^^MRN", false},
		{"PID.5", "Smith^John", false},
		{"PID.5.1", "Smith", false},
		{"PID.5.2", "John", false},
		{"PID.7", "19900101", false},
		{"PID.8", "M", false},
		{"PV1.3", "ICU^Room101^Bed1", false},
		{"PV1.3.1", "ICU", false},
		{"PV1.3.2", "Room101", false},
		{"PV1.3.3", "Bed1", false},
		{"XXX.1", "", true},
		{"MSH.100", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.location, func(t *testing.T) {
			got, err := msg.Get(tt.location)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get(%q) error = %v, wantErr %v", tt.location, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Get(%q) = %q, want %q", tt.location, got, tt.want)
			}
		})
	}
}

func TestMessageGetRepeated(t *testing.T) {
	data := []byte("MSH|^~\\&|APP|FAC|||20240101120000||ADT^A01|MSG001|P|2.5\r" +
		"OBX|1|NM|WBC||7.5|10*3/uL||||F|\r" +
		"OBX|2|NM|RBC||4.8|10*6/uL||||F|\r" +
		"OBX|3|NM|HGB||14.2|g/dL||||F|")

	msg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	tests := []struct {
		location string
		want     string
		wantErr  bool
	}{
		{"OBX[1].3", "WBC", false},
		{"OBX[2].3", "RBC", false},
		{"OBX[3].3", "HGB", false},
		{"OBX[1].5", "7.5", false},
		{"OBX[2].5", "4.8", false},
		{"OBX[3].5", "14.2", false},
	}

	for _, tt := range tests {
		t.Run(tt.location, func(t *testing.T) {
			got, err := msg.Get(tt.location)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get(%q) error = %v, wantErr %v", tt.location, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Get(%q) = %q, want %q", tt.location, got, tt.want)
			}
		})
	}
}

func TestMessageMustGet(t *testing.T) {
	data := []byte("MSH|^~\\&|APP|FAC|||20240101120000||ADT^A01|MSG001|P|2.5")

	msg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	val := msg.MustGet("MSH.3")
	if val != "APP" {
		t.Errorf("MustGet(MSH.3) = %q, want %q", val, "APP")
	}
}

func TestMessageIterate(t *testing.T) {
	data := []byte("MSH|^~\\&|APP|FAC|||20240101120000||ADT^A01|MSG001|P|2.5\r" +
		"PID|1||12345||Smith||19900101|M\r" +
		"PV1|1|I|ICU")

	msg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	var names []string
	for iter := range msg.Iterate() {
		names = append(names, iter.Name)
	}

	if len(names) != 3 {
		t.Errorf("Iterate() returned %d segments, want 3", len(names))
	}

	if names[0] != "MSH" || names[1] != "PID" || names[2] != "PV1" {
		t.Errorf("Segment names = %v, want [MSH PID PV1]", names)
	}
}

func TestSegmentIterator(t *testing.T) {
	data := []byte("MSH|^~\\&|APP|FAC|||20240101120000||ADT^A01|MSG001|P|2.5")

	msg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	for iter := range msg.Iterate() {
		if iter.Name == "MSH" {
			if iter.Count() != 12 {
				t.Errorf("MSH Count() = %d, want 12", iter.Count())
			}

			if iter.ValueAt(3) != "APP" {
				t.Errorf("MSH ValueAt(3) = %q, want %q", iter.ValueAt(3), "APP")
			}

			if iter.ValueAt(4) != "FAC" {
				t.Errorf("MSH ValueAt(4) = %q, want %q", iter.ValueAt(4), "FAC")
			}
		}
	}
}

func TestMessageStats(t *testing.T) {
	data := []byte("MSH|^~\\&|APP|FAC|||20240101120000||ADT^A01|MSG001|P|2.5\r" +
		"PID|1||12345||Smith||19900101|M\r" +
		"PV1|1|I|ICU")

	msg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	stats := msg.Stats()

	if stats.SegmentCount != 3 {
		t.Errorf("SegmentCount = %d, want 3", stats.SegmentCount)
	}

	if !stats.HasMSH {
		t.Error("HasMSH should be true")
	}

	if !stats.HasPID {
		t.Error("HasPID should be true")
	}

	if !stats.HasPV1 {
		t.Error("HasPV1 should be true")
	}

	if stats.MessageType != "ADT^A01" {
		t.Errorf("MessageType = %q, want %q", stats.MessageType, "ADT^A01")
	}

	if stats.Version != "2.5" {
		t.Errorf("Version = %q, want %q", stats.Version, "2.5")
	}

	if stats.SegmentTypes["MSH"] != 1 {
		t.Errorf("SegmentTypes[MSH] = %d, want 1", stats.SegmentTypes["MSH"])
	}
}

func TestMessageSummary(t *testing.T) {
	data := []byte("MSH|^~\\&|APP|FAC|||20240101120000||ADT^A01|MSG001|P|2.5\r" +
		"PID|1||12345||Smith||19900101|M")

	msg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	summary := msg.Summary()

	if summary == "" {
		t.Error("Summary() returned empty string")
	}
}
