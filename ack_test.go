package hl7

import (
	"testing"
)

func TestGenerateACK(t *testing.T) {
	data := []byte("MSH|^~\\&|APP|FAC|||20240101120000||ADT^A01|MSG001|P|2.5")

	msg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	ack, err := Generate(msg)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if !IsACK(ack) {
		t.Error("IsACK() should return true")
	}

	code, err := GetACKCode(ack)
	if err != nil {
		t.Fatalf("GetACKCode() error = %v", err)
	}
	if code != AA {
		t.Errorf("ACK code = %q, want %q", code, AA)
	}
}

func TestGenerateACKWithError(t *testing.T) {
	data := []byte("MSH|^~\\&|APP|FAC|||20240101120000||ADT^A01|MSG001|P|2.5")

	msg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	ack, err := Generate(msg, Error("Test error message"))
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	code, err := GetACKCode(ack)
	if err != nil {
		t.Fatalf("GetACKCode() error = %v", err)
	}
	if code != AE {
		t.Errorf("ACK code = %q, want %q", code, AE)
	}

	text, err := GetACKText(ack)
	if err != nil {
		t.Fatalf("GetACKText() error = %v", err)
	}
	if text != "Test error message" {
		t.Errorf("ACK text = %q, want %q", text, "Test error message")
	}
}

func TestGenerateACKWithReject(t *testing.T) {
	data := []byte("MSH|^~\\&|APP|FAC|||20240101120000||ADT^A01|MSG001|P|2.5")

	msg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	ack, err := Generate(msg, Reject("Rejected"))
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	code, err := GetACKCode(ack)
	if err != nil {
		t.Fatalf("GetACKCode() error = %v", err)
	}
	if code != AR {
		t.Errorf("ACK code = %q, want %q", code, AR)
	}
}

func TestGenerateACKWithERR(t *testing.T) {
	data := []byte("MSH|^~\\&|APP|FAC|||20240101120000||ADT^A01|MSG001|P|2.5")

	msg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	ack, err := Generate(msg,
		Error("Validation failed"),
		WithERR("207", "E", "Invalid patient ID", "PID.3"),
	)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	errs := GetErrors(ack)
	if len(errs) != 1 {
		t.Fatalf("GetErrors() returned %d errors, want 1", len(errs))
	}

	if errs[0].Diagnostics != "Invalid patient ID" {
		t.Errorf("ERR Diagnostics = %q, want %q", errs[0].Diagnostics, "Invalid patient ID")
	}

	if errs[0].ErrorLocation != "PID.3" {
		t.Errorf("ERR ErrorLocation = %q, want %q", errs[0].ErrorLocation, "PID.3")
	}
}

func TestGenerateACKWithValidationErrors(t *testing.T) {
	data := []byte("MSH|^~\\&|APP|FAC|||20240101120000||ADT^A01|MSG001|P|2.5")

	msg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	validationErrs := []*ValidationError{
		{Location: "PID.3", Message: "Patient ID is required"},
		{Location: "PID.5", Message: "Patient name is required"},
	}

	ack, err := GenerateACKWithValidation(msg, validationErrs)
	if err != nil {
		t.Fatalf("GenerateACKWithValidation() error = %v", err)
	}

	code, err := GetACKCode(ack)
	if err != nil {
		t.Fatalf("GetACKCode() error = %v", err)
	}
	if code != AR {
		t.Errorf("ACK code = %q, want %q", code, AR)
	}

	errs := GetErrors(ack)
	if len(errs) != 2 {
		t.Errorf("GetErrors() returned %d errors, want 2", len(errs))
	}
}

func TestIsAccepted(t *testing.T) {
	data := []byte("MSH|^~\\&|APP|FAC|||20240101120000||ADT^A01|MSG001|P|2.5")

	msg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	tests := []struct {
		opts     []ACKOption
		accepted bool
		rejected bool
		hasError bool
	}{
		{[]ACKOption{Accept()}, true, false, false},
		{[]ACKOption{Error("error")}, false, false, true},
		{[]ACKOption{Reject("reject")}, false, true, false},
		{[]ACKOption{CommitAccept()}, true, false, false},
		{[]ACKOption{CommitError("error")}, false, false, true},
		{[]ACKOption{CommitReject("reject")}, false, true, false},
	}

	for _, tt := range tests {
		ack, err := Generate(msg, tt.opts...)
		if err != nil {
			t.Fatalf("Generate() error = %v", err)
		}

		accepted, err := IsAccepted(ack)
		if err != nil {
			t.Fatalf("IsAccepted() error = %v", err)
		}
		if accepted != tt.accepted {
			t.Errorf("IsAccepted() = %v, want %v", accepted, tt.accepted)
		}

		rejected, err := IsRejected(ack)
		if err != nil {
			t.Fatalf("IsRejected() error = %v", err)
		}
		if rejected != tt.rejected {
			t.Errorf("IsRejected() = %v, want %v", rejected, tt.rejected)
		}

		hasError, err := HasError(ack)
		if err != nil {
			t.Fatalf("HasError() error = %v", err)
		}
		if hasError != tt.hasError {
			t.Errorf("HasError() = %v, want %v", hasError, tt.hasError)
		}
	}
}

func TestIsACKWithNil(t *testing.T) {
	if IsACK(nil) {
		t.Error("IsACK(nil) should return false")
	}
}

func TestGetACKCodeWithNonACK(t *testing.T) {
	data := []byte("MSH|^~\\&|APP|FAC|||20240101120000||ADT^A01|MSG001|P|2.5")
	msg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	_, err = GetACKCode(msg)
	if err == nil {
		t.Error("GetACKCode() should return error for non-ACK message")
	}
}

func TestCommitACKCodes(t *testing.T) {
	data := []byte("MSH|^~\\&|APP|FAC|||20240101120000||ADT^A01|MSG001|P|2.5")

	msg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	ack, err := Generate(msg, CommitAccept())
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	code, _ := GetACKCode(ack)
	if code != CA {
		t.Errorf("ACK code = %q, want %q", code, CA)
	}

	ack, err = Generate(msg, CommitError("CE error"))
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	code, _ = GetACKCode(ack)
	if code != ACKCE {
		t.Errorf("ACK code = %q, want %q", code, ACKCE)
	}

	ack, err = Generate(msg, CommitReject("CR reject"))
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	code, _ = GetACKCode(ack)
	if code != CR {
		t.Errorf("ACK code = %q, want %q", code, CR)
	}
}

func TestERRInfoString(t *testing.T) {
	info := ERRInfo{
		ErrorLocation: "PID.3",
		Diagnostics:   "Patient ID is required",
	}

	s := info.String()
	if s != "PID.3: Patient ID is required" {
		t.Errorf("ERRInfo.String() = %q, want %q", s, "PID.3: Patient ID is required")
	}
}
