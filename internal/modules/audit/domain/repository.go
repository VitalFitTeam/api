package auditdomain

import (
	"context"
)

type AuditRepository interface {
	CreateLog(ctx context.Context, log *AuditLog) error
}
