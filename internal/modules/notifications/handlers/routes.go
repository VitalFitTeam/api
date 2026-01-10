package notihandlers

import (
	appservices "github.com/vitalfit/api/internal/app/services"
)

type NotificationHandlersInterface interface {
}

type NotificationHandlers struct {
	services appservices.Services
}

func NewNotificationHandlers(services appservices.Services) *NotificationHandlers {
	return &NotificationHandlers{
		services: services,
	}
}
