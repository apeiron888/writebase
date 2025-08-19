package usecase

import (
	"context"
	"testing"
	"write_base/internal/domain"
)

type reactionRepoMock struct {
	AddReactionFn        func(ctx context.Context, r *domain.Reaction) error
	RemoveReactionFn     func(ctx context.Context, id string) error
	GetReactionsByPostFn func(ctx context.Context, postID string) ([]*domain.Reaction, error)
	GetReactionsByUserFn func(ctx context.Context, userID string) ([]*domain.Reaction, error)
	CountReactionsFn     func(ctx context.Context, postID string, t domain.ReactionType) (int, error)
}

func (m *reactionRepoMock) AddReaction(ctx context.Context, r *domain.Reaction) error {
	return m.AddReactionFn(ctx, r)
}
func (m *reactionRepoMock) RemoveReaction(ctx context.Context, id string) error {
	return m.RemoveReactionFn(ctx, id)
}
func (m *reactionRepoMock) GetReactionsByPost(ctx context.Context, postID string) ([]*domain.Reaction, error) {
	return m.GetReactionsByPostFn(ctx, postID)
}
func (m *reactionRepoMock) GetReactionsByUser(ctx context.Context, userID string) ([]*domain.Reaction, error) {
	return m.GetReactionsByUserFn(ctx, userID)
}
func (m *reactionRepoMock) CountReactions(ctx context.Context, postID string, t domain.ReactionType) (int, error) {
	return m.CountReactionsFn(ctx, postID, t)
}

func TestReactionService_Basic(t *testing.T) {
	repo := &reactionRepoMock{
		AddReactionFn:    func(ctx context.Context, r *domain.Reaction) error { return nil },
		RemoveReactionFn: func(ctx context.Context, id string) error { return nil },
		GetReactionsByPostFn: func(ctx context.Context, postID string) ([]*domain.Reaction, error) {
			return []*domain.Reaction{{ID: "r1"}}, nil
		},
		GetReactionsByUserFn: func(ctx context.Context, userID string) ([]*domain.Reaction, error) {
			return []*domain.Reaction{{ID: "r2"}}, nil
		},
		CountReactionsFn: func(ctx context.Context, postID string, t domain.ReactionType) (int, error) { return 2, nil },
	}
	s := NewReactionService(repo)
	if err := s.AddReaction(context.Background(), &domain.Reaction{ID: "r1"}); err != nil {
		t.Fatal(err)
	}
	if err := s.RemoveReaction(context.Background(), "r1"); err != nil {
		t.Fatal(err)
	}
	if list, err := s.GetReactionsByPost(context.Background(), "p1"); err != nil || len(list) != 1 {
		t.Fatalf("byPost bad")
	}
	if list, err := s.GetReactionsByUser(context.Background(), "u1"); err != nil || len(list) != 1 {
		t.Fatalf("byUser bad")
	}
	if n, err := s.CountReactions(context.Background(), "p1", domain.ReactionLike); err != nil || n != 2 {
		t.Fatalf("count bad")
	}
}
