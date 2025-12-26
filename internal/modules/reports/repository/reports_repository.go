package reportrepository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	accessdomain "github.com/vitalfit/api/internal/modules/access/domain"
	authdomain "github.com/vitalfit/api/internal/modules/auth/domain"
	billingdomain "github.com/vitalfit/api/internal/modules/billing/domain"
	branchdomain "github.com/vitalfit/api/internal/modules/branches/domain"
	reportdomain "github.com/vitalfit/api/internal/modules/reports/domain"
	scheduledomain "github.com/vitalfit/api/internal/modules/schedule/domain"
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

func (rs *ReportStore) GetTotalSales(ctx context.Context) (*reportdomain.TotalSalesStats, error) {
	var totalLifetime, currentTotal, prevTotal decimal.Decimal
	validStatuses := []billingdomain.InvoiceStatus{
		billingdomain.InvoiceStatusPaid,
	}

	now := time.Now()
	currentMonthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)

	err := rs.db.WithContext(ctx).Model(&billingdomain.Invoice{}).
		Where("status IN ?", validStatuses).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&totalLifetime).Error
	if err != nil {
		return nil, err
	}

	// Calculate trend based on current month vs last month (using same valid statuses)
	err = rs.db.WithContext(ctx).Model(&billingdomain.Invoice{}).
		Where("issue_date >= ? AND status IN ?", currentMonthStart, validStatuses).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&currentTotal).Error
	if err != nil {
		return nil, err
	}

	err = rs.db.WithContext(ctx).Model(&billingdomain.Invoice{}).
		Where("issue_date >= ? AND issue_date < ? AND status IN ?", prevMonthStart, currentMonthStart, validStatuses).
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

	return &reportdomain.TotalSalesStats{
		TotalSales:       totalLifetime,
		PercentageChange: percentageChange,
		Trend:            trend,
	}, nil
}

func (rs *ReportStore) GetMonthlySalesKPI(ctx context.Context, branchID *uuid.UUID) (*reportdomain.KPICard, error) {
	now := time.Now()
	currentMonthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)

	var currentTotal, prevTotal decimal.Decimal
	validStatus := billingdomain.InvoiceStatusPaid

	// Query Current Month
	queryCurrent := rs.db.WithContext(ctx).Model(&billingdomain.Invoice{}).
		Where("issue_date >= ? AND status = ?", currentMonthStart, validStatus)

	if branchID != nil {
		queryCurrent = queryCurrent.Where("branch_id = ?", *branchID)
	}

	if err := queryCurrent.Select("COALESCE(SUM(total_amount), 0)").Scan(&currentTotal).Error; err != nil {
		return nil, err
	}

	// Query Previous Month
	queryPrev := rs.db.WithContext(ctx).Model(&billingdomain.Invoice{}).
		Where("issue_date >= ? AND issue_date < ? AND status = ?", prevMonthStart, currentMonthStart, validStatus)

	if branchID != nil {
		queryPrev = queryPrev.Where("branch_id = ?", *branchID)
	}

	if err := queryPrev.Select("COALESCE(SUM(total_amount), 0)").Scan(&prevTotal).Error; err != nil {
		return nil, err
	}

	// Calculate Variation
	percentageChange := 0.0
	if !prevTotal.IsZero() {
		diff := currentTotal.Sub(prevTotal)
		percentageChange, _ = diff.Div(prevTotal).Mul(decimal.NewFromInt(100)).Float64()
	} else if !currentTotal.IsZero() {
		percentageChange = 100.0
	}

	return &reportdomain.KPICard{
		Title:        "Total Sales (Month)",
		Value:        currentTotal,
		TrendPercent: percentageChange,
		TrendLabel:   "vs previous month",
		IsPositive:   percentageChange >= 0,
	}, nil
}

func (rs *ReportStore) GetActiveMembersKPI(ctx context.Context, branchID *uuid.UUID) (*reportdomain.KPICard, error) {
	now := time.Now()
	thirtyDaysAgo := now.AddDate(0, 0, -30)
	sixtyDaysAgo := now.AddDate(0, 0, -60)

	var currentCount, prevCount int64

	buildQuery := func(start, end time.Time) *gorm.DB {
		query := rs.db.WithContext(ctx).Table("attendance_log al").
			Where("al.check_in_time >= ? AND al.check_in_time <= ?", start, end)

		if branchID != nil {
			query = query.Joins("JOIN classes c ON c.class_id = al.schedule_id").
				Where("c.branch_id = ?", *branchID)
		}
		return query
	}

	if err := buildQuery(thirtyDaysAgo, now).Distinct("al.user_id").Count(&currentCount).Error; err != nil {
		return nil, err
	}

	if err := buildQuery(sixtyDaysAgo, thirtyDaysAgo).Distinct("al.user_id").Count(&prevCount).Error; err != nil {
		return nil, err
	}

	percentageChange := 0.0
	if prevCount > 0 {
		percentageChange = float64(currentCount-prevCount) / float64(prevCount) * 100
	} else if currentCount > 0 {
		percentageChange = 100.0
	}

	return &reportdomain.KPICard{
		Title:        "Active Members",
		Value:        decimal.NewFromInt(currentCount),
		TrendPercent: percentageChange,
		TrendLabel:   "vs. the previous 30 days",
		IsPositive:   percentageChange >= 0,
	}, nil
}

func (rs *ReportStore) GetOccupancyKPI(ctx context.Context, branchID *uuid.UUID) (*reportdomain.KPICard, error) {
	var maxCapacity int64
	if branchID != nil {
		err := rs.db.WithContext(ctx).Model(&branchdomain.Branch{}).
			Where("branch_id = ?", *branchID).
			Select("COALESCE(max_capacity, 0)").Scan(&maxCapacity).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := rs.db.WithContext(ctx).Model(&branchdomain.Branch{}).
			Where("status = ?", branchdomain.BranchStatusActive).
			Select("COALESCE(SUM(max_capacity), 0)").Scan(&maxCapacity).Error
		if err != nil {
			return nil, err
		}
	}

	if maxCapacity == 0 {
		return &reportdomain.KPICard{
			Title:        "Average Occupation",
			Value:        decimal.Zero,
			TrendPercent: 0,
			TrendLabel:   "Undefined capacity",
		}, nil
	}

	now := time.Now()
	currentMonthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	daysInCurrent := float64(now.Day())

	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)
	prevMonthEnd := currentMonthStart.Add(-time.Nanosecond)
	daysInPrev := float64(prevMonthEnd.Day())

	countAttendance := func(start, end time.Time) (int64, error) {
		var count int64
		query := rs.db.WithContext(ctx).Table("attendance_log al").
			Where("al.check_in_time >= ? AND al.check_in_time <= ?", start, end)

		if branchID != nil {
			query = query.Joins("JOIN classes c ON c.class_id = al.schedule_id").
				Where("c.branch_id = ?", *branchID)
		}
		err := query.Count(&count).Error
		return count, err
	}

	currentCount, err := countAttendance(currentMonthStart, now)
	if err != nil {
		return nil, err
	}
	avgDailyCurrent := float64(currentCount) / daysInCurrent
	occupancyCurrent := (avgDailyCurrent / float64(maxCapacity)) * 100

	prevCount, err := countAttendance(prevMonthStart, prevMonthEnd)
	if err != nil {
		return nil, err
	}
	avgDailyPrev := float64(prevCount) / daysInPrev
	occupancyPrev := (avgDailyPrev / float64(maxCapacity)) * 100

	trendPercent := occupancyCurrent - occupancyPrev

	return &reportdomain.KPICard{
		Title:        "Average Occupation",
		Value:        decimal.NewFromFloat(occupancyCurrent).Round(2),
		TrendPercent: decimal.NewFromFloat(trendPercent).Round(2).InexactFloat64(),
		TrendLabel:   "vs previous month (pts %)",
		IsPositive:   trendPercent >= 0,
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
			Label:         res.BranchName,
			Value:         res.CurrentTotal,
			PercentChange: percentChange,
			Status:        label,
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

func (rs *ReportStore) GetWeeklySalesChart(ctx context.Context, branchID *uuid.UUID) ([]reportdomain.ChartData, error) {
	now := time.Now()

	// Calcular inicio de semana (Lunes)
	offset := int(now.Weekday())
	if offset == 0 {
		offset = 7
	}
	offset-- // Ajustar para que Lunes sea 0 offset desde el inicio
	startOfWeek := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -offset)
	endOfWeek := startOfWeek.AddDate(0, 0, 7).Add(-time.Nanosecond)

	validStatuses := []billingdomain.InvoiceStatus{
		billingdomain.InvoiceStatusPaid,
		billingdomain.InvoiceStatusUnpaid,
		billingdomain.InvoiceStatusOverdue,
	}

	type dailyResult struct {
		Date  time.Time       `gorm:"column:date"`
		Total decimal.Decimal `gorm:"column:total"`
	}

	var queryResults []dailyResult

	query := rs.db.WithContext(ctx).Model(&billingdomain.Invoice{}).
		Select("DATE(issue_date) as date, SUM(total_amount) as total").
		Where("issue_date BETWEEN ? AND ?", startOfWeek, endOfWeek).
		Where("status IN ?", validStatuses)

	if branchID != nil {
		query = query.Where("branch_id = ?", *branchID)
	}

	err := query.Group("DATE(issue_date)").Scan(&queryResults).Error
	if err != nil {
		return nil, err
	}

	// Mapear resultados a un mapa para acceso rápido
	salesMap := make(map[string]decimal.Decimal)
	for _, r := range queryResults {
		salesMap[r.Date.Format("2006-01-02")] = r.Total
	}

	// Construir respuesta completa para los 7 días (Lunes a Domingo)
	var chartData []reportdomain.ChartData
	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

	for i := 0; i < 7; i++ {
		currentDay := startOfWeek.AddDate(0, 0, i)
		val := salesMap[currentDay.Format("2006-01-02")] // Será 0 si no existe (decimal.Decimal zero value)
		chartData = append(chartData, reportdomain.ChartData{Label: days[i], Value: val})
	}

	return chartData, nil
}

func (rs *ReportStore) GetActivityHeatmap(ctx context.Context, branchID *uuid.UUID) ([]reportdomain.HeatmapPoint, error) {
	var results []reportdomain.HeatmapPoint

	start := time.Now().AddDate(0, 0, -30)

	query := rs.db.WithContext(ctx).Table("attendance_log al").
		Select("EXTRACT(ISODOW FROM al.check_in_time) as day_of_week, FLOOR(EXTRACT(HOUR FROM al.check_in_time) / 3) * 3 as hour, COUNT(*) as value").
		Where("al.check_in_time >= ?", start)

	if branchID != nil {
		query = query.Joins("JOIN classes c ON c.class_id = al.schedule_id").
			Where("c.branch_id = ?", *branchID)
	}

	if err := query.Group("day_of_week, hour").Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (rs *ReportStore) GetClassOccupancyChart(ctx context.Context, branchID *uuid.UUID) ([]reportdomain.ChartData, error) {
	var results []reportdomain.ChartData

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	query := rs.db.WithContext(ctx).Table("classes c").
		Select("sc.name as label, ROUND(AVG((CAST((SELECT COUNT(*) FROM attendance_log al WHERE al.schedule_id = c.class_id AND al.status = 'Attended') AS DECIMAL) / NULLIF(c.max_capacity, 0)) * 100), 2) as value").
		Joins("JOIN services s ON s.service_id = c.service_id").
		Joins("JOIN service_categories sc ON sc.category_id = s.category_id").
		Where("c.starts_at BETWEEN ? AND ?", startOfMonth, endOfMonth)

	if branchID != nil {
		query = query.Where("c.branch_id = ?", *branchID)
	}

	err := query.Group("sc.name").Order("value DESC").Scan(&results).Error
	return results, err
}

func (rs *ReportStore) GetFinancialSummary(ctx context.Context, branchID *uuid.UUID) (*reportdomain.FinancialSummary, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	type result struct {
		Category string          `gorm:"column:category"`
		Amount   decimal.Decimal `gorm:"column:amount"`
	}

	var queryResults []result

	query := rs.db.WithContext(ctx).Table("invoice_items ii").
		Select(`
			CASE 
				WHEN ii.membership_type_id IS NOT NULL THEN 'Memberships'
				WHEN ii.service_id IS NOT NULL THEN 'Services'
				WHEN ii.package_id IS NOT NULL THEN 'Combos'
				ELSE 'Other'
			END as category,
			COALESCE(SUM(ii.total_line), 0) as amount
		`).
		Joins("JOIN invoices i ON i.invoice_id = ii.invoice_id").
		Where("i.status = ?", billingdomain.InvoiceStatusPaid).
		Where("i.issue_date BETWEEN ? AND ?", startOfMonth, endOfMonth)

	if branchID != nil {
		query = query.Where("i.branch_id = ?", *branchID)
	}

	if err := query.Group("category").Scan(&queryResults).Error; err != nil {
		return nil, err
	}

	// Initialize map with fixed categories to ensure they always appear
	summaryMap := map[string]decimal.Decimal{
		"Memberships": decimal.Zero,
		"Services":    decimal.Zero,
		"Products":    decimal.Zero,
	}

	total := decimal.Zero
	for _, r := range queryResults {
		summaryMap[r.Category] = r.Amount
		total = total.Add(r.Amount)
	}

	var items []reportdomain.FinancialSummaryItem
	// Fixed order for consistency
	order := []string{"Memberships", "Services", "Combos"}
	for _, cat := range order {
		items = append(items, reportdomain.FinancialSummaryItem{Category: cat, Amount: summaryMap[cat]})
	}

	return &reportdomain.FinancialSummary{Items: items, Total: total}, nil
}

func (rs *ReportStore) GetTodayCheckInsStat(ctx context.Context, branchID *uuid.UUID) (int64, error) {
	var count int64
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.AddDate(0, 0, 1).Add(-time.Nanosecond)

	query := rs.db.WithContext(ctx).Table("attendance_log al").
		Where("al.check_in_time >= ? AND al.check_in_time <= ?", startOfDay, endOfDay)

	if branchID != nil {
		query = query.Joins("JOIN classes c ON c.class_id = al.schedule_id").
			Where("c.branch_id = ?", *branchID)
	}

	err := query.Count(&count).Error
	return count, err
}

func (rs *ReportStore) GetCurrentOccupancyStat(ctx context.Context, branchID *uuid.UUID) (decimal.Decimal, error) {
	var maxCapacity int64
	if branchID != nil {
		err := rs.db.WithContext(ctx).Model(&branchdomain.Branch{}).
			Where("branch_id = ?", *branchID).
			Select("COALESCE(max_capacity, 0)").Scan(&maxCapacity).Error
		if err != nil {
			return decimal.Zero, err
		}
	} else {
		err := rs.db.WithContext(ctx).Model(&branchdomain.Branch{}).
			Where("status = ?", branchdomain.BranchStatusActive).
			Select("COALESCE(SUM(max_capacity), 0)").Scan(&maxCapacity).Error
		if err != nil {
			return decimal.Zero, err
		}
	}

	if maxCapacity == 0 {
		return decimal.Zero, nil
	}

	var count int64
	now := time.Now()
	startWindow := now.Add(-3 * time.Hour)
	endWindow := now.Add(1 * time.Hour)

	query := rs.db.WithContext(ctx).Table("attendance_log al").
		Where("al.check_in_time >= ? AND al.check_in_time <= ?", startWindow, endWindow)

	if branchID != nil {
		query = query.Joins("JOIN classes c ON c.class_id = al.schedule_id").
			Where("c.branch_id = ?", *branchID)
	}

	if err := query.Count(&count).Error; err != nil {
		return decimal.Zero, err
	}

	occupancy := (float64(count) / float64(maxCapacity)) * 100
	return decimal.NewFromFloat(occupancy).Round(2), nil
}

func (rs *ReportStore) GetClassCapacityRatio(ctx context.Context, classID uuid.UUID) (*reportdomain.ClassCapacityStats, error) {
	var class scheduledomain.Class
	// Get class info (capacity and service name)
	if err := rs.db.WithContext(ctx).Joins("Service").First(&class, "class_id = ?", classID).Error; err != nil {
		return nil, err
	}

	var count int64
	// Count attendees
	if err := rs.db.WithContext(ctx).Model(&accessdomain.AttendanceLog{}).
		Where("schedule_id = ?", classID).
		Where("status = ?", accessdomain.AttendanceStatusAttended).
		Count(&count).Error; err != nil {
		return nil, err
	}

	return &reportdomain.ClassCapacityStats{
		ClassName:    class.Service.Name,
		CurrentCount: count,
		MaxCapacity:  class.MaxCapacity,
		Ratio:        fmt.Sprintf("%d / %d", count, class.MaxCapacity),
	}, nil
}

func (rs *ReportStore) GetUpcomingClassesToday(ctx context.Context, branchID *uuid.UUID) ([]reportdomain.ClassScheduleItem, error) {
	var results []reportdomain.ClassScheduleItem
	now := time.Now()
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())

	query := rs.db.WithContext(ctx).Table("classes c").
		Select("c.class_id, s.name as class_name, u.first_name || ' ' || u.last_name as instructor_name, c.starts_at as start_time, c.ends_at as end_time, c.max_capacity").
		Joins("JOIN services s ON s.service_id = c.service_id").
		Joins("JOIN instructors i ON i.instructor_id = c.instructor_id").
		Joins("JOIN users u ON u.user_id = i.user_id").
		Where("c.ends_at > ? AND c.starts_at <= ?", now, endOfDay)

	if branchID != nil {
		query = query.Where("c.branch_id = ?", *branchID)
	}

	err := query.Order("c.starts_at ASC").Scan(&results).Error
	return results, err
}

func (rs *ReportStore) GetRecentCheckIns(ctx context.Context, branchID *uuid.UUID) ([]reportdomain.RecentAttendanceItem, error) {
	var results []reportdomain.RecentAttendanceItem

	query := rs.db.WithContext(ctx).Table("attendance_log al").
		Select("u.first_name || ' ' || u.last_name as user_name, al.check_in_time, s.name as service_name").
		Joins("JOIN users u ON u.user_id = al.user_id").
		Joins("JOIN services s ON s.service_id = al.service_id").
		Where("al.status = ?", accessdomain.AttendanceStatusAttended)

	if branchID != nil {
		query = query.Joins("JOIN classes c ON c.class_id = al.schedule_id").
			Where("c.branch_id = ?", *branchID)
	}

	err := query.Order("al.check_in_time DESC").Limit(4).Scan(&results).Error
	return results, err
}
