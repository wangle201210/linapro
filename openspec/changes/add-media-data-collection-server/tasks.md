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
