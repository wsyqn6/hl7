# HL7 库路线图

## 当前版本: v0.3.0 ✅

### 已完成功能
- 核心解析/编码/验证
- MLLP 网络传输
- 对象池/零拷贝/并行解析
- CI/CD 集成

---

## v0.4.0 - FHIR 集成

### 目标
连接传统 HL7 v2.x 系统与现代 FHIR 应用

### 功能清单
- [ ] 创建 `fhir/` 子包
- [ ] 实现 HL7 v2 → FHIR Bundle 转换器
- [ ] 实现 Patient 资源映射 (PID → Patient)
- [ ] 实现 Observation 资源映射 (OBX → Observation)
- [ ] 实现 Encounter 资源映射 (PV1 → Encounter)
- [ ] 实现 DiagnosticReport 资源映射 (OBR → DiagnosticReport)
- [ ] 实现 Condition 资源映射 (DG1 → Condition)
- [ ] FHIR 资源验证
- [ ] 单元测试覆盖

### 技术选型
- 使用标准 FHIR 类型定义
- 无外部 FHIR 库依赖（保持轻量）

---

## v0.5.0 - 开发者体验

### 目标
降低使用门槛，提升开发效率

### 功能清单
- [ ] CLI 工具
  - `hl7 validate` - 验证消息
  - `hl7 convert` - 格式转换
  - `hl7 serve` - 本地测试服务
- [ ] 在线 Playground 网页
- [ ] 更多示例 (10+ 场景)
- [ ] 交互式教程

---

## v1.0.0 - 稳定版本

### 目标
正式发布，稳定 API

### 功能清单
- [ ] API 冻结（向后兼容承诺）
- [ ] 性能基准测试公开
- [ ] 完整英文/中文文档
- [ ] API 文档网站
- [ ] 安全审计

---

## 长期目标

- [ ] HL7 FHIR 双向网关
- [ ] HL7 v3 支持
- [ ] 更多语言版本文档
- [ ] 性能对比页面
- [ ] 开源贡献者社区

---

## 参考资料

- [FHIR 官网](https://www.hl7.org/fhir/)
- [HL7 国际](https://www.hl7.org/)
- [FHIR 资源类型](https://www.hl7.org/fhir/resourcelist.html)
