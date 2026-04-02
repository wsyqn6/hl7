# HL7

[![Go Reference](https://pkg.go.dev/badge/github.com/wsyqn6/hl7.svg)](https://pkg.go.dev/github.com/wsyqn6/hl7)
[![Go Report Card](https://goreportcard.com/badge/github.com/wsyqn6/hl7)](https://goreportcard.com/report/github.com/wsyqn6/hl7)

A comprehensive Go library for parsing, encoding, and manipulating HL7 v2.x healthcare messages.

## Features

- **Struct Marshaling**: Map HL7 data to Go structs using tags
- **Streaming Parser**: Memory-efficient parsing for large messages
- **MLLP Support**: Client/server implementation for HL7 transport
- **Message Validation**: Rule-based validation with built-in and custom rules
- **ACK Generation**: Automatic acknowledgment message creation
- **Escape Sequence Handling**: Full support for HL7 escape sequences

## Installation

```bash
go get github.com/wsyqn6/hl7
```

Requires Go 1.26 or later.

## Quick Start

### Parsing a Message

```go
package main

import (
    "fmt"
    "log"

    "github.com/yourusername/hl7"
)

func main() {
    data := []byte(`MSH|^~\&|SENDING|FACILITY|||202401151200||ADT^A01|MSG001|P|2.5
PID|1||12345^^^MRN||Smith^John^A||19800115|M`)

    msg, err := hl7.Parse(data)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Message Type: %s\n", msg.Type())
}
```

### Struct Marshaling

```go
type Patient struct {
    ID        string `hl7:"PID.3.1"`
    LastName  string `hl7:"PID.5.1"`
    FirstName string `hl7:"PID.5.2"`
    DOB       string `hl7:"PID.7"`
    Gender    string `hl7:"PID.8"`
}

var patient Patient
err := hl7.Unmarshal(data, &patient)
```

## Documentation

See [pkg.go.dev/github.com/wsyqn6/hl7](https://pkg.go.dev/github.com/wsyqn6/hl7) for full documentation.

## License

MIT License - see [LICENSE](LICENSE) for details.
