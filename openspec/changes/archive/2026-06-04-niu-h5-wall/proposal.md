## Why

寻牛活动需要一个对外传播载体:**H5 数字纪念墙**。它无需登录即可访问,展示首发激活的 120 人、校史精选与活动荣誉/统计,并可扫码回跳小程序参与,把活动从小程序内向外站传播。复用 C3 首发激活数据、C1 玩家、C2 卡片/金句。

H5 在本仓库交付(插件公开静态资源 + 公开只读 API)。不实现小程序 UI。

## What Changes

- 新增**公开只读纪念墙 API**(无需登录,免鉴权):
  - 首发墙:首发激活的玩家名单(昵称、身份标签、所属牛名称/序号、到场顺序),按序号有界返回。
  - 校史精选:从已上线卡片/校史金句中取样展示。
  - 活动统计:已激活牛数、参与玩家数等公开统计。
- 新增**H5 数字纪念墙静态页**:通过插件 `public_assets` 托管,移动端展示首发墙/校史精选/统计,并提供**扫码/链接回跳小程序**入口。
- 公开 API 不返回手机号、openid 等隐私字段;仅展示昵称与活动相关信息。

不在本 change 范围:看板/导出/结算(C7);荣誉持久发放(C7);小程序 UI;首发数据产生(C3 已做)。

## Capabilities

### New Capabilities
- `niu-memorial-wall`: 公开免登录的数字纪念墙数据 API(首发墙/校史精选/统计)与 H5 静态页(含扫码回跳小程序)。

### Modified Capabilities
<!-- 无:复用 C3 首发激活、C1 玩家昵称、C2 卡片/金句,不改其需求级行为。 -->

## Impact

- **后端**:`apps/lina-plugins/sicau-niu/backend/`
  - 新增 `api/wall`(首发墙/校史精选/统计 公开 DTO)、`internal/controller/wall`、`internal/service/wall`。
  - `plugin.go`:在**公开子路由组**(无 Auth/无 player-auth,与玩家登录同级)绑定纪念墙 API。
  - 无新增数据库表(数据派生自 C3 `activation`(is_first)+ C1 `user` + C2 `card`/`quote`)。
- **公开静态资源**:新增 `frontend/wall/index.html`(H5 纪念墙页),`plugin.yaml` 声明 `public_assets`(托管到 `/x-assets/sicau-niu/{version}/`)。
- **数据权限**:纪念墙为公开活动展示;**只暴露昵称与活动信息,绝不暴露手机号/openid 等隐私**;有界返回(首发 ≤120、精选/统计有上限)。
- **i18n**:单语言(中文),不启用 —— 无 i18n 影响。
- **验证**:服务层 DB 门控单测(首发墙派生、隐私字段不外泄、统计)+ 公开 API 冒烟(无 token 可访问);H5 页面为公开静态资源(统一 E2E 阶段验证)。
