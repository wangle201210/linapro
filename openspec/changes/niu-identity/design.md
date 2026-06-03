## Context

C1 是寻牛活动的身份地基（见 `apps/lina-plugins/sicau-niu/design.md` 程序级总体设计的 C1）。它要解决一个核心结构问题：

- **寻牛玩家与 LinaPro 管理员是两个不同人群**。宿主 `auth`/`session`/`bizctx`/`notify` 契约都围绕**管理员/租户用户**设计（host user ID、tenant token、online session、host inbox），无法直接承载「微信 openid + 手机号」的玩家账户。
- 因此玩家身份必须**插件自有**：插件维护玩家表、自签玩家 token、自建玩家鉴权中间件；而运营后台（院系录入、用户查询）继续走宿主 `Auth + Tenancy + Permission`。

当前 `sicau-niu` 是最小示例（`niu` 牛只占位），C1 用真实身份能力取代它，并把插件作用域从 `tenant_aware` 改为 `platform_only`（活动不涉及多租户）。

约束：遵循 `.agents/rules/` 中 `plugin`、`backend-go`、`api-contract`、`database`、`architecture`、`data-permission` 规则；单语言（中文），不启用 i18n。

## Goals / Non-Goals

**Goals:**
- 玩家可通过微信小程序登录（`code` 换 openid）并完成手机号一键授权绑定，后端签发玩家会话 token。
- 玩家可维护身份资料：身份标签（在校生/校友/川农好友）、院系（从字典选择）、年级（填数字）、毕业年（校友填数字）、昵称。
- 运营可在后台维护院系字典并查询玩家基础信息。
- 一机一号：手机号唯一，记录设备指纹。
- 插件作用域调整为 `platform_only` / `supports_multi_tenant=false` / `global`。
- 提供面向小程序的后端 API（**不实现小程序 UI**）。

**Non-Goals:**
- 地图/激活/卡片/喂草/偷草/送草/排行榜/荣誉/海报/H5/看板（后续 change）。
- 真实微信 appid/secret 联调（C1 用可 Mock 的接缝，凭证后续配置）。
- 完整异常风控告警视图（C7），C1 只落「一机一号 + 设备指纹」基础数据。

## Decisions

### D1 玩家身份插件自有，与宿主管理员身份分离
玩家不是宿主用户。插件维护 `plugin_sicau_niu_user`（openid/phone/昵称/身份资料/设备指纹），不复用宿主用户表与宿主 `AuthService/SessionService`。
- **理由**：宿主契约绑定 host user ID/tenant token；玩家是外部微信人群，规模与生命周期独立。
- **备选**：把玩家塞进宿主用户体系 —— 否决，会污染宿主管理员/租户语义。

### D2 双鉴权面
- **玩家面**（小程序 API，`/player/*`）：插件自有 **player-auth 中间件** 校验玩家 token，注入玩家上下文。
- **运营面**（后台 API，`/admin/*`）：沿用宿主 `Auth + Tenancy + Permission`（与现有 `platform_only` 插件一致），权限标识 `sicau-niu:*`。
- **理由**：两类调用者、两套凭证，必须分离；运营面复用宿主治理链。

### D3 微信登录走可 Mock 的网关接缝
定义插件内 `WeChatAuthGateway` 接缝：`Code2Session(code) → openid`、`DecodePhone(...) → phone`。C1 提供 **Mock 实现**（开发/测试返回可配置值），真实实现调用微信 `jscode2session` 等，凭证 `appid/secret` 走插件 `manifest/config`。
- **理由**：与铁牛定位同样的 Mock-seam 策略，C1 不被微信联调阻塞，接口边界稳定后续直接替换。
- **备选**：直连微信 —— C1 阶段无凭证、不可测，否决。

### D4 玩家 token 为插件自签 JWT
登录成功签发插件 JWT（密钥/有效期来自插件配置），`Authorization: Bearer` 携带；player-auth 中间件验签 + 载入玩家。C1 只发 access token（有效期内重登换新），刷新机制列入开放问题。
- **理由**：无状态、跨节点一致，不依赖宿主会话存储；`platform_only` 单部署足够。

### D5 一机一号
`phone` 唯一约束 = 一手机号一账户；`device_fingerprint` 记录但 C1 不做强拦截，仅留数据给 C7 风控。

### D6 身份类型用命名常量，院系用自有表
- 身份类型（在校生/校友/川农好友）：插件内 Go 命名类型 + 常量（稳定小枚举、插件私有、单语言），满足 `backend-go` 枚举规则；不入宿主字典。
- 院系：可增长业务数据，建 `plugin_sicau_niu_college` 表 + 后台 CRUD 录入，供身份选择与后续院系榜聚合关联。
- 年级/毕业年：用户填写的整数，不自动从学号带出。

### D7 插件作用域改 platform_only
`plugin.yaml`：`scope_nature=platform_only`、`supports_multi_tenant=false`、`default_install_mode=global`；移除 i18n 块（保持单语言）。

### D8 复用示例脚手架，移除占位业务
保留 `sicau-niu` 目录与构建/注册结构，移除 `niu` 牛只 api/controller/service/前端示例，替换为玩家身份与院系能力。

### D9 显式依赖注入
- service 构造函数逐项显式接收依赖（如 `WeChatAuthGateway`、token 签发/验签器、玩家 DAO、院系 DAO）；纯值配置（JWT 密钥、有效期、Mock 开关）用配置结构体。
- controller 在 `NewV1()` 持有 service 字段；player-auth 中间件显式接收 token 验签器与玩家查询能力。
- 不在请求路径临时 `New()` 关键服务。

### 数据模型（C1）
| 表 | 关键列 | 约束 |
|----|-------|------|
| `plugin_sicau_niu_user` | `openid`、`phone`、`nickname`、`avatar`、`identity_type`、`college_id`、`grade`、`graduation_year`、`device_fingerprint`、时间列 | `openid` 唯一、`phone` 唯一、`college_id` 可空软外键 |
| `plugin_sicau_niu_college` | `name`、`sort`、时间列 | `name` 唯一 |

> 软删除与时间字段遵循 `database` 规则；DAO/DO/Entity 由 `make dao` 生成。

### API 契约（C1，RESTful，时间字段 Unix 毫秒 `int64`）
玩家面（`/{api}/api/v1/plugins/sicau-niu/...`）：
| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| `POST` | `/player/login` | 公开 | 微信 `code` 登录，返回 token + 是否新用户 |
| `POST` | `/player/phone` | player | 手机号授权绑定（一机一号） |
| `GET` | `/player/profile` | player | 当前玩家资料 |
| `PUT` | `/player/profile` | player | 更新昵称/身份标签/院系/年级/毕业年 |
| `GET` | `/colleges` | player | 院系下拉（供资料选择） |

运营面（后台，`Auth+Permission`）：
| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| `GET` | `/admin/colleges` | `sicau-niu:college:list` | 院系列表（分页） |
| `POST` | `/admin/colleges` | `sicau-niu:college:create` | 新增院系 |
| `PUT` | `/admin/colleges/{id}` | `sicau-niu:college:update` | 修改院系 |
| `DELETE` | `/admin/colleges/{id}` | `sicau-niu:college:delete` | 删除院系 |
| `GET` | `/admin/players` | `sicau-niu:player:list` | 玩家基础信息查询（分页只读） |

## Risks / Trade-offs

- [微信外部依赖在 C1 无法联调] → `WeChatAuthGateway` Mock 接缝，开发/测试可跑通全流程，真实凭证后续配置替换。
- [作用域由 `tenant_aware` 改 `platform_only`] → 示例插件无真实数据，风险低；`platform_only/global` 模式已有 `linapro-monitor-server` 等先例。
- [玩家自签 token 安全] → JWT 密钥入插件配置、设置有效期；C1 不引入刷新令牌，避免过度设计。
- [手机号唯一 vs 一人多微信] → 一手机号一账户，符合「一机一号」；同人换微信需用同手机号登录绑定。
- [院系删除影响已选用户] → 删除做软删除并校验是否被引用（或置空引用），在 specs/tasks 明确。

## Migration Plan

1. 改 `plugin.yaml`：作用域三项 + 院系录入菜单与权限按钮；移除 i18n 块。
2. 新增 `manifest/sql/` 建表（user、college）与 `manifest/sql/uninstall/`；`plugin_embed.go` 嵌入 `manifest`。
3. 新增 `manifest/config`：微信 `appid/secret`、JWT 密钥/有效期、Mock 开关。
4. 移除示例 `niu` 代码，新增玩家/院系 `api/controller/service`、player-auth 中间件，更新 `plugin.go` 路由（公开登录 + player 面 + 运营面）。
5. `make dao` 生成 DAO/DO/Entity；`make ctrl` 生成 controller 骨架。
6. 编译门禁（plugins workspace 构建）+ 单测；用户可观察行为按需补 E2E。
7. 回滚：作用域与代码改动可整体回退；插件未上线、无生产数据。

## Open Questions

1. 微信 `appid/secret` 与手机号授权方式（`getPhoneNumber` 新版 `code` vs 旧版加密串）的最终形态（Mock 不阻塞 C1）。
2. 玩家 token 是否需要刷新令牌（C1 暂不做，按运营时长评估）。
3. 院系删除策略：禁止删除被引用院系，还是删除后置空用户引用（specs 决策）。

> 影响判断：i18n —— 单语言不启用，**无影响**；缓存一致性 —— C1 玩家 token 无状态、无新增跨节点缓存，**无影响**；数据权限 —— 玩家只能读写自身资料，院系/玩家查询为运营受权限接口；开发工具跨平台 —— 仅新增插件内 SQL/代码，无脚本变更。
