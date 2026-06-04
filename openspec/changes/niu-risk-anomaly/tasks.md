## 1. 后端服务与配置

- [x] 1.1 读取命中规则文件:`backend-go.md`、`api-contract.md`、`database.md`、`data-permission.md`、`architecture.md`、`frontend-ui.md`、`testing.md`、`documentation.md`、`cache-consistency.md`、`i18n.md`(逐一记录影响/无影响)。
- [x] 1.2 settlement 服务 `New(Config{FeedDailyThreshold, StealDailyThreshold, AnomalyLimit})`;`serviceImpl` 新增阈值字段;非正值回退默认。
- [x] 1.3 `internal/service/settlement/settlement_anomaly.go`:喂草按 `(user_id, DATE(created_at))`、偷草按 `(actor_user_id, steal_date)` 分组 `HAVING COUNT > 阈值`,合并 + 批量昵称,有界倒序;集合化无 N+1。
- [x] 1.4 `plugin.go`:读取 `anomaly.feedDailyThreshold`/`anomaly.stealDailyThreshold`/上限配置,装配 `settlementsvc.New(...)`;`config.example.yaml` 增配置项。

## 2. 接口与前端

- [x] 2.1 `api/settlement/v1` 异常告警 DTO(`GET .../settlement/risk/anomalies`,带 `permission:"sicau-niu:settlement:view"`);`api/settlement/settlement.go` 接口追加方法;controller 处理方法。
- [x] 2.2 `plugin.go` 管理端子组绑定异常告警接口。
- [x] 2.3 前端结算页「异常告警」表(玩家/类型/日期/次数/阈值),data-testid 锚点,中文,无 i18n。

## 3. 测试与验证

- [x] 3.1 服务层 DB 门控单测:喂草超阈值命中、偷草超阈值命中、临界(等于阈值)不报、合并与昵称装配;沿用共享临时库 + TRUNCATE + TestMain。
- [x] 3.2 编译门禁:`temp/go.work.sicau-verify` 下 `go build ./...` + `go vet ./...`。
- [x] 3.3 `openspec validate niu-risk-anomaly --strict` 通过。
- [x] 3.4 `lina-review` 通过(记录数据权限/无新增表/无 N+1/缓存无影响/i18n 无影响)。

> E2E:管理端异常告警表端到端验证并入最终统一 E2E 阶段(已并入 TC004 运营结算页)。
