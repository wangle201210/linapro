## Context

C4 接通寻牛的草经济与社交玩法(见程序级设计 C4)。复用 C1 玩家身份(player-auth+`playerctx`)、C2 铁牛(`plugin_sicau_niu_iron`)与金句、C3 激活态(只能给已激活牛喂草)。核心是**账本式草账户的事务一致性**与**铁牛实时定位加成**。

约束:遵循 `.agents/rules/` 的 `plugin`、`backend-go`、`api-contract`、`database`、`data-permission`、`architecture`、`cache-consistency`、`testing`;单语言中文不启用 i18n;`platform_only`。

## Goals / Non-Goals

**Goals:** 每日签到领草;账本式草账户(余额+流水,所有增减事务化);喂任意已激活牛(扣草、随机金句、最近轨迹、铁牛 `<12m` ×1.5 加成并记录原始量/系数/实际效果);偷草(每日随机 12 人名单、每日次数上限、被偷站内信);送草(每日 12 次/每次 ≥12 份、双方事务记账、收草站内信);站内信查询。

**Non-Goals:** 排行榜/荣誉计分(C5,仅产出喂草数据);H5(C6);看板/结算(C7);铁牛运营登记(C2 已做);小程序 UI。

## Decisions

### D1 账本式草账户
`grass_account`(玩家唯一,`balance`)+ `grass_txn`(追加式流水:`delta`、`txn_type`、`ref_id`、`created_at`)。**任何草增减都在事务内同时写流水 + 改余额**,保证可追溯可对账。`txn_type` 用 Go 命名常量:`checkin`/`feed`/`steal_gain`/`stolen_loss`/`gift_out`/`gift_in`。余额不足的扣减拒绝。

### D2 每日签到
`checkin` 唯一 `(user_id, checkin_date)`;每日 1 次随机 20–50 份(`grand`),入账为 `checkin` 流水。重复签到返回 bizerr。

### D3 喂草与铁牛加成
- 仅可喂**已激活**的牛(C3 niu.status=active);否则 bizerr。
- 玩家投喂 `baseAmount` 份草:校验余额 ≥ base → 扣 base(`feed` 流水)。
- **铁牛加成**:取各铁牛**当前经纬**(D4 接缝)与被喂牛锚点的 Haversine 距离,任一 `<12m` 则 `coefficient=1.5`(集中常量),否则 `1.0`。`effectAmount = base × coefficient`。
- 喂草记录保存 `baseAmount`、`coefficientBasis`(100 或 150)、`effectAmount`、`isIronBonus`。`effectAmount` 供 C5 榜单计分。
- 响应返回牛信息 + 随机校史金句(`quote` enabled=1 随机);最近喂养轨迹取最近 10 条。

### D4 铁牛定位外部接缝(Mock)
`IronLocationGateway.Locations(ctx) -> []{ironID, lat, lng, ok}`:Mock 实现返回铁牛 `last_lat/lng`(或配置注入的模拟坐标),真实实现调用外部 API。喂草加成判定按需调用,接口边界稳定后续替换。

### D5 偷草
- **每日随机 12 人名单**:对 `(user_id, 自然日)` 用确定性随机种子,从其他玩家中抽 12 人(无需存表,可复算);偷草时校验目标在请求者当日名单内。
- 偷草:从目标账户偷随机草量 → 目标 `stolen_loss`、本人 `steal_gain`(同事务);每人每日偷草次数上限(`steal` 按 `actor+date` 计数)。
- **被偷站内信**:给目标写 `inbox_msg`(被偷通知)。

### D6 送草
- 玩家将自有草赠送指定玩家:**每人每天最多 12 次、每次至少 12 份**(集中常量),余额不足拒绝。
- 双方在一个事务:本人 `gift_out` 扣减、对方 `gift_in` 增加。
- **收草站内信**:给接收方写 `inbox_msg`(收到赠草)。

### D7 站内信
`inbox_msg`(玩家、类型、内容、已读、时间)。玩家查询本人站内信(分页)、标记已读。类型 Go 常量(`stolen`/`gift_received`)。

### D8 数据权限(玩家自隔离)
账户/流水/签到/喂草/轨迹/送出/站内信均为**本人**(`currentPlayerID`);偷草目标是他人但经"当日可偷名单"授权,且仅能改对方草账户与发其站内信,不暴露对方其他隐私;被偷/收草通知只对接收者可见。

### D9 事务一致性
草增减与流水耦合写入(D1);偷草、送草涉及双方账户的扣/加在**同一事务**内完成,失败整体回滚。

### 数据模型(C4)
| 表 | 关键列 | 约束/索引 |
|----|-------|-----------|
| `plugin_sicau_niu_grass_account` | `user_id`、`balance` | `user_id` 活跃唯一 |
| `plugin_sicau_niu_grass_txn` | `user_id`、`delta`、`txn_type`、`ref_id`、时间 | `user_id` 索引(追加式) |
| `plugin_sicau_niu_checkin` | `user_id`、`checkin_date`、`amount` | `(user_id,checkin_date)` 活跃唯一 |
| `plugin_sicau_niu_feeding` | `user_id`、`niu_id`、`base_amount`、`coefficient_basis`、`effect_amount`、`is_iron_bonus`、`fed_at` | `user_id`、`niu_id` 索引 |
| `plugin_sicau_niu_steal` | `actor_user_id`、`target_user_id`、`amount`、`steal_date` | `(actor_user_id,steal_date)` 索引 |
| `plugin_sicau_niu_gift` | `from_user_id`、`to_user_id`、`amount`、`gift_date` | `(from_user_id,gift_date)` 索引 |
| `plugin_sicau_niu_inbox_msg` | `user_id`、`msg_type`、`content`、`is_read`、时间 | `user_id` 索引 |

铁牛实时经纬复用 C2 `plugin_sicau_niu_iron.last_lat/lng/located_at`。

### API 契约(C4,玩家面,player-auth,时间 Unix 毫秒)
| 方法 | 路径 | 说明 |
|------|------|------|
| `POST` | `/plugins/sicau-niu/player/checkin` | 每日签到领草 |
| `GET` | `/plugins/sicau-niu/player/grass` | 草余额 + 近期流水 |
| `POST` | `/plugins/sicau-niu/player/feedings` | 喂草(body: niuId, baseAmount),返回牛信息+金句+加成 |
| `GET` | `/plugins/sicau-niu/player/feedings` | 最近喂养轨迹(最近 10) |
| `GET` | `/plugins/sicau-niu/player/steal-targets` | 当日可偷 12 人名单 |
| `POST` | `/plugins/sicau-niu/player/steals` | 偷草(body: targetUserId) |
| `POST` | `/plugins/sicau-niu/player/gifts` | 送草(body: toUserId, amount) |
| `GET` | `/plugins/sicau-niu/player/messages` | 站内信列表(分页) |
| `PUT` | `/plugins/sicau-niu/player/messages/{id}/read` | 标记已读 |

## Risks / Trade-offs

- [账本一致性] → 余额改动一律与流水同事务;单测以扣减/转账验证余额=流水累加。
- [偷草名单可复算 vs 存表] → 确定性随机(种子=玩家+日)免存表;若需审计可后续落表。
- [铁牛加成实时性] → 喂草时按需取铁牛当前位置(Mock);真实接口频率/鉴权见开放问题。
- [并发账户操作] → 同玩家并发扣减用账户行锁或乐观校验,避免超扣;单测覆盖。

## Migration Plan

1. `manifest/sql/004-sicau-niu-grass.sql`(7 表+索引)+ `uninstall/004`;`make dao`。
2. `manifest/config`:签到区间、偷草人数/次数、送草限额、铁牛加成阈值/系数(系数常量在代码)、铁牛定位 Mock 开关。
3. `api/player` DTO → `make ctrl` → 填 controller/service;`ironlocation` 接缝;`plugin.go` 玩家面绑定。
4. 服务层 DB 门控单测 + 编译门禁 + API 冒烟;`openspec validate --strict`。
5. 回滚:代码与 SQL 可整体回退;无生产数据。

## Open Questions

1. 喂草单次投喂量是固定还是玩家可选(默认玩家传 baseAmount,设上下限)。
2. 偷草单次草量与每日次数上限的默认值(默认:随机小额,每日 5 次)。
3. 铁牛定位真实 API 形态(Mock 不阻塞)。

> 影响判断:i18n 单语言不启用,**无影响**;字典 —— 流水/站内信类型用 Go 常量,**无宿主字典影响**;缓存一致性 —— 本期无跨节点缓存,账户/流水事务内一致,**无影响**;数据权限 —— 玩家自隔离 + 偷草经名单授权;跨平台 —— 仅插件内文件,**无影响**。
