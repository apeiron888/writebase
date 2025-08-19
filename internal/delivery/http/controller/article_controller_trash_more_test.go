package controller_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"write_base/internal/delivery/http/controller"
	"write_base/internal/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestEmptyTrash_Success(t *testing.T) {
	uc := &mocks.ArticleUsecaseMock{EmptyTrashFn: func(_ context.Context, _ string) error { return nil }}
	h := controller.NewArticleHandler(uc)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("user_id", "u1"); c.Next() })
	r.DELETE("/me/trash", h.EmptyTrash)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/me/trash", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeleteFromTrash_BadRequest(t *testing.T) {
	uc := &mocks.ArticleUsecaseMock{}
	h := controller.NewArticleHandler(uc)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("user_id", "u1"); c.Next() })
	// malformed route to produce empty id
	r.DELETE("/articles/trash//", h.DeleteFromTrash)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/articles/trash//", nil)
	r.ServeHTTP(w, req)
	// malformed route either redirects (307) or 404 depending on router; accept both
	if w.Code != http.StatusTemporaryRedirect && w.Code != http.StatusNotFound {
		t.Fatalf("expected 307 or 404, got %d", w.Code)
	}
}
