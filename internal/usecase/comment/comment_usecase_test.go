package usecase_test

import (
	"context"
	"testing"
	"write_base/internal/domain"
	"write_base/internal/mocks"
	usecase "write_base/internal/usecase/comment"
)

func TestCommentUsecase_BasicFlows(t *testing.T) {
	ctx := context.Background()
	repo := mocks.NewMockICommentRepository(t)
	uc := usecase.NewCommentUsecase(repo)

	c := &domain.Comment{ID: "c1", PostID: "p1", UserID: "u1", Content: "hi"}
	repo.EXPECT().Create(ctx, c).Return(nil)
	if err := uc.CreateComment(ctx, c); err != nil { t.Fatalf("CreateComment err: %v", err) }

	c.Content = "upd"
	repo.EXPECT().Update(ctx, c).Return(nil)
	if err := uc.UpdateComment(ctx, c); err != nil { t.Fatalf("UpdateComment err: %v", err) }

	repo.EXPECT().GetByID(ctx, "c1").Return(c, nil)
	got, err := uc.GetCommentByID(ctx, "c1")
	if err != nil || got.ID != "c1" { t.Fatalf("GetCommentByID mismatch: %v %v", got, err) }

	repo.EXPECT().GetByPostID(ctx, "p1").Return([]*domain.Comment{c}, nil)
	list, err := uc.GetCommentsByPostID(ctx, "p1")
	if err != nil || len(list) != 1 { t.Fatalf("GetCommentsByPostID: %v %v", list, err) }

	repo.EXPECT().GetByUserID(ctx, "u1").Return([]*domain.Comment{c}, nil)
	ulist, err := uc.GetCommentsByUserID(ctx, "u1")
	if err != nil || len(ulist) != 1 { t.Fatalf("GetCommentsByUserID: %v %v", ulist, err) }

	repo.EXPECT().GetReplies(ctx, "c1").Return([]*domain.Comment{}, nil)
	replies, err := uc.GetReplies(ctx, "c1")
	if err != nil || len(replies) != 0 { t.Fatalf("GetReplies: %v %v", replies, err) }

	repo.EXPECT().Delete(ctx, "c1").Return(nil)
	if err := uc.DeleteComment(ctx, "c1"); err != nil { t.Fatalf("DeleteComment err: %v", err) }
}
