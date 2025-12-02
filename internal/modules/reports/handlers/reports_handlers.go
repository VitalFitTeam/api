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
// @Router			/reports/stats/global [get]
func (h *ReportHanlders) GetGlobalSalesStatsHandler(c *gin.Context) {
	ctx := c.Request.Context()

	stats, err := h.services.ReportServices.GetGlobalSalesStats(ctx)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}

// @Summary		Get Top 5 Branches Performance
// @Description	Retrieves the performance of the top 5 branches based on current month's sales.
// @Tags			Reports
// @Security		ApiKeyAuth
// @Produce		json
// @Success		200	{object}	object{data=[]reportdomain.BranchPerformance}	"Top 5 branches performance"
// @Failure		500	{object}	object{error=string}							"Internal Server Error"
// @Router			/reports/stats/top-branches [get]
func (h *ReportHanlders) GetTopBranchesPerformanceHandler(c *gin.Context) {
	ctx := c.Request.Context()

	topBranches, err := h.services.ReportServices.GetTopBranchesPerformance(ctx)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": topBranches})
}

// @Summary		Get Total Clients Stat
// @Description	Retrieves the total number of users with the 'client' role.
// @Tags			Reports
// @Security		ApiKeyAuth
// @Produce		json
// @Success		200	{object}	object{data=int64}		"Total number of clients"
// @Failure		500	{object}	object{error=string}	"Internal Server Error"
// @Router			/reports/stats/total-clients [get]
func (h *ReportHanlders) GetTotalClientsStatHandler(c *gin.Context) {
	ctx := c.Request.Context()
	totalClients, err := h.services.ReportServices.GetTotalClientsStat(ctx)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": totalClients})
}

// @Summary		Get Active Branches Count
// @Description	Retrieves the total number of active branches.
// @Tags			Reports
// @Security		ApiKeyAuth
// @Produce		json
// @Success		200	{object}	object{data=int64}		"Total number of active branches"
// @Failure		500	{object}	object{error=string}	"Internal Server Error"
// @Router			/reports/stats/total-active-branches [get]
func (h *ReportHanlders) GetActiveBranchesCountHandler(c *gin.Context) {
	ctx := c.Request.Context()
	count, err := h.services.ReportServices.GetActiveBranchesCount(ctx)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": count})
}
