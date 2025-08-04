package controller

import (
	"net/http"
	dtodlv "write_base/internal/delivery/http/controller/dto"
	"write_base/internal/domain"

	"github.com/gin-gonic/gin"
)

type CommentController struct {
	usecase domain.ICommentUsecase
}

func NewCommentController(usecase domain.ICommentUsecase) *CommentController {
	return &CommentController{usecase: usecase}
}

func (cc *CommentController) Create(c *gin.Context) {
	   var req dtodlv.CommentRequest
	   if err := c.ShouldBindJSON(&req); err != nil {
			   c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			   return
	   }
	   comment := &domain.Comment{
			   PostID:   req.PostID,
			   UserID:   req.UserID,
			   ParentID: req.ParentID,
			   Content:  req.Content,
	   }
	   if err := cc.usecase.CreateComment(c, comment); err != nil {
			   c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			   return
	   }
	   c.JSON(http.StatusCreated, gin.H{"message": "Comment created"})
}

func (cc *CommentController) Update(c *gin.Context) {
	   var req dtodlv.UpdateCommentRequest
	   id := c.Param("id")
	   if err := c.ShouldBindJSON(&req); err != nil {
			   c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			   return
	   }
	   comment := &domain.Comment{
			   ID:      id,
			   Content: req.Content,
	   }
	   if err := cc.usecase.UpdateComment(c, comment); err != nil {
			   c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			   return
	   }
	   c.JSON(http.StatusOK, gin.H{"message": "Comment updated"})
}

func (cc *CommentController) Delete(c *gin.Context) {
	   id := c.Param("id")
	   if err := cc.usecase.DeleteComment(c, id); err != nil {
			   c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			   return
	   }
	   c.JSON(http.StatusOK, gin.H{"message": "Comment deleted"})
}

func (cc *CommentController) GetByID(c *gin.Context) {
	   id := c.Param("id")
	   comment, err := cc.usecase.GetCommentByID(c, id)
	   if err != nil {
			   c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			   return
	   }
	   c.JSON(http.StatusOK, comment)
}

func (cc *CommentController) GetByPostID(c *gin.Context) {
	   postID := c.Param("post_id")
	   comments, err := cc.usecase.GetCommentsByPostID(c, postID)
	   if err != nil {
			   c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			   return
	   }
	   c.JSON(http.StatusOK, comments)
}

func (cc *CommentController) GetByUserID(c *gin.Context) {
	   userID := c.Param("user_id")
	   comments, err := cc.usecase.GetCommentsByUserID(c, userID)
	   if err != nil {
			   c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			   return
	   }
	   c.JSON(http.StatusOK, comments)
}

func (cc *CommentController) GetReplies(c *gin.Context) {
	   parentID := c.Param("parent_id")
	   replies, err := cc.usecase.GetReplies(c, parentID)
	   if err != nil {
			   c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			   return
	   }
	   c.JSON(http.StatusOK, replies)
}
