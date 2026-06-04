## ADDED Requirements

### Requirement: 真实微信登录换取 openid
系统 SHALL 在非 mock 模式下通过真实调用微信 `jscode2session` 接口,用配置的 AppID/Secret 与登录 code 换取玩家 openid;微信返回错误码或 openid 为空时 MUST 映射为登录无效业务错误。mock 模式继续返回 mock openid 供开发与测试。

#### Scenario: 真实换取 openid 成功
- **WHEN** 非 mock 模式下提交有效登录 code
- **THEN** 系统通过微信接口换得非空 openid 并据此登录

#### Scenario: 微信返回错误码
- **WHEN** 微信 `jscode2session` 返回非零 errcode 或空 openid
- **THEN** 系统返回登录无效业务错误

### Requirement: 真实手机号授权解码
系统 SHALL 在非 mock 模式下通过微信新版 getPhoneNumber code 流程解码手机号:先获取并缓存 `access_token`,再调用 `getuserphonenumber` 用手机号 code 换取手机号;微信返回错误码或手机号为空时 MUST 映射为手机号解码失败业务错误。

#### Scenario: 真实解码手机号成功
- **WHEN** 非 mock 模式下提交有效手机号 code
- **THEN** 系统通过微信接口换得手机号并完成绑定

#### Scenario: 手机号解码失败
- **WHEN** 微信 `getuserphonenumber` 返回非零 errcode 或空手机号
- **THEN** 系统返回手机号解码失败业务错误

### Requirement: access_token 缓存
系统 SHALL 在进程内缓存微信 `access_token` 并在过期前复用,避免每次手机号解码都重新换取令牌而触发微信全局限频;令牌过期后 SHALL 自动刷新。

#### Scenario: 复用未过期令牌
- **WHEN** 已缓存的 access_token 仍在有效期内
- **THEN** 系统复用该令牌而不重新换取

#### Scenario: 过期后刷新令牌
- **WHEN** 缓存的 access_token 已过期
- **THEN** 系统重新换取并更新缓存
