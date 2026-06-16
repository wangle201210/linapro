## 1. 实现

- [x] 1.1 新增按策略 ID 查询租户策略绑定的 API DTO 和控制器方法。
- [x] 1.2 在`media`服务层实现只读`media_strategy_tenant`的分页查询。
- [x] 1.3 补充后端路由测试，验证接口只返回租户策略绑定。
- [x] 1.4 同步媒体管理端列表接口`pageSize`校验上限为`10000`，并覆盖服务层上限一致性验证。

## 2. 验证

- [x] 2.1 运行覆盖`media`插件后端的 Go 测试。
- [x] 2.2 运行`openspec validate add-media-strategy-tenant-lookup --strict`。

## Feedback

- [x] **FB-1**: 新增按策略 ID 查询设备策略绑定列表接口，并验证只读取设备策略表。
- [x] **FB-2**: `UserDeviceStrategyByToken`按节点 ID 检查租户在该节点未关闭会话数量，达到租户流配置限制时返回稳定错误码。
- [x] **FB-3**: `water`插件改为通过配置调用远端`mediaopen`内部策略解析接口，解除对`media`源码插件的本地运行时依赖。

### FB-2 反馈修复记录

- 根因：`UserDeviceStrategyByToken`原请求和服务输入只有`token`和`deviceId`，无法识别租户在具体媒体节点上的流数量限制；策略解析路径也没有读取`media_tenant_stream_config`和`media_report_session`中按租户、节点、未关闭会话的活跃数量，因此达到节点限流后仍会返回策略内容。
- 修复：兼容接口新增必填`nodeId`；服务层在铁塔设备权限通过且存在匹配策略时，将`nodeId`解析为租户流配置的`node_num`，仅当该租户节点存在启用且`max_concurrent > 0`的配置时，统计`media_report_session`中同租户、同节点且`close_time IS NULL`的会话数量，达到上限时返回`MEDIA_TENANT_STREAM_LIMIT_EXCEEDED`；未配置或未启用限流的租户节点保持不限数量。
- 插件本地规范：`apps/lina-plugins/media/AGENTS.md`不存在，按仓库顶层`AGENTS.md`和命中规则执行。
- 架构边界：修改限定在`media`源码插件公开兼容接口和插件内服务层，不修改`apps/lina-core`核心宿主契约，不新增跨模块内部实现依赖。
- 数据权限：该接口仍先通过铁塔 token 得到租户身份并校验租户设备权限；限流统计只使用该 token 租户对应的`tenant_id`和请求节点，不接受调用方覆盖租户，不泄露其他租户会话存在性。
- 缓存一致性：限流判断直接读取数据库投影和租户流配置，不新增缓存、快照或失效逻辑。
- i18n：`media`插件`plugin.yaml`未配置`i18n.enabled: true`，按单语言插件处理；新增业务错误码通过`bizerr.MustDefine`提供稳定`errorCode`和 fallback 文案，不新增插件`manifest/i18n`或`apidoc`翻译资源。
- 数据库：不新增本变更自己的 SQL 文件；复用`add-media-data-collection-server`补充的`media_report_session.close_time`和`idx_media_report_session_active_tenant_node`索引，租户流配置读取复用既有`media_tenant_stream_config`主键。
- 开发工具跨平台：不修改`Makefile`、脚本、CI或`linactl`入口。
- DI 来源检查：未新增运行期接口依赖；限流逻辑复用`media.Service`已有数据库 DAO 和已有铁塔 token 解析链路。
- 性能：租户节点限流为一次配置投影查询加一次按索引计数查询，只有存在启用限流配置时才访问报表会话表；不会随策略列表、会话明细或返回行数产生`N+1`查询。
- 测试策略：本次为功能行为反馈，已新增服务层单元测试覆盖关闭会话不计入限流、达到未关闭会话上限时返回稳定错误码；未涉及前端页面或 E2E 资产，未触发 E2E 质量审查。
- 规则加载：已读取`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/plugin.md`、`.agents/rules/backend-go.md`、`.agents/rules/api-contract.md`、`.agents/rules/database.md`、`.agents/rules/data-permission.md`、`.agents/rules/cache-consistency.md`、`.agents/rules/documentation.md`、`.agents/rules/testing.md`、`.agents/rules/i18n.md`、`.agents/rules/architecture.md`、`.agents/rules/dev-tooling.md`，并使用`lina-feedback`、`goframe-v2`和`karpathy-guidelines`技能。
- 验证：`go test ./backend/internal/service/media -count=1`通过；`go test ./backend/api/media/v1 ./backend/api/mediaopen/v1 -count=1`通过；`go test ./backend/internal/controller/media ./backend/internal/controller/mediaopen -count=1`通过；`go test ./backend -count=1`通过；`openspec validate add-media-strategy-tenant-lookup --strict`通过。

### FB-3 反馈修复记录

- 根因：`water`插件在`plugin.yaml`中声明依赖`media`，并在启动路由时通过`lina-plugin-media/backend/provider/strategy`构造本地`media`策略 resolver；该路径要求`water`所在集群同时安装并运行`media`源码插件，不适合`media`单集群部署、`water`多集群部署的拓扑。
- 修复：`mediaopen`新增受`X-Inner-Api-Key`保护的`GET /api/v1/strategies/resolve`只读接口，复用`media.Service.ResolveStrategy`按租户设备、设备、租户、全局优先级解析策略；`water`移除`media`插件依赖声明和 Go import，改为从插件运行时配置`mediaStrategy.baseUrl`、`mediaStrategy.apiKey`、`mediaStrategy.timeout`构造 HTTP resolver。
- 插件本地规范：`apps/lina-plugins/media/AGENTS.md`和`apps/lina-plugins/water/AGENTS.md`均不存在，按仓库顶层`AGENTS.md`和命中规则执行。
- 架构边界：本次新增跨集群 HTTP DTO 作为稳定契约，未让`water`访问`media`内部 DAO/DO/Entity 或内部 service；`media`仍是策略权威源，`water`只消费远端投影。
- 数据权限：新增接口位于`mediaopen`内部 API Key 鉴权链下，读取租户和设备策略投影，不接收调用方覆盖权限上下文，不新增管理端列表、导出、批量或租户可见性入口。
- 缓存一致性：不新增缓存、快照或失效机制；`water`每次处理按配置调用远端策略接口，异步任务状态缓存仍复用宿主 cache。
- i18n：`media`与`water`插件均未配置`i18n.enabled: true`，按单语言插件处理；新增和修改的 API 文档源文本、插件清单说明与 README 不新增插件多语言资源。
- 数据库：不新增或修改 SQL、DAO、DO、Entity；`media.Service.ResolveStrategy`继续使用现有策略表和既有查询路径。
- 开发工具跨平台：仅使用既有`make -C apps/lina-plugins ctrl p=media`生成接口绑定，不修改`Makefile`、脚本、CI或`linactl`入口。
- DI 来源检查：`water`新增运行期依赖为插件配置服务读取出的纯值配置和 HTTP resolver；owner 为`water`插件，创建位置为`backend/plugin.go`路由注册阶段，传递路径为`watersvc.New(cacheSvc, strategyResolver)`，不新增共享缓存敏感服务实例。
- 性能：单次水印处理最多一次远端策略解析 HTTP GET，不在列表、批量或聚合路径循环调用；`media`端解析沿用固定优先级查询链，不随返回行数产生`N+1`。
- 测试策略：本次为跨模块运行时行为反馈，新增`water`远端 resolver 单元测试覆盖请求头、查询参数、统一响应解码和未配置错误；新增`media`路由测试覆盖内部策略解析接口必须经过 Inner API Key 鉴权。未涉及前端页面或 E2E 资产，未触发 E2E 质量审查。
- 规则加载：已读取`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/architecture.md`、`.agents/rules/data-permission.md`、`.agents/rules/plugin.md`、`.agents/rules/api-contract.md`、`.agents/rules/backend-go.md`、`.agents/rules/cache-consistency.md`、`.agents/rules/documentation.md`、`.agents/rules/testing.md`、`.agents/rules/i18n.md`，并使用`lina-feedback`和`goframe-v2`技能。
- 验证：`go test lina-plugin-water/backend/internal/service/water -count=1`通过；`go test lina-plugin-water/backend -count=1`通过；`go test lina-plugin-water/... -count=1`通过；`go test ./backend/api/mediaopen/v1 ./backend/internal/controller/mediaopen ./backend -run 'TestMediaOpenResolveStrategyRequiresInnerAPIAuth|TestMediaOpenRoutesUseInnerAPIAuth' -count=1`通过；`go test ./backend -count=1`通过；`go test ./... -count=1`在`apps/lina-plugins/media`通过；`openspec validate add-media-strategy-tenant-lookup --strict`通过。
- 宿主启动绑定补充验证：`go test ./internal/cmd -count=1`未计入通过验证；失败点为当前工作区 panic allowlist 计数不匹配，`apps/lina-plugins/media/backend/plugin.go:init`期望 2 实际 3，`apps/lina-plugins/sicau-niu/backend/plugin.go:init`期望 2 实际 3。两个`plugin.go`文件本次均未修改。
