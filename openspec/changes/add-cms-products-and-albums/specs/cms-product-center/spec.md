# cms-product-center 规范增量

## ADDED Requirements

### Requirement: 产品管理接口
CMS 插件 SHALL 提供产品管理接口：分页列表（`GET /api/v1/cms/products`，支持按栏目、状态、名称筛选，每页上限 100）、详情（`GET /cms/products/{id}`）、创建（`POST /cms/products`）、更新（`PUT /cms/products/{id}`）、删除（`DELETE /cms/products/{id}`），权限分别为`cms:product:query/add/edit/remove`。产品 MUST 包含名称、slug（全局唯一）、所属栏目、封面、多图（上限 9 张）、价格展示文本、规格摘要、富文本详情、置顶/推荐标识、草稿/已发布状态与浏览量。产品状态为已发布时发布时间规则 MUST 与文章一致：显式传入值生效、未传时首发写入当前时刻且再次保存保留原值。

#### Scenario: 创建并发布产品
- **WHEN** 管理员以已发布状态创建带多图与价格的产品
- **THEN** 产品保存成功，管理列表返回该产品且发布时间已写入

#### Scenario: 重复 slug 被拒绝
- **WHEN** 创建产品时 slug 与既有产品重复
- **THEN** 接口返回产品 slug 已存在的业务错误，不写入数据

#### Scenario: 列表按栏目与状态筛选
- **WHEN** 管理员按产品栏目和草稿状态查询产品列表
- **THEN** 仅返回该栏目下的草稿产品，分页信息正确

### Requirement: 产品公开可见性
公开产品接口与公开 HTML 渲染 SHALL 仅展示状态为已发布、发布时间非空且不晚于当前时刻、且归属栏目启用的产品。公开详情按 slug 读取并 MUST 自增浏览量。过滤 MUST 在数据库查询阶段注入，不得内存过滤。

#### Scenario: 草稿与定时产品不可见
- **WHEN** 存在草稿产品和发布时间在未来的产品
- **THEN** 公开产品列表与详情均不返回它们，详情访问表现为未找到

#### Scenario: 公开详情自增浏览量
- **WHEN** 匿名按 slug 访问公开产品详情两次
- **THEN** 第二次返回的浏览量比第一次大

### Requirement: 产品栏目公开渲染
栏目类型字典 SHALL 新增`4=产品栏目`。访问产品栏目路径时，公开站点 MUST 按该栏目的列表模板（缺省`product-list.html`）渲染产品分页列表；`?product=<slug>`访问产品详情页（缺省`product-detail.html`）。模板引擎 SHALL 支持`cms:product`循环（`limit`上限 100，`[product:link/name/image/price/summary/index]`标签）与`{product:name/price/spec/summary/content/views/date/gallery}`详情标签。列表与详情的数据装配 MUST 保持数据库访问次数恒定，多图随产品行自带不产生关联查询。

#### Scenario: 产品栏目列表页渲染
- **WHEN** 匿名访问启用的产品栏目路径
- **THEN** 返回产品列表 HTML，包含已发布产品的名称、封面与价格，分页标签可用

#### Scenario: 产品详情页渲染
- **WHEN** 匿名以`?product=<slug>`访问已发布产品
- **THEN** 返回产品详情 HTML，包含名称、价格、规格、富文本详情与多图

### Requirement: 产品进入站点地图
`/cms-site/sitemap.xml` SHALL 收录启用的产品栏目页与已发布且到点的产品详情页，沿用既有文章收录的数量上限与可见性边界；草稿与未到点产品 MUST 不出现。

#### Scenario: sitemap 收录产品
- **WHEN** 存在已发布产品与发布时间在未来的产品
- **THEN** sitemap 包含前者的详情 URL 且不包含后者
