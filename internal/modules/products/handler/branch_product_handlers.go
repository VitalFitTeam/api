package productshandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	productsdomain "github.com/vitalfit/api/internal/modules/products/domain"
)

// @Summary		Assign services to a branch
// @Description	Assigns a list of services with their specific details to a branch.
// @Tags			Branch Services
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			id			path		string						true	"Branch UUID"
// @Param			services	body		AssignBranchServicePayload	true	"Payload with services to assign"
// @Success		201			{object}	map[string]interface{}		"Branch services assigned successfully"
// @Failure		400			{object}	map[string]interface{}		"Bad Request: Invalid UUID or payload"
// @Failure		500			{object}	map[string]interface{}		"Internal Server Error"
// @Router			/branches/{id}/services [post]
func (h *ProductsHandler) AssignBranchServiceHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var payload AssignBranchServicePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	branchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	branchServices := make([]*productsdomain.ServiceBranchDetail, 0, len(payload.Services))
	for _, service := range payload.Services {
		branchService, err := service.ToServiceBranchDetail()
		if err != nil {
			h.services.LogErrors.BadRequestResponse(c, err)
			return
		}
		branchService.BranchID = branchID
		branchServices = append(branchServices, branchService)
	}
	err = h.services.ProductsServices.AssignBranchService(ctx, branchServices)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Branch services assigned successfully"})

}

// @Summary		List services for a branch
// @Description	Retrieves a list of all services assigned to a specific branch.
// @Tags			Branch Services
// @Security		ApiKeyAuth
// @Produce		json
// @Param			id	path		string												true	"Branch UUID"
// @Success		200	{object}	object{data=[]productsdomain.ServiceBranchDetail}	"List of branch services"
// @Failure		400	{object}	map[string]interface{}								"Bad Request: Invalid UUID format"
// @Failure		404	{object}	map[string]interface{}								"Not Found: Branch not found"
// @Failure		500	{object}	map[string]interface{}								"Internal Server Error"
// @Router			/branches/{id}/services [get]
func (h *ProductsHandler) GetBranchServiceHandler(c *gin.Context) {
	ctx := c.Request.Context()
	branchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	services, err := h.services.ProductsServices.GetBranchService(ctx, branchID)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": services})
}

// @Summary		Update a branch service
// @Description	Updates the details of a specific service within a branch.
// @Tags			Branch Services
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			id			path		string						true	"Branch UUID"
// @Param			service_id	path		string						true	"Service UUID"
// @Param			service		body		UpdateBranchServicePayload	true	"Payload with fields to update"
// @Success		204			{object}	nil							"Service updated successfully"
// @Failure		400			{object}	map[string]interface{}		"Bad Request: Invalid UUID or payload"
// @Failure		500			{object}	map[string]interface{}		"Internal Server Error"
// @Router			/branches/{id}/services/{service_id} [put]
func (h *ProductsHandler) UpdateBranchServiceHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var payload UpdateBranchServicePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	service, err := payload.ToServiceBranchDetail()
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	service.BranchID, err = uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return

	}
	service.ServiceID, err = uuid.Parse(c.Param("service_id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return

	}
	err = h.services.ProductsServices.UpdateBranchService(ctx, service)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// @Summary		Remove a service from a branch
// @Description	Removes a specific service from a specific branch.
// @Tags			Branch Services
// @Security		ApiKeyAuth
// @Produce		json
// @Param			id			path		string					true	"Branch UUID"
// @Param			service_id	path		string					true	"Service UUID"
// @Success		204			{object}	nil						"Service removed successfully"
// @Failure		400			{object}	map[string]interface{}	"Bad Request: Invalid UUID format"
// @Failure		404			{object}	map[string]interface{}	"Not Found: Branch, service, or assignment not found"
// @Failure		500			{object}	map[string]interface{}	"Internal Server Error"
// @Router			/branches/{id}/services/{service_id} [delete]
func (h *ProductsHandler) DeleteBranchServiceHandler(c *gin.Context) {
	ctx := c.Request.Context()
	branchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	serviceID, err := uuid.Parse(c.Param("service_id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	err = h.services.ProductsServices.DeleteBranchService(ctx, branchID, serviceID)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// @Summary		Get branch service by ID
// @Description	Retrieves detailed information about a specific service for a specific branch.
// @Tags			Branch Services
// @Security		ApiKeyAuth
// @Produce		json
// @Param			id			path		string											true	"Branch UUID"
// @Param			service_id	path		string											true	"Service UUID"
// @Success		200			{object}	object{data=productsdomain.ServiceBranchDetail}	"Branch service details"
// @Failure		400			{object}	map[string]interface{}							"Bad Request: Invalid UUID format"
// @Failure		404			{object}	map[string]interface{}							"Not Found: Branch or service not found"
// @Failure		500			{object}	map[string]interface{}							"Internal Server Error"
// @Router			/branches/{id}/services/{service_id} [get]
func (h *ProductsHandler) GetBranchServiceByIDHandler(c *gin.Context) {
	ctx := c.Request.Context()
	branchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	serviceID, err := uuid.Parse(c.Param("service_id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	service, err := h.services.ProductsServices.GetBranchServiceByID(ctx, branchID, serviceID)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": service})
}
