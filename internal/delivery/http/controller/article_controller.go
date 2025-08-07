package controller

import (
	"strconv"
	"net/http"
	"write_base/internal/domain"

	"github.com/gin-gonic/gin"

)

type ArticleHandler struct {
	UseCase domain.IArticleUsecase
}

func NewArticleHandler(usecase domain.IArticleUsecase) *ArticleHandler {
	return &ArticleHandler{
		UseCase: usecase,
	}
}

// ===========================================================================//
//	                Article Lifecycle (Author only)                           //
// ===========================================================================//
func (h *ArticleHandler) CreateArticle(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}
	userID := userIDVal.(string)
	var input CreateArticleRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	article := &domain.Article{
		Title:         input.Title,
		Slug:          input.Slug,
		Excerpt:       input.Excerpt,
		AuthorID:      userID,
		Language:      input.Language,
		Tags:          input.Tags,
		ContentBlocks: mapContentBlocks(input.ContentBlocks),
		Status:        domain.StatusDraft,
	}
	res, err := h.UseCase.CreateArticle(ctx, userID, article)
	if err != nil {
		switch err{
		case domain.ErrUnauthorized:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case domain.ErrInvalidRequest:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case domain.ErrUnauthorizedArticleEdit:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case domain.ErrArticleTagLimitExceeded:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"id": article.ID, "data": res})
}

func (h *ArticleHandler) UpdateArticle(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}
	userID := userIDVal.(string)
	var input UpdateArticleRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	articleID := ctx.Param("id")
	article := &domain.Article{
		Title:         input.Title,
		Slug:          input.Slug,
		Excerpt:       input.Excerpt,
		AuthorID:      userID,
		Language:      input.Language,
		Tags:          input.Tags,
		ContentBlocks: mapContentBlocks(input.ContentBlocks),
		Status:        domain.StatusDraft,
	}
	 article.ID = articleID
	res, err := h.UseCase.UpdateArticle(ctx,userID,article)
	if err != nil {
		switch err{
		case domain.ErrNoChangesDetected:
			ctx.JSON(http.StatusNoContent,gin.H{"error":err.Error()})
		case domain.ErrArticleNotFound:
			ctx.JSON(http.StatusNotFound,gin.H{"error": err.Error()})
		case domain.ErrUnauthorized:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case domain.ErrInvalidRequest:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case domain.ErrUnauthorizedArticleEdit:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case domain.ErrDuplicateArticleSlug:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case domain.ErrArticleTagLimitExceeded:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"id": article.ID, "data": res})
}

// Soft delete (status=deleted)
func (h *ArticleHandler) DeleteArticle(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}
	userID := userIDVal.(string)
	articleID := ctx.Param("id")
	if err:=h.UseCase.DeleteArticle(ctx,userID,articleID);err!=nil{
		switch err{
		case domain.ErrArticleNotFound:
			ctx.JSON(http.StatusNotFound,gin.H{"error": err.Error()})
		case domain.ErrUnauthorized:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case domain.ErrUnauthorizedArticleDelete:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusNoContent,nil)
}

// Restore from trash
func (h *ArticleHandler) RestoreArticle(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}
	userID := userIDVal.(string)
	articleID := ctx.Param("id")
	article, err:=h.UseCase.RestoreArticle(ctx,userID,articleID)
	if err!=nil{
		switch err{
		case domain.ErrArticleNotFound:
			ctx.JSON(http.StatusNotFound,gin.H{"error": err.Error()})
		case domain.ErrUnauthorized:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case domain.ErrUnauthorizedArticleEdit:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusNoContent,gin.H{"id": article.ID, "data": article})
}

//===========================================================================//
//	              Article Statistics (Author only)                           //
//===========================================================================//
// Views/claps analytics
// ===========================================================================//
// Article Statistics (Author only)                           //
// ===========================================================================//
func (h *ArticleHandler) GetArticleStats(ctx *gin.Context) {
	articleID := ctx.Param("id")
	if articleID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "article ID is required"})
		return
	}

	stats, err := h.UseCase.GetArticleStats(ctx, articleID)
	if err != nil {
		switch err {
		case domain.ErrArticleNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case domain.ErrUnauthorized:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, stats)
}

// ===========================================================================//
// Article State Management (Author only)                        //
// ===========================================================================//
func (h *ArticleHandler) PublishArticle(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}
	userID := userIDVal.(string)

	articleID := ctx.Param("id")
	if articleID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "article ID is required"})
		return
	}

	article, err := h.UseCase.PublishArticle(ctx, userID, articleID)
	if err != nil {
		switch err {
		case domain.ErrArticleNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case domain.ErrUnauthorized, domain.ErrUnauthorizedArticleEdit:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case domain.ErrArticleAlreadyPublished:
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, article)
}

func (h *ArticleHandler) UnpublishArticle(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}
	userID := userIDVal.(string)

	articleID := ctx.Param("id")
	if articleID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "article ID is required"})
		return
	}

	article, err := h.UseCase.UnpublishArticle(ctx, userID, articleID, false)
	if err != nil {
		switch err {
		case domain.ErrArticleNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case domain.ErrUnauthorized, domain.ErrUnauthorizedArticleEdit:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case domain.ErrArticleNotPublished:
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, article)
}

func (h *ArticleHandler) ArchiveArticle(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}
	userID := userIDVal.(string)

	articleID := ctx.Param("id")
	if articleID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "article ID is required"})
		return
	}

	article, err := h.UseCase.ArchiveArticle(ctx, userID, articleID)
	if err != nil {
		switch err {
		case domain.ErrArticleNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case domain.ErrUnauthorized, domain.ErrUnauthorizedArticleEdit:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case domain.ErrArticleNotPublished:
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, article)
}

func (h *ArticleHandler) UnarchiveArticle(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}
	userID := userIDVal.(string)

	articleID := ctx.Param("id")
	if articleID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "article ID is required"})
		return
	}

	article, err := h.UseCase.UnarchiveArticle(ctx, userID, articleID)
	if err != nil {
		switch err {
		case domain.ErrArticleNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case domain.ErrUnauthorized, domain.ErrUnauthorizedArticleEdit:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case domain.ErrArticleNotArchived:
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, article)
}

// ===========================================================================//
// Article Retrieval                                 //
// ===========================================================================//
func (h *ArticleHandler) GetArticleByID(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}
	userID := userIDVal.(string)

	roleVal, exists := ctx.Get("user_role")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_role not found in context"})
		return
	}
	userRole := roleVal.(string)

	articleID := ctx.Param("id")
	if articleID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "article ID is required"})
		return
	}

	article, err := h.UseCase.GetArticleByID(ctx, userID, articleID, userRole)
	if err != nil {
		switch err {
		case domain.ErrArticleNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case domain.ErrUnauthorized, domain.ErrForbidden:
			ctx.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, article)
}

func (h *ArticleHandler) ViewArticleBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	if slug == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "slug is required"})
		return
	}

	clientIP := ctx.ClientIP()
	article, err := h.UseCase.ViewArticleBySlug(ctx, slug, clientIP)
	if err != nil {
		if err == domain.ErrArticleNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Don't return content for unpublished articles
	if article.Status != domain.StatusPublished {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
		return
	}

	ctx.JSON(http.StatusOK, article)
}

// ===========================================================================//
//                            Article Lists                                   //
// ===========================================================================//
type PaginationRequest struct {
	Page     int `form:"page" binding:"min=1"`
	PageSize int `form:"page_size" binding:"min=1,max=100"`
}

func (h *ArticleHandler) ListUserArticles(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}
	userID := userIDVal.(string)
	var pagReq PaginationRequest
	if err := ctx.ShouldBindQuery(&pagReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pagination := domain.Pagination{
		Page:     pagReq.Page,
		PageSize: pagReq.PageSize,
	}
	// List articles for the current user
	articles, total, err := h.UseCase.ListUserArticles(ctx, userID, userID, pagination)
	if err != nil {
		if err == domain.ErrArticleNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := gin.H{
		"data":       articles,
		"total":      total,
		"page":       pagination.Page,
		"page_size":  pagination.PageSize,
		"total_pages": (total + pagination.PageSize - 1) / pagination.PageSize,
	}
	ctx.JSON(http.StatusOK, response)
}

func (h *ArticleHandler) ListTrendingArticles(ctx *gin.Context) {
	var pagReq PaginationRequest
	if err := ctx.ShouldBindQuery(&pagReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	windowDays := 7 // default window
	if daysStr := ctx.Query("days"); daysStr != "" {
		if days, err := strconv.Atoi(daysStr); err == nil && days > 0 {
			windowDays = days
		}
	}
	pagination := domain.Pagination{
		Page:     pagReq.Page,
		PageSize: pagReq.PageSize,
	}
	articles, total, err := h.UseCase.ListTrendingArticles(ctx, pagination, windowDays)
	if err != nil {
		if err == domain.ErrArticleNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := gin.H{
		"data":       articles,
		"total":      total,
		"page":       pagination.Page,
		"page_size":  pagination.PageSize,
		"total_pages": (total + pagination.PageSize - 1) / pagination.PageSize,
	}
	ctx.JSON(http.StatusOK, response)
}

func (h *ArticleHandler) ListArticlesByTag(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}
	userID := userIDVal.(string)

	tag := ctx.Param("tag")
	if tag == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "tag is required"})
		return
	}

	var pagReq PaginationRequest
	if err := ctx.ShouldBindQuery(&pagReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pagination := domain.Pagination{
		Page:     pagReq.Page,
		PageSize: pagReq.PageSize,
	}

	articles, total, err := h.UseCase.ListArticlesByTag(ctx, userID, tag, pagination)
	if err != nil {
		if err == domain.ErrArticleNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"data":       articles,
		"total":      total,
		"page":       pagination.Page,
		"page_size":  pagination.PageSize,
		"total_pages": (total + pagination.PageSize - 1) / pagination.PageSize,
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *ArticleHandler) SearchArticles(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}
	userID := userIDVal.(string)

	query := ctx.Query("q")
	if query == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "search query is required"})
		return
	}

	var pagReq PaginationRequest
	if err := ctx.ShouldBindQuery(&pagReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pagination := domain.Pagination{
		Page:     pagReq.Page,
		PageSize: pagReq.PageSize,
	}

	articles, total, err := h.UseCase.SearchArticles(ctx, userID, query, pagination)
	if err != nil {
		if err == domain.ErrArticleNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"data":       articles,
		"total":      total,
		"page":       pagination.Page,
		"page_size":  pagination.PageSize,
		"total_pages": (total + pagination.PageSize - 1) / pagination.PageSize,
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *ArticleHandler) FilterArticles(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}
	userID := userIDVal.(string)

	var filter domain.ArticleFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var pagReq PaginationRequest
	if err := ctx.ShouldBindQuery(&pagReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pagination := domain.Pagination{
		Page:     pagReq.Page,
		PageSize: pagReq.PageSize,
	}

	articles, total, err := h.UseCase.FilterArticles(ctx, userID, filter, pagination)
	if err != nil {
		if err == domain.ErrArticleNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"data":       articles,
		"total":      total,
		"page":       pagination.Page,
		"page_size":  pagination.PageSize,
		"total_pages": (total + pagination.PageSize - 1) / pagination.PageSize,
	}

	ctx.JSON(http.StatusOK, response)
}

// ===========================================================================//
// Engagement                                    //
// ===========================================================================//
func (h *ArticleHandler) ClapArticle(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}
	userID := userIDVal.(string)

	articleID := ctx.Param("id")
	if articleID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "article ID is required"})
		return
	}

	stats, err := h.UseCase.ClapArticle(ctx, articleID, userID)
	if err != nil {
		switch err {
		case domain.ErrArticleNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case domain.ErrUnauthorized:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case domain.ErrClapLimitExceeded:
			ctx.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, stats)
}

// ===========================================================================//
// Trash Management (Author only)                            //
// ===========================================================================//
func (h *ArticleHandler) EmptyTrash(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}
	userID := userIDVal.(string)

	if err := h.UseCase.EmptyTrash(ctx, userID); err != nil {
		switch err {
		case domain.ErrArticleNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case domain.ErrUnauthorized:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// ===========================================================================//
// Admin Operations                                  //
// ===========================================================================//
func (h *ArticleHandler) AdminListAllArticles(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}
	userID := userIDVal.(string)

	roleVal, exists := ctx.Get("user_role")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_role not found in context"})
		return
	}
	userRole := roleVal.(string)

	var pagReq PaginationRequest
	if err := ctx.ShouldBindQuery(&pagReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pagination := domain.Pagination{
		Page:     pagReq.Page,
		PageSize: pagReq.PageSize,
	}

	articles, total, err := h.UseCase.AdminListAllArticles(ctx, userID, userRole, pagination)
	if err != nil {
		switch err {
		case domain.ErrUnauthorized:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case domain.ErrArticleNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	response := gin.H{
		"data":       articles,
		"total":      total,
		"page":       pagination.Page,
		"page_size":  pagination.PageSize,
		"total_pages": (total + pagination.PageSize - 1) / pagination.PageSize,
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *ArticleHandler) AdminHardDeleteArticle(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}
	userID := userIDVal.(string)

	roleVal, exists := ctx.Get("user_role")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_role not found in context"})
		return
	}
	userRole := roleVal.(string)

	articleID := ctx.Param("id")
	if articleID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "article ID is required"})
		return
	}

	if err := h.UseCase.AdminHardDeleteArticle(ctx, userID, userRole, articleID); err != nil {
		switch err {
		case domain.ErrUnauthorized:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case domain.ErrArticleNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (h *ArticleHandler) AdminUnpublishArticle(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}
	userID := userIDVal.(string)
	roleVal, exists := ctx.Get("user_role")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_role not found in context"})
		return
	}
	userRole := roleVal.(string)
	articleID := ctx.Param("id")
	if articleID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "article ID is required"})
		return
	}
	article, err := h.UseCase.AdminUnpublishArticle(ctx, userID, userRole, articleID)
	if err != nil {
		switch err {
		case domain.ErrUnauthorized:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case domain.ErrArticleNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case domain.ErrArticleNotPublished:
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusOK, article)
}