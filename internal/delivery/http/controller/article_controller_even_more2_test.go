package controller_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"write_base/internal/delivery/http/controller"
	"write_base/internal/domain"
	"write_base/internal/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestListArticlesByAuthor_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uc := &mocks.ArticleUsecaseMock{ListArticlesByAuthorFn: func(ctx context.Context, uid, aid string, p domain.Pagination) ([]domain.Article, int, error) {
		return []domain.Article{}, 0, nil
	}}
	h := controller.NewArticleHandler(uc)
	r := gin.New()
	r.Use(withAuth())
	r.GET("/authors/:author_id/articles", h.ListArticlesByAuthor)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/authors/u1/articles?page=1&page_size=10", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestGetTrending_New_Popular_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	list := []domain.Article{{ID: "a1"}}
	uc := &mocks.ArticleUsecaseMock{
		GetTrendingArticlesFn: func(ctx context.Context, uid string, p domain.Pagination) ([]domain.Article, int, error) {
			return list, 1, nil
		},
		GetNewArticlesFn: func(ctx context.Context, uid string, p domain.Pagination) ([]domain.Article, int, error) {
			return list, 1, nil
		},
		GetPopularArticlesFn: func(ctx context.Context, uid string, p domain.Pagination) ([]domain.Article, int, error) {
			return list, 1, nil
		},
	}
	h := controller.NewArticleHandler(uc)
	r := gin.New()
	r.Use(withAuth())
	r.GET("/articles/trending", h.GetTrendingArticles)
	r.GET("/articles/new", h.GetNewArticles)
	r.GET("/articles/popular", h.GetPopularArticles)
	for _, path := range []string{"/articles/trending", "/articles/new", "/articles/popular"} {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, path+"?page=1&page_size=10", nil)
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
	}
}

func TestSearchArticles_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uc := &mocks.ArticleUsecaseMock{SearchArticlesFn: func(ctx context.Context, uid, q string, p domain.Pagination) ([]domain.Article, int, error) {
		return []domain.Article{{ID: "a1"}}, 1, nil
	}}
	h := controller.NewArticleHandler(uc)
	r := gin.New()
	r.Use(withAuth())
	r.GET("/search", h.SearchArticles)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/search?q=go&page=1&page_size=5", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestListArticlesByTags_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uc := &mocks.ArticleUsecaseMock{ListArticlesByTagsFn: func(ctx context.Context, uid string, tags []string, p domain.Pagination) ([]domain.Article, int, error) {
		return []domain.Article{}, 0, nil
	}}
	h := controller.NewArticleHandler(uc)
	r := gin.New()
	r.Use(withAuth())
	r.GET("/article/tags", h.ListArticlesByTags)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/article/tags?tags=go&page=1&page_size=5", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestEmptyTrash_And_DeleteFromTrash_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uc := &mocks.ArticleUsecaseMock{EmptyTrashFn: func(ctx context.Context, uid string) error { return nil }, DeleteArticleFromTrashFn: func(ctx context.Context, id, uid string) error { return nil }}
	h := controller.NewArticleHandler(uc)
	r := gin.New()
	r.Use(withAuth())
	r.DELETE("/me/trash", h.EmptyTrash)
	r.DELETE("/articles/trash/:id", h.DeleteFromTrash)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/me/trash", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusNoContent, w.Code)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodDelete, "/articles/trash/a1", nil)
	r.ServeHTTP(w2, req2)
	require.Equal(t, http.StatusNoContent, w2.Code)
}

func TestGetArticleBySlug_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uc := &mocks.ArticleUsecaseMock{GetArticleBySlugFn: func(ctx context.Context, slug, ip string) (*domain.Article, error) {
		return &domain.Article{ID: "a1", Slug: slug}, nil
	}}
	h := controller.NewArticleHandler(uc)
	r := gin.New()
	r.GET("/:slug", h.GetArticleBySlug)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/hello-world", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
