package domain

import "context"

// Repository interface for follow operations


type IFollowRepository interface {
	FollowUser(ctx context.Context, followerID, followeeID string) error
	UnfollowUser(ctx context.Context, followerID, followeeID string) error
	GetFollowers(ctx context.Context, userID string) ([]*User, error)
	GetFollowing(ctx context.Context, userID string) ([]*User, error)
	IsFollowing(ctx context.Context, followerID, followeeID string) (bool, error)
}

//Usecase interface for follow operations
type IFollowUsecase interface {
	FollowUser(ctx context.Context, followerID, followeeID string) error
	UnfollowUser(ctx context.Context, followerID, followeeID string) error
	GetFollowers(ctx context.Context, userID string) ([]*User, error)
	GetFollowing(ctx context.Context, userID string) ([]*User, error)
	IsFollowing(ctx context.Context, followerID, followeeID string) (bool, error)
}