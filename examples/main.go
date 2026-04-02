package main

import (
	"fmt"
	"log"

	"github.com/wsyqn6/hl7"
)

func main() {
	// Example HL7 message
	data := []byte(`MSH|^~\&|SENDING_APP|SENDING_FACILITY|RECEIVING_APP|RECEIVING_FACILITY|20240115120000||ADT^A01|MSG00001|P|2.5|||
PID|1||12345^^^MRN||Smith^John^A||19800115|M|||123 Main St^^Springfield^IL^62701||555-1234|||M||ACCT001|123-45-6789`)

	// Parse the message
	msg, err := hl7.Parse(data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== Parsed HL7 Message ===")
	fmt.Printf("Message Type: %s\n", msg.Type())
	fmt.Printf("Control ID: %s\n", msg.ControlID())

	// Get PID segment
	if pid, ok := msg.Segment("PID"); ok {
		fmt.Println("\n=== PID Segment ===")
		fmt.Printf("Patient ID: %s\n", pid.Component(3, 1))
		fmt.Printf("Last Name: %s\n", pid.Component(5, 1))
		fmt.Printf("First Name: %s\n", pid.Component(5, 2))
		fmt.Printf("DOB: %s\n", pid.Field(7))
		fmt.Printf("Gender: %s\n", pid.Field(8))
	}

	// Encode the message back
	fmt.Println("\n=== Encoded Message ===")
	encoded, err := hl7.EncodeString(msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(encoded)

	// Unmarshal into struct
	fmt.Println("\n=== Unmarshaling to Struct ===")
	type Patient struct {
		MRN       string `hl7:"PID.3.1"`
		LastName  string `hl7:"PID.5.1"`
		FirstName string `hl7:"PID.5.2"`
		DOB       string `hl7:"PID.7"`
		Gender    string `hl7:"PID.8"`
	}

	var patient Patient
	if err := hl7.Unmarshal(data, &patient); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("MRN: %s\n", patient.MRN)
	fmt.Printf("Name: %s, %s\n", patient.LastName, patient.FirstName)
	fmt.Printf("DOB: %s\n", patient.DOB)
	fmt.Printf("Gender: %s\n", patient.Gender)

	// Generate ACK
	fmt.Println("\n=== Generating ACK ===")
	ack, err := hl7.Generate(msg, hl7.Accept())
	if err != nil {
		log.Fatal(err)
	}
	ackData, err := hl7.EncodeString(ack)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ackData)

	// Validation example
	fmt.Println("\n=== Validation Example ===")
	validator := hl7.NewValidator(
		hl7.Required("MSH.9"),
		hl7.Required("PID.3.1"),
		hl7.OneOf("PID.8", "M", "F", "O", "U"),
		hl7.Pattern("PID.7", `^\d{8}$`),
	)

	errors := validator.Validate(msg)
	if len(errors) > 0 {
		fmt.Println("Validation errors:")
		for _, err := range errors {
			fmt.Printf("  - %s\n", err)
		}
	} else {
		fmt.Println("Message is valid!")
	}
}
