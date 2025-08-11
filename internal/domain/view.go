package domain

import (
	"time"
	"context"
)

type View struct {
	ID        string
	UserID    string
	ArticleID string
	ClientIP  string
	CreatedAt time.Time
}


// type View struct {
// 	ID        string    `json:"id"`
// 	UserID    string    `json:"user_id,omitempty"` // Optional for authenticated users
// 	ArticleID string    `json:"article_id"`
// 	ClientIP  string    `json:"client_ip"`
// 	CreatedAt time.Time `json:"created_at"`
// }


type ViewRepository interface {
	Create(ctx context.Context, view *View) error
}

type ViewUsecase interface {
	RecordView(ctx context.Context, userID, articleID, clientIP string) error
}