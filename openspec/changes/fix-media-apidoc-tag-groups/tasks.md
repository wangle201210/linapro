## Feedback

- [x] **FB-1**: 媒体插件接口文档中管理端接口全部聚合在`媒体管理`分组下，左侧导航过于混乱
- [x] **FB-2**: mediaopen 内部对接接口中`媒体鉴权`分组命名不准确，需要改为`内部接口`

### 根因记录

- `FB-1`根因：媒体插件管理端所有请求 DTO 的`g.Meta`都声明为`tags:"媒体管理"`，Stoplight Elements 按 OpenAPI `tags`渲染一级导航时只能生成单个大分组，无法进一步展示资源级分类。
- `FB-2`根因：mediaopen 的 token 策略查询和租户白名单 IP 查询 DTO 仍使用`tags:"媒体鉴权"`，该分组实际属于内部对接接口，不是管理端鉴权模块，导致接口文档左侧分组语义不准确。

### 影响分析

- 已读取规则：`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/documentation.md`、`.agents/rules/architecture.md`、`.agents/rules/data-permission.md`、`.agents/rules/plugin.md`、`.agents/rules/api-contract.md`、`.agents/rules/backend-go.md`、`.agents/rules/cache-consistency.md`、`.agents/rules/testing.md`、`.agents/rules/i18n.md`、`.agents/rules/dev-tooling.md`。
- 插件本地规范：`apps/lina-plugins/media/AGENTS.md`不存在，按仓库顶层规范和规则文件执行。
- 宿主边界影响：本次只调整媒体源码插件自有 API 文档元数据，不修改`apps/lina-core`核心宿主服务、通用接口契约或工作台页面实现。
- API 契约影响：仅调整 OpenAPI `tags`展示分组；不修改 HTTP 方法、资源路径、请求响应字段、权限标签或数据装配成本。
- 后端 Go 影响：修改插件 API DTO 元数据和插件后端测试；不新增运行期依赖、构造函数或服务图。
- `i18n`影响：`media`插件`plugin.yaml`未配置`i18n.enabled: true`，按单语言插件处理，不维护插件`manifest/i18n`或`apidoc`翻译资源。
- 数据权限影响：无业务数据读取、写入、导出、聚合统计或可见性边界变化。
- 缓存一致性影响：无缓存、快照、失效、刷新或集群一致性变化。
- 开发工具跨平台影响：无`Makefile`、脚本、`CI`、代码生成或`linactl`入口变更。
- 测试策略：更新插件后端 OpenAPI 文档测试，断言管理端接口按资源输出多个 tag，旧的单一`媒体管理`分组不再承载管理端路径，且 mediaopen 内部对接接口使用`内部接口`分组。

### 执行记录

- `FB-1`修复：将媒体插件管理端 42 个请求 DTO 的 OpenAPI `tags`从粗粒度`媒体管理`拆分为`媒体策略`、`节点管理`、`设备节点`、`流别名`、`策略绑定`、`租户流配置`和`租户白名单`。
- `FB-1`测试：更新`apps/lina-plugins/media/backend/media_plugin_routes_test.go`，通过插件自有`/api/v1/media/openapi.json`断言管理端接口按资源输出多个一级 tag，并避免回退到单一`媒体管理`分组。
- `FB-1`规范：新增`openspec/changes/fix-media-apidoc-tag-groups/specs/media-api-docs/spec.md`，记录媒体插件接口文档按资源分组的用户可观察要求。
- `FB-2`修复：将 mediaopen 的 token 策略查询和租户白名单 IP 查询 DTO 从`tags:"媒体鉴权"`改为`tags:"内部接口"`；保留`路由记忆`、`媒体配置`等其他 mediaopen 分组不变。
- `FB-2`测试：扩展`apps/lina-plugins/media/backend/media_plugin_routes_test.go`，断言`/api/v1/strategies/user-device`和`/api/v1/tenant-whites/ips`输出`内部接口`tag，并扫描插件自有 OpenAPI 文档不再输出旧的`媒体鉴权`tag。

### 验证记录

- 已执行但失败：`GOWORK=off go test ./backend -run 'TestMediaPluginOpenAPIDocumentOnlyContainsMediaRoutes|TestMediaPluginAPIDocsPageLoadsMediaDocument' -count=1`（在`apps/lina-plugins/media`内）。失败原因是关闭 Go workspace 后无法解析`lina-core`工作区依赖，属于命令选择问题，不是本次代码失败。
- 已执行并通过：`go test ./backend -run 'TestMediaPluginOpenAPIDocumentOnlyContainsMediaRoutes|TestMediaPluginAPIDocsPageLoadsMediaDocument' -count=1`（在`apps/lina-plugins/media`内）。
- 已执行并通过：`go test ./backend -count=1`（在`apps/lina-plugins/media`内）。
- 已执行并通过：`openspec validate fix-media-apidoc-tag-groups --strict`。
- 已执行并通过：`rg -n 'tags:"媒体管理"' apps/lina-plugins/media/backend/api/media/v1 || true`，确认管理端旧分组无残留。
- 已执行并通过：`rg -n 'tags:"(媒体策略|节点管理|设备节点|流别名|策略绑定|租户流配置|租户白名单)"' apps/lina-plugins/media/backend/api/media/v1 | wc -l`，确认 42 个管理端 DTO 已归入新分组。
- 已执行并通过（`FB-2`）：`go test ./backend -run 'TestMediaPluginOpenAPIDocumentOnlyContainsMediaRoutes|TestMediaPluginAPIDocsPageLoadsMediaDocument' -count=1`（在`apps/lina-plugins/media`内）。
- 已执行并通过（`FB-2`）：`go test ./backend -count=1`（在`apps/lina-plugins/media`内）。
- 已执行并通过（`FB-2`）：`openspec validate fix-media-apidoc-tag-groups --strict`。
- 已执行并通过（`FB-2`）：`rg -n 'tags:"媒体鉴权"|媒体鉴权|tags:"内部接口"' apps/lina-plugins/media/backend/api/mediaopen apps/lina-plugins/media/backend/media_plugin_routes_test.go`，确认生产 DTO 仅保留`内部接口`，旧 tag 只出现在测试的禁止出现断言中。
