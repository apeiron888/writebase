package infrastructure

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"write_base/internal/domain"

	"github.com/gin-gonic/gin"
)

func TestRequireRole_Allowed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("role", string(domain.RoleAdmin)); c.Next() })
	r.GET("/p", RequireRole(domain.RoleAdmin, domain.RoleSuperAdmin), func(c *gin.Context) { c.Status(200) })
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/p", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("want 200 got %d", w.Code)
	}
}

func TestRequireRole_Forbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("role", string(domain.RoleUser)); c.Next() })
	r.GET("/p", RequireRole(domain.RoleAdmin), func(c *gin.Context) { c.Status(200) })
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/p", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("want 403 got %d", w.Code)
	}
}
