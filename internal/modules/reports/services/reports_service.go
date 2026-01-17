package reportservices

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	authdomain "github.com/vitalfit/api/internal/modules/auth/domain"
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

func (s *ReportService) GetTotalSales(ctx context.Context) (*reportdomain.TotalSalesStats, error) {
	return s.store.Reports.GetTotalSales(ctx)
}

func (s *ReportService) GetMonthlySalesKPI(ctx context.Context, branchID *uuid.UUID) (*reportdomain.KPICard, error) {
	return s.store.Reports.GetMonthlySalesKPI(ctx, branchID)
}

func (s *ReportService) GetActiveMembersKPI(ctx context.Context, branchID *uuid.UUID) (*reportdomain.KPICard, error) {
	return s.store.Reports.GetActiveMembersKPI(ctx, branchID)
}

func (s *ReportService) GetOccupancyKPI(ctx context.Context, branchID *uuid.UUID) (*reportdomain.KPICard, error) {
	return s.store.Reports.GetOccupancyKPI(ctx, branchID)
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

func (s *ReportService) GetSalesByCategory(ctx context.Context, branchID *uuid.UUID, start, end time.Time) ([]reportdomain.ChartData, error) {
	return s.store.Reports.GetSalesByCategory(ctx, branchID, start, end)
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

func (s *ReportService) GetWeeklySalesChart(ctx context.Context, branchID *uuid.UUID) ([]reportdomain.ChartData, error) {
	return s.store.Reports.GetWeeklySalesChart(ctx, branchID)
}

func (s *ReportService) GetActivityHeatmap(ctx context.Context, branchID *uuid.UUID) ([]reportdomain.HeatmapPoint, error) {
	return s.store.Reports.GetActivityHeatmap(ctx, branchID)
}

func (s *ReportService) GetClassOccupancyChart(ctx context.Context, branchID *uuid.UUID) ([]reportdomain.ChartData, error) {
	return s.store.Reports.GetClassOccupancyChart(ctx, branchID)
}

func (s *ReportService) GetFinancialSummary(ctx context.Context, branchID *uuid.UUID) (*reportdomain.FinancialSummary, error) {
	return s.store.Reports.GetFinancialSummary(ctx, branchID)
}

func (s *ReportService) GetTodayCheckInsStat(ctx context.Context, branchID *uuid.UUID) (int64, error) {
	return s.store.Reports.GetTodayCheckInsStat(ctx, branchID)
}

func (s *ReportService) GetCurrentOccupancyStat(ctx context.Context, branchID *uuid.UUID) (decimal.Decimal, error) {
	return s.store.Reports.GetCurrentOccupancyStat(ctx, branchID)
}

func (s *ReportService) GetClassCapacityRatio(ctx context.Context, classID uuid.UUID) (*reportdomain.ClassCapacityStats, error) {
	return s.store.Reports.GetClassCapacityRatio(ctx, classID)
}

func (s *ReportService) GetUpcomingClassesToday(ctx context.Context, branchID *uuid.UUID) ([]reportdomain.ClassScheduleItem, error) {
	return s.store.Reports.GetUpcomingClassesToday(ctx, branchID)
}

func (s *ReportService) GetRecentCheckIns(ctx context.Context, branchID *uuid.UUID) ([]reportdomain.RecentAttendanceItem, error) {
	return s.store.Reports.GetRecentCheckIns(ctx, branchID)
}

func (s *ReportService) GetInstructorNextClass(ctx context.Context, instructorID uuid.UUID) (string, error) {
	nextClassTime, err := s.store.Reports.GetInstructorNextClass(ctx, instructorID)
	if err != nil {
		return "", err
	}
	if nextClassTime == nil {
		return "Sin pendientes", nil
	}
	return nextClassTime.Format("03:04 PM"), nil
}

func (s *ReportService) GetInstructorIDByUserID(ctx context.Context, userID uuid.UUID) (*uuid.UUID, error) {
	return s.store.Reports.GetInstructorIDByUserID(ctx, userID)
}

func (s *ReportService) GetInstructorStudentCountKPI(ctx context.Context, instructorID uuid.UUID) (*reportdomain.KPICard, error) {
	return s.store.Reports.GetInstructorStudentCountKPI(ctx, instructorID)
}

func (s *ReportService) GetInstructorMonthlyClassesCount(ctx context.Context, instructorID uuid.UUID) (*reportdomain.KPICard, error) {
	return s.store.Reports.GetInstructorMonthlyClassesCount(ctx, instructorID)
}

func (s *ReportService) GetInstructorClassesToday(ctx context.Context, instructorID uuid.UUID) ([]reportdomain.ClassScheduleItem, error) {
	return s.store.Reports.GetInstructorClassesToday(ctx, instructorID)
}

func (s *ReportService) GetWeeklyRevenueKPI(ctx context.Context, branchID *uuid.UUID) (*reportdomain.KPICard, error) {
	return s.store.Reports.GetWeeklyRevenueKPI(ctx, branchID)
}

func (s *ReportService) GetAverageTicketKPI(ctx context.Context, branchID *uuid.UUID) (*reportdomain.KPICard, error) {
	return s.store.Reports.GetAverageTicketKPI(ctx, branchID)
}

func (s *ReportService) GetAccountsReceivableKPI(ctx context.Context, branchID *uuid.UUID) (*reportdomain.KPICard, error) {
	return s.store.Reports.GetAccountsReceivableKPI(ctx, branchID)
}

func (s *ReportService) GetMonthlyRecurringRevenueKPI(ctx context.Context, branchID *uuid.UUID) (*reportdomain.KPICard, error) {
	return s.store.Reports.GetMonthlyRecurringRevenueKPI(ctx, branchID)
}

func (s *ReportService) GetMonthlyRevenueChart(ctx context.Context, branchID *uuid.UUID) ([]reportdomain.ChartData, error) {
	return s.store.Reports.GetMonthlyRevenueChart(ctx, branchID)
}

func (s *ReportService) GetBillingByBranchMatrix(ctx context.Context, start, end time.Time) (*reportdomain.BillingMatrix, error) {
	return s.store.Reports.GetBillingByBranchMatrix(ctx, start, end)
}

func (s *ReportService) GetTotalTransactions(ctx context.Context, branchID *uuid.UUID) (*reportdomain.KPICard, error) {
	return s.store.Reports.GetTotalTransactions(ctx, branchID)
}

func (s *ReportService) GetNewClientsKPI(ctx context.Context, branchID *uuid.UUID) (*reportdomain.KPICard, error) {
	return s.store.Reports.GetNewClientsKPI(ctx, branchID)
}

func (s *ReportService) GetRetentionRateKPI(ctx context.Context, branchID *uuid.UUID) (*reportdomain.KPICard, error) {
	return s.store.Reports.GetRetentionRateKPI(ctx, branchID)
}

func (s *ReportService) GetAverageCLVKPI(ctx context.Context, branchID *uuid.UUID) (*reportdomain.KPICard, error) {
	return s.store.Reports.GetAverageCLVKPI(ctx, branchID)
}

func (s *ReportService) GetNewVsRecurringChart(ctx context.Context, branchID *uuid.UUID) ([]reportdomain.StackedChartData, error) {
	return s.store.Reports.GetNewVsRecurringChart(ctx, branchID)
}

func (s *ReportService) GetCohortAnalysis(ctx context.Context, branchID *uuid.UUID) ([]reportdomain.CohortRetention, error) {
	return s.store.Reports.GetCohortAnalysis(ctx, branchID)
}

func (s *ReportService) GetMonthlyCashFlowChart(ctx context.Context, branchID *uuid.UUID) ([]reportdomain.ChartData, error) {
	return s.store.Reports.GetMonthlyCashFlowChart(ctx, branchID)
}

func (s *ReportService) GetSalesByDemographics(ctx context.Context, branchID *uuid.UUID, start, end time.Time, dimension string) ([]reportdomain.ChartData, error) {
	return s.store.Reports.GetSalesByDemographics(ctx, branchID, start, end, dimension)
}

func (s *ReportService) DetectAndFlagChurnRisk(ctx context.Context) ([]reportdomain.ChurnRiskAnalysis, []reportdomain.ChurnRiskAnalysis, error) {
	// 1. Get historical data for all active clients
	metrics, err := s.store.Reports.GetClientsChurnMetrics(ctx)
	if err != nil {
		return nil, nil, err
	}

	// 2. Get Managers for all branches to identify who to notify
	managersMap, err := s.store.Reports.GetBranchManagers(ctx)
	if err != nil {
		return nil, nil, err
	}

	var newlyDetected []reportdomain.ChurnRiskAnalysis
	var allAtRisk []reportdomain.ChurnRiskAnalysis

	calculateRisk := func(m reportdomain.ClientChurnMetrics) (int, []string) {
		riskScore := 0
		var factors []string

		daysSinceCheckIn := 30.0 // Default high if never checked in
		if m.LastCheckIn != nil {
			daysSinceCheckIn = time.Since(*m.LastCheckIn).Hours() / 24
		}

		if daysSinceCheckIn > 14 {
			riskScore += 40
			factors = append(factors, "Absent > 14 days")
		}

		// B. Trend Factor: Visits this month vs last month
		if m.LastMonthVisits > 0 {
			if float64(m.CurrentMonthVisits) < float64(m.LastMonthVisits)*0.5 {
				riskScore += 30
				factors = append(factors, "Attendance drop > 50%")
			}
		}

		// C. Expiration Factor: Membership expiring soon
		if m.MembershipEndDate != nil {
			daysUntilExpiration := time.Until(*m.MembershipEndDate).Hours() / 24
			if daysUntilExpiration > 0 && daysUntilExpiration < 5 {
				riskScore += 30
				factors = append(factors, "Membership expires < 5 days")
			}
		}
		return riskScore, factors
	}

	for _, m := range metrics {
		riskScore, factors := calculateRisk(m)

		// Action 1: Flag High Risk Users
		if riskScore >= 70 {
			analysis := reportdomain.ChurnRiskAnalysis{
				UserID:            m.UserID,
				Name:              m.FirstName + " " + m.LastName,
				Email:             m.Email,
				Phone:             m.Phone,
				RiskScore:         riskScore,
				Factors:           factors,
				PreferredBranchID: m.PreferredBranchID,
			}

			if m.PreferredBranchID != nil {
				if manager, ok := managersMap[*m.PreferredBranchID]; ok {
					analysis.ManagerID = &manager.UserID
					analysis.ManagerEmail = manager.Email
				}
			}

			// Only update if not already AtRisk to avoid redundant DB writes
			if m.CurrentCategory != string(authdomain.ClientCategoryAtRisk) {
				_ = s.store.Users.UpdateClientCategory(ctx, m.UserID, authdomain.ClientCategoryAtRisk)
				newlyDetected = append(newlyDetected, analysis)
			}

			allAtRisk = append(allAtRisk, analysis)
		} else {
			if m.CurrentCategory == string(authdomain.ClientCategoryAtRisk) {
				newCategory := authdomain.ClientCategoryRegular
				_ = s.store.Users.UpdateClientCategory(ctx, m.UserID, newCategory)
			}
		}
	}

	return newlyDetected, allAtRisk, nil
}

func (s *ReportService) GetChurnRateKPI(ctx context.Context, branchID *uuid.UUID) (*reportdomain.KPICard, error) {
	return s.store.Reports.GetChurnRateKPI(ctx, branchID)
}

func (s *ReportService) GetFinancialReportData(ctx context.Context, branchID *uuid.UUID, start, end time.Time) ([]reportdomain.FinancialReportRow, error) {
	return s.store.Reports.GetFinancialReportData(ctx, branchID, start, end)
}
