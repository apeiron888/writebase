package repository

import (
	"context"
	"starter/internal/domain"
)

type FollowRepositoryStub struct {
	follows map[string]map[string]bool // followerID -> followeeID -> true
}

func NewFollowRepositoryStub() *FollowRepositoryStub {
	return &FollowRepositoryStub{follows: make(map[string]map[string]bool)}
}

func (r *FollowRepositoryStub) FollowUser(ctx context.Context, followerID, followeeID string) error {
	if r.follows[followerID] == nil {
		r.follows[followerID] = make(map[string]bool)
	}
	r.follows[followerID][followeeID] = true
	return nil
}

func (r *FollowRepositoryStub) UnfollowUser(ctx context.Context, followerID, followeeID string) error {
	if r.follows[followerID] != nil {
		delete(r.follows[followerID], followeeID)
	}
	return nil
}

func (r *FollowRepositoryStub) GetFollowers(ctx context.Context, userID string) ([]*domain.User, error) {
	var followers []*domain.User
	for followerID, followees := range r.follows {
		if followees[userID] {
			followers = append(followers, &domain.User{ID: followerID})
		}
	}
	return followers, nil
}

func (r *FollowRepositoryStub) GetFollowing(ctx context.Context, userID string) ([]*domain.User, error) {
	var following []*domain.User
	for followeeID := range r.follows[userID] {
		following = append(following, &domain.User{ID: followeeID})
	}
	return following, nil
}

func (r *FollowRepositoryStub) IsFollowing(ctx context.Context, followerID, followeeID string) (bool, error) {
	if r.follows[followerID] != nil {
		return r.follows[followerID][followeeID], nil
	}
	return false, nil
}
