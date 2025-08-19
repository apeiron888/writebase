package infrastructure

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"write_base/internal/domain"

	"github.com/gin-gonic/gin"
)

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	m := NewMiddleware(&fakeTokenService{validate: func(string) (*domain.AuthClaims, error) { return nil, domain.ErrInvalidToken }})
	r := gin.New()
	r.Use(m.Authmiddleware())
	r.GET("/x", func(c *gin.Context) { c.Status(200) })
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Authorization", "Bearer bad")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("want 401 got %d", w.Code)
	}
}
