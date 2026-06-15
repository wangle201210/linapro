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
AND 不应从`MachineMetric`直接读取或覆盖实例直播流数和会话数。

WHEN server 端收到包含有效业务键的`NetworkMetric`
THEN 应更新对应节点的网络速率和延迟矩阵字段。

WHEN server 端收到包含有效业务键的`StreamMetric`
THEN 应更新`media_report_stream`最新流投影。

WHEN server 端收到包含有效业务键的`SessionMetric`
THEN 应更新`media_report_session`最新会话投影。

WHEN 上报的`StreamMetric`或`SessionMetric`缺少可由服务端推导的时长字段
THEN server 端应根据`start_time`和`timestamp`计算`duration`或`play_duration`。

WHEN 上报的`SessionMetric.link_hops`包含链路跳点
THEN server 端应写入与看板响应一致的`hop_index`、`node_id`和`latency_ms`结构
AND `total_link_latency`缺省时应由链路跳点延迟求和得到。

#### Scenario: 根据生命周期事件维护实时计数

WHEN server 端收到同一实例的`STREAM_ADD`
THEN 应按`stream_id`幂等记录该流归属并递增该实例的实时直播流数量。
AND 该归属和计数 MUST 存储在宿主发布给`media`插件的共享 cache 中，不得仅存储在当前 Pod 进程内存中。

WHEN server 端收到同一实例的`STREAM_DELETE`
THEN 应按`stream_id`幂等删除该流归属并递减该实例的实时直播流数量。
AND 多个 Pod 处理同一`stream_id`事件时不应重复增减。
AND 不应继续在`media_report_stream`保留该流的活跃最新投影。

WHEN server 端收到同一实例的`SESSION_ADD`
THEN 应按`session_id`幂等记录该会话归属并递增该实例的实时会话数量。
AND 该归属和计数 MUST 存储在宿主发布给`media`插件的共享 cache 中，不得仅存储在当前 Pod 进程内存中。

WHEN server 端收到同一实例的`SESSION_DELETE`
THEN 应按`session_id`幂等删除该会话归属并递减该实例的实时会话数量。
AND 多个 Pod 处理同一`session_id`事件时不应重复增减。
AND 不应继续在`media_report_session`保留该会话的活跃最新投影。

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

### Requirement: 媒体插件应提供数据看板读接口

`media`插件 MUST 提供受`media:management:query`权限保护的数据看板查询接口，读取`media_report_*`最新投影表并返回`数据看板.md`中节点总览、实例列表、流列表和会话列表所需字段。

接口响应的`data`对象 MUST 与`apps/lina-plugins/media/数据看板.md`示例中的`data`对象结构保持一致；响应`data`不得额外返回文档未定义的列表包装、`total`、`pageNum`或`pageSize`字段。不返回`total`的看板列表接口 MUST 不暴露分页查询参数，服务端在数据库查询阶段使用固定`10000`条上限保护数据库。

#### Scenario: 查询节点总览树

WHEN 管理端请求节点总览接口
THEN 服务端应从`media_report_node`一次性读取有界节点集合
AND 按`parent_node_id`在内存中组装`child_nodes`
AND 响应`data`应直接为当前节点对象，不得额外包裹`nodes`数组字段
AND 响应字段应覆盖`node_id`、`node_name`、`region`、`status`、`parent_node_id`、资源指标、实时计数、`avg_delay`、`last_heartbeat`、`report_time`、`node_latency_map`和`child_nodes`。

#### Scenario: 查询实例列表

WHEN 管理端按`node_id`、状态或关键词查询实例列表
THEN 服务端应在数据库查询阶段过滤、排序并最多读取`10000`条`media_report_instance`
AND 返回`node_info`和`instance_list`
AND 响应`data`不得返回`total`
AND 请求参数不得包含`pageNum`或`pageSize`
AND 不应通过逐实例回查节点配置表装配展示字段。

#### Scenario: 查询流列表

WHEN 管理端按`source_type + source_id`、`tenant_id`、`node_id`、`instance_id`、状态或关键词查询流列表
THEN 服务端应在数据库查询阶段过滤、排序并最多读取`10000`条`media_report_stream`
AND 返回`source_type`、`source_id`和`stream_list`
AND 响应`data`不得返回`total`
AND 请求参数不得包含`pageNum`或`pageSize`
AND `protocol_summary`应直接由`media_report_stream.protocol_summary`解析得到
AND `protocol_count`和`total_sessions_lifetime`应优先由`protocol_summary`聚合计算，缺少协议摘要时才回退使用读模型字段
AND `current_active_sessions`和`protocol_summary[].current_sessions`应由当前`media_report_session`活跃会话投影聚合计算。

#### Scenario: 查询会话列表

WHEN 管理端按`stream_id`查询会话列表
THEN 服务端应读取`media_report_stream`获取`stream_info`和协议摘要
AND 在数据库查询阶段按`stream_id`、`tenant_id`、`protocol_type`、`node_id`或`instance_id`过滤、排序并最多读取`10000`条`media_report_session`
AND 按`protocol_type`聚合当前查询范围内的会话数并组装`protocol_list`
AND 响应`data`不得返回`total`
AND 请求参数不得包含`pageNum`或`pageSize`
AND `link_hops`应直接由`media_report_session.link_hops`解析得到
AND `total_link_latency`应优先由`link_hops[].latency_ms`求和，缺少链路跳点时才回退使用读模型字段。

#### Scenario: 看板接口数据权限与性能边界

WHEN 查询看板流或会话数据时携带`tenant_id`筛选
THEN 服务端 MUST 在数据库查询条件中注入该筛选，不得先查询全量再在内存中过滤。

WHEN `media`插件以`platform_only`模式运行
THEN 看板接口的数据权限边界由宿主认证权限和业务筛选条件共同约束
AND 不额外引入租户安装态隔离。

WHEN 任一看板列表接口返回动态结果集
THEN 接口 MUST 提供数量上限
AND 后端查询次数不得随返回行数线性增长。
