# AGENTS.md

> HL7 v2.x 医疗消息解析库开发规范 | Go 1.26+

## 构建与测试

```bash
make test    # 测试
make build   # 构建
make coverage # 覆盖率
make bench   # 基准测试
```

## 命名约定

| 类型 | 规则 | 示例 |
|------|------|------|
| 导出 | CamelCase | `Message`, `Parse` |
| 私有 | 小写开头 | `parseField` |
| 常量 | SCREAMING_SNAKE | `MLLP_START` |
| 接口 | `-er` 结尾 | `Validator` |

## 代码规范

- 错误：`errors.Is()` 比较，`fmt.Errorf("context: %w", err)` 包装
- 文档注释：所有导出符号必须注释，以符号名开头
- 测试：`_test.go` 结尾，表驱动，`Test`/`Benchmark`/`Fuzz` 开头

## 提交前检查

```bash
go fmt ./... && go vet ./... && go test ./...
```

- 新功能必须有测试
- 导出符号必须有文档注释
- 保持公开 API 向后兼容
