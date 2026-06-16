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
