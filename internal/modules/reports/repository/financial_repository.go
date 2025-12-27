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
