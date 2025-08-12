package controller

import (
	"net/http"
	"strings"
	"write_base/internal/domain"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Usecase domain.IArticleUsecase
}

func NewArticleHandler(uc domain.IArticleUsecase) *Handler {
	return &Handler{Usecase: uc}
}
// =============================== Article Create ================================
func (h *Handler) CreateArticle(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
	    ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	    return
	}
	userID := userIDVal.(string)
    

    var articleReq *ArticleRequest
    if err := ctx.ShouldBindJSON(&articleReq); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    article := articleReq.ToDomain()
    res, err := h.Usecase.CreateArticle(ctx, userID, article)
    if err != nil {
        code := http.StatusInternalServerError
        switch err {
		case domain.ErrInvalidTagName:
			code = http.StatusBadRequest
        case domain.ErrInvalidArticlePayload:
            code = http.StatusBadRequest
        case domain.ErrUnauthorized:
            code = http.StatusUnauthorized
        }
        ctx.JSON(code, gin.H{"error": err.Error()})
        return
    }

    // Used response format
    ctx.JSON(http.StatusCreated, gin.H{
        "data": gin.H{
            "id": res,
        },
    })
}
// =============================== Article Update ================================
func (h *Handler) UpdateArticle(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
	    ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	    return
	}
	userID := userIDVal.(string)
    
	
    articleID := ctx.Param("id")

    var articleReq ArticleUpdateRequest
    if err := ctx.ShouldBindJSON(&articleReq); err != nil {
        ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    article := articleReq.ToDomain()
    article.ID = articleID
    article.AuthorID = userID

    if err := h.Usecase.UpdateArticle(ctx, userID, article); err != nil {
        code := http.StatusInternalServerError
        switch err {
		case domain.ErrInvalidTagName:
			code = http.StatusBadRequest
        case domain.ErrInvalidArticlePayload:
            code = http.StatusBadRequest
        case domain.ErrUnauthorized:
            code = http.StatusUnauthorized
        case domain.ErrArticleNotFound:
            code = http.StatusNotFound
        }
        ctx.IndentedJSON(code, gin.H{"error": err.Error()})
        return
    }

    articleDTO := new(ArticleResponse)
    articleDTO.ToDTO(article)

    ctx.JSON(http.StatusOK, gin.H{
        "data": articleDTO,
    })
}
// =============================== Article Delete ================================
func (h *Handler) DeleteArticle(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
	    ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	    return
	}
	userID := userIDVal.(string)
    
    articleID := ctx.Param("id")
    if err := h.Usecase.DeleteArticle(ctx,articleID,userID);err!=nil {
        code := http.StatusInternalServerError
        switch err {
        case domain.ErrUnauthorized:
            code = http.StatusUnauthorized
        case domain.ErrArticleNotFound:
            code = http.StatusNotFound
        }
        ctx.IndentedJSON(code, gin.H{"error": err.Error()})
        return
    }
    ctx.IndentedJSON(http.StatusOK, gin.H{"message":"successfully deleted"})
}
// =============================== Article Restore ================================
func (h *Handler) RestoreArticle(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
	    ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	    return
	}
	userID := userIDVal.(string)
    
    articleID := ctx.Param("id")
    if err := h.Usecase.RestoreArticle(ctx,articleID,userID);err!=nil {
        code := http.StatusInternalServerError
        switch err {
        case domain.ErrUnauthorized:
            code = http.StatusUnauthorized
        case domain.ErrArticleNotFound:
            code = http.StatusNotFound
        }
        ctx.IndentedJSON(code, gin.H{"error": err.Error()})
        return
    }
    ctx.IndentedJSON(http.StatusOK, gin.H{"message":"successfully restored"})
}
//===============================================================================//
//                   Article State Management                                    //
//===============================================================================//
// ======================== Article Publish =======================================
func (h *Handler) PublishArticle(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
	    ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	    return
	}
	userID := userIDVal.(string)
    
	articleID := ctx.Param("id")
	if articleID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrArticleInvalidID})
		return
	}

	article, err := h.Usecase.PublishArticle(ctx, articleID, userID)
	if err != nil {
		switch err {
		case domain.ErrUnapprovedTags:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case domain.ErrArticleNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case domain.ErrUnauthorized:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case domain.ErrArticlePublished:
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
    articleDTO := new(ArticleResponse)
    articleDTO.ToDTO(article)

    ctx.JSON(http.StatusOK, gin.H{
        "data": articleDTO,
    })
}
// ======================== Article Unpublish =====================================
func (h *Handler) UnpublishArticle(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
	    ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	    return
	}
	userID := userIDVal.(string)
    
	articleID := ctx.Param("id")
	if articleID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrArticleInvalidID})
		return
	}

	article, err := h.Usecase.UnpublishArticle(ctx, articleID, userID)
	if err != nil {
		switch err {
		case domain.ErrArticleNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case domain.ErrUnauthorized:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case domain.ErrArticleNotPublished:
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
    articleDTO := new(ArticleResponse)
    articleDTO.ToDTO(article)

    ctx.JSON(http.StatusOK, gin.H{
        "data": articleDTO,
    })
}
// ======================== Article Archive =======================================
func (h *Handler) ArchiveArticle(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
	    ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	    return
	}
	userID := userIDVal.(string)
    
	articleID := ctx.Param("id")
	if articleID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrArticleInvalidID})
		return
	}
	article, err := h.Usecase.ArchiveArticle(ctx, articleID, userID)
	if err != nil {
		switch err {
		case domain.ErrArticleNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case domain.ErrUnauthorized:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case domain.ErrArticleArchived:
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
    articleDTO := new(ArticleResponse)
    articleDTO.ToDTO(article)

    ctx.JSON(http.StatusOK, gin.H{
        "data": articleDTO,
    })
}
// ======================== Article Unarchive =====================================
func (h *Handler) UnarchiveArticle(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
	    ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	    return
	}
	userID := userIDVal.(string)
    
	articleID := ctx.Param("id")
	if articleID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrArticleInvalidID})
		return
	}

	article, err := h.Usecase.UnarchiveArticle(ctx, articleID, userID)
	if err != nil {
		switch err {
		case domain.ErrArticleNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case domain.ErrUnauthorized:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case domain.ErrArticleNotArchived:
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
    articleDTO := new(ArticleResponse)
    articleDTO.ToDTO(article)

    ctx.JSON(http.StatusOK, gin.H{
        "data": articleDTO,
    })
}
//===============================================================================//
//                               Retrieve                                        //
//===============================================================================//
// =============================== Article GetByID ================================
func (h *Handler) GetArticleByID(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
	    ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	    return
	}
	userID := userIDVal.(string)
    
    articleID := ctx.Param("id")
    res, err := h.Usecase.GetArticleByID(ctx,articleID,userID)
    if err!=nil {
        code := http.StatusInternalServerError
        switch err {
        case domain.ErrInvalidArticlePayload:
            code = http.StatusBadRequest
        case domain.ErrUnauthorized:
            code = http.StatusUnauthorized
        case domain.ErrArticleNotFound:
            code = http.StatusNotFound
        }
        ctx.IndentedJSON(code, gin.H{"error": err.Error()})
        return
    }
    articleDTO := new(ArticleResponse)
    articleDTO.ToDTO(res)

    ctx.JSON(http.StatusOK, gin.H{
        "data": articleDTO,
    })
}
// =================== Article Stats of Author ====================================
func (h *Handler) GetArticleStats(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
	    ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	    return
	}
	userID := userIDVal.(string)
    
	articleID := ctx.Param("id")
	if articleID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrArticleInvalidID})
		return
	}

	stats, err := h.Usecase.GetArticleStats(ctx,articleID,userID)
	if err != nil {
        code := http.StatusInternalServerError
        switch err {
        case domain.ErrInvalidArticlePayload:
            code = http.StatusBadRequest
        case domain.ErrArticleNotFound:
            code = http.StatusNotFound
        }
        ctx.IndentedJSON(code, gin.H{"error": err.Error()})
        return
	}

	ctx.JSON(http.StatusOK, stats)
}
// =================== All Article Stats of Author ================================
func (h *Handler) GetAllArticleStats(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
	    ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	    return
	}
	userID := userIDVal.(string)
    

	stats,total, err := h.Usecase.GetAllArticleStats(ctx,userID)
	if err != nil {
        code := http.StatusInternalServerError
        switch err {
        case domain.ErrInvalidArticlePayload:
            code = http.StatusBadRequest
        case domain.ErrArticleNotFound:
            code = http.StatusNotFound
        }
        ctx.IndentedJSON(code, gin.H{"error": err.Error()})
        return
	}

	ctx.JSON(http.StatusOK, 
        gin.H{
            "data":stats,
            "total":total,
        },
    )
}
// =================== All Article Stats of Author ================================
func (h *Handler) GetArticleBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	if slug == "" || strings.Contains(slug," ") {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrArticleInvalidSlug})
		return
	}
	//clientIP := ctx.ClientIP()
	clientIP := "192.1.1.1"
    article,err := h.Usecase.GetArticleBySlug(ctx,slug,clientIP)
    if err!=nil {
        code := http.StatusInternalServerError
        switch err {
        case domain.ErrArticleNotFound:
            code = http.StatusNotFound
        }
        ctx.IndentedJSON(code, gin.H{"error": err.Error()})
        return
    }
     articleDTO := new(ArticleResponse)
    articleDTO.ToDTO(article)

    ctx.JSON(http.StatusOK, gin.H{
        "data": articleDTO,
    })
}

//===========================================================================//
//                           Article Lists                                   //
//===========================================================================//
//======================= List user articles ==================================
func (h *Handler) ListArticlesByAuthor(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
	    ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	    return
	}
	userID := userIDVal.(string)
    
	var pagReq PaginationRequest
	if err := ctx.ShouldBindQuery(&pagReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
    articleID := ctx.Param("author_id")
	if articleID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrArticleInvalidID})
		return
	}
	pagination := domain.Pagination{
		Page:     pagReq.Page,
		PageSize: pagReq.PageSize,
	}
	
	// Validate pagination parameters
	pagination.ValidatePagination()
	
	// List articles for the current user
	articles, total, err := h.Usecase.ListArticlesByAuthor(ctx, userID, articleID, pagination)
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
//======================= List Trending articles ==================================
// Trending articles (last 7 days)
func (h *Handler) GetTrendingArticles(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
	    ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	    return
	}
	userID := userIDVal.(string)
    
	var pagReq PaginationRequest
	if err := ctx.ShouldBindQuery(&pagReq); err != nil {
        if err == domain.ErrArticleNotFound {}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pagination := domain.Pagination{
		Page:     pagReq.Page,
		PageSize: pagReq.PageSize,
	}
	
	// Validate pagination parameters
	pagination.ValidatePagination()

	articles, total, err := h.Usecase.GetTrendingArticles(ctx, userID, pagination)
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

//======================= List New articles ==================================

func (h *Handler) GetNewArticles(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
	    ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
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
	
	// Validate pagination parameters
	pagination.ValidatePagination()

	articles, total, err := h.Usecase.GetNewArticles(ctx, userID, pagination)
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

//======================= List Popular articles ==================================

func (h *Handler) GetPopularArticles(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
	    ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
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
	
	// Validate pagination parameters
	pagination.ValidatePagination()

	articles, total, err := h.Usecase.GetPopularArticles(ctx, userID, pagination)
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
//======================= List Author articles ==================================
func (h *Handler) FilterAuthorArticles(c *gin.Context) {
    callerIDVal, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": domain.ErrUnauthorized.Error()})
        return
    }
    callerID := callerIDVal.(string)
    

    authorID := c.Param("author_id")
    if authorID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "author_id required"})
        return
    }

    var req struct {
        Filter     domain.ArticleFilter `json:"filter"`
        Pagination domain.Pagination    `json:"pagination"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrInvalidArticlePayload.Error()})
        return
    }

    // Defensive pagination defaults
    if req.Pagination.Page < 1 {
        req.Pagination.Page = 1
    }
    if req.Pagination.PageSize <= 0 {
        req.Pagination.PageSize = 20
    }

    articles, total, err := h.Usecase.FilterAuthorArticles(c.Request.Context(), callerID, authorID, req.Filter, req.Pagination)
    if err != nil {
        switch err {
        case domain.ErrUnauthorized:
            c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        case domain.ErrAuthorNotFound:
            c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        case domain.ErrArticleNotFound:
            c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"error": domain.ErrInternalServer.Error()})
        }
        return
    }

    resp := make([]ArticleListResponse, 0, len(articles))
    for _, a := range articles {
        var dto ArticleListResponse
        dto.ToListDTO(a)
        resp = append(resp, dto)
    }

    c.JSON(http.StatusOK, gin.H{
        "data":      resp,
        "total":     total,
        "page":      req.Pagination.Page,
        "page_size": req.Pagination.PageSize,
    })
}

//======================= List Author articles (all users) =============================
func (h *Handler) FilterArticles(ctx *gin.Context) {
    var req struct {
        Filter     domain.ArticleFilter `json:"filter"`
        Pagination domain.Pagination    `json:"pagination"`
    }
    _, exists := ctx.Get("user_id")
	if !exists {
	    ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	    return
	}

    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrInvalidArticlePayload.Error()})
        return
    }

    // Defensive pagination defaults
    if req.Pagination.Page < 1 {
        req.Pagination.Page = 1
    }
    if req.Pagination.PageSize <= 0 {
        req.Pagination.PageSize = 20
    }

    articles, total, err := h.Usecase.FilterArticles(ctx.Request.Context(), req.Filter, req.Pagination)
    if err != nil {
        switch err {
        case domain.ErrArticleNotFound:
            ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        case domain.ErrInvalidArticlePayload:
            ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        default:
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": domain.ErrInternalServer.Error()})
        }
        return
    }

    resp := make([]ArticleListResponse, 0, len(articles))
    for _, a := range articles {
        var dto ArticleListResponse
        dto.ToListDTO(a)
        resp = append(resp, dto)
    }

    ctx.JSON(http.StatusOK, gin.H{
        "data":      resp,
        "total":     total,
        "page":      req.Pagination.Page,
        "page_size": req.Pagination.PageSize,
    })
}
//=========================== Search ==============================================
func (h *Handler) SearchArticles(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}
	userID := userIDVal.(string)
    

	query := strings.TrimSpace(ctx.Query("q"))
	if query == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "search query is required"})
		return
	}

    var pagReq PaginationRequest
    if err := ctx.ShouldBindQuery(&pagReq); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    pagination := domain.Pagination{Page: pagReq.Page, PageSize: pagReq.PageSize}
    pagination.ValidatePagination()

    articles, total, err := h.Usecase.SearchArticles(ctx.Request.Context(), userID, query, pagination)
    if err != nil {
        if err == domain.ErrArticleNotFound {
            ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
            return
        }
        if err == domain.ErrUnauthorized {
            ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
            return
        }
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "data":        articles,
        "total":       total,
        "page":        pagination.Page,
        "page_size":   pagination.PageSize,
        "total_pages": (total + pagination.PageSize - 1) / pagination.PageSize,
    })
}
//======================== List By Tags =======================================
func (h *Handler) ListArticlesByTags(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
	    c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	    return
	}
	userID := userIDVal.(string)
	 // temporary until auth is implemented
	
	// Get tags from query params
	tags := c.QueryArray("tags")
	if len(tags) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one tag is required"})
		return
	}
	
	var pagReq PaginationRequest
	if err := c.ShouldBindQuery(&pagReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	pagination := domain.Pagination{
		Page:     pagReq.Page,
		PageSize: pagReq.PageSize,
	}
	pagination.ValidatePagination()
	
	articles, total, err := h.Usecase.ListArticlesByTags(c.Request.Context(), userID, tags, pagination)
	if err != nil {
		code := http.StatusInternalServerError
		switch err {
		case domain.ErrUnapprovedTags:
			code = http.StatusBadRequest
		case domain.ErrUnauthorized:
			code = http.StatusUnauthorized
		case domain.ErrUnapprovedTags:
			code = http.StatusBadRequest
		case domain.ErrArticleNotFound:
			code = http.StatusNotFound
		}
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}
	
	response := gin.H{
		"data":        articles,
		"total":       total,
		"page":        pagination.Page,
		"page_size":   pagination.PageSize,
		"total_pages": (total + pagination.PageSize - 1) / pagination.PageSize,
	}
	c.JSON(http.StatusOK, response)
}
// ===========================================================================//
//                  Trash Management (Author only)                            //
// ===========================================================================//
func (h *Handler) EmptyTrash(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}
	userID := userIDVal.(string)
    

	if err := h.Usecase.EmptyTrash(ctx, userID); err != nil {
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
	ctx.Status(http.StatusNoContent)
}
//======================= Delete From Trash ====================================
func (h *Handler) DeleteFromTrash(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}
	userID := userIDVal.(string)
    
    articleID := ctx.Param("id")
    if articleID == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "article_id required"})
        return
    }

	if err := h.Usecase.DeleteArticleFromTrash(ctx, articleID, userID); err != nil {
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
	ctx.Status(http.StatusNoContent)
}
// ===========================================================================//
//                            Admin Operations                                //
// ===========================================================================//
func (h *Handler) AdminListAllArticles(ctx *gin.Context) {
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
	
	// Validate pagination parameters
	pagination.ValidatePagination()

	articles, total, err := h.Usecase.AdminListAllArticles(ctx, userID, userRole, pagination)
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

func (h *Handler) AdminHardDeleteArticle(ctx *gin.Context) {
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

	if err := h.Usecase.AdminHardDeleteArticle(ctx, userID, userRole, articleID); err != nil {
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

func (h *Handler) AdminUnpublishArticle(ctx *gin.Context) {
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
	article, err := h.Usecase.AdminUnpublishArticle(ctx, userID, userRole, articleID)
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

//=================================== CLAPPING ===================================

// Add new handler method
func (h *Handler) AddClap(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	
	articleID := c.Param("id")
	if articleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Article ID required"})
		return
	}
	
	stats, err := h.Usecase.AddClap(c.Request.Context(), userID, articleID)
	if err != nil {
		code := http.StatusInternalServerError
		switch err {
		case domain.ErrUnauthorized:
			code = http.StatusUnauthorized
		case domain.ErrClapLimitExceeded:
			code = http.StatusTooManyRequests
		case domain.ErrArticleNotFound:
			code = http.StatusNotFound
		}
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, stats)
}
