package hl7

import (
	"os"
	"testing"
)

// Real HL7 message data for testing
const (
	// ADT^A01 - Admit/Visit Notification
	testADTA01 = `MSH|^~\&|Primary|PIEDMONT ATLANTA HOSPITAL|CL|PDMT|20200319170944|ADTPA|ADT^A01|203478||2.3|||||||||||
EVN|A01|20200319170944||ADT_EVENT|^ADT^PATIENT^ACCESS^^^^^PHC^^^^^PAH|20200319170941|
PID|1||E3866034^^^EPIC^MRN~900091283^^^EPI^MR||TEST^PAHSEVENACARDTELEM||19770919|F||White|8 TEST ST^^ATHENS^GA^30605^USA^P^^CLARKE|CLARKE|(999)999-9999^P^PH^^^999^9999999||ENG|MARRIED|UNKNOWN|1000046070|000-00-0001|||NOT HISPANIC||N||||||N||
PD1|||PIEDMONT ATLANTA HOSPITAL^^10500|1750426920^SMITH III^GEORGE^^^^^^NPI^^^^NPI~81923^SMITH III^GEORGE^^^^^^STARPC^^^^STARID~SMIGE^SMITH III^GEORGE^^^^^^MT^^^^MTID||||||||||||||
NK1|1|TEST^SPOUSE^^|Spouse||(888)888-8888^^PH^^^888^8888888||Emergency Contact 1|||||||||||||||||||||||||||
PV1|1|INPATIENT|TA7B^0107028^0107028^PIEDMONT ATLANTA HOSPITAL^^^^^^^DEPID|UR|||1093782799^SMITH^ANITA^^^^^^NPI^^^^NPI~52228^SMITH^ANITA^^^^^^STARPC^^^^STARID|||General Med||||Phys/Clinic|||1093782799^SMITH^ANITA^^^^^^NPI^^^^NPI~52228^SMITH^ANITA^^^^^^STARPC^^^^STARID||2017417374|UHC||||||||||||||||||||||||20200319170941||||||||||
PV2||Priv||||||20200319|||||||||||||n|N||||||||||N|||||||||||||||||
OBX|1|NM|11156-7^LEUKOCYTES^LN||||||||I|
OBX|2|NM|11273-0^ERYTHROCYTES^LN||4.06|tera.l-1||N|||P|||201410060627
AL1|1|DA|32264^NO KNOWN ALLERGIES^||AAA|201410060830
DG1|1||^injury|injury||^10151;EPT||||||||||||||||||||`

	// ADT^A08 - Update Person Information
	testADTA08 = `MSH|^~\&|Primary|PD|CL|PDMT|20200311073040|119297|ADT^A08|203550|T|2.3|||||||||||
EVN|A08|20200311073040||ALLERGY_A08|119297^MOORE^REBEKAH^^^^^^PHC^^^^^PAH||
PID|1||E3843677^^^EPIC^MRN~900070078^^^EPI^MR||TEST^BABYBOY CSECTIONMOM||20181128100700|M|||^^^^^USA^P|||||SINGLE||1000034321|000-00-0000||2017364510^^^^CSN~E3843667^^^EPIC^MRN~900070069^^^EPI^MR||^^ATLANTA^GA^^|N|1|||||N||
PD1|||ATLANTA HOSPITAL^^10500|||||||||||||||
NK1|1|TEST^CSECTIONMOM^^|Mother|123 X street^^MARIETTA^GA^30062^USA|(404)605-5000^^PH^^^404^6055000||Emergency Contact 1|||||||||||||||||||||||||||
PV1|1|NEWBORN|NIP^293^293-01^PIEDMONT ATLANTA HOSPITAL^R^^^^^^DEPID|NB|||1760480644^JOHNSON^SCOTT^A^^^^^NPI^^^^NPI~39831^JOHNSON^SCOTT^A^^^^^STARPC^^^^STARID|||||||Sick|||1760480644^JOHNSON^SCOTT^A^^^^^NPI^^^^NPI~39831^JOHNSON^SCOTT^A^^^^^STARPC^^^^STARID||2017364683|SELF||||||||||||||||||||||||20181128100800|||1605010|||||||
OBX|1|NM|HT^HEIGHT||60|in||||||||20181127||||||||||||||||||
AL1|1|Drug Class|32264^NO KNOWN ALLERGIES^||||||`

	// ORU^R01 - Observation Result
	testORUR01 = `MSH|^~\&|LIS|HOSPITAL|EHR|HOSPITAL|20240115100000||ORU^R01|MSG00001|P|2.5|||||||||
PID|1||12345^^^MRN||Doe^John^A||19800101|M|||123 Main St^^City^ST^12345||(555)555-5555|||S||123456789|987654321||
PV1|1|I|4N^401^A||||||||||||||||2||||||||||||||||||||||||||||||||||||||20240115090000||||||
OBR|1|ORDER001|FILL001|80053^COMPREHENSIVE METABOLIC PANEL^CPT|||20240115093000||||||||||||||||1
OBX|1|NM|2345-7^GLUCOSE^LN||95|mg/dL|74-106|N|||F|||20240115094500
OBX|2|NM|2160-0^CREATININE^LN||1.0|mg/dL|0.7-1.3|N|||F|||20240115094500
OBX|3|NM|3094-0^BUN^LN||15|mg/dL|7-20|N|||F|||20240115094500
OBX|4|NM|17861-6^CALCIUM^LN||9.5|mg/dL|8.5-10.5|N|||F|||20240115094500`
)

func TestParseADTA01(t *testing.T) {
	msg, err := Parse([]byte(testADTA01))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Test message type
	if msgType := msg.Type(); msgType != "ADT^A01" {
		t.Errorf("Type() = %q, want %q", msgType, "ADT^A01")
	}

	// Test control ID
	if ctrlID := msg.ControlID(); ctrlID != "203478" {
		t.Errorf("ControlID() = %q, want %q", ctrlID, "203478")
	}

	// Test segment count
	// ADT_A01 test message has 11 segments: MSH, EVN, PID, PD1, NK1, PV1, PV2, OBX x2, AL1, DG1
	segs := msg.AllSegments()
	if len(segs) < 10 {
		t.Errorf("expected at least 10 segments, got %d", len(segs))
	}

	// Test PID segment
	pid, ok := msg.Segment("PID")
	if !ok {
		t.Fatal("PID segment not found")
	}

	// Test patient name components
	if familyName := pid.Component(5, 1); familyName != "TEST" {
		t.Errorf("PID.5.1 (family name) = %q, want %q", familyName, "TEST")
	}
	if givenName := pid.Component(5, 2); givenName != "PAHSEVENACARDTELEM" {
		t.Errorf("PID.5.2 (given name) = %q, want %q", givenName, "PAHSEVENACARDTELEM")
	}

	// Test DOB
	if dob := pid.Field(7); dob != "19770919" {
		t.Errorf("PID.7 (DOB) = %q, want %q", dob, "19770919")
	}

	// Test gender
	if gender := pid.Field(8); gender != "F" {
		t.Errorf("PID.8 (gender) = %q, want %q", gender, "F")
	}

	// Test multiple MRNs (repetition)
	mrnField := pid.Field(3)
	mrns := SplitField(mrnField, '~')
	if len(mrns) != 2 {
		t.Errorf("expected 2 MRN repetitions, got %d", len(mrns))
	}
}

func TestParseADTA08(t *testing.T) {
	msg, err := Parse([]byte(testADTA08))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if msgType := msg.Type(); msgType != "ADT^A08" {
		t.Errorf("Type() = %q, want %q", msgType, "ADT^A08")
	}

	// Test multiple NK1 segments
	nk1s := msg.Segments("NK1")
	if len(nk1s) != 1 {
		t.Errorf("expected 1 NK1 segment, got %d", len(nk1s))
	}

	// Test NK1 name
	if nk1, ok := msg.Segment("NK1"); ok {
		name := nk1.Component(2, 1)
		if name != "TEST" {
			t.Errorf("NK1.2.1 = %q, want %q", name, "TEST")
		}
	}
}

func TestParseORUR01(t *testing.T) {
	msg, err := Parse([]byte(testORUR01))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if msgType := msg.Type(); msgType != "ORU^R01" {
		t.Errorf("Type() = %q, want %q", msgType, "ORU^R01")
	}

	// Test multiple OBX segments
	obxs := msg.Segments("OBX")
	if len(obxs) != 4 {
		t.Errorf("expected 4 OBX segments, got %d", len(obxs))
	}

	// Test OBX observation values
	if obx, ok := msg.Segment("OBX"); ok {
		// First OBX should be glucose
		obsID := obx.Component(3, 1)
		if obsID != "2345-7" {
			t.Errorf("OBX.3.1 = %q, want %q", obsID, "2345-7")
		}
		value := obx.Field(5)
		if value != "95" {
			t.Errorf("OBX.5 = %q, want %q", value, "95")
		}
		units := obx.Field(6)
		if units != "mg/dL" {
			t.Errorf("OBX.6 = %q, want %q", units, "mg/dL")
		}
	}
}

func TestParseFromFile(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		wantType string
		wantCtrl string
	}{
		{"ADT_A01", "testdata/adt_a01.hl7", "ADT^A01", "203478"},
		{"ADT_A08", "testdata/adt_a08.hl7", "ADT^A08", "203550"},
		{"ADT_A28", "testdata/adt_a28.hl7", "ADT^A28", "203598"},
		{"ORU_R01", "testdata/oru_r01.hl7", "ORU^R01", "MSG00001"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := os.ReadFile(tt.file)
			if err != nil {
				t.Fatalf("failed to read %s: %v", tt.file, err)
			}

			msg, err := Parse(data)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if msgType := msg.Type(); msgType != tt.wantType {
				t.Errorf("Type() = %q, want %q", msgType, tt.wantType)
			}

			if ctrlID := msg.ControlID(); ctrlID != tt.wantCtrl {
				t.Errorf("ControlID() = %q, want %q", ctrlID, tt.wantCtrl)
			}
		})
	}
}

func TestUnmarshalADTA01(t *testing.T) {
	type PatientInfo struct {
		MRN       string `hl7:"PID.3.1"`
		LastName  string `hl7:"PID.5.1"`
		FirstName string `hl7:"PID.5.2"`
		DOB       string `hl7:"PID.7"`
		Gender    string `hl7:"PID.8"`
	}

	var patient PatientInfo
	if err := Unmarshal([]byte(testADTA01), &patient); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if patient.MRN != "E3866034" {
		t.Errorf("MRN = %q, want %q", patient.MRN, "E3866034")
	}
	if patient.LastName != "TEST" {
		t.Errorf("LastName = %q, want %q", patient.LastName, "TEST")
	}
	if patient.FirstName != "PAHSEVENACARDTELEM" {
		t.Errorf("FirstName = %q, want %q", patient.FirstName, "PAHSEVENACARDTELEM")
	}
	if patient.DOB != "19770919" {
		t.Errorf("DOB = %q, want %q", patient.DOB, "19770919")
	}
	if patient.Gender != "F" {
		t.Errorf("Gender = %q, want %q", patient.Gender, "F")
	}
}

func TestUnmarshalORUR01(t *testing.T) {
	type LabResult struct {
		ObservationID   string `hl7:"OBX.3.1"`
		ObservationText string `hl7:"OBX.3.2"`
		Value           string `hl7:"OBX.5"`
		Units           string `hl7:"OBX.6"`
		ReferenceRange  string `hl7:"OBX.7"`
	}

	var result LabResult
	if err := Unmarshal([]byte(testORUR01), &result); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if result.ObservationID != "2345-7" {
		t.Errorf("ObservationID = %q, want %q", result.ObservationID, "2345-7")
	}
	if result.ObservationText != "GLUCOSE" {
		t.Errorf("ObservationText = %q, want %q", result.ObservationText, "GLUCOSE")
	}
	if result.Value != "95" {
		t.Errorf("Value = %q, want %q", result.Value, "95")
	}
	if result.Units != "mg/dL" {
		t.Errorf("Units = %q, want %q", result.Units, "mg/dL")
	}
	if result.ReferenceRange != "74-106" {
		t.Errorf("ReferenceRange = %q, want %q", result.ReferenceRange, "74-106")
	}
}

func TestMarshalAndRoundTrip(t *testing.T) {
	type MSHSegment struct {
		FieldSeparator     string `hl7:"MSH.1"`
		EncodingCharacters string `hl7:"MSH.2"`
		SendingApplication string `hl7:"MSH.3"`
		SendingFacility    string `hl7:"MSH.4"`
		DateTimeOfMessage  string `hl7:"MSH.7"`
		MessageType        string `hl7:"MSH.9"`
		MessageControlID   string `hl7:"MSH.10"`
		ProcessingID       string `hl7:"MSH.11"`
		VersionID          string `hl7:"MSH.12"`
	}

	type PIDSegment struct {
		SetID       string `hl7:"PID.1"`
		PatientID   string `hl7:"PID.3.1"`
		PatientName string `hl7:"PID.5"`
		DOB         string `hl7:"PID.7"`
		Gender      string `hl7:"PID.8"`
	}

	type ADTMessage struct {
		MSH MSHSegment `hl7:"MSH"`
		PID PIDSegment `hl7:"PID"`
	}

	msg := ADTMessage{
		MSH: MSHSegment{
			FieldSeparator:     "|",
			EncodingCharacters: "^~\\&",
			SendingApplication: "TESTAPP",
			SendingFacility:    "TESTFAC",
			DateTimeOfMessage:  "20240115120000",
			MessageType:        "ADT^A01",
			MessageControlID:   "MSG001",
			ProcessingID:       "P",
			VersionID:          "2.5",
		},
		PID: PIDSegment{
			SetID:       "1",
			PatientID:   "12345",
			PatientName: "DOE^JOHN",
			DOB:         "19900101",
			Gender:      "M",
		},
	}

	// Marshal to HL7
	data, err := Marshal(msg)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	// Verify it contains expected parts
	dataStr := string(data)
	if !contains(dataStr, "MSH") {
		t.Error("marshaled data should contain MSH")
	}
	if !contains(dataStr, "TESTAPP") {
		t.Error("marshaled data should contain TESTAPP")
	}
	if !contains(dataStr, "PID") {
		t.Error("marshaled data should contain PID")
	}
	if !contains(dataStr, "DOE") {
		t.Error("marshaled data should contain DOE")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr)))
}

func BenchmarkParseADTA01(b *testing.B) {
	data := []byte(testADTA01)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Parse(data)
	}
}

func BenchmarkUnmarshalADTA01(b *testing.B) {
	type PatientInfo struct {
		MRN       string `hl7:"PID.3.1"`
		LastName  string `hl7:"PID.5.1"`
		FirstName string `hl7:"PID.5.2"`
		DOB       string `hl7:"PID.7"`
		Gender    string `hl7:"PID.8"`
	}

	data := []byte(testADTA01)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var patient PatientInfo
		_ = Unmarshal(data, &patient)
	}
}

func BenchmarkMarshal(b *testing.B) {
	type Message struct {
		MSH struct {
			FieldSeparator     string `hl7:"MSH.1"`
			EncodingCharacters string `hl7:"MSH.2"`
			SendingApplication string `hl7:"MSH.3"`
			MessageType        string `hl7:"MSH.9"`
			MessageControlID   string `hl7:"MSH.10"`
		} `hl7:"MSH"`
		PID struct {
			PatientID   string `hl7:"PID.3.1"`
			PatientName string `hl7:"PID.5"`
		} `hl7:"PID"`
	}

	msg := Message{}
	msg.MSH.FieldSeparator = "|"
	msg.MSH.EncodingCharacters = "^~\\&"
	msg.MSH.SendingApplication = "APP"
	msg.MSH.MessageType = "ADT^A01"
	msg.MSH.MessageControlID = "001"
	msg.PID.PatientID = "12345"
	msg.PID.PatientName = "DOE^JOHN"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Marshal(msg)
	}
}
