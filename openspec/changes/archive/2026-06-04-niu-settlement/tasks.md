## 1. 数据库与 DAO

- [x] 1.1 读取命中规则文件:`database.md`、`plugin.md`、`api-contract.md`、`backend-go.md`、`data-permission.md`、`architecture.md`、`frontend-ui.md`、`testing.md`、`documentation.md`、`cache-consistency.md`、`i18n.md`(逐一记录影响/无影响)。
- [x] 1.2 新增 `manifest/sql/006-sicau-niu-settlement.sql`:表 `plugin_sicau_niu_settlement`(id/title/snapshot/operator_id/archived_at/created_at/updated_at/deleted_at),幂等 `CREATE TABLE IF NOT EXISTS` + `archived_at` 倒序索引(`WHERE deleted_at IS NULL`)。
- [x] 1.3 新增 `manifest/sql/uninstall/006-sicau-niu-settlement.sql`(幂等 `DROP TABLE IF EXISTS`)。
- [x] 1.4 `backend/hack/config.yaml` `tables` 追加 `plugin_sicau_niu_settlement`;`env -u GOWORK make dao p=sicau-niu` 生成 dao/do/entity。

## 2. 后端 API 与服务

- [x] 2.1 新增 `backend/api/settlement/v1` + `api/settlement/settlement.go` 接口:看板/导出/发证/风控/归档(创建+列表)DTO,均带 `permission:"sicau-niu:settlement:*"`;时间字段 Unix 毫秒(`*int64`);导出含 `truncated` 标注。
- [x] 2.2 `internal/service/settlement`:看板(集合聚合)、导出(有界 + 批量装配 + 截断标注)、发证(证书校验 + 集合筛选 + 事务幂等写 user_honor)、风控(设备指纹分组簇 + 批量昵称)、归档(快照持久化 + 有界列表);显式注入,杜绝 N+1。
- [x] 2.3 `internal/controller/settlement`:委派 service。
- [x] 2.4 `plugin.go`:管理端子组绑定结算控制器;装配期注入 settlement service。

## 3. 前端与菜单

- [x] 3.1 `plugin.yaml`:顶级"寻牛活动"下新增"运营结算"菜单页与按钮权限(`sicau-niu:settlement:{view,export,issue,archive}`)。
- [x] 3.2 前端"运营结算"页 `sicau-niu-settlement.vue` + `settlement-client.ts`:看板卡片 + 导出按钮(CSV) + 批量发证(证书 ID + 执行) + 风控告警表 + 归档列表/创建;data-testid 锚点,中文,无 i18n。

## 4. 测试与验证

- [x] 4.1 服务层 DB 门控单测(8 用例全绿):看板聚合计数、发证幂等(重复发新发=0)+ feed_count 队列筛选 + 拒绝非证书/图鉴完成类/不存在、风控簇、归档持久化与列表;沿用共享临时库 + TRUNCATE + TestMain。
- [x] 4.2 编译门禁:`temp/go.work.sicau-verify` 下 `go build ./...` + `go vet ./...` 均通过。
- [x] 4.3 `openspec validate niu-settlement --strict` 通过。
- [x] 4.4 `lina-review` 通过(记录幂等/数据权限/无 N+1/缓存无影响/i18n 无影响)。

> 前端 vite 构建与管理端结算页 E2E:按用户指示并入最终统一 E2E 阶段执行。

> E2E:按用户指示,管理端结算页端到端验证并入最终统一 E2E 阶段执行。
