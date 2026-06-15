# 设计：为 CMS 插件补全常见站点能力

## Context

`cms`是源码插件，公开 HTML 站点由`backend/internal/controller/cms/cms_public_frontend*.go`中的原始`ghttp`处理器渲染，管理接口走`api/cms/v1`声明的 DTO 与`make ctrl`生成的控制器骨架。文章表`plugin_cms_article`已具备`published_at`字段、`(status, published_at)`索引和`PublishedAt *int64`（Unix 毫秒）响应投影，但发布时间只能由系统在状态切换时写入，公开侧可见性也只看`status`。管理端没有批量操作，公开站点没有`sitemap.xml`、`rss.xml`、`robots.txt`。

本设计全部属于源码插件能力，不触碰`lina-core`宿主通用契约、存储模型或工作台适配层。

## Goals / Non-Goals

**Goals:**

- 公开站点提供搜索引擎和订阅器可消费的`sitemap.xml`、`rss.xml`、`robots.txt`。
- 文章支持显式发布时间与定时发布：到点前公开侧不可见，到点后自动可见，无需任何调度任务。
- 管理端支持文章批量发布、批量下线、批量删除。
- 所有新增查询路径保持有界：固定上限、投影查询、零`N+1`。

**Non-Goals:**

- 不做标签管理闭环、文章评论、回收站、草稿预览（后续迭代）。
- 不做多站点、不改数据库结构、不引入定时任务和缓存组件。
- 不改变既有公开模板语法与渲染流程。

## Decisions

### D1：SEO 端点用原始 ghttp 处理器，挂在既有公开路由组

`sitemap.xml`、`rss.xml`、`robots.txt`与`/cms-site`HTML 页面一样属于公开匿名资源，复用`plugin.go`中既有公开路由组（`NeverDoneCtx/CORS/RequestBodyLimit/Ctx`中间件），新增三个原始处理器方法，直接写`application/xml`/`text/plain`响应。不走`g.Meta`DTO：这三个端点输出的是 XML/纯文本协议文件而非 JSON API，与既有`PublicFrontendPage`的处理方式一致，也不进入 OpenAPI 文档。

替代方案：声明成`api/`DTO 接口——被否决，`HandlerResponse`中间件会包 JSON 信封，XML 协议文件无法表达。

### D2：XML 由控制器拼装，数据由 service 提供有界投影

service 新增一个 SEO 数据装配方法，返回 sitemap 与 RSS 需要的最小投影（文章：`id/slug/title/summary/published_at/updated_at`；栏目：`code/path/type/outlink/updated_at`）。约束：

- sitemap 文章上限固定`5000`条（常量），按`published_at`倒序取最新；栏目全量（栏目树本身是小集合，既有管理接口已全量加载）。
- RSS 固定取最新`50`条已发布文章（常量）。
- 各一条 SQL：文章一条投影查询、栏目一条投影查询，数据库访问次数恒定，与行数无关。
- XML 转义使用标准库`encoding/xml`的转义能力，禁止手工字符串替换转义。

数据规模假设：单站文章万级以内；5000 条 sitemap 上限覆盖常见站点，超出部分（更老的文章）不进入 sitemap，属可接受边界并写入 README。

### D3：绝对链接优先用站点 domain，缺省回退相对路径

RSS 与 sitemap 的`<loc>`/`<link>`优先使用`plugin_cms_site.domain`（去尾部`/`）拼接`/cms-site/...`路径；`domain`为空时输出相对路径。robots.txt 的`Sitemap:`行仅在`domain`非空时输出绝对地址，否则输出相对地址。该回退保证未配置域名的环境不会输出错误主机名。

### D4：定时发布 = 显式 publishedAt + 公开侧时间过滤，无调度任务

- `ArticleCreateReq`/`ArticleUpdateReq`新增可选`publishedAt *int64`（Unix 毫秒）。`status=已发布`且传入该字段时按传入值写库（可为未来时间）；未传时保持现有规则（首次发布写`now`，再次保存保留原值）；`status=草稿`时忽略该字段且不改库中已有值（与现状一致）。
- `applyPublicArticleVisibility`追加`published_at IS NOT NULL AND published_at <= now`过滤。所有公开路径（公开 JSON API、HTML 渲染、搜索、sitemap、RSS）都经由该方法，单点收敛。时间比较用`WhereLTE`等链式方法，不拼 SQL 字符串，不用数据库内置函数。
- 不引入新状态枚举：`cms_article_status`字典保持`草稿/已发布`两值；“定时”是`status=已发布 && published_at > now`的派生展示态，由前端用既有响应字段推导并打标签，不进字典、不进存储。这样避免状态机膨胀和到点改状态的调度任务。
- 索引：`idx_plugin_cms_article_status_publish (status, published_at)`已存在，恰好匹配新过滤路径，无新增迁移。

替代方案：新增`status=2(定时)`+定时任务到点改状态——被否决：需要 cron 组件、状态迁移和字典变更，且到点延迟取决于调度精度；时间过滤方案零调度、秒级精确。

### D5：批量操作遵循宿主既有批量契约形态

- 批量状态变更：`PUT /cms/articles`，请求`{ids: []int64, status: 0|1}`，权限`cms:article:edit`，镜像宿主`PUT /user`批量更新形态。
- 批量删除：`DELETE /cms/articles`，请求`{ids: []int64}`，权限`cms:article:remove`，镜像宿主`DELETE /role`批量删除形态。
- `ids`上限`100`（校验规则`max-length:100`），文档说明超限被拒绝。
- service 在事务中先一条`WHERE id IN`计数校验全部存在，任一不存在整体拒绝（`CodeArticleNotFound`）；再以集合化语句完成操作：
  - 批量发布：一条`UPDATE ... SET status=1 WHERE id IN`，再一条`UPDATE ... SET published_at=now WHERE id IN (...) AND published_at IS NULL`（保留已有发布时间，与单篇规则一致）。
  - 批量下线：一条`UPDATE ... SET status=0 WHERE id IN`，保留`published_at`。
  - 批量删除：一条`Delete()`软删除。
- 数据库访问次数恒定（2~3 条语句），与`ids`数量无关。

### D6：前端在既有单文件管理页内扩展

`cms-management.vue`是插件自带的单文件管理页（非宿主 vben 工作台页面，不使用`useVbenVxeGrid`体系，页面内自有表格实现）。沿用页面既有交互语言：

- 文章列表加复选框列与“批量发布/批量下线/批量删除”按钮，按钮按权限串显隐，删除有确认。
- 文章表单加“发布时间”选择器（`datetime-local`输入，与页面现有原生表单控件风格一致），仅状态为已发布时可编辑。
- 列表对`status=已发布 && publishedAt > now`的行展示“定时”标签。
- 所有新文案走`$t()`并补三语语言包；选项与列定义遵守“不得在模块顶层求值`$t()`”的约束。

### D7：能力归属与契约边界

三项能力均为源码插件自有能力：路由、控制器、service、模板、i18n、文档全部收敛在`apps/lina-plugins/cms/`内；对宿主仅使用既有发布中间件与`bizerr`公共组件，无新增跨模块契约、无新增运行期依赖（控制器仍只依赖既有`cmssvc.Service`，DI 路径不变：`plugin.go → cmssvc.New(bizctx) → cmscontroller.NewV1(svc)`）。

### D8：数据权限与公开资源例外

- 公开 SEO 端点输出的内容与既有公开文章/栏目接口完全同源（`applyPublicArticleVisibility`+栏目`status=启用`），属于既有“公开资源”例外：仅暴露已发布且到时的内容，不泄露草稿、定时未到或停用栏目的存在性。
- 批量操作是平台管理能力（插件`scope_nature: platform_only`，不支持多租户），与既有单篇删除一致按权限串控制，无行级数据权限维度；任一目标不存在时整体拒绝。

### D9：缓存与一致性

本插件无缓存组件，三个 SEO 端点每次请求实时查询（各 1~2 条有界 SQL），不引入缓存与失效问题；模板缓存（`sync.Once`编译）不受影响。确认无缓存一致性影响。

## Risks / Trade-offs

- [sitemap 每请求实时查询，遭爬虫高频抓取时有数据库压力] → 查询为索引覆盖的投影查询且有 5000 上限；如未来成为瓶颈，可在后续迭代加 TTL 缓存（届时按缓存一致性规则设计），本次不预置。
- [RSS 摘要含用户富文本，存在注入风险] → 摘要输出经`encoding/xml`转义为纯文本（剥离 HTML 后转义），正文不进 RSS。
- [domain 未配置时 RSS 链接为相对路径，部分阅读器不识别] → README 明确说明配置`domain`后获得绝对链接；这是配置缺失下的最大努力输出。
- [既有演示数据中文章`published_at`可能为 NULL 但状态已发布] → 公开过滤要求`published_at IS NOT NULL`，starter 数据需核对；mock SQL 中已发布文章均带发布时间（实施时验证），异常数据表现为公开不可见而非报错。
- [批量删除不可恢复（软删除但无回收站 UI）] → 前端二次确认；回收站列为后续迭代。

## Migration Plan

无数据库迁移、无配置变更。源码插件随宿主编译发布即生效；回滚即回退代码。`published_at`过滤对存量已发布文章（`published_at`均为历史时间）无可见性影响。

## Open Questions

（无。实施中若发现 starter 数据存在`status=1`且`published_at IS NULL`的行，在 mock SQL 中一并修正。）
