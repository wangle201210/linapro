## Context

补齐 M5 看板的日活(DAU)与留存。复用 C7 结算管理端鉴权与 C1–C4 行为表(activation/feeding/checkin/steal/gift)。活动仅 2 个月、数据有界,采用**按需聚合**(不新增表、不引定时任务)。约束遵循 `.agents/rules/` 的 `backend-go`、`api-contract`、`database`、`data-permission`、`architecture`、`frontend-ui`、`testing`;单语言中文不启用 i18n;`platform_only`。

## Goals / Non-Goals

**Goals:** 日活(最近 N 天每日去重活跃玩家)、整体次日留存与 7 日留存(cohort/回访/留存率);DB 侧集合聚合、有界、无 N+1、无新增表。

**Non-Goals:** 每 cohort 留存曲线明细;30 日等长档留存;实时秒级刷新;埋点 SDK。

## Decisions

### D1 活跃定义与去重活跃集合
「活跃」= 当日在任一行为表有记录的玩家。行为表:`activation`、`feeding`、`checkin`、`steal`(actor)、`gift`(from)。构造一个**去重活跃关系** `active(user_id, day)`,由各表 `(user_id, DATE(created_at))` 的 `UNION` 派生。SQL 用**参数化**写法,表名/列名引用 DAO 常量(`dao.Feeding.Columns()...`)防漂移;日期范围参数化,杜绝注入与全表无界扫描。

### D2 日活(DAU)序列
对 `active` 关系按 `day` 分组 `COUNT(DISTINCT user_id)`,限定最近 `days` 天(默认 14、上限 60)。一条集合查询返回 `[{date, activeUsers}]`,缺口天补 0 在服务层用日历填充(玩家级数据有界)。

### D3 留存(次日 / 7 日)
cohort 按注册日 `DATE(user.created_at)`。对档位 `n∈{1,7}`:
- 合格 cohort = 注册日满足 `register_day + n <= 今日`(窗口已过)的玩家。
- 回访 = 这些玩家在 `register_day + n` 当天出现在 `active` 关系中。
- `留存率 = 回访数 / cohort 数`(cohort 为 0 时率为 0)。
整体口径(跨所有合格 cohort 汇总),DB 侧集合计算,不逐玩家循环。

### D4 接口契约(管理端鉴权,日期为 yyyy-mm-dd 字符串)
| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| `GET` | `/plugins/sicau-niu/settlement/activity` | `sicau-niu:settlement:view` | 日活序列 + 次日/7 日留存 |

请求 `days`(默认 14,上限 60);响应 `dau: []{date, activeUsers}`、`retentionD1`/`retentionD7: {cohortUsers, returnedUsers, rate}`(rate 为 0–1 浮点)。

### D5 分层与依赖
活跃度聚合作为 settlement service 的新增方法;controller 委派;`plugin.go` 在管理端子组绑定。显式注入,装配期构建,不在请求路径 `New()`。

### D6 索引评估
活跃集合按 `DATE(created_at)` 过滤/分组。各行为表已含 `created_at`;范围有界(≤60 天、玩家级)。`DATE()` 函数转换会绕过普通 `created_at` 索引,但**数据规模有界(2 个月活动)**,且为运营低频查询,不在玩家高频路径 —— 评估后**不为此新增函数索引**,在设计中记录该取舍(符合 database 规则「确需派生值时明确可查询字段/索引/缓存方案」:此处选择有界全扫的运营查询,不引入低价值索引)。

## Risks / Trade-offs

- [跨表 UNION 成本] → 数据有界(活动 2 个月)、运营低频;集合化一次查询,无 N+1。
- [DATE() 绕过索引] → 有界数据可接受;已记录取舍,后续如量级变化可加表达式索引或每日快照。
- [留存窗口未过的 cohort] → 显式排除 `register_day+n > 今日` 的 cohort,避免低估。

## Migration Plan

1. `api/settlement/v1` 活跃度 DTO → 填 controller/service 聚合;`plugin.go` 绑定。
2. 前端结算页「活跃度」区。
3. 服务层 DB 门控单测(DAU 计数、次日/7 日留存口径,含窗口未过排除)+ 编译门禁 + `openspec validate --strict`;E2E 并入统一阶段。
4. 回滚:纯新增读取路径,无表、无数据,可整体回退。

## Open Questions

1. DAU 默认展示天数(暂定 14)。

> 影响判断:字典 —— 无新增枚举,**无影响**;缓存一致性 —— 实时按需聚合,**无新增缓存失效路径**;数据权限 —— 管理端鉴权 + `permission` 标签,只读有界(已记录);跨平台 —— 仅插件内文件,**无影响**;i18n —— 单语言,**无影响**。
