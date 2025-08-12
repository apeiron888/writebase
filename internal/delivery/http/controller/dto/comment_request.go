package dto

type CommentRequest struct {
	PostID   string  `json:"post_id"`
	UserID   string  `json:"user_id"`
	ParentID *string `json:"parent_id,omitempty"`
	Content  string  `json:"content"`
}

type UpdateCommentRequest struct {
	Content string `json:"content"`
}
