package auditservice

import (
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
