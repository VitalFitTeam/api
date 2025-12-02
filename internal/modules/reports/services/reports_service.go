package reportservices

import (
	"context"

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
