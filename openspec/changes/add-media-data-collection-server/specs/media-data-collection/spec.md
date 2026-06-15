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

#### Scenario: 根据生命周期事件维护实时计数

WHEN server 端收到同一实例的`STREAM_ADD`
THEN 应按`stream_id`幂等记录该流归属并递增该实例的实时直播流数量。
AND 该归属和计数 MUST 存储在宿主发布给`media`插件的共享 cache 中，不得仅存储在当前 Pod 进程内存中。

WHEN server 端收到同一实例的`STREAM_DELETE`
THEN 应按`stream_id`幂等删除该流归属并递减该实例的实时直播流数量。
AND 多个 Pod 处理同一`stream_id`事件时不应重复增减。

WHEN server 端收到同一实例的`SESSION_ADD`
THEN 应按`session_id`幂等记录该会话归属并递增该实例的实时会话数量。
AND 该归属和计数 MUST 存储在宿主发布给`media`插件的共享 cache 中，不得仅存储在当前 Pod 进程内存中。

WHEN server 端收到同一实例的`SESSION_DELETE`
THEN 应按`session_id`幂等删除该会话归属并递减该实例的实时会话数量。
AND 多个 Pod 处理同一`session_id`事件时不应重复增减。

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
