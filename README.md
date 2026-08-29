# golden-field · 产线化实验下游仓（VP-023）

`schema-ui-core` 包消费产线化的**金标准试验田**：仅经 cli+包 形态（Go `go get` / npm registry）消费上游基架，验证真实发布 → 安装 → 升级 → 运行的完整回路。

- 关联：上游 [schema-ui-core](https://github.com/magicvr/schema-ui-core)（愿景 `schema-ui-core-admin-foundation@0.3.0`）· VP-023（组合层平台波）· workspace-023。
- 建立：2026-08-29（用户建仓 · 本仓初始化由编排器完成）。

## 结构

```
cmd/server/   Go 组合根（apps/api kernel + assembly + 模块；装配/迁移冒烟）
web/          前端骨架（@schema-ui/* tarball 消费 + 主题覆盖 + 探针）
```

## 消费实态（registry 语义 · R1 起）

| 通道 | 现状 | 移除条件（VP-023 判据 #1） |
|------|------|------------------------------|
| Go | `go.mod` require `apps/api v0.1.0`（公共 proxy · 无 replace） | registry 语义已实证（R1）；随 R5 升 v0.2.0 |
| npm | `web/package.json` 六包 registry（`@magicvr/schema-ui-*`） | 已实证（R3） |

## 验证探针

- Go：`go run ./cmd/server -dialect sqlite -dsn <db>` → golden-field kernel=... 装配/迁移冒烟
- Web：`pnpm install` 后 `node probe.mjs` / `probe-render.mjs` / `token-check.mjs`

## 目标（R5 验收走查）

从零（本仓）经 cli+包 到上线：create 骨架 → 装配 → registry 升级（含 breaking 场景）→ 运维路径文档对照。