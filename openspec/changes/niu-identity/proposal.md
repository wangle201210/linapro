## Why

寻牛活动面向所有人（在校生、校友、社会人员）开放，手机号即可注册。后续地图激活、喂草/偷草/送草、排行榜、荣誉等全部玩法都依赖一个统一的玩家账户与身份体系作为数据与权限地基。C1 先把「登录 + 身份 + 资料 + 院系字典 + 一机一号」这层地基打稳，让后续 change 可以直接挂接玩家身份。

本期先实现面向微信小程序的后端 API，不实现小程序 UI（小程序 App 由本团队后续交付）。

## What Changes

- 新增寻牛**玩家账户体系**：微信小程序登录（`code` 换 openid）+ 手机号一键授权登录，后端签发会话 token。
- **用户建档与身份资料**：昵称、身份标签（在校生 / 校友 / 川农好友）、院系（字典选择）、年级（用户填数字）、毕业年（校友填数字）。
- **院系字典 + 运营后台录入入口**：院系清单由运营在后台手动维护，供玩家身份资料选择。
- **轻量风控基础**：手机号唯一（一机一号），记录设备指纹，作为后续异常风控的基础数据。
- **插件作用域调整**：`sicau-niu` 由 `tenant_aware` 调整为 `platform_only` / `supports_multi_tenant=false` / `default_install_mode=global`（本活动不涉及多租户）。
- 用真实身份能力**取代当前示例占位代码**（`niu` 牛只示例列表）。

不在本 change 范围：地图/激活/喂草/偷草/送草/排行榜/荣誉/海报/H5/看板（见各自后续 change）。

## Capabilities

### New Capabilities
- `niu-player-identity`: 寻牛玩家登录（微信 / 手机号）、会话签发、用户建档、身份标签与资料（院系/年级/毕业年）维护、一机一号约束。
- `niu-college-directory`: 院系字典数据与运营后台录入/维护入口，供玩家身份资料选择与院系榜聚合复用。

### Modified Capabilities
<!-- 无：本 change 不修改任何宿主基线能力的需求级行为。插件作用域调整属于插件自身 manifest 配置，不改动宿主 spec。 -->

## Impact

- **后端**：`apps/lina-plugins/sicau-niu/backend/`
  - 新增 `api/`（登录、资料、院系）、`internal/controller/`、`internal/service/`、`internal/dao`+`model/{do,entity}`。
  - 新增 `manifest/sql/` 建表：玩家用户表、院系字典表；`manifest/sql/uninstall/` 对应清理。
  - 调整 `plugin.yaml`（作用域 + 院系录入菜单）、`plugin.go`（公开登录路由 + 受保护资料/院系后台路由）、`plugin_embed.go`（嵌入 `manifest`）。
  - 移除示例 `niu` 牛只 api/controller/service 与前端示例页。
- **运营后台前端**：新增院系字典录入页面（LinaPro 插件前端页面）。
- **复用宿主能力**：登录会话/鉴权/权限中间件、字典治理与日志组件（具体接缝在 design.md 确认）。
- **数据库**：新增插件自有表 `plugin_sicau_niu_user`、`plugin_sicau_niu_college`（命名以 design/database 规则为准）。
- **i18n**：单语言（中文），插件不启用 i18n —— 无 i18n 影响。
