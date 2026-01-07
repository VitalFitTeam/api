package auditrepository

import (
	"gorm.io/gorm"
)

type AuditStore struct {
	db *gorm.DB
}

func NewAuditStore(db *gorm.DB) *AuditStore {
	return &AuditStore{
		db: db,
	}
}
