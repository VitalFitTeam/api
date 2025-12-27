package reportrepository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	billingdomain "github.com/vitalfit/api/internal/modules/billing/domain"
	reportdomain "github.com/vitalfit/api/internal/modules/reports/domain"
)

func (rs *ReportStore) GetWeeklyRevenueKPI(ctx context.Context, branchID *uuid.UUID) (*reportdomain.KPICard, error) {
	now := time.Now()

	// Calcular inicio de semana actual (Lunes)
	offset := int(now.Weekday())
	if offset == 0 {
		offset = 7
	}
	offset--
	startOfWeek := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -offset)
	startOfPrevWeek := startOfWeek.AddDate(0, 0, -7)

	var currentTotal, prevTotal decimal.Decimal

	// Query Current Week
	queryCurrent := rs.db.WithContext(ctx).Table("payments p").
		Joins("JOIN invoices i ON i.invoice_id = p.invoice_id").
		Where("p.payment_date >= ?", startOfWeek).
		Where("p.status = ?", billingdomain.PaymentStatusCompleted)

	if branchID != nil {
		queryCurrent = queryCurrent.Where("i.branch_id = ?", *branchID)
	}

	if err := queryCurrent.Select("COALESCE(SUM(p.amount_base), 0)").Scan(&currentTotal).Error; err != nil {
		return nil, err
	}

	// Query Previous Week
	queryPrev := rs.db.WithContext(ctx).Table("payments p").
		Joins("JOIN invoices i ON i.invoice_id = p.invoice_id").
		Where("p.payment_date >= ? AND p.payment_date < ?", startOfPrevWeek, startOfWeek).
		Where("p.status = ?", billingdomain.PaymentStatusCompleted)

	if branchID != nil {
		queryPrev = queryPrev.Where("i.branch_id = ?", *branchID)
	}

	if err := queryPrev.Select("COALESCE(SUM(p.amount_base), 0)").Scan(&prevTotal).Error; err != nil {
		return nil, err
	}

	return calculateTrend(currentTotal, prevTotal, "Weekly Revenue", "vs previous week")
}

func (rs *ReportStore) GetAverageTicketKPI(ctx context.Context, branchID *uuid.UUID) (*reportdomain.KPICard, error) {
	now := time.Now()
	currentMonthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)

	calcAvg := func(start time.Time, end *time.Time) (decimal.Decimal, error) {
		var result struct {
			Total decimal.Decimal
			Count int64
		}

		query := rs.db.WithContext(ctx).Model(&billingdomain.Invoice{}).
			Where("issue_date >= ? AND status = ?", start, billingdomain.InvoiceStatusPaid)

		if end != nil {
			query = query.Where("issue_date < ?", end)
		}

		if branchID != nil {
			query = query.Where("branch_id = ?", *branchID)
		}

		err := query.Select("COALESCE(SUM(total_amount), 0) as total, COUNT(*) as count").Scan(&result).Error
		if err != nil {
			return decimal.Zero, err
		}

		if result.Count == 0 {
			return decimal.Zero, nil
		}
		return result.Total.Div(decimal.NewFromInt(result.Count)), nil
	}

	currentAvg, err := calcAvg(currentMonthStart, nil)
	if err != nil {
		return nil, err
	}

	prevAvg, err := calcAvg(prevMonthStart, &currentMonthStart)
	if err != nil {
		return nil, err
	}

	return calculateTrend(currentAvg, prevAvg, "Average Ticket", "vs previous month")
}

func (rs *ReportStore) GetAccountsReceivableKPI(ctx context.Context, branchID *uuid.UUID) (*reportdomain.KPICard, error) {
	var totalDebt decimal.Decimal

	// Suma (Total Factura - Total Pagado) para facturas pendientes de usuarios activos
	query := rs.db.WithContext(ctx).Table("invoices i").
		Select("COALESCE(SUM(i.total_amount - COALESCE(paid_sum.paid, 0)), 0)").
		Joins("JOIN users u ON u.user_id = i.user_id").
		Joins(`LEFT JOIN (
			SELECT invoice_id, SUM(amount_base) as paid 
			FROM payments 
			WHERE status = ? 
			GROUP BY invoice_id
		) paid_sum ON paid_sum.invoice_id = i.invoice_id`, billingdomain.PaymentStatusCompleted).
		Where("i.status IN ?", []billingdomain.InvoiceStatus{billingdomain.InvoiceStatusUnpaid, billingdomain.InvoiceStatusOverdue}).
		Where("u.deleted_at IS NULL")

	if branchID != nil {
		query = query.Where("i.branch_id = ?", *branchID)
	}

	if err := query.Scan(&totalDebt).Error; err != nil {
		return nil, err
	}

	return &reportdomain.KPICard{
		Title:      "Accounts Receivable",
		Value:      totalDebt,
		TrendLabel: "Total Outstanding Debt",
	}, nil
}

func (rs *ReportStore) GetMonthlyRecurringRevenueKPI(ctx context.Context, branchID *uuid.UUID) (*reportdomain.KPICard, error) {
	now := time.Now()
	currentMonthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)

	calcMRR := func(start time.Time, end *time.Time) (decimal.Decimal, error) {
		var total decimal.Decimal
		query := rs.db.WithContext(ctx).Table("invoice_items ii").
			Joins("JOIN invoices i ON i.invoice_id = ii.invoice_id").
			Where("i.issue_date >= ? AND i.status = ?", start, billingdomain.InvoiceStatusPaid).
			Where("ii.membership_type_id IS NOT NULL") // Filtro clave para MRR: Solo membresías

		if end != nil {
			query = query.Where("i.issue_date < ?", end)
		}

		if branchID != nil {
			query = query.Where("i.branch_id = ?", *branchID)
		}

		err := query.Select("COALESCE(SUM(ii.total_line), 0)").Scan(&total).Error
		return total, err
	}

	currentMRR, err := calcMRR(currentMonthStart, nil)
	if err != nil {
		return nil, err
	}

	prevMRR, err := calcMRR(prevMonthStart, &currentMonthStart)
	if err != nil {
		return nil, err
	}

	return calculateTrend(currentMRR, prevMRR, "MRR (Memberships)", "vs previous month")
}

func (rs *ReportStore) GetBillingByBranchMatrix(ctx context.Context, start, end time.Time) (*reportdomain.BillingMatrix, error) {
	type queryResult struct {
		BranchName string          `gorm:"column:branch_name"`
		Category   string          `gorm:"column:category"`
		Amount     decimal.Decimal `gorm:"column:amount"`
	}

	var results []queryResult

	err := rs.db.WithContext(ctx).Table("invoice_items ii").
		Select(`
			b.name as branch_name,
			CASE 
				WHEN ii.membership_type_id IS NOT NULL THEN 'Memberships'
				WHEN ii.service_id IS NOT NULL THEN 'Services'
				WHEN ii.package_id IS NOT NULL THEN 'Combos'
				ELSE 'Other'
			END as category,
			COALESCE(SUM(ii.total_line), 0) as amount
		`).
		Joins("JOIN invoices i ON i.invoice_id = ii.invoice_id").
		Joins("JOIN branch b ON b.branch_id = i.branch_id").
		Where("i.status = ?", billingdomain.InvoiceStatusPaid).
		Where("i.issue_date BETWEEN ? AND ?", start, end).
		Group("b.name, category").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// Initialize structures
	matrix := &reportdomain.BillingMatrix{
		Branches:   []string{},
		Rows:       []reportdomain.BillingMatrixRow{},
		Totals:     make(map[string]decimal.Decimal),
		GrandTotal: decimal.Zero,
	}

	// Concepts are fixed to ensure order
	concepts := []string{"Memberships", "Services", "Combos", "Other"}
	rowsMap := make(map[string]*reportdomain.BillingMatrixRow)
	branchesSet := make(map[string]bool)

	for _, c := range concepts {
		rowsMap[c] = &reportdomain.BillingMatrixRow{
			Concept: c,
			Values:  make(map[string]decimal.Decimal),
			Total:   decimal.Zero,
		}
	}

	for _, r := range results {
		if !branchesSet[r.BranchName] {
			branchesSet[r.BranchName] = true
			matrix.Branches = append(matrix.Branches, r.BranchName)
			matrix.Totals[r.BranchName] = decimal.Zero
		}

		row := rowsMap[r.Category]
		row.Values[r.BranchName] = r.Amount
		row.Total = row.Total.Add(r.Amount)

		// Update vertical totals
		matrix.Totals[r.BranchName] = matrix.Totals[r.BranchName].Add(r.Amount)
		matrix.GrandTotal = matrix.GrandTotal.Add(r.Amount)
	}

	// Convert map to slice in order
	for _, c := range concepts {
		matrix.Rows = append(matrix.Rows, *rowsMap[c])
	}

	return matrix, nil
}

func calculateTrend(current, prev decimal.Decimal, title, label string) (*reportdomain.KPICard, error) {
	percentageChange := 0.0
	if !prev.IsZero() {
		diff := current.Sub(prev)
		percentageChange, _ = diff.Div(prev).Mul(decimal.NewFromInt(100)).Float64()
	} else if !current.IsZero() {
		percentageChange = 100.0
	}

	return &reportdomain.KPICard{
		Title:        title,
		Value:        current,
		TrendPercent: percentageChange,
		TrendLabel:   label,
		IsPositive:   percentageChange >= 0,
	}, nil
}
