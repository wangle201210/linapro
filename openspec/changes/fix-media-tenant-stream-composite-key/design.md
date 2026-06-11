# 设计说明

## 资源标识

租户流配置的业务唯一键调整为`tenantId + nodeNum`。详情、更新和删除接口使用嵌套资源路径`/media/tenant-stream-configs/{tenantId}/nodes/{nodeNum}`定位记录，避免同一租户在多个媒体节点上配置流并发时发生误删或误改。

创建接口仍使用集合路径`/media/tenant-stream-configs`，由请求体提供`tenantId`和`nodeNum`。更新接口保留原资源键在路径中，允许请求体同时修改租户 ID 或节点编号；当新复合键已存在时拒绝更新。

## 数据和性能

`media_tenant_stream_config`表主键调整为`("tenant_id", "node_num")`，并保留租户 ID 查询索引，支持列表筛选和单条定位。列表查询继续在数据库侧完成过滤、排序和分页，节点名称通过节点编号集合批量装配，避免逐行查询。

已安装源码插件只有在检测到版本升级时才会执行新增 SQL 迁移，因此`media/plugin.yaml`版本同步提升到`v0.1.1`，确保现有`v0.1.0`安装在插件同步和升级流程中应用复合键迁移。

## 边界和影响

本变更仅修改`media`源码插件，不改变`apps/lina-core`宿主核心契约，不新增运行期依赖、缓存、模块启停逻辑或跨模块调用。权限标签沿用现有`media:management:*`管理权限；`media`插件未启用插件多语言资源，因此不新增`manifest/i18n`或`apidoc`翻译文件。
