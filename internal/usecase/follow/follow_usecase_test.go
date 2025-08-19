package usecase

import (
	"context"
	"testing"
	"write_base/internal/domain"
)

type followRepoMock struct {
	FollowUserFn   func(ctx context.Context, followerID, followeeID string) error
	UnfollowUserFn func(ctx context.Context, followerID, followeeID string) error
	GetFollowersFn func(ctx context.Context, userID string) ([]*domain.User, error)
	GetFollowingFn func(ctx context.Context, userID string) ([]*domain.User, error)
	IsFollowingFn  func(ctx context.Context, followerID, followeeID string) (bool, error)
}

func (m *followRepoMock) FollowUser(ctx context.Context, followerID, followeeID string) error {
	return m.FollowUserFn(ctx, followerID, followeeID)
}
func (m *followRepoMock) UnfollowUser(ctx context.Context, followerID, followeeID string) error {
	return m.UnfollowUserFn(ctx, followerID, followeeID)
}
func (m *followRepoMock) GetFollowers(ctx context.Context, userID string) ([]*domain.User, error) {
	return m.GetFollowersFn(ctx, userID)
}
func (m *followRepoMock) GetFollowing(ctx context.Context, userID string) ([]*domain.User, error) {
	return m.GetFollowingFn(ctx, userID)
}
func (m *followRepoMock) IsFollowing(ctx context.Context, followerID, followeeID string) (bool, error) {
	return m.IsFollowingFn(ctx, followerID, followeeID)
}

func TestFollowService_Basic(t *testing.T) {
	repo := &followRepoMock{
		FollowUserFn:   func(ctx context.Context, f, e string) error { return nil },
		UnfollowUserFn: func(ctx context.Context, f, e string) error { return nil },
		GetFollowersFn: func(ctx context.Context, u string) ([]*domain.User, error) { return []*domain.User{{ID: "u2"}}, nil },
		GetFollowingFn: func(ctx context.Context, u string) ([]*domain.User, error) { return []*domain.User{{ID: "u3"}}, nil },
		IsFollowingFn:  func(ctx context.Context, f, e string) (bool, error) { return true, nil },
	}
	s := NewFollowService(repo)
	if err := s.FollowUser(context.Background(), "u1", "u2"); err != nil {
		t.Fatal(err)
	}
	if err := s.UnfollowUser(context.Background(), "u1", "u2"); err != nil {
		t.Fatal(err)
	}
	if list, err := s.GetFollowers(context.Background(), "u1"); err != nil || len(list) != 1 {
		t.Fatalf("followers bad")
	}
	if list, err := s.GetFollowing(context.Background(), "u1"); err != nil || len(list) != 1 {
		t.Fatalf("following bad")
	}
	if ok, err := s.IsFollowing(context.Background(), "u1", "u2"); err != nil || !ok {
		t.Fatalf("isFollowing bad")
	}
}
