package notiservice

import (
	"context"

	"github.com/google/uuid"
	notidomain "github.com/vitalfit/api/internal/modules/notifications/domain"
	"github.com/vitalfit/api/internal/shared/notifications"
	"github.com/vitalfit/api/internal/store"
	"github.com/vitalfit/api/pkg/pagination"
)

type NotificationService struct {
	store store.Storage
	noti  notifications.PushService
}

func NewNotificationService(store store.Storage, noti notifications.PushService) *NotificationService {
	return &NotificationService{
		store: store,
		noti:  noti,
	}
}

func (s *NotificationService) CreateNotification(ctx context.Context, n *notidomain.Notification) error {
	return s.store.Notification.CreateNotification(ctx, n)
}

func (s *NotificationService) CreateBatchNotifications(ctx context.Context, notifications []notidomain.Notification) error {
	return s.store.Notification.CreateBatchNotifications(ctx, notifications)
}

func (s *NotificationService) GetNotificationsByUserID(ctx context.Context, userID uuid.UUID, fq pagination.PaginatedFeedQuery) ([]notidomain.Notification, error) {
	return s.store.Notification.GetNotificationsByUserID(ctx, userID, fq)
}

func (s *NotificationService) CountUnread(ctx context.Context, userID uuid.UUID) (int64, error) {
	return s.store.Notification.CountUnread(ctx, userID)
}

func (s *NotificationService) MarkAsRead(ctx context.Context, notificationID uuid.UUID) error {
	return s.store.Notification.MarkAsRead(ctx, notificationID)
}

func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	return s.store.Notification.MarkAllAsRead(ctx, userID)
}

func (s *NotificationService) SendBroadcast(ctx context.Context, title, message string) error {
	users, err := s.store.Users.GetAllClients(ctx)
	if err != nil {
		return err
	}

	for _, user := range users {
		session, err := s.store.Session.GetUserSessions(ctx, user.UserID)
		if err != nil {
			return err
		}
		for _, ses := range session {
			if ses.DeviceToken != "" {
				err := s.noti.SendPush(ctx, ses.DeviceToken, title, message, nil)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil

}
