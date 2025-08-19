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

// fake usecase implementing domain.IFollowUsecase
type fakeFollowUC struct{ err error }

func (f *fakeFollowUC) FollowUser(_ context.Context, _ string, _ string) error   { return f.err }
func (f *fakeFollowUC) UnfollowUser(_ context.Context, _ string, _ string) error { return f.err }
func (f *fakeFollowUC) GetFollowers(_ context.Context, _ string) ([]*domain.User, error) {
	return []*domain.User{{ID: "u1"}}, f.err
}
func (f *fakeFollowUC) GetFollowing(_ context.Context, _ string) ([]*domain.User, error) {
	return []*domain.User{{ID: "u2"}}, f.err
}
func (f *fakeFollowUC) IsFollowing(_ context.Context, _ string, _ string) (bool, error) {
	return true, f.err
}

func setupFollowRouter(uc domain.IFollowUsecase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewFollowController(uc)
	r.POST("/follow", h.FollowUser)
	r.DELETE("/unfollow", h.UnfollowUser)
	r.GET("/users/:user_id/followers", h.GetFollowers)
	r.GET("/users/:user_id/following", h.GetFollowing)
	// Avoid conflict with /users/:user_id by using a different prefix
	r.GET("/followcheck/:follower_id/:followee_id", h.IsFollowing)
	return r
}

func TestFollowController_HappyPaths(t *testing.T) {
	uc := &fakeFollowUC{}
	r := setupFollowRouter(uc)

	// Follow
	body, _ := json.Marshal(dtodlv.FollowRequest{FollowerID: "u1", FolloweeID: "u2"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/follow", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	// Unfollow
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodDelete, "/unfollow", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	// Followers
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/users/u2/followers", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	// Following
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/users/u1/following", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	// IsFollowing
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/followcheck/u1/u2", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
