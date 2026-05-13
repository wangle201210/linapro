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

公开模板位于`public/templates`。模板标签使用 CMS 自有的`cms:`前缀，参考站点使用的静态资源位于`public/assets`。

### 模板文件

| 文件 | 用途 |
| --- | --- |
| `partials.html` | 统一的`head`、导航、侧栏与页脚片段 |
| `index.html` | 公开首页 |
| `list.html` | 普通文章列表页 |
| `list-card.html` | 卡片式文章列表页 |
| `search.html` | 公开搜索结果页 |
| `detail.html` | 文章详情页 |
| `single.html` | 单页栏目 |
| `message.html` | 访客留言页 |

### 基础语法

| 语法 | 说明 |
| --- | --- |
| `{include file=comm/head.html}` | 引入公共片段 |
| `{cms:if(条件)}...{/cms:if}` | 条件渲染，支持`{else}`分支 |
| `{cms:list num=10 order=date}...{/cms:list}` | 循环输出文章列表 |
| `{cms:search num=10 order=date}...{/cms:search}` | 循环输出搜索结果 |
| `{cms:nav num=15 parent=0}...{/cms:nav}` | 循环输出栏目导航 |
| `{cms:slide num=5 gid=1}...{/cms:slide}` | 循环输出轮播图 |
| `{cms:link num=10 gid=1}...{/cms:link}` | 循环输出友情链接 |
| `{cms:sort scode=25}...{/cms:sort}` | 按栏目编号读取一个栏目 |

支持的公共片段文件名包括`head.html`、`foot.html`、`sidebar.html`和
`sortnav.html`。

循环的`num`最大值为`100`。`cms:list`未指定`num`时默认每页`12`条。

### 全局站点标签

| 标签 | 输出 |
| --- | --- |
| `{cms:sitepath}` | 公开站点根路径，固定为`/cms-site` |
| `{cms:sitetplpath}` | 公开资源路径，固定为`/cms-site/assets` |
| `{cms:scaction}` | 搜索表单地址，固定为`/cms-site/search` |
| `{cms:msgpage}` | 留言页面地址 |
| `{cms:msgaction}` | 留言提交地址 |
| `{cms:sitetitle}`、`{cms:companyname}` | 站点名称 |
| `{cms:sitesubtitle}`、`{cms:siteslogan}` | 站点标语 |
| `{cms:sitelogo}` | 站点 Logo 地址 |
| `{cms:sitekeywords}` | 站点关键词 |
| `{cms:sitedescription}` | 站点描述 |
| `{cms:companyaddress}` | 联系地址 |
| `{cms:companyphone}` | 联系电话 |
| `{cms:companyemail}` | 联系邮箱 |
| `{cms:companycontact}` | 联系人 |
| `{cms:siteicp}` | 备案信息 |
| `{cms:companyweixin}` | 微信二维码地址 |
| `{cms:keyword}` | 当前搜索关键词 |
| `{cms:year}` | 当前年份 |
| `{cms:firstsortlink}` | 第一个公开栏目的链接 |
| `{cms:primaryslide}` | 第一张轮播图图片 |
| `{cms:primaryslidetitle}` | 第一张轮播图标题 |
| `{cms:pagetitle}` | 当前页面标题 |
| `{cms:pagekeywords}` | 当前页面关键词 |
| `{cms:pagedescription}` | 当前页面描述 |

### 栏目标签

当前栏目可直接使用`{sort:*}`标签：

| 标签 | 输出 |
| --- | --- |
| `{sort:name}` | 当前栏目名称 |
| `{sort:link}` | 当前栏目链接 |
| `{sort:scode}` | 当前栏目编号 |
| `{sort:tcode}` | 当前顶级栏目编号 |
| `{sort:description}` | 当前栏目描述 |

在`cms:nav`、`cms:2nav`、`cms:3nav`和`cms:sort`循环内使用方括号标签：

| 标签 | 输出 |
| --- | --- |
| `[nav:link]`、`[2nav:link]`、`[3nav:link]`、`[sort:link]` | 栏目链接 |
| `[nav:scode]`、`[2nav:scode]`、`[3nav:scode]`、`[sort:scode]` | 栏目编号 |
| `[nav:name]`、`[2nav:name]`、`[3nav:name]`、`[sort:name]` | 栏目名称 |
| `[nav:title]`、`[sort:title]` | 栏目标题 |
| `[nav:keywords]`、`[sort:keywords]` | 栏目关键词 |
| `[nav:description]`、`[sort:description]` | 栏目描述 |
| `[nav:soncount]`、`[sort:soncount]` | 子栏目数量 |
| `[nav:active]`、`[sort:active]` | 当前栏目或其父栏目是否处于激活状态 |

常用导航示例：

```html
{cms:nav num=15 parent=0}
<li>
  <a href="[nav:link]" class="{cms:if('[nav:scode]'=='{sort:tcode}')}active{/cms:if}">[nav:name]</a>
  {cms:if([nav:soncount]>0)}
  <ul>
    {cms:2nav parent=[nav:scode]}
    <li><a href="[2nav:link]">[2nav:name]</a></li>
    {/cms:2nav}
  </ul>
  {/cms:if}
</li>
{/cms:nav}
```

### 文章列表标签

`cms:list`支持以下属性：

| 属性 | 说明 |
| --- | --- |
| `scode` | 栏目编号；不指定时使用当前栏目 |
| `num` | 输出数量，最大`100` |
| `order` | 排序方式，支持`id`、`date`、`sorting`、`visits` |

`cms:list`循环内可用标签：

| 标签 | 输出 |
| --- | --- |
| `[list:id]` | 文章 ID |
| `[list:i]` | 当前序号，从`1`开始 |
| `[list:link]` | 文章详情链接 |
| `[list:title]` | 文章标题 |
| `[list:subtitle]` | 文章副标题 |
| `[list:description]` | 文章摘要 |
| `[list:content]` | 文章摘要 |
| `[list:ico]` | 封面图地址 |
| `[list:date]` | 发布时间 |
| `[list:visits]` | 浏览量 |
| `[list:sortname]` | 所属栏目名称 |

标题、摘要等文本标签支持`len`、`lencn`和`more`参数：

```html
{cms:list scode=19 num=6 order=date}
<a href="[list:link]">[list:title lencn=34]</a>
<span>[list:date style=Y-m-d]</span>
{/cms:list}
```

### 搜索结果标签

公开搜索页位于`/cms-site/search`，模板文件为`public/templates/search.html`。
表单应使用`GET`方式提交到`{cms:scaction}`，关键词参数名为`keyword`：

```html
<form action="{cms:scaction}" method="get">
  <input type="text" name="keyword" value="{cms:keyword}">
  <button type="submit">搜索</button>
</form>
```

搜索会在全部已发布文章中检索标题、副标题、摘要、正文、标签、SEO 关键词和
SEO 描述。搜索结果页不限定当前栏目；如需限定栏目，可在模板循环上使用
`scode`参数。

`cms:search`支持以下属性：

| 属性 | 说明 |
| --- | --- |
| `scode` | 栏目编号；不指定时搜索全部栏目 |
| `num` | 每页输出数量，最大`100` |
| `order` | 排序方式，支持`id`、`date`、`sorting`、`visits` |

`cms:search`循环内可用标签与`cms:list`一致，只是前缀为`search`：

| 标签 | 输出 |
| --- | --- |
| `[search:id]` | 文章 ID |
| `[search:i]` | 当前序号，从`1`开始 |
| `[search:link]` | 文章详情链接 |
| `[search:title]` | 文章标题 |
| `[search:subtitle]` | 文章副标题 |
| `[search:description]` | 文章摘要 |
| `[search:content]` | 文章摘要 |
| `[search:preview]` | 命中关键词的安全预览片段 |
| `[search:ico]` | 封面图地址 |
| `[search:date]` | 发布时间 |
| `[search:visits]` | 浏览量 |
| `[search:sortname]` | 所属栏目名称 |

```html
{cms:search num=10 order=date}
<a href="[search:link]">[search:title lencn=48]</a>
<p>[search:preview]</p>
<span>[search:date style=Y-m-d]</span>
{/cms:search}
```

### 文章详情标签

详情页可用以下`{content:*}`标签：

| 标签 | 输出 |
| --- | --- |
| `{content:title}` | 文章标题 |
| `{content:subtitle}` | 文章副标题 |
| `{content:description}` | 文章摘要 |
| `{content:ico}` | 封面图地址 |
| `{content:author}` | 作者 |
| `{content:source}` | 来源 |
| `{content:date}` | 发布时间 |
| `{content:visits}` | 浏览量 |
| `{content:content}` | 正文 HTML |
| `{content:precontent}` | 上一篇链接 |
| `{content:nextcontent}` | 下一篇链接 |

`{content:date style=Y-m-d}`可用于兼容模板写法，当前输出与`{content:date}`一致。

### 轮播图与友情链接标签

轮播图和友情链接数据在后台 CMS 管理页的“轮播图”和“友情链接”页签维护。
标签中的`gid`对应后台表单里的“分组编码”，例如首页轮播默认使用
`gid=1`，页脚友情链接可以按`gid=1`、`gid=2`等分组展示。

`cms:slide`支持以下属性：

| 属性 | 说明 |
| --- | --- |
| `num` | 输出数量，最大`100` |
| `gid` | 分组编码；不指定时读取全部已启用轮播图 |

`cms:slide`循环内可用标签：

| 标签 | 输出 |
| --- | --- |
| `[slide:i]` | 当前序号，从`1`开始 |
| `[slide:link]` | 跳转链接 |
| `[slide:src]`、`[slide:ico]` | 轮播图片地址 |
| `[slide:title]` | 轮播标题 |
| `[slide:subtitle]` | 轮播副标题 |

`cms:link`支持以下属性：

| 属性 | 说明 |
| --- | --- |
| `num` | 输出数量，最大`100` |
| `gid` | 分组编码；不指定时读取全部已启用友情链接 |

`cms:link`循环内可用标签：

| 标签 | 输出 |
| --- | --- |
| `[link:link]` | 链接地址 |
| `[link:name]` | 链接名称 |
| `[link:logo]` | 链接标识图片地址 |

示例：

```html
{cms:slide num=5 gid=1}
<a href="[slide:link]"><img src="[slide:src]" alt="[slide:title]"></a>
{/cms:slide}

{cms:link num=10 gid=1}
<a href="[link:link]" target="_blank">[link:name]</a>
{/cms:link}
```

### 分页、面包屑与留言状态

| 标签 | 输出 |
| --- | --- |
| `{cms:position separator='&gt;'}` | 面包屑导航 |
| `{page:rows}` | 当前列表总记录数 |
| `{page:index}` | 首页分页链接 |
| `{page:pre}` | 上一页链接 |
| `{page:next}` | 下一页链接 |
| `{page:last}` | 末页链接 |
| `{page:numbar}` | 数字分页 HTML |
| `{cms:submitted}` | 留言提交成功状态 |
| `{cms:invalidmessage}` | 留言参数非法状态 |
| `{cms:messageerror}` | 留言提交失败状态 |

### 条件标签

条件标签只支持模板编译器内置的表达式，不支持任意脚本。常用写法如下：

```html
{cms:if(0=='{sort:scode}')}active{/cms:if}
{cms:if([list:ico])}<img src="[list:ico]" alt="[list:title]">{/cms:if}
{cms:if([nav:soncount]>0)}...{/cms:if}
{cms:if({page:rows}>0)}...{else}暂无内容{/cms:if}
{cms:if({cms:submitted})}<p>提交成功</p>{/cms:if}
```

## 开发

该插件按独立源码插件生成与测试。数据库代码生成使用
`backend/hack/config.yaml`，生成的 DAO、DO、Entity 文件均保留在插件自己的
`backend/internal`目录下。
