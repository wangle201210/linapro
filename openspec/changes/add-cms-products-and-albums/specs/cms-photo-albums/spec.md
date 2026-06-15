# cms-photo-albums 规范增量

## ADDED Requirements

### Requirement: 相册管理接口
CMS 插件 SHALL 提供相册管理接口：分页列表（`GET /api/v1/cms/albums`，支持按栏目、状态、名称筛选，每页上限 100）、详情（`GET /cms/albums/{id}`，含图片）、创建（`POST /cms/albums`）、更新（`PUT /cms/albums/{id}`）、删除（`DELETE /cms/albums/{id}`），权限分别为`cms:album:query/add/edit/remove`。相册 MUST 包含名称、所属栏目、封面、描述、排序、启停状态；相册图片（`url`、可选标题、排序）随创建/更新请求整册提交，单册上限 100 张，服务端 MUST 在一个事务中以整册替换语义保存（先删后插），不得逐图发起独立请求。删除相册 MUST 级联删除其图片行。

#### Scenario: 创建带图片的相册
- **WHEN** 管理员创建含 3 张图片的相册
- **THEN** 相册与图片保存成功，详情接口按排序返回 3 张图片

#### Scenario: 更新整册替换图片
- **WHEN** 管理员把相册图片从 3 张改为 2 张新图并保存
- **THEN** 详情仅返回新的 2 张图片，旧图片行不再存在

#### Scenario: 超限图片被拒绝
- **WHEN** 保存请求携带超过 100 张图片
- **THEN** 接口返回图片数量超限的业务错误，相册数据保持不变

### Requirement: 相册公开可见性与列表装配
公开相册接口与公开 HTML 渲染 SHALL 仅展示启用状态且归属栏目启用的相册。公开相册列表 MUST 返回封面与图片数量，图片数量 MUST 通过单条聚合查询对当前页批量装配；公开相册详情 MUST 以单条查询返回该相册全部图片。数据库访问次数不得随相册数或图片数线性增长。

#### Scenario: 停用相册不可见
- **WHEN** 存在停用相册或归属栏目被停用的相册
- **THEN** 公开相册列表与详情均不返回它们

#### Scenario: 列表返回图片计数
- **WHEN** 匿名访问公开相册列表
- **THEN** 每个相册条目包含封面与正确的图片数量

### Requirement: 相册栏目公开渲染
栏目类型字典 SHALL 新增`5=相册栏目`。访问相册栏目路径时，公开站点 MUST 按该栏目的列表模板（缺省`album-list.html`）渲染相册分页列表；`?album=<id>`访问相册详情页（缺省`album-detail.html`）。模板引擎 SHALL 支持`cms:album`循环（`limit`上限 100，`[album:link/name/cover/count/description]`标签）与相册详情作用域内的`cms:photo`循环（`[photo:url/title/index]`标签）及`{album:name/description}`详情标签。sitemap SHALL 收录启用的相册栏目页，相册详情页不收录。

#### Scenario: 相册栏目列表页渲染
- **WHEN** 匿名访问启用的相册栏目路径
- **THEN** 返回相册列表 HTML，包含启用相册的名称、封面与图片数量

#### Scenario: 相册详情页渲染
- **WHEN** 匿名以`?album=<id>`访问启用相册
- **THEN** 返回相册详情 HTML，按排序展示全部图片及标题

#### Scenario: 不存在的相册详情
- **WHEN** 匿名访问不存在或停用的相册详情
- **THEN** 返回公开内容未找到（404 页面）
