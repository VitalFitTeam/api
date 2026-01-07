package auditrepository

import (
	"context"

	auditdomain "github.com/vitalfit/api/internal/modules/audit/domain"
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

func (s *AuditStore) CreateLog(ctx context.Context, log *auditdomain.AuditLog) error {
	return s.db.WithContext(ctx).Create(log).Error
}
