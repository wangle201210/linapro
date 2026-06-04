## Context

C3 接通寻牛核心玩法(见 `apps/lina-plugins/sicau-niu/design.md` 程序级设计 C3)。复用 C1 玩家身份(player-auth 中间件 + `playerctx`)与 C2 牛/卡片数据。最难点是**共享池首发并发一致性**:120 头牛任一玩家首发激活后全员可见,两人同抢一牛只能一个首发。

约束:遵循 `.agents/rules/` 的 `plugin`、`backend-go`、`api-contract`、`database`、`data-permission`、`architecture`、`cache-consistency`、`testing`;单语言中文不启用 i18n;插件 `platform_only`。

## Goals / Non-Goals

**Goals:**
- 玩家可拉取当前可见牛(按放出节奏过滤),含共享池激活态与本人激活态。
- 玩家到牛锚点附近(服务端 LBS 判距)可激活;每人每天 1 次;首发并发安全;激活发卡。
- 玩家可看个人图鉴(已激活牛卡片)与激活海报数据。

**Non-Goals:**
- 喂草/签到/偷草/送草(C4);铁牛实时定位与加成(C4);排行榜/荣誉(C5);H5(C6);看板/结算(C7);图像识别;小程序 UI。

## Decisions

### D1 共享池首发并发用事务 + 牛行锁
激活在 `dao.Niu.Transaction` 内:对目标牛行 `LockUpdate()`(`SELECT ... FOR UPDATE`)→ 统计该牛已激活数 `n` → 本次 `orderNo=n+1`、`isFirst=(n==0)` → 若首发则将牛 `status` 从 `inactive` 置 `active` → 插入激活记录。行锁串行化同一头牛的并发激活,保证唯一首发。
- **备选**:仅靠唯一约束 —— 不足以原子地同时定首发+序号+牛状态,否决。

### D2 LBS 服务端判距,不做图像识别
服务端用 Haversine 计算玩家上报经纬与牛锚点距离,`≤` 阈值(配置,默认 50m)才允许激活;超出返回 `bizerr`。拍照路径仅记录(取证),不做识别。

### D3 每日限次 + 同牛不重复
`plugin_sicau_niu_activation` 唯一约束:`(user_id, activity_date)` 活跃唯一(每日 1 次)、`(user_id, niu_id)` 活跃唯一(同牛不重复)。服务层先查后插,并以约束兜底。

### D4 发卡与图鉴为派生
激活即"获得"该牛主卡:不新建领取表;**个人图鉴 = 激活记录 JOIN 卡片**(`activation.user_id=me` → 取各 `niu_id` 的卡片),按分类浏览。激活响应直接返回该牛主卡。

### D5 可见牛列表批量装配
玩家可见牛 = 放出节奏命中(`release_stage` 非空且 `online_at<=now`,可选周几/时段匹配)。列表 DB 侧过滤后,**批量**装配:一次查这批牛的共享池态(牛自带 `status`)、一次查本人对这批牛的激活记录(`WhereIn niu_id AND user_id=me`),内存合并出 `activatedByMe`。≤120 头,无 N+1。

### D6 海报数据 + 出图接缝
`GET /player/poster` 返回海报合成数据:昵称、身份标签(在校生/校友/川农好友)、第 N 只牛(到场序号)、随机校史金句、校庆标识。**服务端 PNG 出图**以 `PosterRenderer` 接缝提供:C3 给出基础实现(可 Mock/占位),按设计稿的精细排版与 CJK 字体渲染留待后续替换,接口边界稳定。

### D7 数据权限(玩家自隔离)
激活、图鉴、海报为**本人**数据(`playerctx` 注入的 playerID 限定);可见牛列表与共享池激活态为公开读(全员可见的活动状态,非隐私)。

### D8 缓存
可见牛/共享池态为有界数据(≤120),C3 用 DB 批量读取;集群下若引入快照/缓存,激活时失效,遵循 `cache-consistency`。本期不引入跨节点缓存,记录无缓存一致性影响(单节点正确)。

### D9 显式依赖注入
service 构造函数显式注入 DAO 与跨能力(牛/卡片读取);controller 持有 service 字段;玩家 ID 由 C1 `playerctx` 提供;不在请求路径 `New()`。

### 数据模型(C3)
| 表 | 关键列 | 约束/索引 |
|----|-------|-----------|
| `plugin_sicau_niu_activation` | `user_id`、`niu_id`、`activity_date`、`activated_at`、`is_first`、`order_no`、`photo_path`、软删/时间 | `(user_id,activity_date)` 活跃唯一;`(user_id,niu_id)` 活跃唯一;`niu_id` 索引 |

牛共享池态复用 C2 `plugin_sicau_niu_niu.status`(事务内原子更新)。

### API 契约(C3,玩家面,player-auth,时间 Unix 毫秒)
| 方法 | 路径 | 说明 |
|------|------|------|
| `GET` | `/plugins/sicau-niu/player/niu` | 可见牛列表(地图),含共享池态 + 本人激活态 |
| `POST` | `/plugins/sicau-niu/player/activations` | 激活(body: niuId, lat, lng, photoPath?),返回首发/序号/主卡 |
| `GET` | `/plugins/sicau-niu/player/cards` | 个人图鉴(已激活牛卡片,可按分类过滤) |
| `GET` | `/plugins/sicau-niu/player/poster` | 激活海报合成数据(query: niuId) |

## Risks / Trade-offs

- [并发首发正确性] → 牛行 `FOR UPDATE` 串行化 + 唯一约束兜底;单测以并发激活验证唯一首发。
- [LBS 防造假] → C3 做服务端判距;虚拟定位/漂移的进阶风控属 C7(风控),此处仅基础判距。
- [海报 PNG 精细渲染] → 以接缝 + 基础实现交付,设计稿级渲染后续替换,不阻塞核心激活闭环。
- [activity_date 时区] → 以服务端时区的自然日为准,`activity_date` 存 `YYYY-MM-DD` 字符串,日界清晰。

## Migration Plan

1. 新增 `manifest/sql/003-sicau-niu-activation.sql`(激活表+索引)与 `uninstall/003`;`make dao`。
2. `manifest/config`:LBS 阈值、海报配置。
3. 新增 `api/player` DTO → `make ctrl` → 填 controller/service;`plugin.go` 玩家面绑定。
4. 服务层单测(LBS、并发首发、每日限次、同牛不重复、发卡、图鉴、海报数据)+ 编译门禁 + API 冒烟;`openspec validate --strict`。
5. 回滚:代码与 SQL 可整体回退;无生产数据。

## Open Questions

1. LBS 判距阈值最终值(默认 50m,运营可配)。
2. 海报 PNG 是否本期就要设计稿级渲染,还是先数据 + 基础出图(默认后者)。
3. 拍照是否需要留存图片(默认记录路径即可,图片走宿主文件上传或小程序本地)。

> 影响判断:i18n 单语言不启用,**无影响**;缓存一致性 —— 本期无跨节点缓存,单节点正确,**无影响**(集群快照留作后续);数据权限 —— 激活/图鉴/海报为玩家自隔离,牛列表/共享池态为公开读;跨平台 —— 仅插件内 SQL/代码,**无影响**;字典 —— 无新增枚举(复用 C2 卡片分类常量)。
