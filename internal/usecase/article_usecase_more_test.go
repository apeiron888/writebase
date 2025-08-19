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

func newArticleUC() (*usecase.ArticleUsecase, *mocks.ArticleRepositoryMock, *mocks.PolicyMock, *mocks.UtilsMock, *mocks.TagUsecaseMock, *mocks.ViewUsecaseMock, *mocks.ClapUsecaseMock) {
	repo := &mocks.ArticleRepositoryMock{}
	policy := &mocks.PolicyMock{UserExistsFn: func(string) bool { return true }, ArticleCreateValidFn: func(*domain.Article) bool { return true }, UserOwnsArticleFn: func(uid string, a *domain.Article) bool { return a.AuthorID == uid }}
	utils := &mocks.UtilsMock{GenerateUUIDFn: func() string { return "aid" }, GenerateSlugFn: func(t string) string { return "slug" }, GenerateShortUUIDFn: func() string { return "x1" }, ValidateContentFn: func([]domain.ContentBlock) bool { return true }}
	tagUC := &mocks.TagUsecaseMock{ValidateTagsFn: func([]string) error { return nil }, IsTagApprovedFn: func(string) bool { return true }}
	viewUC := &mocks.ViewUsecaseMock{}
	clapUC := &mocks.ClapUsecaseMock{}
	uc := usecase.NewArticleUsecase(repo, policy, utils, tagUC, viewUC, clapUC, nil).(*usecase.ArticleUsecase)
	return uc, repo, policy, utils, tagUC, viewUC, clapUC
}

func TestUpdateArticle_Success(t *testing.T) {
	uc, repo, _, _, _, _, _ := newArticleUC()
	old := &domain.Article{ID: "a1", AuthorID: "u1", Title: "T", ContentBlocks: []domain.ContentBlock{{Type: domain.BlockParagraph, Content: domain.BlockContent{Paragraph: &domain.ParagraphContent{Text: "x"}}}}, Tags: []string{"go"}}
	repo.GetByIDFn = func(ctx context.Context, id string) (*domain.Article, error) { return old, nil }
	repo.GetBySlugFn = func(ctx context.Context, slug string) (*domain.Article, error) { return nil, domain.ErrArticleNotFound }
	repo.UpdateFn = func(ctx context.Context, a *domain.Article) error { return nil }

	a := &domain.Article{ID: "a1", AuthorID: "u1", Title: "New", Tags: []string{"go"}, ContentBlocks: old.ContentBlocks}
	err := uc.UpdateArticle(context.Background(), "u1", a)
	require.NoError(t, err)
}

func TestDeleteArticle_Unauthorized(t *testing.T) {
	uc, repo, policy, _, _, _, _ := newArticleUC()
	policy.UserOwnsArticleFn = func(uid string, a *domain.Article) bool { return false }
	repo.GetByIDFn = func(ctx context.Context, id string) (*domain.Article, error) {
		return &domain.Article{ID: id, AuthorID: "owner"}, nil
	}
	err := uc.DeleteArticle(context.Background(), "a1", "attacker")
	// Current implementation maps underlying unauthorized from GetArticleByID to ErrArticleNotFound
	require.ErrorIs(t, err, domain.ErrArticleNotFound)
}

func TestDeleteArticle_Success(t *testing.T) {
	uc, repo, _, _, _, _, _ := newArticleUC()
	repo.GetByIDFn = func(ctx context.Context, id string) (*domain.Article, error) {
		return &domain.Article{ID: id, AuthorID: "u1"}, nil
	}
	repo.DeleteFn = func(ctx context.Context, id string) error { return nil }
	err := uc.DeleteArticle(context.Background(), "a1", "u1")
	require.NoError(t, err)
}

func TestUnpublish_NotPublished(t *testing.T) {
	uc, repo, _, _, _, _, _ := newArticleUC()
	repo.GetByIDFn = func(ctx context.Context, id string) (*domain.Article, error) {
		return &domain.Article{ID: id, AuthorID: "u1", Status: domain.StatusDraft}, nil
	}
	_, err := uc.UnpublishArticle(context.Background(), "a1", "u1")
	require.ErrorIs(t, err, domain.ErrArticleNotPublished)
}

func TestUnpublish_Success(t *testing.T) {
	uc, repo, _, _, _, _, _ := newArticleUC()
	repo.GetByIDFn = func(ctx context.Context, id string) (*domain.Article, error) {
		return &domain.Article{ID: id, AuthorID: "u1", Status: domain.StatusPublished}, nil
	}
	repo.UnpublishFn = func(ctx context.Context, id string) error { return nil }
	a, err := uc.UnpublishArticle(context.Background(), "a1", "u1")
	require.NoError(t, err)
	require.Equal(t, domain.StatusDraft, a.Status)
}

func TestArchive_Unarchive(t *testing.T) {
	uc, repo, _, _, _, _, _ := newArticleUC()
	// archive
	repo.GetByIDFn = func(ctx context.Context, id string) (*domain.Article, error) {
		return &domain.Article{ID: id, AuthorID: "u1", Status: domain.StatusDraft}, nil
	}
	archivedAt := time.Time{}
	repo.ArchiveFn = func(ctx context.Context, id string, at time.Time) error { archivedAt = at; return nil }
	a, err := uc.ArchiveArticle(context.Background(), "a1", "u1")
	require.NoError(t, err)
	require.Equal(t, domain.StatusArchived, a.Status)
	require.NotNil(t, a.Timestamps.ArchivedAt)

	// unarchive
	repo.GetByIDFn = func(ctx context.Context, id string) (*domain.Article, error) {
		return &domain.Article{ID: id, AuthorID: "u1", Status: domain.StatusArchived, Timestamps: domain.ArticleTimes{ArchivedAt: &archivedAt}}, nil
	}
	repo.UnarchiveFn = func(ctx context.Context, id string) error { return nil }
	a2, err := uc.UnarchiveArticle(context.Background(), "a1", "u1")
	require.NoError(t, err)
	require.Equal(t, domain.StatusDraft, a2.Status)
}

func TestGetArticleByID_Published_ViewIncrement(t *testing.T) {
	uc, repo, _, _, _, viewUC, _ := newArticleUC()
	repo.GetByIDFn = func(ctx context.Context, id string) (*domain.Article, error) {
		return &domain.Article{ID: id, AuthorID: "author", Status: domain.StatusPublished}, nil
	}
	called := false
	viewUC.RecordViewFn = func(ctx context.Context, uid, aid, ip string) error { called = true; return nil }
	inc := false
	repo.IncrementViewFn = func(ctx context.Context, id string) error { inc = true; return nil }
	_, err := uc.GetArticleByID(context.Background(), "a1", "viewer")
	require.NoError(t, err)
	require.True(t, called)
	require.True(t, inc)
}

func TestGetArticleBySlug_Increments(t *testing.T) {
	uc, repo, _, _, _, viewUC, _ := newArticleUC()
	repo.GetBySlugFn = func(ctx context.Context, slug string) (*domain.Article, error) { return &domain.Article{ID: "a1"}, nil }
	called := false
	viewUC.RecordViewFn = func(ctx context.Context, uid, aid, ip string) error { called = true; return nil }
	inc := false
	repo.IncrementViewFn = func(ctx context.Context, id string) error { inc = true; return nil }
	_, err := uc.GetArticleBySlug(context.Background(), "hello", "1.1.1.1")
	require.NoError(t, err)
	require.True(t, called)
	require.True(t, inc)
}

func TestFilterArticles_DefaultStatus(t *testing.T) {
	uc, repo, _, _, _, _, _ := newArticleUC()
	captured := []domain.ArticleStatus{}
	repo.FilterFn = func(ctx context.Context, filter domain.ArticleFilter, pag domain.Pagination) ([]domain.Article, int, error) {
		captured = filter.Statuses
		return nil, 0, nil
	}
	_, _, err := uc.FilterArticles(context.Background(), domain.ArticleFilter{}, domain.Pagination{})
	require.NoError(t, err)
	require.NotEmpty(t, captured)
	require.Contains(t, captured, domain.StatusPublished)
}

func TestSearchArticles_Unauthorized(t *testing.T) {
	uc, _, policy, _, _, _, _ := newArticleUC()
	policy.UserExistsFn = func(string) bool { return false }
	_, _, err := uc.SearchArticles(context.Background(), "u1", "q", domain.Pagination{})
	require.ErrorIs(t, err, domain.ErrUnauthorized)
}

func TestListArticlesByTags_Unapproved(t *testing.T) {
	uc, _, _, _, tagUC, _, _ := newArticleUC()
	tagUC.IsTagApprovedFn = func(name string) bool { return false }
	_, _, err := uc.ListArticlesByTags(context.Background(), "u", []string{"go"}, domain.Pagination{})
	require.ErrorIs(t, err, domain.ErrUnapprovedTags)
}

func TestAdminListAllArticles_Unauthorized(t *testing.T) {
	uc, _, policy, _, _, _, _ := newArticleUC()
	policy.IsAdminFn = func(string, string) bool { return false }
	_, _, err := uc.AdminListAllArticles(context.Background(), "u1", "user", domain.Pagination{})
	require.ErrorIs(t, err, domain.ErrUnauthorized)
}

func TestAdminHardDelete_Unauthorized(t *testing.T) {
	uc, _, policy, _, _, _, _ := newArticleUC()
	policy.IsAdminFn = func(string, string) bool { return false }
	err := uc.AdminHardDeleteArticle(context.Background(), "u1", "user", "a1")
	require.ErrorIs(t, err, domain.ErrUnauthorized)
}

func TestAdminUnpublish_Success(t *testing.T) {
	uc, repo, policy, _, _, _, _ := newArticleUC()
	policy.IsAdminFn = func(string, string) bool { return true }
	repo.GetByIDFn = func(ctx context.Context, id string) (*domain.Article, error) {
		return &domain.Article{ID: id, Status: domain.StatusPublished}, nil
	}
	repo.UnpublishFn = func(ctx context.Context, id string) error { return nil }
	_, err := uc.AdminUnpublishArticle(context.Background(), "u1", "admin", "a1")
	require.NoError(t, err)
}
