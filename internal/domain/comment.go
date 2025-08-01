package domain

import "context"

type Comment struct {
	ID        string
	PostID    string
	UserID    string
	ParentID  *string // nil if top-level comment
	Content   string
	CreatedAt int64
	UpdatedAt int64
}

// ICommentRepository defines the interface for comment repository operations.
type ICommentRepository interface {
	Create(ctx context.Context, comment *Comment) error
	Update(ctx context.Context, comment *Comment) error
	Delete(ctx context.Context, commentID string) error
	GetByID(ctx context.Context, commentID string) (*Comment, error)
	GetByPostID(ctx context.Context, postID string) ([]*Comment, error)
	GetByUserID(ctx context.Context, userID string) ([]*Comment, error)
	GetReplies(ctx context.Context, parentID string) ([]*Comment, error)
}