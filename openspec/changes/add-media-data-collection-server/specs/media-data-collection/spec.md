## ADDED Requirements

### Requirement: 媒体插件应提供可配置的数据采集 server

`media`插件 MUST 在宿主启动后按插件配置启动一个兼容`net-flux` client 上报协议的 TCP server。

#### Scenario: 配置启用采集 server

WHEN 插件配置`collectionServer.enabled`为`true`
AND 宿主触发`system.started` hook
THEN `media`插件应在`collectionServer.addr`指定的地址启动采集 TCP server
AND 启动逻辑应复用宿主 hook 的上下文取消链完成停止。

#### Scenario: 默认关闭采集 server

WHEN 插件配置未显式启用`collectionServer.enabled`
AND 宿主触发`system.started` hook
THEN `media`插件不应启动采集 TCP server。

### Requirement: 采集 server 应兼容基础上报协议

采集 server MUST 处理`net-flux` client 端发送的基础系统心跳和数据上报命令。

#### Scenario: 响应系统心跳

WHEN client 端发送`Ping`
THEN server 端应返回包含相同时间戳的`Pong`。

#### Scenario: 接受指标上报

WHEN client 端上报`MachineMetric`、`NetworkMetric`、`StreamMetric`或`SessionMetric`
THEN server 端应接受请求、记录采集日志并按上报业务键写入`media`数据看板读模型
AND 不应因为未配置 discovery 服务而拒绝数据上报。

#### Scenario: 写入上报读模型

WHEN server 端收到包含有效业务键的`MachineMetric`
THEN 应将其按实例/容器指标处理并更新`media_report_instance`最新实例投影
AND `instance_id`为空时应回退使用`machine_id`作为实例/容器业务键
AND `memory_allocated`和`memory_used`应按`MB`写入
AND `disk_io_read`、`disk_io_write`、`network_in`和`network_out`应按`KB/S`写入
AND 不应从`MachineMetric`直接读取或覆盖实例直播流数和会话数。

WHEN server 端收到包含有效业务键的`NetworkMetric`
THEN 应更新对应节点的网络速率和延迟矩阵字段。
AND `throughput`写入`network_out`前应按`KB/S`归一化。

WHEN server 端收到包含有效业务键的`StreamMetric`
THEN 应更新`media_report_stream`最新流投影。
AND `media_report_stream`不得再存储`source_id`字段
AND 应写入`protocol_type`保存源流协议类型
AND 应写入`device_id`作为流所属设备业务标识，支持后续按设备维度统计。

WHEN server 端收到包含有效业务键的`SessionMetric`
THEN 应更新`media_report_session`最新会话投影。
AND `client_type`应按`1-mobile`、`2-pc`、`0-未知`枚举值存储
AND 应写入`device_id`作为会话所属设备业务标识，支持后续按设备维度统计。

WHEN 上报的`StreamMetric`或`SessionMetric`缺少可由服务端推导的时长字段
THEN server 端应根据`start_time`和`timestamp`计算`duration`或`play_duration`。

WHEN 上报的`SessionMetric.link_hops`包含链路跳点
THEN server 端应写入与看板响应一致的`hop_index`、`node_id`和`latency_ms`结构
AND `total_link_latency`缺省时应由链路跳点延迟求和得到。

### Requirement: 四类上报应完整保留最新有效报文

`media` MUST 在既有实例、网络链路、流和会话投影表中保存未经业务归一化的完整`report_payload`，保证`net-flux v0.0.21`中四类消息的全部已定义字段和嵌套字段可还原。该字段保存最新有效报文，不提供逐次历史归档；既有投影字段仍用于看板查询。

#### Scenario: 保留原始字段和扩展内容

WHEN 收到业务键有效的`MachineMetric`、`NetworkMetric`、`StreamMetric`或`SessionMetric`
THEN 服务端 MUST 使用`protobuf JSON`保存原始字段名、全部默认值、原始枚举数值、完整`extra`和完整`LinkHop`
AND `int64`字段 MUST 按 protobuf JSON 规则使用字符串保存，避免大整数精度损失
AND 原始 ID、单位、宽高、客户端类型和时间值 MUST 在归一化之前保留，不得由看板派生字段反向构造原始报文
AND 完整内容 MUST 使用 PostgreSQL `JSON`保存，以支持合法原始字符串中的转义`U+0000`，不得因使用`JSONB`而拒绝完整扩展内容
AND 现有无原始报文的行 MUST 标记为`NULL`，不得伪造历史数据或声称已回填。

#### Scenario: 保存删除事件的完整内容

WHEN 收到有效`STREAM_DELETE`或`SESSION_DELETE`
THEN 原始报文及`report_subcommand` MUST 保存，包括首次只收到删除事件的资源
AND 已存在的看板字段 MUST 不因稀疏删除报文而被清空，原有关闭与计数行为继续执行
AND 关闭流和会话仍按现有`30`天清理规则处理，不提供永久历史保留。

#### Scenario: 写入边界与查询成本

WHEN 收到缺少业务键的机器、流或会话报文，或缺少来源和目的地址业务键的网络报文
THEN 继续遵守无效投影不写入的边界，原始报文存储不作为无效报文归档
AND 网络链路的原始内容 MUST 与现有采样时间条件更新保持一致，较旧采样不得覆盖较新原始内容
AND 原始内容 MUST 仅通过明确声明返回该数据的受保护 HTTP DTO 暴露，其他聚合及摘要查询不得无必要地传输完整原始 JSON
AND 写入次数 MUST 保持固定且有界，不按字段或链路跳点逐项写库。

### Requirement: 四类展示查询应返回完整最新上报

`media` MUST 在既有实例、节点网络链路、流和会话查询响应中返回`report_payload`，满足第三方展示完整上报字段的需求，不开发第三方页面或增加重复查询入口。

#### Scenario: 返回完整 JSON 及命令

WHEN 第三方通过既有认证调用四类数据查询
THEN `/media/dashboard/instances`的`instance_list[]`、`/media/dashboard/nodes/overview`各节点的`network_links[]`、`/media/dashboard/streams`的`stream_list[]`及`/media/dashboard/sessions`的`protocol_list[].active_session_list[]` MUST 返回`report_payload`
AND `report_payload` MUST 为完整 JSON 对象，保留已有原始字段、完整`extra`、嵌套跳点、零值、枚举数值、大整数及转义空字符，不得返回二次编码的 JSON 字符串
AND 流和会话对象 MUST 同时返回`report_subcommand`，沿用保存的命令数值，未记录时为`0`
AND SQL NULL对应的旧记录 MUST 返回`report_payload: null`，不得使用归一化展示字段伪造原始报文
AND 原始内容中的时间和64位整数 MUST 保留外部 protobuf JSON 原始类型及单位；既有展示时间字段仍保持现有规范化契约。

#### Scenario: 保持查询范围与装配边界

WHEN 四类接口装配返回数据
THEN 原始内容 MUST 与原记录同批查询，不得逐行额外查询或调用外部服务
AND 原有查询过滤、权限声明、认证链、会话活跃状态及`10000`条上限 MUST 保持原语义
AND 节点树中每条原始网络报文 MUST 仍归属于对应来源节点，截断行为不变
AND 与本次四类列表无关的摘要和拓扑查询 MUST 继续排除不需要的原始 JSON
AND 流缺少协议摘要时，会话统计和会话明细中的相同协议 MUST 合并为一个分组，不得重复返回会话
AND 非法或非对象原始内容 MUST 返回明确查询错误，不得静默丢弃字段或伪装成无原始报文。

### Requirement: 生命周期事件应维护实时计数

`media` MUST 使用共享缓存和既有投影行维护幂等的流、会话实时计数及关闭状态。

#### Scenario: 根据生命周期事件维护实时计数

WHEN server 端收到同一实例的`STREAM_ADD`
THEN 应按`stream_id`幂等记录该流归属并递增该实例的实时直播流数量。
AND 该归属和计数 MUST 存储在宿主发布给`media`插件的共享 cache 中，不得仅存储在当前 Pod 进程内存中。

WHEN server 端收到同一实例的`STREAM_DELETE`
THEN 应按`stream_id`幂等删除该流归属并递减该实例的实时直播流数量。
AND 多个 Pod 处理同一`stream_id`事件时不应重复增减。
AND 应在`media_report_stream.close_time`记录关闭时间，而不是删除流投影。

WHEN server 端收到同一实例的`SESSION_ADD`
THEN 应按`session_id`幂等记录该会话归属并递增该实例的实时会话数量。
AND 该归属和计数 MUST 存储在宿主发布给`media`插件的共享 cache 中，不得仅存储在当前 Pod 进程内存中。

WHEN server 端收到同一实例的`SESSION_DELETE`
THEN 应按`session_id`幂等删除该会话归属并递减该实例的实时会话数量。
AND 多个 Pod 处理同一`session_id`事件时不应重复增减。
AND 应在`media_report_session.close_time`记录关闭时间，而不是删除会话投影。

WHEN server 端收到重复的`STREAM_ADD`、`STREAM_DELETE`、`SESSION_ADD`或`SESSION_DELETE`
THEN 实时直播流数和会话数不应重复增减。

#### Scenario: 重放同一采样点

WHEN 同一个实例的`MachineMetric`使用相同`instance_id`重复上报
THEN 最新实例投影应覆盖更新
AND 不应把容器 CPU、内存和磁盘指标写入`media_report_node`物理节点快照。

#### Scenario: 业务键缺失

WHEN server 端收到同时缺少`instance_id`和`machine_id`的`MachineMetric`
THEN 不应写入实例投影。

WHEN server 端收到缺少`session_id`的`SessionMetric`
THEN 不应写入会话投影。

### Requirement: 采集 server 应支持可配置的 Nacos discovery

采集 server MUST 在启用 discovery 配置后兼容`net-flux`示例 server 端的实例注册、注销和查询命令。

#### Scenario: 默认不连接 Nacos

WHEN 插件配置未显式启用`collectionServer.discovery.enabled`
AND client 端只发送系统心跳或指标上报命令
THEN server 端不应创建 Nacos discovery client
AND 应继续接受心跳和指标上报。

#### Scenario: 注册实例

WHEN 插件配置`collectionServer.discovery.enabled`为`true`
AND client 端发送`Instance` discovery 命令
THEN server 端应使用配置的 Nacos discovery client 注册该实例。

#### Scenario: 注销实例

WHEN 插件配置`collectionServer.discovery.enabled`为`true`
AND client 端发送`Deregister` discovery 命令
THEN server 端应按命令中的实例名、节点、IP 和端口注销该实例。

#### Scenario: 查询实例

WHEN 插件配置`collectionServer.discovery.enabled`为`true`
AND client 端发送`Lookup` discovery 命令
THEN server 端应按服务名和节点查询 Nacos 实例
AND 返回包含查询结果的`LookupAck`。

### Requirement: 媒体插件独立文档页应说明 TCP 采集协议

`media`插件 MUST 在独立接口文档页中提供 TCP 采集协议说明，帮助开发者区分 HTTP OpenAPI 接口与`net-flux`兼容 TCP 采集能力。

#### Scenario: 浏览媒体插件独立接口文档

WHEN 用户打开`/api/v1/media/apidocs.html`
THEN 页面应提供 HTTP API 视图
AND 页面应提供 TCP 采集协议说明视图
AND TCP 采集协议说明应覆盖监听配置、系统心跳、指标上报、生命周期事件、Nacos discovery 命令和联调入口。

#### Scenario: 生成媒体插件 OpenAPI 文档

WHEN 服务端生成`/api/v1/media/openapi.json`
THEN OpenAPI 文档只应包含真实 HTTP 路由
AND 不应为 TCP 采集 server 伪造 HTTP 路径。

### Requirement: 媒体插件应提供数据看板读接口

`media`插件 MUST 提供受`media:management:query`权限保护的数据看板查询接口，读取`media_report_*`最新投影表并返回`数据看板.md`中节点总览、实例列表、流列表和会话列表所需字段。

接口响应的`data`对象 MUST 保持`apps/lina-plugins/media/数据看板.md`示例中的直接对象结构，并支持本规范明确新增的原始链路和直属实例资源字段；响应`data`不得额外返回列表包装、`total`、`pageNum`或`pageSize`字段。不返回`total`的看板列表接口 MUST 不暴露分页查询参数，服务端在数据库查询阶段使用固定`10000`条上限保护数据库。

#### Scenario: 查询节点总览树

WHEN 管理端请求节点总览接口
THEN 服务端应从`media_report_node`一次性读取有界节点集合
AND 按`parent_node_id`在内存中组装`child_nodes`
AND 响应`data`应直接为当前节点对象，不得额外包裹`nodes`数组字段
AND 响应字段应覆盖`node_id`、`node_name`、`region`、`status`、`parent_node_id`、资源指标、实时计数、`avg_delay`、`last_heartbeat`、`report_time`、`node_latency_map`、`network_links`、`network_links_limited`和`child_nodes`。

#### Scenario: 查询节点间原始网络指标及资源

WHEN 采集端上报`NetworkMetric`后管理端请求`GET /api/v1/media/dashboard/nodes/overview`
THEN 服务端 MUST 按`node_id + machine_id + source_ip + destination_ip`保存最新网络指标，避免同节点不同实例或目的端相互覆盖
AND 每个节点的`network_links` MUST 返回来源实例、源地址、目的地址、`rtt`、`jitter`、`packet_loss`、`status_code`、可选`reachable`和毫秒时间戳`report_time`
AND 延迟和抖动单位为毫秒，丢包率使用`0～1`原始比例；当前上游探测器的丢包率口径为`HTTP`探测失败比例
AND `reachable`缺省或无法解析时 MUST 返回`null`，不得把未上报状态推断为成功；探测失败时的原始延迟、抖动和丢包率必须保留
AND 链路查询 MUST 在数据库侧限定到当前返回子树的节点，按稳定顺序一次性最多返回`10000`条，截断时在响应根节点设置`network_links_limited=true`
AND 每个节点现有`cpu_allocated`、`cpu_load`、内存、磁盘和网络资源字段 MUST 继续通过实例投影批量聚合，供调用方按途经节点读取资源，不得逐节点查询实例
AND `instance_resources` MUST 单独返回当前节点直属实例的资源汇总和实例数量，不包含子节点，防止路径读取时重复计算子树资源
AND `MachineMetric.cpu_usage`百分比 MUST 按分配核数转换为`cpu_load`已用核数，确保实例与节点资源口径一致
AND 本次 MUST NOT 计算质量评分、质量分级、候选路径排序或综合资源负载百分比。

#### Scenario: 关联上报节点标识不同的实例资源

WHEN 机器上报的`node_id`为空或等于`instance_id`，而网络上报使用业务节点号
THEN 服务端 MUST 根据`NetworkMetric.machine_id = MachineMetric.instance_id`关联已有实例资源
AND 明确且不同于实例标识的`MachineMetric.node_id` MUST 优先保留
AND 归属 MUST 按该机器全部节点范围的最新采样时间确定，之后再展示当前子树，不能使用已截断的链路明细推断
AND 同一最新时间的多目的端或多源地址指向同一节点时 MUST 仅计算一次实例资源
AND 同一最新时间指向不同节点或没有关联记录时 MUST 视为归属未知，不将实例标识伪装成业务节点
AND 机器指标和网络指标先后到达 MUST 得到相同关联结果，迁出子树的实例 MUST 不再计入旧节点
AND 关联及资源聚合 MUST 使用固定次数的集合查询，不得逐实例回查链路。

#### Scenario: 查询实例列表

WHEN 管理端按`node_id`、状态或关键词查询实例列表
THEN 服务端应在数据库查询阶段过滤、排序并最多读取`10000`条`media_report_instance`
AND 返回`node_info`和`instance_list`
AND 响应`data`不得返回`total`
AND 请求参数不得包含`pageNum`或`pageSize`
AND 不应通过逐实例回查节点配置表装配展示字段。

#### Scenario: 查询流列表

WHEN 管理端按`source_type + source_id`、`tenant_id`、`device_id`、`node_id`、`instance_id`、状态或关键词查询流列表
THEN 服务端应在数据库查询阶段过滤、排序并最多读取`10000`条`media_report_stream`
AND 返回`source_type`、`source_id`和`stream_list`
AND 响应`data`不得返回`total`
AND 请求参数不得包含`pageNum`或`pageSize`
AND `source_type + source_id`筛选应映射到`node_id`或`instance_id`数据库字段，不依赖已删除的`media_report_stream.source_id`
AND `stream_list[].duration`应通过`start_time`计算；若`close_time`不为空则计算`close_time - start_time`，否则计算当前时间与`start_time`的差值
AND `protocol_summary`应直接由`media_report_stream.protocol_summary`解析得到
AND `protocol_count`和`total_sessions_lifetime`应优先由`protocol_summary`聚合计算，缺少协议摘要时才回退使用读模型字段
AND `current_active_sessions`和`protocol_summary[].current_sessions`应由`close_time`为空的`media_report_session`活跃会话投影聚合计算。

#### Scenario: 查询会话列表

WHEN 管理端按`stream_id`查询会话列表
THEN 服务端应读取`media_report_stream`获取`stream_info`和协议摘要
AND 在数据库查询阶段按`stream_id`、`tenant_id`、`device_id`、`protocol_type`、`node_id`或`instance_id`过滤、排序并最多读取`10000`条`media_report_session`
AND 仅返回`close_time`为空的活跃会话明细
AND 按`protocol_type`聚合当前查询范围内未关闭会话数并组装`protocol_list`
AND 响应`data`不得返回`total`
AND 请求参数不得包含`pageNum`或`pageSize`
AND `active_session_list[].play_duration`应通过`start_time`计算；若`close_time`不为空则计算`close_time - start_time`，否则计算当前时间与`start_time`的差值
AND `link_hops`应直接由`media_report_session.link_hops`解析得到
AND `total_link_latency`应优先由`link_hops[].latency_ms`求和，缺少链路跳点时才回退使用读模型字段。

#### Scenario: 查询视频网关服务监测拓扑

WHEN 管理端请求视频网关服务监测拓扑接口
THEN 服务端应提供单个`GET`查询接口，受`media:management:query`权限保护
AND 请求入参应仅包含单个必填`device_id`
AND 服务端应在数据库查询阶段按`device_id`过滤并最多读取`10000`条流投影和`10000`条会话明细
AND 数据库应提供以`device_id`为前导列并限定`close_time`为空的活跃流和活跃会话索引，以支撑单设备拓扑查询、协议租户聚合和稳定排序
AND 服务端应一次性读取有界流集合、当前活跃会话明细和协议租户聚合统计，不得让前端按流、协议或租户逐项补查
AND 响应应按设备、流、推流协议、租户和用户会话组织拓扑树
AND 设备和流节点应包含设备编码、基础平台流参数、视频网关流参数和当前协议复用数
AND 租户节点应包含当前查询范围内的并发会话数
AND 用户节点应包含会话、用户、客户端、播放协议、节点和实例信息
AND 响应业务字段必须来自`media_report_stream`、`media_report_session`、`media_report_node`的真实投影或数据库聚合结果
AND 在上报协议和数据库投影提供对应字段前，响应不得暴露以`stream_name`替代的`device_name`、以`stream_id`替代的`video_code`、以`client_ip`替代的`server_ip`，也不得暴露始终为空的`browser`或`operating_system`。

#### Scenario: 查询租户媒体节点拓扑

WHEN 管理端请求租户媒体节点拓扑接口
THEN 服务端应提供单个`GET`查询接口，受`media:management:query`权限保护
AND 请求入参应仅包含单个必填`tenant_id`，不得包含行政区划筛选
AND 使用铁塔 token 回退鉴权时，请求`tenant_id`必须与 token 用户所属租户一致，缺失或不一致时应拒绝查询
AND 服务端应在数据库查询阶段按`tenant_id`和`close_time`为空过滤活跃流与活跃会话
AND 响应应返回租户当前并发会话数、当前活跃流数、流协议种类及各协议活跃流数
AND 响应应按`node_id`返回最多`10000`个节点，每个节点包含节点名称、当前并发会话数、当前活跃流数、流协议种类及各协议活跃流数
AND 服务端应使用数据库聚合和节点名称批量投影完成数据装配，不得按节点逐项查询
AND 数据库应提供以`tenant_id`、`node_id`和`protocol_type`为前导维度并限定`close_time`为空的活跃流索引；活跃会话应复用现有租户节点部分索引
AND 响应业务字段必须来自`media_report_stream`、`media_report_session`和`media_report_node`的真实投影或数据库聚合结果
AND 在报表投影提供租户名称前，响应不得暴露`tenant_name`，租户显示名称应由调用方使用已选择的租户数据展示。

#### Scenario: 定期清理已关闭流和会话

WHEN `media`插件定时任务运行每周清理流程
THEN 应删除`close_time`早于当前时间`30`天的`media_report_stream`记录
AND 应删除`close_time`早于当前时间`30`天的`media_report_session`记录
AND `close_time`为空的未关闭流和会话不得被清理。

#### Scenario: 定时任务注册失败

WHEN `media`插件启动时注册每周清理定时任务失败
THEN `system.started` hook 应返回启动错误
AND 不得将清理定时任务标记为已启动
AND 后续启动尝试应允许重新注册定时任务。

#### Scenario: 看板接口数据权限与性能边界

WHEN 查询看板流或会话数据时携带`tenant_id`筛选
THEN 服务端 MUST 在数据库查询条件中注入该筛选，不得先查询全量再在内存中过滤。

WHEN `media`插件以`platform_only`模式运行
THEN 看板接口的数据权限边界由宿主认证权限和业务筛选条件共同约束
AND 不额外引入租户安装态隔离。

WHEN 任一看板列表接口返回动态结果集
THEN 接口 MUST 提供数量上限
AND 后端查询次数不得随返回行数线性增长。


### Requirement: 提供全局活跃流和会话汇总

`media` MUST 提供受保护的`GET /api/v1/media/dashboard/summary`，无需查询参数，返回`active_stream_count`、`active_session_count`及`alive_node_count`三个整数，供第三方展示全局总量。

#### Scenario: 精确统计当前未关闭记录

WHEN 已认证调用方读取全局汇总
THEN 系统 MUST 统计流表中`status=running`且`close_time IS NULL`的流记录数
AND 会话总数 MUST 仅统计这些运行中流通过`stream_id`关联的未关闭会话；非运行中流、已关闭流和不存在流的会话不得计入
AND 已关闭记录 MUST 被排除，空表 MUST 返回三个`0`
AND 总量 MUST 不受列表`10000`条上限、协议分组、实例上报累计计数或节点汇总重复值影响
AND 三个统计 MUST 在同一条数据库语句中完成，不加载明细或原始报文。

#### Scenario: 保持媒体平台监控授权边界

WHEN 调用方请求汇总
THEN 接口 MUST 声明`media:management:query`并沿用媒体管理接口的既有宿主／铁塔鉴权链
AND 认证均失败时 MUST 拒绝，不能仅凭`X-Inner-Api-Key`访问
AND 返回值表示平台监控总量，不按租户筛选，也不作为租户自助接口。


### Requirement: 看板查询文档应反映实际筛选与生命周期语义

`media` MUST 在`OpenAPI`中准确说明流列表、会话列表和总量统计的不同范围，保持既有运行时行为。

#### Scenario: 状态和协议与采集映射一致

WHEN 调用方查看流状态筛选及返回字段
THEN 文档 MUST 列明`universal`、`running`、`inactive`、`error`、`closed`、`timeout`、`cancelled`、`failed`及其含义，并采用`running`示例
AND 拓扑中的流状态说明 MUST 与同一映射一致
AND 协议筛选 MUST 列出采集端的标准协议字符串，说明按存储值匹配，不以文档修正新增枚举校验。

#### Scenario: 明确范围和统计口径

WHEN 调用方查看三个接口说明
THEN 流列表 MUST 说明默认包含未关闭和已关闭记录，状态不传时不筛选
AND 会话列表 MUST 说明必填`streamId`、只返回`close_time IS NULL`记录，以及明细共享`10000`条上限而分组计数不受该上限影响
AND 汇总 MUST 明确仅统计`running`且未关闭的流及其未关闭会话，不根据心跳超时判断是否播放
AND 文档 MUST 说明可选参数默认行为、既有排序和统一响应的`data`位置。


### Requirement: 节点和实例网络速率以 Kbps 输出

节点总览和实例列表的展示字段`network_in`、`network_out` MUST 使用`Kbps`，定义为每秒千比特，`1000 bps = 1 Kbps`。

#### Scenario: 输出层统一换算

WHEN 调用方查询`/media/dashboard/nodes/overview`或`/media/dashboard/instances`
THEN 节点、所有子节点、各节点`instance_resources`及`instance_list`的网络速率 MUST 从已有`KiB/s`投影转换为`Kbps`，换算系数为`1024 × 8 ÷ 1000`
AND 零值和小数 MUST 保留，不得取整或重复转换
AND 原始`report_payload`中的同名字段 MUST 保留客户端上报的`bps`值
AND 磁盘速率、数据库存储和采集写入单位 MUST 保持原有语义，不增加查询或迁移。


### Requirement: 实例列表独立返回采集与更新时间

实例列表`instance_list[]` MUST 返回可空的`report_time`和`updated_at`，数值均为`Unix`毫秒时间戳。

#### Scenario: 提取原始采集时间

WHEN 实例具有完整机器上报且`report_payload.timestamp`为有效正整数
THEN `report_time` MUST 取该原始时间并转换为毫秒，支持`protobuf`字符串整数及`JSON`整数，秒值乘以`1000`，毫秒值保持不变
AND `report_payload` MUST 保持原始字段及时间值不变
AND 机器采集时间 MUST 不回退到可能被生命周期计数更新的数据库`report_time`或当前时间。

#### Scenario: 返回更新时间并明确未知值

WHEN 实例列表装配数据
THEN `updated_at` MUST 使用现有记录的数据库更新时间并转换为毫秒
AND 缺失、非正或无法解析的原始采集时间 MUST 返回`report_time: null`
AND 缺失数据库更新时间 MUST 返回`updated_at: null`
AND 两个字段 MUST 与原始报文及实体同批读取，不增加查询或数据库迁移。


### Requirement: 汇总返回存活节点数

`GET /api/v1/media/dashboard/summary` MUST 返回`alive_node_count`，表示全平台上报节点中的存活数量。

#### Scenario: 最近三分钟上报即存活

WHEN 计算存活节点数
THEN 系统 MUST 按服务端维护的`media_report_node.updated_at`是否位于当前时间前`3`分钟至当前时间的闭区间计数，不依赖客户端采样时间
AND 节点`status` MUST 不参与判断，近期上报的离线或空状态节点也计入，旧的`healthy`节点不计入
AND 缺失、超时和未来更新时间 MUST 不计入
AND 有目的地址和无目的地址的节点网络上报 MUST 都刷新服务端更新时间，迟到采样时间不影响存活判断
AND 每个节点记录 MUST 只计一次，不得求和各行的`alive_nodes`，不要求节点承载`running`流
AND 空数据 MUST 返回`0`，总量不受节点树`10000`条上限影响
AND 三个总量 MUST 在同一条聚合语句中读取，不装配节点树或原始报文。
