package reportdomain

import (
	"context"
	"time"
)

type ReportServiceInterface interface {
	GetGlobalSalesStats(ctx context.Context) (*GlobalSalesStats, error)
	GetTopBranchesPerformance(ctx context.Context) ([]BranchPerformance, error)
	GetTotalClientsStat(ctx context.Context) (int64, error)
	GetActiveBranchesCount(ctx context.Context) (int64, error)
	GetSalesByCategory(ctx context.Context, start, end time.Time) ([]ChartData, error)
	GetTopInstructorsByAttendance(ctx context.Context, start, end time.Time) ([]ChartData, error)
	GetSalesByPaymentMethod(ctx context.Context, start, end time.Time) ([]ChartData, error)
	GetSalesByHour(ctx context.Context, start, end time.Time) ([]ChartData, error)
	GetMostUsedServices(ctx context.Context, start, end time.Time) ([]ChartData, error)
}
