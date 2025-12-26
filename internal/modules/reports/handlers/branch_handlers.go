package reporthandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary		Get Monthly Sales KPI
// @Description	Retrieves total accumulated sales for the current month and the trend vs previous month.
// @Tags			Reports Branch
// @Security		ApiKeyAuth
// @Produce		json
// @Param			branch_id	query		string								false	"Filter by Branch UUID"
// @Success		200			{object}	object{data=reportdomain.KPICard}	"Monthly sales KPI"
// @Failure		500			{object}	object{error=string}				"Internal Server Error"
// @Router			/reports/kpi/monthly-sales [get]
func (h *ReportHanlders) GetMonthlySalesKPIHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var branchID *uuid.UUID

	if idStr := c.Query("branch_id"); idStr != "" {
		if id, err := uuid.Parse(idStr); err == nil {
			branchID = &id
		}
	}

	kpi, err := h.services.ReportServices.GetMonthlySalesKPI(ctx, branchID)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": kpi})
}
