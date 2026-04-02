# AGENTS.md

> 本文件为 agentic coding agents（如 GitHub Copilot、Cursor）提供开发规范指引。

## 项目概述

HL7 v2.x 医疗消息解析库 - 支持结构体序列化、流式解析、MLLP传输、消息验证、ACK生成。

**Go 版本要求**: 1.26+

## 构建与测试命令

### 基本命令

```bash
go build ./...              # 构建所有包
go test ./...               # 运行测试
go test -v ./...            # 详细输出
go test -cover ./...        # 显示覆盖率
go test -coverprofile=coverage.out  # 生成覆盖率文件
go tool cover -html=coverage.out    # 查看覆盖率报告
go vet ./...                # 静态分析
go fmt ./...                # 格式化代码
go mod tidy                 # 整理依赖
```

### 运行单个测试

```bash
# 运行特定测试函数
go test -v -run TestUnmarshalMessage

# 运行特定包的测试
go test -v ./internal/escape

# 使用通配符匹配测试
go test -v -run "TestParse.*"

# 运行基准测试
go test -bench=. -benchmem

# 运行模糊测试
go test -fuzz=FuzzParse -fuzztime=10s
```

### 性能分析

```bash
go test -cpu 1,2,4 -bench=. ./...    # 多CPU基准测试
go test -trace trace.out ./...       # 生成追踪文件
go tool trace trace.out              # 查看追踪
```

## 代码风格指南

### 命名约定

- **包名**: 小写单词，无下划线 (`hl7`, `mllp`, `types`)
- **类型/函数**: CamelCase 导出 (`Message`, `Parse`, `GenerateACK`)
- **变量**: camelCase (`msg`, `err`, `parseErr`)
- **私有符号**: 小写开头 (`parseField`, `validateSegment`)
- **常量**: CamelCase 或 SCREAMING_SNAKE_CASE
- **接口**: 形容词或 `-er` 结尾 (`Validator`, `Parser`, `Marshaler`)

### 导入组织（goimports 标准）

```go
package hl7

import (
    // 标准库
    "context"
    "fmt"
    "io"
    "strings"
    "time"

    // 第三方库
    // "github.com/some/pkg"

    // 内部包
    // "github.com/username/hl7/internal/escape"
)
```

### 错误处理

- 总是检查并处理错误
- 使用 `errors.Is()` 和 `errors.As()` 进行错误比较
- 包装错误添加上下文：`fmt.Errorf("parsing message: %w", err)`
- 定义自定义错误类型用于可预见的错误条件

### 文档注释

- 所有导出的类型、函数、变量、常量必须有注释
- 注释以符号名开头
- 包注释放在 `doc.go` 文件

```go
// Message represents an HL7 v2.x message containing multiple segments.
type Message struct { ... }

// Parse parses raw HL7 data into a Message.
func Parse(data []byte) (*Message, error) { ... }
```

### 测试规范

- 测试文件以 `_test.go` 结尾
- 使用表驱动测试处理多个测试用例
- 测试函数以 `Test` 开头
- 基准测试以 `Benchmark` 开头
- 模糊测试以 `Fuzz` 开头

```go
func TestParseField(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected Field
        wantErr  bool
    }{
        {"simple", "value", Field{Value: "value"}, false},
        {"empty", "", Field{}, false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ParseField(tt.input)
            if (err != nil) != tt.wantErr {
                t.Fatalf("ParseField() error = %v", err)
            }
            if !reflect.DeepEqual(got, tt.expected) {
                t.Errorf("ParseField() = %v, want %v", got, tt.expected)
            }
        })
    }
}
```

### 性能考虑

- 避免不必要的内存分配
- 使用 `strings.Builder` 构建字符串
- 对于大文件使用流式处理 (`io.Reader` / `io.Writer`)
- 使用 `sync.Pool` 复用临时对象

## 提交前检查清单

在提交代码前，确保：

1. ✅ `go fmt ./...` - 代码格式化
2. ✅ `go mod tidy` - 依赖整理
3. ✅ `go test ./...` - 所有测试通过
4. ✅ `go vet ./...` - 静态分析通过
5. ✅ 新功能有对应测试
6. ✅ 导出符号有文档注释
7. ✅ 错误处理正确（没有忽略错误）

## 常见任务

### 添加新功能
1. 实现功能代码
2. 添加单元测试
3. 添加基准测试（如需要）
4. 更新文档注释
5. 运行 `go test ./...` 验证

### 修复 Bug
1. 添加复现 bug 的失败测试
2. 修复 bug
3. 验证测试通过
4. 检查是否引入新问题

## 重要提示

- 保持向后兼容性（公开 API）
- 不要提交敏感信息或密钥
- 遵循 HL7 v2.x 规范 ([hl7.org](https://www.hl7.org/))
