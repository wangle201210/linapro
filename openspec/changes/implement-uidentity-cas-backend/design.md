## Context

源项目`uidentity/admin`是 GoAdmin/Gin/GORM 实现，包含默认后台系统管理能力和统一身份业务能力。用户要求旧前端后续只切换接口基址、不修改前端代码，因此本次插件需要同时承载统一身份域业务后端和旧 admin 非 job 系统管理接口兼容面；页面仍不迁移，`lina-core`核心框架不因旧前端参数而改动。

本次设计遵守核心宿主边界：`apps/lina-core`继续提供框架级认证、权限、租户、插件生命周期和通用模块接口；统一身份域以源码插件`linapro-uidentity-cas`落地，插件自有数据、API 和服务均位于`apps/lina-plugins/linapro-uidentity-cas/`。

## Goals / Non-Goals

**Goals:**

- 完整承接源项目统一身份后端功能：账号、详情、单位、容器、分组、应用、授权、黑名单、密码策略、短信记录、认证日志、OAuth token、变更日志、CAS 登录、密码自助变更和统计聚合。
- 兼容迁入旧 admin 运行时后端能力：CAS 服务端票据链、Token 链路、账号激活、用户自助、OAuth 授权码、账号导入检查/导入、LDAP 同步边界、短信发送、文件上传、作业和监控类接口。
- 保持旧`uidentity/admin`非 job HTTP 路由、请求参数和响应字段契约完全一致，使旧前端能够直接切换到 LinaPro 宿主基址。
- 在保留`plugin_linapro_uidentity_cas_`表前缀的前提下，使插件 PostgreSQL 表后缀和字段结构对齐旧 MySQL/GORM 模型；本项目不做租户级旧表隔离，按全局级数据语义承载。
- 使用插件同构后端结构：`backend/api/`、`backend/plugin.go`、`backend/internal/controller/`、`backend/internal/service/`、`backend/internal/dao/`、`backend/internal/model/`。
- 使用插件自有 PostgreSQL SQL 和 DAO/DO/Entity 生成物，不依赖宿主私有 DAO/DO/Entity。
- 高频列表、批量详情、日志和统计接口采用数据库侧过滤、分页、聚合和批量装配，避免`N+1`查询。
- RESTful API 方法语义、时间字段 Unix 毫秒响应、`g.Meta`/`dc`/`eg`/`permission`标签和`bizerr`错误治理满足项目规则。

**Non-Goals:**

- 不实现或修改前端页面、路由页面组件、表格或表单。
- 不把源项目默认系统管理页面能力写入`lina-core`。旧 admin 后端入口属于`sys_user`、`sys_role`、`sys_menu`、`sys_dict`、`sys_config`、`sys_post`、`sys_dept`等通用能力时，只在`linapro-uidentity-cas`插件内提供旧路径、旧参数、旧响应和旧表结构兼容面，不污染宿主核心模型。
- 不把真实 LDAP、外部文件存储或外部监控后端作为默认必需依赖。相关能力在插件内抽象为可配置外部同步或执行边界；未配置时返回结构化错误或使用本地插件数据完成可验证的默认流程。

## Decisions

### 新建源码插件而不是修改`lina-core`

统一身份、CAS 和 OAuth 是业务域能力，不属于框架核心宿主的默认通用契约。将能力放入`linapro-uidentity-cas`源码插件可以复用宿主认证、权限、租户和插件生命周期，同时避免把应用、分组、黑名单、授权等业务展示结构写进`lina-core`。

备选方案是直接修改`lina-core`用户认证模块。该方案会把具体身份域数据模型和 CAS/OAuth 业务接口变成核心宿主契约，降低框架复用性，因此不采用。

### 插件自有表使用稳定前缀和集合化查询

插件表统一使用`plugin_linapro_uidentity_cas_*`前缀。统一身份业务表和旧 admin 非 job 系统管理表均保留旧表后缀和字段名，PostgreSQL 类型按旧 MySQL/GORM 语义转换；列表接口在数据库侧完成过滤、排序和分页，关联名称通过批量查询和内存映射装配；统计接口通过聚合查询一次性获取各维度结果。

备选方案是完全沿用源项目无前缀原表名。该方案容易与宿主或其他插件表冲突，也不符合插件数据归属规则；用户已确认表前缀可以保留，因此不采用无前缀表名。

### 账号域保留独立账号模型

源项目统一身份账号不是 LinaPro 宿主登录用户的简单扩展，包含工号、单位、容器、账号详情、分组、密码强度、授权和外部应用登录等业务字段。插件保留独立账号模型，并通过必要字段与宿主当前用户上下文做审计关联。

备选方案是扩展宿主`sys_user`。该方案会把统一身份域字段写入核心用户模型，并增加默认用户管理接口的复杂度，因此不采用。

### CAS/OAuth 运行时采用插件内服务

CAS ticket 校验、应用状态/白名单/黑名单判断、登录日志记录和 OAuth token 存储由插件服务完成。外部 CAS validate URL 作为插件配置或运行参数读取，缺失时返回结构化配置错误。默认不在插件内创建新的宿主认证 token，避免绕过宿主认证边界。

备选方案是把 CAS 登录直接接入宿主`auth`服务签发 LinaPro JWT。本阶段需求是迁移`uidentity/admin`后端能力且不做页面，直接签发宿主 token 会改变宿主登录语义，需要单独 OpenSpec 设计，因此暂不纳入。

### 旧 admin 运行时兼容接口仍由插件承载

反馈确认当前实现范围不足，旧 admin 的 CAS 服务端票据、Token、激活、用户自助、OAuth 授权码、账号导入、LDAP、短信、文件、作业和监控类服务端能力需要继续迁入。兼容入口必须使用旧`/api/v1/*`、根`/sso/*`、根`/ssologin/*`和旧回调路径，HTTP 方法、请求参数、响应 envelope、`msg`和字段名均保持旧 GoAdmin 契约；job 路由按用户补充不作为本插件兼容目标。

账号、应用、授权、黑名单和日志仍以插件自有表为权威数据源。CAS TGT/ST、访问 token、激活 challenge 和授权码优先复用插件 OAuth token 表，通过类型化 payload 区分用途，避免新增短期状态表；读取和验证按 ticket/token 单点索引查询，不进入列表逐项查询链路。账号导入采用批量检查和批量写入，LDAP 同步默认作为外部目录边界保留结构化错误，不引入宿主核心依赖。

## Risks / Trade-offs

- [Risk] 源项目功能范围大，单次实现容易失控。→ 按可提交切片推进：插件骨架与 schema、基础 CRUD、认证运行时、统计日志、验证和审查。
- [Risk] 外部 LDAP/CAS 环境不可用导致测试不稳定。→ CAS validate 通过接口化 HTTP client 或测试替身覆盖；LDAP 写入本阶段作为可配置外部边界，不要求集成环境。
- [Risk] 日志和统计接口容易产生`N+1`查询。→ 设计中要求列表先分页，再批量装配账号/应用/容器/单位/分组投影；统计使用聚合查询和映射。
- [Risk] 插件启用`i18n`后文案资源维护量增加。→ 插件清单明确启用状态；若启用，API 源文本使用英文并维护插件自身`manifest/i18n`资源。

## Migration Plan

1. 新建`linapro-uidentity-cas`源码插件骨架、`plugin.yaml`、嵌入入口和空后端注册。
2. 新增插件 SQL schema、卸载 SQL、DAO/DO/Entity 生成物或按当前仓库插件范式维护生成结果。
3. 实现基础管理 API 和 service，优先覆盖账号、应用、分组、单位、容器、密码策略、黑名单、授权和日志列表。
4. 实现 CAS/OAuth 运行时服务、密码强度校验、自助改密和统计聚合。
5. 运行 OpenSpec、SQL、Go 编译和单元测试门禁，完成阶段提交。

回滚策略：源码插件可以通过插件生命周期卸载 SQL 删除插件自有表；若尚未启用插件，可直接移除插件目录和`go.mod`引用。

## Open Questions

- 是否需要把 CAS 登录成功接入宿主`auth`签发 LinaPro JWT？当前实现默认保持插件业务登录结果，不修改宿主认证语义。
- 是否需要真实 LDAP 写入能力？当前先保留外部同步边界，避免引入未配置环境导致默认流程不可用。
