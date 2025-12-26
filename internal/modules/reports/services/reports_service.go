package reportservices

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	reportdomain "github.com/vitalfit/api/internal/modules/reports/domain"
	"github.com/vitalfit/api/internal/store"
)

type ReportService struct {
	store store.Storage
}

func NewReportService(store store.Storage) *ReportService {
	return &ReportService{
		store: store,
	}

}

func (s *ReportService) GetGlobalSalesStats(ctx context.Context) (*reportdomain.GlobalSalesStats, error) {
	return s.store.Reports.GetGlobalSalesStats(ctx)
}

func (s *ReportService) GetTotalSales(ctx context.Context) (*reportdomain.TotalSalesStats, error) {
	return s.store.Reports.GetTotalSales(ctx)
}

func (s *ReportService) GetMonthlySalesKPI(ctx context.Context, branchID *uuid.UUID) (*reportdomain.KPICard, error) {
	return s.store.Reports.GetMonthlySalesKPI(ctx, branchID)
}

func (s *ReportService) GetActiveMembersKPI(ctx context.Context, branchID *uuid.UUID) (*reportdomain.KPICard, error) {
	return s.store.Reports.GetActiveMembersKPI(ctx, branchID)
}

func (s *ReportService) GetOccupancyKPI(ctx context.Context, branchID *uuid.UUID) (*reportdomain.KPICard, error) {
	return s.store.Reports.GetOccupancyKPI(ctx, branchID)
}

func (s *ReportService) GetTopBranchesPerformance(ctx context.Context) ([]reportdomain.BranchPerformance, error) {
	return s.store.Reports.GetTopBranchesPerformance(ctx)
}

func (s *ReportService) GetTotalClientsStat(ctx context.Context) (int64, error) {
	return s.store.Reports.GetTotalClientsStat(ctx)
}

func (s *ReportService) GetActiveBranchesCount(ctx context.Context) (int64, error) {
	return s.store.Reports.GetActiveBranchesCount(ctx)
}

func (s *ReportService) GetSalesByCategory(ctx context.Context, start, end time.Time) ([]reportdomain.ChartData, error) {
	return s.store.Reports.GetSalesByCategory(ctx, start, end)
}

func (s *ReportService) GetTopInstructorsByAttendance(ctx context.Context, start, end time.Time) ([]reportdomain.ChartData, error) {
	return s.store.Reports.GetTopInstructorsByAttendance(ctx, start, end)
}

func (s *ReportService) GetSalesByPaymentMethod(ctx context.Context, start, end time.Time) ([]reportdomain.ChartData, error) {
	return s.store.Reports.GetSalesByPaymentMethod(ctx, start, end)
}

func (s *ReportService) GetSalesByHour(ctx context.Context, start, end time.Time) ([]reportdomain.ChartData, error) {
	return s.store.Reports.GetSalesByHour(ctx, start, end)
}

func (s *ReportService) GetMostUsedServices(ctx context.Context, start, end time.Time) ([]reportdomain.ChartData, error) {
	return s.store.Reports.GetMostUsedServices(ctx, start, end)
}

func (s *ReportService) GetWeeklySalesChart(ctx context.Context, branchID *uuid.UUID) ([]reportdomain.ChartData, error) {
	return s.store.Reports.GetWeeklySalesChart(ctx, branchID)
}

func (s *ReportService) GetActivityHeatmap(ctx context.Context, branchID *uuid.UUID) ([]reportdomain.HeatmapPoint, error) {
	return s.store.Reports.GetActivityHeatmap(ctx, branchID)
}

func (s *ReportService) GetClassOccupancyChart(ctx context.Context, branchID *uuid.UUID) ([]reportdomain.ChartData, error) {
	return s.store.Reports.GetClassOccupancyChart(ctx, branchID)
}

func (s *ReportService) GetFinancialSummary(ctx context.Context, branchID *uuid.UUID) (*reportdomain.FinancialSummary, error) {
	return s.store.Reports.GetFinancialSummary(ctx, branchID)
}

func (s *ReportService) GetTodayCheckInsStat(ctx context.Context, branchID *uuid.UUID) (int64, error) {
	return s.store.Reports.GetTodayCheckInsStat(ctx, branchID)
}

func (s *ReportService) GetCurrentOccupancyStat(ctx context.Context, branchID *uuid.UUID) (decimal.Decimal, error) {
	return s.store.Reports.GetCurrentOccupancyStat(ctx, branchID)
}

func (s *ReportService) GetClassCapacityRatio(ctx context.Context, classID uuid.UUID) (*reportdomain.ClassCapacityStats, error) {
	return s.store.Reports.GetClassCapacityRatio(ctx, classID)
}

func (s *ReportService) GetUpcomingClassesToday(ctx context.Context, branchID *uuid.UUID) ([]reportdomain.ClassScheduleItem, error) {
	return s.store.Reports.GetUpcomingClassesToday(ctx, branchID)
}
