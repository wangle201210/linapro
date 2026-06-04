## Context

把 C1 的微信网关从 TODO 占位换成真实生产流程,只把 AppID/Secret 作为外部配置。复用既有 `Gateway` seam 与 `DecodePhoneInput`(已含新版 `Code`)。约束遵循 `.agents/rules/` 的 `backend-go`、`cache-consistency`、`testing`、`api-contract`;`platform_only`。

## Goals / Non-Goals

**Goals:** 真实 `jscode2session`、真实手机号(新版 code 流程)、`access_token` 缓存、错误码映射;解析逻辑可单测;mock 保留;换真凭据即可跑。

**Non-Goals:** 提供真实凭据/正式上线;旧版 encryptedData+IV AES 解密(需 session_key 持久化);跨节点共享 token 缓存。

## Decisions

### D1 真实 HTTP 流程,凭据走配置
`realGateway{appID, secret}` 用 `g.Client()` 调微信开放接口:
- 登录:`GET https://api.weixin.qq.com/sns/jscode2session?appid=&secret=&js_code=&grant_type=authorization_code`。
- 手机号:`POST https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=` body `{code}`。
- 令牌:`GET https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=&secret=`。
基址为常量;凭据来自 `wechat.appId/secret`(已有配置)。`mock=true` 时仍走 mock 网关(开发/E2E)。

### D2 新版 code 手机号流程(不依赖 session_key)
`DecodePhone` 走微信推荐的新版 `getPhoneNumber` code 流程:取 `access_token` → `getuserphonenumber(code)` → `phone_info.phoneNumber`。自包含、无需 session_key。旧版 encryptedData+IV 需 session_key,不在本期(新版覆盖主路径)。

### D3 access_token 进程内缓存
`access_token` 微信侧 2 小时有效且**全局限频**,禁止每次换取。`realGateway` 持 `sync.Mutex + token + expiry`,命中未过期则复用,否则刷新并记 `expiry = now + (expires_in - 提前量)`。单实例进程内缓存:语义简单、由 TTL 驱动失效,不引入跨节点一致性(记录于缓存影响)。

### D4 解析与 HTTP 分离,保证可测
微信响应解析抽为纯函数 `parseSessionResponse/parseTokenResponse/parsePhoneResponse(body string)`,用 `gjson` 解析并映射 `errcode`;HTTP 取文为薄封装。纯函数**无网络可单测**(成功、errcode、空/坏响应);真实往返需真凭据,记录为 CI 外部依赖。

### D5 错误映射
`jscode2session` `errcode!=0` 或 `openid` 空 → `CodeWeChatCodeInvalid`(携 errmsg)。`getuserphonenumber` `errcode!=0` 或 phone 空 → `CodeWeChatPhoneDecodeFailed`。`token` `errcode!=0` 或空 → 内部错误(包裹)。HTTP/网络错误一律包裹返回,不静默吞。

## Risks / Trade-offs

- [无真实凭据无法端到端自测] → 解析纯函数单测覆盖逻辑;真实往返记录为外部依赖。
- [access_token 限频] → 进程内缓存 + 提前刷新;多实例各自缓存(微信允许中心化 token,后续可接宿主缓存)。
- [旧版手机号流程未覆盖] → 采用微信推荐新版 code 流程;旧版需 session_key,留待后续。

## Migration Plan

1. `wechat_real.go`:真实 `Code2Session`/`DecodePhone` + token 缓存;纯解析函数。
2. `wechat_real_test.go`:解析函数单测(无网络)。
3. 编译门禁 + `openspec validate --strict`。
4. 上线:运维填 `wechat.appId/secret` 并设 `wechat.mock=false` 即切真实流程;回滚设 `mock=true`。

## Open Questions

1. 是否将 access_token 改为宿主中心化缓存(多实例共享)—— 本期进程内,后续可演进。

> 影响判断:字典 —— 无新增枚举,**无影响**;缓存一致性 —— 新增进程内 access_token 缓存(TTL 驱动、单实例,已记录);数据权限 —— 登录公开、手机号本人绑定,**不变**;跨平台 —— 仅插件内文件,**无影响**;i18n —— 无新增文案,**无影响**。
