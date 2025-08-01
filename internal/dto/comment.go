package dto

type CommentRequest struct {
	PostID   string  `json:"post_id"`
	UserID   string  `json:"user_id"`
	ParentID *string `json:"parent_id,omitempty"`
	Content  string  `json:"content"`
}

type CommentResponse struct {
	ID        string  `json:"id"`
	PostID    string  `json:"post_id"`
	UserID    string  `json:"user_id"`
	ParentID  *string `json:"parent_id,omitempty"`
	Content   string  `json:"content"`
	CreatedAt int64   `json:"created_at"`
	UpdatedAt int64   `json:"updated_at"`
}
