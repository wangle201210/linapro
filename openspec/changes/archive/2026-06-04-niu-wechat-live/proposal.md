## Why

需求 M1 要求「微信小程序端用户登录 + 手机号一键授权登录」。C1 交付了可 mock 的微信网关 seam,但真实实现是 TODO,`Code2Session`/`DecodePhone` 直接返回 `CodeWeChatNotImplemented`。本变更**实现真实生产流程**:真实 HTTP 调用、响应解析、错误映射全部做实,**唯一外部依赖是 AppID/Secret(走配置)**;替换真实凭据并关闭 mock 后即可直接跑。mock 网关保留供开发与 E2E。

## What Changes

- **真实 `Code2Session`**:调用微信 `sns/jscode2session`(appid/secret/js_code),解析 `openid`,把 `errcode`(40029 无效 code、45011 限频等)映射为 `CodeWeChatCodeInvalid`。
- **真实 `DecodePhone`(新版 code 流程)**:获取并缓存 `access_token`(`cgi-bin/token`),调用 `wxa/business/getuserphonenumber`(手机号 code),解析 `phone_info.phoneNumber`,`errcode` 映射为 `CodeWeChatPhoneDecodeFailed`。新版 code 流程自包含、不依赖 session_key,是微信推荐的正式流程。
- **`access_token` 进程内缓存**(互斥保护 + 过期前刷新),避免每次调用都换取、触发微信全局限频。
- 解析逻辑抽为**纯函数**(parseSession/parseToken/parsePhone),可在无网络/无凭据下单测;HTTP 调用为薄封装。

不在本 change 范围:正式接入(不提供真实 AppID/Secret,继续默认 mock=true);旧版 encryptedData+IV(AES-CBC)解密(需 session_key 持久化,新版 code 流程已覆盖主路径)。

## Capabilities

### New Capabilities
- `niu-wechat-live`: 真实微信登录与手机号授权的生产 HTTP 流程(凭据走配置,mock 保留)。

### Modified Capabilities
<!-- 无需求级行为变更:登录/绑定手机号的对外契约不变,仅把真实网关从 TODO 占位换成可运行实现。 -->

## Impact

- **后端**:`apps/lina-plugins/sicau-niu/backend/internal/service/wechat/`
  - `wechat_real.go`:真实 `Code2Session`/`DecodePhone` + `access_token` 缓存;`g.Client()` HTTP。
  - 新增纯解析函数 + 单测(无网络);`wechat.go` 契约不变。
- **配置**:沿用既有 `wechat.appId`/`wechat.secret`/`wechat.mock`;开发默认 `mock=true`,生产填真实凭据并设 `mock=false`。
- **数据权限**:登录为公开端点(不变);手机号仅本人绑定(不变)。
- **缓存一致性**:新增**进程内 `access_token` 缓存**(互斥 + TTL 提前刷新);单实例语义,非跨节点共享状态,失效由 TTL 驱动 —— 记录该缓存影响。
- **i18n**:错误码已存在,无新增用户可见文案 —— 无影响。
- **验证**:纯解析函数单测(成功 + errcode + 空响应)+ 编译门禁 + `openspec validate --strict`;真实微信往返需真凭据,记录为外部依赖不可在 CI 执行。
