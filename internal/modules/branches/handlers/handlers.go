package branchhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	appservices "github.com/vitalfit/api/internal/app/services"
	branchdomain "github.com/vitalfit/api/internal/modules/branches/domain"
	shared_errors "github.com/vitalfit/api/internal/shared/errors"
	"github.com/vitalfit/api/internal/shared/middleware/auth"
	"github.com/vitalfit/api/pkg/pagination"
)

type BranchHandlersInterface interface {
	BranchRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware)
	CreateBranchHandler(c *gin.Context)
	GetPaymentMethodsHandler(c *gin.Context)
	GetBranchesHandler(c *gin.Context)
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
func (h *BranchHandlers) CreateBranchHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var payload CreateBranchPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	user, err := h.services.UserServices.GetByID(ctx, payload.ManagerID)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	if user.Role.Name != "branch_admin" {
		h.services.LogErrors.ForbiddenResponse(c)
		return
	}

	branch, err := payload.toBranch()
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	//operating hours
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
	//payment methods
	var branchPaymentMethods []branchdomain.PaymentMethodsBranch
	for _, paymentMethod := range payload.PaymentMethods {
		branchPaymentMethod := branchdomain.PaymentMethodsBranch{
			MethodID: paymentMethod,
		}
		branchPaymentMethods = append(branchPaymentMethods, branchPaymentMethod)
	}
	branch.PaymentMethodsLinks = branchPaymentMethods
	//state
	state, err := h.services.LocationsServices.FindOrCreateStateByCountry(ctx, payload.State, payload.Country)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	branch.StateID = state.StateID
	//create branch
	if err := h.services.BranchesServices.CreateBranch(ctx, branch); err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "branch created",
	})

}

// @Summary		Get a branch list
// @Description	Get a branch list paginated with optional filters
// @Tags			Branches
// @Security		ApiKeyAuth
// @Param			limit		query	int		false	"limit results per page"					default(10)			minimum(1)	maximum(100)
// @Param			offset		query	int		false	"number of results to skip (paginación)"	default(0)			minimum(0)
// @Param			sort		query	string	false	"clasification order (asc or desc)"			enums(asc, desc)	default(desc)
// @Param			status		query	string	false	"filter by status"							enums(Active, Inactive, Maintenance)
// @Param			search		query	string	false	"global search across branch, state and country names"
// @Param			location	query	string	false	"filter by location (state or country)"
// @Param			tax_id		query	string	false	"filter by exact tax id"
// @Accept			json
// @Produce		json
// @Success		200	{object}	object{data=[]BranchListResponse}	"succed response"
// @Failure		400	{object}	object{error=string}				"Error: invalid params"
// @Failure		500	{object}	object{error=string}				"Error: internal server error"
// @Router			/branches [get]
func (h *BranchHandlers) GetBranchesHandler(c *gin.Context) {
	ctx := c.Request.Context()
	fq := pagination.PaginatedFeedQuery{
		Limit:  10,
		Offset: 0,
		Sort:   "desc",
		Search: "",
		Status: "Active",
	}

	fq, err := fq.Parse(c.Request)
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	branches, err := h.services.BranchesServices.GetBranches(ctx, fq)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	responseList := make([]BranchListResponse, 0, len(branches.Branches))

	for _, branch := range branches.Branches {

		resp := BranchListResponse{
			BranchID:        branch.BranchID,
			Name:            branch.Name,
			TaxID:           branch.TaxID,
			Status:          string(branch.Status),
			StateName:       branch.State.Name,
			CountryName:     branch.State.Country.Name,
			ManagerName:     branch.Manager.FirstName,
			ManagerLastName: branch.Manager.LastName,
		}

		responseList = append(responseList, resp)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": responseList,
		"count": gin.H{
			"active":      branches.ActiveCount,
			"inactive":    branches.InactiveCount,
			"maintenance": branches.ManteinanceCount,
			"total":       branches.ActiveCount + branches.InactiveCount + branches.ManteinanceCount,
		},
		"pagination": fq,
	})
}

// @Summary		Delete a branch
// @Description	Deletes a specific branch by its UUID. Requires "super_admin" role.
// @Tags			Branches
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			id	path		string					true	"Branch UUID"
// @Success		204	{object}	nil						"No Content"
// @Failure		400	{object}	object{error=string}	"Error: Bad Request - Malformed ID"
// @Failure		401	{object}	object{error=string}	"Error: Unauthorized - Token required"
// @Failure		403	{object}	object{error=string}	"Error: Forbidden - Not super_admin"
// @Failure		404	{object}	object{error=string}	"Error: Not Found - Branch not found"
// @Failure		500	{object}	object{error=string}	"Error: Internal Server Error"
// @Router			/branches/{id} [delete]
func (h *BranchHandlers) DeleteBranchHandler(c *gin.Context) {
	idStr := c.Param("id")

	branchID, err := uuid.Parse(idStr)
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	ctx := c.Request.Context()
	if err := h.services.BranchesServices.DeleteBranch(ctx, branchID); err != nil {
		switch err {
		case shared_errors.ErrNotFound:
			h.services.LogErrors.NotFoundResponse(c)
		default:
			h.services.LogErrors.InternalServerError(c, err)
		}
		return

	}

	c.JSON(http.StatusNoContent, nil)
}

// @Summary		Get all payment methods
// @Description	Retrieves a list of all available payment methods in the system.
// @Tags			Branches
// @Security		ApiKeyAuth
// @Produce		json
// @Success		200	{object}	object{data=[]branchdomain.PaymentMethods}	"A list of payment methods"
// @Failure		500	{object}	object{error=string}						"Error: internal server error"
// @Router			/branches/payment-methods [get]
func (h *BranchHandlers) GetPaymentMethodsHandler(c *gin.Context) {
	paymentMethods, err := h.services.BranchesServices.GetPaymentMethods(c)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(200, gin.H{
		"data": paymentMethods,
	})

}
