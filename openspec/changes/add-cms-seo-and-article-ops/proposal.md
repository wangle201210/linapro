# 提案：为 CMS 插件补全常见站点能力

## Why

`cms`插件对照常规 CMS 框架仍缺少几项标配能力：公开站点没有`sitemap.xml`、`rss.xml`和`robots.txt`，搜索引擎收录和订阅分发无从谈起；文章发布时间永远等于操作时刻，公开侧也不按发布时间过滤，无法实现常见的定时发布；管理端文章只能逐篇发布、下线和删除，内容量稍大时运营效率低。本次变更补齐这三项常见能力，使插件达到常规 CMS 的基础功能预期。

## What Changes

- 新增公开 SEO 端点：`/cms-site/sitemap.xml`输出已发布文章与可见栏目的站点地图，`/cms-site/rss.xml`输出最新已发布文章的 RSS 2.0 订阅源，`/cms-site/robots.txt`输出指向 sitemap 的爬虫提示文件；三者均为匿名只读端点，输出数量有固定上限。
- 新增文章定时发布：管理端创建和编辑文章时可显式指定发布时间`publishedAt`；公开接口、公开 HTML 页面、搜索、sitemap 和 RSS 只展示`status=已发布`且`published_at`不晚于当前时刻的文章；管理列表对未到发布时间的文章给出“定时”视觉标识。
- 新增管理端文章批量操作：批量发布、批量下线（改为草稿）、批量删除，单次操作的 ID 数量有上限；管理列表支持多选并提供批量操作入口。批量操作复用既有`cms:article:edit`与`cms:article:remove`权限，与宿主用户管理批量更新的权限惯例一致。
- 不改动数据库表结构：`published_at`字段与`(status, published_at)`索引已存在，无新增迁移 SQL。
- 标签管理闭环、文章评论、回收站、草稿预览不在本次范围内，留待后续迭代。

## Capabilities

### New Capabilities

- `cms-public-seo-endpoints`：CMS 公开站点的`sitemap.xml`、`rss.xml`、`robots.txt`输出行为，包括可见性边界、数量上限和缓存友好的输出方式。
- `cms-article-scheduled-publishing`：文章发布时间的管理端维护规则与公开侧按发布时间过滤的可见性规则。
- `cms-article-batch-operations`：管理端文章批量发布、批量下线、批量删除的接口契约、上限与失败语义。

### Modified Capabilities

（无。`openspec/specs/`下不存在 cms 相关基线规范，本次全部为新增能力。）

## Impact

- `apps/lina-plugins/cms/backend/api/cms/v1/`：新增批量操作 DTO 文件；文章创建、更新 DTO 增加`publishedAt`字段；文章响应增加`publishedAt`投影。
- `apps/lina-plugins/cms/backend/internal/controller/cms/`：新增批量操作控制器方法与公开`sitemap.xml`、`rss.xml`、`robots.txt`处理方法。
- `apps/lina-plugins/cms/backend/internal/service/cms/`：文章服务新增批量状态变更、批量删除与 SEO 数据装配；公开可见性过滤追加`published_at <= now`条件。
- `apps/lina-plugins/cms/backend/plugin.go`：注册新公开路由与新管理接口绑定。
- `apps/lina-plugins/cms/frontend/pages/cms-management.vue`：文章表单增加发布时间字段，文章列表增加多选与批量操作入口、定时状态标识。
- `apps/lina-plugins/cms/manifest/i18n/`：三语运行时文案与`apidoc`翻译资源同步更新。
- `apps/lina-plugins/cms/README.md`、`README.zh-CN.md`：补充新能力说明。
- `apps/lina-plugins/cms/hack/tests/`：新增对应 E2E 用例；`backend/internal/`新增单元测试。
- 无数据库迁移、无宿主`lina-core`代码改动、无缓存结构变更。
