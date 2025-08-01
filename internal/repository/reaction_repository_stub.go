package repository

import (
	"context"
	"starter/internal/domain"
)

type ReactionRepositoryStub struct {
	reactions map[string]*domain.Reaction
}

func NewReactionRepositoryStub() *ReactionRepositoryStub {
	return &ReactionRepositoryStub{reactions: make(map[string]*domain.Reaction)}
}

func (r *ReactionRepositoryStub) AddReaction(ctx context.Context, reaction *domain.Reaction) error {
	r.reactions[reaction.ID] = reaction
	return nil
}

func (r *ReactionRepositoryStub) RemoveReaction(ctx context.Context, reactionID string) error {
	delete(r.reactions, reactionID)
	return nil
}

func (r *ReactionRepositoryStub) GetReactionsByPost(ctx context.Context, postID string) ([]*domain.Reaction, error) {
	var res []*domain.Reaction
	for _, react := range r.reactions {
		if react.PostID == postID {
			res = append(res, react)
		}
	}
	return res, nil
}

func (r *ReactionRepositoryStub) GetReactionsByUser(ctx context.Context, userID string) ([]*domain.Reaction, error) {
	var res []*domain.Reaction
	for _, react := range r.reactions {
		if react.UserID == userID {
			res = append(res, react)
		}
	}
	return res, nil
}

func (r *ReactionRepositoryStub) CountReactions(ctx context.Context, postID string, reactionType domain.ReactionType) (int, error) {
	count := 0
	for _, react := range r.reactions {
		if react.PostID == postID && react.Type == reactionType {
			count++
		}
	}
	return count, nil
}
