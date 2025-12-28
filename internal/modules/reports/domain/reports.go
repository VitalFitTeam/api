package reportdomain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type GlobalSalesStats struct {
	TotalCurrentMonth decimal.Decimal `json:"total_current_month"`
	TotalLastMonth    decimal.Decimal `json:"total_last_month"`
	PercentageChange  float64         `json:"percentage_change"`
	Trend             string          `json:"trend"`
}

type TotalSalesStats struct {
	TotalSales       decimal.Decimal `json:"total_sales"`
	PercentageChange float64         `json:"percentage_change"`
	Trend            string          `json:"trend"`
}

type BranchPerformance struct {
	Label         string          `json:"label"`
	Value         decimal.Decimal `json:"value"`
	PercentChange float64         `json:"percent_change"`
	Status        string          `json:"status"`
	Trend         string          `json:"trend"`
}

type ChartData struct {
	Label string          `json:"label"`
	Value decimal.Decimal `json:"value"`
	Hour  int             `json:"hour,omitempty"`
}

type KPICard struct {
	Title        string          `json:"title"`
	Value        decimal.Decimal `json:"value"`
	TrendPercent float64         `json:"trend_percent"`
	TrendLabel   string          `json:"trend_label"`
	IsPositive   bool            `json:"is_positive"`
	Target       decimal.Decimal `json:"target,omitempty"`
}

type HeatmapPoint struct {
	DayOfWeek int   `json:"day_of_week"` // 1=Monday, 7=Sunday
	Hour      int   `json:"hour"`        // Start hour of the bucket (0, 3, 6...)
	Value     int64 `json:"value"`
}

type FinancialSummaryItem struct {
	Category string          `json:"category"`
	Amount   decimal.Decimal `json:"amount"`
}

type FinancialSummary struct {
	Items []FinancialSummaryItem `json:"items"`
	Total decimal.Decimal        `json:"total"`
}

type ClassCapacityStats struct {
	ClassName    string `json:"class_name"`
	CurrentCount int64  `json:"current_count"`
	MaxCapacity  int    `json:"max_capacity"`
	Ratio        string `json:"ratio"` // "X / Y"
}

type ClassScheduleItem struct {
	ClassID        uuid.UUID `json:"class_id"`
	ClassName      string    `json:"class_name"`
	InstructorName string    `json:"instructor_name"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	MaxCapacity    int       `json:"max_capacity"`
}

type RecentAttendanceItem struct {
	UserName    string    `json:"user_name"`
	CheckInTime time.Time `json:"check_in_time"`
	ServiceName string    `json:"service_name"`
}

type BillingMatrixRow struct {
	Concept string                     `json:"concept"`
	Values  map[string]decimal.Decimal `json:"values"` // BranchName -> Amount
	Total   decimal.Decimal            `json:"total"`  // Global Total for this concept
}

type BillingMatrix struct {
	Branches   []string                   `json:"branches"`    // List of branch names (columns)
	Rows       []BillingMatrixRow         `json:"rows"`        // Data rows
	Totals     map[string]decimal.Decimal `json:"totals"`      // BranchName -> Total Vertical
	GrandTotal decimal.Decimal            `json:"grand_total"` // Total of Totals
}

type StackedChartData struct {
	Label     string `json:"label"`     // Month Name
	New       int64  `json:"new"`       // Registered this month
	Recurring int64  `json:"recurring"` // Active this month but registered previously
}
