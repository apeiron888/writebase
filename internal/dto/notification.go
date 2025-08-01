package dto

type NotificationRequest struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

type NotificationResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Message   string `json:"message"`
	Read      bool   `json:"read"`
	CreatedAt int64  `json:"created_at"`
}
