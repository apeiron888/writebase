package repository

import (
	"context"
	"starter/internal/domain"
)

type NotificationRepositoryStub struct {
	notifications map[string][]*domain.Notification // userID -> notifications
}

func NewNotificationRepositoryStub() *NotificationRepositoryStub {
	return &NotificationRepositoryStub{notifications: make(map[string][]*domain.Notification)}
}

func (r *NotificationRepositoryStub) CreateNotification(ctx context.Context, notification *domain.Notification) error {
	r.notifications[notification.UserID] = append(r.notifications[notification.UserID], notification)
	return nil
}

func (r *NotificationRepositoryStub) GetNotificationsByUser(ctx context.Context, userID string) ([]*domain.Notification, error) {
	return r.notifications[userID], nil
}

func (r *NotificationRepositoryStub) MarkAsRead(ctx context.Context, notificationID string) error {
	for _, notifs := range r.notifications {
		for _, n := range notifs {
			if n.ID == notificationID {
				n.Read = true
				break
			}
		}
	}
	return nil
}
