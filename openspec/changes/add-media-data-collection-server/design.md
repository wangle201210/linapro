# 设计说明

## 宿主边界

采集 server 属于`media`源码插件自有运行时能力，不修改`apps/lina-core`的核心宿主契约。插件通过现有`system.started` hook 接入宿主启动生命周期，并通过 hook payload 的`Services()`读取插件配置。

## 模块归属

当前将采集 server 作为`media`插件内置子组件实现，而不是拆成独立插件。原因是本次明确需求来源于`media`模块数据采集，当前唯一消费者是媒体节点、流和策略相关能力。代码实现应保持独立 service 边界，避免把采集协议处理散落在`plugin.go`或管理接口中。

当采集 server 后续需要被监控、告警、运维看板、设备资产或其他插件共同消费，或者需要完整承载`net-flux`的 config、event 和 control 横向协议能力时，应新建独立采集插件并通过公开契约向`media`提供数据。

## 启动和配置

新增插件配置段`collectionServer`：

| 配置项 | 默认值 | 说明 |
| ------ | ------ | ---- |
| `collectionServer.enabled` | `false` | 是否启动采集 TCP server |
| `collectionServer.addr` | `:1911` | TCP 监听地址 |
| `collectionServer.discovery.enabled` | `false` | 是否启用 Nacos discovery 命令处理 |
| `collectionServer.discovery.host` | `127.0.0.1` | Nacos 服务地址 |
| `collectionServer.discovery.port` | `8848` | Nacos 服务端口 |
| `collectionServer.discovery.namespace` | `public` | Nacos namespace |
| `collectionServer.discovery.logDir` | `./logs` | Nacos SDK 日志目录 |
| `collectionServer.discovery.cacheDir` | `./cache` | Nacos SDK 本地缓存目录 |
| `collectionServer.discovery.notLoadCacheAtStart` | `true` | Nacos SDK 创建 client 时是否跳过加载本地磁盘缓存 |
| `collectionServer.discovery.timeout` | `5000` | Nacos SDK 请求超时时间，单位毫秒 |
| `collectionServer.discovery.username` | `nacos` | Nacos 用户名 |
| `collectionServer.discovery.password` | `nacos` | Nacos 密码 |
| `collectionServer.discovery.node` | `1` | 默认节点 ID，用于本地配置兜底；discovery 注册、注销和查询时映射为 Nacos group |

采集 server 默认关闭是为了避免在未声明运维端口时改变部署暴露面。discovery 默认关闭是为了避免未配置 Nacos 时影响纯指标采集；启用后才创建 Nacos discovery client。插件在宿主 HTTP server 启动后异步启动 TCP server，并复用宿主传入的上下文取消链完成停止。

## 协议边界

采集 server 复用`github.com/dellinger2023/net-flux`发布的`network.NewTcpServer`和`gen`协议类型。为补齐数据看板所需业务键，本次已将上报协议字段提交并合并到上游`github.com/dellinger2023/net-flux`，`media`插件保留原 import 路径，并将 Go module 依赖指向包含这些字段的上游合并 commit。

实现范围限定为：

- `Ping`请求返回`Pong`。
- `MachineMetric`、`NetworkMetric`、`StreamMetric`和`SessionMetric`上报被接受、记录日志并写入`media`数据看板上报读模型；其中`MachineMetric`字段名沿用`net-flux`历史命名，但语义收敛为实例/容器指标。
- `STREAM_ADD`、`STREAM_DELETE`、`SESSION_ADD`和`SESSION_DELETE`作为流和会话生命周期事件处理，采集 server 使用宿主发布给`media`插件的共享`cachecap.Service`维护资源归属和实例实时计数，并写回`media_report_instance.live_streams`和`media_report_instance.sessions`。
- discovery 命令在`collectionServer.discovery.enabled=true`时对接 Nacos：`Instance`注册实例，`Deregister`注销实例，`Lookup`查询实例并返回`LookupAck`；discovery client 不应跨 TCP 命令复用 Nacos SDK 本地订阅缓存，且默认不加载 Nacos SDK 磁盘缓存，避免多 Pod 部署下注销后仍读取旧实例。
- discovery 未启用时返回明确错误，不影响系统心跳和数据上报命令。
- config、event 和 control 命令暂不提供业务动作，避免在没有 LinaPro 侧契约时引入额外状态变更。

## 上报读模型写入

本次沿用`apps/lina-plugins/media/数据看板.md`中反推的上报读模型，新增插件安装 SQL 与 DAO 生成配置。表归属仍在`media`源码插件内，不修改`lina-core`宿主数据模型。

扩展后的`net-flux`上报包字段与看板模型的映射范围如下：

| 上报包 | 写入表 | 映射策略 |
| ------ | ------ | -------- |
| `MachineMetric` | `media_report_instance` | `instance_id`作为实例业务键，缺省时回退使用`machine_id`承载的容器/实例标识；节点、状态、CPU、版本和启动时间写入实例最新投影；内存`MB`、磁盘和网络`KB/s`六项指标按上报原值保存，不重复单位转换；`live_streams`和`sessions`不由客户端直接上报。 |
| `NetworkMetric` | `media_report_node`、`plugin_media_report_network_link` | 优先使用`node_id`作为节点业务键，缺省时回退`machine_id`和`source_ip`；可选扩展吞吐量仅解析数值后原值保存，由客户端统一单位，`destination_ip -> rtt`在节点行级锁下合并进`node_latency_map`。另按节点、来源实例、源和目的地址保存完整原始链路指标。 |
| `StreamMetric` | `media_report_stream`、`media_report_instance` | `stream_id`作为流业务键；`tenant_id`、`device_id`、`node_id`、`instance_id`写入租户、设备、节点和实例维度，`instance_id`为空时回退使用`machine_id`承载的容器/实例标识；协议、状态、码率、分辨率、帧率、丢包率、延迟、会话摘要和流路径写入最新流投影；`STREAM_ADD`和`STREAM_DELETE`按`stream_id`幂等维护实例实时直播流数量。 |
| `SessionMetric` | `media_report_session`、`media_report_instance` | `session_id`作为会话业务键；流、租户、设备、客户端、播放协议、播放质量、`node_id`、实例和链路跳点写入会话最新投影，`instance_id`为空时回退使用`machine_id`承载的容器/实例标识；`SESSION_ADD`和`SESSION_DELETE`按`session_id`幂等维护实例实时会话数量。 |

`media_report_instance`和`media_report_session`不再伪造业务键，必须分别由`MachineMetric.instance_id`或`MachineMetric.machine_id`、`SessionMetric.session_id`驱动写入；缺少业务键的上报包会被忽略。采集端无法稳定获取宿主机器信息时，不应把容器 CPU、内存和磁盘指标伪装成`media_report_node`物理节点快照。

## 性能和数据边界

上报处理只按单个上报包执行固定次数写入，不在动态结果集中循环查询数据库或调用远程插件能力。`MachineMetric`按实例业务键执行一次最新投影 upsert；`NetworkMetric`按节点业务键在事务内锁定单行并合并当前延迟矩阵，避免同节点并发目的端覆盖；`StreamMetric`和`SessionMetric`分别按流和会话业务键执行一次最新投影 upsert。看板后续查询可直接读取投影表并依赖`node_id`、`report_time`、`source_type/source_id`、`tenant_id/status`和`tenant_id/device_id`等索引，避免实时回查配置表和`N+1`装配。discovery 查询只按单个服务名和节点 ID 访问 Nacos，不在 LinaPro 数据库中装配动态结果集。

实例实时流数和会话数的权威运行时来源是宿主共享 cache 中维护的生命周期事件状态：`stream_id -> instance_id`、`session_id -> instance_id`、`instance_id -> live_streams`和`instance_id -> sessions`。源码插件只能通过`cachecap.Service`访问该状态，不能直接依赖宿主私有 Redis、`kvcache`或 coordination 实现；在`cluster.enabled=true`时，该服务由宿主统一切换到 Redis/coordination KV 等共享后端，在单机模式下可继续使用宿主 SQL KV 后端。采集 server 在处理生命周期事件时使用对应流或会话投影行的数据库行级锁串行化同一资源的事件，再更新共享 cache 和实例投影，保证多 Pod 下重复`ADD`、重复`DELETE`和跨实例迁移不会重复增减。cache 是实时计数权威状态，`media_report_instance`中的`live_streams`和`sessions`是供数据看板读取的最新投影；若共享 cache 写入失败，当前上报处理返回错误并保留下一次生命周期事件重试恢复的路径。

## 影响评估

### FB-54 查询返回时计算会话播放时长

按用户最新要求，`play_duration`在会话列表和设备拓扑返回时计算为服务端当前时间减去`StartTime`的整秒数。在会话映射使用处调用`time.Now()`，直接交给包内时长函数，不沿查询、分组和树形装配逐层传递时间参数。缺失、非正或未来开始时间返回零，不回退到上报时长，也不按关闭时间冻结计算。

计算只发生在既有响应组装路径，不写库、不增加查询或缓存，不改变上报报文、存储的`play_duration`、流`duration`或链路总延迟。两类查询仍仅返回原范围内的未关闭会话，沿用既有权限、过滤和数量上限。复用一个包内纯函数，不增加运行期依赖、DI或宿主改动。公开`dc`继续仅说明会话播放时长及秒单位，不添加内部计算说明；无API结构、i18n资源或前端页面变化。

验证使用固定当前时间覆盖整秒截断、空／零／未来时间、时区等价、时间推进和存储对象不变，并覆盖实际列表及拓扑查询。流时长与客户端总延迟仍验证原值返回，避免扩大本轮例外范围。

### FB-53 客户端度量指标原值存取

客户端负责计算和统一度量单位，服务端仅保存、读取并传递结果。清除`CPU`百分比转已用核数、网络扩展吞吐量缩放以及看板`KB/s`转`Kbps`的输出函数；`cpu_load`直接使用`cpu_usage`，网络和磁盘字段保留`v0.0.25`的`KB/s`，内存保留`MB`。节点直属资源与子树继续按已有关系累加原始数值，`CPU`汇总表示客户端利用率数值之和，不伪称总体加权利用率。

可选网络扩展字段保留既有键提取顺序，只解析数字，不依据`throughput_bps`等键名推断或转换单位；自定义上报端负责与机器资源指标保持一致单位。标准采集端在`MachineMetric`中直接给出已计算的网络`KB/s`。

流`duration`、会话`play_duration`和`total_link_latency`按客户端度量值保存，包括零值。流时长及链路总延迟直接输出；会话播放时长按后续`FB-54`规则在列表和拓扑返回时计算。分辨率字符串组装、业务标识与状态映射、时间戳序列化仍负责协议结构适配；流与会话生命周期计数、协议分组和节点聚合保留已有业务职责，包括既有协议计数缺省行为，不将这些服务器统计伪称为单条客户端指标。

字段结构、路由、权限、数据库列和精度不变，无DDL、DAO生成或历史回填；老投影等待新上报刷新。没有新增查询、远程调用、缓存或依赖注入路径，不新增`N+1`。资源API只说明`CPU`百分比、`KB/s`速率及秒、毫秒时长等业务含义和单位，不在`dc`中描述原值传递、单位换算或内部汇总；处理细节保留在源码注释与技术说明。插件未启用多语言，按单语言维护元数据并同步双语说明。插件管理页面没有看板指标消费，无前端页面改动，使用写入回读、DTO映射和实际OpenAPI生成的等价自动化验证。


### FB-52 新版采集协议与原值存储

`media`依赖按用户最终要求升级到`net-flux v0.0.25`，保留`github.com/dellinger2023/net-flux`模块和导入路径，通过`replace`从`github.com/dellinger2026/net-flux v0.0.25`获取源码。相较原`v0.0.21`，新版机器协议将内存、磁盘速率、磁盘耗时及网络速率共八个字段改为`double`；客户端已将内存计算为`MB`、磁盘和网络差分计算为`KB/s`、磁盘耗时从毫秒转为秒。按用户最新要求，实例投影的两项内存、两项磁盘速率及两项网络吞吐量均使用上报数字直接赋值，不进行再次除以`1024`、`1000`或`8`等处理。

`v0.0.25`机器网络速率已是按`1024`字节进制计算的`KB/s`，不再采用`v0.0.24`的`bps`单位。后续`FB-53`统一为度量原值输出，网络保留`KB/s`，移除原`Kbps`输出换算。节点`NetworkMetric.extra`扩展吞吐量也仅解析数值并原值保存，由客户端保证单位一致。旧换算函数和孤立常量移除，示例采集工具改用新版上报单位，机器指标日志适配浮点类型，原始`report_payload`继续由新版协议完整序列化。

字段编号未变但机器协议的八个字段从`varint`变为`fixed64`，采集客户端和服务端必须按新版协议配套升级，不增加旧版解析或单位猜测。既有数据库数值列及其两位小数约束保留，不改写历史记录；旧投影在后续新上报时刷新。无新增数据库查询、迁移、缓存、跨模块契约、运行期服务实例或依赖注入路径。权限边界沿用原平台采集和看板授权，不新增数据可见范围。插件未启用多语言，本轮不新增用户界面、接口文案或翻译资源。工具仅调整Go结构体示例数字和单位，不新增操作系统依赖；通过插件测试及Windows、Linux交叉编译验证。

验证覆盖新版浮点字段的协议往返、六项原值入库、计数事件不覆盖资源、节点吞吐量原值兜底、节点／子节点／直属资源和实例接口的网络原值展示、完整原始报文。采用后端写入与实际DTO映射的等价自动化验证，无页面交互或截图变更。

### FB-50／FB-51 全量节点配置内部入口

按`FB-51`收窄后的需求，仅保留`GET /api/v1/nodes/all`，归入`内部接口`。移除`FB-50`新增的看板全量节点接口及其`DTO`、控制器方法和生成接口签名，`OpenAPI`不再发布`/api/v1/media/dashboard/nodes/all`。内部入口继续复用`ListAllNodes`读取`media_node`，通过已有`NodeInfo`返回节点管理的全部`10`个字段。返回统一响应的`data.list`，无分页或筛选参数，空结果为数组。原分页管理列表和采集节点总览维持原有职责。

节点编号有数据库`UNIQUE`及`0..255`约束，正常数据最多`256`条；按编号升序完整返回，不人为添加截断。现有服务执行固定多表元数据检查后一次节点查询，不逐行查库，依赖编号唯一索引排序；无需新索引或迁移。该边界可支持无分页全量返回，复用已有服务和启动注入实例，不新增缓存或远程访问。

节点配置属于`platform_only`插件的共享基础设施配置，没有租户／部门归属，创建人和修改人仅为审计字段而非所有权；内部入口沿用`X-Inner-Api-Key`规则和既有显式空密钥配置语义，不额外按组织过滤，不增加新的放行方式，认证失败由原链拒绝。测试验证内部授权与拒绝访问、完整字段、空值、`256`条全量结果，以及已删除看板路由不存在。`media`未启用多语言，中文`OpenAPI`标签及说明按单语言维护，双语`README`同步。无缓存、依赖注入、宿主、前端、数据库迁移或跨平台工具变更。

### FB-49 存活节点汇总

在`summary`增加`alive_node_count int64`，与流、会话数量同一条语句读取。最终按用户要求的最近三分钟收报口径，使用服务端自动维护的`media_report_node.updated_at`判断`[now-3min, now]`闭区间，不看`status`，不依赖客户端`report_time`的单位、时钟或迟到数据。每条节点记录只计一次，不求和派生的`alive_nodes`，不要求节点承载`running`流。

当前节点表由网络上报路径维护；机器指标仅写实例表，不能单独刷新节点表。普通更新由`GoFrame`自动更新时间；无目的地址上报使用`Save/OnConflict`时，显式冲突更新字段此前漏了`updated_at`，本轮只把该列加入`OnDuplicate`，值仍由`GoFrame`自动产生，不手工赋时间。非上报来源若直接修改节点表，也会影响更新时间，应按该存储边界理解统计。

查询对`updated_at`直接做范围比较，当前迭代SQL增加幂等`idx_media_report_node_updated_at`索引，覆盖最近上报范围，避免对时间列做函数转换或加载明细。查询只读取一次当前时间，缺失、超时和未来时间不计入。无需新增字段或生成模型，但部署应执行新增索引DDL。三项统计仍为一次数据库往返。

新增时间窗口仅影响`summary`的节点计数，节点总览既有`alive_nodes`状态口径不在本次修改范围。无缓存、权限、依赖注入、前端或宿主变化，保留平台监控例外和既有双鉴权边界。通过时间上下界、状态与客户端采样时间独立、未来时间、数量上限及真实PostgreSQL重复上报更新时间回归验证。

### FB-48 实例采集与更新时间

实例列表新增可空`report_time`和`updated_at`。前者直接从已解码的`report_payload.timestamp`读取，支持`protobuf`整数字符串和`JSON`整数；沿用当前采集协议的秒／毫秒识别规则，正数小于`1000000000000`时乘`1000`，其余正整数按毫秒输出。不存在、非正或非法值返回`null`。不复用实例投影表的`report_time`，因为该列也会随流／会话计数事件改变，不等同于机器采集时间。

`updated_at`复用实体已查询的更新时间，通过已有`formatTime`转换为毫秒，缺失为`null`。两个字段在服务输出与控制器DTO中显式传递，原始报文保持完整。字段从同一批已读取内容装配，不改变查询、权限、存储、缓存、依赖注入或索引；无需迁移。测试覆盖秒／毫秒、字符串／数值、缺失、无效值、生命周期时间分离以及生成文档的可空数值契约。

### FB-47／FB-53 资源输出单位

按`FB-53`最终要求，节点总览和实例列表的`network_in/network_out`直接返回上报并保存的`KB/s`，覆盖根节点／子节点汇总、直属`instance_resources`和实例行。控制器仅映射字段，不再包含单位换算函数，不修改服务对象或原始报文，也不增加查询。

数据库结构、原始报文、磁盘展示单位和查询范围不变，无迁移、索引、缓存或权限变动；`report_payload`中的原始机器网络字段为`KB/s`。`OpenAPI`明确客户端原值及`KB/s`，双语文档同步更新。通过固定物理量、零值、小数、子节点和原始报文保真断言验证输出，防止重复单位换算。

### FB-44 全局活跃数量汇总

新增`GET /api/v1/media/dashboard/summary`，返回`active_stream_count`和`active_session_count`，类型为`int64`。无输入筛选。按后续`FB-46`用户修正，流数仅统计`status=running`且`close_time IS NULL`的记录；会话没有运行状态字段，只计入通过`stream_id`关联这些流的未关闭会话。两个标量子查询置于同一条`SELECT`，会话使用`EXISTS`限定来源流，避免独立查询在并发更新期间使用不同语句快照和联结行数乘积。不引入心跳超时判断或实例缓存累计值。只读取聚合结果，不加载`report_payload`或列表明细，没有`10000`条截断。

权限沿用`platform_only`媒体平台监控边界和现有管理接口双鉴权，声明`media:management:query`。这是用户要求的全平台总量，指标无人员／部门归属，不按组织数据权限过滤；不作为租户自助接口，也不新增匿名或内部密钥入口。鉴权失败由现有中间件拒绝，验证实际路由及权限元数据。

复用两表已有`close_time`索引、活跃记录部分索引、流`stream_id`主键及会话`stream_id`索引，固定一次数据查询，数据库扫描成本随活跃记录数增长，响应固定两个整数；不引入缓存、新索引、迁移或依赖注入。后端测试覆盖空值、混合关闭状态、孤立会话、不同租户／节点和超过列表上限的数据。同步插件自有`OpenAPI`与双语说明。插件未启用多语言，中文元数据不要求翻译资源；无前端、宿主或跨平台工具改动。

### FB-42 第三方查询返回完整上报

用户明确要求将已保存的完整上报返回给第三方展示。本轮扩展四类现有受保护查询的行级 DTO，增加`report_payload`对象；流和会话增加`report_subcommand`。原始内容和既有展示字段并列，第三方无需逐条补查。SQL NULL或没有记录时返回`null`，命令未记录返回`0`；不增加页面、匿名入口、独立重复接口或数据库迁移。

服务层从同批查询中读取原始 JSON，使用保留数值字面量的标准 JSON 解码，保证大整数不被浮点转换、扩展和嵌套内容不丢失；非法或非对象 JSON 作为查询错误返回。原始时间与64位整数属于外部 protobuf JSON 载荷的保真语义，保持字符串及原上报单位，API文档明确区别于规范化展示字段。原请求过滤、会话仅返回活跃记录、上限和网络截断保持不变；摘要及其他拓扑继续排除原始字段。返回大小随所选行及载荷增长，查询次数不随行数增长，仍建议第三方使用节点、租户、设备、实例和流过滤控制查询范围。

权限边界沿用现有平台监控和媒体鉴权链及`media:management:query`声明，不新增放行路径。此次由用户明确授权同一可见记录的完整`extra`等字段返回，原始内容不传播到原本不查询这些记录的聚合接口。无新运行期依赖、缓存、SQL、生成模型或跨平台工具变化。`media`为单语言插件，API字段说明使用中文；更新OpenAPI形状和完整性自动化验证。第三方自行展示，不修改宿主或插件前端。

### FB-41 四类上报完整性

逐字段审查发现当前投影仅覆盖看板所需内容，机器磁盘容量和 I/O 耗时/次数、流编码和服务器类型、全部原始`extra`、跳点名称/类型等存在遗漏；单位转换、精度舍入、默认值回退和 ID 合并也无法还原原始值。删除事件跳过 upsert，携带字段和首次 DELETE 数据被丢弃。

本轮按现有最新状态语义补齐，在`media_report_instance`、`plugin_media_report_network_link`、`media_report_stream`和`media_report_session`增加可空`report_payload JSON`，流和会话另存`report_subcommand`。使用`protojson`的`UseProtoNames`、`EmitUnpopulated`和`UseEnumNumbers`保存未归一化的全部已定义字段；64位整数按官方 JSON 规范保留字符串，嵌套跳点和`extra`完整保存。未接收新报文的旧行保持`NULL`，不回填伪造内容；JSON不承诺保留protobuf unknown wire字段，不新增逐条历史或新的采集协议。

正常上报把原始 JSON 与原查询投影一起写入；删除事件先以资源主键保存报文和命令元数据，已有投影不被稀疏删除内容覆盖，再执行既有计数和关闭逻辑，首次删除也保留完整报文。网络记录继续使用已有复合键和时间条件拒绝旧采样；空目的等无效业务键仍按原边界跳过。本轮不改变既有按到达顺序更新的机器、流、会话语义和关闭数据30天清理策略。

归属仍为`media`插件，沿用平台采集服务和受保护看板的数据权限。`FB-41`仅补存储时排除原始内容；后续`FB-42`按用户要求在指定四类查询中返回完整内容，其他摘要及聚合继续排除。新增列无查询索引需求，不增加每行查询或远程调用。SQL追加在当前活跃迭代`003`中，仅幂等 DDL；使用`make dao`生成模型。不新增服务实例、缓存、宿主能力、前端或业务枚举。

### FB-39 节点原始链路指标

节点总览沿用现有`GET /api/v1/media/dashboard/nodes/overview`，新增每节点`network_links`和响应根节点的`network_links_limited`。新增`plugin_media_report_network_link`最新投影表，使用`node_id + machine_id + source_ip + destination_ip`复合业务键，不使用软删除；该表为采样最新状态，不是业务操作历史。按当前返回子树的节点集合一次性查询，固定`10000`条返回上限并显式报告截断，使用以前导`node_id`组成的业务键索引支撑读取。网络指标不进入业务缓存，写入与查询使用共享数据库；不新增运行期依赖、跨模块契约或宿主修改。

保存上报原始`rtt`、`jitter`、`packet_loss`、`status_code`，以及`extra.reachable`明确提供的布尔状态，缺省状态为`null`。新字段`report_time`在采集边界统一为`Unix`毫秒时间戳。空目的地址不创建链路记录；没有链路上报时返回空数组。已有`node_latency_map`继续提供按目的地址组织的延迟摘要，完整链路列表保留来源实例维度。当前上游通过`HTTP`探测失败比例生成`packet_loss`，接口说明必须写明该口径，不能将失败固定值当作正常测量值或自动评分。

节点资源继续使用实例表批量聚合。`instance_resources`返回当前节点直属实例的资源汇总和实例数量，不包含子节点；既有顶层资源字段仍表示子树汇总，避免调用方在途经节点之间重复计算资源。按后续`FB-53`要求，上报的`cpu_usage`百分比原值保存为`cpu_load`，不换算已用核数；节点仅累加已有实例原值，接口明确其为原始利用率数值之和。已有样本由后续实例上报刷新，不增加质量评分、质量等级、路径推荐或综合负载计算。

`FB-40`补齐当前`net-flux v0.0.19`的真实标识差异：机器上报的`node_id`可能等于实例标识，而网络上报的`node_id`为业务节点号。读适配层保留明确实例节点；仅对空节点或实例标识占位的记录，按`machine_id = instance_id`联结全部节点范围的最新网络采样。最新采样同刻指向多个节点时不推断归属；多源和多目的记录属于同一节点时合并成一个实例归属。该联结作为原实例聚合的一部分在数据库内完成，先确定全局最新归属再展示所选子树，独立于链路明细的`10000`条上限。增加`machine_id + report_time + node_id`索引支撑查询，不修改上报原值或新增写入、缓存和运行期依赖。

数据权限沿用平台基础设施监控边界：`media`为`platform_only`插件，管理接口仍受既有认证链和`media:management:query`保护；指标无人员、部门归属，不新增按组织权限隔离的业务实体。新增链路查询在数据库侧按响应节点集合过滤，不扩大到子树外数据，不提供公开查询入口。测试覆盖子树过滤、多个实例与目的端、空结果、截断和原始失败值；路由测试验证现有权限声明。`media`未启用`i18n`，接口源文案使用中文，不新增翻译资源；生成配置使用既有跨平台`make dao`入口，无新增脚本。数据库变更追加在本活跃迭代的`003-add-media-data-collection-server.sql`，使用幂等建表，仅包含`DDL`。

- `i18n`影响：新增节点总览的中文 API 文档源文本；插件未启用`i18n`，无需翻译资源，不修改运行时 UI 文案、菜单或语言包。
- 缓存影响：新增 collection server 生命周期事件共享 cache，用于维护实例实时流数和会话数；权威运行时来源为客户端增删事件，写回表为`media_report_instance`；跨实例同步由宿主`cachecap.Service`背后的运行时 KV 后端负责，`cluster.enabled=true`时必须使用宿主 Redis/coordination KV 等共享后端，不允许退化为当前 Pod 内存；最大陈旧窗口为最后一次成功处理的生命周期事件到下一次事件重试之间；`collectionServer.discovery.cacheDir`仅透传给 Nacos SDK 作为其本地客户端缓存目录。
- 数据权限影响：扩展既有受保护节点总览接口，不新增公开读取入口。采集 server 写入插件自有基础设施投影；链路查询在数据库阶段按返回节点集合过滤，平台监控边界及例外依据见`FB-39`。
- 开发工具跨平台影响：不新增或修改 LinaPro 脚本、Makefile、CI 或`linactl`入口；外部`net-flux`协议代码按其现有`go generate ./...`入口生成，该入口依赖`protoc`。
- 数据库影响：本变更维护`media_report_node`、`media_report_node_snapshot`、`media_report_instance`、`media_report_stream`和`media_report_session`读模型；`FB-39`新增符合完整插件前缀的`plugin_media_report_network_link`表，并使用独立 DAO 生成配置维护可空布尔字段。


### FB-55 流列表节点标识

`GET /api/v1/media/dashboard/streams`在每个`stream_list[]`返回字符串`node_id`，直接由已查询的`media_report_stream.node_id`经服务输出和控制器DTO映射得到。空节点返回`""`，`"0"`保留为字符串；不以请求`source_id`或原始报文替代，所以实例来源和跨节点列表也能逐条识别所属节点。

此能力属于`media`插件的看板响应适配，无新抽象、查询、关联回查、缓存、依赖注入或存储变更，无`N+1`和数据库迁移。保留既有查询范围、排序、`10000`条上限、平台监控数据权限例外与管理双鉴权，不新增组织权限例外或公开入口。插件未启用`i18n`，字段使用中文说明并同步双语README。测试从数据库回读到DTO及OpenAPI覆盖节点字符串、空值和实际零标识；无前端页面或跨平台工具变更。
