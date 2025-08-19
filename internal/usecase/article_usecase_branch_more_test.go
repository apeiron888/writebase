package usecase

import (
	"context"
	"testing"
	"time"
	"write_base/internal/domain"
	"write_base/internal/mocks"
)

func TestGetArticleByID_UnauthorizedWhenNotAuthorAndNotPublished(t *testing.T) {
	repo := &mocks.ArticleRepositoryMock{GetByIDFn: func(ctx context.Context, id string) (*domain.Article, error) {
		return &domain.Article{ID: id, AuthorID: "author-1", Status: domain.StatusDraft}, nil
	}}
	pol := &mocks.PolicyMock{UserExistsFn: func(string) bool { return true }}
	uc := &ArticleUsecase{Repo: repo, Policy: pol, Utils: &mocks.UtilsMock{}, TagUsecase: &mocks.TagUsecaseMock{}, ViewUsecase: &mocks.ViewUsecaseMock{}, ClapUsecase: &mocks.ClapUsecaseMock{}}

	_, err := uc.GetArticleByID(context.Background(), "a1", "other-user")
	if err == nil || err != domain.ErrUnauthorized {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestPublishArticle_UnapprovedTags(t *testing.T) {
	repo := &mocks.ArticleRepositoryMock{GetByIDFn: func(ctx context.Context, id string) (*domain.Article, error) {
		return &domain.Article{ID: id, AuthorID: "u1", Status: domain.StatusDraft, Tags: []string{"x"}}, nil
	}, PublishFn: func(ctx context.Context, id string, at time.Time) error { return nil }}
	pol := &mocks.PolicyMock{UserOwnsArticleFn: func(uid string, a *domain.Article) bool { return a.AuthorID == uid }}
	tag := &mocks.TagUsecaseMock{IsTagApprovedFn: func(name string) bool { return false }}
	uc := &ArticleUsecase{Repo: repo, Policy: pol, Utils: &mocks.UtilsMock{}, TagUsecase: tag, ViewUsecase: &mocks.ViewUsecaseMock{}, ClapUsecase: &mocks.ClapUsecaseMock{}}

	_, err := uc.PublishArticle(context.Background(), "a1", "u1")
	if err == nil || err != domain.ErrUnapprovedTags {
		t.Fatalf("expected ErrUnapprovedTags, got %v", err)
	}
}

func TestUnpublishArticle_NotPublished(t *testing.T) {
	repo := &mocks.ArticleRepositoryMock{GetByIDFn: func(ctx context.Context, id string) (*domain.Article, error) {
		return &domain.Article{ID: id, AuthorID: "u1", Status: domain.StatusDraft}, nil
	}, UnpublishFn: func(ctx context.Context, id string) error { return nil }}
	pol := &mocks.PolicyMock{UserOwnsArticleFn: func(uid string, a *domain.Article) bool { return a.AuthorID == uid }}
	uc := &ArticleUsecase{Repo: repo, Policy: pol, Utils: &mocks.UtilsMock{}, TagUsecase: &mocks.TagUsecaseMock{}, ViewUsecase: &mocks.ViewUsecaseMock{}, ClapUsecase: &mocks.ClapUsecaseMock{}}

	_, err := uc.UnpublishArticle(context.Background(), "a1", "u1")
	if err == nil || err != domain.ErrArticleNotPublished {
		t.Fatalf("expected ErrArticleNotPublished, got %v", err)
	}
}
