package reportrepository

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	billingdomain "github.com/vitalfit/api/internal/modules/billing/domain"
	reportdomain "github.com/vitalfit/api/internal/modules/reports/domain"
	"gorm.io/gorm"
)

type ReportStore struct {
	db *gorm.DB
}

func NewReportStore(db *gorm.DB) *ReportStore {
	return &ReportStore{
		db: db}
}

func (rs *ReportStore) GetGlobalSalesStats(ctx context.Context) (*reportdomain.GlobalSalesStats, error) {
	now := time.Now()

	currentMonthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)

	var currentTotal, prevTotal decimal.Decimal

	validStatus := billingdomain.InvoiceStatusPaid

	err := rs.db.WithContext(ctx).Model(&billingdomain.Invoice{}).
		Where("issue_date >= ? AND status = ?", currentMonthStart, validStatus).
		Select("COALESCE(SUM(total_amount), 0)"). // COALESCE evita NULL si no hay ventas
		Scan(&currentTotal).Error
	if err != nil {
		return nil, err
	}

	err = rs.db.WithContext(ctx).Model(&billingdomain.Invoice{}).
		Where("issue_date >= ? AND issue_date < ? AND status = ?", prevMonthStart, currentMonthStart, validStatus).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&prevTotal).Error
	if err != nil {
		return nil, err
	}

	percentageChange := 0.0
	trend := "neutral"

	if !prevTotal.IsZero() {
		diff := currentTotal.Sub(prevTotal)
		percDecimal := diff.Div(prevTotal).Mul(decimal.NewFromInt(100))
		percentageChange, _ = percDecimal.Float64()
	} else if !currentTotal.IsZero() {
		percentageChange = 100.0
	}

	if percentageChange > 5 {
		trend = "up"
	} else if percentageChange < -5 {
		trend = "down"
	}

	return &reportdomain.GlobalSalesStats{
		TotalCurrentMonth: currentTotal,
		TotalLastMonth:    prevTotal,
		PercentageChange:  percentageChange,
		Trend:             trend,
	}, nil

}
