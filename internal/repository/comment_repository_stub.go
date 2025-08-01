package repository

import (
	"context"
	"starter/internal/domain"
)

type CommentRepositoryStub struct {
	comments map[string]*domain.Comment
}

func NewCommentRepositoryStub() *CommentRepositoryStub {
	return &CommentRepositoryStub{comments: make(map[string]*domain.Comment)}
}

func (r *CommentRepositoryStub) Create(ctx context.Context, comment *domain.Comment) error {
	r.comments[comment.ID] = comment
	return nil
}

func (r *CommentRepositoryStub) Update(ctx context.Context, comment *domain.Comment) error {
	r.comments[comment.ID] = comment
	return nil
}

func (r *CommentRepositoryStub) Delete(ctx context.Context, commentID string) error {
	delete(r.comments, commentID)
	return nil
}

func (r *CommentRepositoryStub) GetByID(ctx context.Context, commentID string) (*domain.Comment, error) {
	c, ok := r.comments[commentID]
	if !ok {
		return nil, nil
	}
	return c, nil
}

func (r *CommentRepositoryStub) GetByPostID(ctx context.Context, postID string) ([]*domain.Comment, error) {
	var res []*domain.Comment
	for _, c := range r.comments {
		if c.PostID == postID {
			res = append(res, c)
		}
	}
	return res, nil
}

func (r *CommentRepositoryStub) GetByUserID(ctx context.Context, userID string) ([]*domain.Comment, error) {
	var res []*domain.Comment
	for _, c := range r.comments {
		if c.UserID == userID {
			res = append(res, c)
		}
	}
	return res, nil
}

func (r *CommentRepositoryStub) GetReplies(ctx context.Context, parentID string) ([]*domain.Comment, error) {
	var res []*domain.Comment
	for _, c := range r.comments {
		if c.ParentID != nil && *c.ParentID == parentID {
			res = append(res, c)
		}
	}
	return res, nil
}
