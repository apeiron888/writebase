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

// Helpers
func setupAdminRouter(h *controller.Handler, asAdmin bool) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "admin-id")
		role := string(domain.RoleAdmin)
		if !asAdmin {
			role = string(domain.RoleUser)
		}
		c.Set("user_role", role)
		c.Next()
	})
	return r
}

func TestAdminListAllArticles_NotFound(t *testing.T) {
	uc := &mocks.ArticleUsecaseMock{AdminListAllArticlesFn: func(_ context.Context, _ string, _ string, _ domain.Pagination) ([]domain.Article, int, error) {
		return nil, 0, domain.ErrArticleNotFound
	}}
	h := controller.NewArticleHandler(uc)
	r := setupAdminRouter(h, true)
	r.GET("/admin/articles", h.AdminListAllArticles)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/admin/articles?page=1&page_size=10", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestAdminHardDeleteArticle_Success(t *testing.T) {
	uc := &mocks.ArticleUsecaseMock{AdminHardDeleteArticleFn: func(_ context.Context, _ string, _ string, _ string) error { return nil }}
	h := controller.NewArticleHandler(uc)
	r := setupAdminRouter(h, true)
	r.DELETE("/admin/articles/:id/delete", h.AdminHardDeleteArticle)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/admin/articles/a1/delete", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusNoContent, w.Code)
}

func TestAdminHardDeleteArticle_BadRequest(t *testing.T) {
	uc := &mocks.ArticleUsecaseMock{}
	h := controller.NewArticleHandler(uc)
	r := setupAdminRouter(h, true)
	r.DELETE("/admin/articles//delete", h.AdminHardDeleteArticle)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/admin/articles//delete", nil)
	r.ServeHTTP(w, req)
	// malformed path yields 404 (route not matched)
	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestAdminUnpublishArticle_NotFound_Alt(t *testing.T) {
	uc := &mocks.ArticleUsecaseMock{AdminUnpublishArticleFn: func(_ context.Context, _ string, _ string, _ string) (*domain.Article, error) {
		return nil, domain.ErrArticleNotFound
	}}
	h := controller.NewArticleHandler(uc)
	r := setupAdminRouter(h, true)
	r.POST("/admin/articles/:id/unpublish", h.AdminUnpublishArticle)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/admin/articles/a1/unpublish", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestAdminUnpublishArticle_Unauthorized(t *testing.T) {
	uc := &mocks.ArticleUsecaseMock{AdminUnpublishArticleFn: func(_ context.Context, _ string, _ string, _ string) (*domain.Article, error) {
		return nil, domain.ErrUnauthorized
	}}
	h := controller.NewArticleHandler(uc)
	r := setupAdminRouter(h, false) // not admin role
	r.POST("/admin/articles/:id/unpublish", h.AdminUnpublishArticle)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/admin/articles/a1/unpublish", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)
}
