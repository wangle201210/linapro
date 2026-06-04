## Context

C7 是寻牛活动的运营结算能力(程序级设计 C7:M5 看板 / M13 告警 / M14 公示结算)。复用 C1 玩家、C2 牛、C3 激活、C4 喂草/偷草/送草/签到、C5 荣誉。核心约束是**集合化聚合杜绝 N+1**、**发证幂等**与**结算快照可留存**。

约束:遵循 `.agents/rules/` 的 `plugin`、`backend-go`、`api-contract`、`database`、`data-permission`、`architecture`、`frontend-ui`、`testing`;单语言中文不启用 i18n;`platform_only`。

## Goals / Non-Goals

**Goals:** 运营看板(集合聚合)、名册导出(有界批量)、证书批量幂等发放、风控告警(一机多号)派生视图、结算快照归档与查询;均经宿主鉴权 + 数据权限。

**Non-Goals:** 小程序 UI;皮肤;赠送撤销;风控事件持久化为独立表(本期派生只读)。

## Decisions

### D1 看板:集合化实时聚合
看板 = 一组 DB 侧聚合(`Count`/`Sum`/`GROUP BY`),每项一条有界查询:玩家数、激活牛数、全部牛数、首发数、喂草次数与总效果、偷草次数、送草次数、签到次数、已发证书数。**不**逐玩家循环,**不**引入缓存(实时读取),杜绝随玩家/牛数线性增长的查询。

### D2 导出:有界批量装配
名册导出取玩家页(硬上限 `exportMaxRows`,如 5000),再以 `WhereIn(userIDs)` 的两条分组查询批量装配 激活数 与 喂草总量,组装为行。超过上限时**截断并在响应标注 `truncated`**(no silent cap)。本期返回结构化行,前端导出 CSV。

### D3 证书批量发放:集合筛选 + 幂等写入
入参 `honorId` 必须指向 `honor_type=certificate` 的荣誉定义。按其 `unlock_type` **集合化**筛出达标玩家:
- `participation` → 全部玩家;
- `activation_count` → `GROUP BY user HAVING COUNT(activation) >= threshold`;
- `feed_count` → `GROUP BY user HAVING COUNT(feeding) >= threshold`;
- `category_complete`/`full_complete` → 图鉴完成类按个人在小程序内解锁,**不纳入批量结算**,返回业务错误明确拒绝。
对达标集合,排除 `user_honor` 中已授予(`WhereNotIn` 既有 user_id),其余在一个事务内批量 `Insert`(依赖 `uk_sicau_niu_user_honor` 唯一索引保证幂等)。返回 `eligible/issued/skipped` 计数。

### D4 风控告警:一机多号派生视图
集合化:`plugin_sicau_niu_user` 按 `device_fingerprint` 分组(非空)`HAVING COUNT(*)>1`,取簇;再 `WhereIn` 批量装配簇内玩家昵称。只读派生,无持久化。

### D5 结算归档:快照持久化(新表)
新增 `plugin_sicau_niu_settlement`(`id`/`title`/`snapshot`(看板指标 JSON 文本)/`operator_id`/`archived_at`/时间戳/`deleted_at`)。归档动作:实时计算看板 → 序列化为 `snapshot` → 写一行。列表查询按 `archived_at` 倒序、有界。GoFrame 自动维护时间与软删除;`archived_at` 由服务显式写当前时间(业务字段,非框架时间字段)。

### D6 显式依赖注入与分层
settlement service 显式注入所需读取(或复用 C5 honor service 的发证能力);controller 持 service 字段,装配期 `New`,不在请求路径构造。api/controller/service 三层;DTO 带 `permission` 标签。

### API 契约(C7,管理端鉴权,时间 Unix 毫秒)
| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| `GET` | `/plugins/sicau-niu/settlement/dashboard` | `sicau-niu:settlement:view` | 运营数据看板 |
| `GET` | `/plugins/sicau-niu/settlement/export/players` | `sicau-niu:settlement:export` | 玩家名册导出(有界) |
| `POST` | `/plugins/sicau-niu/settlement/certificates/issue` | `sicau-niu:settlement:issue` | 证书批量幂等发放 |
| `GET` | `/plugins/sicau-niu/settlement/risk/device-clusters` | `sicau-niu:settlement:view` | 一机多号风控簇 |
| `POST` | `/plugins/sicau-niu/settlement/archives` | `sicau-niu:settlement:archive` | 创建结算归档快照 |
| `GET` | `/plugins/sicau-niu/settlement/archives` | `sicau-niu:settlement:view` | 归档快照列表(有界) |

## Risks / Trade-offs

- [发证误发/重复] → 限定 `certificate` 类型 + 集合规则筛选 + 唯一索引幂等 + 事务;单测覆盖幂等与队列。
- [导出过大] → 硬上限 + `truncated` 标注 + 批量装配,杜绝 N+1。
- [看板性能] → 全部集合化聚合,无逐行循环;铁牛/牛数为小常量级,玩家级聚合走分组查询。
- [风控误报] → 仅作派生告警视图供人工研判,不自动处置。
- [归档快照口径] → 快照存 JSON 文本冻结当时指标;口径随看板定义,版本演进由新归档行体现。

## Migration Plan

1. `006-sicau-niu-settlement.sql` + uninstall;`hack/config.yaml` tables 追加;`env -u GOWORK make dao p=sicau-niu`。
2. `api/settlement` DTO → 填 controller/service;`plugin.go` 管理端子组绑定;`plugin.yaml` 菜单+按钮权限。
3. 前端"运营结算"页。
4. 服务层 DB 门控单测(看板/发证幂等/风控/归档)+ 编译门禁 + `openspec validate --strict`;E2E 并入统一阶段。
5. 回滚:卸载 SQL 删表;代码/菜单可整体回退;新表数据为运营快照,回滚不影响 C1–C5 业务数据。

## Open Questions

1. 导出硬上限默认值(暂定 5000,可后续配置化)。
2. 归档快照是否需要版本号/活动周期标识(本期用 `title` + `archived_at` 区分)。

> 影响判断:i18n 单语言不启用,**无影响**;字典 —— 无新增枚举(证书类型复用 C5 honor_type),**无影响**;缓存一致性 —— 看板实时聚合、归档为一次性持久化,**无新增缓存失效路径**;数据权限 —— 全部管理端动作经宿主 `Auth+Tenancy+Permission` 且带 `permission` 标签,发证幂等记录授予(已记录);跨平台 —— 仅插件内文件,**无影响**。
