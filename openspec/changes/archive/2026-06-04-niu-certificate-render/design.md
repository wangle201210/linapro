## Context

补齐 M11「电子证书生成与展示」的服务端生成。复用 C5 `honor_def`(certificate 类型)、`user_honor`(持有)、C1 `user`(昵称),与激活海报 `posterrender` 同构落地一个证书渲染 seam。约束遵循 `.agents/rules/` 的 `backend-go`、`api-contract`、`data-permission`、`architecture`、`plugin`、`testing`;单语言中文不启用 i18n;`platform_only`。

## Goals / Non-Goals

**Goals:** 证书渲染 seam(可替换设计级渲染);玩家按需获取本人已获得证书的个性化 PNG(base64)+ 结构化字段;严格持有校验。

**Non-Goals:** 运营批量导出证书图片;CJK 字体打包与排版(设计级渲染层);新增表/列;小程序 UI。

## Decisions

### D1 证书渲染 seam(同构 posterrender)
新增 `internal/service/honor/internal/certrender`:`CertRenderer` 接口 + `CertData{Nickname, HonorName, HonorCode, CampusBadge}` 值对象 + 默认 `basicRenderer` + `CodeCertRenderFailed`。默认实现输出**真实证书框图 PNG**(边框 + 校庆底色 + **强调色由 `hash(nickname+honorCode)` 派生**,使每位持有者证书视觉可区分),并**实际返回 PNG 字节**(区别于 poster 默认丢弃)。CJK 文字排版作为设计级渲染层在 seam 之后替换,契约不变。

### D2 玩家证书获取 + 持有校验(数据权限)
`PlayerCertificate(ctx, playerID, honorID)`:
1. 加载荣誉定义;不存在 → 业务错误;**`honor_type != certificate` → 业务错误**。
2. **校验 `user_honor` 存在 `(playerID, honorID)`**;不持有 → 业务错误(杜绝越权获取他人/未获得证书)。
3. 装配字段(本人昵称 + 荣誉名/编码 + 获得时间 + 校庆标识)→ 渲染 PNG → 返回 base64 + 字段。

### D3 显式依赖注入
honor 服务新增 `certRenderer` 与 `campusBadge` 字段,`New(certRenderer, Config{CampusBadge})` 装配期注入;`plugin.go` 复用既有 `poster.campusBadge` 配置构造,玩家子组绑定证书接口。re-export `CertRenderer` 与 `NewBasicCertRenderer` 供装配层构造默认渲染器。

### API 契约(玩家令牌保护)
| 方法 | 路径 | 说明 |
|------|------|------|
| `GET` | `/plugins/sicau-niu/player/certificates` | 入参 `honorId`;返回本人该证书的 base64 PNG + 昵称/荣誉名/编码/获得时间(Unix 毫秒) |

## Risks / Trade-offs

- [越权获取证书] → 强制 `user_honor` 持有校验 + 令牌内玩家身份;单测覆盖未持有/非证书拒绝。
- [无 CJK 字体导致文字不可读] → 默认渲染为框图 + 个性化强调色 + 同时返回结构化字段供前端叠加;设计级字体渲染留在 seam 之后(与 poster 一致)。
- [按需渲染成本] → 单张小图、低频(玩家查看自己证书),无 N+1。

## Migration Plan

1. `certrender` seam → honor `PlayerCertificate` → DTO/controller → `plugin.go` 绑定 + `honorsvc.New(renderer, cfg)`。
2. 服务层 DB 门控单测(类型/持有校验、PNG 非空、不同持有者输出不同)+ 编译门禁 + `openspec validate --strict`。
3. 回滚:纯新增 seam/读取路径,无表/列/数据,可整体回退。

## Open Questions

1. 证书 PNG 画布尺寸与默认配色(暂定 600×400、校庆绿+金)。

> 影响判断:字典 —— 无新增枚举(复用 honor_type),**无影响**;缓存一致性 —— 按需渲染,**无新增缓存失效路径**;数据权限 —— 玩家仅可取本人已获证书(已记录持有校验);跨平台 —— 仅插件内文件,**无影响**;i18n —— 单语言,**无影响**。
