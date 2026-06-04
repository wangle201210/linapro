## Why

审计指出两处"走了形式但未做实"的交付点,且都**不需要真实凭据**,应按正式流程补实:
1. **M10 激活海报**:`Poster` 只返回字段、渲染出的 PNG 被丢弃,默认渲染器是 64×64 占位图 —— 未真正"生成海报图片"。
2. **M12 H5 回跳小程序**:H5 里 `MINIAPP_URL` 是硬编码占位串 —— 未"扫码回跳小程序"。

本变更把两者做实:海报**真生成并返回 PNG**(与证书 C9 一致),回跳 URL **走配置**(运维填真实小程序链接即可),都不依赖真实凭据。

## What Changes

- **海报真生成 PNG**:增强海报渲染器产出真实海报框图(竖版边框 + 由玩家数据派生的强调色,每人可区分),`Poster` 接口**返回 base64 PNG + 结构化字段**(停止丢弃),供小程序保存/分享。
- **回跳 URL 配置化**:新增 `miniapp.url` 配置;wall 服务暴露**公开 `wall/config` 端点**返回 `miniappUrl`;H5 改为 `fetch` 该配置作为回跳目标(配置为空时隐藏回跳按钮)。

不在本 change 范围:海报/证书的设计级 CJK 字体排版(留作可替换渲染层);小程序回跳 URL 的服务端生成(`generatescheme` 需真实凭据,本期由配置提供链接)。

## Capabilities

### New Capabilities
- `niu-poster-image`: 激活海报服务端真生成并返回个性化 PNG。
- `niu-miniapp-link`: H5 回跳小程序 URL 配置化(公开配置端点 + H5 读取)。

### Modified Capabilities
<!-- 行为增强,对外契约新增字段/端点,不破坏既有调用。 -->

## Impact

- **后端**:`apps/lina-plugins/sicau-niu/backend/`
  - 海报:增强 `posterrender` 默认渲染器;`PosterOutput`/`PosterRes` 增 `ImageBase64`;`activation_poster.go` 返回渲染 PNG(不再丢弃)。
  - 回跳:`wall` 服务 `New(Config{MiniappURL})`;新增公开 `wall/config` DTO/controller;`plugin.go` 读 `miniapp.url` 装配并绑定。
- **前端**:H5 `frontend/wall/index.html` 改 `fetch wall/config` 取回跳 URL,替换硬编码占位。
- **配置**:新增 `miniapp.url`(写入 `config.example.yaml`,默认空)。
- **数据权限**:海报仍仅本人已激活牛可取(不变);`wall/config` 为公开只读配置(仅回跳 URL,无隐私)。
- **缓存一致性**:无新增缓存 —— 无影响。
- **i18n**:无新增用户可见文案 —— 无影响。
- **验证**:海报服务层 DB 门控单测(返回非空 PNG、不同玩家不同);wall config 单测;编译门禁 + `openspec validate --strict`;H5 与公开端点冒烟。
