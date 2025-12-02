package reportdomain

import "context"

type ReportRepository interface {
	GetGlobalSalesStats(ctx context.Context) (*GlobalSalesStats, error)
}
