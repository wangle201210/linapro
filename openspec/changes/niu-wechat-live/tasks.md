## 1. 真实微信网关

- [x] 1.1 读取命中规则文件:`backend-go.md`、`cache-consistency.md`、`testing.md`、`api-contract.md`、`documentation.md`、`i18n.md`(逐一记录影响/无影响)。
- [x] 1.2 `wechat_real.go`:`realGateway` 增 `httpClient`/`sync.Mutex`/`token`/`expiry` 字段;微信基址常量。
- [x] 1.3 真实 `Code2Session`:`g.Client()` GET `sns/jscode2session`,解析 openid,errcode→`CodeWeChatCodeInvalid`。
- [x] 1.4 真实 `DecodePhone`(新版 code 流程):`accessToken()`(缓存+刷新)+ POST `getuserphonenumber{code}`,解析 `phone_info.phoneNumber`,errcode→`CodeWeChatPhoneDecodeFailed`。
- [x] 1.5 纯解析函数 `parseSessionResponse/parseTokenResponse/parsePhoneResponse`(gjson),HTTP 取文为薄封装。

## 2. 测试与验证

- [x] 2.1 `wechat_real_test.go`:解析函数单测(成功、errcode、空/坏响应),无网络、无凭据。
- [x] 2.2 编译门禁:`temp/go.work.sicau-verify` 下 `go build ./...` + `go vet ./...`。
- [x] 2.3 `openspec validate niu-wechat-live --strict` 通过。
- [x] 2.4 `lina-review` 通过(记录 access_token 进程内缓存影响/契约不变/i18n 无影响/真实往返为 CI 外部依赖)。

> 真实微信往返需真实 AppID/Secret,不在 CI 执行;运维填凭据 + `wechat.mock=false` 即切真实流程。
