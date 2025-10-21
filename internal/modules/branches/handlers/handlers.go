package branchhandlers

import (
	"github.com/gin-gonic/gin"
	appservices "github.com/vitalfit/api/internal/app/services"
	branchdomain "github.com/vitalfit/api/internal/modules/branches/domain"
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
	c.JSON(201, gin.H{
		"message": "branch created",
	})

}

// @Summary		Get a branch list
// @Description	Get a branch list paginated with optional filters
// @Tags			Branches
// @Security		ApiKeyAuth
// @Param			limit	query	int		false	"limit results per page"					default(10)			minimum(1)	maximum(100)
// @Param			offset	query	int		false	"number of results to skip (paginación)"	default(0)			minimum(0)
// @Param			sort	query	string	false	"clasification order (asc or desc)"			enums(asc, desc)	default(desc)
// @Param			search	query	string	false	"search terms (filter by name)"
// @Param			status	query	string	false	"filter by status"	enums(Active, Inactiv
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
		Sort:   "asc",
		Search: "",
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

	responseList := make([]BranchListResponse, 0, len(branches))

	for _, branch := range branches {

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

	c.JSON(200, gin.H{
		"data":       responseList,
		"pagination": fq,
	})
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
