package router

import (
	"write_base/internal/delivery/http/controller"
	"write_base/internal/domain"
	"write_base/internal/infrastructure"

	"github.com/gin-gonic/gin"
)

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