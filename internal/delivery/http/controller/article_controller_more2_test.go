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

func setupRouterWithAuth(h *controller.Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("user_id", "u1"); c.Set("user_role", string(domain.RoleAdmin)); c.Next() })
	return r
}

func TestPublishArticle_Success(t *testing.T) {
	uc := &mocks.ArticleUsecaseMock{PublishArticleFn: func(ctx context.Context, id, uid string) (*domain.Article, error) {
		return &domain.Article{ID: id, Status: domain.StatusPublished}, nil
	}}
	h := controller.NewArticleHandler(uc)
	r := setupRouterWithAuth(h)
	r.POST("/articles/:id/publish", h.PublishArticle)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/articles/a1/publish", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestPublishArticle_Conflict(t *testing.T) {
	uc := &mocks.ArticleUsecaseMock{PublishArticleFn: func(ctx context.Context, id, uid string) (*domain.Article, error) {
		return nil, domain.ErrArticlePublished
	}}
	h := controller.NewArticleHandler(uc)
	r := setupRouterWithAuth(h)
	r.POST("/articles/:id/publish", h.PublishArticle)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/articles/a1/publish", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusConflict, w.Code)
}

func TestUnpublishArticle_Conflict(t *testing.T) {
	uc := &mocks.ArticleUsecaseMock{UnpublishArticleFn: func(ctx context.Context, id, uid string) (*domain.Article, error) {
		return nil, domain.ErrArticleNotPublished
	}}
	h := controller.NewArticleHandler(uc)
	r := setupRouterWithAuth(h)
	r.POST("/articles/:id/unpublish", h.UnpublishArticle)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/articles/a1/unpublish", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusConflict, w.Code)
}

func TestGetArticleByID_NotFound(t *testing.T) {
	uc := &mocks.ArticleUsecaseMock{GetArticleByIDFn: func(ctx context.Context, id, uid string) (*domain.Article, error) {
		return nil, domain.ErrArticleNotFound
	}}
	h := controller.NewArticleHandler(uc)
	r := setupRouterWithAuth(h)
	r.GET("/articles/:id", h.GetArticleByID)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/articles/missing", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestSearchArticles_BadQuery(t *testing.T) {
	uc := &mocks.ArticleUsecaseMock{}
	h := controller.NewArticleHandler(uc)
	r := setupRouterWithAuth(h)
	r.GET("/search", h.SearchArticles)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/search?q=", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestListArticlesByTags_Unapproved(t *testing.T) {
	uc := &mocks.ArticleUsecaseMock{ListArticlesByTagsFn: func(ctx context.Context, uid string, tags []string, p domain.Pagination) ([]domain.Article, int, error) {
		return nil, 0, domain.ErrUnapprovedTags
	}}
	h := controller.NewArticleHandler(uc)
	r := setupRouterWithAuth(h)
	r.GET("/article/tags", h.ListArticlesByTags)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/article/tags?tags=go", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddClap_Success(t *testing.T) {
	uc := &mocks.ArticleUsecaseMock{AddClapFn: func(ctx context.Context, uid, aid string) (domain.ArticleStats, error) {
		return domain.ArticleStats{ClapCount: 3}, nil
	}}
	h := controller.NewArticleHandler(uc)
	r := setupRouterWithAuth(h)
	r.POST("/articles/:id/clap", h.AddClap)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/articles/a1/clap", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
