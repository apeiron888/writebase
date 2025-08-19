package controller_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"write_base/internal/delivery/http/controller"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestGetArticleBySlug_BadSlug(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := controller.NewArticleHandler(nil)
	r := gin.New()
	r.GET("/:slug", h.GetArticleBySlug)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/has space", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}
