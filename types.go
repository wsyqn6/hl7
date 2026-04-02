package hl7

import (
	"fmt"
	"time"
)

// Timestamp represents an HL7 timestamp.
type Timestamp struct {
	time.Time
}

// NewTimestamp creates a new Timestamp from a time.Time.
func NewTimestamp(t time.Time) Timestamp {
	return Timestamp{Time: t}
}

// UnmarshalHL7 unmarshals an HL7 timestamp string.
func (t *Timestamp) UnmarshalHL7(data []byte) error {
	str := string(data)
	if str == "" {
		return nil
	}

	// Try various HL7 timestamp formats
	formats := []string{
		"20060102150405.9999-0700", // Full with timezone and fractional seconds
		"20060102150405-0700",      // With timezone
		"20060102150405",           // Full timestamp
		"200601021504",             // Minute precision
		"2006010215",               // Hour precision
		"20060102",                 // Date only
		"200601",                   // Year and month
		"2006",                     // Year only
	}

	for _, format := range formats {
		parsed, err := time.Parse(format, str)
		if err == nil {
			t.Time = parsed
			return nil
		}
	}

	return fmt.Errorf("unable to parse timestamp: %s", str)
}

// MarshalHL7 marshals the Timestamp to HL7 format.
func (t Timestamp) MarshalHL7() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte{}, nil
	}
	return []byte(t.Time.Format("20060102150405")), nil
}

// String returns the string representation of the Timestamp.
func (t Timestamp) String() string {
	if t.Time.IsZero() {
		return ""
	}
	return t.Time.Format("20060102150405")
}

// NM (Numeric) represents an HL7 numeric value.
type NM struct {
	Value float64
}

// UnmarshalHL7 unmarshals an HL7 numeric value.
func (n *NM) UnmarshalHL7(data []byte) error {
	str := string(data)
	if str == "" {
		return nil
	}
	_, err := fmt.Sscanf(str, "%f", &n.Value)
	return err
}

// MarshalHL7 marshals the NM to HL7 format.
func (n NM) MarshalHL7() ([]byte, error) {
	return []byte(fmt.Sprintf("%g", n.Value)), nil
}

// String returns the string representation of the NM.
func (n NM) String() string {
	return fmt.Sprintf("%g", n.Value)
}

// ID (Identifier) represents an HL7 identifier.
type ID struct {
	Value string
}

// UnmarshalHL7 unmarshals an HL7 identifier.
func (i *ID) UnmarshalHL7(data []byte) error {
	i.Value = string(data)
	return nil
}

// MarshalHL7 marshals the ID to HL7 format.
func (i ID) MarshalHL7() ([]byte, error) {
	return []byte(i.Value), nil
}

// String returns the string representation of the ID.
func (i ID) String() string {
	return i.Value
}

// CE (Coded Element) represents an HL7 coded element.
type CE struct {
	Identifier    string
	Text          string
	CodingSystem  string
	AltIdentifier string
	AltText       string
	AltCodingSys  string
}

// UnmarshalHL7 unmarshals an HL7 coded element.
func (c *CE) UnmarshalHL7(data []byte) error {
	parts := SplitField(string(data), '^')
	if len(parts) > 0 {
		c.Identifier = parts[0]
	}
	if len(parts) > 1 {
		c.Text = parts[1]
	}
	if len(parts) > 2 {
		c.CodingSystem = parts[2]
	}
	if len(parts) > 3 {
		c.AltIdentifier = parts[3]
	}
	if len(parts) > 4 {
		c.AltText = parts[4]
	}
	if len(parts) > 5 {
		c.AltCodingSys = parts[5]
	}
	return nil
}

// MarshalHL7 marshals the CE to HL7 format.
func (c CE) MarshalHL7() ([]byte, error) {
	result := c.Identifier
	if c.Text != "" {
		result += "^" + c.Text
	}
	if c.CodingSystem != "" {
		result += "^" + c.CodingSystem
	}
	if c.AltIdentifier != "" {
		result += "^" + c.AltIdentifier
	}
	if c.AltText != "" {
		result += "^" + c.AltText
	}
	if c.AltCodingSys != "" {
		result += "^" + c.AltCodingSys
	}
	return []byte(result), nil
}

// XPN (Extended Person Name) represents an HL7 extended person name.
type XPN struct {
	FamilyName  string
	GivenName   string
	MiddleName  string
	Suffix      string
	Prefix      string
	Degree      string
	Type        string
	Replication string
}

// UnmarshalHL7 unmarshals an HL7 extended person name.
func (x *XPN) UnmarshalHL7(data []byte) error {
	parts := SplitField(string(data), '^')
	if len(parts) > 0 {
		x.FamilyName = parts[0]
	}
	if len(parts) > 1 {
		x.GivenName = parts[1]
	}
	if len(parts) > 2 {
		x.MiddleName = parts[2]
	}
	if len(parts) > 3 {
		x.Suffix = parts[3]
	}
	if len(parts) > 4 {
		x.Prefix = parts[4]
	}
	if len(parts) > 5 {
		x.Degree = parts[5]
	}
	if len(parts) > 6 {
		x.Type = parts[6]
	}
	if len(parts) > 7 {
		x.Replication = parts[7]
	}
	return nil
}

// MarshalHL7 marshals the XPN to HL7 format.
func (x XPN) MarshalHL7() ([]byte, error) {
	result := x.FamilyName
	if x.GivenName != "" {
		result += "^" + x.GivenName
	}
	if x.MiddleName != "" {
		result += "^" + x.MiddleName
	}
	if x.Suffix != "" {
		result += "^" + x.Suffix
	}
	if x.Prefix != "" {
		result += "^" + x.Prefix
	}
	if x.Degree != "" {
		result += "^" + x.Degree
	}
	if x.Type != "" {
		result += "^" + x.Type
	}
	if x.Replication != "" {
		result += "^" + x.Replication
	}
	return []byte(result), nil
}

// FullName returns the full name as "GivenName FamilyName".
func (x XPN) FullName() string {
	if x.GivenName != "" && x.FamilyName != "" {
		return x.GivenName + " " + x.FamilyName
	}
	if x.FamilyName != "" {
		return x.FamilyName
	}
	return x.GivenName
}

// XAD (Extended Address) represents an HL7 extended address.
type XAD struct {
	Street       string
	Other        string
	City         string
	State        string
	ZipCode      string
	Country      string
	AddressType  string
	StartDate    string
	EndDate      string
	AssignmentID string
}

// UnmarshalHL7 unmarshals an HL7 extended address.
func (a *XAD) UnmarshalHL7(data []byte) error {
	parts := SplitField(string(data), '^')
	if len(parts) > 0 {
		a.Street = parts[0]
	}
	if len(parts) > 1 {
		a.Other = parts[1]
	}
	if len(parts) > 2 {
		a.City = parts[2]
	}
	if len(parts) > 3 {
		a.State = parts[3]
	}
	if len(parts) > 4 {
		a.ZipCode = parts[4]
	}
	if len(parts) > 5 {
		a.Country = parts[5]
	}
	if len(parts) > 6 {
		a.AddressType = parts[6]
	}
	if len(parts) > 7 {
		a.StartDate = parts[7]
	}
	if len(parts) > 8 {
		a.EndDate = parts[8]
	}
	if len(parts) > 9 {
		a.AssignmentID = parts[9]
	}
	return nil
}

// MarshalHL7 marshals the XAD to HL7 format.
func (a XAD) MarshalHL7() ([]byte, error) {
	result := a.Street
	if a.Other != "" {
		result += "^" + a.Other
	}
	if a.City != "" {
		result += "^" + a.City
	}
	if a.State != "" {
		result += "^" + a.State
	}
	if a.ZipCode != "" {
		result += "^" + a.ZipCode
	}
	if a.Country != "" {
		result += "^" + a.Country
	}
	if a.AddressType != "" {
		result += "^" + a.AddressType
	}
	if a.StartDate != "" {
		result += "^" + a.StartDate
	}
	if a.EndDate != "" {
		result += "^" + a.EndDate
	}
	if a.AssignmentID != "" {
		result += "^" + a.AssignmentID
	}
	return []byte(result), nil
}
