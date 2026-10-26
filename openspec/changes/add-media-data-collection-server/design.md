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
| `MachineMetric` | `media_report_instance` | `instance_id`作为实例业务键，缺省时回退使用`machine_id`承载的容器/实例标识；节点、状态、CPU、内存、磁盘、网络、版本和启动时间写入实例最新投影；内存字段统一归一化为`MB`，磁盘和网络速率字段统一归一化为`KB/S`；`live_streams`和`sessions`不由客户端直接上报。 |
| `NetworkMetric` | `media_report_node`、`plugin_media_report_network_link` | 优先使用`node_id`作为节点业务键，缺省时回退`machine_id`和`source_ip`；吞吐量按`KB/S`归一化，`destination_ip -> rtt`在节点行级锁下合并进`node_latency_map`。另按节点、来源实例、源和目的地址保存完整原始链路指标。 |
| `StreamMetric` | `media_report_stream`、`media_report_instance` | `stream_id`作为流业务键；`tenant_id`、`device_id`、`node_id`、`instance_id`写入租户、设备、节点和实例维度，`instance_id`为空时回退使用`machine_id`承载的容器/实例标识；协议、状态、码率、分辨率、帧率、丢包率、延迟、会话摘要和流路径写入最新流投影；`STREAM_ADD`和`STREAM_DELETE`按`stream_id`幂等维护实例实时直播流数量。 |
| `SessionMetric` | `media_report_session`、`media_report_instance` | `session_id`作为会话业务键；流、租户、设备、客户端、播放协议、播放质量、`node_id`、实例和链路跳点写入会话最新投影，`instance_id`为空时回退使用`machine_id`承载的容器/实例标识；`SESSION_ADD`和`SESSION_DELETE`按`session_id`幂等维护实例实时会话数量。 |

`media_report_instance`和`media_report_session`不再伪造业务键，必须分别由`MachineMetric.instance_id`或`MachineMetric.machine_id`、`SessionMetric.session_id`驱动写入；缺少业务键的上报包会被忽略。采集端无法稳定获取宿主机器信息时，不应把容器 CPU、内存和磁盘指标伪装成`media_report_node`物理节点快照。

## 性能和数据边界

上报处理只按单个上报包执行固定次数写入，不在动态结果集中循环查询数据库或调用远程插件能力。`MachineMetric`按实例业务键执行一次最新投影 upsert；`NetworkMetric`按节点业务键在事务内锁定单行并合并当前延迟矩阵，避免同节点并发目的端覆盖；`StreamMetric`和`SessionMetric`分别按流和会话业务键执行一次最新投影 upsert。看板后续查询可直接读取投影表并依赖`node_id`、`report_time`、`source_type/source_id`、`tenant_id/status`和`tenant_id/device_id`等索引，避免实时回查配置表和`N+1`装配。discovery 查询只按单个服务名和节点 ID 访问 Nacos，不在 LinaPro 数据库中装配动态结果集。

实例实时流数和会话数的权威运行时来源是宿主共享 cache 中维护的生命周期事件状态：`stream_id -> instance_id`、`session_id -> instance_id`、`instance_id -> live_streams`和`instance_id -> sessions`。源码插件只能通过`cachecap.Service`访问该状态，不能直接依赖宿主私有 Redis、`kvcache`或 coordination 实现；在`cluster.enabled=true`时，该服务由宿主统一切换到 Redis/coordination KV 等共享后端，在单机模式下可继续使用宿主 SQL KV 后端。采集 server 在处理生命周期事件时使用对应流或会话投影行的数据库行级锁串行化同一资源的事件，再更新共享 cache 和实例投影，保证多 Pod 下重复`ADD`、重复`DELETE`和跨实例迁移不会重复增减。cache 是实时计数权威状态，`media_report_instance`中的`live_streams`和`sessions`是供数据看板读取的最新投影；若共享 cache 写入失败，当前上报处理返回错误并保留下一次生命周期事件重试恢复的路径。

## 影响评估

### FB-41 四类上报完整性

逐字段审查发现当前投影仅覆盖看板所需内容，机器磁盘容量和 I/O 耗时/次数、流编码和服务器类型、全部原始`extra`、跳点名称/类型等存在遗漏；单位转换、精度舍入、默认值回退和 ID 合并也无法还原原始值。删除事件跳过 upsert，携带字段和首次 DELETE 数据被丢弃。

本轮按现有最新状态语义补齐，在`media_report_instance`、`plugin_media_report_network_link`、`media_report_stream`和`media_report_session`增加可空`report_payload JSON`，流和会话另存`report_subcommand`。使用`protojson`的`UseProtoNames`、`EmitUnpopulated`和`UseEnumNumbers`保存未归一化的全部已定义字段；64位整数按官方 JSON 规范保留字符串，嵌套跳点和`extra`完整保存。未接收新报文的旧行保持`NULL`，不回填伪造内容；JSON不承诺保留protobuf unknown wire字段，不新增逐条历史或新的采集协议。

正常上报把原始 JSON 与原查询投影一起写入；删除事件先以资源主键保存报文和命令元数据，已有投影不被稀疏删除内容覆盖，再执行既有计数和关闭逻辑，首次删除也保留完整报文。网络记录继续使用已有复合键和时间条件拒绝旧采样；空目的等无效业务键仍按原边界跳过。本轮不改变既有按到达顺序更新的机器、流、会话语义和关闭数据30天清理策略。

归属仍为`media`插件，沿用平台采集服务和受保护看板的数据权限，不新增公开原始报文接口；原始`extra`可能包含业务敏感内容，不透传到看板 DTO。新增列无查询索引需求，常用列表显式排除原始 JSON，不增加每行查询或远程调用。SQL追加在当前活跃迭代`003`中，仅幂等 DDL；使用`make dao`生成模型。不新增服务实例、缓存、宿主能力、前端、业务枚举或运行时文案，`i18n`、DI与跨平台生产工具无新增影响。

### FB-39 节点原始链路指标

节点总览沿用现有`GET /api/v1/media/dashboard/nodes/overview`，新增每节点`network_links`和响应根节点的`network_links_limited`。新增`plugin_media_report_network_link`最新投影表，使用`node_id + machine_id + source_ip + destination_ip`复合业务键，不使用软删除；该表为采样最新状态，不是业务操作历史。按当前返回子树的节点集合一次性查询，固定`10000`条返回上限并显式报告截断，使用以前导`node_id`组成的业务键索引支撑读取。网络指标不进入业务缓存，写入与查询使用共享数据库；不新增运行期依赖、跨模块契约或宿主修改。

保存上报原始`rtt`、`jitter`、`packet_loss`、`status_code`，以及`extra.reachable`明确提供的布尔状态，缺省状态为`null`。新字段`report_time`在采集边界统一为`Unix`毫秒时间戳。空目的地址不创建链路记录；没有链路上报时返回空数组。已有`node_latency_map`继续提供按目的地址组织的延迟摘要，完整链路列表保留来源实例维度。当前上游通过`HTTP`探测失败比例生成`packet_loss`，接口说明必须写明该口径，不能将失败固定值当作正常测量值或自动评分。

节点资源继续使用实例表批量聚合。`instance_resources`返回当前节点直属实例的资源汇总和实例数量，不包含子节点；既有顶层资源字段仍表示子树汇总，避免调用方在途经节点之间重复计算资源。上报的`cpu_usage`为百分比，入库转换为`cpu_count * cpu_usage / 100`已用核数，与既有接口单位一致；已有样本由后续实例上报刷新。此处仅统一指标单位，不增加质量评分、质量等级、路径推荐或综合负载计算。

`FB-40`补齐当前`net-flux v0.0.19`的真实标识差异：机器上报的`node_id`可能等于实例标识，而网络上报的`node_id`为业务节点号。读适配层保留明确实例节点；仅对空节点或实例标识占位的记录，按`machine_id = instance_id`联结全部节点范围的最新网络采样。最新采样同刻指向多个节点时不推断归属；多源和多目的记录属于同一节点时合并成一个实例归属。该联结作为原实例聚合的一部分在数据库内完成，先确定全局最新归属再展示所选子树，独立于链路明细的`10000`条上限。增加`machine_id + report_time + node_id`索引支撑查询，不修改上报原值或新增写入、缓存和运行期依赖。

数据权限沿用平台基础设施监控边界：`media`为`platform_only`插件，管理接口仍受既有认证链和`media:management:query`保护；指标无人员、部门归属，不新增按组织权限隔离的业务实体。新增链路查询在数据库侧按响应节点集合过滤，不扩大到子树外数据，不提供公开查询入口。测试覆盖子树过滤、多个实例与目的端、空结果、截断和原始失败值；路由测试验证现有权限声明。`media`未启用`i18n`，接口源文案使用中文，不新增翻译资源；生成配置使用既有跨平台`make dao`入口，无新增脚本。数据库变更追加在本活跃迭代的`003-add-media-data-collection-server.sql`，使用幂等建表，仅包含`DDL`。

- `i18n`影响：新增节点总览的中文 API 文档源文本；插件未启用`i18n`，无需翻译资源，不修改运行时 UI 文案、菜单或语言包。
- 缓存影响：新增 collection server 生命周期事件共享 cache，用于维护实例实时流数和会话数；权威运行时来源为客户端增删事件，写回表为`media_report_instance`；跨实例同步由宿主`cachecap.Service`背后的运行时 KV 后端负责，`cluster.enabled=true`时必须使用宿主 Redis/coordination KV 等共享后端，不允许退化为当前 Pod 内存；最大陈旧窗口为最后一次成功处理的生命周期事件到下一次事件重试之间；`collectionServer.discovery.cacheDir`仅透传给 Nacos SDK 作为其本地客户端缓存目录。
- 数据权限影响：扩展既有受保护节点总览接口，不新增公开读取入口。采集 server 写入插件自有基础设施投影；链路查询在数据库阶段按返回节点集合过滤，平台监控边界及例外依据见`FB-39`。
- 开发工具跨平台影响：不新增或修改 LinaPro 脚本、Makefile、CI 或`linactl`入口；外部`net-flux`协议代码按其现有`go generate ./...`入口生成，该入口依赖`protoc`。
- 数据库影响：本变更维护`media_report_node`、`media_report_node_snapshot`、`media_report_instance`、`media_report_stream`和`media_report_session`读模型；`FB-39`新增符合完整插件前缀的`plugin_media_report_network_link`表，并使用独立 DAO 生成配置维护可空布尔字段。
