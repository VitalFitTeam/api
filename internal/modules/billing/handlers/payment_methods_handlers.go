package billinghandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
