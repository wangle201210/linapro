# Tasks

- [x] 1. 更新租户流配置 OpenSpec 增量规范，明确`tenantId + nodeNum`复合唯一键和删改定位语义。
- [x] 2. 修正`media_tenant_stream_config`表约束、服务查询、详情、更新、删除、控制器与前端 API 调用，使同一租户可配置多个节点且删改只影响指定节点，并提升`media`插件版本以触发已安装插件执行新增迁移。
- [x] 3. 补充或更新自动化测试，覆盖同一租户不同节点的创建、详情、更新、删除和列表保留行为。
- [x] 4. 运行 OpenSpec、后端 Go、前端类型或相关 E2E 验证，并记录影响分析。

## 影响分析

- 架构边界：本次修改限定在`media`源码插件的租户流配置能力内，未修改`apps/lina-core`宿主核心契约，未新增模块启停或跨模块依赖。
- API契约：租户流配置详情、更新、删除路径从单`tenantId`定位调整为`tenantId + nodeNum`复合定位；创建、更新响应返回`tenantId`和`nodeNum`，便于前端按复合键刷新。
- 数据库：新增`002-media-tenant-stream-composite-key.sql`，幂等删除旧租户单字段约束并重建`("tenant_id", "node_num")`主键，同时保留租户查询索引；旧数据在原单租户唯一约束下不会产生复合键冲突。
- 插件升级：`media/plugin.yaml`版本从`v0.1.0`提升到`v0.1.1`，用于触发源码插件同步后的升级流程执行新增 SQL；根因是已安装同版本插件不会自动运行后续新增迁移文件。
- 数据权限：权限标签未扩展，仍沿用`media:management:*`管理权限；本次只修正资源定位键，不新增跨租户数据可见入口。
- 性能：列表继续使用分页查询，节点名称仍按节点编号集合批量装配，未引入随返回行数增长的`N+1`查询；删改查按主键复合条件定位。
- 缓存一致性：租户流配置当前未接入缓存或快照，本次无缓存失效、刷新和集群一致性影响。
- 开发工具跨平台：未修改`Makefile`、脚本、CI、代码生成或`linactl`入口，无跨平台执行影响。
- i18n：`media`插件未启用插件语言包资源，本次未新增运行时语言包；API文档源文本和现有中文界面文案随复合键语义局部更新。
- 测试策略：补充后端服务级回归测试覆盖同租户多节点创建、重复校验、更新隔离和删除隔离；更新媒体插件烟测用例中的 API/UI 路径、行键和清理逻辑。
- DI来源：未新增运行期依赖、服务构造函数参数、启动装配、插件宿主服务适配器或`WASM host service`依赖，DI无影响。

## 反馈处理记录

- [x] FB-1：已安装`media`插件同版本不会执行新增复合键 SQL。处理方式为提升插件版本到`v0.1.1`，并通过重启服务后的插件升级流程验证复合键迁移生效。
- [x] FB-2：媒体插件剩余 E2E 页面用例存在入口和清理等待不稳定问题。处理方式为将插件页面导航统一走工作台路径、修正文档页断言为当前 OpenAPI 标签，并在 UI 清理段等待活动表格稳定。

## 验证记录

- `make dev`，已重启后端`http://127.0.0.1:9120/`和前端`http://127.0.0.1:5666/`
- `openspec validate fix-media-tenant-stream-composite-key --strict`
- `go test ./backend/internal/service/media -run TestTenantStreamConfigUsesTenantAndNodeKey -count=1`
- `go test ./backend/internal/service/media -count=1`
- `go test ./... -count=1`
- `pnpm -C apps/lina-vben -F @lina/web-antd run typecheck`
- `pnpm -C hack/tests test:validate`
- `E2E_HOST_PORT=5666 pnpm -C hack/tests test:module -- plugin:media -- --grep "TC-1c"`，通过复合键 REST 语义验证。
- `E2E_HOST_PORT=5666 pnpm -C hack/tests test:module -- plugin:media -- --grep "TC-1f"`，通过租户流配置 UI 编辑、删除和新增验证。
- `E2E_HOST_PORT=5666 pnpm -C hack/tests test:module -- plugin:media`，6 个媒体插件 E2E 全部通过。
- `git diff --check`
- `go test ./internal/cmd -count=1`未计入通过验证；失败点为当前工作区`apps/lina-plugins/sicau-niu/backend/plugin.go`的 panic allowlist 数量变化，与本次`media`复合键变更无关。
- `go test ./internal/cmd/internal/httpstartup -count=1`未计入通过验证；失败点为现有动态插件路由测试命中`runtime-config/global`缓存新鲜度错误，与本次`media`复合键变更无关。
- `go test ./pkg/dialect -count=1`未计入通过验证；失败点覆盖宿主和其他插件历史 SQL 治理项，以及当前工作区`sicau-niu` SQL 改动，与本次新增`media/manifest/sql/002-media-tenant-stream-composite-key.sql`无直接关系。
