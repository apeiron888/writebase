package controller

import (
	"net/http"
	dtodlv "write_base/internal/delivery/http/controller/dto"
	"write_base/internal/domain"

	"github.com/gin-gonic/gin"
)

type FollowController struct {
	usecase domain.IFollowUsecase
}

func NewFollowController(usecase domain.IFollowUsecase) *FollowController {
	return &FollowController{usecase: usecase}
}

func (fc *FollowController) FollowUser(c *gin.Context) {
	   var req dtodlv.FollowRequest
	   if err := c.ShouldBindJSON(&req); err != nil {
			   c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			   return
	   }
	   if err := fc.usecase.FollowUser(c, req.FollowerID, req.FolloweeID); err != nil {
			   c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			   return
	   }
	   c.JSON(http.StatusCreated, gin.H{"message": "Followed user"})
}

func (fc *FollowController) UnfollowUser(c *gin.Context) {
	   var req dtodlv.FollowRequest
	   if err := c.ShouldBindJSON(&req); err != nil {
			   c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			   return
	   }
	   if err := fc.usecase.UnfollowUser(c, req.FollowerID, req.FolloweeID); err != nil {
			   c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			   return
	   }
	   c.JSON(http.StatusOK, gin.H{"message": "Unfollowed user"})
}

func (fc *FollowController) GetFollowers(c *gin.Context) {
	   userID := c.Param("user_id")
	   followers, err := fc.usecase.GetFollowers(c, userID)
	   if err != nil {
			   c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			   return
	   }
	   c.JSON(http.StatusOK, followers)
}

func (fc *FollowController) GetFollowing(c *gin.Context) {
	   userID := c.Param("user_id")
	   following, err := fc.usecase.GetFollowing(c, userID)
	   if err != nil {
			   c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			   return
	   }
	   c.JSON(http.StatusOK, following)
}

func (fc *FollowController) IsFollowing(c *gin.Context) {
	   followerID := c.Param("follower_id")
	   followeeID := c.Param("followee_id")
	   isFollowing, err := fc.usecase.IsFollowing(c, followerID, followeeID)
	   if err != nil {
			   c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			   return
	   }
	   c.JSON(http.StatusOK, gin.H{"is_following": isFollowing})
}
