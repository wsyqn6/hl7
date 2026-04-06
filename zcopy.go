package hl7

import (
	"bytes"
	"strings"
	"unsafe"
)

func bytesToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return unsafe.String(&b[0], len(b))
}

func stringToBytes(s string) []byte {
	if s == "" {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

type bytesBuffer struct {
	data []byte
}

func newBytesBuffer(size int) *bytesBuffer {
	return &bytesBuffer{
		data: make([]byte, 0, size),
	}
}

func (b *bytesBuffer) WriteString(s string) {
	b.data = append(b.data, s...)
}

func (b *bytesBuffer) WriteByte(c byte) error {
	b.data = append(b.data, c)
	return nil
}

func (b *bytesBuffer) Bytes() []byte {
	return b.data
}

func (b *bytesBuffer) String() string {
	return bytesToString(b.data)
}

func (b *bytesBuffer) Reset() {
	b.data = b.data[:0]
}

type zeroCopyParser struct {
	parser *Parser
}

func newZeroCopyParser() *zeroCopyParser {
	return &zeroCopyParser{
		parser: NewParser(),
	}
}

func (p *zeroCopyParser) Parse(data []byte) (*Message, error) {
	if len(data) == 0 {
		return nil, &ParseError{Message: "empty data"}
	}

	msg := GetMessage()

	data = bytes.ReplaceAll(data, []byte{'\r', '\n'}, []byte{'\r'})
	data = bytes.ReplaceAll(data, []byte{'\n'}, []byte{'\r'})
	data = bytes.TrimSuffix(data, []byte{'\r'})

	segments := bytes.Split(data, []byte{'\r'})
	if len(segments) > p.parser.config.MaxSegments {
		PutMessage(msg)
		return nil, &ParseError{Message: "too many segments"}
	}

	for _, segData := range segments {
		if len(segData) == 0 {
			continue
		}
		segData = bytes.TrimSpace(segData)
		if len(segData) == 0 {
			continue
		}
		seg := p.parseSegmentZeroCopy(segData)
		msg.segments = append(msg.segments, seg)
	}

	return msg, nil
}

func (p *zeroCopyParser) parseSegmentZeroCopy(data []byte) Segment {
	seg := Segment{}

	if len(data) == 0 {
		return seg
	}

	delims := DefaultDelimiters()
	segFieldSep := delimByIndex(data, 0)

	if len(data) >= 3 {
		seg.name = bytesToString(data[:3])
	} else {
		seg.name = bytesToString(data)
		return seg
	}

	seg.fields = seg.fields[:0]
	if seg.name == "MSH" {
		seg.fields = append(seg.fields, string(rune(segFieldSep)))

		fields := splitBytes(data[3:], byte(delims.Field))
		for _, field := range fields {
			seg.fields = append(seg.fields, bytesToString(field))
		}
	} else {
		seg.fields = seg.fields[:0]
		fields := splitBytes(data[3:], byte(delims.Field))
		for _, field := range fields {
			seg.fields = append(seg.fields, bytesToString(field))
		}
	}

	return seg
}

func delimByIndex(data []byte, idx int) byte {
	if idx < len(data) {
		return data[idx]
	}
	return 0
}

func splitBytes(data []byte, sep byte) [][]byte {
	if len(data) == 0 {
		return nil
	}

	var result [][]byte
	start := 0

	for i := 0; i < len(data); i++ {
		if data[i] == sep {
			result = append(result, data[start:i])
			start = i + 1
		}
	}

	result = append(result, data[start:])
	return result
}

func splitStringZeroCopy(s string, sep rune) []string {
	if s == "" {
		return nil
	}

	var result []string
	start := 0
	chars := []rune(s)

	for i, r := range chars {
		if r == sep {
			result = append(result, s[start:i])
			start = i + 1
		}
	}

	result = append(result, s[start:])
	return result
}

func NormalizeLineEndings(data []byte) []byte {
	data = bytes.ReplaceAll(data, []byte{'\r', '\n'}, []byte{'\r'})
	data = bytes.ReplaceAll(data, []byte{'\n'}, []byte{'\r'})
	return bytes.TrimSuffix(data, []byte{'\r'})
}

func BuildSegmentString(name string, fields []string, delims Delimiters) string {
	if name == "" {
		return ""
	}

	var buf strings.Builder
	buf.WriteString(name)

	if name == "MSH" {
		buf.WriteRune(delims.Field)
		for i, field := range fields {
			if i > 0 {
				buf.WriteRune(delims.Field)
			}
			buf.WriteString(field)
		}
	} else {
		for _, field := range fields {
			buf.WriteRune(delims.Field)
			buf.WriteString(field)
		}
	}

	return buf.String()
}

func BuildSegmentBytes(name string, fields []string, delims Delimiters) []byte {
	if name == "" {
		return nil
	}

	buf := newBytesBuffer(64)
	buf.WriteString(name)

	if name == "MSH" {
		buf.WriteByte(byte(delims.Field))
		for i, field := range fields {
			if i > 0 {
				buf.WriteByte(byte(delims.Field))
			}
			buf.WriteString(field)
		}
	} else {
		for _, field := range fields {
			buf.WriteByte(byte(delims.Field))
			buf.WriteString(field)
		}
	}

	return buf.Bytes()
}
