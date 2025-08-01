package domain

import "context"

type ReactionType string

const (
	ReactionLike    ReactionType = "like"
	ReactionDislike ReactionType = "dislike"
)

type Reaction struct {
	ID        string
	PostID    string
	UserID    string
	CommentID *string
	Type      ReactionType
	CreatedAt int64
}

type IReactionRepository interface {
	AddReaction(ctx context.Context, reaction *Reaction) error
	RemoveReaction(ctx context.Context, reactionID string) error
	GetReactionsByPost(ctx context.Context, postID string) ([]*Reaction, error)
	GetReactionsByUser(ctx context.Context, userID string) ([]*Reaction, error)
	CountReactions(ctx context.Context, postID string, reactionType ReactionType) (int, error)
}
