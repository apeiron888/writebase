package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"write_base/internal/delivery/http/controller"
	"write_base/internal/domain"
	"write_base/internal/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestArticleController_Create_BadPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uc := &mocks.ArticleUsecaseMock{}
	uc.CreateArticleFn = func(_ context.Context, _ string, _ *domain.Article) (string, error) {
		return "", domain.ErrInvalidArticlePayload
	}
	h := controller.NewArticleHandler(uc)

	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("user_id", "u1"); c.Next() })
	r.POST("/articles/new", h.CreateArticle)

	// missing required fields
	payload := map[string]any{
		"title":          "",
		"content_blocks": []map[string]any{},
	}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/articles/new", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}
