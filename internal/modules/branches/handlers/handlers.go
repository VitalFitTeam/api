package branchhandlers

import (
	"github.com/gin-gonic/gin"
	appservices "github.com/vitalfit/api/internal/app/services"
	branchdomain "github.com/vitalfit/api/internal/modules/branches/domain"
	"github.com/vitalfit/api/internal/shared/middleware/auth"
)

type BranchHandlersInterface interface {
	BranchRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware)
	CreateBranch(c *gin.Context)
}

type BranchHandlers struct {
	services appservices.Services
}

func NewBranchHandlers(services appservices.Services) *BranchHandlers {
	return &BranchHandlers{services: services}
}

// @Summary		create a new branch
// @Description	adds a new branc to the system with his ubication and configuration
// @Tags			Branches
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			branch	body		CreateBranchPayload		true	"Payload needed for the branch creation"
// @Success		201		{object}	object{message=string}	"suceed message"
// @Failure		400		{object}	object{error=string}	"Error: bad request"
// @Failure		500		{object}	object{error=string}	"Error: internal server error"
// @Router			/branches [post]
func (h *BranchHandlers) CreateBranch(c *gin.Context) {
	var payload CreateBranchPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	branch, err := payload.toBranch()
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	var branchOperatingHours []branchdomain.OperatingHours
	for _, operatingHour := range payload.OperatingHours {
		operatingHour, err := operatingHour.toOperatingHour()
		if err != nil {
			h.services.LogErrors.BadRequestResponse(c, err)
			return
		}
		branchOperatingHours = append(branchOperatingHours, *operatingHour)
	}
	branch.OperatingHours = branchOperatingHours
	state, err := h.services.LocationsServices.FindOrCreateStateByCountry(c, payload.State, payload.Country)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	branch.StateID = state.StateID

	if err := h.services.BranchesServices.CreateBranch(c, branch); err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(201, gin.H{
		"message": "branch created",
	})

}
