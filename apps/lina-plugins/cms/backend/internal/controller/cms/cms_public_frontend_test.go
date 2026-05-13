// This file verifies CMS public frontend template compilation behavior.

package cms

import (
	"bytes"
	"html/template"
	"strconv"
	"strings"
	"testing"

	entitymodel "lina-plugin-cms/backend/internal/model/entity"
	cmssvc "lina-plugin-cms/backend/internal/service/cms"
)

// TestPublicFrontendListNumLimitsRenderedItems verifies cms:list num limits
// rendered article items instead of falling back to the default page size.
func TestPublicFrontendListNumLimitsRenderedItems(t *testing.T) {
	const source = `{cms:list num=5 order=date}<span data-testid="item">[list:title]</span>{/cms:list}`

	compiled := compilePublicFrontendTemplate(source, publicFrontendRootScope)
	tpl, err := template.New("case").Funcs(template.FuncMap{
		"cmsLimit":         publicFrontendLimit,
		"cmsOrderArticles": publicFrontendOrderArticles,
	}).Parse(compiled)
	if err != nil {
		t.Fatalf("parse compiled public frontend template: %v", err)
	}

	articles := make([]*publicFrontendArticle, 0, 6)
	for index := 1; index <= 6; index++ {
		articles = append(articles, &publicFrontendArticle{Title: "Article " + strconv.Itoa(index)})
	}

	var buffer bytes.Buffer
	if err = tpl.Execute(&buffer, &publicFrontendView{Articles: articles}); err != nil {
		t.Fatalf("execute compiled public frontend template: %v", err)
	}
	html := buffer.String()
	if count := strings.Count(html, `data-testid="item"`); count != 5 {
		t.Fatalf("expected 5 rendered articles, got %d in %s", count, html)
	}
	if strings.Contains(html, "Article 6") {
		t.Fatalf("expected article beyond num limit to be hidden, got %s", html)
	}
}

// TestPublicFrontendSearchLoopCompiles verifies cms:search renders the public
// search result article tags with the same article projection as list loops.
func TestPublicFrontendSearchLoopCompiles(t *testing.T) {
	const source = `{cms:search num=2 order=date}<a href="[search:link]" data-testid="item">[search:title]<em>[search:preview]</em></a>{/cms:search}`

	compiled := compilePublicFrontendTemplate(source, publicFrontendRootScope)
	tpl, err := template.New("case").Funcs(template.FuncMap{
		"cmsLimit":         publicFrontendLimit,
		"cmsOrderArticles": publicFrontendOrderArticles,
	}).Parse(compiled)
	if err != nil {
		t.Fatalf("parse compiled public frontend search template: %v", err)
	}

	articles := []*publicFrontendArticle{
		{Title: "First", Href: "/cms-site?article=first", SearchPreview: template.HTML(`命中<mark>关键词</mark>`)},
		{Title: "Second", Href: "/cms-site?article=second", SearchPreview: template.HTML(`第二条预览`)},
		{Title: "Third", Href: "/cms-site?article=third"},
	}
	var buffer bytes.Buffer
	if err = tpl.Execute(&buffer, &publicFrontendView{Articles: articles}); err != nil {
		t.Fatalf("execute compiled public frontend search template: %v", err)
	}
	html := buffer.String()
	if count := strings.Count(html, `data-testid="item"`); count != 2 {
		t.Fatalf("expected 2 rendered search results, got %d in %s", count, html)
	}
	if strings.Contains(html, "Third") {
		t.Fatalf("expected search result beyond num limit to be hidden, got %s", html)
	}
	if !strings.Contains(html, `<mark>关键词</mark>`) {
		t.Fatalf("expected search preview markup to render, got %s", html)
	}
}

// TestPublicFrontendSearchActionTargetsSearchPage verifies public templates
// submit keyword searches to the dedicated search result page.
func TestPublicFrontendSearchActionTargetsSearchPage(t *testing.T) {
	compiled := compilePublicFrontendTemplate(`<form action="{cms:scaction}"></form>`, publicFrontendRootScope)
	if !strings.Contains(compiled, `action="/cms-site/search"`) {
		t.Fatalf("expected search action to target /cms-site/search, got %s", compiled)
	}
}

// TestPublicFrontendSearchPageTitle verifies search pages expose a useful
// document title with and without a keyword.
func TestPublicFrontendSearchPageTitle(t *testing.T) {
	if got, want := publicFrontendSearchPageTitle("算力"), "算力-搜索结果"; got != want {
		t.Fatalf("expected keyword search page title %q, got %q", want, got)
	}
	if got, want := publicFrontendSearchPageTitle("  "), "搜索结果"; got != want {
		t.Fatalf("expected empty search page title %q, got %q", want, got)
	}
}

// TestPublicFrontendEmbeddedTemplatesParse verifies embedded CMS public
// templates compile after template files add or change tags.
func TestPublicFrontendEmbeddedTemplatesParse(t *testing.T) {
	tpl, err := publicFrontendTemplate()
	if err != nil {
		t.Fatalf("parse embedded public frontend templates: %v", err)
	}
	if tpl.Lookup(publicFrontendSearchName) == nil {
		t.Fatalf("expected embedded public frontend template %q to be registered", publicFrontendSearchName)
	}
}

// TestPublicFrontendSearchPreviewBuildsSafeHighlightedExcerpt verifies public
// search snippets are plain-text, clipped around the body hit, and highlighted.
func TestPublicFrontendSearchPreviewBuildsSafeHighlightedExcerpt(t *testing.T) {
	item := &cmssvc.ArticleItem{
		CmsArticle: &entitymodel.CmsArticle{
			Title:   "产业动态",
			Summary: "不包含目标词的摘要",
			Content: `<p>这一段正文包含算力网络建设的阶段成果。</p><script>alert(1)</script>`,
		},
	}

	preview := string(publicFrontendSearchPreview(item, "算力网络"))
	if !strings.Contains(preview, `<mark>算力网络</mark>`) {
		t.Fatalf("expected highlighted keyword in search preview, got %s", preview)
	}
	if strings.Contains(preview, "<script") || strings.Contains(preview, "alert") {
		t.Fatalf("expected script content to be stripped from search preview, got %s", preview)
	}
	if strings.Contains(preview, "<p>") {
		t.Fatalf("expected HTML tags to be stripped from search preview, got %s", preview)
	}
}

// TestPublicFrontendStaticAssetPathNormalization verifies seeded static media
// paths resolve through the CMS public asset URL space.
func TestPublicFrontendStaticAssetPathNormalization(t *testing.T) {
	if got, want := publicFrontendNormalizeAssetPath("/static/logo.svg"), "/cms-site/assets/static/logo.svg"; got != want {
		t.Fatalf("expected normalized static path %q, got %q", want, got)
	}
	if got, want := publicFrontendNormalizeContentHTML(`<img src="/static/wechat.jpg">`), `<img src="/cms-site/assets/static/wechat.jpg">`; got != want {
		t.Fatalf("expected normalized static content %q, got %q", want, got)
	}
}
