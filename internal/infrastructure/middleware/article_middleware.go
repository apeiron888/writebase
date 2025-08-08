package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Article ID Validation middleware
func ValidateArticleID() gin.HandlerFunc {
    return func(c *gin.Context) {
        id := c.Param("id")
        if !isValidUUID(id) {
            c.AbortWithStatusJSON(400, gin.H{"error": "Invalid article ID"})
            return
        }
        c.Set("articleID", id)
        c.Next()
    }
}