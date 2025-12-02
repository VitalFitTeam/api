package reportdomain

import "context"

type ReportServiceInterface interface {
	GetGlobalSalesStats(ctx context.Context) (*GlobalSalesStats, error)
}
