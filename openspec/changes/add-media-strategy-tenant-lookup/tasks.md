## 1. 实现

- [x] 1.1 新增按策略 ID 查询租户策略绑定的 API DTO 和控制器方法。
- [x] 1.2 在`media`服务层实现只读`media_strategy_tenant`的分页查询。
- [x] 1.3 补充后端路由测试，验证接口只返回租户策略绑定。
- [x] 1.4 同步媒体管理端列表接口`pageSize`校验上限为`10000`，并覆盖服务层上限一致性验证。

## 2. 验证

- [x] 2.1 运行覆盖`media`插件后端的 Go 测试。
- [x] 2.2 运行`openspec validate add-media-strategy-tenant-lookup --strict`。

## Feedback

- [x] **FB-1**: 新增按策略 ID 查询设备策略绑定列表接口，并验证只读取设备策略表。
