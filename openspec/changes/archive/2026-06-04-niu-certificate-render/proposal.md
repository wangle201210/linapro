## Why

需求 M11 明确要求「电子证书**生成**与展示」。C5 已做证书类荣誉的定义、解锁判定与 C7 批量授予,但**证书图片的服务端生成尚未实现**(此前仅保存模板路径 image_path)。本变更补齐服务端个性化证书 PNG 生成,使玩家可获取并展示自己获得的电子证书。

## What Changes

- 新增**证书渲染 seam**(`certrender`,与激活海报 `posterrender` 同构):接收证书构图数据(昵称/荣誉名/荣誉编码/校庆标识),输出 PNG;默认实现产出**每位持有者视觉可区分**的证书框图(强调色由持有者数据派生),CJK 文字排版作为可替换的设计级渲染层留在 seam 之后。
- 新增**玩家证书接口**(玩家令牌保护):玩家请求自己持有的某证书类荣誉 → 服务端**校验是证书类型且玩家已获得** → 渲染个性化 PNG → 返回 **base64 图片 + 结构化字段**(昵称/荣誉名/获得时间),供小程序展示/保存。
- honor 服务装配期注入证书渲染器与校庆标识(复用既有 poster.campusBadge 配置)。

不在本 change 范围:运营端批量导出证书图片(C7 已有批量授予与名册导出;本期聚焦"生成与展示");CJK 字体打包(留作设计级渲染层)。

## Capabilities

### New Capabilities
- `niu-certificate-render`: 证书渲染 seam + 玩家个性化电子证书 PNG 生成与展示接口。

### Modified Capabilities
<!-- 无需求级行为变更:复用 C5 honor_def/user_honor 与 C1 玩家昵称,仅新增证书渲染与展示读取路径。 -->

## Impact

- **后端**:`apps/lina-plugins/sicau-niu/backend/`
  - 新增 `internal/service/honor/internal/certrender`(seam + 默认渲染器 + 错误码)。
  - honor 服务新增 `PlayerCertificate`:证书类型校验 + 持有校验 + 字段装配 + 渲染;`New()` 改为注入证书渲染器与校庆标识。
  - 新增 `api/player/v1` 证书 DTO、player 控制器处理方法;`plugin.go` 玩家子组绑定并改 `honorsvc.New(...)` 装配。
- **数据库**:无新增表/列;复用 `honor_def`、`user_honor`、`user`。
- **数据权限**:玩家令牌保护,**仅能渲染本人已获得的证书**(校验 user_honor 持有),杜绝越权获取他人证书。
- **缓存一致性**:按需渲染,无新增缓存/失效路径 —— 无影响。
- **i18n**:单语言中文 —— 无影响。
- **验证**:服务层 DB 门控单测(证书类型校验、持有校验拒绝、渲染非空 PNG、个性化区分)+ 编译门禁 + `openspec validate --strict`;小程序 UI 不在本期,接口冒烟并入统一阶段。
