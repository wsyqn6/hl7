package hl7

import (
	"testing"
	"time"
)

func TestMSHHelperMethods(t *testing.T) {
	data := []byte("MSH|^~\\&|SENDING_APP|SENDING_FACILITY|||20240115120000||ADT^A01|MSG001|P|2.5")

	msg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	msh := msg.MSH()

	if !msh.Exists() {
		t.Error("MSH() should exist")
	}

	if msh.FieldSeparator() != "|" {
		t.Errorf("FieldSeparator() = %q, want %q", msh.FieldSeparator(), "|")
	}

	if msh.EncodingCharacters() != "^~\\&" {
		t.Errorf("EncodingCharacters() = %q, want %q", msh.EncodingCharacters(), "^~\\&")
	}

	if msh.SendingApplication() != "SENDING_APP" {
		t.Errorf("SendingApplication() = %q, want %q", msh.SendingApplication(), "SENDING_APP")
	}

	if msh.SendingFacility() != "SENDING_FACILITY" {
		t.Errorf("SendingFacility() = %q, want %q", msh.SendingFacility(), "SENDING_FACILITY")
	}

	if msh.MessageType() != "ADT^A01" {
		t.Errorf("MessageType() = %q, want %q", msh.MessageType(), "ADT^A01")
	}

	if msh.MessageTypeCode() != "ADT" {
		t.Errorf("MessageTypeCode() = %q, want %q", msh.MessageTypeCode(), "ADT")
	}

	if msh.MessageTypeTrigger() != "A01" {
		t.Errorf("MessageTypeTrigger() = %q, want %q", msh.MessageTypeTrigger(), "A01")
	}

	if msh.MessageControlID() != "MSG001" {
		t.Errorf("MessageControlID() = %q, want %q", msh.MessageControlID(), "MSG001")
	}

	if msh.ProcessingID() != "P" {
		t.Errorf("ProcessingID() = %q, want %q", msh.ProcessingID(), "P")
	}

	if msh.VersionID() != "2.5" {
		t.Errorf("VersionID() = %q, want %q", msh.VersionID(), "2.5")
	}
}

func TestMSHHelperDateTime(t *testing.T) {
	data := []byte("MSH|^~\\&|APP|FAC|||20240115120000||ADT^A01|MSG001|P|2.5")

	msg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	msh := msg.MSH()

	dt := msh.DateTime()
	expected := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	if !dt.Equal(expected) {
		t.Errorf("DateTime() = %v, want %v", dt, expected)
	}
}

func TestPIDHelperMethods(t *testing.T) {
	data := []byte("MSH|^~\\&|APP|FAC||||||MSG001|P|2.5\r" +
		"PID|1||12345^^^MRN||Smith^John^A||19800115|M")

	msg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	pid := msg.PID()

	if !pid.Exists() {
		t.Error("PID() should exist")
	}

	if pid.PatientID() != "12345^^^MRN" {
		t.Errorf("PatientID() = %q, want %q", pid.PatientID(), "12345^^^MRN")
	}

	if pid.PatientName() != "Smith^John^A" {
		t.Errorf("PatientName() = %q, want %q", pid.PatientName(), "Smith^John^A")
	}

	if pid.LastName() != "Smith" {
		t.Errorf("LastName() = %q, want %q", pid.LastName(), "Smith")
	}

	if pid.FirstName() != "John" {
		t.Errorf("FirstName() = %q, want %q", pid.FirstName(), "John")
	}

	if pid.MiddleName() != "A" {
		t.Errorf("MiddleName() = %q, want %q", pid.MiddleName(), "A")
	}

	if pid.DateOfBirth() != "19800115" {
		t.Errorf("DateOfBirth() = %q, want %q", pid.DateOfBirth(), "19800115")
	}

	if pid.Sex() != "M" {
		t.Errorf("Sex() = %q, want %q", pid.Sex(), "M")
	}

	if pid.Gender() != "M" {
		t.Errorf("Gender() = %q, want %q", pid.Gender(), "M")
	}
}

func TestPIDHelperDateOfBirthTime(t *testing.T) {
	data := []byte("MSH|^~\\&|APP|FAC||||||MSG001|P|2.5\r" +
		"PID|1||12345||Smith||19800115000000|M")

	msg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	pid := msg.PID()

	dt, err := pid.DateOfBirthTime()
	if err != nil {
		t.Fatalf("DateOfBirthTime() error = %v", err)
	}

	expected := time.Date(1980, 1, 15, 0, 0, 0, 0, time.UTC)
	if !dt.Equal(expected) {
		t.Errorf("DateOfBirthTime() = %v, want %v", dt, expected)
	}
}

func TestPIDHelperAddress(t *testing.T) {
	data := []byte("MSH|^~\\&|APP|FAC||||||MSG001|P|2.5\r" +
		"PID|1||12345||Smith||19800115|M|||123 Main St^Apt 4^City^ST^12345^USA")

	msg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	pid := msg.PID()

	addr := pid.PatientAddress()
	if addr == "" {
		t.Error("PatientAddress() should not be empty")
	}

	components := SplitField(addr, '^')
	if len(components) < 6 {
		t.Errorf("PatientAddress() has %d components, want at least 6", len(components))
	}

	city := pid.City()
	if city != "City" {
		t.Errorf("City() = %q, want %q", city, "City")
	}

	state := pid.State()
	if state != "ST" {
		t.Errorf("State() = %q, want %q", state, "ST")
	}

	postal := pid.PostalCode()
	if postal != "12345" {
		t.Errorf("PostalCode() = %q, want %q", postal, "12345")
	}
}

func TestPV1HelperMethods(t *testing.T) {
	data := []byte("MSH|^~\\&|APP|FAC||||||MSG001|P|2.5\r" +
		"PV1|1|I|ICU^Room101^Bed1")

	msg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	pv1 := msg.PV1()

	if !pv1.Exists() {
		t.Error("PV1() should exist")
	}

	if pv1.SetID() != "1" {
		t.Errorf("SetID() = %q, want %q", pv1.SetID(), "1")
	}

	if pv1.PatientClass() != "I" {
		t.Errorf("PatientClass() = %q, want %q", pv1.PatientClass(), "I")
	}

	if pv1.LocationPointOfCare() != "ICU" {
		t.Errorf("LocationPointOfCare() = %q, want %q", pv1.LocationPointOfCare(), "ICU")
	}

	if pv1.LocationRoom() != "Room101" {
		t.Errorf("LocationRoom() = %q, want %q", pv1.LocationRoom(), "Room101")
	}

	if pv1.LocationBed() != "Bed1" {
		t.Errorf("LocationBed() = %q, want %q", pv1.LocationBed(), "Bed1")
	}
}

func TestPV1HelperAttendingDoctor(t *testing.T) {
	data := []byte("MSH|^~\\&|APP|FAC||||||MSG001|P|2.5\r" +
		"PV1|1|I|||||^AttendingDoctor^12345")

	msg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	pv1 := msg.PV1()

	doc := pv1.AttendingDoctor()
	if doc == "" {
		t.Error("AttendingDoctor() should not be empty")
	}
}

func TestOBRHelperMethods(t *testing.T) {
	data := []byte("MSH|^~\\&|APP|FAC||||||MSG001|P|2.5\r" +
		"OBR|1|ORDER001|FILL001|80053^COMPREHENSIVE METABOLIC PANEL^CPT")

	msg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	obr := msg.OBR()

	if !obr.Exists() {
		t.Error("OBR() should exist")
	}

	if obr.SetID() != "1" {
		t.Errorf("SetID() = %q, want %q", obr.SetID(), "1")
	}

	if obr.PlacerOrderNumber() != "ORDER001" {
		t.Errorf("PlacerOrderNumber() = %q, want %q", obr.PlacerOrderNumber(), "ORDER001")
	}

	if obr.FillerOrderNumber() != "FILL001" {
		t.Errorf("FillerOrderNumber() = %q, want %q", obr.FillerOrderNumber(), "FILL001")
	}

	if obr.UniversalServiceID() != "80053^COMPREHENSIVE METABOLIC PANEL^CPT" {
		t.Errorf("UniversalServiceID() = %q", obr.UniversalServiceID())
	}

	if obr.ServiceIdentifier() != "80053" {
		t.Errorf("ServiceIdentifier() = %q, want %q", obr.ServiceIdentifier(), "80053")
	}

	if obr.ServiceText() != "COMPREHENSIVE METABOLIC PANEL" {
		t.Errorf("ServiceText() = %q, want %q", obr.ServiceText(), "COMPREHENSIVE METABOLIC PANEL")
	}
}

func TestOBXHelperMethods(t *testing.T) {
	data := []byte("MSH|^~\\&|APP|FAC||||||MSG001|P|2.5\r" +
		"OBX|1|NM|WBC^White Blood Cell Count^LN||7.5|10*3/uL|4.5-11.0|N|||F")

	msg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	obx := msg.OBX()

	if !obx.Exists() {
		t.Error("OBX() should exist")
	}

	if obx.SetID() != "1" {
		t.Errorf("SetID() = %q, want %q", obx.SetID(), "1")
	}

	if obx.ValueType() != "NM" {
		t.Errorf("ValueType() = %q, want %q", obx.ValueType(), "NM")
	}

	if obx.ObservationIdentifier() != "WBC^White Blood Cell Count^LN" {
		t.Errorf("ObservationIdentifier() = %q", obx.ObservationIdentifier())
	}

	if obx.ObservationIdentifierCode() != "WBC" {
		t.Errorf("ObservationIdentifierCode() = %q, want %q", obx.ObservationIdentifierCode(), "WBC")
	}

	if obx.ObservationIdentifierText() != "White Blood Cell Count" {
		t.Errorf("ObservationIdentifierText() = %q, want %q", obx.ObservationIdentifierText(), "White Blood Cell Count")
	}

	if obx.ObservationValue() != "7.5" {
		t.Errorf("ObservationValue() = %q, want %q", obx.ObservationValue(), "7.5")
	}

	if obx.Units() != "10*3/uL" {
		t.Errorf("Units() = %q, want %q", obx.Units(), "10*3/uL")
	}

	if obx.ReferenceRange() != "4.5-11.0" {
		t.Errorf("ReferenceRange() = %q, want %q", obx.ReferenceRange(), "4.5-11.0")
	}

	if obx.ResultStatus() != "F" {
		t.Errorf("ResultStatus() = %q, want %q", obx.ResultStatus(), "F")
	}
}

func TestOBXHelperAbnormalFlags(t *testing.T) {
	data := []byte("MSH|^~\\&|APP|FAC||||||MSG001|P|2.5\r" +
		"OBX|1|NM|WBC||7.5|10*3/uL||H^L|||F")

	msg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	obx := msg.OBX()

	flags := obx.AbnormalFlags()
	if len(flags) != 2 {
		t.Errorf("AbnormalFlags() returned %d flags, want 2", len(flags))
		return
	}

	if flags[0] != "H" {
		t.Errorf("AbnormalFlags()[0] = %q, want %q", flags[0], "H")
	}

	if flags[1] != "L" {
		t.Errorf("AbnormalFlags()[1] = %q, want %q", flags[1], "L")
	}
}

func TestAllRepeatingSegmentsMethods(t *testing.T) {
	data := []byte("MSH|^~\\&|APP|FAC||||||MSG001|P|2.5\r" +
		"OBX|1|NM|WBC||7.5|10*3/uL||||F|\r" +
		"OBX|2|NM|RBC||4.8|10*6/uL||||F|\r" +
		"OBX|3|NM|HGB||14.2|g/dL||||F|")

	msg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	allOBX := msg.AllOBX()
	if len(allOBX) != 3 {
		t.Errorf("AllOBX() returned %d segments, want 3", len(allOBX))
	}

	if allOBX[0].ObservationValue() != "7.5" {
		t.Errorf("OBX[0].ObservationValue() = %q, want %q", allOBX[0].ObservationValue(), "7.5")
	}

	if allOBX[1].ObservationValue() != "4.8" {
		t.Errorf("OBX[1].ObservationValue() = %q, want %q", allOBX[1].ObservationValue(), "4.8")
	}

	if allOBX[2].ObservationValue() != "14.2" {
		t.Errorf("OBX[2].ObservationValue() = %q, want %q", allOBX[2].ObservationValue(), "14.2")
	}
}

func TestNK1HelperMethods(t *testing.T) {
	data := []byte("MSH|^~\\&|APP|FAC||||||MSG001|P|2.5\r" +
		"NK1|1|Smith^Jane^SP|Sister|123 Main St^City^ST^12345||555-9999")

	msg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	nk1 := msg.NK1()

	if !nk1.Exists() {
		t.Error("NK1() should exist")
	}

	if nk1.SetID() != "1" {
		t.Errorf("SetID() = %q, want %q", nk1.SetID(), "1")
	}

	if nk1.Name() != "Smith^Jane^SP" {
		t.Errorf("Name() = %q, want %q", nk1.Name(), "Smith^Jane^SP")
	}

	if nk1.Relationship() != "Sister" {
		t.Errorf("Relationship() = %q, want %q", nk1.Relationship(), "Sister")
	}

	addr := nk1.Address()
	if addr == "" {
		t.Error("Address() should not be empty")
	}
}

func TestDG1HelperMethods(t *testing.T) {
	data := []byte("MSH|^~\\&|APP|FAC||||||MSG001|P|2.5\r" +
		"DG1|1||ICD10||Pneumonia|20240115||W|")

	msg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	dg1 := msg.DG1()

	if !dg1.Exists() {
		t.Error("DG1() should exist")
	}

	if dg1.SetID() != "1" {
		t.Errorf("SetID() = %q, want %q", dg1.SetID(), "1")
	}

	t.Logf("DG1 segment fields:")
	for i := 1; i <= 10; i++ {
		t.Logf("  DG1.%d = %q", i, dg1.seg.Field(i))
	}

	code := dg1.DiagnosisCode()
	if code == "" {
		t.Error("DiagnosisCode() should not be empty")
	}
}

func TestSegmentHelperNotExists(t *testing.T) {
	msg := NewMessage()

	tests := []struct {
		name   string
		helper interface{ Exists() bool }
	}{
		{"PID", msg.PID()},
		{"MSH", msg.MSH()},
		{"PV1", msg.PV1()},
		{"OBR", msg.OBR()},
		{"OBX", msg.OBX()},
		{"NK1", msg.NK1()},
		{"DG1", msg.DG1()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.helper.Exists() {
				t.Errorf("%s() should not exist for empty message", tt.name)
			}
		})
	}
}
