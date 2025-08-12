package router

import (
	"write_base/internal/delivery/http/controller"

	"github.com/gin-gonic/gin"
)

func RegisterArticleRouter(r *gin.Engine, h  *controller.Handler)  {
	userAuthGroup := r.Group("/")
	{
		userAuthGroup.POST("/articles/new",h.CreateArticle)
		userAuthGroup.PUT("/articles/:id", h.UpdateArticle)
		userAuthGroup.DELETE("/articles/:id", h.DeleteArticle)
		userAuthGroup.PATCH("/articles/:id/restore", h.RestoreArticle)

		userAuthGroup.GET("/articles/:id", h.GetArticleByID)
		// Article state management
		userAuthGroup.POST("/articles/:id/publish", h.PublishArticle)
		userAuthGroup.POST("/articles/:id/unpublish", h.UnpublishArticle)
		userAuthGroup.POST("/articles/:id/archive", h.ArchiveArticle)
		userAuthGroup.POST("/articles/:id/unarchive", h.UnarchiveArticle)
		// Statistics
		userAuthGroup.GET("/articles/:id/stats", h.GetArticleStats)
		userAuthGroup.GET("/articles/stats/all", h.GetAllArticleStats)
		userAuthGroup.GET(":slug",h.GetArticleBySlug)
		// List Articles
		userAuthGroup.GET("/authors/:author_id/articles", h.ListArticlesByAuthor)
		userAuthGroup.GET("/articles/trending", h.GetTrendingArticles)
		userAuthGroup.GET("/articles/new", h.GetNewArticles)
		userAuthGroup.GET("/articles/popular", h.GetPopularArticles)
		// Author filter
		userAuthGroup.POST("/authors/:author_id/articles/filter", h.FilterAuthorArticles)

		// For all users
		userAuthGroup.POST("/articles/filter", h.FilterArticles)
		// Search
		userAuthGroup.GET("/search", h.SearchArticles)

		userAuthGroup.GET("/article/tags", h.ListArticlesByTags)
		// Trash management
		userAuthGroup.DELETE("/me/trash", h.EmptyTrash)
		userAuthGroup.DELETE("/articles/trash/:id", h.DeleteFromTrash)

		userAuthGroup.POST("articles/:id/clap", h.AddClap)

		userAuthGroup.POST("/generateslug", h.GenerateSlug)
		userAuthGroup.POST("/articles/generatecontent", h.GenerateContent)
	}
	adminGroup := r.Group("/admin")
	{
		adminGroup.GET("/articles", h.AdminListAllArticles)
		adminGroup.DELETE("/articles/:id/delete", h.AdminHardDeleteArticle)
		adminGroup.POST("/articles/:id/unpublish", h.AdminUnpublishArticle)
	}
}

