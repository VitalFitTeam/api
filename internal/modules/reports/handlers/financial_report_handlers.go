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
