# UIdentity Admin 服务端能力审计

## 审计基线

本轮审计针对旧项目`/Users/wanna/mine/ecoding/uidentity/admin/app`，只判断服务端能力，不迁移页面、静态资源和旧 GoAdmin 工作台展示结构。插件根目录`apps/lina-plugins/linapro-uidentity-cas/`未发现本地`AGENTS.md`，本轮继续遵守仓库顶层规则与命中规则文件。

审计命令覆盖旧项目 Go 文件清单、旧 router 注册、旧 service/model 副作用、旧 jobs registry、插件 REST 路由、插件通用资源、插件 cron 执行器和插件测试文件。结论按“实现 / 测试或验证 / 排除理由”逐项记录。

## 逐项能力映射

| 旧服务端能力 | 旧项目证据 | 插件映射 | 测试或验证 | 结论 |
| --- | --- | --- | --- | --- |
| 账号、账号详情、单位、容器、分组、应用、密码策略、短信记录 | `app/admin/router/account*.go`、`units.go`、`containers.go`、`groups.go`、`applications.go`、`pass_ruler.go`、`sms.go` | `GET/POST/PUT/DELETE /uidentity/{resource}`，资源包括`accounts`、`account-details`、`units`、`containers`、`groups`、`applications`、`pass-rules`、`sms-records` | `uidentity_resource*`、`uidentity_account_import_test.go`、插件`go test` | 已实现 |
| 授权、黑名单、日志、OAuth token | `account_unit.go`、`account_app_role.go`、`account_app_blacklist.go`、`group_app_blacklist.go`、`cas_login_log.go`、`oauth_log.go`、`oauthtoken.go`、`account_change_log.go`、`job_log.go` | 通用资源`account-units`、`account-app-roles`、`account-app-blacklists`、`group-app-blacklists`、`cas-login-logs`、`oauth-logs`、`oauth-tokens`、`account-change-logs`；旧`/api/v1/job-log`保留兼容地址，真实任务日志由宿主任务管理承载 | 资源注册静态审计、`plugin_legacy_routes_test.go`、插件`go test` | 已实现 |
| 账号列表筛选、详情关联投影 | `Account.GetPage`预加载`Unit/Detail/Groups/Container`，按`groupIds/passLevels/containerId/unitId/status`筛选 | `ListResource`数据库侧分页过滤，账号记录批量装配单位、容器、分组投影，`groupIds/passLevels`过滤 | `uidentity_resource.go`静态审计，插件`go test` | 已实现 |
| 账号创建和更新的`groupIds`维护`account_group` | `dto.AccountInsertReq.Generate`写入`model.Groups`，`Account.Update`显式增删`account_group` | 新增`uidentity_account_groups.go`，账号创建/更新同步替换`account_group`，创建前校验组租户可见性 | `TestAccountResourceSyncsGroupIDs` | 本轮发现缺口并已修补 |
| 账号详情自动创建、账号/详情变更审计 | `models/account.go`、`models/account_details.go`的 GORM hook 和`OperateErr` | `uidentity_account_audit.go`显式写`account_change_log`，账号创建自动创建详情并审计，写路径脱敏`passwordHash` | `TestResourceAccountAuditLifecycle` | 已实现 |
| 管理员改密、密码自助、密码错误次数解锁 | `account/updatePassword`、`updatePasswordGetUser`、`updatePasswordBySelfPhone`、`updatePasswordBySelf`、`unlockPassword` | `/uidentity/accounts/{id}/password`、`/uidentity/password-challenges*`、`/uidentity/accounts/password-unlocks` | `uidentity_password_test.go`、插件`go test` | 已实现 |
| 账号导入检查与导入 | `app/admin/service/account_import.go` | `/uidentity/accounts/import-checks`、`/uidentity/accounts/imports`，按`number`幂等新增/更新账号和详情 | `uidentity_account_import_test.go` | 已实现 |
| 统计聚合 | `app/admin/router/stat.go`、`service/stat.go` | `/uidentity/stats`数据库侧聚合账号、应用、CAS/OAuth、密码等级、登录类型和应用维度 | 插件路由/服务静态审计，插件`go test` | 已实现 |
| CAS 登录、TGT/ST、ST 校验、登出、XML serviceValidate | `app/cas/router/api.go`、`service/cas.go`、`models/ticket.go` | `/uidentity/cas/password-logins`、`phone-logins`、`union-id-logins`、`service-tickets`、`service-validations`、`tickets/{ticket}`、`legacy/cas/service-validations.xml` | `uidentity_legacy_cas_xml_test.go`、插件`go test` | 已实现 |
| runtime token 签发与按 token 取用户信息 | `app/cas/router/api.go`的`/token/get`和`getUserInfoByToken` | `/uidentity/runtime-tokens`、`/uidentity/runtime-tokens/{accessToken}/user-info` | `uidentity_legacy_runtime.go`静态审计，插件`go test` | 已实现 |
| 激活 baseInfo、face、password、phone、wechat 和 state | `app/cas/service/active.go` | `/uidentity/activations`、`face/password/phone/wechat`、`wechat-states`、`wechat-callbacks`、`state` | `uidentity_activation_user_test.go` | 已实现 |
| UnionID 查询、绑定、迁移旧绑定和激活日志 | `app/cas/service/user.go`的`GetByUnionID`、`BindUnionID`、`account_active_log` | `/uidentity/users/union-id-lookups`、`union-id-bindings`，插件`account-active-logs`资源 | `TestRebindUnionIDToAccountMigratesExistingBinding` | 已实现 |
| 用户自助改密、改手机、改邮箱、改 QQ、解绑/换绑微信、个人信息、登录日志、应用、委托角色 | `app/cas/router/api.go`的`/user/*` | `/uidentity/users/{number}/password|phone|email|qq|wechat`、`wechat-rebind-*`、`applications`、`account-app-roles` | `uidentity_activation_user.go`静态审计、插件`go test` | 已实现 |
| 微信扫码登录 state、回调、轮询 | `app/cas/router/wechat.go`和`cas/router/api.go`扫码接口 | `/uidentity/cas/wechat-login-qrs`、`wechat-login-callbacks`、`wechat-login-qrs/{state}/result` | `uidentity_wechat_runtime_test.go` | 已实现 |
| OAuth 授权码、换 token、用户信息和授权日志 | `app/oauth/router/api.go`、`oauth/apis/oauth.go` | `/uidentity/oauth/authorization-codes`、`access-tokens`、`access-tokens/{accessToken}/user-info`，换 token 写`oauth-log` | `uidentity_oauth_runtime_test.go` | 已实现 |
| 静态配置发现 | `app/admin/router/config.go` | `/uidentity/legacy/config/cas|ldap|oauth|token`，只读插件作用域配置 | `uidentity_legacy_config_test.go` | 已实现 |
| 文件上传、健康、服务器监控、日志快照 | `app/other/router/file.go`、`monitor.go`、`sys_server_monitor.go`、`log.go` | `/uidentity/legacy/uploads`、`health`、`server-monitor`、`log-snapshots` | `uidentity_legacy_ops_test.go` | 已实现 |
| 作业定义、启动、移除和执行日志旧地址 | `app/jobs/router/sys_job.go`、`app/admin/router/job_log.go`、`jobs/service/sys_job.go`、`jobbase.go` | 插件通过`plugin.Cron().RegisterCron`向宿主 gf/gcron 任务管理注册九个旧作业 handler；旧`/api/v1/sysjob`投影九个 managed job，旧`/api/v1/job/start/{id}`、`/api/v1/job/remove/{id}`、`/api/v1/job-log`保留兼容地址；不再维护插件自建`sys_job/job_log`表或本地 cron | `jobs_register_test.go`、`plugin_legacy_routes_test.go`、插件`go test` | 已实现 |
| `WannaT`作业 | `app/jobs/wanna.go` | `jobs`服务注册同名 handler，保留轻量日志探针 | `jobs_register_test.go`和静态审计 | 已实现 |
| `ContainerAccount/NewContainerAccount`作业 | `app/jobs/container_account.go` | `jobs`服务集合化刷新容器账号数 | `jobs`服务测试 | 已实现 |
| `ChangeContainer`作业本地 DB 迁移 | `app/jobs/change_container.go` | `jobs`服务把当年毕业账号集合化迁入`xy`容器；LDAP DN 变更属于外部目录边界 | `TestUpdateGraduatingAccountContainer` | 已实现本地副作用 |
| `SyncMysql2Ldap`作业 | `app/jobs/sync_mysql2ldap.go`调用 LDAP | 插件`jobs`服务支持配置化 LDAP executor；配置缺失时返回`UIDENTITY_JOB_EXECUTOR_UNSUPPORTED`并由宿主任务管理记录结果 | `jobs_register_test.go`和静态审计 | 已实现配置化边界 |
| `SyncStudent`、`SyncStudentYJS`、`SyncStudentWJ`、`SyncDept`、`SyncJzg` | `app/jobs/sync_oracle2mysql.go`依赖外部 Oracle schema | 插件`jobs`服务支持配置化 Oracle 同步；配置缺失时返回`UIDENTITY_JOB_EXECUTOR_UNSUPPORTED`并由宿主任务管理记录结果 | `jobs_register_test.go`和静态审计 | 已实现配置化边界 |
| 账号创建/更新、密码同步、容器变更的 LDAP 写入 | `AddOrModifyLdapUser`、`SyncPassword`、`ChangeContainerID` | 插件默认只维护插件自有账号数据；真实 LDAP 写入归入`legacy/external-actions`和 cron unsupported 外部执行器边界 | 设计`Non-Goals`和静态审计 | 明确排除默认实现 |
| GoAdmin 默认系统管理 | `sys_user`、`sys_role`、`sys_menu`、`sys_dept`、`sys_dict`、`sys_config`、`sys_post`、`sys_api`、登录、刷新、登出、WebSocket、Swagger | 由 LinaPro 宿主默认系统管理、认证、权限、配置、字典和工作台能力承载，不复制进 CAS 插件 | 设计`Non-Goals` | 不属于插件目标 |
| 代码生成、数据库表列浏览、验证码、静态页面和 HTML/重定向壳 | `app/other/router/gen_router.go`、`captcha`、`oauth/login`、`oauth/auth`、`cas/login`、`sso/login`、`bindUnionIDCallBack`、`cas-login/index` | 开发工具和页面不迁移；旧重定向壳地址已在兼容层注册并按旧 redirect/空 200/XML/JSON 响应语义处理。`cas-login/index`旧实现最终调用 GoAdmin JWT 中间件写宿主登录态，plugin 侧仅保留 CAS ticket 校验兼容入口，不伪造 LinaPro core auth session | `plugin_legacy_routes_test.go`、`uidentity_legacy_http_test.go`、路由审计 | 已实现兼容边界 |
| 全局 Prometheus `/metrics` | `app/other/router/monitor.go` | 通用监控由 LinaPro 监控插件或宿主治理；CAS 插件仅提供旧`server-monitor`投影 | 路由审计 | 不属于 CAS 插件目标 |

## 第一轮反向审计：旧项目到插件

审计命令：

- `find /Users/wanna/mine/ecoding/uidentity/admin/app -maxdepth 3 -type f -name '*.go'`
- `rg -n '\\.(GET|POST|PUT|DELETE)\\(|v1(auth)?\\.(GET|POST|PUT|DELETE)\\(' .../app/*/router`
- `rg -n 'SyncMysql2Ldap|SyncStudent|SyncStudentYJS|SyncStudentWJ|SyncDept|SyncJzg|ChangeContainer|NewContainerAccount|WannaT' .../app/jobs`

发现与处理：

- 统一身份域管理资源、CAS、OAuth、Token、激活、用户自助、短信、导入、上传、配置、监控快照、作业和旧兼容路由均已在插件内有服务端入口或资源映射。
- 本轮发现账号创建/更新请求中的`groupIds`没有落到插件`account_group`关系表；已新增`uidentity_account_groups.go`并补`TestAccountResourceSyncsGroupIDs`。
- 默认 GoAdmin 系统管理、代码生成、Swagger、WebSocket、验证码、静态页面和全局 Prometheus 不迁入 CAS 插件，理由见逐项表。
- 外部 Oracle/LDAP 执行器由插件`jobs`服务按配置启用；配置缺失时返回结构化 unsupported 错误并由宿主任务管理记录执行结果，避免把外部依赖硬编码进`lina-core`。

## 第二轮反向审计：插件到旧项目与副作用

审计命令：

- `rg -n 'path:"|group\\.(GET|POST|PUT|DELETE)|RegisterRoutes|/uidentity/' apps/lina-plugins/linapro-uidentity-cas/backend`
- `rg -n 'return map\\[string\\]\\*resourceDefinition|account-groups|plugin.Cron|RegisterCron|job-log|legacyRoot' apps/lina-plugins/linapro-uidentity-cas/backend`
- `rg -n 'After|Before|OperateErr|AddOrModifyLdapUser|ChangeContainerID|GroupIDs|SyncPassword' /Users/wanna/mine/ecoding/uidentity/admin/app`
- `rg -n 'unsupported|ExecutorUnsupported|SyncStudent|SyncMysql2Ldap|ChangeContainer|NewContainerAccount' apps/lina-plugins/linapro-uidentity-cas/backend/internal/service/jobs`

发现与处理：

- 插件公开路由和旧兼容路由均能反向对应旧统一身份后端能力、运行时能力、旧兼容配置/运维能力或通用资源管理能力；根`/sso`、根`/ssologin`、Wechat 域名校验、`stat/get`公开性和`job-log`旧地址已在 FB-15 中补齐。
- 旧模型 hook 副作用已由插件显式审计 helper 承接；账号详情、账号状态、密码、UnionID、微信换绑、激活等直接写表路径均不依赖旧 GORM hook。
- 旧`GroupIDs/groupIds`关系副作用是本轮唯一新增缺口，已经补齐创建、更新、清空和跨租户组拒绝逻辑。
- 旧 LDAP/Oracle 直接写入和同步进入插件`jobs`服务的配置化执行器；未配置时有稳定结构化失败边界，未发现需要修改`apps/lina-core`的能力。

## 影响记录

- 插件：当前实现和兼容复核只修改`apps/lina-plugins/linapro-uidentity-cas/`插件后端与 OpenSpec 记录，不修改`apps/lina-core`核心框架。
- API：旧`/api/v1`、根`/sso`、根`/ssologin`、Wechat 根回调和域名校验旧地址在插件兼容层注册；账号创建/更新已有通用资源接口接受并处理`groupIds`字段。
- 数据权限：`groupIds`写入先按当前租户校验组可见性，关系删除也经`TenantFilter`过滤当前租户。
- SQL/DAO：复用既有`account_group`和`group`表、索引和 DAO，无新增 SQL 或生成文件。
- 缓存一致性：无新增缓存。
- i18n：无新增用户可见运行时文案资源；插件仍未启用`i18n.enabled`。
- 测试：新增并维护`TestAccountResourceSyncsGroupIDs`、`jobs_register_test.go`、`plugin_legacy_routes_test.go`和旧字段/重定向 helper 单测；验证继续运行插件全量测试、聚合模块测试、OpenSpec 严格校验和`git diff --check`。

## 当前结论

完成多轮反向审计后，旧 admin 统一身份服务端能力、旧接口路由和主要旧响应/重定向语义已由`linapro-uidentity-cas`插件承载。页面、默认系统管理、开发工具和全局监控有明确排除理由；`cas-login/index`涉及旧 GoAdmin JWT 登录态写入，当前仅保留 CAS ticket 校验兼容入口，不在 plugin 内伪造 LinaPro core auth session。
