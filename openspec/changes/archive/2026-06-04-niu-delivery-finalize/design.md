## Context

把两处不需凭据的交付点做实:激活海报真生成 PNG、H5 回跳 URL 配置化。复用 C3 海报渲染 seam(`posterrender`)与 C6 wall 公开端点。约束遵循 `.agents/rules/` 的 `backend-go`、`api-contract`、`frontend-ui`、`data-permission`、`plugin`、`testing`;单语言中文不启用 i18n;`platform_only`。

## Goals / Non-Goals

**Goals:** 海报真返回个性化 PNG;回跳 URL 走配置并经公开端点供 H5 读取。

**Non-Goals:** CJK 字体设计级排版;服务端生成小程序 scheme(需凭据);新增表。

## Decisions

### D1 海报真生成并返回 PNG
增强 `posterrender` 默认渲染器:从 64×64 占位升级为竖版海报框图(边框 + 由 `hash(nickname+niuCode)` 派生的强调色,每位玩家可区分),**实际返回 PNG 字节**。`PosterOutput`/`PosterRes` 增 `ImageBase64`;`activation_poster.go` 把渲染结果 base64 后返回(**不再丢弃**)。结构化字段保留,供前端叠加可读文案;设计级字体排版仍在 seam 之后替换。

### D2 回跳 URL 配置化 + 公开端点
新增插件配置 `miniapp.url`。wall 服务 `New(Config{MiniappURL})` 装配期注入。新增**公开**(免登录)端点 `GET /plugins/sicau-niu/wall/config` 返回 `{miniappUrl}`。H5 `fetch` 该端点取回跳目标,替换硬编码占位;`miniappUrl` 为空时隐藏回跳按钮。回跳 URL 由运维提供(真实小程序链接/scheme),不需运行期凭据。

### D3 数据权限
海报:仍校验本人已激活目标牛(不变)。`wall/config`:公开只读,仅返回回跳 URL,无隐私字段。

### API 契约
| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| `GET` | `/plugins/sicau-niu/player/poster` | 玩家令牌 | 响应新增 `imageBase64`(海报 PNG) |
| `GET` | `/plugins/sicau-niu/wall/config` | 公开 | `{miniappUrl}` |

## Risks / Trade-offs

- [无 CJK 字体海报文字不可读] → 框图 + 个性化强调色 + 同时返回字段;字体排版留 seam(与证书一致)。
- [回跳 URL 形态多样] → 由配置提供完整链接/scheme,前端直接使用;为空则隐藏入口。

## Migration Plan

1. 海报:增强渲染器 → `PosterOutput`/`PosterRes` 加 `ImageBase64` → service 返回 → controller 透传。
2. 回跳:`wall.New(Config)` + `wall/config` DTO/controller + `plugin.go` 读 `miniapp.url` 绑定 + `config.example.yaml` + H5 fetch。
3. 海报 DB 门控单测(PNG 非空、玩家间不同)+ wall config 单测 + 编译门禁 + `openspec validate --strict`。
4. 回滚:纯增强 + 新增端点/配置,可整体回退。

## Open Questions

1. 海报画布尺寸(暂定 600×800 竖版)。

> 影响判断:字典 —— 无新增枚举,**无影响**;缓存一致性 —— 无新增缓存,**无影响**;数据权限 —— 海报本人校验不变、wall/config 公开无隐私(已记录);跨平台 —— 仅插件内文件,**无影响**;i18n —— 无新增文案,**无影响**。
