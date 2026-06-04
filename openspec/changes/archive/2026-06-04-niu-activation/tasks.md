## 1. 数据库与 DAO

- [x] 1.1 新增 `manifest/sql/003-sicau-niu-activation.sql`:激活表 `plugin_sicau_niu_activation`(user_id、niu_id、activity_date、activated_at、is_first、order_no、photo_path、软删/时间)
- [x] 1.2 索引/约束:`(user_id,activity_date)` 活跃唯一、`(user_id,niu_id)` 活跃唯一、`niu_id` 索引;补 `manifest/sql/uninstall/003-*`(幂等)
- [x] 1.3 执行 `make dao p=sicau-niu` 生成激活表 DAO/DO/Entity

## 2. 配置

- [x] 2.1 `manifest/config` 新增 LBS 判距阈值(米,默认 50)、海报相关配置;`config.example.yaml` 同步

## 3. API 契约与代码生成

- [x] 3.1 定义玩家面 `api/player` DTO:可见牛列表 `GET /player/niu`、激活 `POST /player/activations`、图鉴 `GET /player/cards`、海报 `GET /player/poster`(RESTful、时间 Unix 毫秒、完整 dc/eg;玩家面不加 permission 标签,走 player-auth)
- [x] 3.2 执行 `make ctrl p=sicau-niu` 生成 controller 骨架与 api 接口

## 4. 服务层 — 地图与激活(`niu-map-activation`)

- [x] 4.1 可见牛列表:放出节奏可见性计算(阶段/上线时间/周几时段)+ DB 侧过滤 + 批量装配共享池态与本人激活态(杜绝 N+1)
- [x] 4.2 LBS 判距:Haversine(玩家经纬 vs 牛锚点)≤ 阈值;超出返回 bizerr;不做图像识别
- [x] 4.3 激活事务:牛行 `LockUpdate()` → 统计已激活数定首发/序号 → 首发置牛 status=active → 插入激活记录
- [x] 4.4 每日限次 `(user_id,activity_date)` + 同牛不重复 `(user_id,niu_id)`:先查后插 + 约束兜底,冲突返回 bizerr
- [x] 4.5 发卡:激活响应返回该牛主卡(无主卡时为空)
- [x] 4.6 `*_code.go` 集中 bizerr 码

## 5. 服务层 — 图鉴与海报(`niu-collection-poster`)

- [x] 5.1 个人图鉴:激活记录 JOIN 卡片(本人),按分类过滤;DB 批量,玩家自隔离
- [x] 5.2 海报数据:昵称/身份标签/第 N 只(序号)/随机金句/校庆标识;未激活该牛返回 bizerr
- [x] 5.3 海报 PNG 出图接缝 `PosterRenderer`(基础实现/可 Mock)
- [x] 5.4 `*_code.go` 集中 bizerr 码

## 6. 路由与依赖注入

- [x] 6.1 `plugin.go` 在玩家面路由组(C1 player-auth 中间件)绑定激活/图鉴/海报控制器
- [x] 6.2 controller `NewV1()` 持有 service 字段;service 构造函数显式注入 DAO 与跨能力(牛/卡片读取),不在请求路径 `New()`;玩家 ID 由 `playerctx` 提供

## 7. 测试与验证门禁

- [x] 7.1 服务层单测:LBS 判距边界、首发并发唯一(并发激活)、每日限次、同牛不重复、发卡、图鉴自隔离、海报数据(纯逻辑 + DB 门控,自包含含清理)
- [x] 7.2 后端编译门禁:plugins workspace 构建 + `go vet` + `go test`
- [x] 7.3 玩家面 API 冒烟:登录玩家 token → 可见牛列表/激活/图鉴/海报 主路径(无运营页,无浏览器 E2E;在任务记录说明)
- [x] 7.4 运行 `openspec validate niu-activation --strict`
- [x] 7.5 见下方「影响判断」

## 影响判断(7.5)

- **i18n**:单语言中文,无 i18n 块、无 `$t()` —— **无影响**。
- **字典**:无新增枚举,复用 C2 `cattle.NiuStatus`、`card.Category` Go 常量 —— **无宿主字典影响**。
- **缓存一致性**:本期未引入跨节点缓存,可见牛/共享池态用 DB 批量读取,单节点正确 —— **无影响**(集群快照留作后续)。
- **数据权限**:激活/图鉴/海报为玩家自隔离(按 `currentPlayerID` 限定);可见牛列表与共享池激活态为公开活动读。批量 `WhereIn` 装配,**无 N+1**。
- **DI 来源检查**:`activation.Service` 在 `plugin.go registerRoutes` 一次性构造,注入已构造的 identity service + `PosterRenderer` + LBS 阈值纯值配置;无请求路径 `New()`。
- **数据库**:单迭代单 SQL(`003`),建表/索引幂等;每日限次/同牛不重复用部分唯一索引;首发并发用牛行 `LockUpdate` 事务串行化。
- **跨平台**:仅插件内 SQL/Go/配置 —— **无影响**。
- **修复记录**:`activation_activate_test.go` 读 niu status 改为 `.Value().String()`(原 DO `interface{}`.(string) 误读为空);产品 flip 逻辑正确。
