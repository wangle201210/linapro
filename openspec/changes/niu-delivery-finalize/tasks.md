## 1. 海报真生成 PNG

- [x] 1.1 读取命中规则文件:`backend-go.md`、`api-contract.md`、`frontend-ui.md`、`data-permission.md`、`plugin.md`、`testing.md`、`documentation.md`、`cache-consistency.md`、`i18n.md`(逐一记录影响/无影响)。
- [x] 1.2 增强 `posterrender` 默认渲染器:竖版海报框图 + 由 hash(nickname+niuCode) 派生强调色;实际返回 PNG(避免 `_,_=Write` 用内联 FNV-1a)。
- [x] 1.3 `PosterOutput` 增 `ImageBase64`;`activation_poster.go` 把渲染结果 base64 返回(不再丢弃)。
- [x] 1.4 `PosterRes` DTO 增 `imageBase64`;`player_v1_poster.go` 透传。

## 2. 回跳 URL 配置化

- [x] 2.1 wall 服务 `New(Config{MiniappURL})`;新增 `WallConfig` 方法返回 `{MiniappURL}`。
- [x] 2.2 `api/wall/v1` 新增 `ConfigReq/Res{MiniappURL}`(公开,无 permission);`api/wall/wall.go` 接口追加;controller `Config` 处理方法。
- [x] 2.3 `plugin.go`:读 `miniapp.url` 配置,装配 `wallsvc.New(Config{...})`;公开子组绑定 `wall/config`;`config.example.yaml` 增 `miniapp.url`。
- [x] 2.4 H5 `frontend/wall/index.html`:`fetch wall/config` 取回跳 URL,替换硬编码占位;为空时隐藏回跳按钮。

## 3. 测试与验证

- [x] 3.1 海报渲染器单测(图片逻辑所在层):渲染输出非空且可解码、目标尺寸、**不同玩家 PNG 不同**;service 仅 base64 包装,未激活拒绝由既有 activation 测试覆盖。
- [x] 3.2 wall config 服务单测/冒烟:返回配置 URL。
- [x] 3.3 编译门禁:`temp/go.work.sicau-verify` 下 `go build ./...` + `go vet ./...`。
- [x] 3.4 `openspec validate niu-delivery-finalize --strict` 通过。
- [x] 3.5 `lina-review` 通过(记录海报本人校验/公开配置无隐私/缓存无影响/i18n 无影响)。

> H5 与公开端点冒烟并入统一阶段。
