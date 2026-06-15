# 提案：为 CMS 插件补全企业建站产品中心与相册能力

## Why

企业建站场景中，产品展示和相册（案例/荣誉/环境照片墙）是与新闻文章并列的标配内容形态。当前`cms`插件只有文章一种内容模型，产品只能伪装成文章发布，没有价格、规格、多图等产品字段，也没有图集浏览能力，无法支撑常规企业官网交付。本次变更补全产品中心与相册两大内容形态。

## What Changes

- 新增产品中心：`plugin_cms_product`产品表（名称、slug、所属栏目、封面、多图、价格文本、规格摘要、富文本详情、推荐/置顶、草稿/发布状态、浏览量）；管理端新增“产品”页签提供 CRUD 与筛选；公开 JSON 接口提供产品列表与详情；公开 HTML 站点新增产品列表模板与产品详情模板，模板引擎新增`cms:product`循环与`{product:*}`详情标签。
- 新增相册：`plugin_cms_album`相册表与`plugin_cms_album_image`相册图片表；管理端新增“相册”页签提供相册 CRUD 与图片维护（整册图片一次性保存）；公开 JSON 接口提供相册列表与相册详情（含图片）；公开 HTML 站点新增相册列表模板与相册详情模板，模板引擎新增`cms:album`与`cms:photo`循环标签。
- 栏目类型扩展：`cms_category_type`字典新增`4=产品栏目`、`5=相册栏目`；产品/相册归属于对应类型的栏目，公开站点按栏目路径渲染对应列表模板，复用既有导航、面包屑和分页设施。
- 公开 sitemap 扩展：站点地图收录产品/相册栏目页和已发布产品详情页（沿用既有数量上限与可见性边界）。
- 前端新增多图上传组件（基于既有单图组件扩展），产品多图与相册图片复用宿主文件上传能力。
- `plugin.yaml`新增产品与相册的按钮权限项；三语运行时文案与`apidoc`翻译同步补齐。
- 新增数据库迁移`manifest/sql/003-cms-products-and-albums.sql`与卸载 SQL 更新，演示数据补充少量产品与相册 starter 内容。
- 产品评论、产品参数表结构化对比、相册视频、购物车/下单等电商能力不在本次范围内。

## Capabilities

### New Capabilities

- `cms-product-center`：产品的管理端维护、公开列表/详情接口与公开 HTML 渲染行为，包括可见性、排序、模板标签、sitemap 收录与性能边界。
- `cms-photo-albums`：相册与相册图片的管理端维护、公开列表/详情接口与公开 HTML 渲染行为，包括图片批量保存语义、数量上限与 sitemap 收录。

### Modified Capabilities

（无。`openspec/specs/`下尚无 cms 基线规范；上一变更`add-cms-seo-and-article-ops`仍处于活跃状态未归档，其 sitemap 行为由本变更两个新能力规范中的收录需求做增量扩展，归档时合并。）

## Impact

- `apps/lina-plugins/cms/manifest/sql/`：新增`003-cms-products-and-albums.sql`（建表、索引、字典种子）、`uninstall/`同步、`mock-data/`补充 starter 产品与相册。
- `apps/lina-plugins/cms/backend/`：`make dao p=cms`生成产品/相册 DAO；`api/cms/v1/`新增产品与相册 DTO 文件；service 新增`cms_product.go`、`cms_album.go`；controller 新增对应实现与公开渲染扩展；`plugin.go`绑定新路由。
- `apps/lina-plugins/cms/public/templates/`：新增`product-list.html`、`product-detail.html`、`album-list.html`、`album-detail.html`，`partials.html`导航适配新栏目类型。
- `apps/lina-plugins/cms/frontend/`：管理页新增产品、相册页签；新增多图上传组件；`cms-client.ts`新增 API 函数。
- `apps/lina-plugins/cms/manifest/i18n/`：三语运行时文案与`zh-CN`/`zh-TW` apidoc 翻译。
- `apps/lina-plugins/cms/hack/tests/`：新增 E2E `TC003`；后端新增单元测试。
- `apps/lina-plugins/cms/README.md`、`README.zh-CN.md`：补充产品与相册说明及模板标签文档。
- 无宿主`lina-core`代码改动。
