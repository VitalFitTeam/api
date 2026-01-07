package auditservice

import (
	"context"

	auditdomain "github.com/vitalfit/api/internal/modules/audit/domain"
	"github.com/vitalfit/api/internal/store"
)

type AuditService struct {
	store store.Storage
}

func NewAuditService(store store.Storage) *AuditService {
	return &AuditService{
		store: store,
	}
}

func (s *AuditService) CreateLog(ctx context.Context, log *auditdomain.AuditLog) error {
	return s.store.Audit.CreateLog(ctx, log)

}
