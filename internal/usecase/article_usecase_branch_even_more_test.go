package usecase

import (
	"context"
	"testing"
	"time"
	"write_base/internal/domain"
	"write_base/internal/mocks"
)

func TestGetArticleByID_AsAuthorSuccess(t *testing.T) {
	repo := &mocks.ArticleRepositoryMock{GetByIDFn: func(ctx context.Context, id string) (*domain.Article, error) {
		return &domain.Article{ID: id, AuthorID: "u1", Status: domain.StatusDraft}, nil
	}}
	pol := &mocks.PolicyMock{UserExistsFn: func(string) bool { return true }}
	uc := &ArticleUsecase{Repo: repo, Policy: pol, Utils: &mocks.UtilsMock{}, TagUsecase: &mocks.TagUsecaseMock{}, ViewUsecase: &mocks.ViewUsecaseMock{}, ClapUsecase: &mocks.ClapUsecaseMock{}}

	art, err := uc.GetArticleByID(context.Background(), "a1", "u1")
	if err != nil || art == nil || art.ID != "a1" {
		t.Fatalf("expected success, got %v, art=%v", err, art)
	}
}

func TestGetArticleByID_AsOtherPublished_RecordsView(t *testing.T) {
	viewed := false
	incremented := false
	repo := &mocks.ArticleRepositoryMock{GetByIDFn: func(ctx context.Context, id string) (*domain.Article, error) {
		return &domain.Article{ID: id, AuthorID: "author", Status: domain.StatusPublished}, nil
	}, IncrementViewFn: func(ctx context.Context, id string) error { incremented = true; return nil }}
	pol := &mocks.PolicyMock{UserExistsFn: func(string) bool { return true }}
	view := &mocks.ViewUsecaseMock{RecordViewFn: func(ctx context.Context, uid, aid, ip string) error { viewed = true; return nil }}
	uc := &ArticleUsecase{Repo: repo, Policy: pol, Utils: &mocks.UtilsMock{}, TagUsecase: &mocks.TagUsecaseMock{}, ViewUsecase: view, ClapUsecase: &mocks.ClapUsecaseMock{}}

	art, err := uc.GetArticleByID(context.Background(), "a2", "someone")
	if err != nil || art == nil {
		t.Fatalf("expected success, got %v", err)
	}
	if !viewed || !incremented {
		t.Fatalf("expected view recorded and incremented, got viewed=%v incremented=%v", viewed, incremented)
	}
}

func TestFilterArticles_DefaultsToPublished(t *testing.T) {
	sawPublished := false
	repo := &mocks.ArticleRepositoryMock{FilterFn: func(ctx context.Context, f domain.ArticleFilter, p domain.Pagination) ([]domain.Article, int, error) {
		// expect StatusPublished inserted when none provided
		for _, s := range f.Statuses {
			if s == domain.StatusPublished {
				sawPublished = true
			}
		}
		return []domain.Article{}, 0, nil
	}}
	uc := &ArticleUsecase{Repo: repo, Policy: &mocks.PolicyMock{}, Utils: &mocks.UtilsMock{}, TagUsecase: &mocks.TagUsecaseMock{}, ViewUsecase: &mocks.ViewUsecaseMock{}, ClapUsecase: &mocks.ClapUsecaseMock{}}

	_, _, err := uc.FilterArticles(context.Background(), domain.ArticleFilter{}, domain.Pagination{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !sawPublished {
		t.Fatalf("expected default published status in filter")
	}
}

func TestAdminListAllArticles_UnauthorizedRole(t *testing.T) {
	pol := &mocks.PolicyMock{UserExistsFn: func(string) bool { return true }, IsAdminFn: func(_, _ string) bool { return false }}
	uc := &ArticleUsecase{Repo: &mocks.ArticleRepositoryMock{}, Policy: pol, Utils: &mocks.UtilsMock{}, TagUsecase: &mocks.TagUsecaseMock{}, ViewUsecase: &mocks.ViewUsecaseMock{}, ClapUsecase: &mocks.ClapUsecaseMock{}}

	_, _, err := uc.AdminListAllArticles(context.Background(), "u1", "user", domain.Pagination{})
	if err == nil || err != domain.ErrUnauthorized {
		t.Fatalf("expected unauthorized, got %v", err)
	}
}

func TestDeleteArticleFromTrash_NotDeletedStatus(t *testing.T) {
	repo := &mocks.ArticleRepositoryMock{GetByIDFn: func(ctx context.Context, id string) (*domain.Article, error) {
		return &domain.Article{ID: id, AuthorID: "u1", Status: domain.StatusDraft}, nil
	}}
	pol := &mocks.PolicyMock{UserExistsFn: func(string) bool { return true }, UserOwnsArticleFn: func(uid string, a *domain.Article) bool { return a.AuthorID == uid }}
	uc := &ArticleUsecase{Repo: repo, Policy: pol, Utils: &mocks.UtilsMock{}, TagUsecase: &mocks.TagUsecaseMock{}, ViewUsecase: &mocks.ViewUsecaseMock{}, ClapUsecase: &mocks.ClapUsecaseMock{}}

	err := uc.DeleteArticleFromTrash(context.Background(), "a1", "u1")
	if err == nil || err != domain.ErrArticleNotFound {
		t.Fatalf("expected ErrArticleNotFound, got %v", err)
	}
}

func TestEmptyTrash_Unauthorized(t *testing.T) {
	pol := &mocks.PolicyMock{UserExistsFn: func(string) bool { return false }}
	uc := &ArticleUsecase{Repo: &mocks.ArticleRepositoryMock{}, Policy: pol, Utils: &mocks.UtilsMock{}, TagUsecase: &mocks.TagUsecaseMock{}, ViewUsecase: &mocks.ViewUsecaseMock{}, ClapUsecase: &mocks.ClapUsecaseMock{}}
	if err := uc.EmptyTrash(context.Background(), "u1"); err == nil || err != domain.ErrUnauthorized {
		t.Fatalf("expected unauthorized")
	}
}

func TestPublishArticle_AlreadyPublished(t *testing.T) {
	repo := &mocks.ArticleRepositoryMock{GetByIDFn: func(ctx context.Context, id string) (*domain.Article, error) {
		return &domain.Article{ID: id, AuthorID: "u1", Status: domain.StatusPublished}, nil
	}, PublishFn: func(ctx context.Context, id string, at time.Time) error { return nil }}
	pol := &mocks.PolicyMock{UserOwnsArticleFn: func(uid string, a *domain.Article) bool { return a.AuthorID == uid }}
	uc := &ArticleUsecase{Repo: repo, Policy: pol, Utils: &mocks.UtilsMock{}, TagUsecase: &mocks.TagUsecaseMock{}, ViewUsecase: &mocks.ViewUsecaseMock{}, ClapUsecase: &mocks.ClapUsecaseMock{}}

	_, err := uc.PublishArticle(context.Background(), "a1", "u1")
	if err == nil || err != domain.ErrArticlePublished {
		t.Fatalf("expected ErrArticlePublished, got %v", err)
	}
}

func TestArchiveArticle_AlreadyArchived(t *testing.T) {
	repo := &mocks.ArticleRepositoryMock{GetByIDFn: func(ctx context.Context, id string) (*domain.Article, error) {
		return &domain.Article{ID: id, AuthorID: "u1", Status: domain.StatusArchived}, nil
	}, ArchiveFn: func(ctx context.Context, id string, at time.Time) error { return nil }}
	pol := &mocks.PolicyMock{UserOwnsArticleFn: func(uid string, a *domain.Article) bool { return a.AuthorID == uid }}
	uc := &ArticleUsecase{Repo: repo, Policy: pol, Utils: &mocks.UtilsMock{}, TagUsecase: &mocks.TagUsecaseMock{}, ViewUsecase: &mocks.ViewUsecaseMock{}, ClapUsecase: &mocks.ClapUsecaseMock{}}

	_, err := uc.ArchiveArticle(context.Background(), "a1", "u1")
	if err == nil || err != domain.ErrArticleArchived {
		t.Fatalf("expected ErrArticleArchived, got %v", err)
	}
}

func TestFilterAuthorArticles_NotAuthor(t *testing.T) {
	repo := &mocks.ArticleRepositoryMock{}
	pol := &mocks.PolicyMock{UserExistsFn: func(string) bool { return true }}
	uc := &ArticleUsecase{Repo: repo, Policy: pol, Utils: &mocks.UtilsMock{}, TagUsecase: &mocks.TagUsecaseMock{}, ViewUsecase: &mocks.ViewUsecaseMock{}, ClapUsecase: &mocks.ClapUsecaseMock{}}

	_, _, err := uc.FilterAuthorArticles(context.Background(), "caller", "author", domain.ArticleFilter{}, domain.Pagination{})
	if err == nil || err != domain.ErrUnauthorized {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestFilterAuthorArticles_Success(t *testing.T) {
	repo := &mocks.ArticleRepositoryMock{FilterAuthorArticlesFn: func(ctx context.Context, authorID string, f domain.ArticleFilter, p domain.Pagination) ([]domain.Article, int, error) {
		return []domain.Article{}, 0, nil
	}}
	pol := &mocks.PolicyMock{UserExistsFn: func(string) bool { return true }}
	uc := &ArticleUsecase{Repo: repo, Policy: pol, Utils: &mocks.UtilsMock{}, TagUsecase: &mocks.TagUsecaseMock{}, ViewUsecase: &mocks.ViewUsecaseMock{}, ClapUsecase: &mocks.ClapUsecaseMock{}}

	_, _, err := uc.FilterAuthorArticles(context.Background(), "u1", "u1", domain.ArticleFilter{}, domain.Pagination{Page: 0, PageSize: 0})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}
