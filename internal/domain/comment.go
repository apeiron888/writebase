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

// usecase interface for comment operations
type ICommentUsecase interface {
	CreateComment(ctx context.Context, comment *Comment) error
	UpdateComment(ctx context.Context, comment *Comment) error
	DeleteComment(ctx context.Context, commentID string) error
	GetCommentByID(ctx context.Context, commentID string) (*Comment, error)
	GetCommentsByPostID(ctx context.Context, postID string) ([]*Comment, error)
	GetCommentsByUserID(ctx context.Context, userID string) ([]*Comment, error)
	GetReplies(ctx context.Context, parentID string) ([]*Comment, error)
}