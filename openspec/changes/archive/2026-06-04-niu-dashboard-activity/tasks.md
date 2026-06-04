## 1. 后端 API 与服务

- [x] 1.1 读取命中规则文件:`backend-go.md`、`api-contract.md`、`database.md`、`data-permission.md`、`architecture.md`、`frontend-ui.md`、`testing.md`、`documentation.md`、`cache-consistency.md`、`i18n.md`(逐一记录影响/无影响)。
- [x] 1.2 `api/settlement/v1`:活跃度 DTO(`GET .../settlement/activity`,`days` 默认 14/上限 60,响应 dau 序列 + retentionD1/D7),带 `permission:"sicau-niu:settlement:view"`;`api/settlement/settlement.go` 接口追加方法。
- [x] 1.3 `internal/service/settlement/settlement_activity.go`:去重活跃集合(参数化 UNION,引用 DAO 表/列常量)、DAU 分组计数(日历补 0)、次日/7 日留存(排除窗口未过 cohort);集合化、有界、无 N+1。
- [x] 1.4 `internal/controller/settlement`:活跃度处理方法委派 service。
- [x] 1.5 `plugin.go`:管理端子组绑定活跃度接口。

## 2. 前端

- [x] 2.1 结算页新增「活跃度」区:DAU 简表 + 次日/7 日留存卡片;data-testid 锚点,中文,无 i18n。

## 3. 测试与验证

- [x] 3.1 服务层 DB 门控单测:DAU 去重计数(多表行为同日只计一次)、次日留存口径、7 日留存口径、窗口未过 cohort 被排除;沿用共享临时库 + TRUNCATE + TestMain。
- [x] 3.2 编译门禁:`temp/go.work.sicau-verify` 下 `go build ./...` + `go vet ./...`。
- [x] 3.3 `openspec validate niu-dashboard-activity --strict` 通过。
- [x] 3.4 `lina-review` 通过(记录无新增表/无 N+1/缓存无影响/索引取舍/i18n 无影响)。

> E2E:管理端活跃度区端到端验证并入最终统一 E2E 阶段。
