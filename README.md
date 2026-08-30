# golden-field · 包消费实证下游仓（VP-022～VP-024）

`schema-ui-core` 构建期包消费的**金标准试验田**：仅经 cli+包 形态（Go `go get` / npm registry）消费上游基架，实证「真实发布 → 无凭据安装 → 升级 → 运行」的完整回路。

- 关联：上游 [schema-ui-core](https://github.com/magicvr/schema-ui-core)（Charter `schema-ui-core-admin-foundation@0.3.0`）· workspace-024（分发形态正式化）· VP-024 **closed**（2026-08-29 · VRev-053）。
- 本仓 = VP-024 判据 #2/#3/#5「公开消费往返实证」的**持久可复现宿主**（closure-report 证据链的可克隆副本）。
- 建立：2026-08-29（用户建仓 · 初始化为编排器义务）；首次推送：2026-08-30。

## 结构

```
cmd/server/             Go 消费组合根（thin wrapper 调 apps/api kernel `server.Serve`；装配/迁移冒烟）
web/                    前端消费骨架（@magicvr/schema-ui-* 六包 + 探针 + pnpm-workspace.yaml 终值）
.github/workflows/      consumer-regression.yml（hosted CI 回归）
```

## 消费实态（registry 语义 · 免凭据）

| 通道 | 状态 |
|------|------|
| Go | `go.mod` require `github.com/magicvr/schema-ui-core/apps/api v0.4.0`（公共 proxy · **无 replace / 无 file:**） |
| npm | 六包终值（冻结面 v1.4.0）：lib `0.1.10` · protocol `0.2.11` · renderer `0.3.8` · shell `0.1.4` · theme `0.1.4` · ui `0.1.8`；`web/.npmrc` 钉 `@magicvr:registry=https://registry.npmjs.org`；`pnpm-lock.yaml` 全 npmjs tarball（GH Packages 残留 0 · integrity 与 npmjs dist 一致） |

## 验证探针

- **Go**：`go build -o gf-server ./cmd/server` → `./gf-server -dialect sqlite -dsn <db> -addr <port>` → `/healthz` / `/readyz` 200
- **Web**（无凭据：空 userconfig + npmjs 公开 registry）：
  - `node probe.mjs`（protocol 2.9）· `probe-render.mjs`（render 1573 B）· `probe-six.mjs`（六包可消费 · DataTable）· `token-check.mjs`（brand 覆盖）· `probe-r5.mjs`（external 断言 17 imports）
  - **UI-ONLY**：仅装 ui + peer 独立消费 PASS
- **CI**：`.github/workflows/consumer-regression.yml`——免凭据（删除 GH token 面）· `pnpm/action-setup` 先于 `setup-node`（`packageManager: pnpm@11.11.0`）· Go `@latest` 升级 + build · serve 后台 + healthz/readyz 轮询 + 四探针 · SIGTERM 收尾断言 `shutdown.complete`（RT-D02 出口）。触发：`workflow_dispatch`（手动）；或主仓发布后 `repository_dispatch`（types: published）。

## 历史锚点

| commit | 事件 |
|--------|------|
| `ef67f0c` | VP-022 R1：GH Packages registry 消费（私有面 · 三探针） |
| `04e73e4` | VP-023：PG external 实测 · compose/ops-playbook · consumer-regression 初版 |
| `9510023` | VP-023：breaking 迁移演练（v0.2.0→v0.3.0）旧态基准 |
| `7fbabda` | VP-024 R2：v0.3.0→v0.4.0 · thin wrapper（serve 面） |
| `fb957a9` | VP-024 R2：lockfile 重解析 npmjs（F-001 fixed · 全凭据零残留） |
| `c4d14ea` | VP-024 R3：workflow 免凭据化重构 |
| `3f2a5c2` | VP-024 R3：pnpm/action-setup 顺序修复 |
| `235196d` | VP-024 R5：六包终值（冻结面 v1.4.0） |
| `8631d53` | 2026-08-30：初始化收口（README 终态 · 二进制清理 · pnpm-workspace 终值） |
| `ba052e7` | 2026-08-30：hosted 首跑修复（pnpm/action-setup 显式 version） |
| `8ef02e9` | 2026-08-30：收尾清理段加固（trap/kill/wait `|| true` · 显式 exit 0） |

## 残余与登记（主仓权威口径）

- hosted CI 实触发：**2026-08-30 闭环**——`33286154992`（action-setup 版本源 ❌）→ `33286191334`（探针全绿 · 清理段 exit 1 ❌）→ `33286302663` **PASS**（1m9s · 四探针 + `shutdown.complete`）。主仓口径见 `GOAL-008` E-004。
- 本仓仅作消费实证，不承载目标状态；主仓治理真相源 = `docs/workspaces/workspace-024-distribution-formalization/`。