## ADDED Requirements

### Requirement: 按策略 ID 反查租户策略绑定

`media`插件 SHALL 提供只读接口，允许调用方按媒体策略 ID 分页查询绑定到该策略的租户策略绑定记录。

#### Scenario: 查询策略绑定租户

- **WHEN** 调用方使用有效策略 ID 查询租户策略绑定
- **THEN** 系统应只从`media_strategy_tenant`中返回`strategy_id`匹配的租户绑定记录
- **AND** 响应应包含租户 ID、策略 ID 和稳定行键
- **AND** 响应不得包含`media_strategy_device_tenant`中的租户设备策略绑定记录

#### Scenario: 分页查询策略绑定租户

- **WHEN** 匹配同一策略 ID 的租户绑定记录超过请求分页大小
- **THEN** 系统应在数据库查询阶段完成策略过滤、排序和分页
- **AND** 响应总数应表示该策略在`media_strategy_tenant`中的匹配记录数

### Requirement: 按策略 ID 反查设备策略绑定

`media`插件 SHALL 提供只读接口，允许调用方按媒体策略 ID 分页查询绑定到该策略的设备策略绑定记录。

#### Scenario: 查询策略绑定设备

- **WHEN** 调用方使用有效策略 ID 查询设备策略绑定
- **THEN** 系统应只从`media_strategy_device`中返回`strategy_id`匹配的设备绑定记录
- **AND** 响应应包含设备 ID、策略 ID 和稳定行键
- **AND** 响应不得包含`media_strategy_device_tenant`中的租户设备策略绑定记录

#### Scenario: 分页查询策略绑定设备

- **WHEN** 匹配同一策略 ID 的设备绑定记录超过请求分页大小
- **THEN** 系统应在数据库查询阶段完成策略过滤、排序和分页
- **AND** 响应总数应表示该策略在`media_strategy_device`中的匹配记录数

### Requirement: 媒体管理列表分页上限

`media`插件管理端列表接口 SHALL 接受最大`10000`的`pageSize`请求值，并在服务层统一将超过该值的分页大小收敛到`10000`。

#### Scenario: 请求大分页列表

- **WHEN** 调用方请求`pageSize=10000`的媒体管理端列表接口
- **THEN** API DTO 校验应允许该分页大小
- **AND** 服务层应在数据库查询阶段执行过滤、排序和分页，不得先加载全量数据后分页
- **AND** 超过`10000`的服务层分页输入应被收敛到`10000`

### Requirement: Token 策略查询应按节点执行租户流数量限制

`mediaopen`公共策略查询接口 SHALL 要求调用方携带节点 ID，并在返回策略前按租户和节点检查未关闭会话数量是否达到租户流配置限制。

#### Scenario: 节点未配置租户流数量限制

- **WHEN** 调用方携带有效 token、设备 ID 和节点 ID 查询策略
- **AND** `media_tenant_stream_config`中不存在该租户和节点的启用限制配置
- **THEN** 系统应按既有优先级返回策略
- **AND** 不应因为其他节点的流数量限制拒绝本次策略查询

#### Scenario: 节点未关闭会话数量达到限制

- **WHEN** 调用方携带有效 token、设备 ID 和节点 ID 查询策略
- **AND** `media_tenant_stream_config`中存在该租户和节点的启用限制配置
- **AND** `media_report_session`中该租户和节点`close_time`为空的会话数量大于或等于`max_concurrent`
- **THEN** 系统应返回稳定业务错误码，告知租户在该节点的流数量已经超过限制
- **AND** 不应返回媒体策略内容

#### Scenario: 已关闭会话不计入限制

- **WHEN** 该租户和节点存在历史会话
- **AND** 历史会话的`close_time`不为空
- **THEN** 系统统计限制时不得将这些已关闭会话计入当前流数量。
