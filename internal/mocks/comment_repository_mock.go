package mocks

import (
	"context"
	"write_base/internal/domain"
)

type CommentRepositoryMock struct {
	CreateFn      func(ctx context.Context, comment *domain.Comment) error
	UpdateFn      func(ctx context.Context, comment *domain.Comment) error
	DeleteFn      func(ctx context.Context, commentID string) error
	GetByIDFn     func(ctx context.Context, commentID string) (*domain.Comment, error)
	GetByPostIDFn func(ctx context.Context, postID string) ([]*domain.Comment, error)
	GetByUserIDFn func(ctx context.Context, userID string) ([]*domain.Comment, error)
	GetRepliesFn  func(ctx context.Context, parentID string) ([]*domain.Comment, error)
}

func (m *CommentRepositoryMock) Create(ctx context.Context, c *domain.Comment) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, c)
	}
	return nil
}
func (m *CommentRepositoryMock) Update(ctx context.Context, c *domain.Comment) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, c)
	}
	return nil
}
func (m *CommentRepositoryMock) Delete(ctx context.Context, id string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}
func (m *CommentRepositoryMock) GetByID(ctx context.Context, id string) (*domain.Comment, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, domain.ErrCommentNotFound
}
func (m *CommentRepositoryMock) GetByPostID(ctx context.Context, postID string) ([]*domain.Comment, error) {
	if m.GetByPostIDFn != nil {
		return m.GetByPostIDFn(ctx, postID)
	}
	return nil, nil
}
func (m *CommentRepositoryMock) GetByUserID(ctx context.Context, userID string) ([]*domain.Comment, error) {
	if m.GetByUserIDFn != nil {
		return m.GetByUserIDFn(ctx, userID)
	}
	return nil, nil
}
func (m *CommentRepositoryMock) GetReplies(ctx context.Context, parentID string) ([]*domain.Comment, error) {
	if m.GetRepliesFn != nil {
		return m.GetRepliesFn(ctx, parentID)
	}
	return nil, nil
}
