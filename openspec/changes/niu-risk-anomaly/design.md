## Context

补齐 M13 的异常喂草/偷草告警。复用 C7 结算管理端鉴权与 C4 `feeding`、`steal` 表。阈值可配置。约束遵循 `.agents/rules/` 的 `backend-go`、`api-contract`、`data-permission`、`architecture`、`testing`;单语言中文不启用 i18n;`platform_only`。

## Goals / Non-Goals

**Goals:** 单玩家单日喂草/偷草超阈值告警(只读派生视图);阈值配置可调;集合化、有界、无 N+1。

**Non-Goals:** 自动处置/封禁;异常定位细化检测;告警事件持久化表;新增数据表。

## Decisions

### D1 异常口径与数据来源
- **喂草异常**:`feeding` 按 `(user_id, DATE(created_at))` 分组,`HAVING COUNT(*) > feedDailyThreshold`。喂草无每日次数硬上限,是主要异常信号。
- **偷草异常**:`steal` 按 `(actor_user_id, steal_date)` 分组,`HAVING COUNT(*) > stealDailyThreshold`。偷草规则上限为每日 N 次(C4),阈值默认贴近上限以识别持续顶格行为。
两类各一条集合查询,合并为统一告警列表;玩家昵称批量装配(一次 WHERE IN),杜绝 N+1。结果按当日次数倒序、有界(上限)。

### D2 阈值配置可调
settlement 服务 `New(Config{FeedDailyThreshold, StealDailyThreshold, AnomalyLimit})` 注入。`plugin.go` 读取插件配置 `anomaly.feedDailyThreshold`(默认 100)、`anomaly.stealDailyThreshold`(默认 5)与结果上限(默认 200),非正值回退默认。写入 `config.example.yaml`。

### D3 只读派生 + 数据权限
管理端 `Auth+Tenancy+Permission` + DTO `permission:"sicau-niu:settlement:view"`;仅告警供人工研判,无写、无自动处置。

### API 契约(管理端鉴权,日期 yyyy-mm-dd 字符串)
| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| `GET` | `/plugins/sicau-niu/settlement/risk/anomalies` | `sicau-niu:settlement:view` | 单玩家单日喂草/偷草超阈值告警(有界) |

响应 `list: []{userId, nickname, type(feed|steal), date, count, threshold}`。

### D4 显式依赖注入
settlement `serviceImpl` 新增阈值字段,装配期注入;controller 委派;`plugin.go` 绑定。不在请求路径 `New()`。

## Risks / Trade-offs

- [偷草规则已限次,异常信号弱] → 阈值默认贴近上限识别顶格;喂草无上限是主信号;两者均可配置。
- [阈值误报] → 仅告警不处置,人工研判;阈值运营可调。
- [跨两表合并成本] → 各一条集合 HAVING 查询 + 一次批量昵称,有界、无 N+1。

## Migration Plan

1. settlement `Config` + 异常聚合 service → DTO/controller → `plugin.go` 读 `anomaly.*` 配置并装配绑定 + `config.example.yaml`。
2. 前端结算页「异常告警」表。
3. 服务层 DB 门控单测(超阈值命中、临界不报、feed/steal 口径)+ 编译门禁 + `openspec validate --strict`;E2E 并入统一阶段。
4. 回滚:纯新增读取路径 + 配置项,无表/数据,可整体回退。

## Open Questions

1. 默认阈值(暂定 feed 100/日、steal 5/日)与结果上限(200)。

> 影响判断:字典 —— 无新增枚举,**无影响**;缓存一致性 —— 实时聚合,**无新增缓存失效路径**;数据权限 —— 管理端鉴权 + `permission` 标签,只读有界(已记录);跨平台 —— 仅插件内文件,**无影响**;i18n —— 单语言,**无影响**。
