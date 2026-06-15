# 任务：为 CMS 插件补全常见站点能力

## 1. 后端：定时发布

- [x] 1.1 `api/cms/v1`文章创建与更新 DTO 新增可选`publishedAt`（Unix 毫秒）字段，完善`dc`/`eg`标签；service `ArticleSaveInput`透传该字段并实现 D4 写入规则
- [x] 1.2 `applyPublicArticleVisibility`追加`published_at`非空且不晚于当前时刻的过滤（`WhereLTE`对`NULL`天然不命中），公开 JSON、HTML 渲染、搜索路径均经由该方法；已核对两份 mock SQL 共 78 行文章数据，已发布行`published_at`均为历史固定时间，无需修正
- [x] 1.3 新增/更新单元测试：显式未来发布时间保存、默认发布时间规则、公开侧到点前不可见与到点后可见

## 2. 后端：批量操作

- [x] 2.1 新增`api/cms/v1/cms_article_batch_update.go`（`PUT /cms/articles`，ids+status，上限 100）与`cms_article_batch_delete.go`（`DELETE /cms/articles`，ids，上限 100），执行`make ctrl p=cms`生成控制器骨架并填充实现（生成器在`cms_new.go`追加的重复骨架已按既有构造函数收敛）
- [x] 2.2 service 新增`BatchUpdateArticleStatus`/`BatchDeleteArticles`：事务内一条计数校验任一缺失整体拒绝，集合化 UPDATE/DELETE，批量发布对`published_at IS NULL`行补写当前时间，语句数恒定（2~3 条）与 ids 数量无关
- [x] 2.3 `plugin.go`绑定两个新管理接口。路由绑定门禁`cd apps/lina-core && go test ./internal/cmd -count=1`存在与本变更无关的既有失败：`sicau-niu`提交`f844547`引入第 3 处注册期 panic 但宿主 allowlist 仍记 2，`TestProductionPanicsMatchAllowlist`必挂。按用户要求不修改`lina-core`代码，该既有失败留待独立反馈处理；曾验证将 allowlist 计数改为 3 后该包全部通过（含本变更路由绑定），证明本变更未引入新失败
- [x] 2.4 新增单元测试：批量发布/下线/删除成功路径、含不存在 ID 整体拒绝、发布时间保留与补写语义

## 3. 后端：公开 SEO 端点

- [x] 3.1 service 新增 SEO 投影装配方法：sitemap 文章（上限 5000）与启用栏目投影、RSS 最新 50 条文章投影，各一条投影 SQL
- [x] 3.2 控制器新增`sitemap.xml`、`rss.xml`、`robots.txt`原始处理器：`encoding/xml`安全转义、D3 域名回退规则、RSS 摘要剥离 HTML；`plugin.go`注册三条公开路由
- [x] 3.3 新增单元测试：sitemap 排除草稿/未到时/停用栏目/外链栏目、RSS 转义与上限、robots 包含 Sitemap 行、domain 空与非空两种链接形态

## 4. 前端与 i18n

- [x] 4.1 `cms-management.vue`文章表单新增“发布时间”控件（仅已发布状态可编辑），列表新增“定时”派生标识
- [x] 4.2 文章列表新增多选与批量发布/下线/批量删除按钮：按权限显隐、空选禁用、删除二次确认，调用新批量接口并刷新
- [x] 4.3 补齐`manifest/i18n`三语运行时文案（新按钮、字段、标识、提示），新 DTO 文档元数据补`zh-CN`/`zh-TW` apidoc 翻译（en-US 保持空占位）。`make i18n.check`存在与本变更无关的既有失败（`zh-TW`宿主公共键`pages.common.*`/`ui.crop.*`缺口，已用 stash 验证无本变更时同样失败）；本变更新增的`plugin.cms.*`键未出现在缺失列表，新增键校验通过

## 5. 测试、文档与验证

- [x] 5.1 使用`lina-e2e`技能新增`hack/tests/e2e/TC002-cms-article-batch-and-seo.ts`（TC-2a 批量发布/下线/删除闭环含空选禁用、批量删除确认与“已选 N 项”/按钮 zh-CN 文案断言；TC-2b 定时发布公开不可见+到点可见+“定时”标签文案断言；TC-2c sitemap/rss/robots 状态码、Content-Type、包含可见文章且排除定时文章断言），POM 扩展批量/定时辅助方法。E2E 全套 5 用例本地通过（TC001 两用例+TC002 三用例）。过程中发现并修复两个问题：①既有 POM `goto("/cms")`未走`workspacePath`，在`/admin`工作台基路径环境下 TC001 必挂（既有缺陷，顺带修复）；②批量删除前端最初用重复`ids`查询参数传参，GoFrame 仅解析到单值导致只删一条——E2E 抓出真实缺陷，已改为 DELETE JSON 请求体（与宿主`role/{id}/users`惯例一致）并同步 DTO `dc`与 apidoc 翻译
- [x] 5.2 更新插件`README.md`与`README.zh-CN.md`：新增公开 SEO 端点、定时发布、文章批量操作三个章节（中英内容一致）
- [x] 5.3 验证完成：`go test lina-plugin-cms/backend/... -count=1`通过（service+controller 共 17 个测试）；路由绑定门禁见 2.3 记录（既有 sicau-niu allowlist 漂移失败，与本变更无关，已验证非本变更引入）；`openspec validate --strict`通过；`make i18n.check`存在与本变更无关的既有失败（见 4.3 记录）。规则域影响判断：i18n 有影响已处理（三语运行时键+apidoc 翻译+E2E 文案断言）；缓存一致性无影响（插件无缓存组件，SEO 端点为有界实时查询，模板`sync.Once`缓存未触及）；数据权限：公开 SEO 端点与既有公开接口同源复用`applyPublicArticleVisibility`属公开资源例外，批量操作为平台管理能力（插件`platform_only`）复用既有权限串且任一目标缺失整体拒绝，均已在 design D8 记录；dev-tooling 无影响（未改 Makefile/脚本/CI）；DI 来源检查：无新增运行期依赖，控制器仍仅依赖既有`cmssvc.Service`，DI 路径`plugin.go → cmssvc.New(bizctx) → cmscontroller.NewV1(svc)`不变；数据库无迁移，索引复用`(status, published_at)`，批量语句数恒定无 N+1
- [x] 5.4 `lina-review`已执行：读取`AGENTS.md`与全部命中规则文件（openspec/plugin/architecture/data-permission/api-contract/backend-go/database/frontend-ui/testing/i18n/documentation/cache-consistency），插件根目录无本地`AGENTS.md`；审查范围覆盖父仓库与`apps/lina-plugins`子仓库全部已跟踪/未跟踪变更（30 个文件）。结论：未发现阻塞问题；2 个警告（API 层声明式参数校验无直接测试、`make i18n.check`既有 zh-TW 宿主键缺口建议另立反馈）已记录
