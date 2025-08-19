package usecase_test

import (
	"context"
	"testing"
	"time"
	"write_base/internal/domain"

	"github.com/stretchr/testify/require"
)

func TestPublishArticle_UnapprovedTags(t *testing.T) {
	uc, repo, _, _, tagUC, _, _ := newArticleUC()
	repo.GetByIDFn = func(ctx context.Context, id string) (*domain.Article, error) {
		return &domain.Article{ID: id, AuthorID: "u1", Status: domain.StatusDraft, Tags: []string{"x"}}, nil
	}
	tagUC.IsTagApprovedFn = func(tag string) bool { return false }
	_, err := uc.PublishArticle(context.Background(), "a1", "u1")
	require.ErrorIs(t, err, domain.ErrUnapprovedTags)
}

func TestPublishArticle_Success(t *testing.T) {
	uc, repo, _, _, tagUC, _, _ := newArticleUC()
	repo.GetByIDFn = func(ctx context.Context, id string) (*domain.Article, error) {
		return &domain.Article{ID: id, AuthorID: "u1", Status: domain.StatusDraft, Tags: []string{"go"}}, nil
	}
	tagUC.IsTagApprovedFn = func(tag string) bool { return true }
	repo.PublishFn = func(ctx context.Context, id string, at time.Time) error { return nil }
	a, err := uc.PublishArticle(context.Background(), "a1", "u1")
	require.NoError(t, err)
	require.Equal(t, domain.StatusPublished, a.Status)
	require.NotNil(t, a.Timestamps.PublishedAt)
}

func TestTrending_Unauthorized(t *testing.T) {
	uc, _, policy, _, _, _, _ := newArticleUC()
	policy.UserExistsFn = func(string) bool { return false }
	_, _, err := uc.GetTrendingArticles(context.Background(), "u1", domain.Pagination{})
	require.ErrorIs(t, err, domain.ErrUnauthorized)
}

func TestFilterAuthorArticles_Unauthorized(t *testing.T) {
	uc, _, policy, _, _, _, _ := newArticleUC()
	policy.UserExistsFn = func(uid string) bool { return true }
	_, _, err := uc.FilterAuthorArticles(context.Background(), "caller", "author", domain.ArticleFilter{}, domain.Pagination{})
	require.ErrorIs(t, err, domain.ErrUnauthorized)
}

func TestFilterAuthorArticles_Success(t *testing.T) {
	uc, repo, policy, _, _, _, _ := newArticleUC()
	policy.UserExistsFn = func(uid string) bool { return true }
	repo.FilterAuthorArticlesFn = func(ctx context.Context, authorID string, f domain.ArticleFilter, p domain.Pagination) ([]domain.Article, int, error) {
		return nil, 0, nil
	}
	_, _, err := uc.FilterAuthorArticles(context.Background(), "u1", "u1", domain.ArticleFilter{}, domain.Pagination{})
	require.NoError(t, err)
}

func TestEmptyTrash_Success(t *testing.T) {
	uc, repo, _, _, _, _, _ := newArticleUC()
	repo.EmptyTrashFn = func(ctx context.Context, uid string) error { return nil }
	err := uc.EmptyTrash(context.Background(), "u1")
	require.NoError(t, err)
}

func TestDeleteFromTrash_NotDeleted(t *testing.T) {
	uc, repo, _, _, _, _, _ := newArticleUC()
	repo.GetByIDFn = func(ctx context.Context, id string) (*domain.Article, error) {
		return &domain.Article{ID: id, Status: domain.StatusDraft, AuthorID: "u1"}, nil
	}
	err := uc.DeleteArticleFromTrash(context.Background(), "a1", "u1")
	require.ErrorIs(t, err, domain.ErrArticleNotFound)
}

func TestAddClap_Success(t *testing.T) {
	uc, repo, policy, _, _, _, clap := newArticleUC()
	policy.UserExistsFn = func(string) bool { return true }
	clap.AddClapFn = func(ctx context.Context, uid, aid string) (int, error) { return 5, nil }
	repo.UpdateClapCountFn = func(ctx context.Context, id string, count int) error { return nil }
	repo.GetByIDFn = func(ctx context.Context, id string) (*domain.Article, error) {
		return &domain.Article{ID: id, Stats: domain.ArticleStats{ClapCount: 5}}, nil
	}
	stats, err := uc.AddClap(context.Background(), "u1", "a1")
	require.NoError(t, err)
	require.Equal(t, 5, stats.ClapCount)
}
