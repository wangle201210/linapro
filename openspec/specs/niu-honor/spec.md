# niu-honor Specification

## Purpose
TBD - created by archiving change niu-ranking-honor. Update Purpose after archive.
## Requirements
### Requirement: 荣誉定义运营维护
系统 SHALL 在运营后台提供荣誉定义的增删改查,受宿主统一权限校验。荣誉含类型(徽章/头像框/证书)、唯一编码、名称、解锁规则类型(参与即得/喂草次数/激活数/集齐分类/集齐全套)、阈值、图片、排序。类型与解锁规则类型为后端命名常量约束的稳定枚举。响应时间点字段 MUST 返回 Unix 毫秒。

#### Scenario: 新增荣誉定义
- **WHEN** 运营携带 `sicau-niu:honor:create` 权限提交未占用编码与合法类型/解锁规则
- **THEN** 系统创建该荣誉定义

#### Scenario: 编码重复或枚举非法
- **WHEN** 提交的荣誉编码已存在,或类型/解锁规则不在允许枚举内
- **THEN** 系统返回 `bizerr`,不创建

#### Scenario: 无权限访问被拒绝
- **WHEN** 调用方不具备对应 `sicau-niu:honor:*` 权限
- **THEN** 系统由统一权限中间件拒绝访问

### Requirement: 玩家荣誉解锁展示
系统 SHALL 向玩家返回各荣誉的解锁状态(只读),依据其喂草次数、激活数与图鉴集齐情况按荣誉解锁规则计算;玩家只能看到自己的荣誉状态。解锁计数用批量聚合,不逐荣誉重复查询。荣誉的持久发放不在本能力范围内。

#### Scenario: 参与类荣誉默认解锁
- **WHEN** 玩家请求荣誉列表
- **THEN** "参与即得"类荣誉显示为已解锁

#### Scenario: 阈值类荣誉按计数解锁
- **WHEN** 玩家喂草次数达到某喂草次数类荣誉的阈值
- **THEN** 该荣誉显示为已解锁,未达阈值的显示未解锁

#### Scenario: 集齐类荣誉按图鉴解锁
- **WHEN** 玩家集齐某分类的全部已上线主卡
- **THEN** 该分类集齐类荣誉显示为已解锁

#### Scenario: 只读不产生持久发放
- **WHEN** 玩家请求荣誉列表
- **THEN** 系统只返回计算的解锁状态,不写入持久发放记录

