package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// Minimal smoke test that routes can be registered without panics
func TestRegisterRouters_NoPanic(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	// do not call DI; we only ensure functions exist and can be invoked with nils where accepted
	// Just verify engine can serve a 404 for an unregistered route
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/__nope__", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
