package comment

import (
	"context"
	"starter/internal/domain"
)

type CommentUsecase struct {
	repo domain.ICommentRepository
}

func NewCommentUsecase(repo domain.ICommentRepository) *CommentUsecase {
	return &CommentUsecase{repo: repo}
}

func (uc *CommentUsecase) CreateComment(ctx context.Context, comment *domain.Comment) error {
	return uc.repo.Create(ctx, comment)
}

func (uc *CommentUsecase) UpdateComment(ctx context.Context, comment *domain.Comment) error {
	return uc.repo.Update(ctx, comment)
}

func (uc *CommentUsecase) DeleteComment(ctx context.Context, commentID string) error {
	return uc.repo.Delete(ctx, commentID)
}

func (uc *CommentUsecase) GetCommentByID(ctx context.Context, commentID string) (*domain.Comment, error) {
	return uc.repo.GetByID(ctx, commentID)
}

func (uc *CommentUsecase) GetCommentsByPostID(ctx context.Context, postID string) ([]*domain.Comment, error) {
	return uc.repo.GetByPostID(ctx, postID)
}

func (uc *CommentUsecase) GetCommentsByUserID(ctx context.Context, userID string) ([]*domain.Comment, error) {
	return uc.repo.GetByUserID(ctx, userID)
}

func (uc *CommentUsecase) GetReplies(ctx context.Context, parentID string) ([]*domain.Comment, error) {
	return uc.repo.GetReplies(ctx, parentID)
}
