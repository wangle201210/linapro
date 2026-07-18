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
- [x] FB-24 修复多 Pod 部署下 discovery lookup 复用 Nacos SDK 本地缓存导致注销后仍返回旧实例的问题。
- [x] FB-25 为`media_report_stream`和`media_report_session`补充`device_id`设备维度字段，支持后续按设备分组统计。
- [x] FB-26 将`media`插件采集链路切换为`net-flux`强类型`device_id`协议字段。
- [x] FB-27 移除`Extra`中的设备 ID 冗余写入和读取路径。
- [x] FB-28 补齐视频网关服务监测拓扑页所需的单次聚合统计接口。
- [x] FB-29 将视频网关服务监测拓扑接口入参收敛为单个设备 ID。
- [x] FB-30 删除视频网关服务监测拓扑响应中数据库无法真实获取的字段。
- [x] FB-31 为单设备拓扑查询补充有效的活跃流和活跃会话索引。

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

### FB-24 反馈修复记录

- 根因：`nightly-20260629`在 k8s 多副本环境下已经将 Nacos 注册和注销改为 persistent instance，Nacos 权威状态在注销后也能变为空；但某个已经执行过`Lookup`的 Pod 后续仍返回旧实例。根因是该 Pod 复用了长生命周期 Nacos SDK naming client，`SelectOneHealthyInstance`优先读取 SDK `serviceInfoHolder`本地订阅缓存，且 SDK 默认`UpdateCacheWhenEmpty=false`时空实例列表不会覆盖旧本地缓存。旧配置项`collectionServer.discovery.preloadCache`还和 SDK 字段`NotLoadCacheAtStart`语义相反，后续配置维护容易重新引入磁盘旧缓存。
- 修复：`collection` discovery runtime 改为每个 TCP discovery 命令创建短生命周期 Nacos client，并在`Register`、`Lookup`和`Deregister`结束后关闭，避免跨命令复用进程内订阅缓存；Nacos SDK client 配置固定`UpdateCacheWhenEmpty=true`；将配置项从易误解的`preloadCache`改为`notLoadCacheAtStart`并默认`true`；移除未参与注册、查询和注销的`collectionServer.discovery.groupName`配置，文档明确 net-flux 的`node`映射为 Nacos group；同步更新插件配置样例、k8s 模板和`/api/v1/media/apidocs.html`说明。
- 插件本地规范：`apps/lina-plugins/media/AGENTS.md`不存在，按仓库顶层`AGENTS.md`和命中规则执行。
- 架构边界：修改限定在`media`源码插件 TCP discovery、插件自有配置样例、部署模板、插件文档页和 OpenSpec 记录内，不修改`apps/lina-core`核心宿主契约，不新增宿主扩展点或跨模块领域依赖。
- API 契约：不新增 HTTP API、路由、DTO、权限标签或 OpenAPI JSON path；`/api/v1/media/apidocs.html`仅更新插件静态说明，OpenAPI JSON 仍只包含真实 HTTP 路由。
- 数据权限：不新增或修改业务数据读取、写入、导出、聚合、批量信息或租户可见性路径；TCP discovery 只访问 Nacos 服务发现状态，不读写 LinaPro 业务库。
- 缓存一致性：本次明确规避 Nacos SDK 进程内和磁盘本地缓存对 discovery 查询正确性的影响。Nacos 服务端是 discovery 权威数据源；每次 TCP discovery 命令使用短生命周期 SDK client；client 启动默认不加载磁盘 cache；SDK 收到空服务时允许更新本地状态；跨 Pod 同步由 Nacos 服务端承担，LinaPro 不再依赖任一 Pod 的本地 SDK 缓存。
- i18n：`media`插件`plugin.yaml`未配置`i18n.enabled: true`，按单语言插件处理；本次更新中文插件文档页静态说明，不新增插件`manifest/i18n`或`apidoc`翻译资源。
- 数据库：不新增或修改 SQL、表、列、索引、DML、DAO、DO 或 Entity。
- 开发工具跨平台：不修改`Makefile`、脚本、CI、`linactl`或长期维护工具入口；`collection-client`仍为 Go 工具，本次未改变其 CLI 契约。
- DI 来源检查：未新增运行期依赖、服务构造函数参数、启动装配、插件宿主服务适配器或`WASM host service`依赖；Nacos client 仍只在`collectionServer.discovery.enabled=true`且收到 TCP discovery 命令时由`collection`服务按配置创建。
- 性能：每个 discovery TCP 命令新增一次 Nacos client 创建和关闭成本，但 discovery 注册、注销、查询不是高频 dashboard 列表或批量数据装配路径；该取舍用于保证多 Pod 注销后的正确性，不新增数据库查询或`N+1`风险。
- 测试策略：本次为功能行为反馈，已补单元测试覆盖短生命周期 client、注销后空 lookup、SDK 参数不加载本地缓存且空服务更新缓存；真实 Nacos 集成测试覆盖注册、lookup、deregister 以及同一 lookup runtime 注销后查空。变更不涉及前端业务页面或 E2E 资产，未触发 E2E 质量审查。
- 验证环境说明：本机 Docker daemon 不可用，因此未声称本机 Docker Nacos 验证；真实 Nacos 验证通过 6006 服务器 k8s Pod 内执行当前工作区编译出的测试二进制，直连集群内`linapro-nacos:8848`，覆盖 Nacos 2.x 所需的 gRPC 内部端口路径。
- 规则加载：已读取`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/plugin.md`、`.agents/rules/backend-go.md`、`.agents/rules/api-contract.md`、`.agents/rules/frontend-ui.md`、`.agents/rules/testing.md`、`.agents/rules/documentation.md`、`.agents/rules/dev-tooling.md`、`.agents/rules/i18n.md`、`.agents/rules/cache-consistency.md`、`.agents/rules/data-permission.md`和`.agents/rules/architecture.md`，并使用`lina-feedback`、`lina-review`和`goframe-v2`技能。
- 验证：`GOWORK=/tmp/linapro-media-work.s13xtD/go.work go test ./... -count=1`通过；`GOWORK=/tmp/linapro-media-work.s13xtD/go.work go test ./backend/internal/service/collection -race -run 'TestEventHandlerCreatesFreshDiscoveryClientPerLookup|TestEventHandlerRegistersDiscoveryInstance|TestEventHandlerDeregistersDiscoveryInstance|TestEventHandlerLooksUpDiscoveryInstance|TestEventHandlerWritesEmptyLookupAckWhenDiscoveryHasNoInstance' -count=1`通过；在 6006 服务器 k8s Pod 内执行`LINAPRO_TEST_NACOS=1 LINAPRO_TEST_NACOS_HOST=linapro-nacos LINAPRO_TEST_NACOS_PORT=8848 /tmp/linapro-media-collection.test -test.run TestNacosDiscoveryClientIntegration -test.count=1 -test.v`通过；`GOWORK=/tmp/linapro-media-work.s13xtD/go.work go test ./backend/internal/service/collection ./backend/internal/service/media ./backend -count=1`通过；`GOWORK=/tmp/linapro-media-work.s13xtD/go.work go test ./backend/internal/service/cron ./backend/internal/controller/media ./backend/api/media/v1 -count=1`通过；`GOWORK=/tmp/linapro-media-work.s13xtD/go.work go test ./hack/tools/collection-client -count=1`通过；`openspec validate add-media-data-collection-server --strict`通过；`git -C apps/lina-plugins diff --check`和`git diff --check`通过。

### FB-25 反馈修复记录

- 根因：`media_report_stream`和`media_report_session`只保存租户、节点、实例、流和会话维度，`net-flux`指标包也没有固定`DeviceId`访问器，导致后续统计无法稳定按设备国标 ID 分组；采集端虽可通过`Extra`携带设备标识，但服务端未归一并持久化。
- 修复：在`media_report_stream`和`media_report_session`补充`device_id`字段、DAO/DO/Entity、写入投影、dashboard 请求筛选和响应投影；采集上报从`Extra.device_id`优先、`Extra.deviceId`兼容读取设备 ID，`collection-client`新增`-device-id`并在流和会话上报中写入`extra.device_id`；补充单元测试覆盖归一、持久化、筛选和接口返回。
- OpenSpec：本反馈属于活跃变更`add-media-data-collection-server`的数据采集读模型和统计查询契约补全，已同步更新`design.md`、增量规范和本任务记录。
- 插件本地规范：`apps/lina-plugins/media/AGENTS.md`不存在，按仓库顶层`AGENTS.md`和命中规则执行。
- 架构边界：修改限定在`media`源码插件自有 SQL、采集服务、dashboard 服务、API DTO、联调工具和插件文档，不修改`apps/lina-core`核心宿主契约，不新增跨模块领域依赖或抽象层。
- API 契约：`GET /api/v1/media/dashboard/streams`和`GET /api/v1/media/dashboard/sessions`新增可选`deviceId`筛选，并在流和会话响应中返回`device_id`；只读接口仍使用`GET`和既有`media:management:query`权限。
- 数据权限：dashboard 仍在数据库查询阶段按租户、节点、实例、状态、协议和新增设备维度过滤，不放宽原有可见性边界；TCP 上报只写入采集投影，不通过 HTTP 暴露额外租户数据。
- 缓存一致性：不新增缓存键、缓存失效或本地状态；设备 ID 是数据库投影字段，实例实时计数仍由既有共享 cache 维护。
- i18n：`media`插件`plugin.yaml`未配置`i18n.enabled: true`，按单语言插件处理；本次新增中文 API 文档源文本和插件文档说明，不新增插件`manifest/i18n`或`apidoc`翻译资源。
- 数据库：在当前迭代唯一 SQL 文件`003-add-media-data-collection-server.sql`中追加幂等`ALTER TABLE ... ADD COLUMN IF NOT EXISTS`和设备筛选索引；不写入自增主键，不新增 Seed DML 或 Mock 数据；SQL 重放验证通过。
- 开发工具跨平台：`collection-client`仍为 Go 工具，新增`-device-id`参数使用 Go 标准库 flag 解析，不新增 shell、PowerShell、Makefile、CI 或平台专属脚本。
- DI 来源检查：未新增运行期依赖、服务构造函数参数、启动装配、插件宿主服务适配器或`WASM host service`依赖。
- 性能：设备筛选在数据库侧完成，并新增`(tenant_id, device_id, status)`和`(tenant_id, device_id, protocol_type)`索引；dashboard 流协议统计和会话统计仍为集合化聚合，不新增逐行查询或`N+1`路径。
- 测试策略：本次为功能行为反馈，使用单元测试、PostgreSQL 持久化测试、SQL 重放、真实 TCP 上报和真实 HTTP dashboard 断言闭环；不涉及前端页面或 E2E 资产，未触发 E2E 用例新增。
- 实际 TCP 上报准确性验证：使用本地 PostgreSQL、Redis、Nacos 和`media`采集 TCP server，测试 ID`tenant-device-20260718131107`、`device-device-20260718131107`、`stream-device-20260718131107`、`session-device-20260718131107`执行`report`后，库表确认`media_report_stream.device_id=device-device-20260718131107`、`status=running`、`current_active_sessions=1`、`total_sessions_lifetime=1`、`protocol_type=HLS`；`media_report_session.device_id=device-device-20260718131107`、`client_id=client-device-20260718131107`、`client_ip=192.0.2.57`、`protocol_type=HLS`；聚合查询按`tenant_id/stream_id/device_id`得到`active_sessions=1`、`total_sessions=1`。
- 实际 dashboard 统计验证：登录`POST /api/v1/auth/login`获取 token 后，`GET /api/v1/media/dashboard/streams?tenantId=tenant-device-20260718131107&deviceId=device-device-20260718131107`返回 1 条流，`device_id`、`current_active_sessions=1`、`total_sessions_lifetime=1`和`HLS.current_sessions=1/total_sessions=1`准确；错误`deviceId`返回 0 条流；`GET /api/v1/media/dashboard/sessions?streamId=stream-device-20260718131107&tenantId=tenant-device-20260718131107&deviceId=device-device-20260718131107`返回`stream_info.device_id`和会话`device_id`准确，`HLS.session_count=1`；错误`deviceId`会话计数为 0。
- 关闭上报验证：执行`report-close`后，库表确认流和会话`device_id`保留不变，`media_report_stream.status=closed/current_active_sessions=0/close_time IS NOT NULL`，`media_report_session.close_time IS NOT NULL`，实例`live_streams=0/sessions=0`；dashboard 流列表仍按正确设备命中该流，返回`status=closed/current_active_sessions=0/HLS.current_sessions=0/total_sessions=1`，会话列表返回`HLS`为`inactive/session_count=0/active_session_list=[]`。
- TCP 注册发现验证：`ping`返回`pong received`；使用干净本地 Docker Nacos 重新验证持久实例后，`register-lookup -settle 0s`返回包含`service-device-fast-20260718131107`、`group_name=1`和`private_port=19191`的健康实例，`deregister`成功，注销后`lookup`返回`{}`；延迟数秒后该本地 Nacos 会将无健康检查的持久实例置为`healthy:false`，因此`settle=3s`查询返回空，记录为本地 Nacos 健康状态约束，不影响本次设备字段上报和统计验证。
- 规则加载：已读取`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/plugin.md`、`.agents/rules/backend-go.md`、`.agents/rules/api-contract.md`、`.agents/rules/database.md`、`.agents/rules/testing.md`、`.agents/rules/architecture.md`、`.agents/rules/data-permission.md`、`.agents/rules/cache-consistency.md`、`.agents/rules/i18n.md`、`.agents/rules/dev-tooling.md`、`.agents/rules/documentation.md`和`.agents/instructions/markdown-format.instructions.md`，并使用`lina-feedback`、`lina-review`、`goframe-v2`和`karpathy-guidelines`技能。
- 验证：`GOWORK=off go test ./backend/internal/service/collection -count=1`通过；`LINAPRO_TEST_POSTGRES=1 GOWORK=off go test ./backend/internal/service/collection -run TestReportRuntimePersistsMetrics -count=1`通过；`GOWORK=off go test ./backend/internal/service/media -count=1`通过；`GOWORK=off go test ./backend/internal/controller/media -count=1`通过；`GOWORK=off go test ./backend/api/media/v1 -count=1`通过；`GOWORK=off go test ./hack/tools/collection-client -count=1`通过；`GOWORK=off go test ./backend -count=1`通过；`make lint dir=apps/lina-plugins/media plugins=1`通过；SQL 全量重放通过；真实 TCP 上报、关闭上报、dashboard HTTP 断言和 TCP 注册发现联调通过。

### FB-26 反馈修复记录

- 根因：`net-flux`上游已合并`StreamMetric.device_id`和`SessionMetric.device_id`强类型协议字段，但`media`插件仍依赖合并前的旧 pseudo-version，只能从`Extra.device_id`或`Extra.deviceId`读取设备维度；插件内`collection-client`也只写`Extra`，导致插件代码没有真正消费上游协议字段。升级上游 SDK 后还暴露出`NetworkMetric.throughput`已不在当前协议中，旧代码无法编译。
- 修复：将`media`插件依赖升级到上游合并 commit 对应的`github.com/dellinger2023/net-flux v0.0.8-0.20260718061636-316ef59c9e05`；`StreamMetric`和`SessionMetric`归一化只读取`GetDeviceId()`；`collection-client`在流和会话上报、关闭上报中写入强类型`DeviceId`，不再向`Extra`写设备 ID；`NetworkMetric`节点吞吐改为读取`Extra.throughput_bps`等非设备 metadata，适配当前上游协议。
- OpenSpec：本反馈属于活跃变更`add-media-data-collection-server`的采集协议实现闭环补充，已同步更新`design.md`中协议依赖来源说明和本任务记录。
- 插件本地规范：`apps/lina-plugins/media/AGENTS.md`不存在，按仓库顶层`AGENTS.md`和命中规则执行。
- 架构边界：修改限定在`media`源码插件采集服务、插件内 Go 联调工具、插件 Go module 依赖和 OpenSpec 记录，不修改`apps/lina-core`核心宿主契约，不新增跨模块领域依赖、服务构造参数或抽象层。
- API 契约：不新增或修改 HTTP API、路由、DTO、权限标签或响应结构；dashboard 的`deviceId`筛选和`device_id`响应字段沿用`FB-25`既有契约。
- 数据权限：不新增数据读取或写入入口，不放宽 dashboard 既有`media:management:query`权限；本次仅改变 TCP 上报包中设备维度的来源优先级。
- 缓存一致性：不修改共享 cache、缓存键、失效、刷新或集群一致性策略；实例实时计数仍沿用既有共享 cache。
- i18n：不新增或修改运行时 UI 文案、API 文档源文本、插件清单、语言包或翻译缓存；`media`插件`plugin.yaml`未配置`i18n.enabled: true`，按单语言插件处理。
- 数据库：不新增或修改 SQL、表、列、索引、DML、DAO、DO 或 Entity；设备字段仍复用`FB-25`已新增的投影列和索引。
- 开发工具跨平台：修改的是插件内 Go 工具`hack/tools/collection-client`，仍使用 Go 标准库和`net-flux` client，不新增 shell、PowerShell、Makefile、CI 或平台专属脚本。
- DI 来源检查：未新增运行期依赖、服务构造函数参数、启动装配、插件宿主服务适配器或`WASM host service`依赖；Go module 版本升级不改变插件服务实例来源。
- 性能：设备 ID 归一化是单包常量时间字段读取；`NetworkMetric`吞吐从`Extra`解析单个字符串，不新增数据库查询、循环装配或`N+1`风险；dashboard 聚合查询路径不变。
- 测试策略：本次为功能行为反馈，使用单元测试覆盖强类型`DeviceId`归一化和`Extra`设备 ID 忽略语义，使用插件 Go 全量测试和 lint 覆盖新 SDK 编译面，并启动本地服务走真实 TCP 上报、关闭上报、落库和 dashboard HTTP 统计验证；不涉及前端页面或 E2E 资产，未触发 E2E 用例新增。
- 实际 TCP 上报准确性验证：启动`make dev`后日志确认`media collection server started addr=127.0.0.1:1911`；使用测试 ID`tenant-typed-20260718143600`、`device-typed-20260718143600`、`stream-typed-20260718143600`和`session-typed-20260718143600`执行`collection-client -action report`，库表确认`media_report_stream.device_id=device-typed-20260718143600`、`status=running`、`current_active_sessions=1`、`total_sessions_lifetime=1`、`protocol_type=HLS`，`media_report_session.device_id=device-typed-20260718143600`、`client_type=2`、`protocol_type=HLS`，实例`live_streams=1/sessions=1`。
- 实际 dashboard 统计验证：使用本地管理员登录 token 请求`GET /api/v1/media/dashboard/streams?tenantId=tenant-typed-20260718143600&deviceId=device-typed-20260718143600`返回 1 条流，`device_id`准确，`current_active_sessions=1`且`HLS.current_sessions=1/total_sessions=1`；错误`deviceId`返回空流列表；`GET /api/v1/media/dashboard/sessions?tenantId=tenant-typed-20260718143600&deviceId=device-typed-20260718143600&streamId=stream-typed-20260718143600`返回`stream_info.device_id`准确，`HLS.session_count=1`且活跃会话`device_id`准确。
- 关闭上报验证：执行`collection-client -action report-close`后，库表确认流和会话`device_id`保留不变，流`status=closed/current_active_sessions=0/total_sessions_lifetime=1/close_time IS NOT NULL`，会话`close_time IS NOT NULL`，实例`live_streams=0/sessions=0`；dashboard 流列表按正确设备仍命中关闭流且`HLS.current_sessions=0/total_sessions=1`，会话统计返回`status=inactive/session_count=0/active_session_list=[]`。
- 规则加载：已读取`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/plugin.md`、`.agents/rules/backend-go.md`、`.agents/rules/testing.md`、`.agents/rules/documentation.md`、`.agents/rules/i18n.md`、`.agents/rules/data-permission.md`、`.agents/rules/cache-consistency.md`、`.agents/rules/architecture.md`、`.agents/rules/dev-tooling.md`和`.agents/instructions/markdown-format.instructions.md`，并使用`lina-feedback`和`goframe-v2`技能。
- 验证：`GOWORK=off GONOSUMDB=github.com/dellinger2023/net-flux go test ./backend/internal/service/collection -count=1`通过；`GOWORK=off GONOSUMDB=github.com/dellinger2023/net-flux go test ./hack/tools/collection-client -count=1`通过；`GOWORK=off GONOSUMDB=github.com/dellinger2023/net-flux go test ./... -count=1`通过；`make lint dir=apps/lina-plugins/media plugins=1`通过；`openspec validate add-media-data-collection-server --strict`通过；`git diff --check`和`git -C apps/lina-plugins diff --check`通过；真实 TCP 上报、关闭上报、PostgreSQL 投影核对和 dashboard HTTP 统计核对通过。

### FB-27 反馈修复记录

- 根因：`FB-26`初版为了平滑切换保留了`Extra.device_id`写入和读取兜底，但用户明确确认协议已扩展后不需要在`Extra`中继续保留设备 ID。项目当前没有兼容历史负担，继续保留双来源会让设备维度语义不够明确。
- 修复：`collection-client`流和会话上报、关闭上报只写强类型`DeviceId`，不再写`Extra.device_id`；采集服务端`normalizeStreamMetric`和`normalizeSessionMetric`只从`GetDeviceId()`读取设备维度；单元测试改为验证`Extra.device_id`和`Extra.deviceId`不会被当作设备维度来源。
- OpenSpec：本反馈属于活跃变更`add-media-data-collection-server`对`FB-26`的实现纠偏，已同步修正任务标题和`FB-26`记录中关于`Extra`兼容的表述。
- 插件本地规范：`apps/lina-plugins/media/AGENTS.md`不存在，按仓库顶层`AGENTS.md`和命中规则执行。
- 架构边界：修改限定在`media`源码插件采集服务、插件内 Go 联调工具、单元测试和 OpenSpec 记录，不修改`apps/lina-core`核心宿主契约，不新增跨模块领域依赖或抽象层。
- API 契约：不新增或修改 HTTP API、路由、DTO、权限标签或响应结构；dashboard 的`deviceId`筛选和`device_id`响应字段不变。
- 数据权限：不新增数据读取或写入入口，不改变 dashboard 既有`media:management:query`权限和查询过滤边界；只收敛 TCP 上报包内设备字段来源。
- 缓存一致性：不修改共享 cache、缓存键、失效、刷新或集群一致性策略；实例实时计数仍沿用既有共享 cache。
- i18n：不新增或修改运行时 UI 文案、API 文档源文本、插件清单、语言包或翻译缓存；`media`插件`plugin.yaml`未配置`i18n.enabled: true`，按单语言插件处理。
- 数据库：不新增或修改 SQL、表、列、索引、DML、DAO、DO 或 Entity；设备字段仍复用`FB-25`已新增的投影列和索引。
- 开发工具跨平台：修改的是插件内 Go 工具`hack/tools/collection-client`，仍使用 Go 标准库和`net-flux` client，不新增 shell、PowerShell、Makefile、CI 或平台专属脚本。
- DI 来源检查：未新增运行期依赖、服务构造函数参数、启动装配、插件宿主服务适配器或`WASM host service`依赖。
- 性能：设备 ID 读取改为单个强类型字段读取，比双来源归一更简单；不新增数据库查询、循环装配或`N+1`风险。
- 测试策略：本次为功能行为反馈，使用单元测试验证`Extra`设备 ID 被忽略，使用插件 Go 全量测试和 lint 覆盖编译面，并重启本地服务后走真实 TCP 上报、关闭上报、落库和 dashboard HTTP 统计验证；不涉及前端页面或 E2E 资产，未触发 E2E 用例新增。
- 实际 TCP 上报准确性验证：重启`make dev`后，使用测试 ID`tenant-typed-noextra-20260718145100`、`device-typed-noextra-20260718145100`、`stream-typed-noextra-20260718145100`和`session-typed-noextra-20260718145100`执行新版`collection-client -action report`，该 client 不再写`Extra.device_id`；库表确认`media_report_stream.device_id=device-typed-noextra-20260718145100`、`status=running`、`current_active_sessions=1`、`total_sessions_lifetime=1`、`protocol_type=HLS`，`media_report_session.device_id=device-typed-noextra-20260718145100`、`client_type=2`、`protocol_type=HLS`，实例`live_streams=1/sessions=1`。
- 实际 dashboard 统计验证：使用本地管理员登录 token 请求`GET /api/v1/media/dashboard/streams?tenantId=tenant-typed-noextra-20260718145100&deviceId=device-typed-noextra-20260718145100`返回 1 条流，`device_id`准确，`current_active_sessions=1`且`HLS.current_sessions=1/total_sessions=1`；`GET /api/v1/media/dashboard/sessions?tenantId=tenant-typed-noextra-20260718145100&deviceId=device-typed-noextra-20260718145100&streamId=stream-typed-noextra-20260718145100`返回`stream_info.device_id`准确，`HLS.session_count=1`且活跃会话`device_id`准确。
- 关闭上报验证：执行新版`collection-client -action report-close`后，库表确认流和会话`device_id`保留不变，流`status=closed/current_active_sessions=0/total_sessions_lifetime=1/close_time IS NOT NULL`，会话`close_time IS NOT NULL`，实例`live_streams=0/sessions=0`；dashboard 流列表按正确设备仍命中关闭流且`HLS.current_sessions=0/total_sessions=1`，会话统计返回`status=inactive/session_count=0/active_session_list=[]`。
- 规则加载：已读取`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/plugin.md`、`.agents/rules/backend-go.md`、`.agents/rules/testing.md`、`.agents/rules/documentation.md`、`.agents/rules/i18n.md`、`.agents/rules/data-permission.md`、`.agents/rules/cache-consistency.md`、`.agents/rules/architecture.md`、`.agents/rules/dev-tooling.md`和`.agents/instructions/markdown-format.instructions.md`，并使用`lina-feedback`和`goframe-v2`技能。
- 验证：`GOWORK=off GONOSUMDB=github.com/dellinger2023/net-flux go test ./backend/internal/service/collection -count=1`通过；`GOWORK=off GONOSUMDB=github.com/dellinger2023/net-flux go test ./hack/tools/collection-client -count=1`通过；`GOWORK=off GONOSUMDB=github.com/dellinger2023/net-flux go test ./... -count=1`通过；`make lint dir=apps/lina-plugins/media plugins=1`通过；真实 TCP 上报、关闭上报、PostgreSQL 投影核对和 dashboard HTTP 统计核对通过。

### FB-28 反馈修复记录

- 根因：现有`media`看板只提供节点、实例、流和会话四个分散接口。视频网关服务监测拓扑页需要同屏展示设备、视频基础平台、视频网关、推流协议、租户并发和用户会话，如果前端串联现有接口，会形成按流、协议或租户逐项补查和前端自聚合，无法满足单次页面加载的统计接口边界。
- 修复：新增`GET /api/v1/media/dashboard/topology`，受`media:management:query`保护；请求只接受单个必填`deviceId`作为页面入口；响应一次返回`node_id`、`device_id`、设备拓扑、流拓扑、基础平台节点、视频网关节点、协议复用数、租户并发数和用户会话节点。
- 数据来源：接口只读取`media_report_stream`和`media_report_session`活跃投影，并按需读取`media_report_node`获取节点名称；响应业务字段均来自真实投影或数据库聚合结果，不返回缺少独立投影来源的`device_name`、`video_code`、`server_ip`、`browser`和`operating_system`。
- OpenSpec：本反馈属于活跃变更`add-media-data-collection-server`的数据看板接口补齐，已在增量规范中新增“查询视频网关服务监测拓扑”场景。
- 插件本地规范：`apps/lina-plugins/media/AGENTS.md`不存在，按仓库顶层`AGENTS.md`和命中规则执行。
- 架构边界：修改限定在`media`源码插件的 API DTO、controller、service、测试和 OpenSpec 记录，不修改`apps/lina-core`核心宿主契约，不新增跨模块领域依赖、领域能力或抽象层。
- API 契约：新增只读`GET`资源接口；`g.Meta`包含`dc`和`permission`标签；响应时间字段`generated_at`和用户节点`start_time`均为`Unix timestamp in milliseconds`；请求不提供`pageNum`或`pageSize`，服务端固定有界返回。
- 数据权限：接口沿用`media:management:query`权限和`platform_only`插件边界；`deviceId`在数据库查询阶段注入，不先查全量再内存过滤；聚合统计只覆盖当前设备查询范围内可见的活跃投影。
- 缓存一致性：不新增或修改缓存、缓存键、失效、刷新或跨实例同步策略；实时计数仍由采集链路写入投影，拓扑接口仅查询数据库最新投影。
- i18n：`media`插件`plugin.yaml`未配置`i18n.enabled: true`，按单语言插件处理；本次新增 API 文档源文本使用英文以避免增加全仓 i18n 扫描项，不新增插件`manifest/i18n`或`apidoc`翻译资源。
- 数据库性能：接口按单个`deviceId`先有界读取最多`10000`条活跃流投影，再基于该有界`stream_id`集合执行一次协议租户聚合查询和一次最多`10000`条活跃会话明细查询；查询使用以`device_id`为前导列的活跃流和活跃会话部分索引，不存在按流、协议、租户或用户循环查询数据库的`N+1`路径。
- 开发工具跨平台：本次运行`make ctrl dir=apps/lina-plugins/media/backend`刷新 GoFrame 生成接口声明，但不修改`Makefile`、脚本、CI、`linactl`或长期维护工具入口。
- DI 来源检查：未新增运行期依赖、服务构造函数参数、启动装配、插件宿主服务适配器或`WASM host service`依赖；新增 controller 方法继续复用`ControllerV1.mediaSvc`。
- 测试策略：新增 service 单元测试覆盖单设备 ID 归一化、设备 ID 必填、协议租户聚合、用户会话节点、真实投影字段值和空设备结果；新增 DTO 和 controller 测试固定拓扑 JSON 形状、确认只暴露必填`deviceId`请求入参、确认删除无投影字段并确认不暴露分页包装。变更不涉及前端页面或 E2E 资产，未新增 E2E。
- 规则加载：已读取`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/plugin.md`、`.agents/rules/backend-go.md`、`.agents/rules/api-contract.md`、`.agents/rules/data-permission.md`、`.agents/rules/testing.md`、`.agents/rules/documentation.md`、`.agents/rules/i18n.md`、`.agents/rules/cache-consistency.md`、`.agents/rules/database.md`、`.agents/rules/architecture.md`和`.agents/instructions/markdown-format.instructions.md`，并使用`lina-feedback`和`goframe-v2`技能。
- 验证：`GOWORK=off GONOSUMDB=github.com/dellinger2023/net-flux go test ./backend/internal/service/media -run 'TestDashboardTopology|TestDashboardQueriesReadReportProjections|TestDashboardListQueriesUseFixedUpperBound' -count=1`通过；`GOWORK=off GONOSUMDB=github.com/dellinger2023/net-flux go test ./backend/api/media/v1 -count=1`通过；`GOWORK=off GONOSUMDB=github.com/dellinger2023/net-flux go test ./backend/internal/controller/media -run TestDashboardControllerReturnsFrontendDataShape -count=1`通过；`GOWORK=off GONOSUMDB=github.com/dellinger2023/net-flux go test ./backend/internal/service/media -count=1`通过；`GOWORK=off GONOSUMDB=github.com/dellinger2023/net-flux go test ./backend/internal/controller/media -count=1`通过；`GOWORK=off GONOSUMDB=github.com/dellinger2023/net-flux go test ./backend/api/media/v1 ./backend/api/media -count=1`通过；`GOWORK=off GONOSUMDB=github.com/dellinger2023/net-flux go test ./backend -count=1`通过；`GOWORK=off GONOSUMDB=github.com/dellinger2023/net-flux go test ./backend/... -count=1`通过；`GOWORK=off GONOSUMDB=github.com/dellinger2023/net-flux go test ./... -count=1`通过；`make lint dir=apps/lina-plugins/media plugins=1`通过；`openspec validate add-media-data-collection-server --strict`通过。

### FB-29 反馈修复记录

- 根因：FB-28 初版把拓扑接口做成了通用多维筛选接口，暴露了`nodeId`、`deviceIds[]`、`tenantId`、`protocolType`和`status`，但页面拓扑入口实际只应由单个设备 ID 驱动，额外入参会扩大前端误用面并让接口契约偏离页面统计语义。
- 修复：`GET /api/v1/media/dashboard/topology`请求 DTO 只保留带 GoFrame 必填和长度校验的`deviceId`；service 层同步收敛为单个`DeviceId`，空值复用`CodeMediaDeviceIDRequired`返回参数错误；流投影和会话聚合均在数据库阶段按`device_id`过滤；响应从`device_ids`改为单个`device_id`。
- OpenSpec：已将“查询视频网关服务监测拓扑”场景更新为请求入参仅包含单个必填`device_id`，并保留`10000`条流投影和`10000`条会话明细上限要求。
- 插件本地规范：`apps/lina-plugins/media/AGENTS.md`不存在，按仓库顶层`AGENTS.md`和命中规则执行。
- 架构边界：修改限定在`media`源码插件 API DTO、controller、service、测试和 OpenSpec 记录内，不修改`apps/lina-core`核心宿主契约，不新增跨模块依赖或抽象。
- API 契约：只读接口继续使用`GET`和`media:management:query`权限；请求仅暴露带`v:"required|length:1,128"`校验的`deviceId`，不暴露分页、多设备或多维筛选字段；响应时间字段仍为`Unix timestamp in milliseconds`。
- 数据权限：接口按`platform_only`插件边界和`media:management:query`权限访问；`deviceId`在数据库查询阶段注入，聚合统计只覆盖当前设备活跃投影，不先查全量再内存过滤。
- 缓存一致性：不新增或修改缓存、缓存键、失效、刷新或跨实例同步策略；接口继续直接查询数据库投影。
- i18n：`media`插件`plugin.yaml`未配置`i18n.enabled: true`，按单语言插件处理；本次新增 topology API 文档源文本已使用英文，不新增插件`manifest/i18n`或`apidoc`翻译资源。已运行`make i18n.check`，该全仓命令仍因既有`cms`、`linapro-uidentity-cas`、宿主前端缺失键以及未启用 i18n 的插件中文源文本失败；针对`apps/lina-plugins/media/backend/api/media/v1/media_dashboard_topology.go`的过滤复查无输出，确认本次新增 API 文件未继续增加 i18n 扫描项。
- 数据库性能：接口按单个`deviceId`先有界读取最多`10000`条活跃流投影，再基于该有界`stream_id`集合执行一次协议租户聚合查询和一次最多`10000`条活跃会话明细查询；查询使用 FB-31 补充的单设备活跃流和活跃会话索引，不存在按流、协议、租户或用户循环查询数据库的`N+1`路径，不新增 DAO、DO 或 Entity。
- 开发工具跨平台：本次不修改`Makefile`、脚本、CI、`linactl`或长期维护工具入口；仅执行`gofmt`。
- DI 来源检查：未新增运行期依赖、服务构造函数参数、启动装配、插件宿主服务适配器或`WASM host service`依赖；controller 继续复用`ControllerV1.mediaSvc`。
- 测试策略：更新 service 单元测试覆盖单设备 ID trim、设备 ID 必填、协议租户聚合、用户会话节点和空设备结果；更新 DTO 和 controller 测试确认请求只暴露带必填校验的`deviceId`并固定响应 JSON 形状。变更不涉及前端页面或 E2E 资产，未新增 E2E。
- 规则加载：已读取`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/plugin.md`、`.agents/rules/backend-go.md`、`.agents/rules/api-contract.md`、`.agents/rules/data-permission.md`、`.agents/rules/testing.md`、`.agents/rules/documentation.md`、`.agents/rules/i18n.md`、`.agents/rules/cache-consistency.md`、`.agents/rules/database.md`、`.agents/rules/architecture.md`和`.agents/instructions/markdown-format.instructions.md`，并使用`lina-feedback`和`goframe-v2`技能。
- 验证：`GOWORK=off GONOSUMDB=github.com/dellinger2023/net-flux go test ./backend/internal/service/media -run 'TestDashboardTopology|TestDashboardQueriesReadReportProjections|TestDashboardListQueriesUseFixedUpperBound' -count=1`通过；`GOWORK=off GONOSUMDB=github.com/dellinger2023/net-flux go test ./backend/api/media/v1 -count=1`通过；`GOWORK=off GONOSUMDB=github.com/dellinger2023/net-flux go test ./backend/internal/controller/media -run TestDashboardControllerReturnsFrontendDataShape -count=1`通过；`GOWORK=off GONOSUMDB=github.com/dellinger2023/net-flux go test ./... -count=1`通过；`make lint dir=apps/lina-plugins/media plugins=1`通过；`openspec validate add-media-data-collection-server --strict`通过；`git diff --check`和`git -C apps/lina-plugins diff --check`通过；`rg -n "[一-龥]" media/backend/api/media/v1/media_dashboard_topology.go`无输出；`go run ./hack/tools/linactl i18n.check 2>&1 | rg 'apps/lina-plugins/media/backend/api/media/v1/media_dashboard_topology.go' || true`无输出；`make i18n.check`未通过，失败原因见 i18n 影响记录。

### FB-30 反馈修复记录

- 根因：拓扑响应为贴合页面原型混入了当前上报协议和数据库投影没有的业务字段。`device_name`实际取自`stream_name`，`video_code`实际取自`stream_id`，`server_ip`实际回退为`client_ip`，`browser`和`operating_system`则始终返回空字符串，字段名称与真实数据语义不一致。
- 修复：从 API DTO、service 投影、controller 映射和测试数据中删除设备节点与基础平台节点的`device_name`、网关节点的`video_code`，以及用户会话节点的`server_ip`、`browser`和`operating_system`；同时删除对应的替代值、空值和回退逻辑。响应其余业务字段均来自`media_report_stream`、`media_report_session`、`media_report_node`真实列或数据库聚合结果；`generated_at`和`session_detail_limited`保留为来源明确的服务端响应元数据。
- OpenSpec：本反馈直接收敛活跃变更`add-media-data-collection-server`中的拓扑接口响应契约，已更新“查询视频网关服务监测拓扑”场景，禁止返回没有真实投影来源的字段。
- 插件本地规范：`apps/lina-plugins/media/AGENTS.md`不存在，按仓库顶层`AGENTS.md`和命中规则执行。
- 架构边界：修改限定在`media`源码插件拓扑 API、service、controller、测试和 OpenSpec 记录内，不修改`apps/lina-core`核心宿主契约，不新增跨模块依赖、抽象层或前端适配。
- API 契约：接口路径、`GET`方法、单个必填`deviceId`入参、权限标签、时间字段和有界返回策略不变；仅删除无法真实获取的响应字段，不保留兼容字段或空占位。
- 数据权限：不改变`platform_only`插件边界、`media:management:query`权限或数据库查询条件；接口仍在数据库阶段按单个`device_id`约束流和会话范围，不先查询全量数据后过滤。
- 缓存一致性：无缓存影响；接口继续直接查询数据库投影，不新增缓存、快照、失效、刷新或跨实例同步策略。
- i18n：`media/plugin.yaml`未配置`i18n.enabled: true`，按单语言插件处理；本次仅删除英文 API 文档字段及其说明，不新增或修改运行时 UI、插件语言包、`apidoc`翻译资源或翻译缓存，目标 API 文件中文静态扫描无输出。
- 数据库：不新增或修改 SQL、表、列、索引、DML、DAO、DO 或 Entity；仅依据现有`media_report_stream`、`media_report_session`和`media_report_node`投影收敛响应字段，SQL 幂等性、数据分类、自增主键和软删除语义无影响。
- 性能：数据库查询次数、`10000`条流与会话上限、集合化聚合和内存组装路径不变；删除字段后不新增查询，也不存在按动态结果集循环查询的`N+1`路径。
- 开发工具跨平台：不修改`Makefile`、脚本、CI、`linactl`或代码生成入口；仅执行`gofmt`，无跨平台影响。
- DI 来源检查：未新增运行期依赖、服务构造函数参数、启动装配、插件宿主服务适配器或`WASM host service`依赖；controller 继续复用`ControllerV1.mediaSvc`。
- 测试策略：更新 DTO 序列化测试，明确断言五类无投影字段的 JSON key 不存在；更新 controller 和 service 测试，继续验证设备流数量、协议复用数、租户并发数、会话明细及`client_ip`等真实投影字段。变更不涉及现有前端页面、路由或 E2E 资产，未触发页面 E2E 和截图验证。
- 规则加载：已读取`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/plugin.md`、`.agents/rules/backend-go.md`、`.agents/rules/api-contract.md`、`.agents/rules/data-permission.md`、`.agents/rules/testing.md`、`.agents/rules/documentation.md`、`.agents/rules/i18n.md`、`.agents/rules/cache-consistency.md`、`.agents/rules/database.md`、`.agents/rules/architecture.md`和`.agents/instructions/markdown-format.instructions.md`，并使用`lina-feedback`、`lina-review`和`goframe-v2`技能。
- 验证：`GOWORK=off GONOSUMDB=github.com/dellinger2023/net-flux go test ./backend/api/media/v1 -count=1`通过；`GOWORK=off GONOSUMDB=github.com/dellinger2023/net-flux go test ./backend/internal/controller/media -run TestDashboardControllerReturnsFrontendDataShape -count=1`通过；`GOWORK=off GONOSUMDB=github.com/dellinger2023/net-flux go test ./backend/internal/service/media -run 'TestDashboardTopology|TestDashboardQueriesReadReportProjections|TestDashboardListQueriesUseFixedUpperBound' -count=1`通过；`GOWORK=off GONOSUMDB=github.com/dellinger2023/net-flux go test ./... -count=1 -parallel 1`通过；`make lint dir=apps/lina-plugins/media plugins=1`通过；`openspec validate add-media-data-collection-server --strict`通过；`git diff --check`和`git -C apps/lina-plugins diff --check`通过；`rg -n "DeviceName|VideoCode|ServerIp|Browser|OperatingSystem|device_name|video_code|server_ip|browser|operating_system"`针对拓扑生产代码无输出。

### FB-31 反馈修复记录

- 根因：拓扑接口按要求只接收`deviceId`，流查询条件为`device_id`和未关闭状态，会话查询条件为`device_id`、有界`stream_id`集合和未关闭状态；现有设备索引均以`tenant_id`为前导列，单设备入口未提供`tenant_id`时不能有效支撑该查询路径。
- 修复：为`media_report_stream`新增以`device_id`、`report_time DESC`和`stream_id`组成且限定`close_time IS NULL`的部分索引；为`media_report_session`新增以`device_id`、`stream_id`、`protocol_type`、`tenant_id`、`report_time DESC`和`session_id`组成且限定`close_time IS NULL`的部分索引，分别覆盖活跃流读取、活跃会话明细排序和协议租户聚合。
- OpenSpec：已在活跃变更`add-media-data-collection-server`的“查询视频网关服务监测拓扑”场景中增加单设备活跃投影索引要求。
- 插件本地规范：`apps/lina-plugins/media/AGENTS.md`不存在，按仓库顶层`AGENTS.md`和命中规则执行。
- 架构边界：仅补充`media`插件自有报表投影表的查询索引，不修改`apps/lina-core`、模块边界、跨模块契约、前端适配或运行期依赖。
- API 契约：HTTP 路径、方法、DTO、权限标签、单设备入参和响应结构均不变；索引只优化现有数据库执行路径。
- 数据权限：不改变`platform_only`插件边界、`media:management:query`权限或查询过滤条件；索引不包含新的可见性语义，也不扩大查询范围。
- 缓存一致性：无缓存影响；不新增或修改缓存、快照、失效、刷新和跨实例同步策略，权威数据源仍为数据库报表投影。
- i18n：不修改运行时 UI、API 文档源文本、插件清单、语言包、`apidoc`资源或翻译缓存，无 i18n 影响。
- 数据库：索引追加到当前迭代`003-add-media-data-collection-server.sql`，使用`CREATE INDEX IF NOT EXISTS`满足重复执行幂等性；不新增表、列、DML、自增主键写入或软删除语义，也不需要运行`make dao`。
- 性能：本地 PostgreSQL 强制关闭顺序扫描后的执行计划显示，活跃流查询使用`idx_media_report_stream_active_device`的 Index Scan，协议租户聚合使用`idx_media_report_session_active_device_stream`的 Index Only Scan；查询次数仍为有界流读取、一次聚合和一次明细读取，不存在`N+1`。
- 开发工具跨平台：不修改`Makefile`、脚本、CI、`linactl`或工具入口；本地`psql`只用于 PostgreSQL SQL 重放和查询计划验证，不属于交付脚本变更。
- DI 来源检查：未新增运行期依赖、构造函数参数、启动装配、插件宿主服务适配器或`WASM host service`依赖。
- 测试策略：扩展 service 单元测试逐项验证保留的流、基础平台、网关和用户会话字段值来自测试数据库投影；使用 PostgreSQL SQL 重放、索引定义查询和`EXPLAIN`验证索引可用性。无前端页面、路由、E2E 资产或用户可见文案变化，未触发页面 E2E 和截图验证。
- 规则加载：已读取`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/plugin.md`、`.agents/rules/backend-go.md`、`.agents/rules/api-contract.md`、`.agents/rules/data-permission.md`、`.agents/rules/testing.md`、`.agents/rules/documentation.md`、`.agents/rules/i18n.md`、`.agents/rules/cache-consistency.md`、`.agents/rules/database.md`、`.agents/rules/architecture.md`和`.agents/instructions/markdown-format.instructions.md`，并使用`lina-feedback`、`lina-review`和`goframe-v2`技能。
- 验证：`psql 'postgresql://postgres:postgres@127.0.0.1:5432/linapro?sslmode=disable' -v ON_ERROR_STOP=1 -f apps/lina-plugins/media/manifest/sql/003-add-media-data-collection-server.sql`连续执行通过，新增索引第二次执行均提示已存在并跳过；`pg_indexes`确认两个部分索引定义与 SQL 一致；活跃流查询`EXPLAIN`使用`idx_media_report_stream_active_device`，协议租户聚合`EXPLAIN`使用`idx_media_report_session_active_device_stream`；`GOWORK=off GONOSUMDB=github.com/dellinger2023/net-flux go test ./backend/internal/service/media -run 'TestDashboardTopology|TestDashboardQueriesReadReportProjections|TestDashboardListQueriesUseFixedUpperBound' -count=1`通过；`GOWORK=off GONOSUMDB=github.com/dellinger2023/net-flux go test ./... -count=1 -parallel 1`通过；`make lint dir=apps/lina-plugins/media plugins=1`通过；`openspec validate add-media-data-collection-server --strict`通过；`git diff --check`和`git -C apps/lina-plugins diff --check`通过。
