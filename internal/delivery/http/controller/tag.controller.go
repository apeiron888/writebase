package controller

import (
	"net/http"
	"write_base/internal/domain"

	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	tagUsecase domain.TagUsecase
}

func NewTagHandler(tagUsecase domain.TagUsecase) *TagHandler {
	return &TagHandler{tagUsecase: tagUsecase}
}

// CreateTag creates a new tag
func (h *TagHandler) CreateTag(c *gin.Context) {
	// userID := c.GetString("user_id")
	// if userID == "" {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	// 	return
	// }
	userID := "1234"

	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag, err := h.tagUsecase.CreateTag(c.Request.Context(), userID, req.Name)
	if err != nil {
		code := http.StatusInternalServerError
		switch err {
		case domain.ErrTagAlreadyExists, domain.ErrInvalidTagName:
			code = http.StatusBadRequest
		}
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": tag})
}

// ApproveTag approves a pending tag (admin only)
func (h *TagHandler) ApproveTag(c *gin.Context) {
	tagID := c.Param("id")
	if tagID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tag ID required"})
		return
	}

	// roleVal, exists := ctx.Get("user_role")
	// if !exists {
	// 	ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_role not found in context"})
	// 	return
	// }
	// userRole := roleVal.(string)
    userRole := "admin"

	if userRole != "admin" && userRole!="super_admin"{
		c.JSON(http.StatusUnauthorized, gin.H{"error": domain.ErrUnauthorized})
		return 
	} 

	tag, err := h.tagUsecase.ApproveTag(c.Request.Context(), tagID)
	if err != nil {
		code := http.StatusInternalServerError
		if err == domain.ErrTagNotFound {
			code = http.StatusNotFound
		}
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tag})
}

// ListTags lists tags with optional status filter
func (h *TagHandler) ListTags(c *gin.Context) {
	status := domain.TagStatus(c.Query("status"))
	tags, err := h.tagUsecase.ListTags(c.Request.Context(), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tags})
}

// DeleteTag deletes a tag (admin only)
func (h *TagHandler) DeleteTag(c *gin.Context) {
	tagID := c.Param("id")
	if tagID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tag ID required"})
		return
	}

	// roleVal, exists := ctx.Get("user_role")
	// if !exists {
	// 	ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_role not found in context"})
	// 	return
	// }
	// userRole := roleVal.(string)
    userRole := "admin"

	if userRole != "admin" && userRole!="super_admin"{
		c.JSON(http.StatusUnauthorized, gin.H{"error": domain.ErrUnauthorized})
		return 
	} 

	if err := h.tagUsecase.DeleteTag(c.Request.Context(), tagID); err != nil {
		code := http.StatusInternalServerError
		if err == domain.ErrTagNotFound {
			code = http.StatusNotFound
		}
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// RejectTag rejects a pending tag (admin only)
func (h *TagHandler) RejectTag(c *gin.Context) {
	tagID := c.Param("id")
	if tagID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tag ID required"})
		return
	}

	// roleVal, exists := ctx.Get("user_role")
	// if !exists {
	// 	ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_role not found in context"})
	// 	return
	// }
	// userRole := roleVal.(string)
    userRole := "admin"

	if userRole != "admin" && userRole!="super_admin"{
		c.JSON(http.StatusUnauthorized, gin.H{"error": domain.ErrUnauthorized})
		return 
	} 


	tag, err := h.tagUsecase.RejectTag(c.Request.Context(), tagID)
	if err != nil {
		code := http.StatusInternalServerError
		if err == domain.ErrTagNotFound {
			code = http.StatusNotFound
		}
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tag})
}