package controller_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"write_base/internal/delivery/http/controller"
	"write_base/internal/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestGetTrendingArticles_Unauthorized_NoUserID(t *testing.T) {
	uc := &mocks.ArticleUsecaseMock{}
	h := controller.NewArticleHandler(uc)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/articles/trending", h.GetTrendingArticles)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/articles/trending?page=1&page_size=10", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetNewArticles_Unauthorized_NoUserID(t *testing.T) {
	uc := &mocks.ArticleUsecaseMock{}
	h := controller.NewArticleHandler(uc)
	r := gin.New()
	r.GET("/articles/new", h.GetNewArticles)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/articles/new?page=1&page_size=10", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetPopularArticles_Unauthorized_NoUserID(t *testing.T) {
	uc := &mocks.ArticleUsecaseMock{}
	h := controller.NewArticleHandler(uc)
	r := gin.New()
	r.GET("/articles/popular", h.GetPopularArticles)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/articles/popular?page=1&page_size=10", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAdminListAllArticles_MissingRole(t *testing.T) {
	uc := &mocks.ArticleUsecaseMock{}
	h := controller.NewArticleHandler(uc)
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("user_id", "u1"); c.Next() })
	r.GET("/admin/articles", h.AdminListAllArticles)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/admin/articles?page=1&page_size=10", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestSearchArticles_Unauthorized_NoUserID(t *testing.T) {
	uc := &mocks.ArticleUsecaseMock{}
	h := controller.NewArticleHandler(uc)
	r := gin.New()
	r.GET("/search", h.SearchArticles)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/search?q=go", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestListArticlesByTags_Unauthorized_NoUserID(t *testing.T) {
	uc := &mocks.ArticleUsecaseMock{}
	h := controller.NewArticleHandler(uc)
	r := gin.New()
	r.GET("/article/tags", h.ListArticlesByTags)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/article/tags?tags=go", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetArticleBySlug_BadSlug_Alt(t *testing.T) {
	uc := &mocks.ArticleUsecaseMock{}
	h := controller.NewArticleHandler(uc)
	r := gin.New()
	r.GET("/:slug", h.GetArticleBySlug)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/bad slug", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddClap_Unauthorized_NoUserID(t *testing.T) {
	uc := &mocks.ArticleUsecaseMock{}
	h := controller.NewArticleHandler(uc)
	r := gin.New()
	r.POST("/articles/:id/clap", h.AddClap)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/articles/a1/clap", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)
}
