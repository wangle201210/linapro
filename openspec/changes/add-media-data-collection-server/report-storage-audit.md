# 四类采集上报存储核对

本次按`net-flux v0.0.21`的`proto/report.proto`核对`MachineMetric`的`26`个字段、`NetworkMetric`的`10`个字段、`StreamMetric`的`28`个字段和`SessionMetric`的`25`个字段，共`89`个顶层字段，另包含`LinkHop`的`6`个嵌套字段。

用户选择保留每个对象的完整最新上报。本轮增加的`report_payload`以输入`protobuf`消息为源，完整保存这些已定义字段；下列原有投影仍供看板查询，不能单独当作原始报文。

## 原有字段遗漏

| 上报 | 原投影未完整保存的内容 |
| --- | --- |
| `MachineMetric` | `disk_used`、`disk_total`、`disk_read_time`、`disk_write_time`、`disk_read_count`、`disk_write_count`及完整`extra`。 |
| `NetworkMetric` | 完整`extra`，包括探测端当前提供的`url`、`status_code`、`success`、`total`、`error`等键；原有`reachable`只是解析后的布尔值。 |
| `StreamMetric` | `machine_type`、`video_codec`、`audio_codec`及完整`extra`。 |
| `SessionMetric` | 完整`extra`；`link_hops`中的`hop_id`、`hop_name`、`hop_type`、`node_name`。 |

## 原有有损转换

| 范围 | 转换及无法还原的原始内容 |
| --- | --- |
| 机器、流、会话标识 | `machine_id`和`instance_id`合并选择单个实例键；字符串去首尾空白，名称和节点可能从其他字段或`extra`回退。 |
| 机器资源 | `cpu_usage`按核数转已用核数；内存字节转`MB`，网络`bps`转`KB/s`，投影的定点数会舍入。 |
| 网络时间 | 链路`timestamp`统一为毫秒；原始非正值被接收时间替代。 |
| 流枚举与媒体属性 | `status/protocol`转换为字符串，未知值合并为默认类型；`width/height`合成分辨率；原始`stream_path/stream_alias`参与地址与名称回退。 |
| 流统计 | 当前会话数、时长、水印可能由默认值、其他字段或`extra`推导，部分小数投影舍入。 |
| 会话属性 | `client_type`文本归并为整数；协议数值转字符串，宽高合成分辨率，时长和总延迟可能派生。空`client_ip`不会清空已有投影值。 |
| 会话跳点 | 仅保存经过整理的节点 ID、延迟和新生成的序号，原始节点名称等内容缺失。 |
| 时间字段 | 部分原始整数时间转 SQL 时间，原始零值及单位可能丢失。 |

完整`report_payload`保留转换之前的原始 ID、数值、单位、空值语义、枚举值、`extra`和嵌套字段。`int64`按官方`protobuf JSON`规则保存为字符串，避免超过`2^53`时丢失精度；原始枚举以数字保存，包含未知数值。数据库 JSON 的键顺序不属于报文语义。

二次复核发现 PostgreSQL `JSONB`拒绝字符串中的转义`U+0000`，因此`report_payload`改为`JSON`列，并为四类报文的完整扩展内容增加该字符的真实数据库往返测试。字段名和其他存储契约不变。

## 生命周期核对

原有`STREAM_DELETE`和`SESSION_DELETE`跳过完整 upsert，仅处理计数和关闭状态，附带字段不会写入；首次收到删除事件且资源行不存在时，更新影响零行。

本轮在删除路径先保存完整报文和`report_subcommand`。已有记录只更新原始报文元数据，避免稀疏删除清空已有名称、租户等展示字段；首次删除创建已关闭记录。后续计数和关闭写入不覆盖原始内容。计数操作更新实例投影时不会清空最近机器上报，机器原始时间始终保存在`report_payload.timestamp`。

## 保存位置与边界

| 原始报文 | 保存位置 |
| --- | --- |
| `MachineMetric` | `media_report_instance.report_payload` |
| `NetworkMetric` | `plugin_media_report_network_link.report_payload` |
| `StreamMetric` | `media_report_stream.report_payload`及`report_subcommand` |
| `SessionMetric` | `media_report_session.report_payload`及`report_subcommand` |

- 原有业务键与有效性边界继续适用：缺实例标识、流 ID、会话 ID，或缺有效网络链路键的报文不作为完整对象归档。
- 完整性保证以符合既有业务投影约束且成功写入的对象为范围；主键长度、`INET`格式、数值范围等约束仍执行，违反约束会返回数据库错误，不承诺任意合法 protobuf 值都能绕过这些约束。原始 JSON 不能使失败的业务投影自动变成成功。
- 每个现有业务键只保留最新有效内容；机器、流和会话按到达顺序覆盖，网络仍按原采样时间条件拒绝旧样本。多 URL 共用同一个链路键时仍覆盖为一份内容。
- 关闭流和会话超过`30`天仍按既有逻辑清理，连同原始内容一起删除。
- `protobuf JSON`覆盖当前协议已定义字段，不保留未知的二进制 wire 字段，不承诺逐字节还原原始 TCP 帧。
- 旧数据库行的`report_payload`保持`NULL`，需要新上报才能填充，此前丢弃的信息无法回填。
- 不新增原始内容的 HTTP 查询入口，避免将`extra`中的潜在敏感业务内容直接暴露；看板列表排除原始 JSON，避免最多`10000`行查询加载无用大字段。

## 独立发现

上游机器采集器从`disk.IOCounters()`取得累计磁盘读写字节，而原有投影直接除以`1024`并标注为速率，没有通过相邻采样差值除以时间。本轮原始 JSON 已保留累计字节及相关耗时、次数，但该展示单位问题需单独处理，不属于完整原始存储的实现结论。

## 验证方式

完整性测试通过协议 descriptor 核对各类型字段覆盖，构造每个字段均有内容的报文，并在序列化与数据库读取后执行`protojson.Unmarshal`及`proto.Equal`比对。额外验证默认零值、未知枚举数值、大整数、完整扩展键和全部跳点字段，以及删除事件、首次删除、计数更新和旧链路报文边界。最终命令执行结果记录在同目录`tasks.md`的`FB-41`记录中。
