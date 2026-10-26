## 1. 实现

- [x] 1.1 新增`media`插件采集 server 配置加载和校验，默认关闭并支持监听地址配置。
- [x] 1.2 在`media`插件启动 hook 中按配置异步启动`net-flux` TCP server。
- [x] 1.3 实现采集 server handler，支持心跳响应并接受机器、网络和流指标上报。
- [x] 1.4 补充后端单元测试覆盖配置解析、禁用分支和启动幂等分支。

## 2. 验证

- [x] 2.1 运行覆盖`media`插件后端的 Go 测试。
- [x] 2.2 运行`openspec validate add-media-data-collection-server --strict`。

## Feedback

- [x] FB-49 在summary返回全平台存活节点数alive_node_count，按服务端updated_at最近3分钟判断且与status无关，保持单语句聚合并补齐文档和回归。

- [x] FB-48 在实例列表外层返回原始机器采集时间report_time及数据库updated_at，统一Unix毫秒和空值语义，补齐OpenAPI与验证。

- [x] FB-47 将节点总览及实例列表的网络速率展示字段统一转换为Kbps，覆盖子节点和直属资源，保留存储及原始报文单位并补齐验证。

- [x] FB-46 按用户要求将summary收窄为running且未关闭的流及其未关闭会话，保持单条聚合语句，更新文档与状态切换回归。

- [x] FB-45 补齐流、会话、汇总及拓扑流状态的OpenAPI说明，修正状态与协议示例，明确筛选默认值、活跃口径和数量边界，并验证实际生成文档。

- [x] FB-44 增加全局活跃流与活跃会话总数查询，使用同一条聚合语句，补齐权限、OpenAPI、边界测试和双语说明。

- [x] FB-39 补齐节点总览的原始链路延迟、丢包率、抖动和探测状态，复用实例资源并统一 CPU 核数单位，完成持久化、接口和验证，不计算质量评分或分级。
- [x] FB-40 修复真实上报中实例标识与业务节点号不一致导致直属资源为零的问题，以唯一最新链路归属批量关联实例，覆盖上报顺序、迁移和冲突。
- [x] FB-41 核对四类上报全部字段，在既有投影保存完整原始报文和流/会话命令，补齐删除路径与往返一致性验证，记录不能回填的历史边界。
- [x] FB-42 在现有四类第三方展示查询中返回完整`report_payload`及流/会话上报命令，同步`OpenAPI`并验证空值、精度、嵌套内容和原查询边界。
- [x] FB-43 修正四类`report_payload`的`OpenAPI`示例类型，使生成示例为真实对象，并补充类型及内容回归验证。

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
- [x] FB-32 补齐仅按租户 ID 查询的租户媒体节点拓扑聚合接口。
- [x] FB-33 将`media`插件的`github.com/dellinger2023/net-flux`依赖升级到最新正式 tag，并验证现有 TCP 采集与 Nacos discovery 适配。
- [x] FB-34 通过 media TCP 采集接口向指定 Nacos 注册持久测试实例并上报一组活跃看板数据，保留注册和上报结果供人工检查。
- [x] FB-35 将联调 Nacos 命名空间从误读的`qinju`纠正为`qiniu`，通过 media TCP 接口重新注册持久实例并保留结果。
- [x] FB-36 在 Nacos discovery 集成测试中增加显式 opt-in 的只注册不注销入口，供人工联调保留实例且不影响默认测试。
- [x] FB-37 将只注册不注销的 Nacos 测试隔离到`manual_nacos`构建标签，确保任何常规全量测试都不会编译或执行该测试。
- [x] FB-38 适配`net-flux`最新`v0.0.19` discovery example，复用每个 TCP 连接/服务的`naming.DiscoClient`并让 Lookup 返回完整实例列表，修复上报实例`healthy=false`。

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

### FB-32 反馈修复记录

- 根因：租户监测页需要同屏展示租户活跃会话、活跃流、协议统计及各媒体节点对应统计，现有节点、流、会话接口需要前端多次查询和自行聚合，既无法保证统一统计时点，也会形成调用瀑布。
- 修复：新增`GET /api/v1/media/dashboard/tenant-topology`，受`media:management:query`保护，请求仅包含必填且最长`64`字符的`tenantId`，不接收行政区划；响应返回租户 ID、当前并发会话数、当前活跃流数、协议种类和协议流数量，以及最多`10000`个媒体节点的同口径统计、节点明细截断标记和毫秒时间戳。
- 数据来源与口径：所有业务字段均来自`media_report_stream`、`media_report_session`和`media_report_node`真实投影或数据库聚合；活跃口径为`close_time IS NULL`；协议数量只统计非空协议，协议值对应活跃流数量；节点名称优先使用最新节点投影。数据库未提供租户名称，响应不暴露`tenant_name`，也不包含行政区划字段。
- 插件本地规范：`apps/lina-plugins/media/AGENTS.md`不存在，按仓库顶层`AGENTS.md`和命中规则执行。
- 架构边界：变更限定在`media`源码插件 API、controller、service、鉴权上下文、SQL、测试和 OpenSpec 记录内，不修改`apps/lina-core`核心宿主或前端，不新增跨模块依赖或抽象层。
- 数据权限：宿主调用继续由`media:management:query`、宿主认证和租户中间件治理；铁塔 token 回退鉴权会把 token 用户租户写入插件私有 context，service 在任何报表查询和表结构检查前拒绝租户缺失或请求租户不一致的调用，防止跨租户查询。
- 性能：租户总数使用无明细上限的协议聚合和会话计数保证准确；节点明细使用两次最多`10000`行的数据库聚合和一次节点名称批量投影，不按节点循环查询，不存在`N+1`；数据查询次数固定为五次，另复用既有四表就绪检查，不随返回节点数量增长。
- 数据库：在当前迭代 SQL 中新增`idx_media_report_stream_active_tenant_node_protocol`部分索引，列顺序为`tenant_id`、`node_id`、`protocol_type`且限定`close_time IS NULL`；活跃会话复用`idx_media_report_session_active_tenant_node`。SQL 使用`CREATE INDEX IF NOT EXISTS`，连续执行幂等，不新增表、列、DML、DAO、DO、Entity、自增主键或软删除语义。
- 缓存一致性：接口直接读取数据库最新报表投影，不新增缓存、缓存键、快照、失效、刷新或跨实例一致性路径。
- i18n：`media/plugin.yaml`未配置`i18n.enabled: true`，按单语言插件处理；新增 API 文档和业务错误 fallback 使用英文，不新增插件`manifest/i18n`或`apidoc`资源。当前工作区与独立干净 worktree 执行全仓扫描均为`1509`条既有违规、涉及`164`个文件，定向检查确认本次新增文件和错误文案没有增加违规；`make i18n.check`仍被同一批仓库存量问题阻断。
- 开发工具跨平台：执行`make ctrl dir=apps/lina-plugins/media/backend`生成接口声明和 controller 骨架，并按插件既有显式依赖注入构造模式移除生成器附加的重复构造定义；不修改`Makefile`、脚本、CI、`linactl`或长期维护工具，无跨平台交付影响。
- DI 来源检查：不新增运行期依赖或构造函数参数；controller 复用`ControllerV1.mediaSvc`，Tieta 租户通过请求 context 传递，service 接口注释明确了调用方绑定要求和稳定失败语义。
- 测试策略：新增 API DTO、controller 和 service 测试，覆盖请求仅含`tenantId`、响应精确字段、无租户名称和行政区划、租户及节点聚合准确性、关闭数据和其他租户排除、节点名称最新投影、空数据、参数边界以及 Tieta 租户匹配/缺失/不一致。变更不涉及前端页面、路由或 E2E 资产，未新增页面 E2E。
- 审查：已执行`lina-review`反馈级审查，范围来自主仓和插件子仓`git status --short`及未跟踪文件展开，共`15`个文件；发现的参数/授权晚于表检查、service 接口注释缺少租户隔离调用边界两项问题均已修复，复审未发现阻塞问题。
- 规则加载：已读取`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/architecture.md`、`.agents/rules/api-contract.md`、`.agents/rules/backend-go.md`、`.agents/rules/data-permission.md`、`.agents/rules/testing.md`、`.agents/rules/i18n.md`、`.agents/rules/plugin.md`、`.agents/rules/cache-consistency.md`、`.agents/rules/database.md`、`.agents/rules/documentation.md`和`.agents/instructions/markdown-format.instructions.md`，并使用`lina-feedback`、`lina-review`和`goframe-v2`技能；前端、E2E、开发工具长期入口和宿主插件 README 均无变更影响。
- 验证：`GOWORK=off GONOSUMDB=github.com/dellinger2023/net-flux go test ./... -count=1 -parallel 1`通过；`make lint dir=apps/lina-plugins/media plugins=1`通过且为`0 issues`；PostgreSQL SQL 脚本连续执行通过且第二次提示新增索引已存在；`pg_indexes`确认两个租户活跃部分索引定义正确；`EXPLAIN`确认活跃流节点协议聚合使用`idx_media_report_stream_active_tenant_node_protocol`，活跃会话节点聚合使用`idx_media_report_session_active_tenant_node`；`openspec validate add-media-data-collection-server --strict`通过；主仓和插件子仓`git diff --check`通过。

### FB-33 反馈修复记录

- 根因：`media`插件仍依赖`github.com/dellinger2023/net-flux v0.0.8-0.20260718061636-316ef59c9e05`伪版本，上游已经发布`v0.0.18`正式 tag。对比`316ef59c9e05..v0.0.18`确认`gen`上报协议、`network.NewTcpServer`和`network.EventHandler`接口未变，新增内容主要位于 naming、balancer、cache 及 TCP client 重连修复，现有采集 server 无需 API 迁移；但用户提供的`collectionServer.discovery.host`为`http://10.157.225.139/`完整 URL，既有适配直接把该值写入 Nacos SDK 的`IpAddr`，SDK 会在尾部斜杠后再次拼接端口并生成错误地址。
- 修复：将`net-flux`依赖升级到`v0.0.18`并刷新`go.sum`；新增 Nacos server config 规范化，将裸主机名或`http/https` URL 映射为 SDK 独立的`Scheme`、`IpAddr`和`Port`字段，拒绝 URL 凭据、查询参数、片段、额外路径和与`collectionServer.discovery.port`冲突的 URL 端口；配置加载阶段提前校验，避免首次 discovery 报文到达时才暴露错误。
- OpenSpec：本反馈直接维护活跃变更`add-media-data-collection-server`引入的 TCP 采集与 Nacos discovery 实现，不新增能力或改变既有采集协议规范，仅在`tasks.md`记录依赖升级和实现适配。
- 插件本地规范：`apps/lina-plugins/media/AGENTS.md`不存在，按仓库顶层`AGENTS.md`和命中规则执行；根目录不存在`.contributing`，变更全部限定在`media`插件与对应 OpenSpec 记录内，未修改`apps/lina-core`、`apps/lina-vben`或`hack`。
- 架构边界：继续由`media`源码插件内部`collection`服务持有 TCP 与 Nacos 适配，不新增模块、跨模块契约、抽象层或宿主能力；URL 解析是 Nacos SDK 边界前的局部值适配。
- API、数据权限与数据库：不新增或修改 HTTP API、路由、DTO、权限标签、数据读写入口、租户可见性、SQL、DAO、DO、Entity、索引或查询路径，无数据权限、数据库性能和`N+1`影响。
- 缓存一致性：虽然新版`net-flux`新增 cache 包，但`media`未导入或启用该能力；既有宿主共享 cache、生命周期计数、缓存键、失效和跨实例策略均未修改，无缓存一致性影响。
- i18n：不修改运行时用户可见文案、前端 UI、API 文档源文本、插件清单、语言包、`apidoc`资源或翻译缓存，无 i18n 影响。
- 开发工具跨平台：不修改`Makefile`、脚本、CI、`linactl`或开发入口；依赖与代码继续使用插件既有`Go 1.25`基线。定向 lint 初次因环境中的`go1.25.8`可执行文件错误搭配`go1.26.4 GOROOT`失败，显式设置匹配的`GOROOT=/Users/wanna/sdk/go1.25.8`和`GOTOOLCHAIN=local`后通过，确认失败与代码无关。
- DI 来源检查：未新增运行期依赖、服务构造函数参数、启动装配、插件宿主服务适配器或`WASM host service`依赖；仅升级 collection 服务已经直接使用的模块版本。
- 测试策略：新增单元测试覆盖截图中的 URL host 规范化和非法 URL path 启动期拒绝；真实 Nacos 集成测试新增命名空间、用户名和密码环境参数，并使用用户提供的主机、端口、`qinju`命名空间及凭据完成持久实例注册、查询和注销。变更不涉及页面、路由、前端接口联动或 E2E 资产，未触发页面 E2E 和截图验证。
- 审查：已执行`lina-review`反馈级审查，范围来自主仓和插件子仓`git status --short`、已跟踪差异及未跟踪文件扫描，共`7`个文件；未发现严重或警告问题。HTTP Nacos 已完成真实验证，HTTPS Nacos 未单独实测，保留为非阻塞剩余风险。
- 规则加载：已读取`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/architecture.md`、`.agents/rules/plugin.md`、`.agents/rules/backend-go.md`、`.agents/rules/testing.md`、`.agents/rules/documentation.md`和`.agents/instructions/markdown-format.instructions.md`，并使用`lina-feedback`、`lina-review`、`goframe-v2`和`karpathy-guidelines`技能；API contract、数据权限、数据库、缓存一致性、i18n、前端 UI 和开发工具规则域确认无变更影响。
- 验证：`go list -m -versions github.com/dellinger2023/net-flux`和上游 Git tags 均确认最新正式 tag 为`v0.0.18`；`GOWORK=off go test ./backend/internal/service/collection ./hack/tools/collection-client -count=1`通过；使用给定 Nacos 配置执行`TestNacosDiscoveryClientIntegration`通过，完成注册、查询和注销，测试后经认证只读扫描未发现`linapro-media-collection-test-`残留服务；`GOWORK=off go test ./... -count=1 -parallel 1`通过；`GOROOT=/Users/wanna/sdk/go1.25.8 GOTOOLCHAIN=local make lint dir=apps/lina-plugins/media plugins=1`通过；`openspec validate add-media-data-collection-server --strict`通过；主仓与插件子仓`git diff --check`通过。

### FB-34 反馈执行记录

- 现状：执行前本机没有进程监听`1911`，本地数据库中的`media`插件处于未安装状态。为确保联调经过项目真实 TCP 采集入口，而不是直接写 Nacos 或数据库，在`temp/20260818/nacos-report-155822/`创建忽略跟踪的临时运行配置，关闭本机缺失 Redis 时的集群模式、自动安装并启用`media`插件，并使用用户提供的 Nacos 主机、`8848`端口、`qinju`命名空间和账号配置启动临时官方插件后端。
- Nacos 注册：通过`collection-client -action register`向`127.0.0.1:1911`发送 discovery 注册报文，保留命名空间`qinju`、分组`1`、服务`linapro-media-report-20260818-160000`、实例`127.0.0.1:19191`；Nacos 查询确认`enabled=true`且`ephemeral=false`。该地址是仅供记录检查的本机测试端点，Nacos 服务端无法探测，因此显示`healthy=false`，不代表注册失败。
- 指标上报：通过`collection-client -action report`发送实例、网络、流、会话指标及流/会话新增事件。保留实例`linapro-media-report-20260818-160000`、租户`tenant-nacos-demo`、设备`device-nacos-demo-20260818-160000`、节点`node-nacos-demo-1`、流`stream-nacos-demo-20260818-160000`和会话`session-nacos-demo-20260818-160000`；数据库确认实例`live_streams=1`、`sessions=1`，流状态为`running`，流和会话`close_time`均为空。
- 保留策略：遵循用户“不用清理”的要求，未发送 discovery deregister、Nacos 服务删除、`STREAM_DELETE`、`SESSION_DELETE`或数据库删除；只停止占用本机`1911`和`19120`端口的临时后端进程。停止后再次查询确认 Nacos 持久实例、数据库实例计数、活跃流和活跃会话仍存在，本地临时配置和构建产物保留在已忽略的`temp/`目录。
- 插件与架构边界：未修改`apps/lina-core`、`apps/lina-vben`、`hack`或`media`生产代码和运行配置；联调复用`media`插件既有 TCP handler、Nacos discovery adapter、report writer 和宿主发布的 cache 服务，不新增模块、接口、依赖或抽象。
- API、数据权限与数据库：不改变 HTTP API、DTO、权限标签、租户可见性、SQL、DAO、DO、Entity 或索引。此次按用户授权向本地开发数据库写入命名唯一的测试投影，并向指定 Nacos 写入一个持久实例；不存在批量查询、循环数据库访问或`N+1`影响。
- 缓存一致性：临时后端使用`cluster.enabled=false`单机分支，生命周期事件通过既有宿主 cache 能力维护实例实时计数；未修改缓存键、TTL、失效、共享后端或分布式策略。
- i18n：不修改运行时文案、前端 UI、API 文档、插件清单、语言包或翻译缓存，无 i18n 影响。
- 开发工具跨平台：仅使用现有 Go 构建、`collection-client`、`curl`和`psql`执行联调；没有修改`Makefile`、脚本、CI、`linactl`或跨平台入口。临时构建显式使用项目要求的`Go 1.25.8`工具链和`temp/go.work.plugins`官方插件工作区。
- DI 来源检查：不修改运行期依赖或构造函数；临时后端按正常启动装配发布 plugin config、cache 和数据库服务，`media`插件在`system.started`hook 中启动 collection server。
- 测试策略：这是后端 TCP/Nacos/数据库端到端联调，不涉及页面、路由、表单、权限交互或 E2E 测试资产，未触发页面 E2E 与截图验证。验证链路覆盖 TCP 注册报文、Nacos 只读查询、四类指标与生命周期报文、数据库投影查询以及停止临时后端后的持久性复查。
- 审查：已执行`lina-review`反馈级审查，范围为`FB-34`的 OpenSpec 执行记录、忽略跟踪的临时配置和外部写入证据；确认核心目录无差异、临时配置未进入 Git、联调结果与记录一致，未发现严重或警告问题。
- 规则加载：已读取`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/testing.md`、`.agents/rules/documentation.md`和`.agents/instructions/markdown-format.instructions.md`，并使用`lina-feedback`和`lina-review`技能；插件、后端 Go、API contract、数据权限、数据库、缓存一致性、i18n、前端 UI 和开发工具文件均无变更影响。
- 验证：临时后端日志确认`media collection server started addr=:1911`并接收 discovery、MachineMetric、NetworkMetric、StreamMetric 和 SessionMetric 报文；Nacos HTTP 查询确认服务`1@@linapro-media-report-20260818-160000`包含一个持久实例；PostgreSQL 查询确认实例、节点、流和会话投影写入成功；临时后端停止后`1911`和`19120`无监听，Nacos 与数据库复查结果仍保留。

### FB-35 反馈修复记录

- 根因：上一轮根据截图把目标 Nacos 命名空间误读为`qinju`，并按该错误值启动临时 media 后端，因此用户在正确的`qiniu`命名空间中无法看到服务。Nacos 命名空间接口确认`qiniu`真实存在，纠正前在`qiniu`查询服务`linapro-media-report-20260818-160000`返回`0`个实例。
- 修复：将忽略跟踪的临时 media 配置`collectionServer.discovery.namespace`改为`qiniu`，重新启动既有官方插件临时后端，并再次通过`collection-client -action register`向 media TCP `127.0.0.1:1911`发送同名 discovery 注册报文；未直接调用 Nacos 写接口。
- 验证结果：Nacos 在`qiniu`命名空间、分组`1`下返回服务`1@@linapro-media-report-20260818-160000`和实例`127.0.0.1#19191#DEFAULT#1@@linapro-media-report-20260818-160000`，状态为`enabled=true`、`ephemeral=false`。实例使用本机回环测试地址，Nacos 服务端无法探测，因此显示`healthy=false`。
- 数据边界：Nacos discovery 只保存服务与实例注册，不保存`MachineMetric`、`NetworkMetric`、`StreamMetric`或`SessionMetric`业务指标；上一轮这些指标已经通过 TCP 上报并保留在本地 PostgreSQL 的`media_report_*`读模型中。
- 保留策略：遵循用户“不用清理”的要求，没有注销或删除`qiniu`新实例，也没有删除先前误写到`qinju`的旧实例和数据库指标；只停止本地临时后端。停止后再次查询确认`qiniu`持久实例仍存在，`1911`和`19120`均无监听进程。
- 影响分析：未修改生产代码、HTTP API、DTO、权限、数据库结构、缓存实现、运行期依赖、前端 UI、i18n、构建脚本或跨平台入口；只更新 OpenSpec 执行记录和已忽略的临时配置，无数据权限、缓存一致性、`N+1`、DI 或开发工具交付影响。
- 测试策略：使用 Nacos 命名空间只读列表、纠正前实例查询、media TCP 注册、服务端接收日志、纠正后实例查询和停止进程后的持久性复查完成闭环；不涉及页面或 E2E 测试资产，未触发页面 E2E 与截图验证。
- 审查：已执行`lina-review`反馈级审查，范围为`FB-35`的 OpenSpec 记录、忽略跟踪的临时配置和 Nacos 查询证据；确认`qiniu`值与持久实例反查一致、生产文件无新增变更、本地无遗留监听进程，未发现严重或警告问题。
- 规则加载：已读取`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/testing.md`、`.agents/rules/documentation.md`和`.agents/instructions/markdown-format.instructions.md`，并使用`lina-feedback`和`lina-review`技能；其余规则域确认无文件变更影响。

### FB-36 反馈修复记录

- 根因：现有`TestNacosDiscoveryClientIntegration`固定执行注册、查询、注销和服务删除，适合自动回归但无法承载人工联调需要保留 Nacos 实例的场景；此前保留实例只能临时运行`collection-client`命令，没有固化为可重复执行的测试入口。
- 修复：在`collection_discovery_test.go`新增`TestNacosDiscoveryClientRegisterWithoutCleanup`。测试只创建 discovery runtime 并调用`Register`，不创建 deregister runtime、不调用`Deregister`或`deleteNacosTestService`；`runtime.Close()`仅释放 Nacos SDK client，不会注销由服务端强制注册为`ephemeral=false`的持久实例。
- 安全门禁：新测试默认跳过，只有显式设置`LINAPRO_TEST_NACOS_KEEP=1`才执行并接受外部残留；支持`LINAPRO_TEST_NACOS_KEEP_SERVICE`、`LINAPRO_TEST_NACOS_KEEP_IP`和`LINAPRO_TEST_NACOS_KEEP_PORT`指定保留实例，并复用现有 Nacos host、port、namespace、username、password 环境参数。新增`LINAPRO_TEST_NACOS_NODE`设置分组节点，并校验节点为正数、实例端口位于`1..65535`。
- 配置复用：抽取`newNacosIntegrationConfig`统一两个真实 Nacos 集成测试的配置读取、临时日志目录和缓存目录，避免保留测试与清理测试复制环境参数装配逻辑。
- 实测结果：使用`qiniu`命名空间运行新测试通过，保留命名空间`qiniu`、分组`1`、服务`linapro-media-kept-test-20260819-144000`和实例`127.0.0.1:19191`；Nacos HTTP 反查确认实例`enabled=true`、`ephemeral=false`。实例使用回环测试地址，Nacos 服务端无法探测，因此`healthy=false`。
- 插件与架构边界：修改限定在`media`插件既有 Nacos 测试文件和 OpenSpec 记录，不修改生产代码、`apps/lina-core`、模块契约、运行期依赖或启动装配；`apps/lina-plugins/media/AGENTS.md`不存在，继续遵守仓库顶层规范。
- API、数据权限、数据库与缓存：不修改 HTTP API、DTO、权限、业务数据读写、SQL、DAO、DO、Entity、索引、缓存键、TTL、失效或分布式策略，无数据权限、数据库性能、`N+1`和缓存一致性影响。
- i18n 与前端：不修改运行时文案、API 文档、插件清单、语言包、翻译缓存、前端页面或交互，无 i18n 和前端影响。
- 开发工具跨平台：不修改`Makefile`、脚本、CI、`linactl`或工具入口；测试使用标准 Go 环境变量，适用于支持 Go 的平台。首次测试因本机`GOROOT`错误指向`go1.26.4`而未进入编译，显式设置匹配`Go 1.25.8`的`GOROOT`后通过。
- DI 来源检查：未新增生产运行期依赖、构造函数参数、插件宿主服务适配器或`WASM host service`依赖；测试继续直接构造包内 discovery runtime。
- 测试策略：默认运行验证新测试明确跳过且不产生外部残留；opt-in 运行验证真实 Nacos 注册成功并刻意保留实例；media 全量 Go 测试和定向 lint 覆盖当前工作区。变更不涉及页面或 E2E 资产，未触发页面 E2E 与截图验证。
- 审查：已执行`lina-review`反馈级审查，范围为新增 Nacos 保留测试、共享测试配置 helper、OpenSpec 记录和真实 Nacos 反查证据；确认 opt-in 门禁、默认无副作用、显式保留语义和验证覆盖符合本次用户要求，未发现严重或警告问题。操作者显式启用后需自行管理外部残留，这是该入口的预期行为。
- 规则加载：已读取`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/plugin.md`、`.agents/rules/backend-go.md`、`.agents/rules/testing.md`、`.agents/rules/documentation.md`和`.agents/instructions/markdown-format.instructions.md`，并使用`lina-feedback`、`lina-review`、`goframe-v2`和`karpathy-guidelines`技能；其余规则域确认无变更影响。
- 验证：`GOROOT=/Users/wanna/sdk/go1.25.8 GOTOOLCHAIN=local GOWORK=off go test ./backend/internal/service/collection -count=1 -v`通过且保留测试默认跳过；设置完整 Nacos 环境和`LINAPRO_TEST_NACOS_KEEP=1`后定向运行`TestNacosDiscoveryClientRegisterWithoutCleanup`通过；Nacos API 反查保留实例成功；`GOROOT=/Users/wanna/sdk/go1.25.8 GOTOOLCHAIN=local GOWORK=off go test ./... -count=1 -parallel 1`通过；`GOROOT=/Users/wanna/sdk/go1.25.8 GOTOOLCHAIN=local make lint dir=apps/lina-plugins/media plugins=1`通过且为`0 issues`；`openspec validate add-media-data-collection-server --strict`和主仓、插件子仓`git diff --check`通过。

### FB-37 反馈修复记录

- 根因：`FB-36`初版仅依赖`LINAPRO_TEST_NACOS_KEEP=1`在测试运行期决定跳过；正常环境不会执行，但如果全量测试进程误带该环境变量，保留测试仍会被发现并执行，无法满足“只能手动执行、全量测试绝不执行”的强隔离要求。
- 修复：将`TestNacosDiscoveryClientRegisterWithoutCleanup`从常规`collection_discovery_test.go`移动到独立`collection_discovery_manual_test.go`，文件使用`//go:build manual_nacos`构建约束。常规测试不传该 tag 时文件不参与编译，测试列表中不存在该方法。
- 双重门禁：即使操作者显式加入`-tags=manual_nacos`，测试仍要求`LINAPRO_TEST_NACOS_KEEP=1`才执行，防止带 tag 的静态检查或误运行直接写入共享 Nacos。两个条件同时满足后，测试按用户要求只注册持久实例，不执行注销或删除。
- 手工入口：运行命令必须包含`go test -tags=manual_nacos ./backend/internal/service/collection -run '^TestNacosDiscoveryClientRegisterWithoutCleanup$' -count=1 -v`，并提供`LINAPRO_TEST_NACOS_KEEP=1`及 Nacos 连接、命名空间和保留实例环境参数。
- 实测结果：使用双重门禁在`qiniu`命名空间重新运行服务`linapro-media-kept-test-20260819-144000`成功，复用并保留`group=1`、`127.0.0.1:19191`持久实例，没有注销或删除。
- 影响分析：仅调整`media`插件测试文件组织和 OpenSpec 记录，不修改生产代码、API、权限、数据库、缓存、DI、i18n、前端或开发工具入口，无数据权限、缓存一致性、数据库性能、`N+1`和跨平台交付影响；`apps/lina-plugins/media/AGENTS.md`不存在。
- 测试策略：在环境中故意设置`LINAPRO_TEST_NACOS_KEEP=1`但不传 tag，执行`go test -list '^TestNacosDiscoveryClient'`确认列表中不存在保留测试；传 tag 但不设置确认变量，定向测试明确跳过；同时提供 tag 和确认变量，定向真实 Nacos 测试通过。随后在误带确认变量的环境中运行 media 全量测试通过，证明全量测试不会执行保留入口；常规 lint 和带`manual_nacos`tag 的 lint 均为`0 issues`。变更不涉及页面或 E2E 资产，未触发页面 E2E 与截图验证。
- 审查：已执行`lina-review`反馈级审查，范围为常规测试文件、`manual_nacos`构建标签测试文件、OpenSpec 记录和 Nacos 保留实例证据；确认常规测试不编译手工文件、双门禁有效、带 tag 代码可编译且静态检查通过，未发现严重或警告问题。显式手工执行后的 Nacos 残留由操作者按预期管理。
- 规则加载：已读取`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/plugin.md`、`.agents/rules/backend-go.md`、`.agents/rules/testing.md`、`.agents/rules/documentation.md`和`.agents/instructions/markdown-format.instructions.md`，并使用`lina-feedback`、`lina-review`、`goframe-v2`和`karpathy-guidelines`技能；其余规则域确认无变更影响。
- 验证：`LINAPRO_TEST_NACOS_KEEP=1 go test ./backend/internal/service/collection -list '^TestNacosDiscoveryClient'`仅列出清理型`TestNacosDiscoveryClientIntegration`；`go test -tags=manual_nacos ./backend/internal/service/collection -run '^TestNacosDiscoveryClientRegisterWithoutCleanup$' -count=1 -v`在缺少确认变量时跳过；完整手工环境下同一命令通过；`LINAPRO_TEST_NACOS_KEEP=1 GOWORK=off go test ./... -count=1 -parallel 1`通过；常规及`GOFLAGS='-tags=manual_nacos'`的`make lint dir=apps/lina-plugins/media plugins=1`均通过且为`0 issues`。

### FB-38 反馈修复记录

- 根因：media 之前直接创建 Nacos SDK naming client，并在每个 Register、Deregister、Lookup 报文处理后关闭 client。上游`net-flux`最新`v0.0.19` example 已明确通过`discoClientFactory`按 TCP `conn.ID()`和服务名复用一个`naming.DiscoClient`；旧实现关闭 client 后无法维持上游实例心跳，导致 Nacos 中实例出现`healthy=false`。旧 Lookup 还调用`GetServiceInstanceByGroup`，只能返回单个实例。
- 修复：将`github.com/dellinger2023/net-flux`升级到`v0.0.19`，新增`disco_cli_factory.go`，按连接 ID和服务名缓存`naming.DiscoClient`，Register/Deregister/Lookup复用对应 client，TCP `OnClose`统一释放该连接的全部 client；Deregister成功后释放对应服务 client。Lookup改为`GetServiceInstances(serviceName, group, []string{})`，清洗并返回全部非空实例，服务组名使用请求节点组。
- 配置适配：使用上游`naming.DiscoSetting`替代 media 自有 Nacos SDK adapter，保留命名空间、账号、密码、日志目录、缓存目录、超时和 URL host 校验；URL host 在传给上游 gRPC client 前规范化为裸主机名，避免把`http://`前缀当作 gRPC 地址的一部分。Register 对上游 client 初始化期间的`client not connected`启动竞态做最多`10`次、每次`100ms`的有界重试，其他错误立即返回。
- 测试：更新 handler fake 以实现上游`naming.DiscoClient`契约，新增同一 TCP 连接复用 client 测试和两个实例 Lookup 响应测试；真实 Nacos 集成测试使用`qiniu`命名空间验证注册实例返回`healthy=true`，再执行 Lookup 全实例列表和 Deregister；常规测试不涉及手工保留测试。
- 插件与架构边界：修改限定在`media`插件 collection discovery 生产/测试代码、依赖文件和本 OpenSpec 记录，不修改`apps/lina-core`、`apps/lina-vben`或`hack`；`apps/lina-plugins/media/AGENTS.md`不存在，继续遵守仓库顶层规则。新增 factory 是上游 example 要求的生命周期适配，不引入额外跨模块契约。
- API、数据权限、数据库：不修改 HTTP API、DTO、权限、租户边界、SQL、DAO、DO、Entity、索引或报表写入路径，无数据权限和`N+1`影响。
- 缓存一致性：新增的是上游 Nacos naming client 的连接/订阅生命周期复用，不使用 media 业务 cache 保存 discovery 状态；既有宿主共享 cache、生命周期计数、缓存键、TTL、失效和集群策略未修改，无 media cache 一致性影响。
- i18n 与前端：不修改运行时用户文案、API 文档、插件清单、语言包、翻译缓存、前端页面或交互，无 i18n、前端和 E2E 影响。
- 开发工具跨平台：不修改`Makefile`、脚本、CI、`linactl`或工具入口；只更新 Go 依赖、Go 源码和测试，使用标准 Go build tag 与既有插件 lint 入口。
- DI 来源检查：未新增生产运行期依赖或服务构造函数参数；`discoClientFactory`由 collection runtime 创建并在 TCP server 生命周期内持有，底层 owner 是上游`naming.NewNacosDiscoverClient`，通过连接关闭和插件上下文关闭释放，不创建宿主独立服务图。
- 测试策略：后端纯 discovery 生命周期和协议行为使用单元测试；真实 Nacos 集成测试验证健康注册和全实例 Lookup；media 全量测试覆盖常规路径。手工保留测试仍由`manual_nacos` build tag隔离，不参与常规全量测试。
- 审查：已执行`lina-review`反馈级审查，范围为`collection_discovery.go`、新增`disco_cli_factory.go`、handler、相关测试、`go.mod/go.sum`和 OpenSpec 记录；确认 client owner、连接生命周期、列表 Lookup、错误处理和测试覆盖符合规则，未发现严重或警告问题。
- 规则加载：已读取`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/architecture.md`、`.agents/rules/plugin.md`、`.agents/rules/backend-go.md`、`.agents/rules/testing.md`、`.agents/rules/documentation.md`、`.agents/instructions/markdown-format.instructions.md`，并使用`lina-feedback`、`lina-review`、`goframe-v2`和`karpathy-guidelines`技能；API contract、数据权限、数据库、i18n、前端 UI 和开发工具长期入口确认无文件变更影响。
- 验证：`go list -m -versions github.com/dellinger2023/net-flux`和上游 tags确认最新正式版本为`v0.0.19`；`LINAPRO_TEST_NACOS=1`配合`http://10.157.225.139/`、`qiniu`、`nacos/nacos`的真实 Nacos 集成测试通过，并断言 Lookup 实例`healthy=true`；`GOROOT=/Users/wanna/sdk/go1.25.8 GOTOOLCHAIN=local GOWORK=off go test ./... -count=1 -parallel 1`通过；常规和`GOFLAGS='-tags=manual_nacos'`的`make lint dir=apps/lina-plugins/media plugins=1`均通过且为`0 issues`；`openspec validate add-media-data-collection-server --strict`、主仓和插件子仓`git diff --check`通过。


### FB-39 反馈修复记录

- 分诊：`openspec-existing`。反馈补全本活跃变更的采集读模型和节点总览契约，属于规范级功能遗漏；虽然此前任务全部完成，该变更尚未归档。本次先更新增量规范和设计，再实施。
- 根因：`NetworkMetric`中的抖动、丢包率和探测状态仅写日志，节点投影只保留目的地址到延迟的摘要；已有资源字段为包含子节点的汇总，不适合直接作为每个途经节点的直属实例资源；采集端将`CPU`百分比直接保存到声明为核数的`cpu_load`。
- 实现：扩展`GET /api/v1/media/dashboard/nodes/overview`，新增`network_links`、响应根节点的`network_links_limited`和每节点`instance_resources`。链路项保留来源实例、源地址、目的地址、原始延迟、抖动、丢包比例、状态码、可空可达状态和毫秒采样时间。直属资源复用已有实例聚合，不包含子树；采集端将`CPU`使用率转换为已用核数。不新增质量评分、质量分级、候选路径排序或综合负载百分比。
- 数据库与并发：当前活跃迭代的`003-add-media-data-collection-server.sql`追加`plugin_media_report_network_link`最新投影表，复合主键为节点、来源实例、源地址和目的地址。仅新增幂等`DDL`，无种子或演示数据，无自增主键和软删除。链路采用固定两次写入及采样时间条件更新，拒绝更旧样本；空目的地址不创建链路。可空布尔字段通过独立生成配置得到`*bool`，缺省状态显式更新为`SQL NULL`。
- 性能：链路按当前返回子树节点在数据库侧集合过滤，只新增一次投影查询，最多读取`10001`条以检测`10000`条截断；不逐节点、逐实例或逐链路回查。复合主键以前导`node_id`支撑过滤及稳定排序，资源继续使用既有分组查询，无新增`N+1`。
- 权限：沿用既有认证链、`media:management:query`和平台基础设施监控边界；监控样本无人员或部门归属，按组织数据权限隔离的业务实体规则不适用于该投影，例外及拒绝边界已在设计中记录。接口权限元数据测试和现有路由认证测试通过；新增 SQL 只查询响应子树节点，测试覆盖子树外样本不返回且不占用截断额度。
- 缓存与依赖：新增最新投影直接写共享数据库，不新增业务缓存、失效机制、运行期服务依赖或构造参数；既有共享生命周期计数缓存与启动服务图不变，无新增`DI`来源问题。
- 插件与跨平台：所有生产修改限定在`media`插件内部；`media/AGENTS.md`不存在。生成配置和`make dao`、`make ctrl`沿用仓库跨平台入口，无新增脚本；本轮在`macOS`和本地`PostgreSQL`验证，未额外执行其他操作系统测试。主框架插件能力无变化，宿主插件`README`无需更新。独立`media`仓库未开发或同步。
- 文档与`i18n`：同步更新中英文`README`，保留用户预存白名单说明改动。`media/plugin.yaml`未启用`i18n`，按`.agents/rules/i18n.md`的单语言插件例外使用中文接口源文案，无需插件翻译资源。`make i18n.check`已执行但全仓失败，报告既有`1620`项硬编码文案及消息和前端翻译缺失，涉及`cms`等本轮未修改范围；本轮不宣称全仓国际化门禁通过，也不扩展修复无关模块。
- 测试：新增`6`个测试函数并扩展既有断言，覆盖原始失败值、零丢包和零抖动、`true/false/null`、`true`更新为`NULL`、秒转毫秒、四维业务键、空目的端、迟到报告、`12`个并发写入、子树过滤、直属资源与子树资源区分、恰好和超过`10000`条边界，以及接口权限和序列化。后台接口行为通过真实`PostgreSQL`持久化、`SQLite`查询和控制器/路由自动化联合验证；本次没有前端页面或调用方变化，未运行页面`E2E`或截图。
- 生成与验证：本地新表`DDL`连续执行两次通过；`make dao dir=apps/lina-plugins/media/backend`和`make ctrl dir=apps/lina-plugins/media/backend`通过；`LINAPRO_TEST_POSTGRES=1 go test ./backend/internal/service/collection -count=1`通过；插件目录`go test ./... -count=1`全部通过；`make lint dir=apps/lina-plugins/media`通过，`0 issues`；`openspec validate add-media-data-collection-server --strict`及主仓、插件子仓`git diff --check`通过。初次误用`make dao p=media`产生的宿主生成差异已全部恢复，不包含在交付中。
- 审查：按`lina-review`完成反馈级主审及独立子代理审查，涵盖已跟踪改动和新增文件、生成层、SQL、数据权限和测试有效性；文档影响描述警告已修正，本轮代码未发现阻塞问题。全仓国际化检查失败作为范围外风险记录，不将其报告为通过。已读取`AGENTS.md`及`openspec`、`documentation`、`architecture`、`plugin`、`backend-go`、`api-contract`、`database`、`data-permission`、`cache-consistency`、`testing`、`i18n`、`dev-tooling`规则和`markdown-format.instructions.md`；使用`lina-feedback`、`lina-review`、`goframe-v2`、`karpathy-guidelines`和`use-modern-go`技能。
- 交付限制：部署环境需执行新增建表`DDL`，本地已执行且验证；新链路详情需等待后续网络上报，既有延迟摘要无法反推丢包率或抖动。既有`CPU`样本由后续机器上报刷新。保留旧`node_latency_map`接收顺序更新语义，新增`network_links`按采样时间防止旧值覆盖。未部署、提交或推送。


### FB-39 推送前复核

- 范围：只提交节点原始链路指标、直属实例资源、生成模型、迁移及双语说明，主仓提交本变更的规范记录和子模块指针；预存白名单响应修复、`go.work/go.work.sum`、其他 OpenSpec 变更和`lina-10-12.log`保留在工作区。
- 验证：`GOWORK=off go test ./... -count=1`通过，确认交付不依赖本地`go.work`调整；本地`PostgreSQL`定向执行`TestReportRuntimePersistsMetrics`、`TestReportRuntimePreservesLatestRawNetworkLinks`和`TestReportRuntimeRejectsConcurrentStaleNetworkLinks`通过；`make lint dir=apps/lina-plugins/media`通过，`0 issues`。
- 复核：重新读取本轮命中规则，检查暂存内容和新增文件；确认可空状态、原始失败值、子树过滤、截断、直属资源、CPU单位、并发陈旧保护和部署迁移均符合契约。既有全仓`i18n.check`失败仍按前述单语言插件边界记录，无新增翻译资源要求。复用本轮后端行为自动化覆盖，未新增前端或页面验证范围。
- Git 范围：根据用户本次“review后push”授权，在两仓库当前的`origin/feat/wanna`提交和普通推送，不执行公共仓库合并或部署。


### FB-40 推送复核修复记录

- 根因：独立复核追踪到`net-flux v0.0.19`真实发送实现，`MachineMetric.node_id`使用`nodeInfo.ID`实例标识，而`NetworkMetric.node_id`使用业务节点号；直接按两个`node_id`关联会使链路节点资源为零，旧测试人为使用相同节点标识而未覆盖该场景。
- 修复：仅调整`media`看板读适配层。在既有实例聚合 SQL 中联结每个来源机器的唯一最新业务节点；明确实例节点优先，空节点或实例标识占位时通过`machine_id = instance_id`推断。归属在全局节点范围选最新时间，再展示子树；同刻跨节点冲突、缺少关联或仍为实例标识的归属保持未知。多源和多目的端在联结前合并，不重复累加资源；资源不依赖已截断的链路详情。
- 性能和数据库：保持原资源聚合的一次集合 SQL，无逐实例查询；新增`idx_plugin_media_network_link_machine_report`索引覆盖机器、采样时间、节点查询。本地只执行新增索引 DDL，连续执行两次通过，未重置数据库。索引追加到同一活跃迭代 SQL，无新数据字段或生成层变化。
- 测试：新增真实协议的`SQLite`回归，覆盖先网络后机器、先机器后网络、多目的端、明确节点优先、迁出当前子树、同刻冲突及随后恢复唯一归属；新增真实`PostgreSQL`回归，校验实例数量、CPU、内存、磁盘和网络总量。测试发现 GoFrame 的`g.DB`默认驱动缓存不随`gdb.SetConfig`恢复，现将该 PG 用例放入标准`os/exec`启动的独立测试子进程，父进程不切换数据库驱动，确保可与已有 SQLite 测试一起执行；所有测试数据按随机唯一标识清理，验证残留为零。
- 影响：无新增运行期依赖、缓存、跨插件契约、权限入口或可见范围；沿用平台监控边界和认证。无前端或页面变化，采用后端接口等价自动化验证；插件未启用`i18n`，双语说明同步更新，无新增翻译资源。测试子进程使用标准 Go API，无 shell 或平台专属工具，无跨平台生产入口变化。
- 验证：`LINAPRO_TEST_POSTGRES=1 GOWORK=off go test ./... -count=1`全量通过，包括 collection 写入并发与 NULL 验证、media 真实 PG 关联及全部 SQLite 回归；`make lint dir=apps/lina-plugins/media`通过，`0 issues`；OpenSpec 严格校验和两仓差异检查通过。全仓既有 i18n 问题继续按 FB-39 范围说明保留，不宣称全仓门禁全绿。
- 审查：按`lina-review`完成主审与独立复核，此前真实数据关联阻塞已解除，未发现新的本轮代码阻塞问题；本轮重新读取的规则域与 FB-39 相同。保持部署前执行 DDL 的要求，同版本替换二进制不会自动完成插件迁移。
- 提交策略：首个子模块实现提交后发现上述问题，通过后续修复提交补齐，不重写已推送历史；主仓最终指针仅指向修复并验证后的子模块提交。预存白名单、工作区配置和其他变更继续保留。


### FB-41 完整最新上报存储

- 分诊：`openspec-existing`。本轮针对既有采集实现的字段遗漏，用户明确选择“保留每个对象的完整最新上报”，不要求逐次历史。逐字段审计见同目录`report-storage-audit.md`，覆盖`26 + 10 + 28 + 25`个顶层字段及`LinkHop`全部`6`个字段。
- 根因：已有代码只保存看板投影，磁盘容量/耗时/次数、服务器类型、音视频编码、完整`extra`与跳点元数据遗漏；归一化改变部分原始 ID、枚举、宽高、时间、数值精度和默认值。删除事件跳过 upsert，首次只收到删除时没有资源行。
- 实现：四张既有表增加可空`report_payload JSON`；流/会话增加`report_subcommand INTEGER`。从原始`metric`通过官方`protojson`序列化，使用原字段名、全部默认值和数值枚举；64位整数用字符串保存精度，完整保留`extra`和嵌套跳点。普通上报将 JSON 与原有投影同次写入；网络仍使用同一时间条件拒绝旧报文。
- 删除路径：先保存完整报文和子命令，再继续原有计数和关闭处理；已有记录只更新原始元数据，稀疏删除不会清空展示字段；首次删除插入时即为关闭状态。计数失败时返回错误但保留完整删除报文，恢复重试仍可完成；未将整个生命周期操作改造成原子事务，也不作该承诺。
- 查询与性能：实例、流、会话列表及流摘要显式排除原始 JSON，其他拓扑与聚合保留既有明确投影，HTTP DTO不新增原始内容。正常上报没有新增数据库调用，DELETE增加一次固定写入，不随字段或跳点数量增加查询，无新增`N+1`；新列不参与过滤或排序，不额外加索引。
- 数据库：当前活跃迭代`003-add-media-data-collection-server.sql`追加幂等增列和注释，仅 DDL，无种子或演示数据，不改变主键、软删除或清理语义。新增 DDL 在本地`127.0.0.1`的`PostgreSQL`执行两遍，通过幂等验证；`make dao dir=apps/lina-plugins/media/backend`生成4表对应的12个 DAO/DO/Entity文件，未手改生成代码，无无关生成差异。旧行保持 NULL，原始数据无法从有损投影回填。
- 运行边界：保留现有业务键和无效报文跳过策略；空目的网络报文不创建完整链路记录，网络较旧采样不覆盖payload。机器/流/会话按到达顺序替换，链路共享相同现有主键时共用一份最新内容。关闭流/会话仍按30天策略清理。JSON保存当前协议已定义字段，不保存未知 wire 字段或逐字节 TCP 帧。
- 权限与安全边界：原始内容只进入插件自有表，不新增公开原始查询入口或改变租户、平台数据权限；看板DTO不暴露任意`extra`，不会将潜在敏感扩展直接返回前端。
- 缓存、DI、宿主与跨平台：复用原有宿主共享计数缓存，不改缓存键、TTL、失效或跨实例机制，不新增运行期依赖与构造参数。修改全部限定media插件和OpenSpec，插件根无AGENTS.md，不修改独立media仓库或宿主。未新增开发脚本，DAO沿用既有跨平台生成入口；本轮运行环境为macOS，本地PostgreSQL验证。
- 文档与i18n：同步双语README和审计记录；不修改API文案、运行时文案、语言包或插件清单，media未启用i18n，无新增翻译资源要求。前端页面和交互不变，后端行为通过真实数据库与协议往返的等价自动化验证，不新增页面E2E或截图。
- 测试：新增3个测试函数，使用每字段有值的独立fixture、protobuf descriptor字段覆盖检查和`proto.Equal`验证四类型及嵌套JSON往返；覆盖大整数、默认零值、未知枚举、完整extra、非法UTF-8序列化拒绝、正常更新、旧网络报文保护、已有与首次DELETE、计数失败及重试、机器payload不被计数覆盖、未提供payload时保持SQL NULL。PostgreSQL数据使用随机唯一标识并清理。
- 验证：`LINAPRO_TEST_POSTGRES=1 GOWORK=off go test ./... -count=1`全部通过；补充异常断言后再次运行`LINAPRO_TEST_POSTGRES=1 GOWORK=off go test ./backend/internal/service/collection -count=1`通过；`GOWORK=off make lint dir=apps/lina-plugins/media`通过，0 issues；OpenSpec严格校验和两仓diff检查通过。
- 审查：按`lina-review`完成主审与独立审查，覆盖25个本轮文件及未跟踪源码/审计文件，未发现本轮代码阻塞问题。已读取AGENTS及openspec、backend-go、architecture、plugin、database、data-permission、cache-consistency、testing、documentation、i18n、dev-tooling和Markdown格式规则，并应用goframe-v2、Modern Go及最小改动原则。原有go.work/go.work.sum和日志不属于本轮。
- 独立发现：上游磁盘读写字节为累计计数，既有看板换算后误标速率。本轮完整payload保留原始计数及耗时、次数，显示口径另行处理，审计报告已明确，不据此宣称展示指标完全正确。
- 交付：本地代码和DDL已验证，尚未提交、推送、同步独立media仓库或修改部署数据库。线上需先执行新增列DDL，再部署后等待新报文。


### FB-41 二次完整性复核

- 用户要求再次确认全部字段是否补齐。重新审查四类入口、原始序列化、全部 upsert 更新列、DELETE、计数更新及失败返回，限定保证范围为每个符合既有业务投影约束的对象的完整最新有效报文，不承诺无效键归档或任意协议值绕过数据库字段约束。
- 新发现并修复：合法`extra`字符串可包含`U+0000`，`protojson`输出合法的转义形式，但 PostgreSQL `JSONB`拒绝，导致已有测试未覆盖的有效原始内容无法入库。已通过只读类型转换 SQL 和真实写入回归复现`unsupported Unicode escape sequence`，不能只依赖普通字段 fixture 判定完整性。
- 修复：四表`report_payload`改用`JSON`，在同一迭代 SQL 中增加已有原始列转`JSON`的幂等迁移；仅本地执行原始内容增列及类型转换片段两次，确认四列均为`json`、可空。重新运行`make dao dir=apps/lina-plugins/media/backend`，生成层类型仍为字符串，未产生其他表的差异。双语 README、设计、规范和审计报告同步更新。
- 验证：四类型完整fixture的`extra`均增加`before\x00after`，保持 descriptor 全字段检查和`proto.Equal`数据库往返；修改前真实 PostgreSQL 测试失败，类型修复后`LINAPRO_TEST_POSTGRES=1 GOWORK=off go test ./... -count=1`全部通过；`GOWORK=off make lint dir=apps/lina-plugins/media`通过，0 issues；OpenSpec严格校验和两仓差异检查通过。
- 复核结论：当前`v0.0.21`的89个顶层字段及6个跳点字段均由原始JSON保留，没有再发现本轮保证范围内成功处理却遗漏已定义字段的路径。原始字段不是全部拆成独立 SQL 列，已有看板投影仍保留其显示口径；未知 wire字段、被业务约束拒绝的报文、历史版本和旧数据回填不在保证范围。
- 治理与范围：再次读取 AGENTS 及后端、数据库、插件、架构、权限、缓存、测试、文档、OpenSpec、i18n和开发工具规则；无新缓存、DI、权限入口或前端变化。独立审查同样确认类型边界修复后限定范围完整。未提交、推送、同步或部署；线上应用前需使用最终JSON类型的迁移。


### FB-42 第三方查询返回完整上报

- 分诊：`openspec-existing`。用户明确要求“刚才加的数据都在接口里面返回下”，在本活跃变更中补齐既有四类展示查询，不新增页面或另一套查询入口。
- 根因：`FB-41`仅补齐原始报文存储；实例、流和会话查询主动排除了原始列，网络链路投影及四类响应`DTO`也未声明原始字段，因此第三方无法读取新增存储内容。
- 实现：实例列表、节点总览的各级`network_links`、流列表和会话的`active_session_list`统一返回`report_payload`对象；流和会话同时返回`report_subcommand`。未记录原始报文的记录返回`null`，未记录命令返回`0`。完整保留原字段名、嵌套跳点、`extra`、零值、转义字符、枚举值、原始时间单位和`protobuf`64位整数字符串；解码使用`json.Number`避免数字精度损失，非法或非对象内容返回既有查询错误。
- 契约与文档：更新四个`DTO`、控制器映射、插件自有`OpenAPI`及双语说明。实际生成器将`map[string]any`的属性错误限制为对象，本轮在插件文档构建中修正为允许任意属性值且可空的对象，测试直接验证生成文档。已执行`make ctrl dir=apps/lina-plugins/media/backend`，生成入口与路由定义无额外差异。
- 关联修复：查询回归发现流没有协议摘要时，同一会话协议会从计数和会话行重复加入分组。补齐计数阶段的已见集合，回归断言每个协议仅返回一组。
- 性能和权限：原始字段随对应记录同批读取，不增加数据库访问次数或逐项查询；保留既有节点、租户、设备等筛选、`10000`条上限、链路截断标志及鉴权、管理权限声明。会话查询仍必填`streamId`且仅返回活跃会话，不扩大为关闭事件或历史查询。完整报文增加单条响应大小，第三方应使用现有筛选限制返回范围。
- 边界与无影响判断：所有实现位于`media`插件，插件根目录没有`AGENTS.md`；没有修改宿主、前端、权限入口、采集存储逻辑或独立`media`仓库。没有新增`SQL`、迁移、索引、缓存、运行期依赖或依赖注入变更，也没有开发工具和跨平台入口变化。部署仍需先前完整上报存储迁移，本次返回层不需要额外更新数据库。
- 测试与质量审查：本轮为后端接口行为变化，执行等价自动化质量审查；数据库查询、控制器映射、响应序列化和实际`OpenAPI`生成覆盖四类成功路径、空值、嵌套字段、命令、数字精度、非法载荷及原查询范围。没有页面变化，不新增浏览器截图或前端`E2E`资产。新增测试各自构造并清理数据，不依赖其他用例。
- 验证：`LINAPRO_TEST_POSTGRES=1 GOWORK=off go test ./... -count=1 -timeout=60s`全量通过，覆盖真实`PostgreSQL`存储与关联、查询及接口回归；`GOWORK=off make lint dir=apps/lina-plugins/media`通过，`0 issues`。首次未设置短超时的全量运行在既有`backend`路由测试中等待，终止取栈后，单独`backend`测试及上述有界全量复跑均通过，未确认可复现根因，不将首次运行记为通过。
- `i18n`：`media`未启用多语言，新增中文`dc`符合单语言插件规则，无需新增翻译资源。已运行全仓`make i18n.check`，因既有治理问题及扫描器仍收集单语言插件中文而失败，共`1626`项，其中本轮新增的`6`项均为上述允许的中文字段说明；不宣称全仓检查通过。日志为`temp/20261026/112700-media-payload-i18n-check.log`。
- 规则加载：本轮已读取`AGENTS.md`、`.agents/rules/openspec.md`、`documentation.md`、`architecture.md`、`plugin.md`、`data-permission.md`、`api-contract.md`、`backend-go.md`、`testing.md`、`i18n.md`及`.agents/instructions/markdown-format.instructions.md`，使用`lina-feedback`、`lina-review`、`goframe-v2`、`karpathy-guidelines`和`modern-go-guidelines`技能。
- 审查：按`lina-review`完成主审与独立只读复核，范围包含未跟踪的文档构建测试及代码生成入口；未发现本轮阻塞问题。`openspec validate add-media-data-collection-server --strict`和两仓`git diff --check`通过。保留全仓`i18n`检查失败及首次测试等待的如实记录；本轮未提交、推送或同步独立仓库。


### FB-43 OpenAPI 原始报文示例修正

- 分诊：`openspec-existing`。本轮闭环`FB-42`复核发现的文档生成遗漏，补充既有验收，不改变四类接口业务契约。
- 根因：`GoFrame`对映射字段保留`eg`字符串，此前仅修正`type`、`nullable`和`additionalProperties`，四类`report_payload`示例仍被二次编码为字符串。
- 修复：在插件自有`applyReportPayloadSchemas`中，将字符串示例转换为`json.RawMessage`后嵌入生成文档，保留原始值与数字表达。示例来自固定`DTO`标签，序列化与类型由回归测试验证，不引入重复解码、额外错误接口或宿主修改。
- 验证：先补四类实际生成示例的对象类型和源`eg`内容断言，修复前四个子场景全部复现失败；修复后`GOWORK=off go test ./backend/internal/apidoc ./backend -count=1 -timeout=60s`通过，覆盖文档构建及路由。`GOWORK=off make lint dir=apps/lina-plugins/media`通过，`0 issues`；两仓`git diff --check`通过。
- 影响：本轮只修改插件文档构建和对应测试，无数据库、数据权限、查询装配、缓存、运行期依赖、依赖注入、前端或跨平台工具变化。只修正`eg`的输出编码，不修改文案；`media`仍未启用多语言，且`eg`不参与翻译治理，无新增翻译资源影响。此前全仓`i18n`检查失败仍按`FB-42`记录，不宣称全仓门禁通过。本轮采用实际文档生成的等价自动化验证，无页面截图需求。
- 规则与审查：本轮重新读取`AGENTS.md`、`.agents/rules/openspec.md`、`plugin.md`、`api-contract.md`、`backend-go.md`、`testing.md`、`i18n.md`、`documentation.md`及 Markdown 格式指引；插件根目录无`AGENTS.md`。执行`lina-feedback`与`lina-review`，使用`goframe-v2`和现代`Go`指引；主审与独立复核确认原警告已解决，未发现新增阻塞项。
- 交付：`openspec validate add-media-data-collection-server --strict`通过；保留既有未提交工作区，本轮未提交、推送或同步独立仓库。


### FB-44 全局活跃数量汇总

- 分诊：`openspec-existing`。用户要求增加活跃流总量和活跃会话总量接口，补全本活跃采集变更的第三方查询能力；先更新增量规范、设计和任务，再完成实现。
- 实现：新增`GET /api/v1/media/dashboard/summary`，无参数，返回`data.active_stream_count`和`data.active_session_count`，均为`int64`。活跃沿用`close_time IS NULL`，空表返回`0`，不根据展示状态或心跳年龄额外过滤。
- 性能与存储：使用同一条`SELECT`中的两个独立`COUNT(*)`子查询，固定一次数据查询，避免联结放大和两条语句的快照偏移；不加载实体列表、原始报文或实例累计计数，不受`10000`条限制。复用两表的关闭时间索引及活跃部分索引，精确计数扫描成本随匹配记录增长，不宣称常数时间。无新增迁移、索引、生成模型或写入；数据库幂等性和数据分类无新增影响。
- 权限：接口声明`media:management:query`，沿用`platform_only`媒体管理双鉴权。宿主路径执行现有权限中间件；铁塔路径仍按既有认证回退处理，不声称铁塔调用经过宿主角色授权。用户要求的是全平台监控总量，指标无人员／部门归属，不按租户或组织筛选；拒绝策略由现有双鉴权提供，不新增匿名或内部密钥入口。该边界已同步设计、规范及接口说明。
- 生成：执行`make ctrl dir=apps/lina-plugins/media/backend`生成接口签名及控制器骨架；生成器附加的重复空构造已移除，既有`media_new.go`恢复无差异，控制器复用原有注入服务。
- 测试：服务层验证空表、关闭排除、状态不替代关闭标志、多节点／租户、孤立会话、错误不伪装为零及每类超过`10000`条；实际`HTTP`测试验证无参调用、整数形状、关闭后计数变化和双鉴权拒绝。`OpenAPI`测试验证新路径、两个`integer/int64`字段和权限声明。属于后端可观察接口行为，执行等价自动化质量审查，无页面或截图变更。
- 测试隔离修复：首次全量运行暴露新路由测试的临时`SQLite`默认驱动缓存污染后续测试，已将新测试放入独立测试子进程，未修改旧测试排序或共享 helper；两个新服务测试各自注册`Cleanup`，恢复缺表并清空批量数据。所有本轮副作用由当前测试闭环清理。
- 验证：最终`LINAPRO_TEST_POSTGRES=1 GOWORK=off go test ./... -count=1 -timeout=60s`通过，包含既有真实`PostgreSQL`回归和本轮`SQLite`／`HTTP`／`OpenAPI`用例；`GOWORK=off make lint dir=apps/lina-plugins/media`通过，`0 issues`。`openspec validate add-media-data-collection-server --strict`及两仓差异检查通过。
- `i18n`：`media`未启用多语言，新接口中文文档符合单语言插件规则，无新增翻译资源要求。全仓`make i18n.check`仍失败，`1629`项，其中本轮增加的`3`项为允许的中文元数据；日志为`temp/20261026/135102-media-summary-i18n.log`。不将既有全仓失败误报为通过。
- 其他影响：无宿主、前端、缓存一致性、运行期依赖、依赖注入或跨平台工具入口变化；测试隔离使用标准`os/exec`。独立`media`仓库未修改，预存`go.work`、`go.work.sum`和日志保留。本轮不自动提交、推送或同步。
- 审查：按`lina-review`完成主审与独立复核，审查包含`10`个插件变更文件、`4`个规范记录文件及无差异的生成入口；新测试清理警告已关闭，未发现剩余阻塞项。已读取`AGENTS.md`、`.agents/rules/openspec.md`、`architecture.md`、`data-permission.md`、`plugin.md`、`api-contract.md`、`backend-go.md`、`database.md`、`testing.md`、`i18n.md`、`documentation.md`及 Markdown 格式指引，插件根目录无本地`AGENTS.md`；使用`lina-feedback`、`lina-review`、`goframe-v2`、`karpathy-guidelines`与现代`Go`指引。


### FB-45 看板 OpenAPI 完整性

- 分诊：`openspec-existing`。用户指出流状态文档不准确并要求补全，属于当前第三方查询契约的文档遗漏。根因是`DTO`沿用展示示例`playing/paused`和协议`FLV/DASH`，没有跟随采集映射，也没有明确列表范围与计数边界。
- 修复：流筛选、流响应及拓扑流状态列明实际`8`种状态并使用`running`示例；会话协议筛选及响应列明`9`种标准协议。补齐可选参数默认值、精确匹配、来源参数组合和请求值回显、列表排序、共享`10000`条上限、完整分组计数与明细长度区别、累计／当前计数和派生字段来源。文档说明统一`code/message/data`包装位置，未增加参数枚举校验或修改流／会话列表过滤行为。
- 同步：更新双语说明和增量规范；执行`make ctrl dir=apps/lina-plugins/media/backend`，生成接口入口无额外差异。使用实际`Build`生成文档的回归，验证状态／协议值、示例、必填`streamId`和计数说明进入最终文档。

### FB-46 汇总限定运行中数据

- 用户追加明确要求`summary`只看`running`，覆盖此前`FB-44`按所有未关闭记录计数的口径。流数现在要求`status=running AND close_time IS NULL`；会话本身无状态字段，要求自身未关闭，且通过`stream_id`关联上述流。非运行中流、关闭流和孤立会话不计入。
- 实现：保留一条`SELECT`，流计数增加参数化状态过滤，会话计数使用`EXISTS`过滤对应流，避免联结放大和逐行补查。响应字段、平台全局范围、无分页和管理鉴权保持原约定；仅收窄数据范围，既有平台监控数据权限例外和双鉴权拒绝边界继续适用。
- 性能：复用活跃记录部分索引、关闭时间索引、流主键和会话流索引。对本地`PostgreSQL`执行`EXPLAIN`确认语法有效；小数据优化器选择顺序扫描，在事务内临时关闭顺序扫描后确认活跃索引及流主键可用，随后回滚会话设置。没有修改数据库结构或新增索引，不引入缓存。
- 测试：服务层覆盖所有非运行状态、已关闭但状态为`running`、关闭会话、孤立会话、空结果、超过`10000`条和缺表错误；实际`HTTP`覆盖`error`恢复`running`后总数变化及关闭流后关联会话退出统计。新口径与文档同步，不根据心跳超时推断状态。

### FB-45／FB-46 验证与审查

- 验证：`LINAPRO_TEST_POSTGRES=1 GOWORK=off go test ./... -count=1 -timeout=60s`通过；最后文案补正后再次运行`GOWORK=off go test ./backend/internal/apidoc ./backend/api/media/v1 -count=1 -timeout=60s`通过。`GOWORK=off make lint dir=apps/lina-plugins/media`通过，`0 issues`；`OpenSpec`严格校验与两仓差异检查通过。
- `i18n`：插件未启用多语言，沿用其源文案语言，无需补翻译。已运行全仓`make i18n.check`，仍因既有治理问题报`1629`项，与修改前相同；日志`temp/20261026/142242-media-running-summary-i18n.log`，不宣称全仓治理通过。
- 影响与规则：本轮读取`AGENTS.md`和`.agents/rules/openspec.md`、`architecture.md`、`data-permission.md`、`plugin.md`、`api-contract.md`、`backend-go.md`、`database.md`、`testing.md`、`i18n.md`、`documentation.md`及 Markdown 格式指引。插件根目录无本地`AGENTS.md`。无新增迁移、数据写入、运行期依赖、依赖注入、缓存一致性、宿主、前端或跨平台工具变化，生产实现先在`linapro`闭环。执行后端等价自动化质量审查，保留测试独立进程和清理；无页面截图需求。
- 交付：本轮未提交、推送或同步独立`media`仓库；预存工作区配置及日志保留。
- 审查结论：已按`lina-review`完成主审及独立只读复核，覆盖全部当前插件与规范改动、测试及无差异的生成入口；确认文档与最终统计口径一致，未发现剩余阻塞问题。


### FB-47 网络速率以 Kbps 输出

- 分诊：`openspec-existing`。用户要求两个展示接口输出`Kbps`，对应节点总览和实例列表。根因是原展示字段直接返回存储的`KiB/s`数值，文档标为`KB/S`；这与客户端协议`bps`及用户要求的千比特单位不同。
- 实现：在插件控制器`DTO`映射层统一执行`KiB/s × 1024 × 8 ÷ 1000`，覆盖根节点、全部子节点、各级`instance_resources`以及实例列表的六处字段赋值。保留浮点数、不取整、不原地修改服务输出，聚合完成后每个展示字段只换算一次。
- 契约：`network_in/network_out`输出单位为`Kbps`，明确`1000 bps = 1 Kbps`。六个字段的`dc/eg`及双语说明已同步；例如投影`1000 KiB/s`返回`8192 Kbps`。原始`report_payload`仍保留`bps`，已有`NUMERIC(14,2)`存储舍入保持不变，不声称输出能够恢复丢失的精度。磁盘速率没有改变。
- 边界：没有修改采集、存储、`DAO/DO/Entity`、查询、索引、缓存、运行期依赖、依赖注入、权限、宿主或前端。沿用现有平台监控和双鉴权边界，纯输出换算不改变数据库访问次数，无`N+1`或迁移要求。实例采集时间／更新时间仅按用户询问说明现状，本轮未添加时间字段。
- 验证：控制器回归覆盖零、小数、固定物理量、已有两位小数存储、根／子节点和两级直属资源，并断言原始报文、磁盘数值和服务对象不被改变；实际`OpenAPI Build`验证三个资源类型的六个字段具有`Kbps`及`1000`进制说明。两个实际控制器入口均使用这些映射，作为本轮接口行为的等价自动化验证，无页面截图需求。
- 门禁：`make ctrl dir=apps/lina-plugins/media/backend`完成，两个生成接口入口无额外差异；`LINAPRO_TEST_POSTGRES=1 GOWORK=off go test ./... -count=1 -timeout=60s`全量通过，覆盖既有真实`PostgreSQL`回归。`GOWORK=off make lint dir=apps/lina-plugins/media`通过，`0 issues`；`openspec validate add-media-data-collection-server --strict`与两仓差异检查通过。
- `i18n`：插件未启用多语言，中文接口元数据按单语言规则维护。全仓`make i18n.check`仍失败，`1629`项，与本轮修改前相同；日志`temp/20261026/155134-media-kbps-i18n.log`。未新增翻译资源，保留既有全仓治理失败记录。
- 审查：使用`lina-feedback`、`lina-review`、`goframe-v2`、`karpathy-guidelines`和现代`Go`指引，主审与独立只读复核未发现阻断问题。本轮重新读取`AGENTS.md`及`.agents/rules/openspec.md`、`architecture.md`、`data-permission.md`、`plugin.md`、`api-contract.md`、`backend-go.md`、`testing.md`、`i18n.md`、`documentation.md`与 Markdown 格式指引；插件根目录无本地`AGENTS.md`。审查范围为本轮`7`个插件文件与`3`个规范文件，并检查无差异生成入口。无新增数据库、跨平台工具或缓存一致性规则影响。
- 交付：保留此前`FB-45/46`源仓未提交改动和用户工作区配置、日志。本轮未提交、推送或同步独立`media`仓库。


### FB-48 实例采集时间与更新时间

- 分诊：`openspec-existing`。用户要求把原始`report_payload.timestamp`以`report_time`返回到实例列表项外层，并增加`updated_at`，补全当前第三方展示契约。此前采集时间只在原始报文内，记录更新时间虽已落库但没有映射到响应。
- 实现：实例`DTO`和服务输出增加`*int64`字段`report_time/updated_at`，控制器直传，均返回`Unix`毫秒数字或显式`null`。采集时间从同批已解码原始报文提取，支持字符串整数及`json.Number`，正数小于`1000000000000`按秒乘`1000`，其余正整数按毫秒；缺失、非正或无法解析返回`null`。`updated_at`复用现有`formatTime`转换数据库更新时间。
- 数据来源：不以数据库实例`report_time`兜底采集时间，该列可能被流／会话计数事件更新。原始`report_payload.timestamp`保持类型与单位，不删除、不修改；既有列表排序继续使用原投影时间，不改为按新增外层采集时间排序。
- 契约与生成：新字段显式声明`nullable:true`与`format:int64`，修正生成器对`*int64`自动产生非标准格式的问题。同步双语说明；执行`make ctrl dir=apps/lina-plugins/media/backend`，两生成接口入口无额外差异。
- 测试：覆盖秒／毫秒、字符串／JSON整数、超过浮点安全范围的整数、缺失、零、负数、非法、溢出及可空更新时间。SQL联动刻意设置晚于机器采样的投影时间，确认不误用；比对实际数据库更新时间。控制器序列化及真实生成OpenAPI验证数字、空值、完整原始报文、可空`integer/int64`和毫秒说明。
- 验证：`LINAPRO_TEST_POSTGRES=1 GOWORK=off go test ./... -count=1 -timeout=60s`全量通过；`GOWORK=off make lint dir=apps/lina-plugins/media`通过，`0 issues`。`OpenSpec`严格校验及两仓差异检查通过。后端可观察接口行为采用上述等价自动化覆盖，无页面或截图需求。
- `i18n`：插件未启用多语言，新增中文字段说明符合单语言规则。全仓`make i18n.check`仍失败，`1631`项，比此前增加的`2`项均为允许的中文字段说明；日志`temp/20261026/160338-media-instance-times-i18n.log`，不宣称全仓治理通过。
- 无影响判断：从已读取实体和原始报文装配，无额外数据库查询或`N+1`，无新增存储、迁移、索引、数据写入、缓存、运行期依赖或依赖注入变动；查询范围、权限、宿主、前端和跨平台工具入口均无修改。
- 审查：按`lina-review`完成主审与独立复核，未发现阻断问题；范围包含本轮`9`个插件文件、`3`个规范记录文件及无差异生成入口。已读取`AGENTS.md`、`.agents/rules/openspec.md`、`architecture.md`、`data-permission.md`、`plugin.md`、`api-contract.md`、`backend-go.md`、`testing.md`、`i18n.md`、`documentation.md`及 Markdown 格式指引；插件根目录无本地`AGENTS.md`。使用`lina-feedback`、`lina-review`、`goframe-v2`及现代`Go`指引。
- 交付：本轮未提交、推送或同步独立仓库；此前未提交功能改动和用户工作区配置、日志均保留。


### FB-49 存活节点计数与接收时间索引

- 分诊：`openspec-existing`。用户要求`summary`新增存活节点数，并进一步明确最近`3min`上报即存活、不看`status`；随后希望使用`updated_at`改善语义和查询。最终新增`alive_node_count int64`，使用服务端节点记录更新时间的`[now-3min, now]`闭区间。
- 实现：保留同一条`SELECT`返回三个标量计数，节点不关联流、不求和`alive_nodes`，无`10000`条上限；参数化`updated_at >= ? AND updated_at <= ?`，一次捕获当前时间。移除本轮早期状态归一和秒／毫秒分支，不依赖客户端采样时间。空数据为`0`，过期、未来和缺失更新时间不计入；近期的`offline`或空状态节点正常计入。
- 上报时间修复：真实代码核对发现无目的地址网络上报的`Save/OnConflict`显式更新列遗漏`updated_at`，重复上报不会刷新接收时钟。本轮仅将该列加入`OnDuplicate`，值继续由`GoFrame`自动维护。有目的地址路径原有`Update`保持自动时间处理。
- 数据边界：当前只有网络上报更新节点投影，机器指标仅更新实例表，不能单独延长节点窗口；直接修改节点记录也会影响更新时间。接口、双语说明及设计已经注明服务端更新时间口径。节点总览既有`alive_nodes`状态算法不在本次变更范围，不能把两个接口的判定混为一谈。
- SQL与性能：当前活跃迭代`003-add-media-data-collection-server.sql`追加`CREATE INDEX IF NOT EXISTS idx_media_report_node_updated_at ON media_report_node(updated_at DESC)`，仅索引DDL，无新字段、种子、演示数据或生成模型变更。已对本地PostgreSQL连续执行两次确认幂等；诊断性关闭顺序扫描后`EXPLAIN`确认该谓词可使用`Index Only Scan`。这是索引可用性验证，没有宣称线上延迟或完成负载压测。部署时需执行新增索引，生产数据库未代执行。
- 测试：固定时钟服务用例覆盖三分钟前及当前时刻、边界外`1`微秒、状态独立、旧／未来客户端采样时间、真实`0`节点、干扰累计值和超过`10000`节点；测试构造明确通过SQL设置历史更新时间，生产不手工维护时间字段。实际HTTP及OpenAPI验证三字段数值、空值、拒绝访问与文档口径。新增真实PostgreSQL回归先复现无目的地址重复上报时间不刷新，再确认两条网络写入路径均刷新服务端时间，客户端旧采样值保持不变；测试记录全部清理。
- 验证：最终`LINAPRO_TEST_POSTGRES=1 GOWORK=off go test ./... -count=1 -timeout=60s`通过；`GOWORK=off make lint dir=apps/lina-plugins/media`通过，`0 issues`；OpenSpec严格校验和两仓差异检查通过。全仓`make i18n.check`仍失败，`1632`项，其中本轮增加的一项为单语言插件允许的中文字段说明；日志`temp/20261026/163911-media-receive-time-i18n.log`，不宣称全仓门禁通过。
- 审查：使用`lina-feedback`和`lina-review`，主审及独立复核确认当前聚合、时间刷新、索引与测试无阻塞问题。已重新读取`AGENTS.md`和`.agents/rules/openspec.md`、`architecture.md`、`data-permission.md`、`plugin.md`、`api-contract.md`、`backend-go.md`、`database.md`、`testing.md`、`i18n.md`、`documentation.md`及 Markdown 格式指引，使用`goframe-v2`与现代Go指引；插件无本地`AGENTS.md`。
- 其他影响：沿用媒体平台监控和双鉴权边界，无新增公开入口、租户范围或组织数据权限变化；无缓存、运行期依赖、依赖注入、宿主、前端或跨平台工具修改。后端接口变化有等价自动化覆盖，无页面截图需求。本轮未提交、推送或同步独立仓库，保留此前源仓改动及用户工作区配置和日志。
