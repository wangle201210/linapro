# 新增媒体数据采集 Server

## 背景

`media`插件当前已经提供媒体策略、节点、流别名、白名单和公开查询接口，但缺少面向采集客户端的数据上报入口。客户端会使用`net-flux`的 client 端协议向 server 端上报机器、网络、实例、流和会话指标，因此插件需要在宿主启动后提供兼容的 TCP server 端能力。

## 范围

本变更在`media`源码插件内部集成`net-flux` server 端，参考`examples/disco/server/server.go`实现协议处理。采集 server 通过插件配置显式启用，支持配置监听地址，并在宿主`system.started`钩子中异步启动。server 端应响应系统心跳，接受机器、网络、实例、流和会话指标上报，将可映射的上报数据写入`media`插件既有数据看板读模型，并在启用 discovery 配置后对接 Nacos 完成实例注册、注销和查询。

## 非目标

- 不新增查询接口或前端展示；本次持久化范围限定为`media`插件已有数据看板设计中的上报读模型。
- 不把采集 server 能力上移到`lina-core`宿主契约。
- 不新增`net-flux` config、event 或 control 命令对应的 LinaPro 业务动作。
