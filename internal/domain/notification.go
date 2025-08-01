package domain

import "context"

type Notification struct {
	ID        string
	UserID    string
	Message   string
	Read      bool
	CreatedAt int64
}

type INotificationRepository interface {
	CreateNotification(ctx context.Context, notification *Notification) error
	GetNotificationsByUser(ctx context.Context, userID string) ([]*Notification, error)
	MarkAsRead(ctx context.Context, notificationID string) error
}
