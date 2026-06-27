## 1. 实现

- [x] 1.1 新增`media`插件采集 server 配置加载和校验，默认关闭并支持监听地址配置。
- [x] 1.2 在`media`插件启动 hook 中按配置异步启动`net-flux` TCP server。
- [x] 1.3 实现采集 server handler，支持心跳响应并接受机器、网络和流指标上报。
- [x] 1.4 补充后端单元测试覆盖配置解析、禁用分支和启动幂等分支。

## 2. 验证

- [x] 2.1 运行覆盖`media`插件后端的 Go 测试。
- [x] 2.2 运行`openspec validate add-media-data-collection-server --strict`。

## Feedback

- [x] FB-1 补齐`net-flux` discovery server 端能力，支持按配置接入 Nacos 注册、注销和查询。
- [x] FB-2 依赖来源检查：collection server 由`media`插件后端持有，在`system.started` hook 中通过宿主`payload.Services().Plugins().Config()`读取插件配置；Nacos discovery client 仅在`collectionServer.discovery.enabled=true`且收到 discovery 命令时由`collection`服务懒加载创建，并随插件上下文结束或注销动作关闭。
- [x] FB-3 使用本地 Docker Nacos 验证真实 discovery 注册、查询和注销流程，并修正发现的 Nacos group 前缀兼容问题。
- [x] FB-4 结合既有数据看板上报读模型处理`MachineMetric`、`NetworkMetric`和`StreamMetric`业务写入，替代仅记录日志的实现。
- [x] FB-5 fork`net-flux`并扩展上报协议字段，接入会话、流和实例维度字段写入完整数据看板读模型。
- [x] FB-6 按采集端只能获取容器信息修正协议语义，将`MachineMetric`收敛为实例/容器指标并移除独立`InstanceMetric`。
- [x] FB-7 按客户端只通知流和会话增减修正实时统计口径，由采集 server 通过`STREAM_ADD`、`STREAM_DELETE`、`SESSION_ADD`和`SESSION_DELETE`维护实例实时计数。
- [x] FB-8 将实例实时直播流数和会话数迁移到宿主共享 cache，避免多 Pod 部署下使用进程内存计数。
- [x] FB-9 补齐`数据看板.md`中节点总览、实例列表、流列表和会话列表对应的管理端查询接口。
- [x] FB-10 将不返回`total`的看板列表接口改为无分页查询，并保留`10000`条服务端上限。
- [x] FB-11 扩充数据看板接口 mock 数据和自动化验证，覆盖多节点、多实例、多流、多协议、多租户和`10000`条上限场景。
- [x] FB-12 核实并修正数据看板计算字段来源，确保真实采集事件写入后接口字段与`数据看板.md`保持一致。
- [x] FB-13 将采集上报写入的`memory_allocated`、`memory_used`统一为`MB`，并将`disk_io_read`、`disk_io_write`、`network_in`、`network_out`统一为`KB/S`。
- [x] FB-14 全量核查`media`模块旧单位残留，将被排除在 Git 跟踪外的`数据看板字典.xlsx`同步为内存`MB`、速率`KB/S`。
- [x] FB-15 补齐流和会话关闭时间、源流协议、客户端类型枚举、统计时长和30天清理语义。
- [x] FB-16 定时清理 cron 注册失败时必须返回插件启动错误，并允许后续启动重试。
- [x] FB-17 修复`media`插件文件版本低于数据库有效版本导致运行时状态异常，恢复源码插件自动升级入口。
- [x] FB-18 扩展`media`采集 TCP client，覆盖数据上报和 discovery 注册、查询、注销联调路径。
- [x] FB-19 实际启动`media`采集 TCP server 联调数据上报和注册发现，并修复实测发现的关闭投影和空发现响应问题。
- [x] FB-20 重新实测 TCP 上报后的落库准确性和 dashboard 统计接口准确性，并修复实测发现的时长与关闭协议统计问题。
- [x] FB-21 在`/api/v1/media/apidocs.html`中补充 TCP 采集协议说明入口，避免用户误以为 TCP 能力缺失。
- [x] FB-22 检查`media`插件内`collection-client`完整性，并修正插件工具运行入口和根`go.work`工作区范围。
- [x] FB-23 使用本地 Docker Nacos 和本地服务复测远端`lookup/register-lookup`超时与`report-close`流投影异常，修复本地可复现问题。

### FB-13 反馈修复记录

- 根因：`collection_report.go`仍按旧口径将内存写为`GB`、磁盘速率写为`MB/s`、网络速率写为`Mbps`，而看板 API DTO、SQL 注释和数据看板文档也沿用了旧单位，导致真实采集上报写入和接口文档语义不一致。
- 修复：将`MachineMetric.mem_total`和`mem_used`按字节转换为`MB`；将`disk_read_bytes`和`disk_write_bytes`按字节速率转换为`KB/S`；将协议中标注为`bps`的`network_in`、`network_out`和`NetworkMetric.throughput`先换算为字节速率后转换为`KB/S`；同步更新看板 API DTO、插件 SQL 注释、数据看板文档和生成的 DAO/DO/Entity 注释。
- 插件本地规范：`apps/lina-plugins/media/AGENTS.md`不存在，按仓库顶层`AGENTS.md`和命中规则执行。
- 架构边界：修改限定在`media`源码插件采集写入和该插件看板契约文本内，未修改`apps/lina-core`核心宿主契约，未新增模块启停、跨模块调用或抽象层。
- 数据权限：本次只修改采集数据单位归一化和文档说明，不新增读取或写入入口，不改变`platform_only`看板权限边界。
- 缓存一致性：不修改生命周期计数共享 cache、缓存键、失效或集群策略；实时流数和会话数仍沿用既有共享 cache 机制。
- i18n：`media`插件`plugin.yaml`未配置`i18n.enabled: true`，按单语言插件处理；本次仅更新中文 API 文档源文本和插件文档，不新增插件`manifest/i18n`或`apidoc`翻译资源。
- 数据库：未新增表、列、索引或 DML；仅修改当前迭代 SQL 中幂等`COMMENT ON COLUMN`文本，并在本地开发库重放`003-add-media-data-collection-server.sql`后执行`make dao p=media`刷新生成注释。
- 开发工具跨平台：不修改`Makefile`、脚本或工具入口；执行的`psql`仅用于本地验证和 DAO 生成前刷新本地开发库注释，不属于交付脚本变更。
- 性能：上报处理仍为单个报文固定次数转换和写入，不新增随返回行数、节点数或插件数增长的查询，不引入`N+1`路径。
- 规则加载：已读取`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/plugin.md`、`.agents/rules/backend-go.md`、`.agents/rules/api-contract.md`、`.agents/rules/data-permission.md`、`.agents/rules/testing.md`、`.agents/rules/documentation.md`、`.agents/rules/database.md`、`.agents/rules/i18n.md`、`.agents/rules/architecture.md`，并使用`goframe-v2`技能。
- 验证：`psql 'postgresql://postgres:postgres@127.0.0.1:5432/linapro?sslmode=disable' -v ON_ERROR_STOP=1 -f apps/lina-plugins/media/manifest/sql/003-add-media-data-collection-server.sql`通过；`make dao p=media`通过；`rg -n "单位GB|单位MB/s|单位Mbps|MB/s|Mbps|bytesPerGigabyte|bitsPerMegabit|bytesToGigabytes|bitsToMegabits|converted to 5 Mbps" apps/lina-plugins/media -S`无残留；`go test ./backend/internal/service/collection -run 'TestNormalizeMachineMetricBuildsInstanceReport|TestNormalizeNetworkMetricBuildsLatencyUpdate|TestReportRuntimePersistsMetrics' -count=1`通过；`LINAPRO_TEST_POSTGRES=1 go test ./backend/internal/service/collection -run TestReportRuntimePersistsMetrics -count=1`通过；`go test ./backend/internal/service/collection -count=1`通过；`go test ./backend/internal/service/media -count=1`通过；`go test ./backend/api/media/v1 -count=1`通过；`go test ./backend/internal/controller/media -count=1`通过；`go test ./backend -count=1`通过；`openspec validate add-media-data-collection-server --strict`通过。

### FB-14 反馈修复记录

- 根因：常规`rg`扫描无法读取被 Git exclude 的`数据看板字典.xlsx`内部 XML，导致该本地数据字典工作簿仍保留`memory_allocated`、`memory_used`的`单位GB`，以及`disk_io_read`、`disk_io_write`的`单位MB/s`和`network_in`、`network_out`的`单位Mbps`说明。
- 修复：保持工作簿结构和样式不变，仅机械替换`数据看板字典.xlsx`内联字符串中的内存单位为`MB`，磁盘和网络速率单位为`KB/S`；同时按窄口径扫描确认 Git 跟踪范围内没有新的旧单位残留。
- 排除项：`GB28181`和`GB device ID`属于国标协议/设备标识，不是存储容量单位；`bitrate`和`current_bitrate`的`Kbps`属于流码率字段，不属于本次`MB/S -> KB/S`速率口径。
- 插件本地规范：`apps/lina-plugins/media/AGENTS.md`不存在，按仓库顶层`AGENTS.md`和命中规则执行。
- 架构边界：本次仅补齐`media`插件本地数据字典资料，不修改`apps/lina-core`核心宿主契约，不新增模块接口、模块启停或跨模块依赖。
- 数据权限：不新增或修改 HTTP API、读取路径、写入入口或数据权限过滤逻辑。
- 缓存一致性：不修改共享 cache、缓存键、失效、刷新或集群一致性策略。
- i18n：`media`插件`plugin.yaml`未配置`i18n.enabled: true`，按单语言插件处理；本次只更新中文本地资料文件，不新增语言包或`apidoc`翻译资源。
- 数据库：不修改 SQL、表、列、索引、DML 或 DAO 生成结果；`数据看板字典.xlsx`是被 Git exclude 的本地资料工作簿。
- 开发工具跨平台：不修改`Makefile`、脚本、CI、`linactl`或长期维护工具入口；本次命令仅用于一次性本地工作簿内容修正和静态验证。
- 测试策略：生产 Go 代码、API DTO、SQL 和前端 UI 均无新增变更，使用静态扫描、工作簿内部 XML 检查和 OpenSpec 校验作为治理验证；无需新增单元测试或 E2E。
- 规则加载：已读取`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/plugin.md`、`.agents/rules/backend-go.md`、`.agents/rules/api-contract.md`、`.agents/rules/data-permission.md`、`.agents/rules/testing.md`、`.agents/rules/documentation.md`、`.agents/rules/database.md`、`.agents/rules/i18n.md`、`.agents/rules/architecture.md`、`.agents/rules/cache-consistency.md`、`.agents/rules/dev-tooling.md`，并使用`goframe-v2`和`spreadsheets`技能。
- 验证：`unzip -p apps/lina-plugins/media/数据看板字典.xlsx xl/worksheets/sheet1.xml | rg -n "单位GB|单位MB/s|单位MB/S|单位Mbps|MB/S|MB/s|Mbps"`无残留；`rg -n "(存储|容量|内存|磁盘|memory|mem|storage|capacity|disk).{0,40}(GB|GiB|G[Bb])|(GB|GiB|G[Bb]).{0,40}(存储|容量|内存|磁盘|memory|mem|storage|capacity|disk)" apps/lina-plugins/media -S --glob '!go.sum' --glob '!数据看板字典.xlsx'`无有效残留；`rg -n "(速度|速率|流量|吞吐|读|写|入站|出站|speed|rate|throughput|bandwidth|read|write|network|disk).{0,50}(MB/S|MB/s|MBps|Mbps)|(MB/S|MB/s|MBps|Mbps).{0,50}(速度|速率|流量|吞吐|读|写|入站|出站|speed|rate|throughput|bandwidth|read|write|network|disk)" apps/lina-plugins/media -S --glob '!go.sum' --glob '!数据看板字典.xlsx'`无有效残留。

### FB-15 反馈修复记录

- 根因：现有上报处理在`STREAM_DELETE`和`SESSION_DELETE`时直接删除流/会话投影，缺少`close_time`生命周期记录；`media_report_stream`仍保存已不需要的`source_id`且缺少源流`protocol_type`；`media_report_session.client_type`仍为字符串；看板`duration`和`play_duration`依赖上报端持续时间或旧采样时间口径，无法表达未关闭对象的实时持续时间；报表投影没有按`close_time`超过`30`天的周期清理任务。
- 修复：`media_report_stream`移除存储字段`source_id`，新增`protocol_type`和`close_time`；`media_report_session.client_type`迁移为`1-mobile`、`2-pc`、`0-未知`整型枚举并新增`close_time`；采集服务在生命周期删除命令中更新`close_time`，新增 ADD 命令会清空旧关闭时间以支持资源重开；看板流和会话统计按`start_time`到`close_time`或当前时间计算持续时间，并只统计`close_time`为空的活跃会话；新增`service/cron`每周清理`close_time`早于`30`天的关闭流和会话；同步更新`数据看板.md`、增量规范、SQL 注释和 DAO/DO/Entity。
- 插件本地规范：`apps/lina-plugins/media/AGENTS.md`不存在，按仓库顶层`AGENTS.md`和命中规则执行。
- 架构边界：修改限定在`media`源码插件报表投影、采集写入、看板读取和插件内定时任务，不修改`apps/lina-core`核心宿主契约，不新增跨模块领域依赖。
- 数据权限：看板接口仍按既有`media:management:query`权限和`platform_only`插件边界执行；本次只改变同一插件报表投影字段和活跃会话过滤，不新增对外读取入口或扩大可见数据范围。
- 缓存一致性：生命周期实时计数仍沿用既有宿主共享 cache；本次只新增数据库`close_time`投影和数据库清理任务，不新增缓存键、失效策略或本地缓存。
- i18n：`media`插件`plugin.yaml`未配置`i18n.enabled: true`，按单语言插件处理；本次修改中文 API 文档源文本、错误码源文案和中文插件文档，不新增插件`manifest/i18n`或`apidoc`翻译资源。
- 数据库：在当前迭代 SQL 中维护同一文件；DDL 使用`DROP COLUMN IF EXISTS`、`ADD COLUMN IF NOT EXISTS`、`DROP INDEX IF EXISTS`、`CREATE INDEX IF NOT EXISTS`和可重入`ALTER COLUMN`迁移，已本地重复执行通过；新增索引覆盖`source_type/node_id/instance_id/status`、`close_time`清理和`tenant_id/node_id WHERE close_time IS NULL`活跃会话计数路径。
- 开发工具跨平台：不修改`Makefile`、脚本、CI或`linactl`入口；`psql`和`make dao p=media`仅用于本地验证和生成代码刷新。
- DI 来源检查：新增`cron`组件由`media/backend/plugin.go`在`system.started` hook 中创建，业务依赖为同一次启动装配中显式创建的`mediasvc.Service`和宿主传入的共享`cacheSvc`、`BizCtx`；`sharedCronSvc`在源码插件进程内只注册一次，定时任务业务逻辑委托`media.Service.CleanupClosedReports`，未在业务路径临时创建新的服务图。
- 性能：看板列表仍保持`10000`条上限；活跃会话数使用一次按`stream_id + protocol_type`聚合查询，租户节点限流使用`tenant_id/node_id/close_time`索引计数，清理任务使用`close_time`索引范围删除，不引入随返回行逐项查询的`N+1`路径。
- 测试策略：本次为功能行为反馈，已更新单元测试覆盖采集关闭标记、协议类型、客户端类型枚举、动态时长、活跃会话过滤、关闭数据清理；未涉及前端页面或 E2E 资产，未触发 E2E 质量审查。
- 规则加载：已读取`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/plugin.md`、`.agents/rules/backend-go.md`、`.agents/rules/api-contract.md`、`.agents/rules/database.md`、`.agents/rules/data-permission.md`、`.agents/rules/cache-consistency.md`、`.agents/rules/documentation.md`、`.agents/rules/testing.md`、`.agents/rules/i18n.md`、`.agents/rules/architecture.md`、`.agents/rules/dev-tooling.md`，并使用`lina-feedback`、`goframe-v2`和`karpathy-guidelines`技能。
- 验证：`psql "postgresql://postgres:postgres@127.0.0.1:5432/linapro?sslmode=disable" -v ON_ERROR_STOP=1 -f apps/lina-plugins/media/manifest/sql/003-add-media-data-collection-server.sql`重复执行通过；`make dao p=media`通过；`go test ./backend/internal/service/collection -count=1`通过；`LINAPRO_TEST_POSTGRES=1 go test ./backend/internal/service/collection -run TestReportRuntimePersistsMetrics -count=1`通过；`go test ./backend/internal/service/media -count=1`通过；`go test ./backend/internal/service/cron -count=1`通过；`go test ./backend/api/media/v1 ./backend/api/mediaopen/v1 -count=1`通过；`go test ./backend/internal/controller/media ./backend/internal/controller/mediaopen -count=1`通过；`go test ./backend -count=1`通过；`openspec validate add-media-data-collection-server --strict`通过。

### FB-16 反馈修复记录

- 根因：`media`插件每周清理任务注册逻辑只在`gcron.AddSingleton`失败时记录日志，`cron.Start(ctx)`不返回错误，`system.started` hook 因此无法感知 cron 注册失败；同时`sync.Once`会在失败后阻止后续启动流程重试注册。
- 修复：将`cron.Service.Start(ctx)`调整为返回`error`，`plugin.go`启动 hook 直接返回该错误；定时任务注册成功后才标记`started`，失败时允许后续启动重试；新增 cron 单元测试覆盖注册失败返回错误和失败后重试成功。
- 插件本地规范：`apps/lina-plugins/media/AGENTS.md`不存在，按仓库顶层`AGENTS.md`和命中规则执行。
- 架构边界：修改限定在`media`源码插件启动 hook 和插件内 cron 服务，不修改`apps/lina-core`核心宿主契约，不新增模块启停策略或跨模块领域依赖。
- 数据权限：本次仅改变插件启动错误传播和定时任务注册状态，不新增 HTTP API、读取入口、写入入口或数据可见性变化。
- 缓存一致性：不修改共享 cache、缓存键、失效、刷新或集群一致性策略。
- i18n：不新增或修改用户可见运行时文案、API 文档源文本或翻译资源。
- 数据库：不新增表、列、索引或 DML；清理任务仍复用既有`CleanupClosedReports`和`close_time`索引。
- 开发工具跨平台：不修改`Makefile`、脚本、CI或`linactl`入口。
- DI 来源检查：未新增服务构造参数或运行期依赖；cron 服务仍由`plugin.go`使用同一启动流程中的`mediasvc.Service`实例创建。
- 性能：启动阶段只执行一次 cron 注册；清理任务业务查询路径不变，不新增周期内额外查询。
- 测试策略：新增`backend/internal/service/cron`单元测试，验证注册错误会返回给调用方，并且失败不会污染后续启动状态。
- 验证：`go test ./backend/internal/service/cron -count=1`通过；`go test ./backend -count=1`通过；`go test ./backend/... -count=1`通过；`openspec validate add-media-data-collection-server --strict`通过。

### FB-17 反馈修复记录

- 根因：本地开发库中`sys_plugin.media.version`和 active release 已经是`v0.1.2`，但仓库`apps/lina-plugins/media/plugin.yaml`仍停留在`v0.1.1`。运行时升级投影会将源码插件文件发现版本低于数据库有效版本判定为`abnormal`，管理端因此显示“插件文件与数据库有效状态不一致，无法自动升级”。
- 修复：将`media/plugin.yaml`版本提升到`v0.1.2`，与当前`003-add-media-data-collection-server.sql`所属数据采集变更和数据库有效版本对齐，恢复源码插件运行时升级状态为可比较的正常版本序列；新增嵌入清单测试，防止携带`003`迁移的插件再次发布为低版本。
- 插件本地规范：`apps/lina-plugins/media/AGENTS.md`不存在，按仓库顶层`AGENTS.md`和命中规则执行。
- 架构边界：修改限定在`media`源码插件清单版本和测试，不修改`apps/lina-core`核心宿主契约，不新增模块启停、跨模块调用或抽象层。
- 数据权限：不新增或修改 HTTP API、读取入口、写入入口、批量动作或租户可见性路径。
- 缓存一致性：不修改共享 cache、缓存键、失效、刷新或集群一致性策略；修复后插件运行时状态仍由宿主既有插件 runtime cache 和数据库 release 状态派生。
- i18n：不新增或修改运行时用户可见文案、API 文档源文本、插件语言包或翻译缓存；本次只修正触发已有管理端异常文案的版本源数据。
- 数据库：不新增或修改 SQL、DAO、DO、Entity；`003-add-media-data-collection-server.sql`仍是当前数据采集变更的唯一迁移文件，版本提升用于触发已安装源码插件执行该迁移。
- 开发工具跨平台：不修改`Makefile`、脚本、CI、`linactl`或长期维护工具入口。
- DI 来源检查：未新增运行期依赖、服务构造函数参数、启动装配、插件宿主服务适配器或`WASM host service`依赖。
- 性能：不修改列表、批量、导出、聚合或插件扫描装配路径，不引入数据库查询频次变化或`N+1`风险。
- 测试策略：本次为源码插件升级治理缺陷，新增`plugin_embed_test.go`验证嵌入`plugin.yaml`版本与`003`迁移保持一致，并运行插件包测试和 OpenSpec 严格校验。
- 规则加载：已读取`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/plugin.md`、`.agents/rules/backend-go.md`、`.agents/rules/database.md`、`.agents/rules/cache-consistency.md`、`.agents/rules/documentation.md`、`.agents/rules/testing.md`、`.agents/rules/i18n.md`、`.agents/rules/architecture.md`、`.agents/rules/data-permission.md`和`.agents/instructions/markdown-format.instructions.md`，并使用`lina-feedback`和`goframe-v2`技能。
- 验证：`psql 'postgresql://postgres:postgres@127.0.0.1:5432/linapro?sslmode=disable' -v ON_ERROR_STOP=1 -P pager=off -c "SELECT plugin_id, version, installed, status, current_state, release_id FROM sys_plugin WHERE plugin_id='media';" -c "SELECT id, release_version, status FROM sys_plugin_release WHERE plugin_id='media' ORDER BY id;"`确认数据库有效版本为`v0.1.2`且 active release 为`v0.1.2`；`go test . -run TestEmbeddedManifestVersionCoversCollectionServerSQL -count=1`通过；`go test ./backend -count=1`通过；`go test ./internal/service/plugin -run 'TestSourcePluginListMarksLowerDiscoveredVersionAbnormal|TestValidateSourcePluginUpgradeReadinessAllowsPendingUpgrade' -count=1`通过；`openspec validate add-media-data-collection-server --strict`通过；`git diff --check`通过。

### FB-18 反馈修复记录

- 根因：已有`collection-client`只覆盖`net-flux` discovery 注册、查询和注销，缺少通过同一 TCP 协议发送`MachineMetric`、`NetworkMetric`、`StreamMetric`和`SessionMetric`的联调入口，无法直接用 client 验证数据上报链路。
- 修复：扩展`apps/lina-plugins/media/hack/tools/collection-client`，新增`report`、`report-close`、`report-cycle`和`smoke`动作；`report`发送实例、网络、流新增和会话新增报文，`report-close`发送会话关闭和流关闭报文，`report-cycle`顺序发送新增和关闭报文，`smoke`串联 discovery 注册、查询、数据上报和注销；新增 CLI 参数控制租户、节点、实例、流、会话、协议和状态；补充单元测试覆盖参数默认值、非法动作和协议枚举解析。
- 插件本地规范：`apps/lina-plugins/media/AGENTS.md`不存在，按仓库顶层`AGENTS.md`和命中规则执行。
- 架构边界：修改限定在`media`源码插件自有开发联调工具，不修改`apps/lina-core`核心宿主契约，不新增模块边界、跨模块调用或运行期服务依赖。
- 数据权限：不新增 HTTP API、读取入口、写入入口或数据权限过滤逻辑；client 只模拟采集端 TCP 上报。
- 缓存一致性：不修改 server 端共享 cache、缓存键、失效或集群策略；client 仅触发现有生命周期事件处理。
- i18n：不新增或修改运行时 UI 文案、API 文档源文本、插件清单、语言包或翻译缓存。
- 数据库：不新增或修改 SQL、表、列、索引、DAO、DO 或 Entity；数据写入仍由既有 TCP server handler 和 report writer 完成。
- 开发工具跨平台：修改的是 Go 开发工具入口，不新增 shell、PowerShell、Makefile、CI 或平台专属命令；使用 Go 标准库 flag、time、strings 和现有`net-flux` client，验证方式为 Go 单元测试。
- DI 来源检查：未新增运行期依赖、服务构造函数参数、启动装配、插件宿主服务适配器或`WASM host service`依赖。
- 性能：client 每次动作发送固定数量报文，不新增列表、批量、聚合或数据库查询路径，不引入`N+1`风险。
- 测试策略：新增工具层单元测试并运行采集 server 所在`collection`包回归测试；未涉及前端页面、E2E 资产或用户可观察 UI，未触发 E2E 质量审查。
- 规则加载：已读取`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/plugin.md`、`.agents/rules/backend-go.md`、`.agents/rules/testing.md`、`.agents/rules/documentation.md`、`.agents/rules/dev-tooling.md`、`.agents/rules/i18n.md`、`.agents/instructions/markdown-format.instructions.md`，并使用`lina-feedback`、`lina-review`、`goframe-v2`和`karpathy-guidelines`技能。
- 验证：`go test ./hack/tools/collection-client -count=1`通过；`go run ./hack/tools/collection-client -h`通过并确认新增动作和参数出现在帮助输出中；`go test ./backend/internal/service/collection -count=1`通过；`openspec validate add-media-data-collection-server --strict`通过；实际 TCP 联调结果见`FB-19`。

### FB-19 反馈修复记录

- 根因：实际启动宿主和`media`采集 TCP server 后，用`collection-client`调试发现两个服务端行为缺陷：`STREAM_DELETE`只更新`close_time/report_time`，没有同步将流投影`status`置为`closed`、`current_active_sessions`置为`0`；注销后再次`LOOKUP`时，Nacos SDK 返回`instance list is empty!`被当成错误向上返回，TCP handler 未发送空`LookupAck`，导致客户端等待到超时。
- 修复：`markReportStreamClosed`在关闭流时同步写入`status=closed`和`current_active_sessions=0`；`discoveryRuntime.Lookup`将 Nacos 空实例错误归一为空`LookupAck`，使注销后查询稳定返回`{}`；补充单元测试覆盖流关闭投影收敛和空发现响应。
- 实际联调：使用本地 Docker PostgreSQL、Redis 和 Nacos，配置`GF_GCFG_PATH=/tmp/linapro-media-tcp.yMCd9d`启动`apps/lina-core`：`go run -tags official_plugins .`，确认日志包含`media collection server started addr=127.0.0.1:1911`和`[tcp-server] is listening on 127.0.0.1:1911`。
- TCP 数据上报验证：`go run ./hack/tools/collection-client -action ping -timeout 5s`返回`pong received`；`go run ./hack/tools/collection-client -action report ... -timeout 10s -settle 2s`发送实例、网络、流和会话上报后，库表确认`media_report_instance.live_streams=1/sessions=1`、`media_report_stream.status=running/current_active_sessions=1/close_time IS NULL`、`media_report_session.close_time IS NULL`；`go run ./hack/tools/collection-client -action report-close ... -timeout 10s -settle 2s`后，库表确认`media_report_instance.live_streams=0/sessions=0`、`media_report_stream.status=closed/current_active_sessions=0/close_time IS NOT NULL`、`media_report_session.close_time IS NOT NULL`。
- TCP 注册发现验证：`go run ./hack/tools/collection-client -action register-lookup -service linapro-media-client-test-e2e-20260621-empty -instance-id linapro-media-client-test-e2e-20260621-empty -node 1 -private-ip 127.0.0.1 -private-port 19191 -timeout 30s -settle 2s`返回包含健康实例的`LookupAck`；`go run ./hack/tools/collection-client -action deregister ... -timeout 10s -settle 2s`成功；注销后`go run ./hack/tools/collection-client -action lookup -service linapro-media-client-test-e2e-20260621-empty -node 1 -timeout 10s -settle 1s`快速返回`{}`，不再超时。
- 插件本地规范：`apps/lina-plugins/media/AGENTS.md`不存在，按仓库顶层`AGENTS.md`和命中规则执行。
- 架构边界：修改限定在`media`源码插件采集服务和插件内联调工具，不修改`apps/lina-core`核心宿主契约，不新增跨模块领域依赖或抽象层。
- 数据权限：不新增 HTTP API、读取入口、写入入口或数据权限过滤逻辑；本次只修复采集端 TCP 写入投影和 discovery TCP 响应语义。
- 缓存一致性：不新增缓存键、失效或集群策略；实例实时计数仍由既有宿主共享 cache 维护，本次只修复数据库关闭投影字段。
- i18n：不新增或修改运行时 UI 文案、API 文档源文本、插件清单、语言包或翻译缓存。
- 数据库：不新增或修改 SQL、表、列、索引、DAO、DO 或 Entity；本地`psql DELETE/SELECT`仅用于清理和核验本次联调测试 ID。
- 开发工具跨平台：修改 Go 开发工具和 Go 后端代码，不新增 shell、PowerShell、Makefile、CI 或平台专属命令。
- 性能：关闭事件和 lookup 仍为单报文固定次数处理；未新增列表、批量、聚合或循环查询路径，不引入`N+1`风险。
- 测试策略：本次为实际 TCP 行为反馈，使用单元测试覆盖服务端边界行为，并通过本地 Docker 依赖完成真实 TCP 手工联调；未涉及前端页面、E2E 资产或用户可观察 UI，未触发 E2E 质量审查。
- 规则加载：已读取`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/plugin.md`、`.agents/rules/backend-go.md`、`.agents/rules/testing.md`、`.agents/rules/documentation.md`、`.agents/rules/dev-tooling.md`、`.agents/rules/i18n.md`、`.agents/rules/architecture.md`、`.agents/rules/data-permission.md`、`.agents/rules/cache-consistency.md`、`.agents/rules/database.md`、`.agents/instructions/markdown-format.instructions.md`，并使用`lina-feedback`、`lina-review`和`goframe-v2`技能。
- 验证：`go test ./backend/internal/service/collection -count=1`通过；`go test ./hack/tools/collection-client -count=1`通过；`go test ./backend -count=1`通过；上述实际 TCP 联调命令和库表查询通过。

### FB-20 反馈修复记录

- 根因：按实际 TCP 上报样本核对 dashboard 接口时发现两个统计口径问题：活跃流和会话的`duration/play_duration`在本地 PostgreSQL `timestamp without time zone`读写后受时区解释影响，接口动态计算认为`start_time`晚于当前时间并返回`0`，没有回退到投影表已保存的上报端显式时长；关闭流后数据库`current_active_sessions`已为`0`，但`protocol_summary.current_sessions`仍保留旧值，dashboard 流列表又从旧摘要汇总回`1`，导致关闭后接口统计不准确。
- 修复：`dashboardElapsedSeconds`在动态时长计算非正数时回退到投影表显式`duration/play_duration`；`markReportStreamClosed`关闭流时保留协议历史累计数并清零`protocol_summary.current_sessions`；dashboard 流列表的协议合并逻辑改为已加载活跃会话计数时以`media_report_session.close_time IS NULL`聚合结果为准，即使结果为空也覆盖旧`current_sessions`。
- 实际联调环境：使用本地 Docker PostgreSQL、Redis 和 Nacos，配置`GF_GCFG_PATH=/tmp/linapro-media-tcp.yMCd9d`启动`apps/lina-core`：`go run -tags official_plugins .`，确认日志包含`media collection server started addr=127.0.0.1:1911`和`http server started listening on [:9120]`。
- TCP 注册发现验证：`go run ./hack/tools/collection-client -action ping -addr 127.0.0.1:1911 -timeout 10s`返回`pong received`；`go run ./hack/tools/collection-client -action register-lookup -service linapro-media-client-stat-e2e-20260622-0915 -instance-id inst-stat-e2e-20260622-0915 -node 1 -private-ip 127.0.0.1 -private-port 19191 -timeout 30s -settle 2s`返回包含同一 service、`group_name=1`和`private_port=19191`的健康实例；`deregister`后再次`lookup`返回`{}`。
- Active 上报库表验证：`go run ./hack/tools/collection-client -action report ... -tenant-id tenant-stat-e2e-20260622-0915 -node-id node-stat-e2e-20260622-0915 -instance-id inst-stat-e2e-20260622-0915 -stream-id stream-stat-e2e-20260622-0915 -session-id session-stat-e2e-20260622-0915 -client-ip 192.0.2.55 -protocol hls -status running -timeout 15s -settle 2s`后，库表精确核对通过：`media_report_instance`为`cpu_allocated=4`、`cpu_load=32.50`、`memory_allocated=1024.00`、`memory_used=512.00`、`disk_io_read=256.00`、`disk_io_write=128.00`、`network_in=4.00`、`network_out=6.00`、`live_streams=1`、`sessions=1`；`media_report_stream`为`protocol_type=HLS`、`status=running`、`resolution=1280x720`、`fps=25.00`、`bitrate=4096`、`packet_loss=0.0100`、`duration=60`、`avg_delay=35`、`total_sessions_lifetime=1`、`current_active_sessions=1`、`protocol_summary.current_sessions=1`；`media_report_session`为`client_type=2`、`play_duration=60`、`current_bitrate=2048`、`current_resolution=1280x720`、`total_link_latency=18`；`media_report_node.node_latency_map`包含`192.0.2.55:18`。
- Active dashboard 接口验证：登录`POST /api/v1/auth/login`获取 token 后，`GET /api/v1/media/dashboard/instances?nodeId=node-stat-e2e-20260622-0915`返回实例`live_streams=1/sessions=1`和资源指标与库表一致；`GET /api/v1/media/dashboard/streams?tenantId=tenant-stat-e2e-20260622-0915&nodeId=node-stat-e2e-20260622-0915&instanceId=inst-stat-e2e-20260622-0915&status=running`返回流`duration=60/current_active_sessions=1/protocol_summary.current_sessions=1`；`GET /api/v1/media/dashboard/sessions?...`返回`HLS`分组`active/session_count=1`且会话`play_duration=60/total_link_latency=18`；`GET /api/v1/media/dashboard/nodes/overview?rootNodeId=node-stat-e2e-20260622-0915&includeEmpty=true`返回节点聚合`live_streams=1/sessions=1/avg_delay=35`。
- Close 上报库表和接口验证：`go run ./hack/tools/collection-client -action report-close ... -timeout 15s -settle 2s`后，库表确认`media_report_instance.live_streams=0/sessions=0`、`media_report_stream.status=closed/current_active_sessions=0/protocol_summary.current_sessions=0/close_time IS NOT NULL`、`media_report_session.close_time IS NOT NULL/play_duration=60`；dashboard 接口确认实例和节点统计归零，流列表返回`status=closed/current_active_sessions=0/protocol_summary.current_sessions=0`，会话列表`HLS`分组为`inactive/session_count=0/active_session_list=[]`。
- 插件本地规范：`apps/lina-plugins/media/AGENTS.md`不存在，按仓库顶层`AGENTS.md`和命中规则执行。
- 架构边界：修改限定在`media`源码插件采集投影写入、dashboard 统计读取和插件内联调工具验证，不修改`apps/lina-core`核心宿主契约，不新增跨模块领域依赖或抽象层。
- 数据权限：不新增 HTTP API、读取入口、写入入口或权限标签；dashboard 接口仍按既有`media:management:query`权限访问，本次只修正同一权限边界内的统计口径。
- 缓存一致性：生命周期实时计数仍由既有宿主共享 cache 维护；本次只在关闭流时同步数据库投影摘要，并在 dashboard 读取时以数据库活跃会话聚合覆盖旧摘要，不新增缓存键、失效策略或本地缓存。
- i18n：`media`插件`plugin.yaml`未配置`i18n.enabled: true`，按单语言插件处理；本次不新增或修改运行时 UI 文案、API 文档源文本、插件清单、语言包或翻译缓存。
- 数据库：不新增或修改 SQL、表、列、索引、DAO、DO 或 Entity；本次使用既有`media_report_*`投影表和`close_time`活跃会话过滤路径，`psql DELETE/SELECT`仅用于清理和核验本次联调测试 ID。
- 开发工具跨平台：长期维护的`collection-client`仍为 Go 工具；本轮未新增 shell、PowerShell、Makefile、CI 或平台专属脚本。实际调试使用的`psql`、`curl`、`node`和`lsof`仅为本地一次性验证命令，不作为交付开发入口。
- DI 来源检查：未新增运行期依赖、服务构造函数参数、启动装配、插件宿主服务适配器或`WASM host service`依赖；dashboard 统计仍通过既有`media.Service`方法读取 DAO 投影，采集写入仍通过既有`reportRuntime`和宿主传入的共享 cache。
- 性能：dashboard 流列表活跃会话统计仍为一次`WHERE stream_id IN (...) AND close_time IS NULL GROUP BY stream_id, protocol_type`聚合查询；空计数覆盖不增加随流行数逐项查询；关闭流额外读取单条 stream 投影用于保留历史协议累计数，固定单资源操作，不引入`N+1`路径。
- 测试策略：本次为功能行为反馈，新增/更新单元测试覆盖时长回退、关闭协议摘要清零和 dashboard 旧摘要覆盖；变更不涉及前端页面、E2E 资产或用户可观察 UI，未触发 E2E 质量审查。
- 规则加载：已读取`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/plugin.md`、`.agents/rules/backend-go.md`、`.agents/rules/api-contract.md`、`.agents/rules/data-permission.md`、`.agents/rules/cache-consistency.md`、`.agents/rules/database.md`、`.agents/rules/testing.md`、`.agents/rules/documentation.md`、`.agents/rules/dev-tooling.md`、`.agents/rules/i18n.md`、`.agents/rules/architecture.md`，并使用`lina-feedback`、`lina-review`和`goframe-v2`技能。
- 验证：`go test ./backend/internal/service/media -count=1`通过；`LINAPRO_TEST_POSTGRES=1 go test ./backend/internal/service/collection -count=1`通过；`go test ./hack/tools/collection-client -count=1`通过；`go test ./backend/internal/controller/media -count=1`通过；`go test ./backend -count=1`通过；`openspec validate add-media-data-collection-server --strict`通过；`git diff --check`通过；`git -C apps/lina-plugins/media diff --check`通过；上述实际 TCP 联调、库表核对和 dashboard HTTP 断言均通过。

### FB-21 反馈修复记录

- 根因：`/api/v1/media/apidocs.html`原本只渲染 Stoplight 的 HTTP OpenAPI 文档，`media`采集 server 又是`system.started` hook 启动的`net-flux`兼容 TCP 服务，不会作为 GoFrame HTTP route 进入`/api/v1/media/openapi.json`，导致用户在插件独立文档页看不到 TCP 采集能力说明，容易误判为接口缺失。
- 修复：在`media`插件自有`apidocs.html`中增加`HTTP API`和`TCP 采集协议`两个视图；HTTP 视图继续加载`/api/v1/media/openapi.json`，TCP 视图说明`collectionServer.enabled`、`collectionServer.addr`、系统心跳、指标上报、生命周期事件、Nacos discovery 命令、处理边界和`collection-client`联调入口；同时保持 OpenAPI JSON 只包含真实 HTTP 路由，不为 TCP 采集 server 伪造 HTTP path。
- 插件本地规范：`apps/lina-plugins/media/AGENTS.md`不存在，按仓库顶层`AGENTS.md`和命中规则执行。
- 架构边界：修改限定在`media`源码插件自有 API 文档页面和插件后端测试，不修改`apps/lina-core`核心宿主契约，不新增宿主扩展点、跨模块领域依赖或抽象层。
- API 契约：`/api/v1/media/apidocs.html`和`/api/v1/media/openapi.json`路由不变；OpenAPI JSON 仍只发布真实 HTTP API，TCP 能力作为同页协议说明展示，不改变请求 DTO、响应 DTO、权限标签或 HTTP 方法。
- 数据权限：不新增或修改数据读取、写入、导出、聚合、批量信息或下拉候选接口；不改变任何业务数据可见性边界。
- 缓存一致性：不修改共享 cache、缓存键、失效、刷新、集群一致性或采集生命周期计数策略。
- i18n：`media`插件`plugin.yaml`未配置`i18n.enabled: true`，按单语言插件处理；本次新增中文插件文档页静态说明，不新增插件`manifest/i18n`或`apidoc`翻译资源。
- 数据库：不新增或修改 SQL、表、列、索引、DML、DAO、DO 或 Entity。
- 开发工具跨平台：不修改`Makefile`、脚本、CI、`linactl`或长期维护工具入口；页面内仅展示既有 Go 联调命令。
- DI 来源检查：未新增运行期依赖、服务构造函数参数、启动装配、插件宿主服务适配器或`WASM host service`依赖。
- 性能：文档页只增加静态 HTML/CSS/JS 和一次本地视图切换，不新增数据库查询、后端数据装配、HTTP 调用瀑布或`N+1`风险。
- 测试策略：本次为用户可观察文档入口反馈，更新`media/backend`路由测试，断言页面包含 TCP 协议说明和`collection-client`入口，并断言 OpenAPI JSON 不包含伪造的 TCP collection HTTP path；未新增 Vben 前端页面、E2E 资产或业务端到端流程，未触发插件 E2E 用例新增。
- 规则加载：已读取`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/plugin.md`、`.agents/rules/api-contract.md`、`.agents/rules/backend-go.md`、`.agents/rules/frontend-ui.md`、`.agents/rules/testing.md`、`.agents/rules/documentation.md`、`.agents/rules/i18n.md`、`.agents/rules/architecture.md`、`.agents/rules/cache-consistency.md`、`.agents/rules/data-permission.md`、`.agents/rules/dev-tooling.md`、`.agents/instructions/markdown-format.instructions.md`，并使用`lina-feedback`、`lina-review`、`goframe-v2`、`frontend-design`和`karpathy-guidelines`技能。
- 验证：`openspec validate add-media-data-collection-server --strict`通过；使用临时 Go workspace 包含`apps/lina-core`与`apps/lina-plugins/media`后执行`go test /Users/wanna/mine/github/wangle201210/linapro/apps/lina-plugins/media/backend -run 'TestMediaPluginOpenAPIDocumentOnlyContainsMediaRoutes|TestMediaPluginAPIDocsPageLoadsMediaDocument' -count=1`通过；使用同一临时 workspace 执行`go test /Users/wanna/mine/github/wangle201210/linapro/apps/lina-plugins/media/backend -count=1`通过。初次尝试在`apps/lina-plugins/media`内直接执行`go test ./backend ...`失败，原因是根`go.work`包含不存在`go.mod`的`./apps/lina-plugins`工作区条目，后续已用临时 workspace 规避并完成验证。

### FB-22 反馈修复记录

- 根因：`collection-client`实际位于`apps/lina-plugins/media/hack/tools/collection-client`，仓库根目录不存在`/hack/tools/collection-client`；同时根`go.work`错误包含`./apps/lina-plugins`和源码插件模块，违反宿主专用`host-only`工作区约定，并导致在`media`插件目录直接执行`go run ./hack/tools/collection-client`时因`apps/lina-plugins/go.mod`不存在而失败。`/api/v1/media/apidocs.html`中的联调命令也未说明插件目录和`GOWORK`边界，容易继续误导使用者。
- 检查结论：`media`插件内`collection-client`功能覆盖完整，已支持`ping`、`register`、`lookup`、`deregister`、`register-lookup`、`report`、`report-close`、`report-cycle`和`smoke`动作；数据上报覆盖`MachineMetric`、`NetworkMetric`、`StreamMetric`、`SessionMetric`以及`STREAM_ADD`、`STREAM_DELETE`、`SESSION_ADD`、`SESSION_DELETE`生命周期事件；单元测试覆盖参数默认值、非法动作、系统包回调、协议枚举和状态枚举解析。
- 修复：将根`go.work`恢复为宿主专用工作区，仅保留`apps/lina-core`和`hack/tools/linactl`；更新`media`插件`apidocs.html`的联调入口，明确先进入`apps/lina-plugins/media`，并将`GOWORK`设为`off`后运行插件内`collection-client`；同步更新插件后端路由测试断言。
- 插件本地规范：`apps/lina-plugins/media/AGENTS.md`不存在，按仓库顶层`AGENTS.md`和命中规则执行。
- OpenSpec：本反馈属于活跃变更`add-media-data-collection-server`的工具入口和验收闭环补充，不新增规范需求；根`go.work`调整与既有`linactl-build-tool-consolidation`规范中根工作区保持`host-only`的要求一致。
- API 契约：不新增 HTTP API、路由、DTO 或权限标签；`/api/v1/media/apidocs.html`仅更新静态说明，`/api/v1/media/openapi.json`仍只包含真实 HTTP 路由。
- i18n：`media`插件`plugin.yaml`未配置`i18n.enabled: true`，按单语言插件处理；本次更新中文插件文档页静态说明，不新增插件`manifest/i18n`或`apidoc`翻译资源。
- 缓存一致性：不修改共享 cache、缓存键、生命周期计数、失效或集群一致性策略。
- 数据权限：不新增或修改数据读取、写入、导出、聚合、批量信息或下拉候选接口；不改变任何业务数据可见性边界。
- 数据库：不新增或修改 SQL、表、列、索引、DML、DAO、DO 或 Entity。
- 开发工具跨平台：`collection-client`仍为 Go 工具；根`go.work`恢复宿主专用范围，插件完整模式继续由`linactl`生成`temp/go.work.plugins`。插件内单独运行该工具时需要按平台设置`GOWORK=off`环境变量后执行`go run ./hack/tools/collection-client`，未新增 shell、PowerShell、Makefile、CI 或平台专属脚本。
- DI 来源检查：未新增运行期依赖、服务构造函数参数、启动装配、插件宿主服务适配器或`WASM host service`依赖。
- 测试策略：本次为治理和工具入口反馈，不改变 TCP 协议业务行为；使用工具单元测试、帮助输出 smoke、后端路由测试、`linactl`工具测试和 OpenSpec 严格校验闭环，不新增 E2E。
- 规则加载：已读取`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/plugin.md`、`.agents/rules/api-contract.md`、`.agents/rules/backend-go.md`、`.agents/rules/testing.md`、`.agents/rules/documentation.md`、`.agents/rules/i18n.md`和`.agents/rules/dev-tooling.md`，并使用`lina-feedback`、`lina-review`、`goframe-v2`和`karpathy-guidelines`技能。
- 验证：`cd apps/lina-plugins/media && GOWORK=off go test ./hack/tools/collection-client -count=1`通过；`cd apps/lina-plugins/media && GOWORK=off go run ./hack/tools/collection-client -h`通过，帮助输出包含全部动作和参数；使用临时 Go workspace 包含`apps/lina-core`与`apps/lina-plugins/media`后执行`go test /Users/wanna/mine/github/wangle201210/linapro/apps/lina-plugins/media/backend -run 'TestMediaPluginOpenAPIDocumentOnlyContainsMediaRoutes|TestMediaPluginAPIDocsPageLoadsMediaDocument' -count=1`通过；`go test ./hack/tools/linactl/... -count=1`通过；`go test ./apps/lina-core/internal/service/plugin/internal/testutil -count=1`通过；`openspec validate add-media-data-collection-server --strict`通过；`git diff --check`和`git -C apps/lina-plugins diff --check`通过；`make status`通过；`make dev`已重启后端和前端，当前`/api/v1/media/apidocs.html`返回内容包含`TCP 采集协议`、`apps/lina-plugins/media`、`GOWORK=off`和`collection-client`，`/api/v1/media/openapi.json`确认不包含伪造的 TCP collection HTTP path。

### FB-23 反馈修复记录

- 根因：远端`10.157.225.55:27153`实测出现`lookup/register-lookup`等待`LookupAck`超时，以及`report-close`后流列表仍显示`running/current_active_sessions=1`。按用户要求切到本地 Docker Nacos 和本地服务复测时，首先发现当前测试替身未实现`cachecap.Service`新增的`GetMany`、`SetMany`、`DeleteMany`方法，导致`collection`和`media`服务测试无法编译运行，阻断本地行为验证。
- 修复：补齐`collection`包`memoryCollectionCache`和`media`包`memoryRouteMemoryCache`的批量缓存接口实现，使测试替身重新满足宿主发布 cache 契约；本地真实 Nacos、TCP server 和 dashboard API 复测未复现远端异常，生产 TCP discovery 和`report-close`流投影逻辑无需改动。
- 本地服务环境：本地 Docker 已运行`linapro-postgres`、`server-redis-1`和`linapro-nacos-test`；创建 Git 忽略的`apps/lina-plugins/media/manifest/config/config.yaml`启用`collectionServer.enabled=true`和`collectionServer.discovery.enabled=true`，执行`make dev`启动服务，确认`temp/lina-core.log`包含`media collection server started addr=127.0.0.1:1911`且`127.0.0.1:1911`已监听。
- TCP discovery 验证：`go run ./hack/tools/collection-client -addr 127.0.0.1:1911 -action ping -timeout 10s`返回`pong received`；`register-lookup -service linapro-media-local-fb23 -instance-id inst-local-fb23 -node 1 -private-ip 127.0.0.1 -private-port 19191 -timeout 30s -settle 2s`返回包含健康实例的`LookupAck`；空服务`lookup`返回`{}`；`deregister`后再次`lookup`返回`{}`，均无超时。
- TCP 上报和 dashboard 验证：使用测试 ID`tenant-local-fb23-20260627232350`、`node-local-fb23-20260627232350`、`inst-local-fb23-20260627232350`、`stream-local-fb23-20260627232350`、`session-local-fb23-20260627232350`执行`report`后，HTTP dashboard 接口返回实例`live_streams=1/sessions=1`、流`status=running/current_active_sessions=1/protocol_summary.current_sessions=1`、会话`session_count=1`、节点`live_streams=1/sessions=1`；执行`report-close`后，接口返回实例`live_streams=0/sessions=0`、流`status=closed/current_active_sessions=0/protocol_summary.current_sessions=0`、`status=running`筛选为空、`status=closed`筛选命中该流、会话分组`inactive/session_count=0`、节点`live_streams=0/sessions=0`。
- 远端判断：本地当前代码和本地 Nacos 不复现远端`LookupAck`超时和流关闭投影滞留，远端更可能是部署版本未包含当前修复、`collectionServer.discovery.*`/Nacos 配置不一致、或远端服务日志中存在 discovery/关闭事件处理错误。
- 插件本地规范：`apps/lina-plugins/media/AGENTS.md`不存在，按仓库顶层`AGENTS.md`和命中规则执行。
- 架构边界：修改仅限`media`源码插件测试替身和 OpenSpec 反馈记录，不修改`apps/lina-core`核心宿主契约，不新增跨模块领域依赖或抽象层。
- 数据权限：不新增或修改 HTTP API、读取入口、写入入口、权限标签或数据权限过滤逻辑；本地 dashboard 验证仍使用既有登录鉴权和`media:management:query`边界。
- 缓存一致性：生产缓存策略不变；测试替身补齐批量接口以匹配宿主共享 cache 契约，不新增缓存键、失效策略、跨实例同步策略或本地缓存路径。
- i18n：不新增或修改运行时 UI 文案、API 文档源文本、插件清单、语言包或翻译缓存。
- 数据库：不新增或修改 SQL、表、列、索引、DML、DAO、DO 或 Entity；本地测试数据通过唯一业务键写入既有`media_report_*`投影表。
- 开发工具跨平台：不修改`Makefile`、脚本、CI、`linactl`或长期维护工具入口；实际联调使用既有 Go 工具`collection-client`和`make dev`，`curl`、`lsof`仅为本地一次性验证命令。
- DI 来源检查：未新增运行期依赖、服务构造函数参数、启动装配、插件宿主服务适配器或`WASM host service`依赖。
- 性能：不修改生产查询或写入路径；测试替身批量方法为固定输入 key 集合的内存 map 操作，不影响运行时性能。
- 测试策略：本次为功能行为复测反馈，使用本地 Docker Nacos 集成测试、服务包单元测试、`collection-client`工具测试和真实 TCP/API 手工验证闭环；未涉及前端页面、E2E 资产或用户可观察 UI，未触发 E2E 质量审查。
- 规则加载：已读取`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/plugin.md`、`.agents/rules/backend-go.md`、`.agents/rules/testing.md`、`.agents/rules/documentation.md`、`.agents/rules/dev-tooling.md`、`.agents/rules/i18n.md`、`.agents/rules/cache-consistency.md`、`.agents/rules/database.md`、`.agents/rules/data-permission.md`和`.agents/rules/architecture.md`，并使用`lina-feedback`、`goframe-v2`和`karpathy-guidelines`技能。
- 验证：使用临时 Go workspace 包含`apps/lina-core`与`apps/lina-plugins/media`后执行`go test ./backend/internal/service/collection ./backend/internal/service/media -count=1`通过；`LINAPRO_TEST_NACOS=1 go test ./backend/internal/service/collection -run TestNacosDiscoveryClientIntegration -count=1`通过；`cd apps/lina-plugins/media && GOWORK=off go test ./hack/tools/collection-client -count=1`通过；`openspec validate add-media-data-collection-server --strict`通过；上述本地真实 TCP discovery、`report`、`report-close`和 dashboard HTTP 验证通过。
