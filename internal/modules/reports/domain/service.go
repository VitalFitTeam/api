package reportdomain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ReportServiceInterface interface {
	GetGlobalSalesStats(ctx context.Context) (*GlobalSalesStats, error)
	GetTotalSales(ctx context.Context) (*TotalSalesStats, error)
	GetMonthlySalesKPI(ctx context.Context, branchID *uuid.UUID) (*KPICard, error)
	GetActiveMembersKPI(ctx context.Context, branchID *uuid.UUID) (*KPICard, error)
	GetOccupancyKPI(ctx context.Context, branchID *uuid.UUID) (*KPICard, error)
	GetTopBranchesPerformance(ctx context.Context) ([]BranchPerformance, error)
	GetTotalClientsStat(ctx context.Context) (int64, error)
	GetActiveBranchesCount(ctx context.Context) (int64, error)
	GetSalesByCategory(ctx context.Context, start, end time.Time) ([]ChartData, error)
	GetTopInstructorsByAttendance(ctx context.Context, start, end time.Time) ([]ChartData, error)
	GetSalesByPaymentMethod(ctx context.Context, start, end time.Time) ([]ChartData, error)
	GetSalesByHour(ctx context.Context, start, end time.Time) ([]ChartData, error)
	GetMostUsedServices(ctx context.Context, start, end time.Time) ([]ChartData, error)
	GetWeeklySalesChart(ctx context.Context, branchID *uuid.UUID) ([]ChartData, error)
	GetActivityHeatmap(ctx context.Context, branchID *uuid.UUID) ([]HeatmapPoint, error)
	GetClassOccupancyChart(ctx context.Context, branchID *uuid.UUID) ([]ChartData, error)
}
