# niu-miniapp-link Specification

## Purpose
TBD - created by archiving change niu-delivery-finalize. Update Purpose after archive.
## Requirements
### Requirement: H5 回跳小程序 URL 配置化
系统 SHALL 通过插件配置提供 H5 数字纪念墙的回跳小程序 URL,并经一个公开免登录端点返回该 URL,供 H5 读取作为回跳目标;配置为空时 H5 SHALL 隐藏回跳入口。该端点 MUST NOT 暴露任何隐私字段。

#### Scenario: 公开读取回跳 URL
- **WHEN** 访客或 H5 请求公开纪念墙配置端点
- **THEN** 系统返回配置的回跳小程序 URL(未配置时为空)

#### Scenario: 未配置时隐藏回跳
- **WHEN** 回跳 URL 未配置(为空)
- **THEN** H5 不展示回跳小程序入口

