package usecase_test

import (
	"context"
	"testing"
	"write_base/internal/mocks"
	usecase "write_base/internal/usecase/follow"
)

func TestFollowService_Basic(t *testing.T) {
	ctx := context.Background()
	repo := mocks.NewMockIFollowRepository(t)
	svc := usecase.NewFollowService(repo)

	repo.EXPECT().FollowUser(ctx, "u1", "u2").Return(nil)
	if err := svc.FollowUser(ctx, "u1", "u2"); err != nil { t.Fatalf("FollowUser: %v", err) }

	repo.EXPECT().IsFollowing(ctx, "u1", "u2").Return(true, nil)
	ok, err := svc.IsFollowing(ctx, "u1", "u2")
	if err != nil || !ok { t.Fatalf("IsFollowing: %v %v", ok, err) }

	repo.EXPECT().UnfollowUser(ctx, "u1", "u2").Return(nil)
	if err := svc.UnfollowUser(ctx, "u1", "u2"); err != nil { t.Fatalf("UnfollowUser: %v", err) }
}
