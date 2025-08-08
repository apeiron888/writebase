package router
import (
	"write_base/internal/delivery/http/controller"
	"write_base/internal/domain"
	"write_base/internal/infrastructure"

	"github.com/gin-gonic/gin"
)

// Comment Routes
func RegisterCommentRoutes(r *gin.Engine, commentController *controller.CommentController) {
    comments := r.Group("/comments")
    {
        comments.POST("", commentController.Create)
        comments.PUT("/:id", commentController.Update)
        comments.DELETE("/:id", commentController.Delete)
        comments.GET("/:id", commentController.GetByID)
        comments.GET("/post/:post_id", commentController.GetByPostID)
        comments.GET("/user/:user_id", commentController.GetByUserID)
        comments.GET("/replies/:parent_id", commentController.GetReplies)
    }
}

// Reaction Routes
func RegisterReactionRoutes(r *gin.Engine, reactionController *controller.ReactionController) {
    reactions := r.Group("/reactions")
    {
        reactions.POST("", reactionController.AddReaction)
        reactions.DELETE(":id", reactionController.RemoveReaction)
        reactions.GET("/post/:post_id", reactionController.GetReactionsByPost)
        reactions.GET("/user/:user_id", reactionController.GetReactionsByUser)
        reactions.GET("/count/:post_id/:type", reactionController.CountReactions)
    }
}

// Follow Routes
func RegisterFollowRoutes(r *gin.Engine, followController *controller.FollowController) {
    follows := r.Group("/follows")
    {
        follows.POST("/follow", followController.FollowUser)
        follows.POST("/unfollow", followController.UnfollowUser)
        follows.GET("/followers/:user_id", followController.GetFollowers)
        follows.GET("/following/:user_id", followController.GetFollowing)
        follows.GET("/is-following/:follower_id/:followee_id", followController.IsFollowing)
    }
}

// Report Routes
func RegisterReportRoutes(r *gin.Engine, reportController *controller.ReportController) {
    reports := r.Group("/reports")
    {
        reports.POST("", reportController.CreateReport)
        reports.GET("", reportController.GetReports)
        reports.PUT(":id/status", reportController.UpdateReportStatus)
    }
}

// AI Routes
func RegisterAIRoutes(r *gin.Engine, aiController *controller.AIController) {
    ai := r.Group("/ai")
    {
        ai.POST("/suggest", aiController.Suggest)
        ai.POST("/generate-content", aiController.GenerateContent)
    }
}



func UserRouter(r *gin.Engine, userController *controller.UserController, authMiddleware *infrastructure.Middleware){
	auth := r.Group("/auth")
	{
		auth.POST("/register",userController.Register)
		auth.GET("/verify", userController.Verify)
		auth.GET("/verify-Update-Email", userController.VerifyUpdateEmail)
		auth.POST("/login", userController.Login)
		auth.GET("/google/login", userController.GoogleLogin)
		auth.GET("/google/callback", userController.GoogleCallback)
		auth.POST("/forget-password", userController.ForgetPassword)
		auth.POST("/reset-password", userController.ResetPassword)
		auth.POST("/logout", authMiddleware.Authmiddleware(), userController.Logout)
		auth.POST("/refresh",userController.RefreshToken)
	}
	user := r.Group("/users")
	user.Use(authMiddleware.Authmiddleware())
	{
		user.GET("/me", userController.MyProfile)
		user.PATCH("/me", userController.UpdateMyProfile)
		user.PATCH("/username", userController.UpdateMyUsername)
		user.PATCH("/email", userController.UpdateMyEmail)
		user.PUT("/password",userController.ChangeMyPassword)
		
	}
	admin := r.Group("/admin")
	admin.Use(authMiddleware.Authmiddleware(), infrastructure.RequireRole(domain.RoleAdmin, domain.RoleSuperAdmin))
	{
		admin.PUT("/user/:id/promote", userController.PromoteToAdmin)
		admin.PUT("/user/:id/demote", userController.DemoteToUser)
		admin.PUT("/user/:id/disable", userController.DisableUser)
		admin.PUT("/user/:id/enable", userController.EnableUser)
	}

}