package dto

type ReactionResponse struct {
 ID        string  `json:"id" bson:"_id,omitempty"`
 PostID    string  `json:"post_id" bson:"post_id"`
 UserID    string  `json:"user_id" bson:"user_id"`
 CommentID *string `json:"comment_id,omitempty" bson:"comment_id,omitempty"`
 Type      string  `json:"type" bson:"type"`
 CreatedAt int64   `json:"created_at" bson:"created_at"`
}
