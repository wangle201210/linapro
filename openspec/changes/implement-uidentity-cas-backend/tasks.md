## 规则读取与影响记录

- [x] 1.1 读取并记录命中的规范：`documentation.md`、`openspec.md`、`architecture.md`、`data-permission.md`、`plugin.md`、`api-contract.md`、`backend-go.md`、`database.md`、`testing.md`、`i18n.md`和`goframe-v2`技能
- [x] 1.2 维护任务影响记录：i18n（插件未启用 i18n 资源治理，后续运行时文案仍使用结构化错误码）、数据权限（所有租户数据表保留`tenant_id`并通过宿主`TenantFilter`接入）、缓存一致性（当前阶段无缓存状态）、插件（本地未发现插件根`AGENTS.md`）、SQL（PostgreSQL-only 幂等安装/卸载 SQL）、API、后端 Go、测试和开发工具跨平台影响

## 插件骨架

- [x] 2.1 新建`apps/lina-plugins/linapro-uidentity-cas/`源码插件目录、`plugin.yaml`、`plugin_embed.go`和通用资源目录
- [x] 2.2 新建插件`backend/plugin.go`并注册受宿主认证、租户、权限中间件治理的 RESTful 路由
- [x] 2.3 更新插件工作区`go.mod`引用和宿主插件打包/嵌入发现路径

## 数据模型与 SQL

- [x] 3.1 新增插件安装、卸载 SQL，覆盖账号、账号详情、单位、容器、分组、应用、授权、黑名单、密码策略、短信、CAS 日志、OAuth 日志、OAuth token 和变更日志表
- [x] 3.2 生成或维护插件 DAO/DO/Entity，并确保生产代码不手改生成文件
- [x] 3.3 执行 SQL 幂等性、PostgreSQL-only、自增主键和索引静态检查

## 基础管理 API 与 Service

- [x] 4.1 实现账号、账号详情、单位、容器、分组、应用、密码策略和短信记录的列表、详情、创建、更新、删除接口
- [x] 4.2 实现账号应用授权、账号应用黑名单、分组应用黑名单、OAuth token、OAuth 日志、CAS 登录日志和账号变更日志接口
- [x] 4.3 实现列表关联投影批量装配、批量删除可见性校验和结构化`bizerr`错误

## 认证运行时与统计

- [x] 5.1 实现密码强度校验、管理员改密、自助改密、手机号校验改密和账号锁定/解锁
- [x] 5.2 实现 CAS ticket 校验、CAS XML 解析、应用访问规则校验和 CAS 登录日志记录
- [x] 5.3 实现 OAuth token 治理和统一身份统计聚合，确保统计不产生`N+1`查询

## 验证与提交

- [x] 6.1 运行`openspec validate implement-uidentity-cas-backend --strict`
- [x] 6.2 运行插件和受影响宿主包 Go 编译/单元测试门禁
- [x] 6.3 运行`git diff --check`和插件子仓库静态检查
- [x] 6.4 在合适阶段提交主仓库和插件仓库代码

## Feedback

- [x] **FB-1**: 源项目`uidentity/admin`服务端功能未全量迁入插件，当前实现只覆盖管理 CRUD、基础 CAS ticket 校验和统计
- [x] **FB-2**: 插件缺少旧 admin 的 CAS 服务端票据链、Token 链路、账号激活和用户自助后端接口
- [x] **FB-3**: 插件缺少旧 admin 的 OAuth 授权码流程、账号导入检查/导入、LDAP 同步边界、短信发送、文件上传、作业和监控类后端能力
- [x] **FB-4**: 二次核对发现旧 admin 的`account/unlockPassword`密码错误次数解锁后端能力未在插件内落地

## Feedback Impact Notes

- FB-4 根因记录：旧实现的`account/unlockPassword`不属于页面能力，它清理 CAS 密码登录中按`cas:pwd:errnum:<number>`累计的短期密码错误次数；插件已实现账号状态、管理员改密和运行时密码校验，但未保存密码错误次数状态，也没有按账号列表解锁的后端入口。本轮补齐时仅修改`linapro-uidentity-cas`插件，使用插件自有 token 表承载短期错误次数状态，不修改`apps/lina-core`核心框架。
- FB-4 已完成：新增受宿主认证、租户和权限中间件治理的`POST /uidentity/accounts/password-unlocks`，按当前租户可见账号列表批量清理插件 token 表中的`cas:pwd:errnum:<number>`短期错误次数状态；CAS 密码登录、runtime token、OAuth 授权码密码校验和 UnionID 账号密码绑定统一接入错误次数检查，错误达到旧项目同等阈值后返回结构化`UIDENTITY_PASSWORD_FAILURES_LOCKED`。i18n：插件未启用`i18n.enabled`，新增 API 文档源文本和错误码暂不补插件翻译资源；缓存一致性：新增短期状态以插件 OAuth token 表为权威源并依赖`expired_at`失效，不引入进程内缓存；数据权限：管理员解锁接口进入宿主认证、租户和权限中间件，解锁前按账号表租户可见性批量过滤；SQL/DAO：复用既有 token 表和唯一 code 索引，无新增 SQL 或 DAO；API/后端 Go：新增 RESTful DTO、controller、service 接口和实现，运行`make ctrl p=linapro-uidentity-cas`同步 GoFrame 控制器接口生成；开发工具跨平台：仅使用既有 Makefile 入口，无脚本变更；测试：新增纯函数单测覆盖旧 key 前缀、空账号过滤、去重和数量截断。验证补充运行`GOWORK=off go test ./... -count=1`（插件独立模块）、`GOWORK=off go test ./... -count=1`（`apps/lina-plugins`聚合模块）、`openspec validate implement-uidentity-cas-backend --strict`、主仓库和插件子仓库`git diff --check`。
- FB-1 已完成：本轮以源项目`/Users/wanna/mine/ecoding/uidentity/admin/app/**/router/*.go`实际路由为入口，重新审计旧后端能力与插件当前注册路由、资源定义、SQL 表和 service 实现。统一身份域后端能力已由`linapro-uidentity-cas`插件承载；默认系统管理、代码生成、页面/重定向壳和全局监控暴露明确不属于 CAS 插件后端迁移目标。
- FB-2 已完成：CAS TGT/ST 服务端票据链、运行时 token、账号激活、UnionID 绑定和用户自助后端接口均落在`linapro-uidentity-cas`插件内，未修改`apps/lina-core`核心框架；i18n 未启用插件资源治理，缓存无新增状态，数据权限继续通过宿主`TenantFilter`接入。
- 验证记录：已运行`GOWORK=off go test ./... -count=1`（插件独立模块）、`GOWORK=off go test ./... -count=1`（`apps/lina-plugins`聚合模块）、`openspec validate implement-uidentity-cas-backend --strict`、主仓库和插件子仓库`git diff --check`。
- 继续迁移记录：本轮补齐 FB-3 中账号导入检查/导入和短信发送子能力，导入按`number`在插件账号表内幂等新增/更新，短信发送使用插件 SMS 表和本地频控记录验证码；仍未完成 OAuth 授权码完整流程、LDAP 同步边界、文件上传、作业和监控类后端能力。验证补充运行`GOWORK=off go test ./... -count=1`（插件独立模块），新增纯函数单测覆盖导入空行、生日日期归一化和短信类型校验。
- 继续迁移记录：本轮补齐 FB-3 中 OAuth 授权码运行时子能力，新增公开 RESTful 接口`POST /uidentity/oauth/authorization-codes`、`POST /uidentity/oauth/access-tokens`和`GET /uidentity/oauth/access-tokens/{accessToken}/user-info`，复用插件账号、应用、黑名单、OAuth token 和 OAuth log 表；授权码按唯一 code 单点查询并在事务内一次性消费，access token 按唯一 access 单点查询，未修改`apps/lina-core`核心框架。i18n：插件未启用`i18n.enabled`，新增 API 文档源文本不要求补插件 apidoc 翻译资源；缓存一致性：无新增缓存或快照状态；数据权限：公开运行时接口仍通过插件`TenantFilter`上下文限制租户表数据，授权码和 access token 均绑定账号、应用和租户；SQL/DAO：复用既有 OAuth 表和唯一索引，无新增 SQL/DAO；DI：无新增运行期依赖，继续复用启动期注入的`BizCtx`、插件`Config`和`TenantFilter`服务；开发工具跨平台：未修改脚本或工具入口；测试：新增`uidentity_oauth_runtime_test.go`覆盖 URL-escaped client secret、redirect URI 校验和授权码 grant 类型。仍未完成 LDAP 同步边界、文件上传、作业和监控类后端能力。
- 继续迁移记录：本轮补齐 FB-3 中文件上传、健康检查、服务器监控、日志快照和旧 LDAP/外部文件/作业/外部监控动作边界。新增插件本地 RESTful 接口`POST /uidentity/legacy/uploads`、`GET /uidentity/legacy/health`、`GET /uidentity/legacy/server-monitor`、`GET /uidentity/legacy/log-snapshots`和`POST /uidentity/legacy/external-actions`；上传默认写入插件配置的本地目录并返回旧`uploadFile`兼容字段，监控使用本机进程/OS 投影，日志读取改为带上限的 tail snapshot，外部 LDAP/作业执行在未配置执行器时返回结构化 unsupported-flow 错误，避免把外部依赖硬编码进插件或污染`apps/lina-core`。i18n：插件未启用`i18n.enabled`，新增 API 文档源文本不要求补插件 apidoc 翻译资源；缓存一致性：无新增缓存或快照状态；数据权限：健康检查公开且不返回业务数据，上传/监控/日志/外部动作接口进入宿主认证、租户和权限中间件，文件和日志仅访问插件配置边界内的本地资源；SQL/DAO：无新增 SQL/DAO，复用插件现有表；DI：新增`gopsutil/v4`作为插件本地监控依赖，无新增接口型运行期依赖，继续复用启动期注入的`BizCtx`、插件`Config`和`TenantFilter`服务；开发工具跨平台：未修改脚本或工具入口；测试：新增`uidentity_legacy_ops_test.go`覆盖 base64 payload 拆分、安全文件名和日志 tail 上限逻辑。验证补充运行`GOWORK=off go test ./... -count=1`（插件独立模块）。
- 继续迁移记录：本轮补齐旧 admin 静态配置发现接口和 CAS XML `serviceValidate`返回形态。新增插件本地 RESTful 接口`GET /uidentity/legacy/config/cas`、`GET /uidentity/legacy/config/ldap`、`GET /uidentity/legacy/config/oauth`、`GET /uidentity/legacy/config/token`和`POST /uidentity/legacy/cas/service-validations.xml`；配置发现只读取插件作用域配置并提供旧结构字段，LDAP 默认只暴露 unsupported 元数据而不启动执行器；CAS XML 接口复用既有 ST 消费、账号访问和委托授权校验，成功时输出`cas:authenticationSuccess`，失败时按旧协议输出`INVALID_TICKET` XML，而不是 JSON 错误。根因判断：前几轮已迁入运行时 ticket 和外部动作边界，但旧 admin 的`/api/v1/config/*`与`/sso/serviceValidate`XML 协议仍缺少插件侧替代入口。i18n：插件未启用`i18n.enabled`，新增 API 文档源文本不要求补插件 apidoc 翻译资源；缓存一致性：无新增缓存或快照状态，配置读取复用宿主启动期注入的插件配置服务；数据权限：配置发现接口进入宿主认证、租户和权限中间件，CAS XML 为运行时公开协议且通过一次性 ST 自隔离，不返回未绑定 ticket 的业务数据；SQL/DAO：无新增 SQL/DAO，复用既有 token、账号、应用和授权表；DI：无新增接口型运行期依赖，继续复用启动期注入的`BizCtx`、插件`Config`和`TenantFilter`服务；开发工具跨平台：未修改脚本或工具入口；测试：新增`uidentity_legacy_cas_xml_test.go`和`uidentity_legacy_config_test.go`覆盖 XML 成功/失败渲染、配置默认值和插件作用域覆盖读取。验证补充运行`GOWORK=off go test ./... -count=1`（插件独立模块）、`GOWORK=off go test ./... -count=1`（`apps/lina-plugins`聚合模块）和插件子仓库`git diff --check`。
- 继续迁移记录：本轮补齐旧 admin 作业定义、作业日志和本地作业运行状态动作。新增插件表`plugin_linapro_uidentity_cas_sys_job`、`plugin_linapro_uidentity_cas_job_log`及 DAO/DO/Entity，通用资源新增`sys-jobs`、`job-logs`和旧别名`sysjob`、`job-log`；`POST /uidentity/legacy/external-actions`支持`job_start`和`job_remove`，对启用作业设置或清理插件内`entry_id`，不把外部 cron 执行器硬编码进`apps/lina-core`。i18n：插件未启用`i18n.enabled`，新增 API 文档源文本不要求补插件 apidoc 翻译资源；缓存一致性：无新增缓存或快照状态；数据权限：作业和日志表均保留`tenant_id`并通过宿主`TenantFilter`接入，列表/详情/删除继续走可见性校验；SQL/DAO：安装/卸载 SQL 为 PostgreSQL-only 幂等建表/删表，索引覆盖租户、状态、分组、调度 entry、job_id 和执行时间查询；DI：无新增运行期依赖；开发工具跨平台：仅更新 DAO 配置和 Go 生成代码，未修改脚本入口；测试：新增作业目标解析和运行 entry 生成/保留纯函数单测。验证补充运行`GOWORK=off go test ./... -count=1`（插件独立模块）、插件子仓库`git diff --check`。
- 继续迁移记录：本轮补齐旧 admin 微信扫码登录后端状态链。新增插件本地 RESTful 接口`POST /uidentity/cas/wechat-login-qrs`、`POST /uidentity/cas/wechat-login-callbacks`和`GET /uidentity/cas/wechat-login-qrs/{state}/result`；扫码 state、回调结果和轮询结果复用插件 OAuth token 存储，已绑定`unionId`时签发 CAS TGT/ST，未绑定时创建既有 UnionID 绑定 challenge，未配置真实微信 OAuth 解析器时记录`unsupported`状态而不是把第三方依赖硬编码进插件或`apps/lina-core`。i18n：插件未启用`i18n.enabled`，新增 API 文档源文本不要求补插件 apidoc 翻译资源；缓存一致性：无新增缓存或快照状态，扫码状态以数据库 token 表为权威源并随过期时间失效；数据权限：运行时公开接口只通过一次性 state/ticket/challenge 自隔离，账号、应用、绑定和登录日志读写继续经`TenantFilter`接入；SQL/DAO：无新增 SQL/DAO，复用 OAuth token、账号详情、账号、应用和日志表；DI：无新增运行期依赖，继续复用启动期注入的`BizCtx`、插件`Config`和`TenantFilter`服务；开发工具跨平台：本轮运行`make ctrl p=linapro-uidentity-cas`同步 GoFrame 控制器接口生成，修复生成器追加的默认构造以保留插件显式依赖注入；测试：新增`uidentity_wechat_runtime_test.go`覆盖授权 URL 装配和扫码结果投影。验证补充运行`GOWORK=off go test ./... -count=1`（插件独立模块）、`GOWORK=off go test ./... -count=1`（`apps/lina-plugins`聚合模块）、`openspec validate implement-uidentity-cas-backend --strict`和插件子仓库`git diff --check`。FB-3 具名外围能力已完成，FB-1 继续保留全量兼容性总审计。
- 最终总审计记录：源项目统一身份域管理资源`account`、`account_details`、`units`、`containers`、`groups`、`applications`、`account_unit`、`account_app_role`、`account_app_blacklist`、`group_app_blacklist`、`pass_ruler`、`sms`、`cas_login_log`、`oauth_log`、`oauthtoken`、`account_change_log`、`job_log`和`stat`已映射到插件通用资源`/uidentity/{resource}`、`/uidentity/stats`、`sys-jobs`和`job-logs`；账号改密、导入检查/导入、短信发送、CAS TGT/ST、CAS XML `serviceValidate`、runtime token、OAuth authorization code/access token/user-info、账号激活、UnionID 查询/绑定、用户自助、微信扫码登录、上传、健康检查、服务器监控、日志快照、静态配置发现和外部动作边界均存在插件 RESTful 入口。源项目`sys_user`、`sys_role`、`sys_menu`、`sys_dept`、`sys_dict`、`sys_config`、`sys_post`、`sys_api`、登录/刷新/登出、Swagger、WebSocket、静态文件和验证码属于 LinaPro 宿主默认系统管理、页面/开发工具或通用工作台能力，不复制进 CAS 插件；`other/gen_router.go`的代码生成、数据库表/列浏览和菜单/API 生成属于开发工具链，不是统一身份运行时后端；`oauth/login`、`oauth/auth`、`cas/login`、`sso/login`、`bindUnionIDCallBack`等 HTML 或重定向壳由插件已有 RESTful 运行时接口替代，不迁移页面；旧`/metrics`是全局 Prometheus 暴露，LinaPro 已有独立`linapro-monitor-server`插件承载通用监控，CAS 插件仅保留旧`server-monitor`投影。规则影响：本轮只更新 OpenSpec 任务记录，无代码、SQL、API、i18n 资源、缓存、数据权限、DI 或开发工具跨平台新增影响；验证使用源路由`rg`审计、插件路由/资源/SQL 静态检索、`openspec validate implement-uidentity-cas-backend --strict`和当前工作区 Go 测试门禁。
