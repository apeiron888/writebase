package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestVerify_MissingCode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uc := NewUserController(nil, nil)
	r := gin.New()
	r.GET("/verify", uc.Verify)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/verify", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400 got %d", w.Code)
	}
}

func TestGoogleCallback_MissingState(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uc := NewUserController(nil, nil)
	r := gin.New()
	r.GET("/cb", uc.GoogleCallback)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/cb", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400 got %d", w.Code)
	}
}

func TestMyProfile_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uc := NewUserController(nil, nil)
	r := gin.New()
	r.GET("/me", uc.MyProfile)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/me", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("want 401 got %d", w.Code)
	}
}

func TestChangeMyPassword_BadPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uc := NewUserController(nil, nil)
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("user_id", "u1"); c.Next() })
	r.PUT("/pwd", uc.ChangeMyPassword)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/pwd", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400 got %d", w.Code)
	}
}

func TestUpdateMyProfile_BadPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uc := NewUserController(nil, nil)
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("user_id", "u1"); c.Next() })
	r.PATCH("/me", uc.UpdateMyProfile)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/me", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400 got %d", w.Code)
	}
}
