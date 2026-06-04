## ADDED Requirements

### Requirement: 日活统计
系统 SHALL 向运营提供管理端鉴权的日活(DAU)统计:返回最近 N 天(N 有上限)每日去重活跃玩家数,活跃定义为当日在任一行为(激活/喂草/签到/偷草/送草)中出现的玩家。聚合 MUST 在数据库侧完成,禁止 N+1。

#### Scenario: 运营查看日活序列
- **WHEN** 具备 `sicau-niu:settlement:view` 权限的运营请求活跃度统计
- **THEN** 系统返回最近 N 天每日去重活跃玩家数的序列

#### Scenario: 无行为数据时日活为零
- **WHEN** 某日无任何玩家行为
- **THEN** 该日的活跃玩家数为 0

### Requirement: 留存统计
系统 SHALL 提供管理端鉴权的留存统计:按注册日 cohort 计算整体次日留存与 7 日留存,各返回 cohort 人数、回访人数与留存率。仅统计留存窗口已过的 cohort(注册日 + 档位 ≤ 今日),避免低估。聚合 MUST 在数据库侧完成。

#### Scenario: 运营查看次日与 7 日留存
- **WHEN** 运营请求活跃度统计
- **THEN** 系统返回次日留存与 7 日留存的 cohort 人数、回访人数与留存率

#### Scenario: 排除窗口未过的 cohort
- **WHEN** 某玩家注册日加档位天数尚未到达今日
- **THEN** 该玩家不计入对应档位的留存 cohort
