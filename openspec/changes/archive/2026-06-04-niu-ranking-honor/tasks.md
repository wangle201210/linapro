## 1. 数据库与 DAO

- [x] 1.1 新增 `manifest/sql/005-sicau-niu-honor.sql`:`honor_def`(honor_type、code 活跃唯一、name、unlock_type、threshold、category、image_path、sort、软删/时间)、`user_honor`(user_id、honor_id、unlocked_at、`(user_id,honor_id)` 活跃唯一);`uninstall/005`(幂等)
- [x] 1.2 执行 `make dao p=sicau-niu` 生成 DAO/DO/Entity

## 2. API 契约与代码生成

- [x] 2.1 玩家面 `api/player`:三榜 `GET /player/rankings/{feed,college,friend}`、玩家荣誉 `GET /player/honors`(走 player-auth)
- [x] 2.2 运营面 `api/admin`:荣誉配置 `GET/POST /admin/honors`、`GET/PUT/DELETE /admin/honors/{id}`(permission 标签 `sicau-niu:honor:*`)
- [x] 2.3 RESTful、时间 Unix 毫秒、完整 dc/eg;执行 `make ctrl p=sicau-niu`

## 3. 服务层 — 排行榜(`niu-ranking`)

- [x] 3.1 个人喂草榜:`feeding GROUP BY user_id SUM(effect_amount)` 倒序 Top-N + 关联 user 昵称 + 本人名次
- [x] 3.2 院系榜:在校生按 college_id 聚合喂草总效果 Top-N(关联院系名)
- [x] 3.3 川农好友榜:identity=friend 玩家个人喂草效果 Top-N + 本人名次
- [x] 3.4 Top-N 上限配置;`*_code.go` bizerr 码;聚合 DB 侧完成杜绝内存全量

## 4. 服务层 — 荣誉(`niu-honor`)

- [x] 4.1 荣誉定义 CRUD:编码唯一、类型/解锁规则枚举校验;DB 侧分页
- [x] 4.2 荣誉类型/解锁规则 Go 命名常量 + 校验
- [x] 4.3 玩家荣誉解锁计算(只读):批量取玩家喂草次数/激活数/按分类已得卡数与分类卡片总数,对各 honor_def 计算解锁;杜绝逐荣誉 N+1
- [x] 4.4 `*_code.go` bizerr 码

## 5. 路由、依赖注入与运营前端

- [x] 5.1 `plugin.yaml` 在「寻牛活动」下新增荣誉配置菜单 + `sicau-niu:honor:*` 按钮权限
- [x] 5.2 `plugin.go`:玩家面绑定三榜/荣誉;运营面(Auth+Permission)绑定荣誉配置;controller 持 service 字段,service 显式注入 DAO/跨能力,不在请求路径 `New()`
- [x] 5.3 运营荣誉配置前端页(vxe-grid + 表单弹窗:类型/解锁规则/阈值/分类/图片路径)+ API 客户端

## 6. 测试与验证门禁

- [x] 6.1 服务层单测:三榜聚合/名次、荣誉编码唯一/枚举、玩家荣誉解锁计算(参与/阈值/集齐)(纯逻辑 + DB 门控,自包含含清理)
- [x] 6.2 后端编译门禁:plugins workspace 构建 + `go vet` + `go test`
- [x] 6.3 运营荣誉配置页 E2E 资产已编写(`TC003-honor-crud.ts`+`pages/SicauNiuHonorPage.ts`);玩家三榜/荣誉为 API。**按用户指示 E2E 统一在功能全开发完后执行**;注:荣誉 admin API 已 curl 验证(create code:0),TC003 当前前端表单交互待统一 E2E 时修(非产品缺陷)。
- [x] 6.4 运行 `openspec validate niu-ranking-honor --strict`
- [x] 6.5 见下方「影响判断」

## 影响判断(6.5)

- **i18n**:单语言中文,无 i18n 块/`$t()` —— **无影响**。
- **字典**:honor_type/unlock_type 用 Go 命名常量 —— **无宿主字典影响**。
- **缓存一致性**:本期未引入跨节点缓存,三榜为 DB 侧聚合(SUM+GROUP BY+Top-N),单节点正确;**预聚合/快照**留作规模优化(设计 D2) —— **无影响**。
- **数据权限**:玩家荣誉自隔离(`currentPlayerID`,只读不写 user_honor);三榜为公开活动读;荣誉配置受 `sicau-niu:honor:*`。本人名次用聚合子查询计数,**无 N+1**;玩家荣誉解锁批量预取计数后内存计算,**无逐荣誉查询**。
- **DI 来源检查**:`ranking`/`honor` service 在 `plugin.go registerRoutes` 一次性构造、显式注入 DAO/配置;无请求路径 `New()`。
- **数据库**:单迭代单 SQL(`005`,honor_def/user_honor),建表/索引幂等;排行榜聚合自 C4 feeding 无新表。
- **跨平台**:仅插件内 SQL/Go/Vue —— **无影响**。
- **修复记录**:honor/ranking 测试 seed 修正(activation 每行用不同自然日;user 每行用唯一 openid),贴合一机一号/每日一激活约束;产品聚合/解锁逻辑正确。
