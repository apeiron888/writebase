package controller

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
	"write_base/config"
	"write_base/internal/domain"

	// "golang.org/x/oauth2"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)


type UserController struct{
	userUsercase domain.IUserUsecase
}

func NewUserController(usecase domain.IUserUsecase) * UserController{
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
func (uc *UserController) Verify(c *gin.Context){
	ctx := c.Request.Context()
	code := c.Query("code")
	if code == ""{
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": domain.ErrMissingVerifyCode.Error()})
		return
	}
	err := uc.userUsercase.VerifyEmail(ctx, code)
	if err != nil{
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return 
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message":"Account is verified"})

}

func (uc *UserController) Login(c *gin.Context){
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

	c.IndentedJSON(http.StatusOK, gin.H{"message":"user logged in successfully","token":ToLoginResponse(jwtToken)})

}

////-========================google auth====================================////
func (uc *UserController) GoogleLogin(c *gin.Context) {
	stateToken := uuid.New().String()
	http.SetCookie(c.Writer, &http.Cookie{
		Name: "oauthStateToken",
		Value: stateToken,
		Expires: time.Now().Add((10 * time.Minute)),
		HttpOnly: true,
		Secure: false, // make it true for https
		Path: "/auth",
		
	})
	url := config.GoogleOAuthConfig.AuthCodeURL(stateToken) // 
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (uc *UserController) GoogleCallback(c *gin.Context) {
	ctx := c.Request.Context()
	stateFromQurey := c.Query("state")
	if stateFromQurey == ""{
		c.IndentedJSON(http.StatusBadRequest, domain.ErrMissingState)
	}
	cookie, err := c.Request.Cookie("oauthStateToken")
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, domain.ErrMissingOAuthStateToken)
	}
	stateFromcookie := cookie.Value
	if stateFromQurey != stateFromcookie{
		c.IndentedJSON(http.StatusBadRequest, domain.ErrMissingOrExpiredStateCookie)
	}
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": domain.ErrMissingOAuthCode.Error()})
		return

	}

	token, err := config.GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": domain.ErrTokenExchangeFailed.Error()})
		return
	}

	client := config.GoogleOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": domain.ErrFailedToFetchUserInfo.Error()})
		return
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(body, &data)

	email := data["email"].(string)
	name := data["name"].(string)

	// Pass to Usecase
	registerOrLogin := &domain.RegisterInput{Username: name, Email:email}
	iP := c.ClientIP()
	userAgent := c.Request.UserAgent()
	deviceInfo := c.GetHeader("Device-Info")
	metadata := &domain.AuthMetadata{IP: iP, UserAgent: userAgent, DeviceInfo: deviceInfo}
	jwtToken, err := uc.userUsercase.LoginOrRegisterOAuthUser(ctx, registerOrLogin, metadata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": domain.ErrOAuthLoginFailed.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"access_token": ToLoginResponse(jwtToken)})
}
/////====================================================================================////
func (uc *UserController) RefreshToken(c *gin.Context){
	ctx := c.Request.Context()
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil{
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message":err.Error()})
		return
	}
	loginResult, err := uc.userUsercase.RefreshToken(ctx, req.RefreshToken)
	if err!= nil{
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
        return

	}
	c.IndentedJSON(http.StatusOK, gin.H{"token": ToRefreshTokenResponse(loginResult)})

}

func (uc *UserController) Logout(c *gin.Context) {
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
func (uc *UserController) ForgetPassword(c *gin.Context){
	ctx := c.Request.Context()
	var req ForgotPasswordRequest
	if err:= c.ShouldBindJSON(&req); err != nil{
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message":err.Error()})
		return
	}
	err := uc.userUsercase.ForgotPassword(ctx, req.Email)
	if err != nil{
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "check your email"})


}
func (uc *UserController) ResetPassword(c *gin.Context){
	ctx := c.Request.Context()
	var req ResetPasswordRequest
	if err:= c.ShouldBindJSON(&req); err!= nil{
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message":err.Error()})
		return
	}
	err := uc.userUsercase.ResetPassword(ctx, req.Token, req.NewPassword)
	if err != nil{
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message":err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Reset Password successfully"})
}
/////.......... authenticated user.........//

func (uc *UserController) MyProfile(c *gin.Context){
	
}
func (uc *UserController) UpdateMyProfile(c *gin.Context){
	
}
func (uc *UserController) ChangeMyPassword(c *gin.Context){
	
}

////////////............Admin........////////////
func (uc *UserController) DemoteToUser(c *gin.Context){

	
}
func (uc *UserController) PromoteToAdmin(c *gin.Context){
	
}
func (uc *UserController) DisabelUser(c *gin.Context){
	
}
func (uc *UserController) EnableUser(c *gin.Context){
	
}