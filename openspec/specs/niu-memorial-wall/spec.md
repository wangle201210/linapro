# niu-memorial-wall Specification

## Purpose
TBD - created by archiving change niu-h5-wall. Update Purpose after archive.
## Requirements
### Requirement: 公开首发纪念墙
系统 SHALL 提供公开免登录的首发纪念墙数据接口:返回首发激活(`is_first`)的玩家名单(昵称、身份标签、所属牛名称/序号、到场顺序),按到场顺序/牛序号有界返回(上限 120)。响应 MUST NOT 包含手机号、openid、设备指纹等隐私字段。响应时间点字段 MUST 返回 Unix 毫秒。批量装配,禁止 N+1。

#### Scenario: 访客无需登录查看首发墙
- **WHEN** 访客在未携带任何登录令牌的情况下请求首发纪念墙
- **THEN** 系统返回首发激活玩家名单(昵称/身份/牛名称/序号/到场顺序),且响应不含手机号或 openid

#### Scenario: 首发名单有界
- **WHEN** 首发激活记录存在
- **THEN** 系统按到场顺序返回不超过 120 条记录

#### Scenario: 暂无首发数据时名单为空
- **WHEN** 尚无任何首发激活记录
- **THEN** 系统返回空名单

### Requirement: 公开校史精选
系统 SHALL 提供公开免登录的校史精选接口:从已上线的校史卡片与校史金句中取样返回(标题/分类/图片/文案),有数量上限。

#### Scenario: 访客查看校史精选
- **WHEN** 访客请求校史精选
- **THEN** 系统返回有上限的校史卡片与金句取样,供 H5 展示

### Requirement: 公开活动统计
系统 SHALL 提供公开免登录的活动统计接口:返回已激活牛数、参与玩家数等公开聚合,统计在数据库侧完成,不暴露隐私明细。

#### Scenario: 访客查看活动统计
- **WHEN** 访客请求活动统计
- **THEN** 系统返回已激活牛数与参与玩家数等公开聚合值

### Requirement: H5 数字纪念墙静态页
系统 SHALL 通过插件 `public_assets` 托管一个移动端 H5 数字纪念墙静态页,展示首发墙/校史精选/统计,并提供扫码或链接回跳小程序的入口。该静态页独立于管理端 Vue 页面,不进入管理端前端构建。

#### Scenario: 访客访问 H5 纪念墙
- **WHEN** 访客通过宿主托管的公开静态资源地址访问 H5 纪念墙
- **THEN** 页面加载并调用公开纪念墙 API 渲染首发墙、校史精选与统计,并展示回跳小程序入口

