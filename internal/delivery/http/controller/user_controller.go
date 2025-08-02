package controller

import (
	"context"
	"net/http"
	"write_base/internal/domain"

	"github.com/gin-gonic/gin"
)


type UserController struct{
	userUsercase domain.IUserUsecase
}

func NewUserController(ctx context.Context, usecase domain.IUserUsecase) * UserController{
	return &UserController{userUsercase: usecase}
}

func (uc *UserController) Register(c *gin.Context ){
	ctx := c.Request.Context()
	var user RegisterRequest
	if err:= c.ShouldBindJSON(&user); err != nil{
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := uc.userUsercase.Register(ctx, user.ToRegisterInput())
	if err != nil{
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message":err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "successfully registerd"})
}

func (uc *UserController)Login(c *gin.Context){
	ctx := c.Request.Context()

	var user LoginRequest
	if err:= c.ShouldBindJSON(&user); err != nil{
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	iP := c.ClientIP()
	userAgent := c.Request.UserAgent()
	deviceInfo := c.GetHeader("Device-Info")
	metadata := &domain.AuthMetadata{IP: iP, UserAgent: userAgent, DeviceInfo: deviceInfo}

	jwtToken, err := uc.userUsercase.Login(ctx, user.ToLoginInput(), metadata)
	if err != nil{
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message":err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message":"user logged in successfully","token":jwtToken})

}


func (uc *UserController) Logout (c *gin.Context) {
    ctx := c.Request.Context()
    var req LogoutRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }
    err := uc.userUsercase.Logout(ctx, req.RefreshToken)
    if err != nil {
        c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
        return
    }
    c.IndentedJSON(http.StatusOK, gin.H{"message": "User logged out successfully"})
}
