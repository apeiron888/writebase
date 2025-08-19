package controller

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"write_base/internal/domain"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type fakeUserUC struct{}

func (f *fakeUserUC) Register(ctx context.Context, input *domain.RegisterInput) error { return nil }
func (f *fakeUserUC) VerifyEmail(ctx context.Context, token string) error             { return nil }
func (f *fakeUserUC) VerifyUpdateEmail(ctx context.Context, token string) error       { return nil }
func (f *fakeUserUC) Login(ctx context.Context, input *domain.LoginInput, metadata *domain.AuthMetadata) (*domain.LoginResult, error) {
	return &domain.LoginResult{AccessToken: "a", RefreshToken: "r"}, nil
}
func (f *fakeUserUC) Logout(ctx context.Context, refreshToken string) error { return nil }
func (f *fakeUserUC) RefreshToken(ctx context.Context, refreshToken string) (*domain.LoginResult, error) {
	return &domain.LoginResult{AccessToken: "a2", RefreshToken: "r2"}, nil
}
func (f *fakeUserUC) LoginOrRegisterOAuthUser(ctx context.Context, input *domain.RegisterInput, metadata *domain.AuthMetadata) (*domain.LoginResult, error) {
	return &domain.LoginResult{AccessToken: "a3", RefreshToken: "r3"}, nil
}
func (f *fakeUserUC) GetProfile(ctx context.Context, userID string) (*domain.User, error) {
	return &domain.User{ID: userID}, nil
}
func (f *fakeUserUC) UpdateProfile(ctx context.Context, input *domain.UpdateProfileInput) error {
	return nil
}
func (f *fakeUserUC) UpdateUsername(ctx context.Context, input *domain.UpdateAccountInput) error {
	return nil
}
func (f *fakeUserUC) UpdateEmail(ctx context.Context, input *domain.UpdateAccountInput) error {
	return nil
}
func (f *fakeUserUC) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	return nil
}
func (f *fakeUserUC) ForgotPassword(ctx context.Context, email string) error             { return nil }
func (f *fakeUserUC) ResetPassword(ctx context.Context, token, newPassword string) error { return nil }
func (f *fakeUserUC) DemoteToUser(ctx context.Context, userID string) error              { return nil }
func (f *fakeUserUC) PromoteToAdmin(ctx context.Context, userID string) error            { return nil }
func (f *fakeUserUC) DisableUser(ctx context.Context, userID string) error               { return nil }
func (f *fakeUserUC) EnableUser(ctx context.Context, userID string) error                { return nil }

func TestUserController_HappyPaths(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uc := &fakeUserUC{}
	google := &oauth2.Config{ClientID: "x", RedirectURL: "http://localhost/cb"}
	c := NewUserController(uc, google)

	// Register
	r := gin.New()
	r.POST("/register", c.Register)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"username":"neo","email":"e@example.com","password":"P@ssw0rd!"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("register %d", w.Code)
	}

	// Verify
	r = gin.New()
	r.GET("/verify", c.Verify)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/verify?code=abc", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("verify %d", w.Code)
	}

	// Verify update email
	r = gin.New()
	r.GET("/verify-up", c.VerifyUpdateEmail)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/verify-up?code=abc", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("verify up %d", w.Code)
	}

	// Login
	r = gin.New()
	r.POST("/login", c.Login)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"email_or_username":"e@x","password":"p"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("login %d", w.Code)
	}

	// Refresh
	r = gin.New()
	r.POST("/refresh", c.RefreshToken)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/refresh", strings.NewReader(`{"refresh_token":"r"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("refresh %d", w.Code)
	}

	// Logout
	r = gin.New()
	r.POST("/logout", c.Logout)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/logout", strings.NewReader(`{"refresh_token":"r"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("logout %d", w.Code)
	}

	// Forgot password
	r = gin.New()
	r.POST("/forgot", c.ForgetPassword)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/forgot", strings.NewReader(`{"email":"e@example.com"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("forgot %d", w.Code)
	}

	// Reset password
	r = gin.New()
	r.POST("/reset", c.ResetPassword)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/reset", strings.NewReader(`{"token":"t","new_password":"NewP@ssw0rd!"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("reset %d", w.Code)
	}

	// My profile
	r = gin.New()
	r.Use(func(c *gin.Context) { c.Set("user_id", "u1"); c.Next() })
	r.GET("/me", c.MyProfile)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/me", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("me %d", w.Code)
	}

	// Update profile
	r = gin.New()
	r.Use(func(c *gin.Context) { c.Set("user_id", "u1"); c.Next() })
	r.PATCH("/me", c.UpdateMyProfile)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPatch, "/me", strings.NewReader(`{"bio":"hello"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusAccepted {
		t.Fatalf("update profile %d", w.Code)
	}

	// Change password
	r = gin.New()
	r.Use(func(c *gin.Context) { c.Set("user_id", "u1"); c.Next() })
	r.PUT("/password", c.ChangeMyPassword)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPut, "/password", strings.NewReader(`{"old_password":"oldsecret","new_password":"NewP@ssw0rd!"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code < 200 || w.Code >= 300 {
		t.Fatalf("change pwd %d", w.Code)
	}

	// Update username
	r = gin.New()
	r.Use(func(c *gin.Context) { c.Set("user_id", "u1"); c.Next() })
	r.PATCH("/username", c.UpdateMyUsername)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPatch, "/username", strings.NewReader(`{"username":"neo"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code < 200 || w.Code >= 300 {
		t.Fatalf("update username %d", w.Code)
	}

	// Update email
	r = gin.New()
	r.Use(func(c *gin.Context) { c.Set("user_id", "u1"); c.Next() })
	r.PATCH("/email", c.UpdateMyEmail)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPatch, "/email", strings.NewReader(`{"email":"e2@example.com"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code < 200 || w.Code >= 300 {
		t.Fatalf("update email %d", w.Code)
	}

	// Google login (redirect)
	r = gin.New()
	r.GET("/google/login", c.GoogleLogin)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/google/login", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusTemporaryRedirect {
		t.Fatalf("google login %d", w.Code)
	}
}
