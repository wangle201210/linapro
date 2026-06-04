## 1. 插件作用域与脚手架调整

- [x] 1.1 修改 `plugin.yaml`：`scope_nature=platform_only`、`supports_multi_tenant=false`、`default_install_mode=global`；移除 i18n 块（保持单语言）
- [x] 1.2 更新 `plugin.yaml` 菜单：新增院系字典录入菜单页 + 增删改按钮权限（`sicau-niu:college:*`），以及玩家查询菜单（`sicau-niu:player:list`）
- [x] 1.3 移除示例 `niu` 牛只代码：`backend/api/niu`、`backend/internal/controller/niu`、`backend/internal/service/niu`、`frontend/pages/sidebar-entry.vue`、`frontend/pages/niu-client.ts`
- [x] 1.4 调整 `plugin_embed.go` 嵌入 `plugin.yaml frontend manifest`

## 2. 数据库与 DAO

- [x] 2.1 新增 `manifest/sql/` 建表：玩家表 `plugin_sicau_niu_user`（openid 唯一、phone 唯一、身份资料列、device_fingerprint、软删除与时间列）
- [x] 2.2 新增 `manifest/sql/` 建表：院系表 `plugin_sicau_niu_college`（name 唯一、sort、软删除与时间列）
- [x] 2.3 新增 `manifest/sql/uninstall/` 对应删表脚本，保证幂等
- [x] 2.4 配置 `backend/hack/config.yaml` 并执行 `make dao p=sicau-niu` 生成 DAO/DO/Entity

## 3. 配置与微信网关接缝（可 Mock）

- [x] 3.1 新增 `manifest/config` 配置项：微信 `appid/secret`、玩家 JWT 密钥与有效期、微信网关 Mock 开关
- [x] 3.2 在 `backend/internal/service` 定义 `WeChatAuthGateway` 接缝（`Code2Session`、`DecodePhone`）与 **Mock 实现**，真实实现留待后续接入
- [x] 3.3 实现玩家 token 签发/验签（插件自签 JWT，密钥/有效期取自配置）

## 4. API 契约与代码生成

- [x] 4.1 定义玩家面 `api/player/v1` DTO：`/player/login`、`/player/phone`、`GET/PUT /player/profile`、`GET /colleges`（RESTful、时间字段 `int64` Unix 毫秒、完整 `dc`/`eg` 标签）
- [x] 4.2 定义运营面 `api/admin/v1` DTO：院系 CRUD（`/admin/colleges`）、玩家查询（`/admin/players`），受保护接口在 `g.Meta` 声明 `permission` 标签
- [x] 4.3 执行 `make ctrl p=sicau-niu` 生成 controller 骨架与 api 接口文件

## 5. 服务层实现 — 玩家身份（`niu-player-identity`）

- [x] 5.1 实现微信登录：经 `WeChatAuthGateway` 解析 openid，首登建档，签发 token，返回 `isNewUser`；无效 code 返回 `bizerr` 鉴权错误
- [x] 5.2 实现手机号授权绑定：解析手机号、唯一性校验（一机一号）、记录设备指纹；冲突返回 `bizerr`
- [x] 5.3 实现身份资料读取/更新：身份类型 Go 命名常量约束、院系存在性校验、年级/毕业年数字校验
- [x] 5.4 实现玩家鉴权中间件：验签玩家 token、载入玩家、注入玩家上下文，保证只能读写本人数据
- [x] 5.5 实现运营玩家分页查询（只读，先在 DB 侧分页过滤）

## 6. 服务层实现 — 院系字典（`niu-college-directory`）

- [x] 6.1 实现院系 CRUD：名称唯一校验、分页查询（DB 侧分页）、排序维护
- [x] 6.2 实现院系删除引用保护：被玩家引用则拒绝删除并返回 `bizerr`
- [x] 6.3 实现玩家院系下拉：按排序一次性返回有界列表

## 7. 鉴权与路由装配

- [x] 7.1 在 `backend/plugin.go` 装配路由：公开 `/player/login`；玩家面用插件玩家鉴权中间件；运营面用宿主 `Auth+Tenancy+Permission`
- [x] 7.2 controller `NewV1()` 显式持有 service 字段；service 构造函数显式注入依赖（网关、token 器、DAO），不在请求路径临时 `New()`

## 8. 运营后台前端（院系录入）

- [x] 8.1 新增院系字典管理页面（`useVbenVxeGrid` + `Page`，列表/新增/编辑/删除，按钮按 `sicau-niu:college:*` 权限显隐）+ 玩家查询只读页面
- [x] 8.2 新增院系/玩家 API 客户端，时间字段按 Unix 毫秒展示格式化

## 9. 测试与验证门禁

- [x] 9.1 服务层单测：登录建档/复登、手机号一机一号冲突、资料校验、院系唯一与删除引用保护（纯逻辑测试始终运行；DB 集成测试以 `LINA_TEST_PGSQL_LINK` 门控，自包含含清理）
- [x] 9.2 运行后端编译门禁：构建 plugins workspace 覆盖插件后端与路由绑定（`go build`/`go vet`/`go test` 均 rc=0）
- [x] 9.3 用户可观察行为（院系录入页）E2E 已编写并**执行通过**：`hack/tests/e2e/TC001-college-dictionary-crud.ts`（接入 `ensureSourcePluginEnabled` 自动装启用插件）+ `hack/tests/pages/SicauNiuCollegePage.ts`。运行 `pnpm test:module -- plugin:sicau-niu`（浏览器 history 模式 `E2E_BASE_URL=:5666` + 后端 `:9120`）：**4 passed**（新增/编辑/同名拒绝/删除，均断言列表持久态）。
  - 运行前修复了一处**宿主环境表结构漂移**（与本变更无关）：`sys_online_session` 缺 `client_type` 列导致登录写在线会话失败、token 被判无效。已非破坏性补列 `ALTER TABLE sys_online_session ADD COLUMN IF NOT EXISTS client_type VARCHAR(32) NOT NULL DEFAULT ''`（彻底修法为 `make db.init`）。
- [x] 9.4 运行 `openspec validate niu-identity --strict` → valid
- [x] 9.5 见下方「影响判断」小节

## 影响判断（9.5）

- **i18n**：插件单语言（中文），`plugin.yaml` 无 i18n 块，前端无 `$t()`，后端无运行期用户可见文案走翻译 —— **无 i18n 影响**。
- **缓存一致性**：C1 玩家 token 为无状态 HMAC 签名，无新增跨节点缓存/快照/失效逻辑 —— **无缓存一致性影响**（遵循 `cache-consistency` 无需改动）。
- **数据权限**：玩家面经插件玩家中间件注入 playerID，service 一律以该 ID 限定本人数据（BindPhone/GetProfile/UpdateProfile 均按 `playerID` 读写）；院系/玩家查询为运营受权限接口（`sicau-niu:*`），DB 侧分页过滤；铁牛/远程写入不涉及。公开接口仅登录。
- **DI 来源检查**：新增运行期依赖（`wechat.Gateway`、`token.Service`、`identity.Service`、`college.Service`、`PlayerAuth` 中间件）均在 `plugin.go:registerRoutes` 一次性构造并经构造函数显式注入；owner=插件后端装配入口；配置来自 `services.Config()`（宿主插件作用域静态配置）；无临时 `New()` 关键服务、无聚合 Deps 结构体；player token 无状态故无需共享后端实例。
- **跨平台/开发工具**：仅新增插件内 SQL/Go/Vue/配置与 `backend/hack/config.yaml`，未改 `Makefile`/`linactl`/CI/脚本 —— **无跨平台影响**。
- **数据库**：单迭代单 SQL 文件，建表/索引幂等（`IF NOT EXISTS` + 部分唯一索引），软删除走 `deleted_at` 自动维护，未显式写自增 id —— 符合 `database` 规则。
