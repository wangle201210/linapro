## 1. 后端服务与契约

- [x] 1.1 读取命中规则文件:`backend-go.md`、`api-contract.md`、`data-permission.md`、`architecture.md`、`frontend-ui.md`、`testing.md`、`documentation.md`、`cache-consistency.md`、`i18n.md`(逐一记录影响/无影响)。
- [x] 1.2 `internal/service/record`:`Service` + `New()`;6 个分页查询方法(feeding/steal/gift/checkin/activation/grass_txn),DB 侧 Count+Page+OrderDesc;`record_relations.go` 批量装配玩家昵称/牛名(WhereIn,无 N+1);`record_code.go` 查询错误码。
- [x] 1.3 `api/record/v1`:6 列表 DTO(分页/筛选/`{list,total}`,时间 Unix 毫秒,权限 `sicau-niu:record:list`);`api/record/record.go` 接口聚合。
- [x] 1.4 `internal/controller/record`:6 个处理方法委派 service。
- [x] 1.5 `plugin.go`:装配 record service/controller,管理端子组绑定 6 接口。

## 2. 前端与菜单

- [x] 2.1 `plugin.yaml`:新增"活动记录"子目录(D)+ 6 个查询页(M)与权限 `sicau-niu:record:list`。
- [x] 2.2 前端 `record-client.ts` + 6 个 vxe 查询页(筛选 + 分页 + 装配列),data-testid,中文,无 i18n。

## 3. 测试与验证

- [x] 3.1 服务层 DB 门控单测:feeding 按牛筛选 + 分页 + 玩家/牛名装配、steal 双方昵称装配、grass_txn 按玩家分页;空页;沿用共享临时库 + TRUNCATE + TestMain。
- [x] 3.2 编译门禁:`temp/go.work.sicau-verify` 下 `go build ./...` + `go vet ./...`。
- [x] 3.3 `openspec validate niu-activity-records --strict` 通过。
- [x] 3.4 `lina-review` 通过(记录只读数据权限/无新增表/无 N+1/缓存无影响/i18n 无影响)。

> 管理端 6 个记录页 E2E 并入最终统一阶段(mock 数据已含各类记录)。
