package mocks

import (
	"context"
	"write_base/internal/domain"
)

// ArticleUsecaseMock implements domain.IArticleUsecase with pluggable funcs for tests.
type ArticleUsecaseMock struct {
	CreateArticleFn             func(ctx context.Context, userID string, input *domain.Article) (string, error)
	UpdateArticleFn             func(ctx context.Context, userID string, input *domain.Article) error
	DeleteArticleFn             func(ctx context.Context, articleID, userID string) error
	RestoreArticleFn            func(ctx context.Context, userID string, articleID string) error
	GetArticleByIDFn            func(ctx context.Context, articleID, userID string) (*domain.Article, error)
	GetArticleBySlugFn          func(ctx context.Context, slug string, clientIP string) (*domain.Article, error)
	GetArticleStatsFn           func(ctx context.Context, articleID, userID string) (*domain.ArticleStats, error)
	GetAllArticleStatsFn        func(ctx context.Context, userID string) ([]domain.ArticleStats, int, error)
	PublishArticleFn            func(ctx context.Context, articleID, userID string) (*domain.Article, error)
	UnpublishArticleFn          func(ctx context.Context, articleID, userID string) (*domain.Article, error)
	ArchiveArticleFn            func(ctx context.Context, articleID, userID string) (*domain.Article, error)
	UnarchiveArticleFn          func(ctx context.Context, articleID, userID string) (*domain.Article, error)
	ListArticlesByAuthorFn      func(ctx context.Context, userID, authorID string, pag domain.Pagination) ([]domain.Article, int, error)
	GetTrendingArticlesFn       func(ctx context.Context, userID string, pag domain.Pagination) ([]domain.Article, int, error)
	GetNewArticlesFn            func(ctx context.Context, userID string, pag domain.Pagination) ([]domain.Article, int, error)
	GetPopularArticlesFn        func(ctx context.Context, userID string, pag domain.Pagination) ([]domain.Article, int, error)
	FilterAuthorArticlesFn      func(ctx context.Context, callerID, authorID string, filter domain.ArticleFilter, pag domain.Pagination) ([]domain.Article, int, error)
	FilterArticlesFn            func(ctx context.Context, filter domain.ArticleFilter, pag domain.Pagination) ([]domain.Article, int, error)
	SearchArticlesFn            func(ctx context.Context, userID, query string, pag domain.Pagination) ([]domain.Article, int, error)
	ListArticlesByTagsFn        func(ctx context.Context, userID string, tags []string, pag domain.Pagination) ([]domain.Article, int, error)
	EmptyTrashFn                func(ctx context.Context, userID string) error
	DeleteArticleFromTrashFn    func(ctx context.Context, articleID, userID string) error
	AdminListAllArticlesFn      func(ctx context.Context, userID, userRole string, pag domain.Pagination) ([]domain.Article, int, error)
	AdminHardDeleteArticleFn    func(ctx context.Context, userID, userRole, articleID string) error
	AdminUnpublishArticleFn     func(ctx context.Context, userID, userRole, articleID string) (*domain.Article, error)
	AddClapFn                   func(ctx context.Context, userID, articleID string) (domain.ArticleStats, error)
	GenerateContentForArticleFn func(ctx context.Context, article *domain.Article, instructions string) (*domain.Article, error)
	GenerateSlugForTitleFn      func(ctx context.Context, title string) (string, error)
}

func (m *ArticleUsecaseMock) CreateArticle(ctx context.Context, userID string, input *domain.Article) (string, error) {
	if m.CreateArticleFn != nil {
		return m.CreateArticleFn(ctx, userID, input)
	}
	return "", nil
}
func (m *ArticleUsecaseMock) UpdateArticle(ctx context.Context, userID string, input *domain.Article) error {
	if m.UpdateArticleFn != nil {
		return m.UpdateArticleFn(ctx, userID, input)
	}
	return nil
}
func (m *ArticleUsecaseMock) DeleteArticle(ctx context.Context, articleID, userID string) error {
	if m.DeleteArticleFn != nil {
		return m.DeleteArticleFn(ctx, articleID, userID)
	}
	return nil
}
func (m *ArticleUsecaseMock) RestoreArticle(ctx context.Context, userID string, articleID string) error {
	if m.RestoreArticleFn != nil {
		return m.RestoreArticleFn(ctx, userID, articleID)
	}
	return nil
}
func (m *ArticleUsecaseMock) GetArticleByID(ctx context.Context, articleID, userID string) (*domain.Article, error) {
	if m.GetArticleByIDFn != nil {
		return m.GetArticleByIDFn(ctx, articleID, userID)
	}
	return nil, nil
}
func (m *ArticleUsecaseMock) GetArticleBySlug(ctx context.Context, slug string, clientIP string) (*domain.Article, error) {
	if m.GetArticleBySlugFn != nil {
		return m.GetArticleBySlugFn(ctx, slug, clientIP)
	}
	return nil, nil
}
func (m *ArticleUsecaseMock) GetArticleStats(ctx context.Context, articleID, userID string) (*domain.ArticleStats, error) {
	if m.GetArticleStatsFn != nil {
		return m.GetArticleStatsFn(ctx, articleID, userID)
	}
	return nil, nil
}
func (m *ArticleUsecaseMock) GetAllArticleStats(ctx context.Context, userID string) ([]domain.ArticleStats, int, error) {
	if m.GetAllArticleStatsFn != nil {
		return m.GetAllArticleStatsFn(ctx, userID)
	}
	return nil, 0, nil
}
func (m *ArticleUsecaseMock) PublishArticle(ctx context.Context, articleID, userID string) (*domain.Article, error) {
	if m.PublishArticleFn != nil {
		return m.PublishArticleFn(ctx, articleID, userID)
	}
	return nil, nil
}
func (m *ArticleUsecaseMock) UnpublishArticle(ctx context.Context, articleID, userID string) (*domain.Article, error) {
	if m.UnpublishArticleFn != nil {
		return m.UnpublishArticleFn(ctx, articleID, userID)
	}
	return nil, nil
}
func (m *ArticleUsecaseMock) ArchiveArticle(ctx context.Context, articleID, userID string) (*domain.Article, error) {
	if m.ArchiveArticleFn != nil {
		return m.ArchiveArticleFn(ctx, articleID, userID)
	}
	return nil, nil
}
func (m *ArticleUsecaseMock) UnarchiveArticle(ctx context.Context, articleID, userID string) (*domain.Article, error) {
	if m.UnarchiveArticleFn != nil {
		return m.UnarchiveArticleFn(ctx, articleID, userID)
	}
	return nil, nil
}
func (m *ArticleUsecaseMock) ListArticlesByAuthor(ctx context.Context, userID, authorID string, pag domain.Pagination) ([]domain.Article, int, error) {
	if m.ListArticlesByAuthorFn != nil {
		return m.ListArticlesByAuthorFn(ctx, userID, authorID, pag)
	}
	return nil, 0, nil
}
func (m *ArticleUsecaseMock) GetTrendingArticles(ctx context.Context, userID string, pag domain.Pagination) ([]domain.Article, int, error) {
	if m.GetTrendingArticlesFn != nil {
		return m.GetTrendingArticlesFn(ctx, userID, pag)
	}
	return nil, 0, nil
}
func (m *ArticleUsecaseMock) GetNewArticles(ctx context.Context, userID string, pag domain.Pagination) ([]domain.Article, int, error) {
	if m.GetNewArticlesFn != nil {
		return m.GetNewArticlesFn(ctx, userID, pag)
	}
	return nil, 0, nil
}
func (m *ArticleUsecaseMock) GetPopularArticles(ctx context.Context, userID string, pag domain.Pagination) ([]domain.Article, int, error) {
	if m.GetPopularArticlesFn != nil {
		return m.GetPopularArticlesFn(ctx, userID, pag)
	}
	return nil, 0, nil
}
func (m *ArticleUsecaseMock) FilterAuthorArticles(ctx context.Context, callerID, authorID string, filter domain.ArticleFilter, pag domain.Pagination) ([]domain.Article, int, error) {
	if m.FilterAuthorArticlesFn != nil {
		return m.FilterAuthorArticlesFn(ctx, callerID, authorID, filter, pag)
	}
	return nil, 0, nil
}
func (m *ArticleUsecaseMock) FilterArticles(ctx context.Context, filter domain.ArticleFilter, pag domain.Pagination) ([]domain.Article, int, error) {
	if m.FilterArticlesFn != nil {
		return m.FilterArticlesFn(ctx, filter, pag)
	}
	return nil, 0, nil
}
func (m *ArticleUsecaseMock) SearchArticles(ctx context.Context, userID, query string, pag domain.Pagination) ([]domain.Article, int, error) {
	if m.SearchArticlesFn != nil {
		return m.SearchArticlesFn(ctx, userID, query, pag)
	}
	return nil, 0, nil
}
func (m *ArticleUsecaseMock) ListArticlesByTags(ctx context.Context, userID string, tags []string, pag domain.Pagination) ([]domain.Article, int, error) {
	if m.ListArticlesByTagsFn != nil {
		return m.ListArticlesByTagsFn(ctx, userID, tags, pag)
	}
	return nil, 0, nil
}
func (m *ArticleUsecaseMock) EmptyTrash(ctx context.Context, userID string) error {
	if m.EmptyTrashFn != nil {
		return m.EmptyTrashFn(ctx, userID)
	}
	return nil
}
func (m *ArticleUsecaseMock) DeleteArticleFromTrash(ctx context.Context, articleID, userID string) error {
	if m.DeleteArticleFromTrashFn != nil {
		return m.DeleteArticleFromTrashFn(ctx, articleID, userID)
	}
	return nil
}
func (m *ArticleUsecaseMock) AdminListAllArticles(ctx context.Context, userID, userRole string, pag domain.Pagination) ([]domain.Article, int, error) {
	if m.AdminListAllArticlesFn != nil {
		return m.AdminListAllArticlesFn(ctx, userID, userRole, pag)
	}
	return nil, 0, nil
}
func (m *ArticleUsecaseMock) AdminHardDeleteArticle(ctx context.Context, userID, userRole, articleID string) error {
	if m.AdminHardDeleteArticleFn != nil {
		return m.AdminHardDeleteArticleFn(ctx, userID, userRole, articleID)
	}
	return nil
}
func (m *ArticleUsecaseMock) AdminUnpublishArticle(ctx context.Context, userID, userRole, articleID string) (*domain.Article, error) {
	if m.AdminUnpublishArticleFn != nil {
		return m.AdminUnpublishArticleFn(ctx, userID, userRole, articleID)
	}
	return nil, nil
}
func (m *ArticleUsecaseMock) AddClap(ctx context.Context, userID, articleID string) (domain.ArticleStats, error) {
	if m.AddClapFn != nil {
		return m.AddClapFn(ctx, userID, articleID)
	}
	return domain.ArticleStats{}, nil
}
func (m *ArticleUsecaseMock) GenerateContentForArticle(ctx context.Context, article *domain.Article, instructions string) (*domain.Article, error) {
	if m.GenerateContentForArticleFn != nil {
		return m.GenerateContentForArticleFn(ctx, article, instructions)
	}
	return nil, nil
}
func (m *ArticleUsecaseMock) GenerateSlugForTitle(ctx context.Context, title string) (string, error) {
	if m.GenerateSlugForTitleFn != nil {
		return m.GenerateSlugForTitleFn(ctx, title)
	}
	return "", nil
}
