# 任务：为 CMS 插件补全产品中心与相册能力

## 1. 数据库与生成

- [x] 1.1 新增`manifest/sql/003-cms-products-and-albums.sql`：`plugin_cms_product`、`plugin_cms_album`、`plugin_cms_album_image`三表（幂等建表+注释+索引）、`cms_category_type`字典新增 4/5 值、`cms_product_status`字典（Seed `ON CONFLICT DO NOTHING`）；`uninstall/`同步 DROP 三表；执行`make db.init`
- [x] 1.2 执行`make dao p=cms`生成三表 DAO/DO/Entity；`mock-data/002-cms-starter-content.sql`追加产品/相册栏目与示例数据（幂等`NOT EXISTS`，记录修改旧迭代 mock 文件的例外原因）；`ClearSiteData`/`LoadSampleData`清理表清单加入三张新表

## 2. 后端：产品中心

- [x] 2.1 `api/cms/v1/`新增产品管理与公开 DTO 文件（列表/详情/创建/更新/删除/公开列表/公开详情），`make ctrl p=cms`生成骨架并收敛`cms_new.go`重复声明
- [x] 2.2 service 新增`cms_product.go`：CRUD、slug 唯一校验、发布时间规则复用、公开可见性过滤（状态+到点+栏目启用）、公开详情自增浏览量；`cms_code.go`新增产品错误码
- [x] 2.3 `plugin.go`绑定产品管理与公开路由；新增产品单元测试（CRUD、slug 冲突、定时可见性、浏览量自增）

## 3. 后端：相册

- [x] 3.1 `api/cms/v1/`新增相册管理与公开 DTO（图片数组内嵌、单册上限 100），`make ctrl p=cms`生成骨架
- [x] 3.2 service 新增`cms_album.go`：CRUD、整册图片事务替换（先删后插）、删除级联清图、公开列表图片计数聚合装配、公开详情单查图片；新增相册错误码
- [x] 3.3 `plugin.go`绑定相册路由；新增相册单元测试（整册替换、级联删除、停用不可见、计数装配正确）

## 4. 公开渲染与 SEO

- [x] 4.1 新增`product-list.html`、`product-detail.html`、`album-list.html`、`album-detail.html`模板；模板编译器扩展`cms:product`/`cms:album`/`cms:photo`循环与`{product:*}`/`{album:*}`标签；栏目类型 4/5 缺省模板回退
- [x] 4.2 `buildPublicFrontendView`按栏目类型 4/5 装配产品/相册分页视图，`?product=`/`?album=`详情分流；管理端栏目模板下拉补充新模板选项
- [x] 4.3 `GetPublicSeoContent`栏目类型过滤扩展为 1/2/4/5，新增已发布产品投影并进 sitemap；新增渲染与 sitemap 单元测试

## 5. 前端与 i18n

- [x] 5.1 新增`CmsImageListUpload.vue`多图组件；`cms-client.ts`新增产品与相册 API 函数和类型
- [x] 5.2 `cms-management.vue`新增“产品”页签：筛选+表格+表单（多图、富文本、价格、规格、状态、发布时间）
- [x] 5.3 `cms-management.vue`新增“相册”页签：列表+表单（封面、描述、图片列表维护）；`plugin.yaml`新增产品/相册权限项
- [x] 5.4 三语运行时文案（plugin.json/menu.json）与`zh-CN`/`zh-TW` apidoc 翻译补齐（en-US apidoc 维持空占位）；`make i18n.check`本变更新增`plugin.cms.*`键均未出现在缺失列表（zh-TW 宿主公共键缺口为既有问题，见上一变更记录）。另：starter 示例站点已实际使用产品/相册能力——新增产品中心/院区相册两个栏目、3 个带封面与多图的产品、2 个带封面的相册及图片，公开页面渲染与 sitemap 收录已在 dev 环境逐项 curl 验证

## 6. 测试、文档与验证

- [x] 6.1 新增`TC003-cms-products-and-albums.ts`（TC-3a 产品页签 UI 创建/删除+管理与公开 API 校验+产品栏目页/详情页渲染+sitemap 收录；TC-3b 相册整册图片保存与替换+相册栏目页/详情页渲染+不存在相册 404），POM 扩展产品/相册辅助方法（含“产品”页签 zh-CN 文案断言）。E2E 全套 7 用例通过。过程中 E2E 抓出两个真实缺陷并修复：①栏目 DTO `type`校验仍限 1,2,3，已扩为 1~5 并同步 dc 与 apidoc 翻译；②GoFrame `max-length`按字符串而非数组元素数校验，相册 images/产品 gallery/批量 ids 的数组上限全部改为 service 层执行（新增`CodeArticleBatchLimitExceeded`，gallery 超出静默截断并写入 dc 说明），相应规范措辞同步
- [x] 6.2 更新插件`README.md`/`README.zh-CN.md`：新增“产品中心与相册”章节，覆盖栏目类型、公开地址、模板标签表与模板文件（中英一致）
- [x] 6.3 验证完成：`go test lina-plugin-cms/backend/... -count=1`通过（service 28 测试+controller 14 测试）；E2E 7/7 通过；`openspec validate --strict`通过；003 迁移与 mock SQL 双跑验证幂等；dev 环境 curl 验证产品/相册公开页渲染与 sitemap 收录。规则域影响判断：数据权限——管理接口按权限串控制（平台管理能力），公开接口仅暴露已发布/启用且栏目启用内容（design D8）；缓存一致性——无缓存组件，模板`sync.Once`编译缓存包含新模板文件且随宿主编译更新，无失效问题（design D9）；dev-tooling——仅改插件`backend/hack/config.yaml`的 dao 表清单（生成器配置输入，跨平台无影响）；DI——无新增运行期依赖，新方法挂在既有`cmssvc.Service`接口；i18n——三语运行时键+menu 键+apidoc 翻译已补齐；数据库——003 迁移幂等、`album_image`表无`deleted_at`为整册替换硬删除边界（design D9 记录）、索引匹配查询路径、相册图片计数走单条 GROUP BY 聚合无 N+1；修改旧迭代 mock 文件例外原因已在 design 风险记录（LoadSampleData 固定加载该文件名）
- [x] 6.4 `lina-review`已执行：读取`AGENTS.md`与全部命中规则文件（本会话内读取且`.agents/`无变更），插件根目录无本地`AGENTS.md`；范围覆盖父仓库 openspec 文档与`apps/lina-plugins`子仓库 94 项 cms 变更（含 DAO 生成件）。审查发现并修复 1 个警告：`cms_album.go`将 5 类接口堆在单文件违反`api-contract.md`按用途拆分要求，已拆为 list/get(+delete)/create/update 四个文件并复跑 gofmt/build/vet/单测全部通过。结论：无阻塞问题
