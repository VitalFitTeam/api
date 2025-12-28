package reporthandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary		Get Weekly Revenue KPI
// @Description	Retrieves total revenue (sum of completed payments) for the current week compared to the previous week.
// @Tags			Reports Financial
// @Security		ApiKeyAuth
// @Produce		json
// @Param			branch_id	query		string								false	"Filter by Branch UUID"
// @Success		200			{object}	object{data=reportdomain.KPICard}	"Weekly revenue KPI"
// @Failure		500			{object}	object{error=string}				"Internal Server Error"
// @Router			/reports/kpi/weekly-revenue [get]
func (h *ReportHanlders) GetWeeklyRevenueKPIHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var branchID *uuid.UUID
	if idStr := c.Query("branch_id"); idStr != "" {
		if id, err := uuid.Parse(idStr); err == nil {
			branchID = &id
		}
	}
	kpi, err := h.services.ReportServices.GetWeeklyRevenueKPI(ctx, branchID)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": kpi})
}

// @Summary		Get Monthly Recurring Revenue (MRR) KPI
// @Description	Retrieves the sum of revenue generated solely from memberships (excluding one-time products/services) for the current month.
// @Tags			Reports Financial
// @Security		ApiKeyAuth
// @Produce		json
// @Param			branch_id	query		string								false	"Filter by Branch UUID"
// @Success		200			{object}	object{data=reportdomain.KPICard}	"MRR KPI"
// @Failure		500			{object}	object{error=string}				"Internal Server Error"
// @Router			/reports/kpi/mrr [get]
func (h *ReportHanlders) GetMonthlyRecurringRevenueKPIHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var branchID *uuid.UUID
	if idStr := c.Query("branch_id"); idStr != "" {
		if id, err := uuid.Parse(idStr); err == nil {
			branchID = &id
		}
	}
	kpi, err := h.services.ReportServices.GetMonthlyRecurringRevenueKPI(ctx, branchID)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": kpi})
}

// @Summary		Get Accounts Receivable KPI
// @Description	Retrieves the total outstanding debt (Accounts Receivable) from unpaid or overdue invoices of active users.
// @Tags			Reports Financial
// @Security		ApiKeyAuth
// @Produce		json
// @Param			branch_id	query		string								false	"Filter by Branch UUID"
// @Success		200			{object}	object{data=reportdomain.KPICard}	"Accounts Receivable KPI"
// @Failure		500			{object}	object{error=string}				"Internal Server Error"
// @Router			/reports/kpi/accounts-receivable [get]
func (h *ReportHanlders) GetAccountsReceivableKPIHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var branchID *uuid.UUID
	if idStr := c.Query("branch_id"); idStr != "" {
		if id, err := uuid.Parse(idStr); err == nil {
			branchID = &id
		}
	}
	kpi, err := h.services.ReportServices.GetAccountsReceivableKPI(ctx, branchID)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": kpi})
}

// @Summary		Get Average Ticket KPI
// @Description	Retrieves the average transaction value (Total Revenue / Number of Transactions) for the current month compared to the previous month.
// @Tags			Reports Financial
// @Security		ApiKeyAuth
// @Produce		json
// @Param			branch_id	query		string								false	"Filter by Branch UUID"
// @Success		200			{object}	object{data=reportdomain.KPICard}	"Average Ticket KPI"
// @Failure		500			{object}	object{error=string}				"Internal Server Error"
// @Router			/reports/kpi/average-ticket [get]
func (h *ReportHanlders) GetAverageTicketKPIHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var branchID *uuid.UUID
	if idStr := c.Query("branch_id"); idStr != "" {
		if id, err := uuid.Parse(idStr); err == nil {
			branchID = &id
		}
	}
	kpi, err := h.services.ReportServices.GetAverageTicketKPI(ctx, branchID)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": kpi})
}

// @Summary		Get Monthly Revenue Chart
// @Description	Retrieves total revenue per month for the current year.
// @Tags			Reports Financial
// @Security		ApiKeyAuth
// @Produce		json
// @Param			branch_id	query		string									false	"Filter by Branch UUID"
// @Success		200			{object}	object{data=[]reportdomain.ChartData}	"Monthly revenue data"
// @Failure		500			{object}	object{error=string}					"Internal Server Error"
// @Router			/reports/charts/monthly-revenue [get]
func (h *ReportHanlders) GetMonthlyRevenueChartHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var branchID *uuid.UUID
	if idStr := c.Query("branch_id"); idStr != "" {
		if id, err := uuid.Parse(idStr); err == nil {
			branchID = &id
		}
	}
	data, err := h.services.ReportServices.GetMonthlyRevenueChart(ctx, branchID)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

// @Summary		Get Billing Detail by Branch Matrix
// @Description	Retrieves a comparative matrix of revenue by concept (rows) and branch (columns).
// @Tags			Reports Financial
// @Security		ApiKeyAuth
// @Produce		json
// @Param			start	query		string									false	"Start date (YYYY-MM-DD)"
// @Param			end		query		string									false	"End date (YYYY-MM-DD)"
// @Success		200		{object}	object{data=reportdomain.BillingMatrix}	"Billing matrix data"
// @Failure		500		{object}	object{error=string}					"Internal Server Error"
// @Router			/reports/kpi/billing-matrix [get]
func (h *ReportHanlders) GetBillingByBranchMatrixHandler(c *gin.Context) {
	ctx := c.Request.Context()
	start, end := parseTimeRange(c)
	matrix, err := h.services.ReportServices.GetBillingByBranchMatrix(ctx, start, end)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": matrix})
}

// @Summary		Get Total Transactions Stats
// @Description	Retrieves the total lifetime count of completed transactions (payments). Supports filtering by branch. Includes a trend comparison (current month vs previous month).
// @Tags			Reports Financial
// @Security		ApiKeyAuth
// @Produce		json
// @Param			branch_id	query		string								false	"Filter by Branch UUID"
// @Success		200			{object}	object{data=reportdomain.KPICard}	"Total transactions KPI"
// @Failure		500			{object}	object{error=string}				"Internal Server Error"
// @Router			/reports/kpi/total-transactions [get]
func (h *ReportHanlders) GetTotalTransactionsHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var branchID *uuid.UUID
	if idStr := c.Query("branch_id"); idStr != "" {
		if id, err := uuid.Parse(idStr); err == nil {
			branchID = &id
		}
	}
	kpi, err := h.services.ReportServices.GetTotalTransactions(ctx, branchID)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": kpi})
}

// @Summary		Get Monthly Cash Flow Chart
// @Description	Retrieves total actual income (cash flow) from completed payments per month for the current year.
// @Tags			Reports Financial
// @Security		ApiKeyAuth
// @Produce		json
// @Param			branch_id	query		string									false	"Filter by Branch UUID"
// @Success		200			{object}	object{data=[]reportdomain.ChartData}	"Monthly cash flow data"
// @Failure		500			{object}	object{error=string}					"Internal Server Error"
// @Router			/reports/charts/cash-flow [get]
func (h *ReportHanlders) GetMonthlyCashFlowChartHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var branchID *uuid.UUID
	if idStr := c.Query("branch_id"); idStr != "" {
		if id, err := uuid.Parse(idStr); err == nil {
			branchID = &id
		}
	}
	data, err := h.services.ReportServices.GetMonthlyCashFlowChart(ctx, branchID)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

// @Summary		Get Average Customer Lifetime Value (CLV) KPI
// @Description	Retrieves the average revenue generated per customer (Historical Revenue / Unique Paying Customers).
// @Tags			Reports Financial
// @Security		ApiKeyAuth
// @Produce		json
// @Param			branch_id	query		string								false	"Filter by Branch UUID"
// @Success		200			{object}	object{data=reportdomain.KPICard}	"Average CLV KPI"
// @Failure		500			{object}	object{error=string}				"Internal Server Error"
// @Router			/reports/kpi/clv [get]
func (h *ReportHanlders) GetAverageCLVKPIHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var branchID *uuid.UUID
	if idStr := c.Query("branch_id"); idStr != "" {
		if id, err := uuid.Parse(idStr); err == nil {
			branchID = &id
		}
	}
	kpi, err := h.services.ReportServices.GetAverageCLVKPI(ctx, branchID)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": kpi})
}
