package reportdomain

import (
	"context"
)

type ReportServiceInterface interface {
	GetGlobalSalesStats(ctx context.Context) (*GlobalSalesStats, error)
	GetTopBranchesPerformance(ctx context.Context) ([]BranchPerformance, error)
	GetTotalClientsStat(ctx context.Context) (int64, error)
	GetActiveBranchesCount(ctx context.Context) (int64, error)
}
