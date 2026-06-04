# niu-poster-image Specification

## Purpose
TBD - created by archiving change niu-delivery-finalize. Update Purpose after archive.
## Requirements
### Requirement: 激活海报服务端生成并返回 PNG
系统 SHALL 在玩家请求激活海报时,服务端生成个性化海报 PNG 并随响应返回(base64),不同玩家的海报在视觉上可区分;渲染 MUST 通过可替换的渲染 seam 提供,使设计级排版可在不改契约的前提下替换默认实现。海报仍 MUST 仅对已激活目标牛的本人返回。

#### Scenario: 已激活玩家获取海报含图片
- **WHEN** 已激活某牛的玩家请求该牛的激活海报
- **THEN** 系统返回非空的海报 PNG(base64)与结构化字段

#### Scenario: 不同玩家海报可区分
- **WHEN** 两位不同玩家请求各自的激活海报
- **THEN** 两张海报 PNG 不相同

#### Scenario: 未激活拒绝
- **WHEN** 玩家请求自己尚未激活的牛的海报
- **THEN** 系统返回业务错误且不返回图片

