package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"write_base/internal/delivery/http/controller"
	"write_base/internal/domain"
	"write_base/internal/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func withAuth() gin.HandlerFunc {
	return func(c *gin.Context) { c.Set("user_id", "u1"); c.Set("user_role", string(domain.RoleAdmin)); c.Next() }
}

func TestUpdateArticle_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uc := &mocks.ArticleUsecaseMock{UpdateArticleFn: func(ctx context.Context, uid string, a *domain.Article) error { return domain.ErrUnauthorized }}
	h := controller.NewArticleHandler(uc)
	r := gin.New()
	r.Use(withAuth())
	r.PUT("/articles/:id", h.UpdateArticle)
	body := map[string]any{"title": "x", "content_blocks": []map[string]any{{"type": "paragraph", "order": 0, "content": map[string]any{"paragraph": map[string]any{"text": "hi"}}}}}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/articles/a1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUpdateArticle_BadPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uc := &mocks.ArticleUsecaseMock{}
	h := controller.NewArticleHandler(uc)
	r := gin.New()
	r.Use(withAuth())
	r.PUT("/articles/:id", h.UpdateArticle)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/articles/a1", bytes.NewBufferString("not-json"))
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteArticle_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uc := &mocks.ArticleUsecaseMock{DeleteArticleFn: func(ctx context.Context, id, uid string) error { return domain.ErrArticleNotFound }}
	h := controller.NewArticleHandler(uc)
	r := gin.New()
	r.Use(withAuth())
	r.DELETE("/articles/:id", h.DeleteArticle)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/articles/a1", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteArticle_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uc := &mocks.ArticleUsecaseMock{DeleteArticleFn: func(ctx context.Context, id, uid string) error { return nil }}
	h := controller.NewArticleHandler(uc)
	r := gin.New()
	r.Use(withAuth())
	r.DELETE("/articles/:id", h.DeleteArticle)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/articles/a1", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestRestoreArticle_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uc := &mocks.ArticleUsecaseMock{RestoreArticleFn: func(ctx context.Context, uid string, id string) error { return nil }}
	h := controller.NewArticleHandler(uc)
	r := gin.New()
	r.Use(withAuth())
	r.PATCH("/articles/:id/restore", h.RestoreArticle)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/articles/a1/restore", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestArchiveArticle_Conflict(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uc := &mocks.ArticleUsecaseMock{ArchiveArticleFn: func(ctx context.Context, id, uid string) (*domain.Article, error) {
		return nil, domain.ErrArticleArchived
	}}
	h := controller.NewArticleHandler(uc)
	r := gin.New()
	r.Use(withAuth())
	r.POST("/articles/:id/archive", h.ArchiveArticle)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/articles/a1/archive", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusConflict, w.Code)
}

func TestUnarchiveArticle_Conflict(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uc := &mocks.ArticleUsecaseMock{UnarchiveArticleFn: func(ctx context.Context, id, uid string) (*domain.Article, error) {
		return nil, domain.ErrArticleNotArchived
	}}
	h := controller.NewArticleHandler(uc)
	r := gin.New()
	r.Use(withAuth())
	r.POST("/articles/:id/unarchive", h.UnarchiveArticle)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/articles/a1/unarchive", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusConflict, w.Code)
}

func TestGetArticleStats_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uc := &mocks.ArticleUsecaseMock{GetArticleStatsFn: func(ctx context.Context, id, uid string) (*domain.ArticleStats, error) {
		return nil, domain.ErrArticleNotFound
	}}
	h := controller.NewArticleHandler(uc)
	r := gin.New()
	r.Use(withAuth())
	r.GET("/articles/:id/stats", h.GetArticleStats)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/articles/a1/stats", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetAllArticleStats_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uc := &mocks.ArticleUsecaseMock{GetAllArticleStatsFn: func(ctx context.Context, uid string) ([]domain.ArticleStats, int, error) {
		return []domain.ArticleStats{{ClapCount: 1}}, 1, nil
	}}
	h := controller.NewArticleHandler(uc)
	r := gin.New()
	r.Use(withAuth())
	r.GET("/articles/stats/all", h.GetAllArticleStats)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/articles/stats/all", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestFilterArticles_BadPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uc := &mocks.ArticleUsecaseMock{}
	h := controller.NewArticleHandler(uc)
	r := gin.New()
	r.Use(withAuth())
	r.POST("/articles/filter", h.FilterArticles)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/articles/filter", bytes.NewBufferString("{"))
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFilterArticles_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uc := &mocks.ArticleUsecaseMock{FilterArticlesFn: func(ctx context.Context, f domain.ArticleFilter, p domain.Pagination) ([]domain.Article, int, error) {
		return []domain.Article{}, 0, nil
	}}
	h := controller.NewArticleHandler(uc)
	r := gin.New()
	r.Use(withAuth())
	r.POST("/articles/filter", h.FilterArticles)
	payload := map[string]any{"filter": map[string]any{}, "pagination": map[string]any{"page": 1, "page_size": 10}}
	b, _ := json.Marshal(payload)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/articles/filter", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestFilterAuthorArticles_BadPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uc := &mocks.ArticleUsecaseMock{}
	h := controller.NewArticleHandler(uc)
	r := gin.New()
	r.Use(withAuth())
	r.POST("/authors/:author_id/articles/filter", h.FilterAuthorArticles)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/authors/u1/articles/filter", bytes.NewBufferString("{"))
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAdminListAllArticles_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uc := &mocks.ArticleUsecaseMock{AdminListAllArticlesFn: func(ctx context.Context, uid, role string, p domain.Pagination) ([]domain.Article, int, error) {
		return nil, 0, domain.ErrUnauthorized
	}}
	h := controller.NewArticleHandler(uc)
	r := gin.New()
	r.Use(withAuth())
	r.GET("/admin/articles", h.AdminListAllArticles)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/admin/articles?page=1&page_size=10", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAdminHardDeleteArticle_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uc := &mocks.ArticleUsecaseMock{AdminHardDeleteArticleFn: func(ctx context.Context, uid, role, id string) error { return domain.ErrArticleNotFound }}
	h := controller.NewArticleHandler(uc)
	r := gin.New()
	r.Use(withAuth())
	r.DELETE("/admin/articles/:id/delete", h.AdminHardDeleteArticle)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/admin/articles/a1/delete", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestAdminUnpublishArticle_Conflict(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uc := &mocks.ArticleUsecaseMock{AdminUnpublishArticleFn: func(ctx context.Context, uid, role, id string) (*domain.Article, error) {
		return nil, domain.ErrArticleNotPublished
	}}
	h := controller.NewArticleHandler(uc)
	r := gin.New()
	r.Use(withAuth())
	r.POST("/admin/articles/:id/unpublish", h.AdminUnpublishArticle)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/admin/articles/a1/unpublish", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusConflict, w.Code)
}
