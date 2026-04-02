# hl7 🏥

> A comprehensive Go library for parsing, encoding, and manipulating HL7 v2.x healthcare messages.

[English](#english) | [中文](#中文)

---

## English

### Features

- 🔍 **Full HL7 v2.x Support** - Parse and encode all HL7 v2.x message types (2.3, 2.4, 2.5, etc.)
- ⚡ **Struct Marshaling** - Map HL7 data to Go structs using intuitive `hl7` tags
- 🌊 **Streaming Parser** - Memory-efficient parsing with MLLP frame support
- 🌐 **MLLP Network Transport** - Built-in client/server for HL7 over TCP
- ✅ **Message Validation** - Flexible validation with built-in and custom rules
- 📝 **ACK/NAK Generation** - Automatic acknowledgment message creation
- 🔄 **Bidirectional Conversion** - Parse HL7 to structs, marshal structs back to HL7

### Installation

```bash
go get github.com/wsyqn6/hl7
```

Requires Go 1.26 or later.

### Quick Start

#### Parsing a Message

```go
package main

import (
    "fmt"
    "log"

    "github.com/wsyqn6/hl7"
)

func main() {
    data := []byte(`MSH|^~\&|SENDING|FACILITY|||202401151200||ADT^A01|MSG001|P|2.5
PID|1||12345^^^MRN||Smith^John^A||19800115|M`)

    msg, err := hl7.Parse(data)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Message Type: %s\n", msg.Type())
    fmt.Printf("Control ID: %s\n", msg.ControlID())

    // Access PID segment
    if pid, ok := msg.Segment("PID"); ok {
        fmt.Printf("Patient: %s %s\n", pid.Component(5, 2), pid.Component(5, 1))
        fmt.Printf("DOB: %s, Gender: %s\n", pid.Field(7), pid.Field(8))
    }
}
```

#### Struct Marshaling

```go
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

fmt.Printf("Patient: %s, %s (MRN: %s)\n", patient.LastName, patient.FirstName, patient.MRN)
```

#### Validation

```go
validator := hl7.NewValidator(
    hl7.Required("MSH.9"),
    hl7.Required("PID.3.1"),
    hl7.OneOf("PID.8", "M", "F", "O", "U"),
    hl7.Pattern("PID.7", `^\d{8}$`),
)

if errors := validator.Validate(msg); len(errors) > 0 {
    for _, err := range errors {
        fmt.Printf("Validation error at %s: %s\n", err.Location, err.Message)
    }
}
```

#### MLLP Server

```go
handler := func(ctx context.Context, msg *hl7.Message) (*hl7.Message, error) {
    fmt.Printf("Received: %s\n", msg.Type())
    return hl7.Generate(msg, hl7.Accept())
}

server := hl7.NewServer(":2575", handler)
if err := server.ListenAndServe(); err != nil {
    log.Fatal(err)
}
```

#### MLLP Client

```go
client, err := hl7.Dial("localhost:2575")
if err != nil {
    log.Fatal(err)
}
defer client.Close()

ack, err := client.Send(ctx, msg)
if err != nil {
    log.Fatal(err)
}
```

### API Reference

| Function | Description |
|----------|-------------|
| `Parse(data []byte)` | Parse raw HL7 data into a Message |
| `ParseString(s string)` | Parse HL7 from a string |
| `Encode(msg *Message)` | Encode Message to bytes |
| `Unmarshal(data []byte, v interface{})` | Unmarshal HL7 into a struct |
| `Marshal(v interface{})` | Marshal a struct to HL7 bytes |
| `Generate(msg *Message, opts ...ACKOption)` | Generate ACK/NAK response |
| `NewValidator(rules ...Rule)` | Create a message validator |
| `NewServer(addr string, handler Handler)` | Create MLLP server |
| `Dial(addr string)` | Create MLLP client |

### License

MIT License - see [LICENSE](LICENSE) for details.

---

## 中文

### 功能特性

- 🔍 **完整的 HL7 v2.x 支持** - 解析和编码所有 HL7 v2.x 消息类型（2.3、2.4、2.5 等）
- ⚡ **结构体序列化** - 通过直观的 `hl7` 标签将 HL7 数据映射到 Go 结构体
- 🌊 **流式解析器** - 高效内存使用的流式解析，支持 MLLP 帧
- 🌐 **MLLP 网络传输** - 内置的 HL7 over TCP 客户端/服务器
- ✅ **消息验证** - 灵活的验证规则，支持内置和自定义规则
- 📝 **ACK/NAK 生成** - 自动创建确认和拒绝消息
- 🔄 **双向转换** - 解析 HL7 为结构体，将结构体编组回 HL7

### 安装

```bash
go get github.com/wsyqn6/hl7
```

需要 Go 1.26 或更高版本。

### 快速开始

#### 解析消息

```go
package main

import (
    "fmt"
    "log"

    "github.com/wsyqn6/hl7"
)

func main() {
    data := []byte(`MSH|^~\&|发送系统|医院|||202401151200||ADT^A01|MSG001|P|2.5
PID|1||12345^^^MRN||张三^李||19800115|M`)

    msg, err := hl7.Parse(data)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("消息类型: %s\n", msg.Type())
    fmt.Printf("控制ID: %s\n", msg.ControlID())

    // 访问 PID 段
    if pid, ok := msg.Segment("PID"); ok {
        fmt.Printf("患者: %s %s\n", pid.Component(5, 2), pid.Component(5, 1))
        fmt.Printf("出生日期: %s, 性别: %s\n", pid.Field(7), pid.Field(8))
    }
}
```

#### 结构体序列化

```go
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

fmt.Printf("患者: %s, %s (病历号: %s)\n", patient.LastName, patient.FirstName, patient.MRN)
```

#### 消息验证

```go
validator := hl7.NewValidator(
    hl7.Required("MSH.9"),      // 消息类型必填
    hl7.Required("PID.3.1"),    // 患者ID必填
    hl7.OneOf("PID.8", "M", "F", "O", "U"),  // 有效的性别代码
    hl7.Pattern("PID.7", `^\d{8}$`),  // 出生日期格式 YYYYMMDD
)

if errors := validator.Validate(msg); len(errors) > 0 {
    for _, err := range errors {
        fmt.Printf("验证错误 [%s]: %s\n", err.Location, err.Message)
    }
}
```

#### MLLP 服务器

```go
handler := func(ctx context.Context, msg *hl7.Message) (*hl7.Message, error) {
    fmt.Printf("收到消息: %s\n", msg.Type())
    return hl7.Generate(msg, hl7.Accept())
}

server := hl7.NewServer(":2575", handler)
if err := server.ListenAndServe(); err != nil {
    log.Fatal(err)
}
```

#### MLLP 客户端

```go
client, err := hl7.Dial("localhost:2575")
if err != nil {
    log.Fatal(err)
}
defer client.Close()

ack, err := client.Send(ctx, msg)
if err != nil {
    log.Fatal(err)
}
```

### API 参考

| 函数 | 描述 |
|------|------|
| `Parse(data []byte)` | 将原始 HL7 数据解析为 Message |
| `ParseString(s string)` | 从字符串解析 HL7 |
| `Encode(msg *Message)` | 将 Message 编码为字节 |
| `Unmarshal(data []byte, v interface{})` | 将 HL7 反序列化为结构体 |
| `Marshal(v interface{})` | 将结构体编组为 HL7 字节 |
| `Generate(msg *Message, opts ...ACKOption)` | 生成 ACK/NAK 响应 |
| `NewValidator(rules ...Rule)` | 创建消息验证器 |
| `NewServer(addr string, handler Handler)` | 创建 MLLP 服务器 |
| `Dial(addr string)` | 创建 MLLP 客户端 |

### 许可证

MIT 许可证 - 详见 [LICENSE](LICENSE)。

---

<p align="center">
  <a href="https://pkg.go.dev/github.com/wsyqn6/hl7">📚 GoDoc</a> •
  <a href="https://github.com/wsyqn6/hl7/issues">🐛 Issues</a> •
  <a href="https://hl7.org">🌐 HL7 Official</a>
</p>
