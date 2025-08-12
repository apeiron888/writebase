package dto

type NotificationResponse struct {
 ID        string `json:"id" bson:"_id,omitempty"`
 UserID    string `json:"user_id" bson:"user_id"`
 Message   string `json:"message" bson:"message"`
 Read      bool   `json:"read" bson:"read"`
 CreatedAt int64  `json:"created_at" bson:"created_at"`
}
