package dto

type CommentResponse struct {
 ID        string  `json:"id" bson:"_id,omitempty"`
 PostID    string  `json:"post_id" bson:"post_id"`
 UserID    string  `json:"user_id" bson:"user_id"`
 ParentID  *string `json:"parent_id,omitempty" bson:"parent_id,omitempty"`
 Content   string  `json:"content" bson:"content"`
 CreatedAt int64   `json:"created_at" bson:"created_at"`
 UpdatedAt int64   `json:"updated_at" bson:"updated_at"`
}
