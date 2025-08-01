package controller

import (
	"net/http"
	"starter/internal/domain"
	dtodlv "starter/internal/delivery/http/controller/dto"
	"github.com/gin-gonic/gin"
)

type ReactionController struct {
	usecase domain.IReactionUsecase
}

func NewReactionController(usecase domain.IReactionUsecase) *ReactionController {
	return &ReactionController{usecase: usecase}
}

func (rc *ReactionController) AddReaction(c *gin.Context) {
	   var req dtodlv.ReactionRequest
	   if err := c.ShouldBindJSON(&req); err != nil {
			   c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			   return
	   }
	   reaction := &domain.Reaction{
			   PostID:    req.PostID,
			   UserID:    req.UserID,
			   CommentID: req.CommentID,
			   Type:      domain.ReactionType(req.Type),
	   }
	   if err := rc.usecase.AddReaction(c, reaction); err != nil {
			   c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			   return
	   }
	   c.JSON(http.StatusCreated, gin.H{"message": "Reaction added"})
}

func (rc *ReactionController) RemoveReaction(c *gin.Context) {
	   id := c.Param("id")
	   if err := rc.usecase.RemoveReaction(c, id); err != nil {
			   c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			   return
	   }
	   c.JSON(http.StatusOK, gin.H{"message": "Reaction removed"})
}

func (rc *ReactionController) GetReactionsByPost(c *gin.Context) {
	   postID := c.Param("post_id")
	   reactions, err := rc.usecase.GetReactionsByPost(c, postID)
	   if err != nil {
			   c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			   return
	   }
	   c.JSON(http.StatusOK, reactions)
}

func (rc *ReactionController) GetReactionsByUser(c *gin.Context) {
	   userID := c.Param("user_id")
	   reactions, err := rc.usecase.GetReactionsByUser(c, userID)
	   if err != nil {
			   c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			   return
	   }
	   c.JSON(http.StatusOK, reactions)
}

func (rc *ReactionController) CountReactions(c *gin.Context) {
	   postID := c.Param("post_id")
	   reactionType := c.Param("type")
	   count, err := rc.usecase.CountReactions(c, postID, domain.ReactionType(reactionType))
	   if err != nil {
			   c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			   return
	   }
	   c.JSON(http.StatusOK, gin.H{"count": count})
}
