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

- [ ] **FB-1**: 源项目`uidentity/admin`服务端功能未全量迁入插件，当前实现只覆盖管理 CRUD、基础 CAS ticket 校验和统计
- [x] **FB-2**: 插件缺少旧 admin 的 CAS 服务端票据链、Token 链路、账号激活和用户自助后端接口
- [ ] **FB-3**: 插件缺少旧 admin 的 OAuth 授权码流程、账号导入检查/导入、LDAP 同步边界、短信发送、文件上传、作业和监控类后端能力

## Feedback Impact Notes

- FB-1/FB-3 仍在推进：旧 admin 的 OAuth 授权码完整流程、账号导入检查/导入、短信发送、文件上传、健康检查、服务器监控、日志快照和 LDAP/外部文件/作业/外部监控动作边界已迁入插件；仍需继续做旧 admin 服务端接口兼容性总审计，尤其是实际 LDAP 同步执行器、作业 CRUD/调度执行器是否应由插件承载或委托宿主稳定能力。
- FB-2 已完成：CAS TGT/ST 服务端票据链、运行时 token、账号激活、UnionID 绑定和用户自助后端接口均落在`linapro-uidentity-cas`插件内，未修改`apps/lina-core`核心框架；i18n 未启用插件资源治理，缓存无新增状态，数据权限继续通过宿主`TenantFilter`接入。
- 验证记录：已运行`GOWORK=off go test ./... -count=1`（插件独立模块）、`GOWORK=off go test ./... -count=1`（`apps/lina-plugins`聚合模块）、`openspec validate implement-uidentity-cas-backend --strict`、主仓库和插件子仓库`git diff --check`。
- 继续迁移记录：本轮补齐 FB-3 中账号导入检查/导入和短信发送子能力，导入按`number`在插件账号表内幂等新增/更新，短信发送使用插件 SMS 表和本地频控记录验证码；仍未完成 OAuth 授权码完整流程、LDAP 同步边界、文件上传、作业和监控类后端能力。验证补充运行`GOWORK=off go test ./... -count=1`（插件独立模块），新增纯函数单测覆盖导入空行、生日日期归一化和短信类型校验。
- 继续迁移记录：本轮补齐 FB-3 中 OAuth 授权码运行时子能力，新增公开 RESTful 接口`POST /uidentity/oauth/authorization-codes`、`POST /uidentity/oauth/access-tokens`和`GET /uidentity/oauth/access-tokens/{accessToken}/user-info`，复用插件账号、应用、黑名单、OAuth token 和 OAuth log 表；授权码按唯一 code 单点查询并在事务内一次性消费，access token 按唯一 access 单点查询，未修改`apps/lina-core`核心框架。i18n：插件未启用`i18n.enabled`，新增 API 文档源文本不要求补插件 apidoc 翻译资源；缓存一致性：无新增缓存或快照状态；数据权限：公开运行时接口仍通过插件`TenantFilter`上下文限制租户表数据，授权码和 access token 均绑定账号、应用和租户；SQL/DAO：复用既有 OAuth 表和唯一索引，无新增 SQL/DAO；DI：无新增运行期依赖，继续复用启动期注入的`BizCtx`、插件`Config`和`TenantFilter`服务；开发工具跨平台：未修改脚本或工具入口；测试：新增`uidentity_oauth_runtime_test.go`覆盖 URL-escaped client secret、redirect URI 校验和授权码 grant 类型。仍未完成 LDAP 同步边界、文件上传、作业和监控类后端能力。
- 继续迁移记录：本轮补齐 FB-3 中文件上传、健康检查、服务器监控、日志快照和旧 LDAP/外部文件/作业/外部监控动作边界。新增插件本地 RESTful 接口`POST /uidentity/legacy/uploads`、`GET /uidentity/legacy/health`、`GET /uidentity/legacy/server-monitor`、`GET /uidentity/legacy/log-snapshots`和`POST /uidentity/legacy/external-actions`；上传默认写入插件配置的本地目录并返回旧`uploadFile`兼容字段，监控使用本机进程/OS 投影，日志读取改为带上限的 tail snapshot，外部 LDAP/作业执行在未配置执行器时返回结构化 unsupported-flow 错误，避免把外部依赖硬编码进插件或污染`apps/lina-core`。i18n：插件未启用`i18n.enabled`，新增 API 文档源文本不要求补插件 apidoc 翻译资源；缓存一致性：无新增缓存或快照状态；数据权限：健康检查公开且不返回业务数据，上传/监控/日志/外部动作接口进入宿主认证、租户和权限中间件，文件和日志仅访问插件配置边界内的本地资源；SQL/DAO：无新增 SQL/DAO，复用插件现有表；DI：新增`gopsutil/v4`作为插件本地监控依赖，无新增接口型运行期依赖，继续复用启动期注入的`BizCtx`、插件`Config`和`TenantFilter`服务；开发工具跨平台：未修改脚本或工具入口；测试：新增`uidentity_legacy_ops_test.go`覆盖 base64 payload 拆分、安全文件名和日志 tail 上限逻辑。验证补充运行`GOWORK=off go test ./... -count=1`（插件独立模块）。
