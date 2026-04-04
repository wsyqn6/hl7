package hl7

import "bytes"

// Encoder encodes HL7 messages to byte slices.
type Encoder struct {
	delims         Delimiters
	segmentEnd     []byte
	useMLLPFraming bool
	encoding       Encoding
	escapeHandling bool
}

// NewEncoder creates a new Encoder with default settings.
func NewEncoder() *Encoder {
	return &Encoder{
		delims:         DefaultDelimiters(),
		segmentEnd:     []byte{'\r'},
		encoding:       EncodingASCII,
		escapeHandling: true,
	}
}

// NewEncoderWithDelimiters creates a new Encoder with custom delimiters.
func NewEncoderWithDelimiters(delims Delimiters) *Encoder {
	return &Encoder{
		delims:         delims,
		segmentEnd:     []byte{'\r'},
		encoding:       EncodingASCII,
		escapeHandling: true,
	}
}

// WithMLLPFraming enables or disables MLLP framing.
func (e *Encoder) WithMLLPFraming(enabled bool) *Encoder {
	e.useMLLPFraming = enabled
	return e
}

// Encode encodes a Message to a byte slice.
func (e *Encoder) Encode(msg *Message) ([]byte, error) {
	var buf bytes.Buffer

	for i, seg := range msg.AllSegments() {
		// Encode segment name
		buf.WriteString(seg.Name())

		// Encode fields
		fields := seg.Fields()
		for j, field := range fields {
			if i == 0 && j == 0 && seg.Name() == "MSH" {
				// MSH segment: field 1 is the field separator
				buf.WriteByte(byte(e.delims.Field))
			} else {
				buf.WriteByte(byte(e.delims.Field))
				buf.WriteString(field)
			}
		}

		buf.Write(e.segmentEnd)
	}

	data := buf.Bytes()

	if e.useMLLPFraming {
		// Wrap in MLLP framing
		framed := make([]byte, 0, len(data)+3)
		framed = append(framed, 0x0B) // Start block
		framed = append(framed, data...)
		framed = append(framed, 0x1C, 0x0D) // End block + carriage return
		return framed, nil
	}

	return data, nil
}

// EncodeString encodes a Message to a string.
func (e *Encoder) EncodeString(msg *Message) (string, error) {
	data, err := e.Encode(msg)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Encode encodes a Message to a byte slice (package-level function).
func Encode(msg *Message) ([]byte, error) {
	e := NewEncoder()
	return e.Encode(msg)
}

// EncodeString encodes a Message to a string (package-level function).
func EncodeString(msg *Message) (string, error) {
	e := NewEncoder()
	return e.EncodeString(msg)
}

// EncodeWithMLLP encodes a Message with MLLP framing.
func EncodeWithMLLP(msg *Message) ([]byte, error) {
	e := NewEncoder().WithMLLPFraming(true)
	return e.Encode(msg)
}
