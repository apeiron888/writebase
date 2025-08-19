package usecase

import (
	"context"
	"testing"
	"write_base/internal/domain"
	"write_base/internal/mocks"
)

func TestCommentUsecase_CRUD(t *testing.T) {
	repo := &mocks.CommentRepositoryMock{}
	uc := NewCommentUsecase(repo)

	c := &domain.Comment{ID: "c1", PostID: "p1", UserID: "u1", Content: "hello"}

	// Create
	created := false
	repo.CreateFn = func(ctx context.Context, comment *domain.Comment) error { created = true; return nil }
	if err := uc.CreateComment(context.Background(), c); err != nil || !created {
		t.Fatalf("create failed")
	}

	// Update
	updated := false
	repo.UpdateFn = func(ctx context.Context, comment *domain.Comment) error { updated = true; return nil }
	if err := uc.UpdateComment(context.Background(), c); err != nil || !updated {
		t.Fatalf("update failed")
	}

	// GetByID
	repo.GetByIDFn = func(ctx context.Context, id string) (*domain.Comment, error) { return c, nil }
	if out, err := uc.GetCommentByID(context.Background(), "c1"); err != nil || out.ID != "c1" {
		t.Fatalf("get by id failed")
	}

	// GetByPostID
	repo.GetByPostIDFn = func(ctx context.Context, pid string) ([]*domain.Comment, error) { return []*domain.Comment{c}, nil }
	if out, err := uc.GetCommentsByPostID(context.Background(), "p1"); err != nil || len(out) != 1 {
		t.Fatalf("get by post id failed")
	}

	// GetByUserID
	repo.GetByUserIDFn = func(ctx context.Context, uid string) ([]*domain.Comment, error) { return []*domain.Comment{c}, nil }
	if out, err := uc.GetCommentsByUserID(context.Background(), "u1"); err != nil || len(out) != 1 {
		t.Fatalf("get by user id failed")
	}

	// GetReplies
	repo.GetRepliesFn = func(ctx context.Context, parentID string) ([]*domain.Comment, error) {
		return []*domain.Comment{c}, nil
	}
	if out, err := uc.GetReplies(context.Background(), "c0"); err != nil || len(out) != 1 {
		t.Fatalf("get replies failed")
	}

	// Delete
	deleted := false
	repo.DeleteFn = func(ctx context.Context, id string) error { deleted = true; return nil }
	if err := uc.DeleteComment(context.Background(), "c1"); err != nil || !deleted {
		t.Fatalf("delete failed")
	}
}
