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

func TestArticleController_Create_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	uc := &mocks.ArticleUsecaseMock{}
	uc.CreateArticleFn = func(_ context.Context, userID string, input *domain.Article) (string, error) {
		if userID == "" || input == nil || input.Title == "" {
			t.Fatalf("bad input")
		}
		return "id-1", nil
	}

	h := controller.NewArticleHandler(uc)

	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("user_id", "u1"); c.Next() })
	r.POST("/articles/new", h.CreateArticle)

	payload := map[string]any{
		"title":    "Hello",
		"slug":     "hello",
		"excerpt":  "Hello world",
		"language": "en",
		"tags":     []string{"go"},
		"content_blocks": []map[string]any{{
			"type":    "paragraph",
			"order":   0,
			"content": map[string]any{"paragraph": map[string]any{"text": "hi", "style": "normal"}},
		}},
	}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/articles/new", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)
}

func TestArticleController_Create_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	uc := &mocks.ArticleUsecaseMock{}
	h := controller.NewArticleHandler(uc)

	r := gin.New()
	r.POST("/articles/new", h.CreateArticle)

	payload := map[string]any{
		"title":    "Hello",
		"language": "en",
		"content_blocks": []map[string]any{{
			"type":    "paragraph",
			"order":   0,
			"content": map[string]any{"paragraph": map[string]any{"text": "hi"}},
		}},
	}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/articles/new", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}
