## Why

寻牛活动收尾需要运营结算能力(M5 看板 / M13 风控告警 / M14 公示结算):运营在管理端总览活动数据、导出名册、把荣誉证书批量发放给达标玩家、查看风控告警(一机多号),并在活动结束时归档一份结算快照用于公示与留存。C7 复用 C1 玩家、C2 牛、C3 激活、C4 喂草/偷草/送草/签到、C5 荣誉,是面向运营的最后一块拼图。

## What Changes

- 新增**运营数据看板**(管理端,鉴权):集合化聚合的活动总览(玩家数、激活牛数、首发数、喂草次数与总效果、偷草/送草/签到次数、已发证书数等),全部 DB 侧聚合,杜绝 N+1。
- 新增**名册导出**(管理端,鉴权):有界、批量装配的玩家名册投影(昵称/身份/院系/年级/激活数/喂草总量),供运营导出 CSV;有上限并记录截断。
- 新增**证书批量发放**(管理端,鉴权):对指定证书类荣誉(`honor_type=certificate`),按其解锁规则集合化筛出达标玩家队列,幂等写入 `user_honor`(命中既有授予即跳过);返回 达标/新发/跳过 计数。集合规则支持 参与/喂草次数/激活次数;图鉴完成类(category_complete/full_complete)按个人在小程序内解锁,不纳入批量结算并以业务错误明确拒绝。
- 新增**风控告警视图**(管理端,鉴权):集合化列出共享设备指纹的玩家簇(一机多号),含簇内昵称,供人工研判。只读。
- 新增**结算归档**(管理端,鉴权):把看板指标冻结为一份结算快照持久化到新表 `plugin_sicau_niu_settlement`(标题/快照 JSON/操作人/归档时间),并提供归档列表查询。
- 顶级"寻牛活动"菜单下新增"运营结算"页面与按钮权限(`sicau-niu:settlement:*`)。

不在本 change 范围:小程序 UI;皮肤;赠送撤销;把风控告警持久化为独立事件表(本期为派生只读视图)。

## Capabilities

### New Capabilities
- `niu-settlement`: 运营数据看板、名册导出、证书批量幂等发放、风控告警(一机多号)派生视图、结算快照归档与查询。

### Modified Capabilities
<!-- 无需求级行为变更:复用 C1–C5 既有表与契约,仅新增运营聚合/导出/发证/归档读写路径。 -->

## Impact

- **数据库**:新增 `manifest/sql/006-sicau-niu-settlement.sql`(表 `plugin_sicau_niu_settlement`,含 `deleted_at IS NULL` 软删除与必要索引)+ `manifest/sql/uninstall/006-*.sql`;`backend/hack/config.yaml` `tables` 追加该表并 `make dao`。
- **后端**:`apps/lina-plugins/sicau-niu/backend/`
  - 新增 `api/settlement`、`internal/controller/settlement`、`internal/service/settlement`(看板/导出/发证/风控/归档)。
  - `plugin.go`:管理端子组(Auth+Tenancy+Permission)绑定结算控制器;装配期显式注入。
  - 发证复用 C5 `honor_def`/`user_honor`;聚合复用 C3/C4/C5 表;**禁止 N+1**,集合化查询。
- **前端**:管理端新增"运营结算"页(看板卡片 + 导出按钮 + 批量发证 + 风控告警表 + 归档列表/创建),数据测试锚点,中文,无 i18n。
- **数据权限**:管理端动作均经宿主 `Auth+Tenancy+Permission`,DTO 带 `permission:"sicau-niu:settlement:*"`;发证为写操作,幂等且记录授予;风控/导出为运营只读。`platform_only`。
- **i18n**:单语言中文,不启用 —— 无 i18n 影响。
- **缓存一致性**:看板为 DB 侧实时聚合,本期不引入缓存/快照失效路径;归档快照为一次性持久化记录 —— 无新增缓存一致性影响。
- **验证**:服务层 DB 门控单测(看板聚合、发证幂等与队列筛选、风控簇、归档持久化)+ 编译门禁 + `openspec validate --strict`;管理端 E2E 并入最终统一阶段。
