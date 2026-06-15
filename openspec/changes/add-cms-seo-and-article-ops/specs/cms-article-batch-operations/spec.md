# cms-article-batch-operations 规范增量

## ADDED Requirements

### Requirement: 文章批量状态变更接口
CMS 插件 SHALL 提供`PUT /api/v1/cms/articles`批量状态变更接口，请求包含文章 ID 列表和目标状态（草稿或已发布），权限为`cms:article:edit`。ID 列表 MUST 非空且单次不超过 100 个，超限 MUST 拒绝。任一目标文章不存在时 MUST 整体拒绝且不产生任何变更。批量发布时，发布时间为空的文章 MUST 写入当前时刻，已有发布时间的文章 MUST 保留原值；批量下线 MUST 保留发布时间。实现 MUST 使用集合化语句在事务内完成，数据库访问次数不得随 ID 数量线性增长。

#### Scenario: 批量发布草稿
- **WHEN** 管理员选中多篇草稿文章执行批量发布
- **THEN** 所有目标文章状态变为已发布，无发布时间的写入当前时刻，公开站点可见

#### Scenario: 批量下线
- **WHEN** 管理员选中多篇已发布文章执行批量下线
- **THEN** 所有目标文章状态变为草稿、发布时间保留，公开站点不再展示

#### Scenario: 含不存在 ID 时整体拒绝
- **WHEN** 批量请求的 ID 列表中包含不存在的文章 ID
- **THEN** 接口返回文章不存在错误，所有目标文章状态保持不变

### Requirement: 文章批量删除接口
CMS 插件 SHALL 提供`DELETE /api/v1/cms/articles`批量删除接口，请求包含文章 ID 列表，权限为`cms:article:remove`。ID 列表 MUST 非空且单次不超过 100 个。任一目标文章不存在时 MUST 整体拒绝。删除 MUST 复用既有软删除语义并以单条集合化语句完成。

#### Scenario: 批量删除文章
- **WHEN** 管理员选中多篇文章执行批量删除并确认
- **THEN** 所有目标文章被软删除，管理列表与公开站点均不再返回

#### Scenario: 空列表被拒绝
- **WHEN** 批量删除请求的 ID 列表为空
- **THEN** 接口返回参数校验错误，不执行删除

### Requirement: 管理列表批量操作入口
管理端文章列表 SHALL 提供多选能力与“批量发布、批量下线、批量删除”操作入口。批量按钮 MUST 按`cms:article:edit`与`cms:article:remove`权限显隐，未选中任何行时不可执行，批量删除 MUST 有用户确认步骤。

#### Scenario: 多选后执行批量操作
- **WHEN** 管理员在文章列表勾选多行并点击批量发布
- **THEN** 接口按所选 ID 调用批量状态变更，成功后列表刷新且展示成功提示

#### Scenario: 无权限时隐藏入口
- **WHEN** 当前用户缺少`cms:article:remove`权限
- **THEN** 批量删除按钮不渲染，其余有权限的批量按钮不受影响
