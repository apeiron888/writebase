package dto

type ReactionRequest struct {
	PostID    string  `json:"post_id"`
	UserID    string  `json:"user_id"`
	CommentID *string `json:"comment_id,omitempty"`
	Type      string  `json:"type"`
}
