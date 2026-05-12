# LinaPro CMS 插件

CMS 源码插件为 LinaPro 提供内容站点能力。它负责站点配置、栏目树、文章内容、友情链接、轮播图、访客留言以及公开 HTML 页面渲染。

## 范围

- 源码插件目录为 `apps/lina-plugins/cms`
- 插件自有 `plugin_cms_*` 数据表
- 管理端认证接口位于 `/api/v1/cms/*`
- 公开只读接口位于 `/api/v1/cms/public/*`
- 公开 HTML 页面位于 `/cms-site`
- 通过插件运行时嵌入 Vben 管理页面
- 插件自有 E2E 测试位于 `hack/tests`

## 公开模板

公开模板位于 `public/templates`。模板标签使用 CMS 自有的 `cms:` 前缀，例如 `cms:list`、`cms:nav`、`cms:content`、`cms:prev`、`cms:next`。参考站点使用的静态资源位于 `public/assets`。

## 开发

该插件按独立源码插件生成与测试。数据库代码生成使用 `backend/hack/config.yaml`，生成的 DAO、DO、Entity 文件均保留在插件自己的 `backend/internal` 目录下。
