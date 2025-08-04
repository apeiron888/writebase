package usecase

import (
	"context"
	"write_base/internal/domain"
)

type FollowService struct {
	repo domain.IFollowRepository
}

func NewFollowService(repo domain.IFollowRepository) *FollowService {
	return &FollowService{repo: repo}
}

func (s *FollowService) FollowUser(ctx context.Context, followerID, followeeID string) error {
	return s.repo.FollowUser(ctx, followerID, followeeID)
}

func (s *FollowService) UnfollowUser(ctx context.Context, followerID, followeeID string) error {
	return s.repo.UnfollowUser(ctx, followerID, followeeID)
}

func (s *FollowService) GetFollowers(ctx context.Context, userID string) ([]*domain.User, error) {
	return s.repo.GetFollowers(ctx, userID)
}

func (s *FollowService) GetFollowing(ctx context.Context, userID string) ([]*domain.User, error) {
	return s.repo.GetFollowing(ctx, userID)
}

func (s *FollowService) IsFollowing(ctx context.Context, followerID, followeeID string) (bool, error) {
	return s.repo.IsFollowing(ctx, followerID, followeeID)
}
