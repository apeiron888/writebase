package policy

import (
	"testing"
	"write_base/internal/domain"
	imocks "write_base/internal/mocks"
)

func TestArticleCreateValid(t *testing.T) {
	utils := &imocks.UtilsMock{ValidateContentFn: func(_ []domain.ContentBlock) bool { return true }}
	p := NewArticlePolicy(utils)
	a := &domain.Article{Title: "ok", Tags: []string{"go"}, ContentBlocks: []domain.ContentBlock{{Type: domain.BlockParagraph, Content: domain.BlockContent{Paragraph: &domain.ParagraphContent{Text: "x"}}}}}
	if !p.ArticleCreateValid(a) {
		t.Fatalf("expected valid article")
	}
}

func TestArticleCreateValid_Invalid(t *testing.T) {
	utils := &imocks.UtilsMock{ValidateContentFn: func(_ []domain.ContentBlock) bool { return false }}
	p := NewArticlePolicy(utils)
	a := &domain.Article{Title: "", Tags: []string{}, ContentBlocks: nil}
	if p.ArticleCreateValid(a) {
		t.Fatalf("expected invalid article")
	}
}

func TestUserOwnsArticle(t *testing.T) {
	p := NewArticlePolicy(&imocks.UtilsMock{})
	a := &domain.Article{AuthorID: "u1"}
	if !p.UserOwnsArticle("u1", a) {
		t.Fatalf("owner should match")
	}
	if p.UserOwnsArticle("u2", a) {
		t.Fatalf("owner mismatch should be false")
	}
}

func TestCheckArticleChangesAndValid(t *testing.T) {
	p := NewArticlePolicy(&imocks.UtilsMock{})
	old := &domain.Article{ID: "a1", AuthorID: "u1", Title: "t1", Excerpt: "e", Language: "en", Slug: "t1"}
	neu := &domain.Article{ID: "a1", AuthorID: "u1", Title: "t2", Excerpt: "e", Language: "en", Slug: "t2"}
	if !p.CheckArticleChangesAndValid(old, neu) {
		t.Fatalf("title change should be allowed and true")
	}

	// invalid author change
	neu.AuthorID = "other"
	if p.CheckArticleChangesAndValid(old, neu) {
		t.Fatalf("author change should be invalid")
	}
}

func TestIsAdmin(t *testing.T) {
	p := NewArticlePolicy(&imocks.UtilsMock{})
	if !p.IsAdmin("u1", "admin") {
		t.Fatalf("admin should pass")
	}
	if p.IsAdmin("u1", "user") {
		t.Fatalf("non-admin should fail")
	}
}
