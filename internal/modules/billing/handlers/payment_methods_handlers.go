package billinghandlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	billingdomain "github.com/vitalfit/api/internal/modules/billing/domain"
	shared_errors "github.com/vitalfit/api/internal/shared/errors"
)

// @Summary		Get all payment methods
// @Description	Retrieves a list of all available payment methods in the system.
// @Tags			Billing
// @Security		ApiKeyAuth
// @Produce		json
// @Success		200	{object}	object{data=[]billingdomain.PaymentMethods}	"A list of available payment methods"
// @Failure		500	{object}	object{error=string}						"Error: internal server error"
// @Router			/billing/payment-methods [get]
func (h *BillingHandlers) GetPaymentMethodsHandler(c *gin.Context) {
	ctx := c.Request.Context()
	paymentMethods, err := h.services.BillingServices.GetPaymentMethods(ctx)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": paymentMethods,
	})

}

// @Summary		Create a new payment method
// @Description	Adds a new payment method to the system.
// @Tags			Billing
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			paymentMethod	body		CreatePaymentMethodPayload	true	"Payment method creation payload"
// @Success		201				{object}	object{message=string}		"payment method created"
// @Failure		400				{object}	object{error=string}		"Bad Request"
// @Failure		500				{object}	object{error=string}		"Internal Server Error"
// @Router			/billing/payment-methods [post]
func (h *BillingHandlers) CreatePaymentMethodHandler(c *gin.Context) {
	var payload CreatePaymentMethodPayload
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	paymentMethod := payload.toPaymentMethod()
	if err := h.services.BillingServices.CreatePaymentMethod(ctx, paymentMethod); err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "payment method created",
	})
}

// @Summary		Update a payment method
// @Description	Updates an existing payment method's details.
// @Tags			Billing
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			id				path		string						true	"Payment Method UUID"
// @Param			paymentMethod	body		UpdatePaymentMethodPayload	true	"Payment method update payload"
// @Success		204				{object}	nil							"No Content"
// @Failure		400				{object}	object{error=string}		"Bad Request"
// @Failure		404				{object}	object{error=string}		"Not Found"
// @Failure		500				{object}	object{error=string}		"Internal Server Error"
// @Router			/billing/payment-methods/{id} [put]
func (h *BillingHandlers) UpdatePaymentMethodHandler(c *gin.Context) {
	ctx := c.Request.Context()
	methodID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	var payload UpdatePaymentMethodPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	paymentmethod := payload.toPaymentMethod()
	paymentmethod.MethodID = methodID
	err = h.services.BillingServices.UpdatePaymentMethod(ctx, paymentmethod)
	if err != nil {
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

// @Summary		Delete a payment method
// @Description	Deletes a payment method from the system.
// @Tags			Billing
// @Security		ApiKeyAuth
// @Produce		json
// @Param			id	path		string					true	"Payment Method UUID"
// @Success		204	{object}	nil						"No Content"
// @Failure		400	{object}	object{error=string}	"Bad Request"
// @Failure		404	{object}	object{error=string}	"Not Found"
// @Failure		500	{object}	object{error=string}	"Internal Server Error"
// @Router			/billing/payment-methods/{id} [delete]
func (h *BillingHandlers) DeletePaymentMethodHandler(c *gin.Context) {
	ctx := c.Request.Context()
	methodID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	err = h.services.BillingServices.DeletePaymentMethod(ctx, methodID)
	if err != nil {
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

// @Summary		Get payment method by ID
// @Description	Retrieves a single payment method by its UUID.
// @Tags			Billing
// @Security		ApiKeyAuth
// @Produce		json
// @Param			id	path		string										true	"Payment Method UUID"
// @Success		200	{object}	object{data=billingdomain.PaymentMethods}	"Payment method details"
// @Failure		400	{object}	object{error=string}						"Bad Request"
// @Failure		404	{object}	object{error=string}						"Not Found"
// @Failure		500	{object}	object{error=string}						"Internal Server Error"
// @Router			/billing/payment-methods/{id} [get]
func (h *BillingHandlers) GetPaymentMethodByIDHandler(c *gin.Context) {
	ctx := c.Request.Context()
	methodID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	paymentMethod, err := h.services.BillingServices.GetPaymentMethodByID(ctx, methodID)
	if err != nil {
		switch err {
		case shared_errors.ErrNotFound:
			h.services.LogErrors.NotFoundResponse(c)
		default:
			h.services.LogErrors.InternalServerError(c, err)
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": paymentMethod})
}

// @Summary		Add payment methods to a branch
// @Description	Assigns one or more payment methods to a specific branch with custom configurations.
// @Tags			Branch Payment Methods
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			id		path		string						true	"Branch UUID"
// @Param			payload	body		[]MethodBranchConfigPayload	true	"An array of payment method configurations to add to the branch"
// @Success		201		{object}	object{message=string}		"Payment methods added to branch"
// @Failure		400		{object}	object{error=string}		"Bad Request (e.g., invalid UUID, invalid payload)"
// @Failure		404		{object}	object{error=string}		"Not Found (e.g., branch or payment method not found)"
// @Failure		500		{object}	object{error=string}		"Internal Server Error"
// @Router			/branches/{id}/payment-methods [post]
func (h *BillingHandlers) AddPaymentMethodsToBranchHandler(c *gin.Context) {
	ctx := c.Request.Context()
	branchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, errors.New("invalid branch id"))
		return
	}

	var payloads []MethodBranchConfigPayload
	if err := c.ShouldBindJSON(&payloads); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	if len(payloads) == 0 {
		h.services.LogErrors.BadRequestResponse(c, errors.New("payload array cannot be empty"))
		return
	}

	branchMethods := make([]*billingdomain.PaymentMethodsBranch, 0, len(payloads))

	for _, payload := range payloads {
		methodID, err := uuid.Parse(payload.MethodID)
		if err != nil {
			h.services.LogErrors.BadRequestResponse(c, err)
			return
		}

		paymentMethod, err := h.services.BillingServices.GetPaymentMethodByID(ctx, methodID)
		if err != nil {
			h.services.LogErrors.NotFoundResponse(c)
			return
		}

		if err := h.validateBranchConfig(payload.Configuration, paymentMethod); err != nil {
			h.services.LogErrors.BadRequestResponse(c, err)
			return
		}

		branchMethod := &billingdomain.PaymentMethodsBranch{
			BranchID:            branchID,
			MethodID:            paymentMethod.MethodID,
			DisplayName:         payload.DisplayName,
			Configuration:       payload.Configuration,
			Visibility:          billingdomain.BranchPaymentVisibilityEnum(payload.Visibility),
			SurchargeFixed:      payload.SurchargeFixed,
			SurchargePercentage: payload.SurchargePercentage,
		}

		branchMethods = append(branchMethods, branchMethod)
	}

	err = h.services.BillingServices.AddPaymentMethodsToBranch(ctx, branchMethods)
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Payment methods added to branch"})
}

// @Summary		Remove a payment method from a branch
// @Description	Removes a specific payment method configuration from a specific branch.
// @Tags			Branch Payment Methods
// @Security		ApiKeyAuth
// @Produce		json
// @Param			id			path	string	true	"Branch UUID"
// @Param			method_id	path	string	true	"Payment Method UUID"
// @Success		204			"No Content"
// @Failure		400			{object}	object{error=string}	"Bad Request (e.g., invalid UUID)"
// @Failure		500			{object}	object{error=string}	"Internal Server Error"
// @Router			/branches/{id}/payment-methods/{method_id} [delete]
func (h *BillingHandlers) DeletePaymentMethodsFromBranchHandler(c *gin.Context) {
	ctx := c.Request.Context()
	branchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	methodID, err := uuid.Parse(c.Param("method_id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	err = h.services.BillingServices.DeletePaymentMethodsFromBranch(ctx, branchID, methodID)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)

}

// @Summary		List payment methods for a branch
// @Description	Retrieves all payment methods configured for a specific branch.
// @Tags			Branch Payment Methods
// @Security		ApiKeyAuth
// @Produce		json
// @Param			id	path		string	true	"Branch UUID"
// @Success		200	{object}	object{data=[]billingdomain.PaymentMethodsBranch}
// @Failure		400	{object}	object{error=string}	"Bad Request (e.g., invalid UUID)"
// @Failure		500	{object}	object{error=string}	"Internal Server Error"
// @Router			/branches/{id}/payment-methods [get]
func (h *BillingHandlers) GetPaymentMethodsFromBranchHandler(c *gin.Context) {
	ctx := c.Request.Context()
	branchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	Branchmethods, err := h.services.BillingServices.GetPaymentMethodsFromBranch(ctx, branchID)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return

	}
	c.JSON(http.StatusOK, gin.H{"data": Branchmethods})

}

// @Summary		Update a branch payment method configuration
// @Description	Updates the configuration of a specific payment method for a specific branch.
// @Tags			Branch Payment Methods
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			id			path		string						true	"Branch UUID"
// @Param			method_id	path		string						true	"Payment Method UUID"
// @Param			payload		body		UpdateBranchConfigPayload	true	"Configuration update payload"
// @Success		200			{object}	object{message=string}		"Payment method configuration updated"
// @Failure		400			{object}	object{error=string}		"Bad Request (e.g., invalid UUID, invalid payload)"
// @Failure		404			{object}	object{error=string}		"Not Found (e.g., payment method not found)"
// @Failure		500			{object}	object{error=string}		"Internal Server Error"
// @Router			/branches/{id}/payment-methods/{method_id} [put]
func (h *BillingHandlers) UpdatePaymentMethodFromBranchHandler(c *gin.Context) {
	ctx := c.Request.Context()
	branchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, errors.New("invalid branch id"))
		return
	}
	methodID, err := uuid.Parse(c.Param("method_id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, errors.New("invalid method id"))
		return
	}

	var payload UpdateBranchConfigPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	paymentMethod, err := h.services.BillingServices.GetPaymentMethodByID(ctx, methodID)
	if err != nil {
		h.services.LogErrors.NotFoundResponse(c)
		return
	}

	if err := h.validateBranchConfig(payload.Configuration, paymentMethod); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	if payload.IsActive != "" {
		isActive, err := strconv.ParseBool(payload.IsActive)
		if err != nil {
			h.services.LogErrors.BadRequestResponse(c, err)
			return
		}
		payload.IsActive = strconv.FormatBool(isActive)
	}
	methodStatus, err := strconv.ParseBool(payload.IsActive)
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	branchMethod := &billingdomain.PaymentMethodsBranch{
		BranchID:            branchID,
		MethodID:            methodID,
		IsActive:            methodStatus,
		DisplayName:         payload.DisplayName,
		Configuration:       payload.Configuration,
		Visibility:          billingdomain.BranchPaymentVisibilityEnum(payload.Visibility),
		SurchargeFixed:      payload.SurchargeFixed,
		SurchargePercentage: payload.SurchargePercentage,
	}

	err = h.services.BillingServices.UpsertBranchPaymentConfig(ctx, branchMethod)
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment method configuration updated"})
}

// @Summary		Get branch payment method by ID
// @Description	Retrieves a single payment method configuration for a specific branch by its UUID and the method's UUID.
// @Tags			Branch Payment Methods
// @Security		ApiKeyAuth
// @Produce		json
// @Param			id			path		string	true	"Branch UUID"
// @Param			method_id	path		string	true	"Payment Method UUID"
// @Success		200			{object}	object{data=BranchPaymentMethodResponse}
// @Failure		400			{object}	object{error=string}	"Error: Bad Request (e.g., invalid UUID)"
// @Failure		404			{object}	object{error=string}	"Error: Not Found"
// @Failure		500			{object}	object{error=string}	"Error: Internal Server Error"
// @Router			/branches/{id}/payment-methods/{method_id} [get]
func (h *BillingHandlers) GetBranchPaymentMethodByIDHandler(c *gin.Context) {
	ctx := c.Request.Context()
	branchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	methodID, err := uuid.Parse(c.Param("method_id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	branchMethod, err := h.services.BillingServices.GetBranchPaymentMethodByID(ctx, branchID, methodID)
	if err != nil {
		switch err {
		case shared_errors.ErrNotFound:
			h.services.LogErrors.NotFoundResponse(c)
		default:
			h.services.LogErrors.InternalServerError(c, err)
		}
		return
	}

	resp := BranchPaymentMethodResponse{
		BranchID:            branchMethod.BranchID,
		MethodID:            branchMethod.MethodID,
		DisplayName:         branchMethod.DisplayName,
		Configuration:       branchMethod.Configuration,
		Visibility:          string(branchMethod.Visibility),
		SurchargeFixed:      branchMethod.SurchargeFixed,
		SurchargePercentage: branchMethod.SurchargePercentage,
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}
