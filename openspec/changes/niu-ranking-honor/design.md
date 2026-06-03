## Context

C5 把 C4 喂草数据与 C1 身份转化为榜单与荣誉(见程序级设计 C5)。复用 C4 `feeding.effect_amount`、C1 `user`(身份/院系)、C3 激活/图鉴。

约束:遵循 `.agents/rules/` 的 `plugin`、`backend-go`、`api-contract`、`database`、`data-permission`、`architecture`、`cache-consistency`、`frontend-ui`、`testing`;单语言中文不启用 i18n;`platform_only`。

## Goals / Non-Goals

**Goals:** 三榜(喂草/院系/川农好友)聚合 + 本人名次;运营荣誉配置 CRUD;玩家荣誉解锁状态展示(基于喂草次数/激活数/集卡)。

**Non-Goals:** 看板/导出/批量发证(C7);荣誉持久发放(C7,C5 仅展示计算);H5(C6);小程序 UI。

## Decisions

### D1 排行榜为聚合派生,不建权威表
喂草榜 = `feeding GROUP BY user_id SUM(effect_amount)` 倒序 Top-N,关联 `user` 取昵称。院系榜 = 在校生(identity=student)按 `college_id` 聚合各院系喂草总效果。川农好友榜 = `identity=friend` 玩家按个人喂草效果排名。均 Top-N 上限;本人名次单独查(本人聚合值 + 高于本人的计数 + 1)。

### D2 榜单性能与缓存
榜单为高频读。C5 用**数据库侧聚合 + Top-N 上限**(依赖 `feeding.user_id` 索引、`user` 关联),不在内存加载全量。**预聚合/快照**(定时或事件刷新的榜单快照)作为规模优化记录于设计,本期不引入跨节点缓存(单节点聚合查询正确);集群下若引入快照,刷新遵循 `cache-consistency`。记录:本期无新增缓存。

### D3 荣誉定义运营 CRUD
`honor_def`:`honor_type`(badge/avatar_frame/certificate)、`code`(活跃唯一)、`name`、`unlock_type`(participation/feed_count/activation_count/category_complete/full_complete)、`threshold`(int,阈值类用)、`category`(分类集齐类用,卡片分类)、`image_path`、`sort`。运营后台 CRUD,受 `sicau-niu:honor:*` 权限。`honor_type`/`unlock_type` 用 Go 命名常量。

### D4 玩家荣誉解锁为只读计算
玩家荣誉展示对每个 `honor_def` 计算是否解锁:
- `participation`:已注册即解锁;
- `feed_count`:玩家喂草次数 ≥ threshold;
- `activation_count`:玩家激活数 ≥ threshold;
- `category_complete`:玩家图鉴中某分类已集齐(该分类全部已上线卡片均已获得);
- `full_complete`:集齐全部已上线主卡。
计算所需计数用批量聚合(喂草次数、激活数、按分类的已得卡片数 vs 该分类卡片总数),**杜绝逐荣誉 N+1**。C5 为只读展示(GET 无副作用);荣誉的持久发放(写 `user_honor`)与批量发证留待 C7。

### D5 数据权限
玩家荣誉为本人(`currentPlayerID`);榜单为公开活动读(展示昵称/院系聚合,不泄露隐私);荣誉配置为运营受权限全量。

### D6 显式依赖注入
service 构造函数显式注入 DAO 与跨能力(读 feeding/activation/card/user 计数);controller 持 service 字段;不在请求路径 `New()`。

### 数据模型(C5)
| 表 | 关键列 | 约束/索引 |
|----|-------|-----------|
| `plugin_sicau_niu_honor_def` | `honor_type`、`code`、`name`、`unlock_type`、`threshold`、`category`、`image_path`、`sort`、软删/时间 | `code` 活跃唯一 |
| `plugin_sicau_niu_user_honor` | `user_id`、`honor_id`、`unlocked_at`、软删/时间 | `(user_id,honor_id)` 活跃唯一(供 C7 发放) |

排行榜聚合自 C4 `feeding` + C1 `user`,无新表。

### API 契约(C5,时间 Unix 毫秒)
玩家面(player-auth):
| 方法 | 路径 | 说明 |
|------|------|------|
| `GET` | `/plugins/sicau-niu/player/rankings/feed` | 个人喂草榜 Top-N + 本人名次 |
| `GET` | `/plugins/sicau-niu/player/rankings/college` | 院系榜 Top-N |
| `GET` | `/plugins/sicau-niu/player/rankings/friend` | 川农好友榜 Top-N + 本人名次 |
| `GET` | `/plugins/sicau-niu/player/honors` | 玩家荣誉解锁状态列表 |

运营面(`Auth+Permission`):
| 方法 | 路径 | 权限 |
|------|------|------|
| `GET/POST` | `/plugins/sicau-niu/admin/honors`、`GET/PUT/DELETE /admin/honors/{id}` | `sicau-niu:honor:list/create/update/delete` |

菜单挂 C1「寻牛活动」顶级目录下。

## Risks / Trade-offs

- [榜单规模性能] → DB 聚合 + Top-N 上限;预聚合快照留作规模优化(设计记录),本期单节点聚合正确。
- [玩家荣誉计算成本] → 各计数用批量聚合(喂草次数/激活数/分类集卡),避免逐荣誉查询。
- [集齐判定依赖"已上线卡片"口径] → 以当前已存在主卡为分母;运营增删卡片影响集齐判定,属预期。

## Migration Plan

1. `manifest/sql/005-sicau-niu-honor.sql`(honor_def、user_honor)+ `uninstall/005`;`make dao`。
2. `api/player` + `api/admin` DTO → `make ctrl` → 填 controller/service;`plugin.yaml` 荣誉菜单;`plugin.go` 绑定。
3. 运营荣誉配置前端页(vxe-grid + 弹窗)。
4. 服务层 DB 门控单测 + 编译门禁 + 运营页 E2E + 玩家 API 冒烟;`openspec validate --strict`。
5. 回滚:代码与 SQL 可整体回退;无生产数据。

## Open Questions

1. 各榜 Top-N 默认值(默认 100)。
2. 喂草阈值徽章(青铜/白银/黄金/铂金)的阈值是配置(honor_def.threshold)还是固定(用 honor_def 配置)。
3. 荣誉持久发放时机(默认 C7 批量;C5 仅展示)。

> 影响判断:i18n 单语言不启用,**无影响**;字典 —— honor_type/unlock_type 用 Go 常量,**无宿主字典影响**;缓存一致性 —— 本期无跨节点缓存,榜单 DB 聚合,**无影响**(预聚合快照留后续);数据权限 —— 玩家荣誉自隔离,榜单公开读,荣誉配置受权限;跨平台 —— 仅插件内文件,**无影响**。
