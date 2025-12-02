package reportservices

import (
	"context"
	"time"

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
