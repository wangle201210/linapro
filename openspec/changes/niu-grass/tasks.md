## 1. 数据库与 DAO

- [x] 1.1 新增 `manifest/sql/004-sicau-niu-grass.sql`:建表 `grass_account`(user_id 活跃唯一、balance)、`grass_txn`(user_id、delta、txn_type、ref_id、时间)
- [x] 1.2 建表 `checkin`(user_id+checkin_date 活跃唯一、amount)、`feeding`(user_id、niu_id、base_amount、coefficient_basis、effect_amount、is_iron_bonus、fed_at)
- [x] 1.3 建表 `steal`(actor/target/amount/steal_date)、`gift`(from/to/amount/gift_date)、`inbox_msg`(user_id、msg_type、content、is_read、时间);索引;`uninstall/004`(幂等)
- [x] 1.4 执行 `make dao p=sicau-niu` 生成 DAO/DO/Entity

## 2. 配置

- [x] 2.1 `manifest/config` 新增:签到草量区间(20–50)、偷草名单人数(12)/每日次数、送草每日次数(12)/单次下限(12)、铁牛加成距离阈值(12m)、铁牛定位 Mock 开关;系数 1.5 为代码常量;`config.example.yaml` 同步

## 3. API 契约与代码生成

- [x] 3.1 定义玩家面 `api/player` DTO:签到、草账户、喂草(POST)+轨迹(GET)、可偷名单、偷草、送草、站内信列表+标记已读(RESTful、时间 Unix 毫秒、完整 dc/eg;走 player-auth 无 permission 标签)
- [x] 3.2 执行 `make ctrl p=sicau-niu` 生成 controller 骨架与 api 接口

## 4. 服务层 — 草账户与签到(`niu-grass-account`)

- [x] 4.1 账本式草账户:余额+流水,**任意增减同事务写流水+改余额**;扣减余额不足 bizerr;查询余额+近期流水(本人)
- [x] 4.2 每日签到:`(user_id,date)` 唯一,随机 20–50(`grand`),入账;重复签到 bizerr
- [x] 4.3 `*_code.go` bizerr 码;账户增减提供可复用的内部记账方法(供喂草/偷草/送草调用)

## 5. 服务层 — 喂草与铁牛加成(`niu-feeding`)

- [x] 5.1 喂草:仅已激活牛(C3 status);校验余额→扣 base(feed 流水);返回牛信息+随机金句
- [x] 5.2 铁牛定位接缝 `IronLocationGateway`(Mock 实现,返回铁牛 last_lat/lng 或配置坐标);真实实现占位
- [x] 5.3 加成判定:Haversine(牛锚点 vs 各铁牛当前位置)<12m → 系数 1.5;effect=base×系数;记录 base/coefficient_basis/effect/is_iron_bonus
- [x] 5.4 最近喂养轨迹(最近 10,本人);`*_code.go` bizerr 码

## 6. 服务层 — 偷草/送草/站内信(`niu-grass-social`)

- [x] 6.1 每日可偷名单:确定性随机(种子=玩家+日)抽 12 人(其他玩家);偷草校验目标在名单内
- [x] 6.2 偷草:每日次数上限;同事务 目标 stolen_loss / 本人 steal_gain;被偷站内信
- [x] 6.3 送草:每日 12 次/每次 ≥12 份/余额校验;同事务 gift_out/gift_in;收草站内信
- [x] 6.4 站内信:本人列表(分页)+标记已读(校验归属);类型 Go 常量;`*_code.go` bizerr 码

## 7. 路由与依赖注入

- [x] 7.1 `plugin.go` 玩家面路由组(C1 player-auth)绑定签到/账户/喂草/轨迹/偷草/送草/站内信控制器
- [x] 7.2 controller `NewV1()` 持有 service 字段;service 构造函数显式注入 DAO 与跨能力(账户记账、铁牛定位、激活态/牛/金句读取),不在请求路径 `New()`;玩家 ID 由 `playerctx` 提供

## 8. 测试与验证门禁

- [x] 8.1 服务层单测:账本一致性(余额=流水累加)、余额不足拒绝、签到限次、喂草加成判距/系数/effect、给未激活牛拒绝、偷草名单/次数/通知/双方事务、送草限额/事务/通知、站内信归属与已读(纯逻辑 + DB 门控,自包含含清理)
- [x] 8.2 后端编译门禁:plugins workspace 构建 + `go vet` + `go test`
- [x] 8.3 玩家面 API 冒烟:登录玩家 token → 签到/账户/喂草/偷草名单/送草/站内信 主路径(无运营页,无浏览器 E2E,记录说明)
- [x] 8.4 运行 `openspec validate niu-grass --strict`
- [x] 8.5 见下方「影响判断」

## 影响判断(8.5)

- **i18n**:单语言中文,无 i18n 块/`$t()`;站内信内容为运行期用户内容 —— **无影响**。
- **字典**:草流水类型(checkin/feed/steal_gain/stolen_loss/gift_out/gift_in)、站内信类型(stolen/gift_received)用 Go 命名常量 —— **无宿主字典影响**。
- **缓存一致性**:本期未引入跨节点缓存;草账户余额/流水为权威 DB,所有增减事务内一致;偷草/送草双方在同一事务 —— **无影响**。
- **数据权限**:账户/签到/喂草/轨迹/送出/站内信玩家自隔离(`currentPlayerID`);偷草仅经"当日可偷名单"授权,只改对方草账户 + 发其站内信,不暴露对方其他隐私。批量 `WhereIn` 装配,**无 N+1**。
- **DI 来源检查**:`grass`/`feeding`/`grasssocial` service 在 `plugin.go registerRoutes` 一次性构造,显式注入账户记账能力、`IronLocationGateway`(Mock)、cattle/quote 读取、配置纯值;无请求路径 `New()`。
- **数据库**:单迭代单 SQL(`004`,7 表),建表/索引幂等;账本余额改动一律与流水同事务;签到/账户唯一约束;铁牛实时经纬复用 C2 iron 表。
- **跨平台**:仅插件内 SQL/Go/配置 —— **无影响**。
