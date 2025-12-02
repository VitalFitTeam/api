package reporthandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary		Get Global Sales Stats
// @Description	Retrieves global sales statistics, comparing the current month's sales with the previous month.
// @Tags			Reports
// @Security		ApiKeyAuth
// @Produce		json
// @Success		200	{object}	object{data=reportdomain.GlobalSalesStats}	"Global sales statistics"
// @Failure		500	{object}	object{error=string}						"Internal Server Error"
// @Router			/billing/stats/global [get]
func (h *ReportHanlders) GetGlobalSalesStatsHandler(c *gin.Context) {
	ctx := c.Request.Context()

	stats, err := h.services.ReportServices.GetGlobalSalesStats(ctx)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}
