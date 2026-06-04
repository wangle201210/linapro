## 1. 公开纪念墙 API 契约

- [x] 1.1 读取命中规则文件:`plugin.md`、`api-contract.md`、`backend-go.md`、`data-permission.md`、`architecture.md`、`frontend-ui.md`、`testing.md`、`documentation.md`(单语言不启用 i18n,记录无影响)。
- [x] 1.2 新增 `backend/api/wall/v1`:首发墙/校史精选/统计 公开 DTO(无 permission 标签);时间字段 Unix 毫秒(`*int64`,沿用 `apitime.Milli`);响应仅含昵称/身份/牛名称/序号/到场顺序与公开统计,**不含手机号/openid**。
- [x] 1.3 手写 `api/wall/wall.go` 接口聚合(控制器为自定义结构,沿用 player/admin 同构方式,免代码生成器)。

## 2. 服务层与路由

- [x] 2.1 `internal/service/wall`:首发墙(`activation.is_first` 关联 `user`/`niu`,批量装配,有界 ≤120)、校史精选(`card`/`quote` 取样有上限)、统计(DB 侧计数);`New()` 无依赖,隐私字段在 service 投影裁剪。
- [x] 2.2 填充 `internal/controller/wall` 委派 service。
- [x] 2.3 `plugin.go`:在公开子路由组(无 Auth、无 player-auth,与玩家登录同级)绑定纪念墙 API。

## 3. H5 静态页与 public_assets

- [x] 3.1 新增 `frontend/wall/index.html`(移动端 H5:fetch 公开 API 渲染首发墙/精选/统计 + 回跳小程序占位入口)。
- [x] 3.2 `plugin.yaml` 声明 `public_assets`(source `frontend/wall`,mount `/wall`,index `index.html`);`plugin_embed.go` 已嵌入 `frontend`(含 `frontend/wall`)。

## 4. 测试与验证

- [x] 4.1 服务层 DB 门控单测:首发墙派生与有界、**响应不含隐私字段**断言、校史精选有上限(仅启用金句)、统计计数;沿用共享临时库 + TRUNCATE + TestMain(4 用例全绿)。
- [x] 4.2 编译门禁:`temp/go.work.sicau-verify` 下 `go build ./...` + `go vet ./...` 均通过。
- [x] 4.3 公开 API 冒烟 + H5 静态资源可访问:`make dev` 后经 vite(5666)与后端(9120)无 token 访问 `wall/{stats,first-activators,highlights,config}` 均 `code:0`、`/x-assets/sicau-niu/v0.1.0/wall/index.html` HTTP 200;并新增 **E2E `TC005-h5-wall`**(匿名 page:H5 页渲染公开数据 + 无 token 公开 API 返回 `code:0`),2/2 通过。
- [x] 4.4 `openspec validate niu-h5-wall --strict` 通过。
- [x] 4.5 `lina-review` 通过(记录隐私裁剪、无新增表、i18n 无影响)。

> E2E:按用户指示,H5 与公开 API 的端到端验证并入最终统一 E2E 阶段执行。
