// This file verifies CMS public frontend template compilation behavior.

package cms

import (
	"bytes"
	"html/template"
	"strconv"
	"strings"
	"testing"
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
