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

// fake usecase implementing domain.ICommentUsecase
type fakeCommentUC struct {
	createErr  error
	updateErr  error
	deleteErr  error
	getByIDRes *domain.Comment
	getByIDErr error
	listRes    []*domain.Comment
	listErr    error
}

func (f *fakeCommentUC) CreateComment(_ context.Context, _ *domain.Comment) error { return f.createErr }
func (f *fakeCommentUC) UpdateComment(_ context.Context, _ *domain.Comment) error { return f.updateErr }
func (f *fakeCommentUC) DeleteComment(_ context.Context, _ string) error          { return f.deleteErr }
func (f *fakeCommentUC) GetCommentByID(_ context.Context, _ string) (*domain.Comment, error) {
	return f.getByIDRes, f.getByIDErr
}
func (f *fakeCommentUC) GetCommentsByPostID(_ context.Context, _ string) ([]*domain.Comment, error) {
	return f.listRes, f.listErr
}
func (f *fakeCommentUC) GetCommentsByUserID(_ context.Context, _ string) ([]*domain.Comment, error) {
	return f.listRes, f.listErr
}
func (f *fakeCommentUC) GetReplies(_ context.Context, _ string) ([]*domain.Comment, error) {
	return f.listRes, f.listErr
}

func setupCommentRouter(uc domain.ICommentUsecase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewCommentController(uc)
	r.POST("/comments", h.Create)
	r.PUT("/comments/:id", h.Update)
	r.DELETE("/comments/:id", h.Delete)
	r.GET("/comments/:id", h.GetByID)
	r.GET("/posts/:post_id/comments", h.GetByPostID)
	r.GET("/users/:user_id/comments", h.GetByUserID)
	// Use a different prefix to avoid wildcard conflict with /comments/:id
	r.GET("/cmt/:parent_id/replies", h.GetReplies)
	return r
}

func TestCommentController_HappyPaths(t *testing.T) {
	uc := &fakeCommentUC{getByIDRes: &domain.Comment{ID: "1"}, listRes: []*domain.Comment{{ID: "1"}}}
	r := setupCommentRouter(uc)

	// Create
	body, _ := json.Marshal(dtodlv.CommentRequest{PostID: "p1", UserID: "u1", Content: "hi"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/comments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	// Update
	up, _ := json.Marshal(dtodlv.UpdateCommentRequest{Content: "new"})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPut, "/comments/1", bytes.NewReader(up))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	// Delete
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodDelete, "/comments/1", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	// GetByID
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/comments/1", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	// GetByPostID
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/posts/p1/comments", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	// GetByUserID
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/users/u1/comments", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	// GetReplies
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/cmt/parent1/replies", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
