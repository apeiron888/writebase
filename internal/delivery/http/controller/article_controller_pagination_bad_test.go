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

func TestGetTrendingArticles_BadPagination(t *testing.T) {
	h := controller.NewArticleHandler(&mocks.ArticleUsecaseMock{})
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("user_id", "u1"); c.Next() })
	r.GET("/articles/trending", h.GetTrendingArticles)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/articles/trending?page=abc&page_size=10", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSearchArticles_BadPagination(t *testing.T) {
	h := controller.NewArticleHandler(&mocks.ArticleUsecaseMock{})
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("user_id", "u1"); c.Next() })
	r.GET("/search", h.SearchArticles)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/search?q=go&page_size=abc", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddClap_TooMany(t *testing.T) {
	uc := &mocks.ArticleUsecaseMock{AddClapFn: func(_ context.Context, _ string, _ string) (domain.ArticleStats, error) {
		return domain.ArticleStats{}, domain.ErrClapLimitExceeded
	}}
	h := controller.NewArticleHandler(uc)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("user_id", "u1"); c.Next() })
	r.POST("/articles/:id/clap", h.AddClap)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/articles/a1/clap", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusTooManyRequests, w.Code)
}
