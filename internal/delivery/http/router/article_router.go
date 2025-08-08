package router

import (
	"write_base/internal/delivery/http/controller"

	"github.com/gin-gonic/gin"
)

func RegisterArticleRouter(r *gin.Engine, h *controller.ArticleHandler) {
	// Public routes
	r.GET("/p/:slug", h.ViewArticleBySlug)

	// Authenticated routes
	authGroup := r.Group("/")
	{
		// Article lifecycle
		authGroup.POST("/articles", h.CreateArticle)
		authGroup.PUT("/articles/:id", h.UpdateArticle)
		authGroup.DELETE("/articles/:id", h.DeleteArticle)
		authGroup.POST("/articles/:id/restore", h.RestoreArticle)
		
		// Article state management
		authGroup.POST("/articles/:id/publish", h.PublishArticle)
		authGroup.POST("/articles/:id/unpublish", h.UnpublishArticle)
		authGroup.POST("/articles/:id/archive", h.ArchiveArticle)
		authGroup.POST("/articles/:id/unarchive", h.UnarchiveArticle)
		
		// Statistics
		authGroup.GET("/articles/:id/stats", h.GetArticleStats)
		
		// Article retrieval
		authGroup.GET("/articles/:id", h.GetArticleByID)
		
		// Article lists
		authGroup.GET("/me/articles", h.ListUserArticles)
		authGroup.GET("/articles/trending", h.ListTrendingArticles)
		authGroup.GET("/tags/:tag/articles", h.ListArticlesByTag)
		authGroup.GET("/search", h.SearchArticles)
		authGroup.GET("/articles", h.FilterArticles)
		
		// Engagement
		authGroup.POST("/articles/:id/clap", h.ClapArticle)
		
		// Trash management
		authGroup.DELETE("/me/trash", h.EmptyTrash)
	}

	// Admin routes
	adminGroup := r.Group("/admin")
	{
		adminGroup.GET("/articles", h.AdminListAllArticles)
		adminGroup.DELETE("/articles/:id", h.AdminHardDeleteArticle)
		adminGroup.POST("/articles/:id/unpublish", h.AdminUnpublishArticle)
	}
}