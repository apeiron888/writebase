package mocks

import (
	"context"
	"time"
	"write_base/internal/domain"
)

// ArticleRepositoryMock implements domain.IArticleRepository with pluggable funcs.
type ArticleRepositoryMock struct {
	CreateFn               func(ctx context.Context, article *domain.Article) error
	UpdateFn               func(ctx context.Context, article *domain.Article) error
	DeleteFn               func(ctx context.Context, articleID string) error
	RestoreFn              func(ctx context.Context, articleID string) error
	GetByIDFn              func(ctx context.Context, articleID string) (*domain.Article, error)
	GetBySlugFn            func(ctx context.Context, slug string) (*domain.Article, error)
	GetStatsFn             func(ctx context.Context, articleID string) (*domain.ArticleStats, error)
	GetAllArticleStatsFn   func(ctx context.Context, userID string) ([]domain.ArticleStats, int, error)
	PublishFn              func(ctx context.Context, articleID string, publishAt time.Time) error
	UnpublishFn            func(ctx context.Context, articleID string) error
	ArchiveFn              func(ctx context.Context, articleID string, archiveAt time.Time) error
	UnarchiveFn            func(ctx context.Context, articleID string) error
	ListByAuthorFn         func(ctx context.Context, authorID string, pag domain.Pagination) ([]domain.Article, int, error)
	FindTrendingFn         func(ctx context.Context, windowDays int, pag domain.Pagination) ([]domain.Article, int, error)
	FindNewArticlesFn      func(ctx context.Context, pag domain.Pagination) ([]domain.Article, int, error)
	FindPopularArticlesFn  func(ctx context.Context, pag domain.Pagination) ([]domain.Article, int, error)
	FilterAuthorArticlesFn func(ctx context.Context, authorID string, filter domain.ArticleFilter, pag domain.Pagination) ([]domain.Article, int, error)
	FilterFn               func(ctx context.Context, filter domain.ArticleFilter, pag domain.Pagination) ([]domain.Article, int, error)
	SearchFn               func(ctx context.Context, query string, pag domain.Pagination) ([]domain.Article, int, error)
	ListByTagsFn           func(ctx context.Context, tags []string, pag domain.Pagination) ([]domain.Article, int, error)
	EmptyTrashFn           func(ctx context.Context, userID string) error
	DeleteFromTrashFn      func(ctx context.Context, articleID, userID string) error
	AdminListAllArticlesFn func(ctx context.Context, pag domain.Pagination) ([]domain.Article, int, error)
	HardDeleteFn           func(ctx context.Context, articleID string) error
	IncrementViewFn        func(ctx context.Context, articleID string) error
	UpdateClapCountFn      func(ctx context.Context, articleID string, count int) error
}

func (m *ArticleRepositoryMock) Create(ctx context.Context, a *domain.Article) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, a)
	}
	return nil
}
func (m *ArticleRepositoryMock) Update(ctx context.Context, a *domain.Article) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, a)
	}
	return nil
}
func (m *ArticleRepositoryMock) Delete(ctx context.Context, id string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}
func (m *ArticleRepositoryMock) Restore(ctx context.Context, id string) error {
	if m.RestoreFn != nil {
		return m.RestoreFn(ctx, id)
	}
	return nil
}
func (m *ArticleRepositoryMock) GetByID(ctx context.Context, id string) (*domain.Article, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, domain.ErrArticleNotFound
}
func (m *ArticleRepositoryMock) GetBySlug(ctx context.Context, slug string) (*domain.Article, error) {
	if m.GetBySlugFn != nil {
		return m.GetBySlugFn(ctx, slug)
	}
	return nil, domain.ErrArticleNotFound
}
func (m *ArticleRepositoryMock) GetStats(ctx context.Context, id string) (*domain.ArticleStats, error) {
	if m.GetStatsFn != nil {
		return m.GetStatsFn(ctx, id)
	}
	return nil, domain.ErrArticleNotFound
}
func (m *ArticleRepositoryMock) GetAllArticleStats(ctx context.Context, userID string) ([]domain.ArticleStats, int, error) {
	if m.GetAllArticleStatsFn != nil {
		return m.GetAllArticleStatsFn(ctx, userID)
	}
	return nil, 0, nil
}
func (m *ArticleRepositoryMock) Publish(ctx context.Context, id string, at time.Time) error {
	if m.PublishFn != nil {
		return m.PublishFn(ctx, id, at)
	}
	return nil
}
func (m *ArticleRepositoryMock) Unpublish(ctx context.Context, id string) error {
	if m.UnpublishFn != nil {
		return m.UnpublishFn(ctx, id)
	}
	return nil
}
func (m *ArticleRepositoryMock) Archive(ctx context.Context, id string, at time.Time) error {
	if m.ArchiveFn != nil {
		return m.ArchiveFn(ctx, id, at)
	}
	return nil
}
func (m *ArticleRepositoryMock) Unarchive(ctx context.Context, id string) error {
	if m.UnarchiveFn != nil {
		return m.UnarchiveFn(ctx, id)
	}
	return nil
}
func (m *ArticleRepositoryMock) ListByAuthor(ctx context.Context, authorID string, pag domain.Pagination) ([]domain.Article, int, error) {
	if m.ListByAuthorFn != nil {
		return m.ListByAuthorFn(ctx, authorID, pag)
	}
	return nil, 0, nil
}
func (m *ArticleRepositoryMock) FindTrending(ctx context.Context, windowDays int, pag domain.Pagination) ([]domain.Article, int, error) {
	if m.FindTrendingFn != nil {
		return m.FindTrendingFn(ctx, windowDays, pag)
	}
	return nil, 0, nil
}
func (m *ArticleRepositoryMock) FindNewArticles(ctx context.Context, pag domain.Pagination) ([]domain.Article, int, error) {
	if m.FindNewArticlesFn != nil {
		return m.FindNewArticlesFn(ctx, pag)
	}
	return nil, 0, nil
}
func (m *ArticleRepositoryMock) FindPopularArticles(ctx context.Context, pag domain.Pagination) ([]domain.Article, int, error) {
	if m.FindPopularArticlesFn != nil {
		return m.FindPopularArticlesFn(ctx, pag)
	}
	return nil, 0, nil
}
func (m *ArticleRepositoryMock) FilterAuthorArticles(ctx context.Context, authorID string, filter domain.ArticleFilter, pag domain.Pagination) ([]domain.Article, int, error) {
	if m.FilterAuthorArticlesFn != nil {
		return m.FilterAuthorArticlesFn(ctx, authorID, filter, pag)
	}
	return nil, 0, nil
}
func (m *ArticleRepositoryMock) Filter(ctx context.Context, filter domain.ArticleFilter, pag domain.Pagination) ([]domain.Article, int, error) {
	if m.FilterFn != nil {
		return m.FilterFn(ctx, filter, pag)
	}
	return nil, 0, nil
}
func (m *ArticleRepositoryMock) Search(ctx context.Context, query string, pag domain.Pagination) ([]domain.Article, int, error) {
	if m.SearchFn != nil {
		return m.SearchFn(ctx, query, pag)
	}
	return nil, 0, nil
}
func (m *ArticleRepositoryMock) ListByTags(ctx context.Context, tags []string, pag domain.Pagination) ([]domain.Article, int, error) {
	if m.ListByTagsFn != nil {
		return m.ListByTagsFn(ctx, tags, pag)
	}
	return nil, 0, nil
}
func (m *ArticleRepositoryMock) EmptyTrash(ctx context.Context, userID string) error {
	if m.EmptyTrashFn != nil {
		return m.EmptyTrashFn(ctx, userID)
	}
	return nil
}
func (m *ArticleRepositoryMock) DeleteFromTrash(ctx context.Context, articleID, userID string) error {
	if m.DeleteFromTrashFn != nil {
		return m.DeleteFromTrashFn(ctx, articleID, userID)
	}
	return nil
}
func (m *ArticleRepositoryMock) AdminListAllArticles(ctx context.Context, pag domain.Pagination) ([]domain.Article, int, error) {
	if m.AdminListAllArticlesFn != nil {
		return m.AdminListAllArticlesFn(ctx, pag)
	}
	return nil, 0, nil
}
func (m *ArticleRepositoryMock) HardDelete(ctx context.Context, articleID string) error {
	if m.HardDeleteFn != nil {
		return m.HardDeleteFn(ctx, articleID)
	}
	return nil
}
func (m *ArticleRepositoryMock) IncrementView(ctx context.Context, articleID string) error {
	if m.IncrementViewFn != nil {
		return m.IncrementViewFn(ctx, articleID)
	}
	return nil
}
func (m *ArticleRepositoryMock) UpdateClapCount(ctx context.Context, articleID string, count int) error {
	if m.UpdateClapCountFn != nil {
		return m.UpdateClapCountFn(ctx, articleID, count)
	}
	return nil
}
