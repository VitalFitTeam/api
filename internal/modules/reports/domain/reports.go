package reportdomain

import "github.com/shopspring/decimal"

type GlobalSalesStats struct {
	TotalCurrentMonth decimal.Decimal `json:"total_current_month"`
	TotalLastMonth    decimal.Decimal `json:"total_last_month"`
	PercentageChange  float64         `json:"percentage_change"`
	Trend             string          `json:"trend"`
}

type BranchPerformance struct {
	BranchName    string          `json:"branch_name"`
	TotalSales    decimal.Decimal `json:"total_sales"`
	PercentChange float64         `json:"percent_change"`
	Label         string          `json:"label"`
	Trend         string          `json:"trend"`
}

type ChartData struct {
	Label string          `json:"label"`
	Value decimal.Decimal `json:"value"`
	Hour  int             `json:"hour,omitempty"`
}
