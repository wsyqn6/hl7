# hl7 🏥

> A comprehensive Go library for parsing, encoding, and manipulating HL7 v2.x healthcare messages.

[English](#english) | [中文](#中文)

[![Go Version](https://img.shields.io/github/go-mod/go-version/wsyqn6/hl7)](https://github.com/wsyqn6/hl7)
[![License](https://img.shields.io/github/license/wsyqn6/hl7)](https://github.com/wsyqn6/hl7)
[![Test](https://github.com/wsyqn6/hl7/actions/workflows/ci.yml/badge.svg)](https://github.com/wsyqn6/hl7/actions)
[![Coverage](https://codecov.io/gh/wsyqn6/hl7/branch/main/graph/badge.svg)](https://codecov.io/gh/wsyqn6/hl7)

---

## English

### Features

- 🔍 **Full HL7 v2.x Support** - Parse and encode all HL7 v2.x message types (2.3, 2.4, 2.5, etc.)
- ⚡ **Struct Marshaling** - Map HL7 data to Go structs using intuitive `hl7` tags
- 🌊 **Streaming Parser** - Memory-efficient parsing with MLLP frame support
- 🌐 **MLLP Network Transport** - Built-in client/server for HL7 over TCP with TLS support
- ✅ **Message Validation** - Flexible validation with built-in and custom rules
- 📝 **ACK/NAK Generation** - Automatic acknowledgment message creation
- 🔄 **Bidirectional Conversion** - Parse HL7 to structs, marshal structs back to HL7
- 📊 **Schema-less Parsing** - Access fields without predefined structs using `Get()` method
- 🔧 **Segment Helpers** - Convenient methods for common segments (PID, MSH, PV1, OBR, OBX, NK1, DG1)

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

#### Advanced Struct Mapping

**Optional Fields:**

```go
type Patient struct {
    MRN       string `hl7:"PID.3.1"`
    LastName  string `hl7:"PID.5.1"`
    Phone     string `hl7:"PID.13,optional"`  // Optional - won't error if missing
}
```

**Nested Struct with Components:**

```go
type PersonName struct {
    FamilyName  string `hl7:"1"`
    GivenName   string `hl7:"2"`
    MiddleName  string `hl7:"3"`
}

type Patient struct {
    MRN     string     `hl7:"PID.3.1"`
    Name    PersonName `hl7:"PID.5"`
}
```

#### Schema-less Parsing

Access fields without predefined structs:

```go
// Simple field access
lastName, _ := msg.Get("PID.5.1")
firstName, _ := msg.Get("PID.5.2")

// Component access
patientID := msg.MustGet("PID.3")      // Full field
mrn := msg.MustGet("PID.3.1")          // First component

// Repeated segments
wbc := msg.MustGet("OBX[1].5")         // First OBX observation value
rbc := msg.MustGet("OBX[2].5")         // Second OBX observation value
```

#### Segment Helpers

Convenient methods for common segments:

```go
// PID segment
pid := msg.PID()
fmt.Printf("Patient: %s, %s\n", pid.LastName(), pid.FirstName())
fmt.Printf("DOB: %s, Gender: %s\n", pid.DateOfBirth(), pid.Gender())

// MSH segment
msh := msg.MSH()
fmt.Printf("From: %s@%s\n", msh.SendingApplication(), msh.SendingFacility())

// OBX results (repeated)
for _, obx := range msg.AllOBX() {
    fmt.Printf("%s: %s %s\n", obx.ObservationIdentifierCode(), obx.ObservationValue(), obx.Units())
}
```

#### Message Type Examples

**ADT (Admission/Discharge/Transfer) Message:**

```go
// ADT^A01 - Patient Admission
adtMsg := []byte(`MSH|^~\&|ADT_SYS|HOSPITAL|||^202401151200||ADT^A01|CTRL|P|2.5
PID|1||12345^^^MRN||Smith^John^A||19800115|M|||123 Main St^^Springfield^IL^62701||555-1234
PV1|1|I|ICU^Room101^BED1^^^SEC^ICU|||DR001^Dr. Smith^John^MD|||ICU|||||||||ADM`)

msg, _ := hl7.Parse(adtMsg)
pid := msg.PID()
fmt.Printf("Patient: %s %s\n", pid.FirstName(), pid.LastName())
fmt.Printf("Admission Type: %s\n", msg.PV1().AdmissionType())
```

**ORU (Observation Result) Message:**

```go
// ORU^R01 - Laboratory Results
oruMsg := []byte(`MSH|^~\&|LAB_SYS|HOSPITAL|||^202401151200||ORU^R01|CTRL|P|2.5
PID|1||12345^^^MRN||Smith^John^A||19800115|M
OBR|1||12345|LAB001^CBC^L|||202401151000|||LAB||||||DR001
OBX|1|NM|WBC^WBC^LN||7.5|x10^3/uL^3^1^ML||4.5-11.0|N|||F
OBX|2|NM|RBC^RBC^LN||4.8|x10^6/uL^3^2^ML||4.5-5.5|N|||F
OBX|3|NM|HGB^Hemoglobin^LN||14.2|g/dL^3^3^ML||12.0-17.0|N|||F`)

msg, _ := hl7.Parse(oruMsg)
for _, obx := range msg.AllOBX() {
    fmt.Printf("%s: %s %s\n", obx.ObservationIdentifierCode(), obx.ObservationValue(), obx.Units())
}
```

**ORM (Order) Message:**

```go
// ORM^O01 - General Order
ormMsg := []byte(`MSH|^~\&|ORD_SYS|HOSPITAL|||^202401151200||ORM^O01|CTRL|P|2.5
PID|1||12345^^^MRN||Smith^John^A||19800115|M
ORC|RE|12345|ORD001|CBC^LAB^L||1|||DR001^Dr. Smith^John^MD|||202401151200
OBR|1||12345|ORD001|CBC^LAB^L|||202401151200|||LAB||||||DR001`)

msg, _ := hl7.Parse(ormMsg)
orc := msg.ORC()
fmt.Printf("Order Control: %s\n", orc.OrderStatus())
```

#### TLS MLLP Server

Secure MLLP server with TLS:

```go
handler := func(ctx context.Context, msg *hl7.Message) (*hl7.Message, error) {
    return hl7.Generate(msg, hl7.Accept())
}

server := hl7.NewServer(":2575", handler,
    hl7.WithTLS(nil),  // Uses cert/key files
)
if err := server.ListenAndServeTLS("cert.pem", "key.pem"); err != nil {
    log.Fatal(err)
}
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

#### MLLP Client with Retry

```go
client, err := hl7.Dial("localhost:2575",
    hl7.WithRetry(3, 100*time.Millisecond),  // 3 retries with 100ms initial delay
)
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// SendWithRetry automatically retries on failure
ack, err := client.SendWithRetry(ctx, msg)
```

#### MLLP Client Pool

```go
pool := hl7.NewClientPool("localhost:2575",
    hl7.WithPoolSize(10),  // Max 10 clients in pool
    hl7.WithRetry(2, 50*time.Millisecond),
)
defer pool.Close()

// Get client from pool, send, return to pool
ack, err := pool.Send(ctx, msg)
```

### API Reference

| Function | Description |
|----------|-------------|
| `Parse(data []byte)` | Parse raw HL7 data into a Message |
| `ParseString(s string)` | Parse HL7 from a string |
| `Encode(msg *Message)` | Encode Message to bytes |
| `Unmarshal(data []byte, v interface{})` | Unmarshal HL7 into a struct |
| `Marshal(v interface{})` | Marshal a struct to HL7 bytes |
| `msg.Get(location)` | Get field value by location (e.g., "PID.5.1") |
| `msg.MustGet(location)` | Get field value, panic on error |
| `msg.GetAllRepetitions(location)` | Get all repetitions of a field |
| `msg.CountSegment(name)` | Count segments by name |
| `msg.HasSegment(name)` | Check if segment exists |
| `msg.ParseLocation(location)` | Parse location string to Location struct |
| `msg.Iterate()` | Iterate over all segments |
| `msg.Stats()` | Get message statistics |
| `seg.Repetitions(fieldIdx)` | Get all repetitions of a field |
| `seg.Components(fieldIdx)` | Get all components of a field |
| `Generate(msg *Message, opts ...ACKOption)` | Generate ACK/NAK response |
| `NewValidator(rules ...Rule)` | Create a message validator |
| `NewServer(addr string, handler Handler)` | Create MLLP server |
| `Dial(addr string)` | Create MLLP client |
| `DialTLS(addr string, config)` | Create MLLP client with TLS |
| `NewClientPool(addr)` | Create MLLP client pool |

### Segment Helpers

| Helper | Description |
|--------|-------------|
| `msg.PID()` | Get PID segment helper |
| `msg.MSH()` | Get MSH segment helper |
| `msg.PV1()` | Get PV1 segment helper |
| `msg.OBR()` | Get OBR segment helper |
| `msg.OBX()` | Get first OBX segment helper |
| `msg.AllOBX()` | Get all OBX segments |
| `msg.NK1()` | Get NK1 segment helper |
| `msg.DG1()` | Get DG1 segment helper |

### License

MIT License - see [LICENSE](LICENSE) for details.

---

## 中文

### 功能特性

- 🔍 **完整的 HL7 v2.x 支持** - 解析和编码所有 HL7 v2.x 消息类型（2.3、2.4、2.5 等）
- ⚡ **结构体序列化** - 通过直观的 `hl7` 标签将 HL7 数据映射到 Go 结构体
- 🌊 **流式解析器** - 高效内存使用的流式解析，支持 MLLP 帧
- 🌐 **MLLP 网络传输** - 内置的 HL7 over TCP 客户端/服务器，支持 TLS
- ✅ **消息验证** - 灵活的验证规则，支持内置和自定义规则
- 📝 **ACK/NAK 生成** - 自动创建确认和拒绝消息
- 🔄 **双向转换** - 解析 HL7 为结构体，将结构体编组回 HL7
- 📊 **无模式解析** - 使用 `Get()` 方法无需预定义结构体即可访问字段
- 🔧 **段帮助函数** - 常用段（PID、MSH、PV1、OBR、OBX、NK1、DG1）的便捷方法

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

#### 无模式解析

无需预定义结构体即可访问字段：

```go
// 简单字段访问
lastName, _ := msg.Get("PID.5.1")
firstName, _ := msg.Get("PID.5.2")

// 组件访问
patientID := msg.MustGet("PID.3")     // 完整字段
mrn := msg.MustGet("PID.3.1")          // 第一组件

// 重复段
wbc := msg.MustGet("OBX[1].5")        // 第一个 OBX 观察值
rbc := msg.MustGet("OBX[2].5")        // 第二个 OBX 观察值
```

#### 段帮助函数

常用段的便捷方法：

```go
// PID 段
pid := msg.PID()
fmt.Printf("患者: %s, %s\n", pid.LastName(), pid.FirstName())
fmt.Printf("出生日期: %s, 性别: %s\n", pid.DateOfBirth(), pid.Gender())

// MSH 段
msh := msg.MSH()
fmt.Printf("来自: %s@%s\n", msh.SendingApplication(), msh.SendingFacility())

// OBX 结果（重复段）
for _, obx := range msg.AllOBX() {
    fmt.Printf("%s: %s %s\n", obx.ObservationIdentifierCode(), obx.ObservationValue(), obx.Units())
}
```

#### TLS MLLP 服务器

支持 TLS 的安全 MLLP 服务器：

```go
handler := func(ctx context.Context, msg *hl7.Message) (*hl7.Message, error) {
    return hl7.Generate(msg, hl7.Accept())
}

server := hl7.NewServer(":2575", handler,
    hl7.WithTLS(nil),  // 使用证书/密钥文件
)
if err := server.ListenAndServeTLS("cert.pem", "key.pem"); err != nil {
    log.Fatal(err)
}
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

#### MLLP 客户端（重试机制）

```go
client, err := hl7.Dial("localhost:2575",
    hl7.WithRetry(3, 100*time.Millisecond),  // 3次重试，初始延迟100ms
)
defer client.Close()

// 自动重试发送
ack, err := client.SendWithRetry(ctx, msg)
```

#### MLLP 连接池

```go
pool := hl7.NewClientPool("localhost:2575",
    hl7.WithPoolSize(10),  // 最多10个客户端
    hl7.WithRetry(2, 50*time.Millisecond),
)
defer pool.Close()

// 从池中获取客户端发送
ack, err := pool.Send(ctx, msg)
```

### API 参考

| 函数 | 描述 |
|------|------|
| `Parse(data []byte)` | 将原始 HL7 数据解析为 Message |
| `ParseString(s string)` | 从字符串解析 HL7 |
| `Encode(msg *Message)` | 将 Message 编码为字节 |
| `Unmarshal(data []byte, v interface{})` | 将 HL7 反序列化为结构体 |
| `Marshal(v interface{})` | 将结构体编组为 HL7 字节 |
| `msg.Get(location)` | 按位置获取字段值（如 "PID.5.1"）|
| `msg.MustGet(location)` | 获取字段值，错误时 panic |
| `msg.GetAllRepetitions(location)` | 获取字段的所有重复值 |
| `msg.CountSegment(name)` | 按名称统计段数量 |
| `msg.HasSegment(name)` | 检查段是否存在 |
| `msg.ParseLocation(location)` | 解析位置字符串为 Location 结构体 |
| `msg.Iterate()` | 遍历所有段 |
| `msg.Stats()` | 获取消息统计信息 |
| `seg.Repetitions(fieldIdx)` | 获取字段的所有重复值 |
| `seg.Components(fieldIdx)` | 获取字段的所有组件 |
| `Generate(msg *Message, opts ...ACKOption)` | 生成 ACK/NAK 响应 |
| `NewValidator(rules ...Rule)` | 创建消息验证器 |
| `NewServer(addr string, handler Handler)` | 创建 MLLP 服务器 |
| `Dial(addr string)` | 创建 MLLP 客户端 |
| `DialTLS(addr string, config)` | 创建 TLS MLLP 客户端 |
| `NewClientPool(addr)` | 创建 MLLP 客户端连接池 |

### 段帮助函数

| 帮助函数 | 描述 |
|----------|------|
| `msg.PID()` | 获取 PID 段帮助器 |
| `msg.MSH()` | 获取 MSH 段帮助器 |
| `msg.PV1()` | 获取 PV1 段帮助器 |
| `msg.OBR()` | 获取 OBR 段帮助器 |
| `msg.OBX()` | 获取第一个 OBX 段帮助器 |
| `msg.AllOBX()` | 获取所有 OBX 段 |
| `msg.NK1()` | 获取 NK1 段帮助器 |
| `msg.DG1()` | 获取 DG1 段帮助器 |

### 许可证

MIT 许可证 - 详见 [LICENSE](LICENSE)。

---

<p align="center">
  <a href="https://pkg.go.dev/github.com/wsyqn6/hl7">📚 GoDoc</a> •
  <a href="https://github.com/wsyqn6/hl7/issues">🐛 Issues</a> •
  <a href="https://hl7.org">🌐 HL7 Official</a>
</p>
