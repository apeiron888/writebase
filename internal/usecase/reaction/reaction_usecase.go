package usecase

import (
	"context"
	"write_base/internal/domain"
)

type ReactionService struct {
	repo domain.IReactionRepository
}

func NewReactionService(repo domain.IReactionRepository) *ReactionService {
	return &ReactionService{repo: repo}
}

func (s *ReactionService) AddReaction(ctx context.Context, reaction *domain.Reaction) error {
	return s.repo.AddReaction(ctx, reaction)
}

func (s *ReactionService) RemoveReaction(ctx context.Context, reactionID string) error {
	return s.repo.RemoveReaction(ctx, reactionID)
}

func (s *ReactionService) GetReactionsByPost(ctx context.Context, postID string) ([]*domain.Reaction, error) {
	return s.repo.GetReactionsByPost(ctx, postID)
}

func (s *ReactionService) GetReactionsByUser(ctx context.Context, userID string) ([]*domain.Reaction, error) {
	return s.repo.GetReactionsByUser(ctx, userID)
}

func (s *ReactionService) CountReactions(ctx context.Context, postID string, reactionType domain.ReactionType) (int, error) {
	return s.repo.CountReactions(ctx, postID, reactionType)
}
