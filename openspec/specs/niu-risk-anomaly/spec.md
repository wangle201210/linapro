# niu-risk-anomaly Specification

## Purpose
TBD - created by archiving change niu-risk-anomaly. Update Purpose after archive.
## Requirements
### Requirement: 异常行为告警视图
系统 SHALL 提供管理端鉴权的异常行为告警视图:集合化列出单玩家单日喂草次数或偷草次数超过可配置阈值的记录,含玩家、昵称、行为类型、日期、当日次数与命中阈值,供人工研判。告警为只读派生,不自动处置;阈值 MUST 可通过插件配置调整并具备合理默认;结果 MUST 有界且在数据库侧聚合,禁止 N+1。

#### Scenario: 喂草超阈值触发告警
- **WHEN** 某玩家某一天的喂草次数超过配置的喂草日阈值
- **THEN** 告警视图包含该玩家该日的喂草异常记录(含当日次数与阈值)

#### Scenario: 偷草超阈值触发告警
- **WHEN** 某玩家某一天的偷草次数超过配置的偷草日阈值
- **THEN** 告警视图包含该玩家该日的偷草异常记录

#### Scenario: 未超阈值不告警
- **WHEN** 某玩家任一天的喂草与偷草次数均未超过对应阈值
- **THEN** 告警视图不包含该玩家的异常记录

#### Scenario: 阈值随配置调整
- **WHEN** 运营调整插件中的喂草或偷草日阈值
- **THEN** 告警判定按新的阈值生效

