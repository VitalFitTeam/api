package notirepository

import (
	"gorm.io/gorm"
)

type NotificationStore struct {
	db *gorm.DB
}

func NewNotificationStore(db *gorm.DB) *NotificationStore {
	return &NotificationStore{
		db: db}
}
