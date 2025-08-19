package infrastructure

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"write_base/internal/domain"

	"github.com/gin-gonic/gin"
)

type fakeTokenService struct {
	validate func(string) (*domain.AuthClaims, error)
}

func (f *fakeTokenService) GenerateAccessToken(u *domain.User) (string, error)  { return "", nil }
func (f *fakeTokenService) GenerateRefreshToken(u *domain.User) (string, error) { return "", nil }
func (f *fakeTokenService) ValidateAccessToken(t string) (*domain.AuthClaims, error) {
	return f.validate(t)
}
func (f *fakeTokenService) ValidateRefreshToken(t string) (*domain.AuthClaims, error) {
	return f.validate(t)
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	m := NewMiddleware(&fakeTokenService{})
	r := gin.New()
	r.Use(m.Authmiddleware())
	r.GET("/x", func(c *gin.Context) { c.Status(200) })
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/x", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("want 401 got %d", w.Code)
	}
}

func TestAuthMiddleware_BadFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	m := NewMiddleware(&fakeTokenService{})
	r := gin.New()
	r.Use(m.Authmiddleware())
	r.GET("/x", func(c *gin.Context) { c.Status(200) })
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Authorization", "Token abc")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("want 401 got %d", w.Code)
	}
}

func TestAuthMiddleware_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	m := NewMiddleware(&fakeTokenService{validate: func(string) (*domain.AuthClaims, error) {
		return &domain.AuthClaims{UserID: "u1", Role: string(domain.RoleAdmin)}, nil
	}})
	r := gin.New()
	r.Use(m.Authmiddleware())
	r.GET("/x", func(c *gin.Context) {
		if c.GetString("user_id") == "u1" {
			c.Status(200)
		} else {
			c.Status(500)
		}
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Authorization", "Bearer token")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("want 200 got %d", w.Code)
	}
}
