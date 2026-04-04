package hl7

import (
	"strings"
	"testing"
	"time"
)

func TestSchemaLookup(t *testing.T) {
	tests := []struct {
		name      string
		segment   string
		wantFound bool
	}{
		{"MSH definition", "MSH", true},
		{"PID definition", "PID", true},
		{"PV1 definition", "PV1", true},
		{"OBX definition", "OBX", true},
		{"Unknown segment", "XYZ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, found := LookupSegmentDefinition(tt.segment)
			if found != tt.wantFound {
				t.Errorf("LookupSegmentDefinition(%q) = %v, want %v", tt.segment, found, tt.wantFound)
			}
		})
	}
}

func TestMessageStructureLookup(t *testing.T) {
	tests := []struct {
		name      string
		msgType   string
		wantFound bool
	}{
		{"ADT_A01", "ADT_A01", true},
		{"ORU_R01", "ORU_R01", true},
		{"By message type", "ADT^A01", true},
		{"Unknown", "UNKNOWN", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, found := LookupMessageStructure(tt.msgType)
			if found != tt.wantFound {
				t.Errorf("LookupMessageStructure(%q) = %v, want %v", tt.msgType, found, tt.wantFound)
			}
		})
	}
}

func TestHL7TableLookup(t *testing.T) {
	tests := []struct {
		name      string
		tableID   string
		code      string
		wantDesc  string
		wantFound bool
	}{
		{"HL70001 M", "HL70001", "M", "Male", true},
		{"HL70001 F", "HL70001", "F", "Female", true},
		{"HL70001 X", "HL70001", "X", "", false},
		{"Unknown table", "HL99999", "A", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			desc, found := LookupTable(tt.tableID, tt.code)
			if found != tt.wantFound {
				t.Errorf("LookupTable(%q, %q) found = %v, want %v", tt.tableID, tt.code, found, tt.wantFound)
			}
			if found && desc != tt.wantDesc {
				t.Errorf("LookupTable(%q, %q) desc = %q, want %q", tt.tableID, tt.code, desc, tt.wantDesc)
			}
		})
	}
}

func TestFieldDefinition(t *testing.T) {
	field, found := GetFieldDefinition("PID", 3)
	if !found {
		t.Fatal("GetFieldDefinition(PID, 3) not found")
	}
	if field.Name != "PatientIdentifierList" {
		t.Errorf("field.Name = %q, want %q", field.Name, "PatientIdentifierList")
	}
	if field.DataType != "CX" {
		t.Errorf("field.DataType = %q, want %q", field.DataType, "CX")
	}
	if !field.IsRequired {
		t.Error("field.IsRequired should be true")
	}
}

func TestTableValidation(t *testing.T) {
	data := []byte(`MSH|^~\&|SENDING|FACILITY|||202401151200||ADT^A01|MSG001|P|2.5
PID|1||12345^^^MRN||Smith^John^A||19800115|M|||`)

	msg, err := Parse(data)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Valid gender", func(t *testing.T) {
		v := NewValidator(RequiredTable("PID.8", "HL70001"))
		if errs := v.Validate(msg); len(errs) > 0 {
			t.Errorf("unexpected errors: %v", errs)
		}
	})

	t.Run("Invalid gender", func(t *testing.T) {
		v := NewValidator(RequiredTable("PID.8", "HL70001"))
		msg.PID().Gender()
		if errs := v.Validate(msg); len(errs) > 0 {
			if errs[0].Message == "" {
				t.Error("expected error message")
			}
		}
	})
}

func TestVersionRule(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		expected bool
	}{
		{"Valid v2.5", "2.5", true},
		{"Valid v2.4", "2.4", true},
		{"Invalid version", "9.9", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := []byte(`MSH|^~\&|SENDING|FACILITY|||202401151200||ADT^A01|MSG001|P|` + tt.version)
			msg, _ := Parse(data)
			v := NewValidator(SupportedVersion("2.4", "2.5"))
			errs := v.Validate(msg)
			hasError := len(errs) > 0
			if hasError != !tt.expected {
				t.Errorf("version %q: hasError=%v, expected=%v", tt.version, hasError, tt.expected)
			}
		})
	}
}

func TestMessageTypeRule(t *testing.T) {
	data := []byte(`MSH|^~\&|SENDING|FACILITY|||202401151200||ADT^A01|MSG001|P|2.5`)
	msg, _ := Parse(data)

	t.Run("Valid message type", func(t *testing.T) {
		v := NewValidator(SupportedMessageType("ADT^A01", "ADT^A08"))
		if errs := v.Validate(msg); len(errs) > 0 {
			t.Errorf("unexpected errors: %v", errs)
		}
	})

	t.Run("Invalid message type", func(t *testing.T) {
		v := NewValidator(SupportedMessageType("ORU^R01"))
		if errs := v.Validate(msg); len(errs) == 0 {
			t.Error("expected validation error")
		}
	})
}

func TestCompositeRules(t *testing.T) {
	data := []byte(`MSH|^~\&|SENDING|FACILITY|||202401151200||ADT^A01|MSG001|P|2.5
PID|1||12345^^^MRN||Smith^John^A||19800115|M|||`)
	msg, _ := Parse(data)

	t.Run("AllOf - all pass", func(t *testing.T) {
		v := NewValidator(AllOf(
			Required("MSH.9"),
			Required("PID.3.1"),
		))
		if errs := v.Validate(msg); len(errs) > 0 {
			t.Errorf("unexpected errors: %v", errs)
		}
	})

	t.Run("AnyOf - one passes", func(t *testing.T) {
		v := NewValidator(AnyOf(
			Required("PID.100"),
			Required("PID.3"),
		))
		if errs := v.Validate(msg); len(errs) > 0 {
			t.Errorf("unexpected errors: %v", errs)
		}
	})

	t.Run("Not rule", func(t *testing.T) {
		v := NewValidator(Not(Required("PID.100")))
		if errs := v.Validate(msg); len(errs) > 0 {
			t.Errorf("unexpected errors: %v", errs)
		}
	})

	t.Run("When rule - condition true", func(t *testing.T) {
		v := NewValidator(When(
			func(m *Message) bool { return m.HasSegment("PID") },
			Required("PID.3.1"),
		))
		if errs := v.Validate(msg); len(errs) > 0 {
			t.Errorf("unexpected errors: %v", errs)
		}
	})
}

func TestRangeRule(t *testing.T) {
	data := []byte(`MSH|^~\&|SENDING|FACILITY|||202401151200||ADT^A01|MSG001|P|2.5
OBX|1|NM|WBC^WBC||5.0|10^3/uL|4.5-11.0|||N|||`)
	msg, _ := Parse(data)

	t.Run("Value in range", func(t *testing.T) {
		v := NewValidator(Range("OBX.5", 1.0, 10.0))
		if errs := v.Validate(msg); len(errs) > 0 {
			t.Errorf("unexpected errors: %v", errs)
		}
	})

	t.Run("Value out of range", func(t *testing.T) {
		v := NewValidator(Range("OBX.5", 1.0, 3.0))
		if errs := v.Validate(msg); len(errs) == 0 {
			t.Error("expected validation error")
		}
	})
}

func TestDateFormatRule(t *testing.T) {
	data := []byte(`MSH|^~\&|SENDING|FACILITY|||202401151200||ADT^A01|MSG001|P|2.5
PID|1||12345||Smith^John||19800115|||`)
	msg, _ := Parse(data)

	t.Run("Valid format", func(t *testing.T) {
		v := NewValidator(DateFormat("PID.7", "20060102"))
		if errs := v.Validate(msg); len(errs) > 0 {
			t.Errorf("unexpected errors: %v", errs)
		}
	})

	t.Run("Invalid format", func(t *testing.T) {
		v := NewValidator(DateFormat("PID.7", "2006-01-02"))
		if errs := v.Validate(msg); len(errs) == 0 {
			t.Error("expected validation error")
		}
	})
}

func TestValidateWithSchema(t *testing.T) {
	data := []byte(`MSH|^~\&|SENDING|FACILITY|||202401151200||ADT^A01|MSG001|P|2.5
EVN|A01|202401151200|||
PID|1||12345^^^MRN||Smith^John^A||19800115|M|||
PV1|1|I|||Bed 1^Room 100^^Facility 1|||Dr. Johnson^John^M^MD|||`)
	msg, _ := Parse(data)

	errors := ValidateWithSchema(msg, "ADT_A01")
	if len(errors) > 0 {
		t.Errorf("unexpected errors: %v", errors)
	}
}

func TestLocationParse(t *testing.T) {
	tests := []struct {
		name      string
		loc       string
		segment   string
		field     int
		comp      int
		wantError bool
	}{
		{"Simple field", "PID.3", "PID", 3, 0, false},
		{"With component", "PID.3.1", "PID", 3, 1, false},
		{"Segment index", "PID[1].3", "PID", 3, 0, false},
		{"Repetition", "PID.3[0].1", "PID", 3, 1, false},
		{"Full path", "OBX[2].5[1].1.2", "OBX", 5, 1, false},
		{"Empty", "", "", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loc, err := ParseLocation(tt.loc)
			if tt.wantError {
				if err == nil {
					t.Error("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if loc.Segment != tt.segment {
				t.Errorf("Segment = %q, want %q", loc.Segment, tt.segment)
			}
			if loc.Field != tt.field {
				t.Errorf("Field = %d, want %d", loc.Field, tt.field)
			}
			if loc.Component != tt.comp {
				t.Errorf("Component = %d, want %d", loc.Component, tt.comp)
			}
		})
	}
}

func TestMessageGetRepetitions(t *testing.T) {
	data := []byte(`MSH|^~\&|SENDING|FACILITY|||202401151200||ADT^A01|MSG001|P|2.5
PID|1||12345~67890~11111^^^MRN||Smith^John|||`)
	msg, _ := Parse(data)

	reps, err := msg.GetAllRepetitions("PID.3")
	if err != nil {
		t.Fatal(err)
	}
	if len(reps) != 3 {
		t.Errorf("len(reps) = %d, want 3", len(reps))
	}
	if reps[0] != "12345" {
		t.Errorf("reps[0] = %q, want %q", reps[0], "12345")
	}
}

func TestSegmentRepetitions(t *testing.T) {
	data := []byte(`MSH|^~\&|SENDING|FACILITY|||202401151200||ADT^A01|MSG001|P|2.5
PID|1||12345~67890~11111^^^MRN||Smith^John|||`)
	msg, _ := Parse(data)
	pid, _ := msg.Segment("PID")

	reps := pid.Repetitions(3)
	if len(reps) != 3 {
		t.Errorf("len(reps) = %d, want 3", len(reps))
	}

	count := pid.RepetitionCount(3)
	if count != 3 {
		t.Errorf("RepetitionCount = %d, want 3", count)
	}
}

func TestGetNthSegment(t *testing.T) {
	data := []byte(`MSH|^~\&|SENDING|FACILITY|||202401151200||ADT^A01|MSG001|P|2.5
PID|1||12345||Smith|||`)
	msg, _ := Parse(data)

	seg, ok := msg.GetNthSegment("PID", 1)
	if !ok {
		t.Fatal("GetNthSegment failed")
	}
	if seg.Name() != "PID" {
		t.Errorf("seg.Name() = %q, want PID", seg.Name())
	}

	_, ok = msg.GetNthSegment("PID", 10)
	if ok {
		t.Error("expected not found for index 10")
	}
}

func TestHasSegment(t *testing.T) {
	data := []byte(`MSH|^~\&|SENDING|FACILITY|||202401151200||ADT^A01|MSG001|P|2.5
PID|1||12345||Smith|||`)
	msg, _ := Parse(data)

	if !msg.HasSegment("MSH") {
		t.Error("expected HasSegment(MSH) = true")
	}
	if !msg.HasSegment("PID") {
		t.Error("expected HasSegment(PID) = true")
	}
	if msg.HasSegment("OBX") {
		t.Error("expected HasSegment(OBX) = false")
	}
}

func TestEscape(t *testing.T) {
	tests := []struct {
		input    string
		contains string
	}{
		{"Hello|World", "\\S"},
		{"Test^Data", "^"},
		{"New\nLine", "\\R"},
		{"Tab\tHere", "\\T"},
		{"Back\\Slash", "\\\\"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := Escape(tt.input)
			if !strings.Contains(result, tt.contains) {
				t.Errorf("Escape(%q) = %q, should contain %q", tt.input, result, tt.contains)
			}
		})
	}
}

func TestUnescape(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		wantErr  bool
	}{
		{"Hello\\SWorld", "Hello|World", false},
		{"Test\\X0AData", "Test\nData", false},
		{"Invalid\\", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := Unescape(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Unescape(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestEncodingConversion(t *testing.T) {
	t.Run("ASCII to UTF-8", func(t *testing.T) {
		data := []byte("Hello World")
		result, err := ConvertToUTF8(data, EncodingASCII)
		if err != nil {
			t.Fatal(err)
		}
		if string(result) != "Hello World" {
			t.Errorf("conversion failed: %q", result)
		}
	})

	t.Run("UTF-8 validation", func(t *testing.T) {
		data := []byte("Hello")
		_, err := ConvertToUTF8(data, EncodingUTF8)
		if err != nil {
			t.Error("valid UTF-8 should not error")
		}
	})
}

func TestDecoder(t *testing.T) {
	decoder := NewDecoder()

	t.Run("Decode with escape", func(t *testing.T) {
		data := []byte("Hello\\S\\World")
		result, err := decoder.Decode(data)
		if err != nil {
			t.Fatal(err)
		}
		if string(result) != "Hello|World" {
			t.Errorf("got %q, want %q", result, "Hello|World")
		}
	})

	t.Run("With encoding", func(t *testing.T) {
		decoder := NewDecoder().WithEncoding(EncodingUTF8)
		data := []byte("Test")
		_, err := decoder.Decode(data)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestOptionalTag(t *testing.T) {
	type TestStruct struct {
		Required string `hl7:"PID.3.1"`
		Optional string `hl7:"PID.4.1,optional"`
	}

	data := []byte(`MSH|^~\&|SENDING|FACILITY|||202401151200||ADT^A01|MSG001|P|2.5
PID|1||12345|||`)
	msg, _ := Parse(data)

	var s TestStruct
	err := UnmarshalMessage(msg, &s)
	if err != nil {
		t.Fatal(err)
	}
	if s.Required != "12345" {
		t.Errorf("Required = %q, want %q", s.Required, "12345")
	}
	if s.Optional != "" {
		t.Errorf("Optional = %q, want empty", s.Optional)
	}
}

func TestNestedStructMapping(t *testing.T) {
	type Name struct {
		Family string `hl7:"1"`
		Given  string `hl7:"2"`
		Prefix string `hl7:"5"`
	}

	type Patient struct {
		Name Name `hl7:"PID.5"`
	}

	data := []byte(`MSH|^~\&|SENDING|FACILITY|||202401151200||ADT^A01|MSG001|P|2.5
PID|1||12345||Smith^John^Middle^Jr^Dr|||`)
	msg, _ := Parse(data)

	var p Patient
	err := UnmarshalMessage(msg, &p)
	if err != nil {
		t.Fatal(err)
	}
	if p.Name.Family != "Smith" {
		t.Errorf("Name.Family = %q, want %q", p.Name.Family, "Smith")
	}
	if p.Name.Given != "John" {
		t.Errorf("Name.Given = %q, want %q", p.Name.Given, "John")
	}
	if p.Name.Prefix != "Dr" {
		t.Errorf("Name.Prefix = %q, want %q", p.Name.Prefix, "Dr")
	}
}

func TestTimezone(t *testing.T) {
	now := time.Now()
	ts := NewTimestamp(now)
	hl7Str := ts.String()
	if hl7Str == "" {
		t.Error("expected non-empty timestamp string")
	}

	var ts2 Timestamp
	err := ts2.UnmarshalHL7([]byte("20240115120000"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestNMType(t *testing.T) {
	var nm NM
	err := nm.UnmarshalHL7([]byte("123.45"))
	if err != nil {
		t.Fatal(err)
	}
	if nm.Value != 123.45 {
		t.Errorf("nm.Value = %f, want 123.45", nm.Value)
	}

	data, err := nm.MarshalHL7()
	if err != nil {
		t.Fatal(err)
	}
	if string(data) == "" {
		t.Error("expected non-empty marshal output")
	}
}

func TestCETYPE(t *testing.T) {
	var ce CE
	err := ce.UnmarshalHL7([]byte("ICD10^Pneumonia^I10^J18.9^Pneumonia unspecified^I10"))
	if err != nil {
		t.Fatal(err)
	}
	if ce.Identifier != "ICD10" {
		t.Errorf("ce.Identifier = %q, want ICD10", ce.Identifier)
	}
	if ce.Text != "Pneumonia" {
		t.Errorf("ce.Text = %q, want Pneumonia", ce.Text)
	}
}

func TestCountSegment(t *testing.T) {
	data := []byte(`MSH|^~\&|SENDING|FACILITY|||202401151200||ADT^A01|MSG001|P|2.5
PID|1||12345||Smith|||
PID|1||67890||Doe|||`)
	msg, _ := Parse(data)

	count := msg.CountSegment("PID")
	if count != 2 {
		t.Errorf("CountSegment(PID) = %d, want 2", count)
	}

	count = msg.CountSegment("OBX")
	if count != 0 {
		t.Errorf("CountSegment(OBX) = %d, want 0", count)
	}
}

func TestClientPool(t *testing.T) {
	pool := NewClientPool("localhost:9999")
	if pool == nil {
		t.Fatal("NewClientPool returned nil")
	}
	pool.Close()
}

func TestClientWithRetry(t *testing.T) {
	client := &Client{
		retryCount:       3,
		retryDelay:       100 * time.Millisecond,
		maxRetryDelay:    500 * time.Millisecond,
		backoffMultipler: 2.0,
	}

	if client.retryCount != 3 {
		t.Errorf("retryCount = %d, want 3", client.retryCount)
	}
}
