## 1. 证书渲染 seam

- [x] 1.1 读取命中规则文件:`backend-go.md`、`plugin.md`、`api-contract.md`、`data-permission.md`、`architecture.md`、`testing.md`、`documentation.md`、`cache-consistency.md`、`i18n.md`(逐一记录影响/无影响)。
- [x] 1.2 新增 `internal/service/honor/internal/certrender`:`CertRenderer` 接口 + `CertData` 值对象 + 默认 `basicRenderer`(框图 + 强调色由 hash(nickname+honorCode) 派生,实际返回 PNG)+ `CodeCertRenderFailed`。

## 2. honor 服务与装配

- [x] 2.1 `honor.go`:新增 `certRenderer`/`campusBadge` 字段,`New(certRenderer, Config{CampusBadge})`;re-export `CertRenderer` 与 `NewBasicCertRenderer`。
- [x] 2.2 `honor_certificate.go`:`PlayerCertificate(ctx, playerID, honorID)` —— 证书类型校验 + `user_honor` 持有校验 + 字段装配 + 渲染,返回 base64 + 字段;新增持有/非证书错误码。
- [x] 2.3 `plugin.go`:`honorsvc.New(honorsvc.NewBasicCertRenderer(), honorsvc.Config{CampusBadge: ...})`(复用 poster.campusBadge);玩家子组绑定证书接口。

## 3. 接口层

- [x] 3.1 `api/player/v1` 证书 DTO(`GET .../player/certificates`,入参 honorId,响应 base64 PNG + 昵称/荣誉名/编码/获得时间 Unix 毫秒);`api/player/player.go` 接口追加方法。
- [x] 3.2 `internal/controller/player` 证书处理方法委派 honor 服务。

## 4. 测试与验证

- [x] 4.1 服务层 DB 门控单测:持有者生成非空 PNG、不同持有者 PNG 不同、未持有拒绝、非证书类型拒绝;沿用共享临时库 + TRUNCATE + TestMain。
- [x] 4.2 编译门禁:`temp/go.work.sicau-verify` 下 `go build ./...` + `go vet ./...`。
- [x] 4.3 `openspec validate niu-certificate-render --strict` 通过。
- [x] 4.4 `lina-review` 通过(记录持有校验数据权限/无新增表/缓存无影响/i18n 无影响)。

> 小程序展示 UI 不在本期;玩家接口冒烟并入最终统一阶段。
