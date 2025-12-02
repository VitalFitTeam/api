package reportrepository

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	authdomain "github.com/vitalfit/api/internal/modules/auth/domain"
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

func (rs *ReportStore) GetTopBranchesPerformance(ctx context.Context) ([]reportdomain.BranchPerformance, error) {
	now := time.Now()
	currentMonthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)

	// Filtramos las facturas válidas (excluyendo Void)
	validStatuses := []billingdomain.InvoiceStatus{
		billingdomain.InvoiceStatusPaid,
		billingdomain.InvoiceStatusUnpaid,
		billingdomain.InvoiceStatusOverdue,
	}

	type resultRaw struct {
		BranchName   string          `gorm:"column:branch_name"`
		CurrentTotal decimal.Decimal `gorm:"column:current_total"`
		PrevTotal    decimal.Decimal `gorm:"column:prev_total"`
	}

	var rawResults []resultRaw

	// Query optimizada
	err := rs.db.Table("invoices").
		Select(`
            branch.name as branch_name,
            SUM(CASE WHEN invoices.issue_date >= ? THEN invoices.total_amount ELSE 0 END) as current_total,
            SUM(CASE WHEN invoices.issue_date >= ? AND invoices.issue_date < ? THEN invoices.total_amount ELSE 0 END) as prev_total
        `, currentMonthStart, prevMonthStart, currentMonthStart).
		Joins("JOIN branch ON branch.branch_id = invoices.branch_id").
		Where("invoices.issue_date >= ? AND invoices.status IN ?", prevMonthStart, validStatuses).
		Group("branch.name").
		Order("current_total DESC").
		Limit(5).
		Scan(&rawResults).Error

	if err != nil {
		return nil, err
	}

	// Mapeo y lógica de etiquetas
	var performanceList []reportdomain.BranchPerformance

	for _, res := range rawResults {
		percentChange := 0.0
		if !res.PrevTotal.IsZero() {
			diff := res.CurrentTotal.Sub(res.PrevTotal)
			val, _ := diff.Div(res.PrevTotal).Mul(decimal.NewFromInt(100)).Float64()
			percentChange = val
		} else if !res.CurrentTotal.IsZero() {
			percentChange = 100.0
		}

		label := "Bueno"
		trend := "up"

		if percentChange >= 10 {
			label = "Excelent"
		} else if percentChange < 0 {
			label = "Attention"
			trend = "down"
		} else if percentChange < 5 {
			label = "Regular"
		}

		performanceList = append(performanceList, reportdomain.BranchPerformance{
			BranchName:    res.BranchName,
			TotalSales:    res.CurrentTotal,
			PercentChange: percentChange,
			Label:         label,
			Trend:         trend,
		})
	}

	return performanceList, nil
}

func (rs *ReportStore) GetTotalClientsStat(ctx context.Context) (int64, error) {
	var totalClients int64
	err := rs.db.WithContext(ctx).Model(&authdomain.Users{}).
		Joins("JOIN roles ON roles.role_id = users.role_id").
		Where("roles.name = ?", "client").
		Count(&totalClients).Error
	if err != nil {
		return 0, err
	}
	return totalClients, nil
}
