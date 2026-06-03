## Why

C4 产生了喂草行为数据,C1 提供了身份标签。C5 把这些数据转化为**荣誉与竞争**:个人喂草榜、院系榜、川农好友榜驱动持续参与;参与徽章、分类徽章、电子证书、头像框等荣誉提供达成激励。这是活动"动机与奖励"层,为 C7 的结算与批量发证提供荣誉数据。

本期交付面向小程序的玩家榜单/荣誉 API 与运营荣誉配置后台;不实现小程序 UI。

## What Changes

- 新增**排行榜(M8)**:个人喂草榜(按喂草实际效果累计)、院系榜(按身份标签的院系聚合)、川农好友榜(社会人员/非在校玩家)。三榜均有 Top-N 上限,玩家可看自己名次。
- 新增**荣誉定义与玩家荣誉(M11)**:
  - 运营后台**荣誉配置 CRUD**:荣誉类型(徽章/头像框/证书)、名称、解锁规则(参与即得 / 喂草次数阈值 / 激活数阈值 / 集齐分类 / 集齐全套)、图片/模板、排序。
  - 玩家**荣誉展示**:基于其喂草次数、激活数、图鉴集齐情况,计算各荣誉的解锁状态(只读展示)。
- 荣誉的持久发放与批量发证留待 C7(结算)。

不在本 change 范围:看板/导出/批量发证(C7);H5(C6);喂草数据产生(C4 已做)。

## Capabilities

### New Capabilities
- `niu-ranking`: 个人喂草榜、院系榜、川农好友榜的聚合与名次查询。
- `niu-honor`: 荣誉定义运营 CRUD 与玩家荣誉解锁状态展示。

### Modified Capabilities
<!-- 无:复用 C1 身份、C3 激活/图鉴、C4 喂草数据,不改其需求级行为。 -->

## Impact

- **后端**:`apps/lina-plugins/sicau-niu/backend/`
  - 新增 `manifest/sql/005-sicau-niu-honor.sql`:荣誉定义 `honor_def`、玩家荣誉 `user_honor`(供 C7 持久发放);`uninstall/005`。排行榜为聚合派生,不新建权威表。
  - 新增 `api/player`(三榜、玩家荣誉)、`api/admin`(荣誉配置 CRUD)、`internal/controller`、`internal/service`、`internal/dao`+`model`。
  - `plugin.yaml`:「寻牛活动」下新增荣誉配置菜单与 `sicau-niu:honor:*` 权限。
  - `plugin.go`:玩家面绑定三榜/荣誉;运营面绑定荣誉配置。
- **数据库**:新增 2 张表;排行榜用对 `feeding` 的聚合查询(`GROUP BY user_id SUM(effect_amount)`)+ 关联 `user`,Top-N 限制,依赖 feeding/user 索引。
- **接口性能**:榜单为高频读;C5 用数据库侧聚合 + Top-N 上限;**预聚合/快照**作为规模优化在设计记录并留待后续(集群下遵循 `cache-consistency`)。
- **i18n**:单语言(中文),不启用 —— 无 i18n 影响。
- **验证**:服务层 DB 门控单测(三榜聚合/名次、荣誉配置、玩家荣誉解锁计算)+ 运营荣誉配置页 E2E + 玩家 API 冒烟。
