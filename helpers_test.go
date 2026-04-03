package hl7

import (
	"testing"
)

func TestPIDHelper(t *testing.T) {
	data := []byte("MSH|^~\\&|Primary|PIEDMONT ATLANTA HOSPITAL|CL|PDMT|20200319170944|ADTPA|ADT^A01|203478||2.3|||||||||||\r" +
		"PID|1||E3866034^^^EPIC^MRN~900091283^^^EPI^MR||TEST^PAHSEVENACARDTELEM||19770919|F||White|8 TEST ST^^ATHENS^GA^30605^USA^P^^CLARKE|CLARKE|(999)999-9998^P^PH^^^999^9999999||ENG|MARRIED|UNKNOWN|1000046070|000-00-0001")

	msg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	pid := msg.PID()

	if !pid.Exists() {
		t.Error("PID() should exist")
	}

	if pid.PatientID() != "E3866034^^^EPIC^MRN~900091283^^^EPI^MR" {
		t.Errorf("PatientID() = %q, want %q", pid.PatientID(), "E3866034^^^EPIC^MRN~900091283^^^EPI^MR")
	}

	if pid.LastName() != "TEST" {
		t.Errorf("LastName() = %q, want %q", pid.LastName(), "TEST")
	}

	if pid.FirstName() != "PAHSEVENACARDTELEM" {
		t.Errorf("FirstName() = %q, want %q", pid.FirstName(), "PAHSEVENACARDTELEM")
	}

	if pid.Sex() != "F" {
		t.Errorf("Sex() = %q, want %q", pid.Sex(), "F")
	}

	if pid.DateOfBirth() != "19770919" {
		t.Errorf("DateOfBirth() = %q, want %q", pid.DateOfBirth(), "19770919")
	}

	if pid.Race() != "White" {
		t.Errorf("Race() = %q, want %q", pid.Race(), "White")
	}

	if pid.City() != "ATHENS" {
		t.Errorf("City() = %q, want %q", pid.City(), "ATHENS")
	}

	if pid.State() != "GA" {
		t.Errorf("State() = %q, want %q", pid.State(), "GA")
	}

	if pid.PostalCode() != "30605" {
		t.Errorf("PostalCode() = %q, want %q", pid.PostalCode(), "30605")
	}

	if pid.PrimaryLanguage() != "ENG" {
		t.Errorf("PrimaryLanguage() = %q, want %q", pid.PrimaryLanguage(), "ENG")
	}

	if pid.MaritalStatus() != "MARRIED" {
		t.Errorf("MaritalStatus() = %q, want %q", pid.MaritalStatus(), "MARRIED")
	}

	if pid.Religion() != "UNKNOWN" {
		t.Errorf("Religion() = %q, want %q", pid.Religion(), "UNKNOWN")
	}

	if pid.SSN() != "000-00-0001" {
		t.Errorf("SSN() = %q, want %q", pid.SSN(), "000-00-0001")
	}
}

func TestPIDHelperNotExists(t *testing.T) {
	msg := NewMessage()
	pid := msg.PID()
	if pid.Exists() {
		t.Error("PID() should not exist for empty message")
	}
}

func TestMSHHelper(t *testing.T) {
	data := []byte("MSH|^~\\&|SENDING_APP|SENDING_FAC|||20240115120000||ADT^A01|MSG001|P|2.5")

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

	if msh.SendingFacility() != "SENDING_FAC" {
		t.Errorf("SendingFacility() = %q, want %q", msh.SendingFacility(), "SENDING_FAC")
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

func TestMSHHelperNotExists(t *testing.T) {
	msg := NewMessage()
	msh := msg.MSH()
	if msh.Exists() {
		t.Error("MSH() should not exist for empty message")
	}
}

func TestPV1Helper(t *testing.T) {
	data := []byte("MSH|^~\\&|Primary|PIEDMONT ATLANTA HOSPITAL|CL|PDMT|20200319170944|ADTPA|ADT^A01|203478||2.3|||||||||||\r" +
		"PV1|1|INPATIENT|TA7B^0107028^0107028^PIEDMONT ATLANTA HOSPITAL^^^^^^^DEPID|UR|||1093782799^SMITH^ANITA^^^^^^NPI^^^^NPI~52228^SMITH^ANITA^^^^^^STARPC^^^^STARID|||General Med||||Phys/Clinic|||1093782799^SMITH^ANITA^^^^^^NPI^^^^NPI~52228^SMITH^ANITA^^^^^^STARPC^^^^STARID||2017417374|UHC||||||||||||||||||||||||20200319170941||||||||||")

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

	if pv1.PatientClass() != "INPATIENT" {
		t.Errorf("PatientClass() = %q, want %q", pv1.PatientClass(), "INPATIENT")
	}

	if pv1.LocationPointOfCare() != "TA7B" {
		t.Errorf("LocationPointOfCare() = %q, want %q", pv1.LocationPointOfCare(), "TA7B")
	}

	if pv1.LocationRoom() != "0107028" {
		t.Errorf("LocationRoom() = %q, want %q", pv1.LocationRoom(), "0107028")
	}

	if pv1.AdmissionDate() != "20200319170941" {
		t.Errorf("AdmissionDate() = %q, want %q", pv1.AdmissionDate(), "20200319170941")
	}
}

func TestOBRHelper(t *testing.T) {
	data := []byte("MSH|^~\\&|LIS|HOSPITAL|EHR|HOSPITAL|20240115100000||ORU^R01|MSG00001|P|2.5|||||||||\r" +
		"OBR|1|ORDER001|FILL001|80053^COMPREHENSIVE METABOLIC PANEL^CPT|||20240115093000||||||||||||||||1")

	msg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	obr := msg.OBR()
	if !obr.Exists() {
		t.Error("OBR() should exist")
	}

	if obr.PlacerOrderNumber() != "ORDER001" {
		t.Errorf("PlacerOrderNumber() = %q, want %q", obr.PlacerOrderNumber(), "ORDER001")
	}

	if obr.FillerOrderNumber() != "FILL001" {
		t.Errorf("FillerOrderNumber() = %q, want %q", obr.FillerOrderNumber(), "FILL001")
	}

	if obr.UniversalServiceID() != "80053^COMPREHENSIVE METABOLIC PANEL^CPT" {
		t.Errorf("UniversalServiceID() = %q, want %q", obr.UniversalServiceID(), "80053^COMPREHENSIVE METABOLIC PANEL^CPT")
	}

	if obr.ServiceIdentifier() != "80053" {
		t.Errorf("ServiceIdentifier() = %q, want %q", obr.ServiceIdentifier(), "80053")
	}

	if obr.ServiceText() != "COMPREHENSIVE METABOLIC PANEL" {
		t.Errorf("ServiceText() = %q, want %q", obr.ServiceText(), "COMPREHENSIVE METABOLIC PANEL")
	}
}

func TestOBXHelper(t *testing.T) {
	data := []byte("MSH|^~\\&|APP|FAC||||||MSG001|P|2.5\r" +
		"OBX|1|NM|WBC^White Blood Cell Count^LN||7.5|10*3/uL|4.5-11.0|N|||F|||20240101120000")

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

	if obx.ObservationIdentifierCode() != "WBC" {
		t.Errorf("ObservationIdentifierCode() = %q, want %q", obx.ObservationIdentifierCode(), "WBC")
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

func TestAllRepeatingSegments(t *testing.T) {
	data := []byte("MSH|^~\\&|APP|FAC||||||MSG001|P|2.5\r" +
		"OBX|1|NM|WBC^WBC||7.5|10*3/uL||||F|\r" +
		"OBX|2|NM|RBC^RBC||4.8|10*6/uL||||F|\r" +
		"OBX|3|NM|HGB^Hemoglobin||14.2|g/dL||||F|")

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

func TestNK1Helper(t *testing.T) {
	data := []byte("MSH|^~\\&|Primary|PIEDMONT ATLANTA HOSPITAL|CL|PDMT|20200319170944|ADTPA|ADT^A01|203478||2.3|||||||||||\r" +
		"NK1|1|TEST^SPOUSE^^|Spouse||(888)888-8888^^PH^^^888^8888888||Emergency Contact 1|||||||||||||||||||||||||||")

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

	if nk1.Name() != "TEST^SPOUSE^^" {
		t.Errorf("Name() = %q, want %q", nk1.Name(), "TEST^SPOUSE^^")
	}

	if nk1.Relationship() != "Spouse" {
		t.Errorf("Relationship() = %q, want %q", nk1.Relationship(), "Spouse")
	}

	if nk1.PhoneNumber() != "(888)888-8888^^PH^^^888^8888888" {
		t.Errorf("PhoneNumber() = %q, want %q", nk1.PhoneNumber(), "(888)888-8888^^PH^^^888^8888888")
	}
}

func TestDG1Helper(t *testing.T) {
	data := []byte("MSH|^~\\&|Primary|PIEDMONT ATLANTA HOSPITAL|CL|PDMT|20200319170944|ADTPA|ADT^A01|203478||2.3|||||||||||\r" +
		"DG1|1||^injury|injury||^10151;EPT||||||||||||||||||||")

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

	if dg1.DiagnosisDescription() != "injury" {
		t.Errorf("DiagnosisDescription() = %q, want %q", dg1.DiagnosisDescription(), "injury")
	}
}
