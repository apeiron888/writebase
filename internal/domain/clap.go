package domain

import (
	"time"
	"context"
)

const MaxClapsPerUser = 10

type Clap struct {
	ID        string
	UserID    string
	ArticleID string
	Count     int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// type Clap struct {
// 	ID        string    `json:"id"`
// 	UserID    string    `json:"user_id"`
// 	ArticleID string    `json:"article_id"`
// 	Count     int       `json:"count"` // Number of claps (1-10)
// 	CreatedAt time.Time `json:"created_at"`
// 	UpdatedAt time.Time `json:"updated_at"`
// }

type ClapRepository interface {
	Create(ctx context.Context, clap *Clap) error
	Update(ctx context.Context, clap *Clap) error
	GetByUserAndArticle(ctx context.Context, userID, articleID string) (*Clap, error)
	GetArticleClapCount(ctx context.Context, articleID string) (int, error)
}

type ClapUsecase interface {
	AddClap(ctx context.Context, userID, articleID string) (int, error) // Returns total clap count
}