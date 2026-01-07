package auditdomain

import "context"

type AuditService interface {
	CreateLog(ctx context.Context, log *AuditLog) error
}
