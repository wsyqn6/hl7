# AGENTS.md

> 本文件为 agentic coding agents 提供开发规范指引。

## 项目概述

HL7 v2.x 医疗消息解析库 - 支持结构体序列化、流式解析、MLLP传输、消息验证、ACK生成。
**Go 版本要求**: 1.26+

## 构建与测试

```bash
make            # 等同于 make test
make test       # go test ./...
make build      # go build ./...
make coverage   # 覆盖率报告
make bench      # 基准测试
```

## 代码风格

### 命名约定

| 类型 | 规则 | 示例 |
|------|------|------|
| 包名 | 小写，无下划线 | `hl7`, `types` |
| 导出 | CamelCase | `Message`, `Parse` |
| 变量 | camelCase | `msg`, `err` |
| 私有 | 小写开头 | `parseField` |
| 常量 | SCREAMING_SNAKE | `MLLP_START` |
| 接口 | `-er` 结尾 | `Validator`, `Parser` |

### 错误处理

- 使用 `errors.Is()` / `errors.As()` 比较错误
- 包装错误：`fmt.Errorf("parsing: %w", err)`
- 自定义错误类型用于可预见错误

### 文档注释

所有导出符号必须有注释，以符号名开头：

```go
// Message represents an HL7 v2.x message.
type Message struct { ... }

// Parse parses raw HL7 data into a Message.
func Parse(data []byte) (*Message, error) { ... }
```

## 测试规范

- 测试文件以 `_test.go` 结尾
- 使用表驱动测试
- 测试/基准/模糊测试分别以 `Test`/`Benchmark`/`Fuzz` 开头

```go
func TestParseField(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"simple", "value", false},
        {"empty", "", false},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := ParseField(tt.input)
            if (err != nil) != tt.wantErr {
                t.Fatalf("error = %v", err)
            }
        })
    }
}
```

## 提交前检查

1. `go fmt ./...` - 格式化
2. `go test ./...` - 测试通过
3. `go vet ./...` - 静态分析
4. 新功能有对应测试
5. 导出符号有文档注释

## 重要提示

- 保持公开 API 向后兼容
- 不要提交敏感信息
- 遵循 HL7 v2.x 规范 ([hl7.org](https://www.hl7.org/))
