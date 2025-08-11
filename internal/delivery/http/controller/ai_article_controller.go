package controller

import (
	"net/http"
	"write_base/internal/domain"

	"github.com/gin-gonic/gin"
)

// ---------------- DTOs ----------------
type GenerateSlugRequest struct {
	Title string `json:"title" binding:"required"`
}

type GenerateSlugResponse struct {
	Slug string `json:"slug"`
}

type EditContentRequest struct {
	ArticleID    string            `json:"article_id" binding:"required"`
	Instructions string            `json:"instructions" binding:"required"`
	Content      []ContentBlockDTO `json:"content_blocks" binding:"required,min=1,dive"`
}

type EditContentResponse struct {
	ContentBlocks []ContentBlockDTO `json:"content_blocks"`
}

// ------------- Handlers --------------

func (h *Handler) GenerateSlug(ctx *gin.Context) {
	var req GenerateSlugRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slug, err := h.Usecase.GenerateSlugForTitle(ctx.Request.Context(), req.Title)
	if err != nil {
		// Map domain errors to HTTP codes where appropriate
		switch err {
		case domain.ErrInvalidArticlePayload:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, GenerateSlugResponse{Slug: slug})
}

//========================== Generate Content ========================================
type GenerateContentRequest struct {
	ArticleID    string            `json:"article_id,omitempty"`
	Title        string            `json:"title,omitempty"`
	Excerpt      string            `json:"excerpt,omitempty"`
	Language     string            `json:"language,omitempty"`
	Tags         []string          `json:"tags,omitempty"`
	Instructions string            `json:"instructions" binding:"required"`
	Content      []ContentBlockDTO `json:"content_blocks" binding:"required,min=1,dive"`
}

type GenerateContentResponse struct {
	Article ArticleResponse `json:"article"`
}

func (h *Handler) GenerateContent(ctx *gin.Context) {
	var req GenerateContentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Build a domain.Article from request (reuse your mapper)
	article := &domain.Article{
		ID:            req.ArticleID,
		Title:         req.Title,
		Slug:          "", // slug may be produced separately
		Excerpt:       req.Excerpt,
		Language:      req.Language,
		Tags:          req.Tags,
		ContentBlocks: mapContentBlocks(req.Content),
	}

	edited, err := h.Usecase.GenerateContentForArticle(ctx.Request.Context(), article, req.Instructions)
	if err != nil {
		// map known errors
		switch err {
		case domain.ErrArticleNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case domain.ErrUnauthorized:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		default:
			// check usecase-specific policy error
			if err == domain.ErrContentPolicyViolation {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "content violates policy"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Convert to DTO
	var articleDTO ArticleResponse
	articleDTO.ToDTO(edited)

	ctx.JSON(http.StatusOK, GenerateContentResponse{Article: articleDTO})
}