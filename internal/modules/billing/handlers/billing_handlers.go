package billinghandlers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary		Create a new invoice
// @Description	Creates a new invoice for a specific user and branch, containing a list of items. If the user making the request is a client, the invoice is created for them. If the user is staff (with 'billing:create_invoice' permission), the 'user_id' in the payload is required to specify for whom the invoice is being created.
// @Tags			Billing
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			invoice	body		CreateInvoicePayload	true	"Invoice creation payload"
// @Success		201		{object}	object{message=string}	"Invoice created successfully"
// @Failure		400		{object}	object{error=string}	"Bad Request (e.g., invalid payload, missing user_id for staff)"
// @Failure		403		{object}	object{error=string}	"Forbidden (user does not have permission)"
// @Failure		500		{object}	object{error=string}	"Internal Server Error"
// @Router			/billing/invoices [post]
func (h *BillingHandlers) CreateInvoiceHandler(c *gin.Context) {
	var payload CreateInvoicePayload
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	user := h.services.UserServices.GetUserFromContext(c)
	var targetUserID uuid.UUID

	if user.Role.Name != "client" {
		permission := "billing:create_invoice"
		if user.Role.Name != "super_admin" {
			ok, err := h.services.UserServices.RoleHasPermission(ctx, user.RoleID, permission)
			if err != nil {
				h.services.LogErrors.InternalServerError(c, err)
				return
			}
			if !ok {
				h.services.LogErrors.ForbiddenResponse(c)
				return
			}
		}
		if payload.UserID == nil {
			h.services.LogErrors.BadRequestResponse(c, errors.New("user_id is required for non-client users"))
			return
		}
		targetUserID = *payload.UserID
	} else {
		targetUserID = user.UserID
	}

	invoice := payload.ToInvoice()
	invoice.UserID = targetUserID

	items := payload.ToInvoiceItems()

	err := h.services.BillingServices.CreateInvoice(ctx, invoice, items)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(201, gin.H{"message": "Invoice created successfully"})
}
