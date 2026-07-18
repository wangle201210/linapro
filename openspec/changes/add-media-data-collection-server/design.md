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

采集 server 复用`github.com/dellinger2023/net-flux`发布的`network.NewTcpServer`和`gen`协议类型。为补齐数据看板所需业务键，本次 fork 到`github.com/wangle201210/net-flux`并追加上报协议字段，通过`go generate ./...`重新生成协议代码；`media`插件仍保留原 import 路径，并通过 Go module `replace`指向 fork 版本。

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
| `MachineMetric` | `media_report_instance` | `instance_id`作为实例业务键，缺省时回退使用`machine_id`承载的容器/实例标识；节点、状态、CPU、内存、磁盘、网络、版本和启动时间写入实例最新投影；内存字段统一归一化为`MB`，磁盘和网络速率字段统一归一化为`KB/S`；`live_streams`和`sessions`不由客户端直接上报。 |
| `NetworkMetric` | `media_report_node` | `machine_id`作为节点业务键；`throughput`按`KB/S`归一化后写入`network_out`，`destination_ip -> rtt`在事务内按节点行级锁合并进`node_latency_map`，避免同节点并发上报覆盖其他目的端延迟。 |
| `StreamMetric` | `media_report_stream`、`media_report_instance` | `stream_id`作为流业务键；`tenant_id`、`device_id`、`node_id`、`instance_id`写入租户、设备、节点和实例维度，`instance_id`为空时回退使用`machine_id`承载的容器/实例标识；协议、状态、码率、分辨率、帧率、丢包率、延迟、会话摘要和流路径写入最新流投影；`STREAM_ADD`和`STREAM_DELETE`按`stream_id`幂等维护实例实时直播流数量。 |
| `SessionMetric` | `media_report_session`、`media_report_instance` | `session_id`作为会话业务键；流、租户、设备、客户端、播放协议、播放质量、`node_id`、实例和链路跳点写入会话最新投影，`instance_id`为空时回退使用`machine_id`承载的容器/实例标识；`SESSION_ADD`和`SESSION_DELETE`按`session_id`幂等维护实例实时会话数量。 |

`media_report_instance`和`media_report_session`不再伪造业务键，必须分别由`MachineMetric.instance_id`或`MachineMetric.machine_id`、`SessionMetric.session_id`驱动写入；缺少业务键的上报包会被忽略。采集端无法稳定获取宿主机器信息时，不应把容器 CPU、内存和磁盘指标伪装成`media_report_node`物理节点快照。

## 性能和数据边界

上报处理只按单个上报包执行固定次数写入，不在动态结果集中循环查询数据库或调用远程插件能力。`MachineMetric`按实例业务键执行一次最新投影 upsert；`NetworkMetric`按节点业务键在事务内锁定单行并合并当前延迟矩阵，避免同节点并发目的端覆盖；`StreamMetric`和`SessionMetric`分别按流和会话业务键执行一次最新投影 upsert。看板后续查询可直接读取投影表并依赖`node_id`、`report_time`、`source_type/source_id`、`tenant_id/status`和`tenant_id/device_id`等索引，避免实时回查配置表和`N+1`装配。discovery 查询只按单个服务名和节点 ID 访问 Nacos，不在 LinaPro 数据库中装配动态结果集。

实例实时流数和会话数的权威运行时来源是宿主共享 cache 中维护的生命周期事件状态：`stream_id -> instance_id`、`session_id -> instance_id`、`instance_id -> live_streams`和`instance_id -> sessions`。源码插件只能通过`cachecap.Service`访问该状态，不能直接依赖宿主私有 Redis、`kvcache`或 coordination 实现；在`cluster.enabled=true`时，该服务由宿主统一切换到 Redis/coordination KV 等共享后端，在单机模式下可继续使用宿主 SQL KV 后端。采集 server 在处理生命周期事件时使用对应流或会话投影行的数据库行级锁串行化同一资源的事件，再更新共享 cache 和实例投影，保证多 Pod 下重复`ADD`、重复`DELETE`和跨实例迁移不会重复增减。cache 是实时计数权威状态，`media_report_instance`中的`live_streams`和`sessions`是供数据看板读取的最新投影；若共享 cache 写入失败，当前上报处理返回错误并保留下一次生命周期事件重试恢复的路径。

## 影响评估

- `i18n`影响：不新增运行时 UI 文案、菜单、语言包或 API 文档源文本。
- 缓存影响：新增 collection server 生命周期事件共享 cache，用于维护实例实时流数和会话数；权威运行时来源为客户端增删事件，写回表为`media_report_instance`；跨实例同步由宿主`cachecap.Service`背后的运行时 KV 后端负责，`cluster.enabled=true`时必须使用宿主 Redis/coordination KV 等共享后端，不允许退化为当前 Pod 内存；最大陈旧窗口为最后一次成功处理的生命周期事件到下一次事件重试之间；`collectionServer.discovery.cacheDir`仅透传给 Nacos SDK 作为其本地客户端缓存目录。
- 数据权限影响：不新增 HTTP API 或数据读取接口；新增写入路径为采集 server 接收设备上报并写入插件自有读模型，后续查询接口必须在数据库查询阶段使用冗余的`tenant_id`、`node_id`、`instance_id`等维度过滤。
- 开发工具跨平台影响：不新增或修改 LinaPro 脚本、Makefile、CI 或`linactl`入口；外部`net-flux`协议代码按其现有`go generate ./...`入口生成，该入口依赖`protoc`。
- 数据库影响：新增`media_report_node`、`media_report_node_snapshot`、`media_report_instance`、`media_report_stream`和`media_report_session`插件表，并更新`media`插件 DAO 生成配置。
