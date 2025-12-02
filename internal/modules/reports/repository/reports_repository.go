package reportrepository

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	accessdomain "github.com/vitalfit/api/internal/modules/access/domain"
	authdomain "github.com/vitalfit/api/internal/modules/auth/domain"
	billingdomain "github.com/vitalfit/api/internal/modules/billing/domain"
	branchdomain "github.com/vitalfit/api/internal/modules/branches/domain"
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
		Select("COALESCE(SUM(total_amount), 0)").
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

		label := "Good"
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

func (rs *ReportStore) GetActiveBranchesCount(ctx context.Context) (int64, error) {
	var count int64
	err := rs.db.WithContext(ctx).Model(&branchdomain.Branch{}).
		Where("status = ?", branchdomain.BranchStatusActive).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (rs *ReportStore) GetSalesByCategory(ctx context.Context, start, end time.Time) ([]reportdomain.ChartData, error) {
	var results []reportdomain.ChartData
	validStatuses := []billingdomain.InvoiceStatus{
		billingdomain.InvoiceStatusPaid,
		billingdomain.InvoiceStatusUnpaid,
		billingdomain.InvoiceStatusOverdue,
	}

	query := `
        SELECT
            CASE
                WHEN ii.membership_type_id IS NOT NULL THEN 'Memberships'
                WHEN ii.package_id IS NOT NULL THEN 'Packages'
                WHEN sc.name IS NOT NULL THEN sc.name
                ELSE 'Otros'
            END as label,
            SUM(ii.total_line) as value
        FROM invoice_items ii
        JOIN invoices i ON i.invoice_id = ii.invoice_id
        LEFT JOIN services s ON s.service_id = ii.service_id
        LEFT JOIN service_categories sc ON sc.category_id = s.category_id
        WHERE i.issue_date BETWEEN ? AND ?
        AND i.status IN ?
        GROUP BY label
        ORDER BY value DESC
    `

	err := rs.db.WithContext(ctx).Raw(query, start, end, validStatuses).Scan(&results).Error
	return results, err
}

func (rs *ReportStore) GetTopInstructorsByAttendance(ctx context.Context, start, end time.Time) ([]reportdomain.ChartData, error) {
	var results []reportdomain.ChartData

	err := rs.db.WithContext(ctx).Table("attendance_log as al").
		Select("u.first_name || ' ' || u.last_name as label, COUNT(*) as value").
		Joins("JOIN classes c ON c.class_id = al.schedule_id").
		Joins("JOIN instructors i ON i.instructor_id = c.instructor_id").
		Joins("JOIN users u ON u.user_id = i.user_id").
		Where("al.check_in_time BETWEEN ? AND ?", start, end).
		Where("al.status = ?", accessdomain.AttendanceStatusAttended).
		Group("label").
		Order("value DESC").
		Limit(5).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (rs *ReportStore) GetSalesByPaymentMethod(ctx context.Context, start, end time.Time) ([]reportdomain.ChartData, error) {
	var results []reportdomain.ChartData
	validStatuses := []billingdomain.InvoiceStatus{
		billingdomain.InvoiceStatusPaid,
		billingdomain.InvoiceStatusUnpaid,
		billingdomain.InvoiceStatusOverdue,
	}

	err := rs.db.WithContext(ctx).Table("payments p").
		Select("pm.name as label, SUM(p.amount_base) as value").
		Joins("JOIN invoices i ON i.invoice_id = p.invoice_id").
		Joins("JOIN payment_methods pm ON pm.method_id = p.payment_method_id").
		Where("i.issue_date BETWEEN ? AND ?", start, end).
		Where("i.status IN ?", validStatuses).
		Where("p.status = ?", billingdomain.PaymentStatusCompleted).
		Group("pm.name").
		Order("value DESC").
		Scan(&results).Error

	return results, err
}

func (rs *ReportStore) GetSalesByHour(ctx context.Context, start, end time.Time) ([]reportdomain.ChartData, error) {
	var results []reportdomain.ChartData
	validStatuses := []billingdomain.InvoiceStatus{
		billingdomain.InvoiceStatusPaid,
		billingdomain.InvoiceStatusUnpaid,
		billingdomain.InvoiceStatusOverdue,
	}

	err := rs.db.WithContext(ctx).Model(&billingdomain.Invoice{}).
		Select("EXTRACT(HOUR FROM created_at) as hour, SUM(total_amount) as value").
		Where("created_at BETWEEN ? AND ?", start, end).
		Where("status IN ?", validStatuses).
		Group("hour").
		Order("hour ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	hourlyMap := make(map[int]decimal.Decimal)
	for _, r := range results {
		hourlyMap[r.Hour] = r.Value
	}

	fullResults := make([]reportdomain.ChartData, 24)
	for i := 0; i < 24; i++ {
		value := decimal.Zero
		if val, ok := hourlyMap[i]; ok {
			value = val
		}
		fullResults[i] = reportdomain.ChartData{Hour: i, Label: "Sales", Value: value}
	}

	return fullResults, nil
}

func (rs *ReportStore) GetMostUsedServices(ctx context.Context, start, end time.Time) ([]reportdomain.ChartData, error) {
	var results []reportdomain.ChartData

	err := rs.db.WithContext(ctx).Model(&accessdomain.AttendanceLog{}).
		Select("s.name as label, COUNT(attendance_log.service_id) as value").
		Joins("JOIN services s ON s.service_id = attendance_log.service_id").
		Where("attendance_log.check_in_time BETWEEN ? AND ?", start, end).
		Group("s.name").
		Order("value DESC").
		Limit(5).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}
