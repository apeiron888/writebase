package controller_test

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/suite"
// 	"write_base/internal/delivery/http/controller"
// 	"write_base/internal/mocks"
// 	"write_base/internal/domain"
// )

// type ArticleControllerSuite struct {
// 	suite.Suite
// 	router      *gin.Engine
// 	mockUsecase *mocks.MockIArticleUsecase
// }

// func (s *ArticleControllerSuite) SetupTest() {
// 	gin.SetMode(gin.TestMode)
// 	s.mockUsecase = new(mocks.MockIArticleUsecase)
// 	h := controller.NewArticleHandler(s.mockUsecase)
// 	r := gin.New()
// 	r.Use(func(c *gin.Context) {
// 		c.Set("user_id", "test_user")
// 		c.Set("user_role", "admin")
// 		c.Next()
// 	})
// 	RegisterTestArticleRoutes(r, h)
// 	s.router = r
// }

// func RegisterTestArticleRoutes(r *gin.Engine, h *controller.ArticleHandler) {
// 	r.GET("/p/:slug", h.ViewArticleBySlug)
// 	r.POST("/articles", h.CreateArticle)
// 	r.PUT("/articles/:id", h.UpdateArticle)
// 	r.DELETE("/articles/:id", h.DeleteArticle)
// 	r.POST("/articles/:id/restore", h.RestoreArticle)
// 	r.POST("/articles/:id/publish", h.PublishArticle)
// 	r.POST("/articles/:id/unpublish", h.UnpublishArticle)
// 	r.POST("/articles/:id/archive", h.ArchiveArticle)
// 	r.POST("/articles/:id/unarchive", h.UnarchiveArticle)
// 	r.GET("/articles/:id/stats", h.GetArticleStats)
// 	r.GET("/articles/:id", h.GetArticleByID)
// 	r.GET("/me/articles", h.ListUserArticles)
// 	r.GET("/articles/trending", h.ListTrendingArticles)
// 	r.GET("/tags/:tag/articles", h.ListArticlesByTag)
// 	r.GET("/search", h.SearchArticles)
// 	r.GET("/articles", h.FilterArticles)
// 	r.POST("/articles/:id/clap", h.ClapArticle)
// 	r.DELETE("/me/trash", h.EmptyTrash)
// 	r.GET("/admin/articles", h.AdminListAllArticles)
// 	r.DELETE("/admin/articles/:id", h.AdminHardDeleteArticle)
// 	r.POST("/admin/articles/:id/unpublish", h.AdminUnpublishArticle)
// }

// func (s *ArticleControllerSuite) TestCreateArticle_Success() {
// 	input := controller.CreateArticleRequest{
// 		Title:    "Test Article",
// 		Slug:     "test-article",
// 		Excerpt:  "Preview",
// 		Language: "en",
// 		Tags:     []string{"go"},
// 		ContentBlocks: []controller.ContentBlockDTO{
// 			{
// 				Type:  "paragraph",
// 				Order: 0,
// 				Content: controller.BlockContentDTO{
// 					Paragraph: &controller.ParagraphContentDTO{Text: "Body", Style: "normal"},
// 				},
// 			},
// 		},
// 	}
// 	article := &domain.Article{ID: "123"}
// 	s.mockUsecase.On("CreateArticle", mock.Anything, "test_user", mock.AnythingOfType("*domain.Article")).Return(article, nil)
// 	body, _ := json.Marshal(input)
// 	req := httptest.NewRequest(http.MethodPost, "/articles", bytes.NewReader(body))
// 	req.Header.Set("Content-Type", "application/json")
// 	resp := httptest.NewRecorder()
// 	s.router.ServeHTTP(resp, req)
// 	s.Equal(http.StatusCreated, resp.Code)
// }

// func (s *ArticleControllerSuite) TestUpdateArticle_Success() {
// 	input := controller.UpdateArticleRequest{
// 		Title:    "Updated Title",
// 		Slug:     "updated-title",
// 		Excerpt:  "Updated excerpt",
// 		Language: "en",
// 		Tags:     []string{"go"},
// 	}
// 	article := &domain.Article{ID: "article123"}
// 	s.mockUsecase.On("UpdateArticle", mock.Anything, "test_user", mock.AnythingOfType("*domain.Article")).Return(article, nil)
// 	body, _ := json.Marshal(input)
// 	req := httptest.NewRequest("PUT", "/articles/article123", bytes.NewBuffer(body))
// 	req.Header.Set("Content-Type", "application/json")
// 	resp := httptest.NewRecorder()
// 	s.router.ServeHTTP(resp, req)
// 	s.Equal(http.StatusCreated, resp.Code)
// }

// func (s *ArticleControllerSuite) TestDeleteArticle_Success() {
// 	s.mockUsecase.On("DeleteArticle", mock.Anything, "test_user", "article123").Return(nil)
// 	req := httptest.NewRequest("DELETE", "/articles/article123", nil)
// 	resp := httptest.NewRecorder()
// 	s.router.ServeHTTP(resp, req)
// 	s.Equal(http.StatusNoContent, resp.Code)
// }

// func (s *ArticleControllerSuite) TestRestoreArticle_Success() {
// 	article := &domain.Article{ID: "article123"}
// 	s.mockUsecase.On("RestoreArticle", mock.Anything, "test_user", "article123").Return(article, nil)
// 	req := httptest.NewRequest("POST", "/articles/article123/restore", nil)
// 	resp := httptest.NewRecorder()
// 	s.router.ServeHTTP(resp, req)
// 	s.Equal(http.StatusNoContent, resp.Code)
// }

// func (s *ArticleControllerSuite) TestGetArticleStats_Success() {
// 	stats := domain.ArticleStats{ViewCount: 10, ClapCount: 3}
// 	s.mockUsecase.On("GetArticleStats", mock.Anything, "article123").Return(stats, nil)
// 	req := httptest.NewRequest("GET", "/articles/article123/stats", nil)
// 	resp := httptest.NewRecorder()
// 	s.router.ServeHTTP(resp, req)
// 	s.Equal(http.StatusOK, resp.Code)
// }

// func (s *ArticleControllerSuite) TestPublishArticle_Success() {
// 	article := &domain.Article{ID: "article123"}
// 	s.mockUsecase.On("PublishArticle", mock.Anything, "test_user", "article123").Return(article, nil)
// 	req := httptest.NewRequest("POST", "/articles/article123/publish", nil)
// 	resp := httptest.NewRecorder()
// 	s.router.ServeHTTP(resp, req)
// 	s.Equal(http.StatusOK, resp.Code)
// }

// func (s *ArticleControllerSuite) TestViewArticleBySlug_Success() {
// 	article := &domain.Article{ID: "article123", Status: domain.StatusPublished}
// 	s.mockUsecase.On("ViewArticleBySlug", mock.Anything, "my-article", mock.Anything).Return(article, nil)
// 	req := httptest.NewRequest("GET", "/p/my-article", nil)
// 	resp := httptest.NewRecorder()
// 	s.router.ServeHTTP(resp, req)
// 	s.Equal(http.StatusOK, resp.Code)
// }

// func (s *ArticleControllerSuite) TestListUserArticles_Success() {
// 	articles := []domain.Article{{ID: "a1"}, {ID: "a2"}}
// 	s.mockUsecase.On("ListUserArticles", mock.Anything, "test_user", "test_user", mock.Anything).Return(articles, 2, nil)
// 	req := httptest.NewRequest("GET", "/me/articles?page=1&page_size=10", nil)
// 	resp := httptest.NewRecorder()
// 	s.router.ServeHTTP(resp, req)
// 	s.Equal(http.StatusOK, resp.Code)
// }

// func (s *ArticleControllerSuite) TestEmptyTrash_Success() {
// 	s.mockUsecase.On("EmptyTrash", mock.Anything, "test_user").Return(nil)
// 	req := httptest.NewRequest("DELETE", "/me/trash", nil)
// 	resp := httptest.NewRecorder()
// 	s.router.ServeHTTP(resp, req)
// 	s.Equal(http.StatusNoContent, resp.Code)
// }

// func (s *ArticleControllerSuite) TestAdminListAllArticles_Success() {
// 	articles := []domain.Article{{ID: "a1"}, {ID: "a2"}}
// 	s.mockUsecase.On("AdminListAllArticles", mock.Anything, "test_user", "admin", mock.Anything).Return(articles, 2, nil)
// 	req := httptest.NewRequest("GET", "/admin/articles?page=1&page_size=10", nil)
// 	resp := httptest.NewRecorder()
// 	s.router.ServeHTTP(resp, req)
// 	s.Equal(http.StatusOK, resp.Code)
// }

// func (s *ArticleControllerSuite) TestAdminHardDeleteArticle_Success() {
// 	s.mockUsecase.On("AdminHardDeleteArticle", mock.Anything, "test_user", "admin", "article123").Return(nil)
// 	req := httptest.NewRequest("DELETE", "/admin/articles/article123", nil)
// 	resp := httptest.NewRecorder()
// 	s.router.ServeHTTP(resp, req)
// 	s.Equal(http.StatusNoContent, resp.Code)
// }

// func TestArticleControllerSuite(t *testing.T) {
// 	suite.Run(t, new(ArticleControllerSuite))
// }