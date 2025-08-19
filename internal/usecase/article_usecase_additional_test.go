package usecase_test

import (
	"context"
	"testing"
	"write_base/internal/domain"

	"github.com/stretchr/testify/require"
)

func TestListArticlesByAuthor_Unauthorized(t *testing.T) {
	uc, _, policy, _, _, _, _ := newArticleUC()
	policy.UserExistsFn = func(string) bool { return false }
	_, _, err := uc.ListArticlesByAuthor(context.Background(), "caller", "author", domain.Pagination{})
	require.ErrorIs(t, err, domain.ErrUnauthorized)
}

func TestGetArticleStats_Success(t *testing.T) {
	uc, repo, _, _, _, _, _ := newArticleUC()
	repo.GetStatsFn = func(ctx context.Context, id string) (*domain.ArticleStats, error) {
		return &domain.ArticleStats{ViewCount: 1}, nil
	}
	stats, err := uc.GetArticleStats(context.Background(), "a1", "u1")
	require.NoError(t, err)
	require.Equal(t, 1, stats.ViewCount)
}

func TestGetAllArticleStats_Success(t *testing.T) {
	uc, repo, _, _, _, _, _ := newArticleUC()
	repo.GetAllArticleStatsFn = func(ctx context.Context, uid string) ([]domain.ArticleStats, int, error) {
		return []domain.ArticleStats{{ViewCount: 2}}, 1, nil
	}
	stats, total, err := uc.GetAllArticleStats(context.Background(), "u1")
	require.NoError(t, err)
	require.Equal(t, 1, total)
	require.Equal(t, 2, stats[0].ViewCount)
}
