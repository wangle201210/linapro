## Why

`/Users/wanna/mine/ecoding/uidentity/admin`已经沉淀了统一身份、CAS/OAuth 接入、账号治理、授权关系和审计日志等后端能力。LinaPro 需要在不绑定具体页面的前提下，把这些能力作为可持续交付的 AI 原生全栈框架中的源码插件能力落地，避免把身份域业务细节污染到`lina-core`核心宿主契约。

## What Changes

- 新增`linapro-uidentity-cas`源码插件，承载统一身份域后端能力，插件仓库和主仓库均使用`feat/cas`分支。
- 新增插件自有 PostgreSQL schema、安装/卸载 SQL、插件清单、后端 API、Controller、Service 和必要测试。
- 实现账号、账号详情、单位、容器、分组、应用、账号应用授权、账号应用黑名单、分组应用黑名单、密码策略、短信记录、CAS 登录日志、OAuth 日志、OAuth token 与账号变更日志的后端管理能力。
- 实现 CAS ticket 校验、CAS 登录结果记录、OAuth token 管理、密码强度校验、账号密码自助变更和统一身份统计聚合能力。
- 复用宿主认证、权限、租户、插件生命周期、日志、`bizerr`和 GoFrame 数据访问范式，不复制默认工作台页面。

## Capabilities

### New Capabilities

- `uidentity-cas-backend`: 统一身份 CAS/OAuth 源码插件后端能力，包括插件 schema、API、服务、认证运行时和治理验证。

### Modified Capabilities

- 无。

## Impact

- 影响目录：`apps/lina-plugins/linapro-uidentity-cas/`、`apps/lina-plugins/go.mod`、必要的插件打包/嵌入入口。
- API：新增插件命名空间下的 RESTful HTTP API，由宿主插件路由注册并受统一认证、租户和权限中间件治理。
- 数据库：新增插件自有表，表名使用`plugin_linapro_uidentity_cas_*`前缀；SQL 放在插件`manifest/sql/`并满足 PostgreSQL 幂等约束。
- 数据权限：读取类、详情类、写入类和批量删除接口均按插件租户边界与宿主数据权限能力接入或记录明确例外。
- i18n：插件启用`i18n`时，新增运行时菜单和 API 文档源文本需维护插件自身资源；若后续确认关闭插件`i18n`，需在任务记录中说明。
- 缓存：统计聚合可使用短 TTL 缓存或直接聚合；如新增缓存，必须记录失效边界。
- 测试：至少运行 OpenSpec 严格校验、插件 Go 编译门禁、受影响包单元测试和 SQL 静态检查；涉及运行时认证链路时补充服务层测试。
