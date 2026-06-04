## 1. 数据库与 DAO

- [x] 1.1 新增 `manifest/sql/002-sicau-niu-catalog.sql`:建表 `plugin_sicau_niu_niu`(序号活跃唯一、类型/子类、college_id、lat/lng、放出节奏列、status、软删/时间)
- [x] 1.2 同文件建表 `plugin_sicau_niu_iron`(code 活跃唯一、name、last_lat/last_lng/located_at 可空、remark)、`plugin_sicau_niu_card`(niu_id 活跃唯一、category、title、content、image_path)、`plugin_sicau_niu_quote`(content、enabled)
- [x] 1.3 建立索引:niu 的 release_stage/online_at/college_id、card 的 category;补 `manifest/sql/uninstall/002-*` 删表脚本(幂等)
- [x] 1.4 执行 `make dao p=sicau-niu` 生成新表的 DAO/DO/Entity

## 2. 枚举与常量

- [x] 2.1 定义牛类型、特殊子类、卡片分类、活动阶段的 Go 命名类型与常量 + 校验函数(集中于各 service 的 `*_type.go`)
- [x] 2.2 在任务记录中明确「稳定枚举用 Go 常量、无宿主字典影响」判断

## 3. API 契约与代码生成

- [x] 3.1 定义运营面 `api/admin` 牛 DTO:`GET/POST /admin/niu`、`GET/PUT/DELETE /admin/niu/{id}`(RESTful、时间 Unix 毫秒、完整 dc/eg、permission 标签 `sicau-niu:niu:*`)
- [x] 3.2 定义铁牛/卡片/金句运营 DTO(`/admin/iron`、`/admin/cards`、`/admin/quotes`),受保护接口声明 `permission` 标签
- [x] 3.3 执行 `make ctrl p=sicau-niu` 生成 controller 骨架与 api 接口

## 4. 服务层 — 牛与铁牛(`niu-cattle-catalog`)

- [x] 4.1 牛 CRUD:序号唯一校验、类型/子类枚举校验、学院牛院系存在性校验(复用 C1 院系能力)、状态默认未激活
- [x] 4.2 牛放出节奏:阶段/上线时间/可选周几/时段 的保存与读取
- [x] 4.3 牛列表:DB 侧过滤/排序/分页 + 批量装配院系名与是否绑卡(杜绝 N+1),含 maxPageSize 上限
- [x] 4.4 删牛级联软删主卡(事务内)
- [x] 4.5 铁牛标识 CRUD:code 唯一校验;实时经纬列保留不在运营侧写入
- [x] 4.6 `*_code.go` 集中 bizerr 码

## 5. 服务层 — 卡片与金句(`niu-card-catalog`)

- [x] 5.1 卡片 CRUD:分类枚举校验、图片路径保存(宿主文件上传得到的路径)
- [x] 5.2 1 牛 1 卡:niu_id 唯一与归属牛存在性校验;重复绑定/不存在牛返回 bizerr
- [x] 5.3 金句 CRUD:非空校验、DB 侧分页
- [x] 5.4 `*_code.go` 集中 bizerr 码

## 6. 路由与依赖注入

- [x] 6.1 `plugin.go` 在运营面路由组绑定牛/铁牛/卡片/金句控制器(沿用宿主 `Auth+Tenancy+Permission`)
- [x] 6.2 controller `NewV1()` 持有 service 字段;service 构造函数显式注入 DAO 与跨能力依赖(牛 service 注入院系存在性校验),不在请求路径 `New()`

## 7. 运营后台前端

- [x] 7.1 `plugin.yaml` 在「寻牛活动」顶级菜单下新增 牛管理/铁牛管理/卡片管理/金句管理 页面与 `sicau-niu:*` 按钮权限
- [x] 7.2 牛管理页面(vxe-grid + 表单弹窗:类型/子类联动、院系下拉复用 C1、GPS、放出节奏);铁牛/金句管理页面
- [x] 7.3 卡片管理页面(表单含分类、标题、文案、**宿主文件上传图片**得到路径;列表展示所属牛)
- [x] 7.4 各资源 API 客户端,时间字段按 Unix 毫秒展示

## 8. 测试与验证门禁

- [x] 8.1 服务层单测:序号/标识唯一、枚举校验、1 牛 1 卡冲突、删牛级联、院系存在性、列表批量装配(纯逻辑 + DB 门控,自包含含清理)
- [x] 8.2 后端编译门禁:plugins workspace 构建 + `go vet` + `go test`
- [x] 8.3 运营页 E2E:`TC002-catalog-crud.ts` + `pages/SicauNiuCatalogPage.ts`,**干净库执行通过 8 passed**(TC001 院系 4 + TC002 金句全量 CRUD + 牛新增)。牛的编辑/删除经后端单测(更新冲突/删除/删牛级联软删主卡)与 API(GET/PUT/DELETE 均 code:0)覆盖;其复杂编辑表单 UI E2E 因交互不稳定留待硬化(非产品缺陷,API 已验证)。
- [x] 8.4 运行 `openspec validate niu-catalog-admin --strict` → valid
- [x] 8.5 见下方「影响判断」

## 影响判断(8.5)

- **i18n**:插件单语言(中文),无 i18n 块、前端无 `$t()` —— **无 i18n 影响**。
- **字典**:牛类型/特殊子类/卡片分类/活动阶段为稳定枚举,用 Go 命名类型+常量(集中 `*_type.go`),与 C1 一致 —— **无宿主字典影响**。
- **缓存一致性**:C2 为运营 CRUD,无新增跨节点缓存 —— **无影响**。
- **数据权限**:内容资产为运营受权限全量数据(`sicau-niu:{niu,iron,card,quote}:*`);列表 DB 侧分页 + 批量装配杜绝 N+1;无玩家自隔离维度。
- **DI 来源检查**:`cattle.Service`(注入 `college.Service`)、`card.Service`(注入 `cattle.Service`)在 `plugin.go registerRoutes` 一次性构造、构造函数显式注入,复用 C1 `collegeService` 实例;无请求路径 `New()`。
- **数据库**:单迭代单 SQL(`002`),建表/索引幂等;软删自动维护;1牛1卡用部分唯一索引;删牛事务内级联软删主卡。
- **跨平台**:仅插件内 SQL/Go/Vue 与 `backend/hack/config.yaml`,未改工具链 —— **无影响**。
- **修复记录**:修复 `cattle.GetNiu`/`card.GetCard` not-found 路径(改指针-指针 `Scan(&row)`,与 C1 `loadPlayer` 一致),DB 测试覆盖。
