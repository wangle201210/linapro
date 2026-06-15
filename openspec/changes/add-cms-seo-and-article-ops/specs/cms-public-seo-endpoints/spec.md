# cms-public-seo-endpoints 规范增量

## ADDED Requirements

### Requirement: 公开站点地图输出
CMS 插件 SHALL 在`GET /cms-site/sitemap.xml`提供匿名可访问的站点地图，内容类型为`application/xml`，包含站点首页、全部启用栏目页和已发布且发布时间不晚于当前时刻的文章详情页。文章条目数量 MUST 有固定上限（5000 条，按发布时间倒序取最新），数据装配 MUST 通过固定条数的投影查询完成，数据库访问次数不得随文章数或栏目数线性增长。外链类型栏目 MUST 不进入站点地图。

#### Scenario: 输出已发布内容的站点地图
- **WHEN** 匿名访问`/cms-site/sitemap.xml`，站内存在启用栏目和已发布文章
- **THEN** 响应为`application/xml`的合法 sitemap 文档，包含首页、各启用栏目链接和各已发布文章链接

#### Scenario: 隐藏不可见内容
- **WHEN** 站内存在草稿文章、发布时间在未来的文章、停用栏目或外链栏目
- **THEN** 这些内容对应的 URL 不出现在 sitemap 中

#### Scenario: 域名缺省时的链接形态
- **WHEN** 站点配置`domain`非空
- **THEN** sitemap 中的`loc`为`domain`拼接公开路径的绝对地址；`domain`为空时输出以`/cms-site`开头的相对路径

### Requirement: 公开 RSS 订阅源输出
CMS 插件 SHALL 在`GET /cms-site/rss.xml`提供匿名可访问的 RSS 2.0 订阅源，内容类型为`application/xml`，包含最新 50 条已发布且发布时间不晚于当前时刻的文章。条目 MUST 包含标题、链接、描述和发布时间；描述 MUST 输出剥离 HTML 标签并经 XML 转义后的摘要文本，不得输出原始正文 HTML。频道元数据 MUST 来自站点配置（名称、标语、描述）。

#### Scenario: 输出最新文章订阅源
- **WHEN** 匿名访问`/cms-site/rss.xml`
- **THEN** 响应为合法 RSS 2.0 文档，频道信息来自站点配置，条目为最新已发布文章且不含未到发布时间的文章

#### Scenario: 富文本安全输出
- **WHEN** 文章摘要或标题包含 HTML 标签或 XML 特殊字符
- **THEN** RSS 条目中的对应文本已剥离 HTML 并完成 XML 转义，文档仍可被解析器正确解析

### Requirement: 公开 robots.txt 输出
CMS 插件 SHALL 在`GET /cms-site/robots.txt`提供匿名可访问的`text/plain`爬虫提示文件，允许抓取公开站点路径，并包含指向`/cms-site/sitemap.xml`的`Sitemap:`行；站点配置`domain`非空时该行 MUST 为绝对地址。

#### Scenario: 输出 robots 文件
- **WHEN** 匿名访问`/cms-site/robots.txt`
- **THEN** 响应为`text/plain`文本，包含`User-agent`规则和指向 sitemap 的`Sitemap:`行
