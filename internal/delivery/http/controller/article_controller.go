package controller

import (
	"write_base/internal/domain"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	UseCase domain.IArticleUsecase
}

func NewArticleHandler(usecase *domain.IArticleUsecase) *ArticleHandler{
	return &ArticleHandler{
		UseCase: usecase,
	}
}
//===========================================================================//
//                 Article Lifecycle (Author only)                           //
//===========================================================================//
func (h *ArticleHandler) CreateArticle(ctx *gin.Context)
func (h *ArticleHandler) UpdateArticle(ctx *gin.Context)
// Soft delete (status=deleted)
func (h *ArticleHandler) DeleteArticle(ctx *gin.Context)
// Restore from trash
func (h *ArticleHandler) RestoreArticle(ctx *gin.Context)
//===========================================================================//
//                 Article Statistics (Author only)                           //
//===========================================================================//
// Views/claps analytics
func (h *ArticleHandler) GetArticleStats(ctx *gin.Context)    
//===========================================================================//
//             Article State Management (Author only)                        //
//===========================================================================//
func (h *ArticleHandler) PublishArticle(ctx *gin.Context)
func (h *ArticleHandler) UnpublishArticle(ctx *gin.Context)
func (h *ArticleHandler) ArchiveArticle(ctx *gin.Context)
func (h *ArticleHandler) UnarchiveArticle(ctx *gin.Context)
//===========================================================================//
//                         Article Retrieval                                 //
//===========================================================================//
// For authors/admins (with draft access)
func (h *ArticleHandler) GetArticleByID(ctx *gin.Context)    
// Public view with tracking
func (h *ArticleHandler) ViewArticleBySlug(ctx *gin.Context) 

//===========================================================================//
//                           Article Lists                                   //
//===========================================================================//
// Author's articles
func (h *ArticleHandler) ListUserArticles(ctx *gin.Context)      
// Popular content
func (h *ArticleHandler) ListTrendingArticles(ctx *gin.Context)  
// Tag-filtered
func (h *ArticleHandler) ListArticlesByTag(ctx *gin.Context)
// Full-text search     
func (h *ArticleHandler) SearchArticles(ctx *gin.Context) 
// Advanced filtering       
func (h *ArticleHandler) FilterArticles(ctx *gin.Context)        
//===========================================================================//
//                             Engagement                                    //
//===========================================================================//
func (h *ArticleHandler) ClapArticle(ctx *gin.Context)
// Generate share links
func (h *ArticleHandler) ShareArticle(ctx *gin.Context)       

//===========================================================================//
//                 Trash Management (Author only)                            //
//===========================================================================//
// Hard delete trash (cron alternative)
func (h *ArticleHandler) EmptyTrash(ctx *gin.Context)        
//===========================================================================//
//                         Admin Operations                                  //
//===========================================================================//
func (h *ArticleHandler) AdminListAllArticles(ctx *gin.Context)
func (h *ArticleHandler) AdminHardDeleteArticle(ctx *gin.Context)
func (h *ArticleHandler) AdminUnpublishArticle(ctx *gin.Context)