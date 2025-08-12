package usecase_test

import (
	"context"
	"testing"
	"write_base/internal/domain"
	"write_base/internal/mocks"
	usecase "write_base/internal/usecase/reaction"
)

func TestReactionService_Basic(t *testing.T) {
	ctx := context.Background()
	repo := mocks.NewMockIReactionRepository(t)
	svc := usecase.NewReactionService(repo)

	r := &domain.Reaction{ID:"r1", PostID:"p1", UserID:"u1", Type: domain.ReactionLike}
	repo.EXPECT().AddReaction(ctx, r).Return(nil)
	if err := svc.AddReaction(ctx, r); err != nil { t.Fatalf("AddReaction: %v", err) }

	repo.EXPECT().GetReactionsByPost(ctx, "p1").Return([]*domain.Reaction{r}, nil)
	got, err := svc.GetReactionsByPost(ctx, "p1")
	if err != nil || len(got) != 1 { t.Fatalf("GetReactionsByPost: %v %v", got, err) }

	repo.EXPECT().CountReactions(ctx, "p1", domain.ReactionLike).Return(1, nil)
	n, err := svc.CountReactions(ctx, "p1", domain.ReactionLike)
	if err != nil || n != 1 { t.Fatalf("CountReactions: %d %v", n, err) }

	repo.EXPECT().RemoveReaction(ctx, "r1").Return(nil)
	if err := svc.RemoveReaction(ctx, "r1"); err != nil { t.Fatalf("RemoveReaction: %v", err) }
}
