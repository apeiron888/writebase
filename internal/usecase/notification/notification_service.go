package notification

import (
	"context"
	"starter/internal/domain"
)

type NotificationService struct {
	repo domain.INotificationRepository
}

func NewNotificationService(repo domain.INotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) CreateNotification(ctx context.Context, notification *domain.Notification) error {
	return s.repo.CreateNotification(ctx, notification)
}

func (s *NotificationService) GetNotificationsByUser(ctx context.Context, userID string) ([]*domain.Notification, error) {
	return s.repo.GetNotificationsByUser(ctx, userID)
}

func (s *NotificationService) MarkAsRead(ctx context.Context, notificationID string) error {
	return s.repo.MarkAsRead(ctx, notificationID)
}
