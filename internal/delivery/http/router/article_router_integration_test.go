package router

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

func TestRegisterArticleRouter_WiresRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	// Use a mock usecase so handlers won't panic
	uc := &mocks.ArticleUsecaseMock{GetPopularArticlesFn: func(ctx context.Context, userID string, p domain.Pagination) ([]domain.Article, int, error) {
		return nil, 0, nil
	}}
	h := controller.NewArticleHandler(uc)
	// Add auth context so controller passes auth checks
	r.Use(func(c *gin.Context) { c.Set("user_id", "u1"); c.Set("user_role", string(domain.RoleAdmin)); c.Next() })
	RegisterArticleRouter(r, h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/articles/popular", nil)
	r.ServeHTTP(w, req)
	require.NotEqual(t, http.StatusNotFound, w.Code)
}
