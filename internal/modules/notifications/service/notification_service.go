package notiservice

import (
	"github.com/vitalfit/api/internal/store"
)

type NotificationService struct {
	store store.Storage
}
