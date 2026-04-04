package hl7

import (
	"bytes"
	"strings"
	"testing"
)

var (
	adta01Message = `MSH|^~\&|SENDING|FACILITY|||202401151200||ADT^A01|MSG001|P|2.5
EVN|A01|202401151200|||
PID|1||12345^^^MRN^MR^N||Smith^John^A||19800115|M|||123 Main St^^New York^NY^10001||(555)123-4567|||
PV1|1|I|||Bed 1^Room 100^^Facility 1|||Dr. Johnson^John^M^MD|||ED|||||||||||`
	largeMessage string
)

func init() {
	var sb strings.Builder
	sb.WriteString("MSH|^~\\&|SENDING|FACILITY|||202401151200||ORU^R01|MSG001|P|2.5\r\n")
	sb.WriteString("PID|1||12345^^^MRN||Smith^John^A||19800115|M|||123 Main St^^New York^NY^10001||\r\n")
	sb.WriteString("PV1|1|O|||Dr. Johnson^John^M^MD|||||||||||||||||||||||||\r\n")

	for i := 0; i < 50; i++ {
		sb.WriteString("OBX|")
		sb.WriteString(string(rune('1' + i%9)))
		sb.WriteString("|NM|WBC^WBC^LN||")
		sb.WriteString(string(rune('0' + i%10)))
		sb.WriteString(".5|x10^3/uL^3^1^ML||4.5-11.0|N|||F|||202401151200\r\n")
	}
	largeMessage = sb.String()
}

func BenchmarkParseSimple(b *testing.B) {
	data := []byte(adta01Message)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Parse(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseLarge(b *testing.B) {
	data := []byte(largeMessage)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Parse(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseScanner(b *testing.B) {
	data := []byte(adta01Message)
	parser := NewParser()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.Parse(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncode(b *testing.B) {
	msg, err := Parse([]byte(adta01Message))
	if err != nil {
		b.Fatal(err)
	}
	encoder := NewEncoder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := encoder.Encode(msg)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncodeWithMLLP(b *testing.B) {
	msg, err := Parse([]byte(adta01Message))
	if err != nil {
		b.Fatal(err)
	}
	encoder := NewEncoder().WithMLLPFraming(true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := encoder.Encode(msg)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMessageGet(b *testing.B) {
	msg, err := Parse([]byte(adta01Message))
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = msg.Get("PID.3.1")
		_, _ = msg.Get("PID.5.1")
		_, _ = msg.Get("PID.5.2")
		_, _ = msg.Get("PID.7")
		_, _ = msg.Get("PID.8")
	}
}

func BenchmarkMessageGetRepeated(b *testing.B) {
	msg, err := Parse([]byte(largeMessage))
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = msg.Get("OBX[1].5")
		_, _ = msg.Get("OBX[5].5")
		_, _ = msg.Get("OBX[10].5")
	}
}

func BenchmarkSegmentField(b *testing.B) {
	msg, err := Parse([]byte(adta01Message))
	if err != nil {
		b.Fatal(err)
	}
	pid, _ := msg.Segment("PID")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pid.Field(3)
		_ = pid.Field(5)
		_ = pid.Field(7)
		_ = pid.Field(8)
	}
}

func BenchmarkSegmentComponent(b *testing.B) {
	msg, err := Parse([]byte(adta01Message))
	if err != nil {
		b.Fatal(err)
	}
	pid, _ := msg.Segment("PID")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pid.Component(5, 1)
		_ = pid.Component(5, 2)
		_ = pid.Component(3, 1)
	}
}

func BenchmarkUnmarshalStruct(b *testing.B) {
	type Patient struct {
		ID        string `hl7:"PID.3.1"`
		LastName  string `hl7:"PID.5.1"`
		FirstName string `hl7:"PID.5.2"`
		DOB       string `hl7:"PID.7"`
		Gender    string `hl7:"PID.8"`
	}

	data := []byte(adta01Message)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var patient Patient
		_ = Unmarshal(data, &patient)
	}
}

func BenchmarkMarshalToHL7(b *testing.B) {
	type Patient struct {
		ID        string `hl7:"PID.3.1"`
		LastName  string `hl7:"PID.5.1"`
		FirstName string `hl7:"PID.5.2"`
		DOB       string `hl7:"PID.7"`
		Gender    string `hl7:"PID.8"`
	}

	patient := Patient{
		ID:        "12345",
		LastName:  "Smith",
		FirstName: "John",
		DOB:       "19800115",
		Gender:    "M",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Marshal(patient)
	}
}

func BenchmarkSplitField(b *testing.B) {
	field := "Smith^John^A^Jr^Dr^MD"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SplitField(field, '^')
	}
}

func BenchmarkIsHL7Message(b *testing.B) {
	data := []byte(adta01Message)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = IsHL7Message(data)
	}
}

func BenchmarkStripMLLP(b *testing.B) {
	data := []byte{0x0B, 0x0D, 0x0A}
	data = append(data, []byte(adta01Message)...)
	data = append(data, 0x1C, 0x0D)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = StripMLLP(data)
	}
}

func BenchmarkSplitMessages(b *testing.B) {
	var multiMsg bytes.Buffer
	for i := 0; i < 10; i++ {
		multiMsg.WriteString(adta01Message)
		multiMsg.WriteString("\r\r")
	}
	data := multiMsg.Bytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SplitMessages(data)
	}
}

func BenchmarkLocationParse(b *testing.B) {
	locations := []string{
		"PID.3.1",
		"PID.5.2",
		"OBX[1].5",
		"OBX[10].3.1",
		"PV1.3.2",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, loc := range locations {
			_, _ = ParseLocation(loc)
		}
	}
}

func BenchmarkEscape(b *testing.B) {
	testStrings := []string{
		"Hello World",
		"Special | chars ^ here",
		"Newline\nHere",
		"Tab\there",
		"Backslash\\test",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, s := range testStrings {
			_ = Escape(s)
		}
	}
}

func BenchmarkUnescape(b *testing.B) {
	testStrings := []string{
		"Hello\\X20World",
		"Special\\S\\X7C\\chars",
		"Newline\\N\\X0A\\Here",
		"Tab\\T\\X09\\here",
		"Backslash\\\\test",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, s := range testStrings {
			_, _ = Unescape(s)
		}
	}
}

func BenchmarkValidateTable(b *testing.B) {
	msg, _ := Parse([]byte(adta01Message))

	validator := NewValidator(
		Required("MSH.9"),
		Required("PID.3.1"),
		RequiredTable("PID.8", "HL70001"),
		Table("MSH.12", "HL70396"),
		OneOf("PID.8", "M", "F", "O", "U"),
		Pattern("PID.7", `^\d{8}$`),
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validator.Validate(msg)
	}
}

func BenchmarkTableLookup(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = LookupTable("HL70001", "M")
		_, _ = LookupTable("HL70001", "F")
		_, _ = LookupTable("HL70001", "X")
	}
}

func BenchmarkSegmentDefinition(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = LookupSegmentDefinition("PID")
		_, _ = LookupSegmentDefinition("MSH")
		_, _ = LookupSegmentDefinition("OBX")
	}
}
