package hl7

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

var (
	ErrInvalidUTF8 = errors.New("invalid UTF-8 sequence")
	ErrInvalidHL7  = errors.New("invalid HL7 escape sequence")
)

type Encoding string

const (
	EncodingASCII       Encoding = "ASCII"
	EncodingUTF8        Encoding = "UTF-8"
	EncodingUTF16       Encoding = "UTF-16"
	EncodingUTF16LE     Encoding = "UTF-16LE"
	EncodingUTF16BE     Encoding = "UTF-16BE"
	EncodingISO88591    Encoding = "ISO-8859-1"
	EncodingWindows1252 Encoding = "Windows-1252"
	EncodingLatin1      Encoding = "latin1"
)

var defaultEncoding = EncodingASCII

func SetDefaultEncoding(enc Encoding) {
	defaultEncoding = enc
}

func GetDefaultEncoding() Encoding {
	return defaultEncoding
}

type EscapeHandler struct {
	escapeChar rune
}

func NewEscapeHandler() *EscapeHandler {
	return &EscapeHandler{
		escapeChar: '\\',
	}
}

func (h *EscapeHandler) Encode(data string) string {
	var result strings.Builder
	for _, r := range data {
		switch r {
		case '|':
			result.WriteRune(h.escapeChar)
			result.WriteRune('S')
		case '~':
			result.WriteRune(h.escapeChar)
			result.WriteRune('R')
			result.WriteRune('~')
		case '^':
			result.WriteRune(h.escapeChar)
			result.WriteRune('R')
			result.WriteRune('^')
		case '&':
			result.WriteRune(h.escapeChar)
			result.WriteRune('R')
			result.WriteRune('&')
		case '\\':
			result.WriteRune(h.escapeChar)
			result.WriteRune('\\')
		case '\x0B':
			result.WriteRune(h.escapeChar)
			result.WriteRune('H')
		case '\x09':
			result.WriteRune(h.escapeChar)
			result.WriteRune('T')
		case '\x0D':
			result.WriteRune(h.escapeChar)
			result.WriteRune('N')
		case '\x0A':
			result.WriteRune(h.escapeChar)
			result.WriteRune('R')
		case '\x0C':
			result.WriteRune(h.escapeChar)
			result.WriteRune('F')
		case '\x1B':
			result.WriteRune(h.escapeChar)
			result.WriteRune('E')
		default:
			if r < 0x20 || r > 0x7E {
				result.WriteRune(h.escapeChar)
				result.WriteRune('X')
				result.WriteString(fmt.Sprintf("%02X", r))
			} else {
				result.WriteRune(r)
			}
		}
	}
	return result.String()
}

func (h *EscapeHandler) Decode(data string) (string, error) {
	var result strings.Builder
	runes := []rune(data)
	i := 0

	for i < len(runes) {
		if runes[i] != h.escapeChar {
			result.WriteRune(runes[i])
			i++
			continue
		}

		if i+1 >= len(runes) {
			return "", fmt.Errorf("%w: escape at end of string", ErrInvalidHL7)
		}

		next := runes[i+1]
		switch next {
		case 'H':
			result.WriteRune(0x0B)
			i += 2
		case 'N':
			result.WriteRune(0x0D)
			i += 2
		case 'R':
			result.WriteRune(0x0A)
			i += 2
		case 'S':
			result.WriteRune('|')
			i += 2
		case 'T':
			result.WriteRune(0x09)
			i += 2
		case 'E':
			result.WriteRune(0x1B)
			i += 2
		case '\\':
			result.WriteRune('\\')
			i += 2
		case 'F':
			result.WriteRune(0x0C)
			i += 2
		case 'C':
			result.WriteRune(0x0D)
			i += 2
		case 'X':
			if i+3 >= len(runes) {
				return "", fmt.Errorf("%w: incomplete hex escape", ErrInvalidHL7)
			}
			hex := string(runes[i+2]) + string(runes[i+3])
			b, err := parseHex(hex)
			if err != nil {
				return "", err
			}
			result.WriteRune(rune(b))
			i += 4
		case 'Z':
			start := i + 2
			found := false
			for j := start; j < len(runes); j++ {
				if runes[j] == h.escapeChar {
					i = j + 1
					found = true
					break
				}
			}
			if !found {
				i = len(runes)
			}
		default:
			result.WriteRune(next)
			i += 2
		}
	}

	return result.String(), nil
}

func parseHex(s string) (byte, error) {
	var b byte
	for _, c := range s {
		var val byte
		switch {
		case c >= '0' && c <= '9':
			val = byte(c - '0')
		case c >= 'A' && c <= 'F':
			val = byte(c - 'A' + 10)
		case c >= 'a' && c <= 'f':
			val = byte(c - 'a' + 10)
		default:
			return 0, fmt.Errorf("%w: invalid hex digit %c", ErrInvalidHL7, c)
		}
		b = b<<4 | val
	}
	return b, nil
}

func ConvertToUTF8(data []byte, srcEncoding Encoding) ([]byte, error) {
	switch srcEncoding {
	case EncodingUTF8:
		if !utf8.Valid(data) {
			return nil, ErrInvalidUTF8
		}
		return data, nil
	case EncodingASCII, EncodingISO88591, EncodingLatin1:
		return []byte(iso88591ToUTF8(string(data))), nil
	case EncodingWindows1252:
		return []byte(windows1252ToUTF8(string(data))), nil
	case EncodingUTF16, EncodingUTF16LE:
		return utf16LEToUTF8(data)
	case EncodingUTF16BE:
		return utf16BEToUTF8(data)
	default:
		return data, nil
	}
}

func ConvertFromUTF8(data []byte, dstEncoding Encoding) ([]byte, error) {
	if !utf8.Valid(data) {
		return nil, ErrInvalidUTF8
	}

	switch dstEncoding {
	case EncodingUTF8:
		return data, nil
	case EncodingASCII:
		return []byte(utf8ToASCII(string(data))), nil
	case EncodingISO88591, EncodingLatin1:
		return []byte(utf8ToISO88591(string(data))), nil
	default:
		return data, nil
	}
}

func iso88591ToUTF8(s string) string {
	runes := make([]rune, len(s))
	for i, b := range []byte(s) {
		runes[i] = rune(b)
	}
	return string(runes)
}

func windows1252ToUTF8(s string) string {
	var result []rune
	for _, r := range []byte(s) {
		switch r {
		case 0x80:
			result = append(result, 0x20AC)
		case 0x81:
			result = append(result, 0x0081)
		case 0x82:
			result = append(result, 0x201A)
		case 0x83:
			result = append(result, 0x0192)
		case 0x84:
			result = append(result, 0x201E)
		case 0x85:
			result = append(result, 0x2026)
		case 0x86:
			result = append(result, 0x2020)
		case 0x87:
			result = append(result, 0x2021)
		case 0x88:
			result = append(result, 0x02C6)
		case 0x89:
			result = append(result, 0x2030)
		case 0x8A:
			result = append(result, 0x0160)
		case 0x8B:
			result = append(result, 0x2039)
		case 0x8C:
			result = append(result, 0x0152)
		case 0x8D:
			result = append(result, 0x008D)
		case 0x8E:
			result = append(result, 0x017D)
		case 0x8F:
			result = append(result, 0x008F)
		case 0x90:
			result = append(result, 0x0090)
		case 0x91:
			result = append(result, 0x2018)
		case 0x92:
			result = append(result, 0x2019)
		case 0x93:
			result = append(result, 0x201C)
		case 0x94:
			result = append(result, 0x201D)
		case 0x95:
			result = append(result, 0x2022)
		case 0x96:
			result = append(result, 0x2013)
		case 0x97:
			result = append(result, 0x2014)
		case 0x98:
			result = append(result, 0x02DC)
		case 0x99:
			result = append(result, 0x2122)
		case 0x9A:
			result = append(result, 0x0161)
		case 0x9B:
			result = append(result, 0x203A)
		case 0x9C:
			result = append(result, 0x0153)
		case 0x9D:
			result = append(result, 0x009D)
		case 0x9E:
			result = append(result, 0x017E)
		case 0x9F:
			result = append(result, 0x0178)
		default:
			result = append(result, rune(r))
		}
	}
	return string(result)
}

func utf16LEToUTF8(data []byte) ([]byte, error) {
	if len(data)%2 != 0 {
		return nil, fmt.Errorf("invalid UTF-16 data length")
	}
	u16 := make([]uint16, len(data)/2)
	for i := 0; i < len(u16); i++ {
		u16[i] = uint16(data[i*2]) | (uint16(data[i*2+1]) << 8)
	}
	return []byte(string(utf16.Decode(u16))), nil
}

func utf16BEToUTF8(data []byte) ([]byte, error) {
	if len(data)%2 != 0 {
		return nil, fmt.Errorf("invalid UTF-16 data length")
	}
	u16 := make([]uint16, len(data)/2)
	for i := 0; i < len(u16); i++ {
		u16[i] = (uint16(data[i*2]) << 8) | uint16(data[i*2+1])
	}
	return []byte(string(utf16.Decode(u16))), nil
}

func utf8ToASCII(s string) string {
	var result strings.Builder
	for _, r := range s {
		if r < 128 {
			result.WriteRune(r)
		} else {
			result.WriteRune('?')
		}
	}
	return result.String()
}

func utf8ToISO88591(s string) string {
	var result strings.Builder
	for _, r := range s {
		if r < 256 {
			result.WriteRune(r)
		} else {
			result.WriteRune('?')
		}
	}
	return result.String()
}

var escapeHandler = NewEscapeHandler()

type Decoder struct {
	encoding       Encoding
	escapeHandling bool
}

func NewDecoder() *Decoder {
	return &Decoder{
		encoding:       defaultEncoding,
		escapeHandling: true,
	}
}

func (d *Decoder) WithEncoding(enc Encoding) *Decoder {
	d.encoding = enc
	return d
}

func (d *Decoder) WithEscapeHandling(enabled bool) *Decoder {
	d.escapeHandling = enabled
	return d
}

func (d *Decoder) Decode(data []byte) ([]byte, error) {
	if d.escapeHandling {
		str := string(data)
		decoded, err := escapeHandler.Decode(str)
		if err != nil {
			return nil, err
		}
		data = []byte(decoded)
	}

	return ConvertToUTF8(data, d.encoding)
}

func DecodeMessage(data []byte, enc Encoding) ([]byte, error) {
	decoder := NewDecoder().WithEncoding(enc)
	return decoder.Decode(data)
}

func EncodeMessage(data []byte, enc Encoding) ([]byte, error) {
	return ConvertFromUTF8(data, enc)
}

func Escape(data string) string {
	return escapeHandler.Encode(data)
}

func Unescape(data string) (string, error) {
	return escapeHandler.Decode(data)
}

func (e *Encoder) WithEncoding(enc Encoding) *Encoder {
	e.encoding = enc
	return e
}

func (e *Encoder) WithEscapeHandling(enabled bool) *Encoder {
	e.escapeHandling = enabled
	return e
}

func (e *Encoder) EncodeWithOptions(msg *Message, enc Encoding) ([]byte, error) {
	data, err := e.Encode(msg)
	if err != nil {
		return nil, err
	}
	return ConvertFromUTF8(data, enc)
}
