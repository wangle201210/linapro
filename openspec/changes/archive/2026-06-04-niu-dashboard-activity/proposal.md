## Why

需求 M5 后台看板明确要求「激活数、**日活用户、留存**等基础统计数据」。C7 结算看板已覆盖激活数、玩家数、喂草等总量,但**日活(DAU)与留存尚未实现**。本变更补齐这两项基础活跃度指标,使 M5 看板完整。

## What Changes

- 新增**活跃度看板 API**(管理端,鉴权):
  - **日活(DAU)**:返回最近 N 天每日去重活跃玩家数(有任意行为=激活/喂草/签到/偷草/送草之一)。
  - **留存**:按注册日 cohort 计算**次日留存**与**7 日留存**的整体(cohort 人数 / 回访人数 / 留存率)。
- 全部在数据库侧集合化聚合,有界(N 天上限),杜绝 N+1;**不新增数据表、不新增定时任务**(按需从既有行为表的 `created_at` 派生)。

不在本 change 范围:细分到每个 cohort 的留存曲线明细;7 日以外更多留存档(活动仅 2 个月)。

## Capabilities

### New Capabilities
- `niu-dashboard-activity`: 运营看板的日活(DAU 序列)与留存(次日/7 日)活跃度指标 API。

### Modified Capabilities
<!-- 无需求级行为变更:复用 C7 结算的管理端鉴权与 C1-C4 行为表,仅新增活跃度聚合读取路径。 -->

## Impact

- **后端**:`apps/lina-plugins/sicau-niu/backend/`
  - 新增 `api/settlement/v1` 活跃度 DTO、`internal/controller/settlement` 处理方法、`internal/service/settlement` 活跃度聚合。
  - `plugin.go`:管理端子组绑定活跃度接口(权限 `sicau-niu:settlement:view`)。
  - 跨行为表(activation/feeding/checkin/steal/gift)的去重活跃集合用**参数化集合查询**派生,引用 DAO 表/列常量防漂移,杜绝 N+1。
- **数据库**:无新增表、无迁移;派生自既有 `created_at`。
- **前端**:运营结算页新增「活跃度」区(DAU 简表 + 次日/7 日留存卡片),data-testid 锚点,中文,无 i18n。
- **数据权限**:管理端鉴权(`Auth+Tenancy+Permission` + `permission` 标签);只读聚合,有界。
- **缓存一致性**:实时按需聚合,无新增缓存/失效路径 —— 无影响。
- **i18n**:单语言中文 —— 无影响。
- **验证**:服务层 DB 门控单测(DAU 计数、次日/7 日留存口径)+ 编译门禁 + `openspec validate --strict`;管理端 E2E 并入最终统一阶段。
