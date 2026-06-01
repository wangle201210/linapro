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

- [ ] 4.1 实现账号、账号详情、单位、容器、分组、应用、密码策略和短信记录的列表、详情、创建、更新、删除接口
- [ ] 4.2 实现账号应用授权、账号应用黑名单、分组应用黑名单、OAuth token、OAuth 日志、CAS 登录日志和账号变更日志接口
- [ ] 4.3 实现列表关联投影批量装配、批量删除可见性校验和结构化`bizerr`错误

## 认证运行时与统计

- [ ] 5.1 实现密码强度校验、管理员改密、自助改密、手机号校验改密和账号锁定/解锁
- [ ] 5.2 实现 CAS ticket 校验、CAS XML 解析、应用访问规则校验和 CAS 登录日志记录
- [ ] 5.3 实现 OAuth token 治理和统一身份统计聚合，确保统计不产生`N+1`查询

## 验证与提交

- [x] 6.1 运行`openspec validate implement-uidentity-cas-backend --strict`
- [ ] 6.2 运行插件和受影响宿主包 Go 编译/单元测试门禁
- [x] 6.3 运行`git diff --check`和插件子仓库静态检查
- [ ] 6.4 在合适阶段提交主仓库和插件仓库代码
