package domain

import (
	"time"
	"context"
)

type ArticleView struct {
    ArticleID string
    UserID    string
    ViewedAt  time.Time
    IPAddress string
}

type ArticleClap struct {
    ArticleID string
    UserID    string
    Count     int
    LastClap  time.Time
}

type IViewUsecase interface {
    RecordView(ctx context.Context, articleID, userID, ipAddress string) error
    HasViewed(ctx context.Context, userID, articleID string, within time.Duration) (*ArticleView, error)
    CountViews(ctx context.Context, articleID string) (int64, error)
}
type IClapUsecase interface {
    AddClap(ctx context.Context, articleID, userID string) error
    GetClap(ctx context.Context, userID, articleID string) (*ArticleClap, error)
    RemoveClap(ctx context.Context, userID, articleID string, count int) error
}
type IViewRepository interface {
    RecordView(ctx context.Context, view *ArticleView) error
    HasViewed(ctx context.Context, userID, articleID string, within time.Duration) (*ArticleView, error)
    CountViews(ctx context.Context, articleID string) (int64, error)
}

type IClapRepository interface {
	AddClap(ctx context.Context, clap *ArticleClap) error
	GetClap(ctx context.Context, userID, articleID string) (*ArticleClap, error)
	RemoveClap(ctx context.Context, userID, articleID string, count int) error
}