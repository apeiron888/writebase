package usecase_test

import (
	"context"
	"testing"
	"time"
	"write_base/internal/domain"
	"write_base/internal/mocks"
	"write_base/internal/usecase"

	"github.com/stretchr/testify/require"
)

func TestArticleUsecase_CreateArticle_Success(t *testing.T) {
	repo := &mocks.ArticleRepositoryMock{}
	policy := &mocks.PolicyMock{}
	utils := &mocks.UtilsMock{GenerateUUIDFn: func() string { return "aid-1" }, GenerateSlugFn: func(_ string) string { return "hello" }}
	tagUC := &mocks.TagUsecaseMock{}
	viewUC := &mocks.ViewUsecaseMock{}
	clapUC := &mocks.ClapUsecaseMock{}

	// expectations
	policy.ArticleCreateValidFn = func(a *domain.Article) bool { return true }
	repo.GetBySlugFn = func(ctx context.Context, slug string) (*domain.Article, error) { return nil, domain.ErrArticleNotFound }
	repo.CreateFn = func(ctx context.Context, a *domain.Article) error { return nil }

	uc := usecase.NewArticleUsecase(repo, policy, utils, tagUC, viewUC, clapUC, nil)

	input := &domain.Article{Title: "Hello", Tags: []string{"go"}, ContentBlocks: []domain.ContentBlock{{Type: domain.BlockParagraph, Order: 0, Content: domain.BlockContent{Paragraph: &domain.ParagraphContent{Text: "hi"}}}}}
	id, err := uc.CreateArticle(context.Background(), "u1", input)

	require.NoError(t, err)
	require.Equal(t, "aid-1", id)
}

func TestArticleUsecase_PublishArticle_RequiresApprovedTags(t *testing.T) {
	repo := &mocks.ArticleRepositoryMock{}
	policy := &mocks.PolicyMock{UserOwnsArticleFn: func(uid string, a *domain.Article) bool { return true }}
	utils := &mocks.UtilsMock{}
	tagUC := &mocks.TagUsecaseMock{IsTagApprovedFn: func(name string) bool { return name == "go" }}
	viewUC := &mocks.ViewUsecaseMock{}
	clapUC := &mocks.ClapUsecaseMock{}

	art := &domain.Article{ID: "a1", AuthorID: "u1", Status: domain.StatusDraft, Tags: []string{"go"}}
	repo.GetByIDFn = func(ctx context.Context, id string) (*domain.Article, error) { return art, nil }
	repo.PublishFn = func(ctx context.Context, id string, at time.Time) error { return nil }

	uc := usecase.NewArticleUsecase(repo, policy, utils, tagUC, viewUC, clapUC, nil)

	out, err := uc.PublishArticle(context.Background(), "a1", "u1")
	require.NoError(t, err)
	require.Equal(t, domain.StatusPublished, out.Status)
}
