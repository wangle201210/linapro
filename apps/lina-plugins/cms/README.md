# LinaPro CMS Plugin

The CMS source plugin provides content site capabilities for LinaPro. It owns
site settings, category trees, articles, links, slides, visitor messages, and
public HTML rendering APIs.

## Scope

- Source plugin under `apps/lina-plugins/cms`
- Plugin-owned `plugin_cms_*` tables
- Authenticated management APIs under `/api/v1/cms/*`
- Public read-only APIs under `/api/v1/cms/public/*`
- Public HTML pages under `/cms-site`
- Vben management page embedded through the plugin runtime
- Plugin-owned E2E tests under `hack/tests`

## Public Templates

Public templates live in `public/templates`. They use CMS-owned template tags
with the `cms:` prefix. Static assets used by the reference site live in
`public/assets`.

### Template Files

| File | Purpose |
| --- | --- |
| `partials.html` | Shared `head`, navigation, sidebar, and footer fragments |
| `index.html` | Public home page |
| `list.html` | Standard article list page |
| `list-card.html` | Card-style article list page |
| `search.html` | Public search result page |
| `detail.html` | Article detail page |
| `single.html` | Single-page category |
| `message.html` | Visitor message page |

### Basic Syntax

| Syntax | Description |
| --- | --- |
| `{include file=comm/head.html}` | Includes a shared fragment |
| `{cms:if(condition)}...{/cms:if}` | Conditional block with `{else}` |
| `{cms:list num=10 order=date}...{/cms:list}` | Iterates articles |
| `{cms:search num=10 order=date}...{/cms:search}` | Iterates search results |
| `{cms:nav num=15 parent=0}...{/cms:nav}` | Iterates categories |
| `{cms:slide num=5 gid=1}...{/cms:slide}` | Iterates slides |
| `{cms:link num=10 gid=1}...{/cms:link}` | Iterates friendly links |
| `{cms:sort scode=25}...{/cms:sort}` | Loads one category by code |

Supported include names are `head.html`, `foot.html`, `sidebar.html`, and
`sortnav.html`.

Loop `num` is capped at `100`. `cms:list` defaults to `12` items per page when
`num` is omitted.

### Global Site Tags

| Tag | Output |
| --- | --- |
| `{cms:sitepath}` | Public site root path, fixed to `/cms-site` |
| `{cms:sitetplpath}` | Public asset path, fixed to `/cms-site/assets` |
| `{cms:scaction}` | Search form action, fixed to `/cms-site/search` |
| `{cms:msgpage}` | Message page URL |
| `{cms:msgaction}` | Message submit URL |
| `{cms:sitetitle}`, `{cms:companyname}` | Site name |
| `{cms:sitesubtitle}`, `{cms:siteslogan}` | Site slogan |
| `{cms:sitelogo}` | Site logo URL |
| `{cms:sitekeywords}` | Site keywords |
| `{cms:sitedescription}` | Site description |
| `{cms:companyaddress}` | Contact address |
| `{cms:companyphone}` | Contact phone |
| `{cms:companyemail}` | Contact email |
| `{cms:companycontact}` | Contact person |
| `{cms:siteicp}` | ICP record text |
| `{cms:companyweixin}` | WeChat QR code URL |
| `{cms:keyword}` | Current search keyword |
| `{cms:year}` | Current year |
| `{cms:firstsortlink}` | First public category link |
| `{cms:primaryslide}` | First slide image |
| `{cms:primaryslidetitle}` | First slide title |
| `{cms:pagetitle}` | Current page title |
| `{cms:pagekeywords}` | Current page keywords |
| `{cms:pagedescription}` | Current page description |

### Category Tags

Use `{sort:*}` tags for the current category:

| Tag | Output |
| --- | --- |
| `{sort:name}` | Current category name |
| `{sort:link}` | Current category link |
| `{sort:scode}` | Current category code |
| `{sort:tcode}` | Current top-level category code |
| `{sort:description}` | Current category description |

Use bracket tags inside `cms:nav`, `cms:2nav`, `cms:3nav`, and `cms:sort`
loops:

| Tag | Output |
| --- | --- |
| `[*:link]` | Category link |
| `[*:scode]` | Category code |
| `[*:name]` | Category name |
| `[nav:title]`, `[sort:title]` | Category title |
| `[nav:keywords]`, `[sort:keywords]` | Category keywords |
| `[nav:description]`, `[sort:description]` | Category description |
| `[nav:soncount]`, `[sort:soncount]` | Child category count |
| `[nav:active]`, `[sort:active]` | Active category flag |

`*` represents the active loop prefix, such as `nav`, `2nav`, `3nav`, or
`sort`.

Common navigation example:

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

### Article List Tags

`cms:list` supports these attributes:

| Attribute | Description |
| --- | --- |
| `scode` | Category code. When omitted, the current category is used |
| `num` | Item count, capped at `100` |
| `order` | Sort order: `id`, `date`, `sorting`, or `visits` |

Tags available inside `cms:list`:

| Tag | Output |
| --- | --- |
| `[list:id]` | Article ID |
| `[list:i]` | One-based item index |
| `[list:link]` | Article detail link |
| `[list:title]` | Article title |
| `[list:subtitle]` | Article subtitle |
| `[list:description]` | Article summary |
| `[list:content]` | Article summary |
| `[list:ico]` | Cover image URL |
| `[list:date]` | Publish time |
| `[list:visits]` | View count |
| `[list:sortname]` | Category name |

Text tags such as title and description support `len`, `lencn`, and `more` modifiers:

```html
{cms:list scode=19 num=6 order=date}
<a href="[list:link]">[list:title lencn=34]</a>
<span>[list:date style=Y-m-d]</span>
{/cms:list}
```

### Search Result Tags

The public search page is `/cms-site/search`, and its template is
`public/templates/search.html`. Search forms should submit with `GET` to
`{cms:scaction}` and use `keyword` as the query parameter:

```html
<form action="{cms:scaction}" method="get">
  <input type="text" name="keyword" value="{cms:keyword}">
  <button type="submit">Search</button>
</form>
```

Search matches all published articles by title, subtitle, summary, body, tags,
SEO keywords, and SEO description. The search page is not limited to the current
category. Use the loop `scode` attribute when a template needs a category
restriction.

`cms:search` supports these attributes:

| Attribute | Description |
| --- | --- |
| `scode` | Category code. When omitted, all categories are searched |
| `num` | Items per page, capped at `100` |
| `order` | Sort order: `id`, `date`, `sorting`, or `visits` |

Tags available inside `cms:search` match `cms:list` with the `search` prefix:

| Tag | Output |
| --- | --- |
| `[search:id]` | Article ID |
| `[search:i]` | One-based item index |
| `[search:link]` | Article detail link |
| `[search:title]` | Article title |
| `[search:subtitle]` | Article subtitle |
| `[search:description]` | Article summary |
| `[search:content]` | Article summary |
| `[search:preview]` | Safe excerpt around the matched keyword |
| `[search:ico]` | Cover image URL |
| `[search:date]` | Publish time |
| `[search:visits]` | View count |
| `[search:sortname]` | Category name |

```html
{cms:search num=10 order=date}
<a href="[search:link]">[search:title lencn=48]</a>
<p>[search:preview]</p>
<span>[search:date style=Y-m-d]</span>
{/cms:search}
```

### Article Detail Tags

Detail pages support these `{content:*}` tags:

| Tag | Output |
| --- | --- |
| `{content:title}` | Article title |
| `{content:subtitle}` | Article subtitle |
| `{content:description}` | Article summary |
| `{content:ico}` | Cover image URL |
| `{content:author}` | Author |
| `{content:source}` | Source |
| `{content:date}` | Publish time |
| `{content:visits}` | View count |
| `{content:content}` | Body HTML |
| `{content:precontent}` | Previous article link |
| `{content:nextcontent}` | Next article link |

`{content:date style=Y-m-d}` is accepted for template compatibility. It
currently renders the same value as `{content:date}`.

### Slide and Link Tags

Carousel slides and friendly links are maintained in the CMS management page
through the "Slides" and "Links" tabs. The `gid` attribute maps to the
management form's group code. For example, the homepage carousel can use
`gid=1`, while footer links can be rendered as several `gid=1`, `gid=2`, and
similar groups.

`cms:slide` supports these attributes:

| Attribute | Description |
| --- | --- |
| `num` | Item count, capped at `100` |
| `gid` | Group code. When omitted, all enabled slides are used |

Tags available inside `cms:slide`:

| Tag | Output |
| --- | --- |
| `[slide:i]` | One-based item index |
| `[slide:link]` | Target link |
| `[slide:src]`, `[slide:ico]` | Slide image URL |
| `[slide:title]` | Slide title |
| `[slide:subtitle]` | Slide subtitle |

`cms:link` supports these attributes:

| Attribute | Description |
| --- | --- |
| `num` | Item count, capped at `100` |
| `gid` | Group code. When omitted, all enabled friendly links are used |

Tags available inside `cms:link`:

| Tag | Output |
| --- | --- |
| `[link:link]` | Link URL |
| `[link:name]` | Link name |
| `[link:logo]` | Link logo URL |

Example:

```html
{cms:slide num=5 gid=1}
<a href="[slide:link]"><img src="[slide:src]" alt="[slide:title]"></a>
{/cms:slide}

{cms:link num=10 gid=1}
<a href="[link:link]" target="_blank">[link:name]</a>
{/cms:link}
```

### Pagination, Breadcrumb, and Message State

| Tag | Output |
| --- | --- |
| `{cms:position separator='&gt;'}` | Breadcrumb navigation |
| `{page:rows}` | Total rows in the current list |
| `{page:index}` | First page link |
| `{page:pre}` | Previous page link |
| `{page:next}` | Next page link |
| `{page:last}` | Last page link |
| `{page:numbar}` | Numeric pagination HTML |
| `{cms:submitted}` | Message submitted successfully |
| `{cms:invalidmessage}` | Message payload is invalid |
| `{cms:messageerror}` | Message submit failed |

### Conditional Tags

Condition tags only support expressions recognized by the template compiler.
They do not execute arbitrary scripts. Common forms:

```html
{cms:if(0=='{sort:scode}')}active{/cms:if}
{cms:if([list:ico])}<img src="[list:ico]" alt="[list:title]">{/cms:if}
{cms:if([nav:soncount]>0)}...{/cms:if}
{cms:if({page:rows}>0)}...{else}No content{/cms:if}
{cms:if({cms:submitted})}<p>Submitted</p>{/cms:if}
```

## Development

The plugin is generated and tested as an independent source plugin. Database
code generation uses `backend/hack/config.yaml`, and the generated DAO, DO, and
Entity files stay inside the plugin `backend/internal` tree.
