## Why

需求 M13 风控明确要求防范「异常**喂草**或**偷草**」并提供「异常行为基础**告警**」。C7 已做一机多号(共享设备指纹)告警视图,但**异常喂草/偷草频率告警尚未实现**。本变更补齐基于阈值的异常行为告警,使 M13 风控告警完整。

## What Changes

- 新增**异常行为告警视图**(管理端,鉴权,只读):集合化列出**单玩家单日**喂草次数或偷草次数**超过阈值**的记录(玩家、昵称、类型 feed/steal、日期、当日次数、命中阈值),供人工研判。
- 阈值**走插件配置可调**(`anomaly.feedDailyThreshold`、`anomaly.stealDailyThreshold`),给合理默认;结果有界(上限),DB 侧聚合,杜绝 N+1。

不在本 change 范围:自动封禁/处置(仅告警供人工研判);异常定位告警(激活已做 LBS 校验,定位异常的细化检测后续);告警持久化为独立事件表(本期派生只读视图)。

## Capabilities

### New Capabilities
- `niu-risk-anomaly`: 单玩家单日喂草/偷草超阈值的异常行为告警派生视图(阈值可配置)。

### Modified Capabilities
<!-- 无需求级行为变更:复用 C7 结算管理端鉴权与 C4 喂草/偷草表,仅新增异常告警聚合读取路径。 -->

## Impact

- **后端**:`apps/lina-plugins/sicau-niu/backend/`
  - 新增 `api/settlement/v1` 异常告警 DTO、controller 处理方法、`internal/service/settlement` 异常聚合。
  - settlement 服务 `New()` 改为注入异常阈值配置;`plugin.go` 读取 `anomaly.*` 配置并装配,管理端子组绑定接口(权限 `sicau-niu:settlement:view`)。
  - feeding 按 `DATE(created_at)`、steal 按 `steal_date` 分组 `HAVING COUNT > 阈值`,集合化,有界,杜绝 N+1;昵称批量装配。
- **数据库**:无新增表/列;复用 `feeding`、`steal`、`user`。
- **配置**:新增 `anomaly.feedDailyThreshold`、`anomaly.stealDailyThreshold`(及结果上限),写入 `config.example.yaml` 与开发配置。
- **数据权限**:管理端鉴权 + `permission` 标签;只读派生,有界。
- **缓存一致性**:实时聚合,无新增缓存/失效路径 —— 无影响。
- **i18n**:单语言中文 —— 无影响。
- **验证**:服务层 DB 门控单测(超阈值命中、未超阈值不报、feed/steal 各自口径)+ 编译门禁 + `openspec validate --strict`;管理端 E2E 并入最终统一阶段。
