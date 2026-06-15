# 设计：为 CMS 插件补全产品中心与相册能力

## Context

`cms`插件当前只有文章一种内容模型。公开 HTML 渲染由`cms_public_frontend_view.go`按栏目类型选择模板（`1=列表`、`2=单页`、`3=外链`），管理端是单文件`cms-management.vue`的页签式工作台，公开模板引擎为自有`cms:`标签编译器。本次新增产品与相册两个内容模型，全部收敛在插件内，无宿主改动。

上一变更`add-cms-seo-and-article-ops`仍处于活跃状态，本变更在其代码基础上叠加（公开可见性过滤、sitemap 装配等已有设施直接复用）。

## Goals / Non-Goals

**Goals:**

- 产品中心：产品 CRUD、公开列表/详情、产品栏目类型、公开模板标签、sitemap 收录。
- 相册：相册 CRUD + 整册图片维护、公开列表/详情、相册栏目类型、公开模板标签。
- 数据装配全部有界：批量查询、固定上限、零`N+1`。

**Non-Goals:**

- 不做电商交易（价格仅为展示文本，无库存、购物车、订单）。
- 不做产品参数结构化对比、产品评论、相册视频。
- 不做新公开 JSON 端点之外的 RSS 扩展（RSS 维持文章订阅语义）。

## Decisions

### D1：产品与相册作为独立表，不复用文章表

`plugin_cms_product`独立建表（含`gallery`多图 JSON 文本列、`price`展示文本、`spec`规格摘要），`plugin_cms_album`+`plugin_cms_album_image`两表（图片行含`url/title/sort`）。否决"文章表加 type 字段"方案：产品/相册字段语义与文章差异大（多图、价格、图集行），挤在一张表会让文章契约和索引被污染，违背模块边界清晰原则。

### D2：栏目类型扩展为字典新值，渲染按类型路由模板

`cms_category_type`字典新增`4=产品栏目`、`5=相册栏目`（003 迁移 Seed DML，幂等）。`buildPublicFrontendView`在栏目命中后按`category.Type`分流：类型 4 装配产品分页数据走`product-list.html`，类型 5 装配相册分页走`album-list.html`，其余保持现状。产品详情用`?product=<slug>`、相册详情用`?album=<id>`查询参数，与既有`?article=<slug>`形态一致。栏目表新增`list_template`默认值按类型回退（产品栏目缺省`product-list.html`、相册栏目缺省`album-list.html`），不改表结构。

### D3：管理接口契约镜像文章模块

- 产品：`GET/POST /cms/products`、`GET/PUT/DELETE /cms/products/{id}`，权限`cms:product:query/add/edit/remove`；列表筛选栏目、状态、名称；分页上限 100。
- 相册：`GET/POST /cms/albums`、`GET/PUT/DELETE /cms/albums/{id}`，权限`cms:album:query/add/edit/remove`。相册保存请求携带完整图片数组（`url/title/sort`，单册上限 100 张），service 在事务中先删后插整册替换——相册图片是从属值集合而非独立资源，整册替换语义最简单且一次事务两条语句，避免逐图增删接口与前端瀑布调用。
- 公开 JSON：`GET /cms/public/products`、`GET /cms/public/products/{slug}`（详情自增浏览量）、`GET /cms/public/albums`、`GET /cms/public/albums/{id}`。公开列表只返回已发布产品/启用相册，且归属栏目必须启用，复用与文章一致的可见性子查询模式。
- 产品状态复用`cms_article_status`字典语义但独立字典`cms_product_status`（草稿/已发布）；相册仅`cms_status`启停。产品发布时间规则与文章一致（首发写 now、保留已有），公开侧按`published_at <= now`过滤——直接复用既有`publishedAtForStatus`与可见性模式。

### D4：列表多图与图片装配的性能边界

- 产品`gallery`存 JSON 字符串列（上限 9 张），列表与详情都单行自带，无关联查询。
- 相册列表页需要封面与图片计数：封面为相册表自有字段；图片计数用一条`GROUP BY album_id`聚合查询对当前页相册批量装配，公开相册详情一条`WHERE album_id = ?`图片查询。所有路径数据库访问次数恒定（列表 2~3 条、详情 2 条）。
- 公开模板循环`cms:product`/`cms:album`沿用`limit`上限 100；`cms:photo`仅在相册详情作用域内迭代当前相册图片，无额外查询。

### D5：公开模板与标签

新增四个模板：`product-list.html`、`product-detail.html`、`album-list.html`、`album-detail.html`，复用`partials.html`头尾与分页设施。模板编译器扩展：

- 循环：`{cms:product limit=12 code=...}`（`[product:link/name/image/price/summary/index]`）、`{cms:album limit=12}`（`[album:link/name/cover/count/description]`)、`{cms:photo}`（`[photo:url/title/index]`，仅相册详情）。
- 详情标签：`{product:name/price/spec/summary/content/views/date}`、`{product:gallery}`输出多图轮播 HTML、`{album:name/description}`。
- 列表页分页复用既有`{page:*}`标签（产品/相册分页数据走同一`Pagination`视图模型）。

### D6：前端管理页与多图上传

- `cms-management.vue`新增“产品”“相册”两个页签，交互复刻文章页签（筛选栏+表格+弹窗表单）；产品表单含富文本详情（复用`CmsRichTextEditor`）与多图；相册表单含图片列表维护（增删、排序值、标题）。
- 新增`CmsImageListUpload.vue`多图组件：基于既有`CmsImageUpload`的上传通路（`uploadApi`+scene），`value`为`{url,title?}[]`，支持`maxCount`。既有单图组件不改。
- 按钮权限串`cms:product:*`、`cms:album:*`经`plugin.yaml`菜单项下发，前端用`hasAccessByCodes`显隐。

### D7：sitemap 收录扩展

`GetPublicSeoContent`栏目投影的类型过滤从`IN (1,2)`扩展为`IN (1,2,4,5)`；新增一条已发布产品投影查询（沿用 5000 上限与`published_at`过滤），sitemap 输出产品详情`<loc>`。相册详情不进 sitemap（图集页 SEO 价值低且无 slug），仅收录相册栏目页。RSS 不变。

### D8：能力归属、数据权限与缓存

- 全部为源码插件自有能力；管理接口按权限串控制（插件`platform_only`无行级数据权限维度），任一批量/详情操作先校验目标存在；公开接口与既有公开资源例外一致，仅暴露已发布/启用且栏目启用的内容，不泄露草稿存在性。
- 无缓存组件引入；公开渲染仍为实时有界查询。模板`sync.Once`编译缓存因新增模板文件而自然覆盖新模板，无失效问题。
- DI 不变：仍只有`cmssvc.Service`一个运行期依赖，新方法挂在既有接口上。

### D9：SQL 与生成流程

- 新迁移`manifest/sql/003-cms-products-and-albums.sql`：`CREATE TABLE IF NOT EXISTS`三表、索引（产品`(status, published_at)`、`(category_id)`、slug 唯一；相册`(status, sort)`；图片`(album_id, sort)`）、字典 Seed `ON CONFLICT DO NOTHING`。
- 卸载 SQL 在`uninstall/`同步 DROP；starter 内容在`mock-data/002`追加产品栏目/相册栏目、3 个产品、2 个相册（幂等`NOT EXISTS`）。`ClearSiteData`与`LoadSampleData`的清理表清单加入三张新表。
- `make db.init`后`make dao p=cms`生成 DAO/DO/Entity，`make ctrl p=cms`生成控制器骨架（注意生成器会向`cms_new.go`追加重复骨架，需按既有构造函数收敛）。

## Risks / Trade-offs

- [相册整册替换在并发编辑下后写覆盖先写] → 管理操作低频，事务保证原子性；并发编辑冲突提示留待后续迭代。
- [产品 gallery 用 JSON 文本列无法按图查询] → 当前无按图检索需求；如未来需要再拆表迁移。
- [新栏目类型对旧模板`partials.html`导航的影响] → 导航只渲染链接与名称，类型无关；产品/相册栏目路径链接复用`publicFrontendCategoryPathHref`，无需改导航逻辑。
- [`mock-data/002`属上一迭代文件，追加内容违反"不改旧迭代 SQL"默认规则] → starter 内容文件本身定位是"当前交付示例站点"且`LoadSampleData`固定加载该文件名，新增独立 mock 文件无法被该功能加载；按规则记录例外原因，仅追加幂等内容不改既有行。

## Migration Plan

随插件安装 SQL 自动建表；已有环境执行`make db.init`增量建表（幂等）。回滚即回退代码并执行 uninstall SQL。无既有数据迁移。

## Open Questions

（无）
