package billinghandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
