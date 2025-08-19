package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	dtodlv "write_base/internal/delivery/http/controller/dto"
	"write_base/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// fake usecase implementing domain.IReactionUsecase
type fakeReactionUC struct {
	err  error
	list []*domain.Reaction
}

func (f *fakeReactionUC) AddReaction(_ context.Context, _ *domain.Reaction) error { return f.err }
func (f *fakeReactionUC) RemoveReaction(_ context.Context, _ string) error        { return f.err }
func (f *fakeReactionUC) GetReactionsByPost(_ context.Context, _ string) ([]*domain.Reaction, error) {
	return f.list, f.err
}
func (f *fakeReactionUC) GetReactionsByUser(_ context.Context, _ string) ([]*domain.Reaction, error) {
	return f.list, f.err
}
func (f *fakeReactionUC) CountReactions(_ context.Context, _ string, _ domain.ReactionType) (int, error) {
	return 1, f.err
}

func setupReactionRouter(uc domain.IReactionUsecase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewReactionController(uc)
	r.POST("/reactions", h.AddReaction)
	r.DELETE("/reactions/:id", h.RemoveReaction)
	r.GET("/posts/:post_id/reactions", h.GetReactionsByPost)
	r.GET("/users/:user_id/reactions", h.GetReactionsByUser)
	r.GET("/posts/:post_id/reactions/:type/count", h.CountReactions)
	return r
}

func TestReactionController_HappyPaths(t *testing.T) {
	uc := &fakeReactionUC{list: []*domain.Reaction{{ID: "1"}}}
	r := setupReactionRouter(uc)

	// Add
	body, _ := json.Marshal(dtodlv.ReactionRequest{PostID: "p1", UserID: "u1", Type: "like"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/reactions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	// Remove
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodDelete, "/reactions/1", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	// ByPost
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/posts/p1/reactions", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	// ByUser
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/users/u1/reactions", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	// Count
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/posts/p1/reactions/like/count", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
