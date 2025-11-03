package membershipshandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	shared_errors "github.com/vitalfit/api/internal/shared/errors"
)

// @Summary		Create a new membership type
// @Description	Adds a new membership type (plan) to the system.
// @Tags			Memberships
// @Accept			json
// @Produce		json
// @Security		ApiKeyAuth
// @Param			membership	body		CreateMembershipPayload				true	"Membership creation payload"
// @Success		201			{object}	membershipsdomain.MembershipType	"Membership type created successfully"
// @Failure		400			{object}	object{error=string}				"error: Bad Request"
// @Failure		409			{object}	object{error=string}				"error: Conflict"
// @Failure		500			{object}	object{error=string}				"error: Internal Server Error"
// @Router			/memberships/types [post]
func (h *MembershipHandler) CreateMembershipHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var payload CreateMembershipPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	membership, err := payload.toMembership()
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	err = h.services.MembershipServices.CreateMembershipType(ctx, membership)
	if err != nil {
		switch err {
		case shared_errors.ErrConflict:
			h.services.LogErrors.ConflictResponse(c, err)
		default:
			h.services.LogErrors.InternalServerError(c, err)
		}
		return
	}

	c.JSON(http.StatusCreated, membership)
}

// @Summary		Update a membership type
// @Description	Updates an existing membership type’s details.
// @Tags			Memberships
// @Accept			json
// @Produce		json
// @Security		ApiKeyAuth
// @Param			id			path		string								true	"Membership Type UUID"
// @Param			membership	body		UpdateMembershipPayload				true	"Membership update payload"
// @Success		200			{object}	membershipsdomain.MembershipType	"Membership type updated successfully"
// @Failure		400			{object}	object{error=string}				"error: Bad Request - Invalid ID or payload"
// @Failure		404			{object}	object{error=string}				"error: Not Found - Membership type not found"
// @Failure		500			{object}	object{error=string}				"error: Internal Server Error"
// @Router			/memberships/types/{id} [put]
func (h *MembershipHandler) UpdateMembershipHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var payload UpdateMembershipPayload

	membershipID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	membership, err := payload.toMembership()
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	membership.MembershipTypeID = membershipID

	err = h.services.MembershipServices.UpdateMembershipType(ctx, membership)
	if err != nil {
		switch err {
		case shared_errors.ErrNotFound:
			h.services.LogErrors.NotFoundResponse(c)
		default:
			h.services.LogErrors.InternalServerError(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, membership)
}

// @Summary		Delete a membership type
// @Description	Soft deletes a membership type from the system (sets inactive).
// @Tags			Memberships
// @Produce		json
// @Security		ApiKeyAuth
// @Param			id	path		string					true	"Membership Type UUID"
// @Success		204	{object}	nil						"No Content"
// @Failure		400	{object}	object{error=string}	"error: Bad Request - Invalid ID"
// @Failure		404	{object}	object{error=string}	"error: Not Found - Membership type not found"
// @Failure		500	{object}	object{error=string}	"error: Internal Server Error"
// @Router			/memberships/types/{id} [delete]
func (h *MembershipHandler) DeleteMembershipHandler(c *gin.Context) {
	ctx := c.Request.Context()
	membershipID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	err = h.services.MembershipServices.DeleteMembershipType(ctx, membershipID)
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

// @Summary		Get membership type by ID
// @Description	Retrieves a single membership type by its UUID.
// @Tags			Memberships
// @Produce		json
// @Security		ApiKeyAuth
// @Param			id	path		string	true	"Membership Type UUID"
// @Success		200	{object}	object{data=MembershipResponse}
// @Failure		400	{object}	object{error=string}	"error: Bad Request - Invalid ID"
// @Failure		404	{object}	object{error=string}	"error: Not Found - Membership type not found"
// @Failure		500	{object}	object{error=string}	"error: Internal Server Error"
// @Router			/memberships/types/{id} [get]
func (h *MembershipHandler) GetMembershipByIDHandler(c *gin.Context) {
	ctx := c.Request.Context()
	membershipID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	membership, err := h.services.MembershipServices.GetMembershipTypeByID(ctx, membershipID)
	if err != nil {
		switch err {
		case shared_errors.ErrNotFound:
			h.services.LogErrors.NotFoundResponse(c)
		default:
			h.services.LogErrors.InternalServerError(c, err)
		}
		return
	}

	resp := &MembershipResponse{
		MembershipTypeID: membership.MembershipTypeID,
		Name:             membership.Name,
		Description:      membership.Description,
		DurationDays:     membership.DurationDays,
		Price:            membership.Price,
		IsActive:         membership.IsActive,
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// @Summary		List all membership types
// @Description	Retrieves a list of all membership types.
// @Tags			Memberships
// @Produce		json
// @Security		ApiKeyAuth
// @Success		200	{object}	object{data=[]MembershipResponse}
// @Failure		500	{object}	object{error=string}	"error: Internal Server Error"
// @Router			/memberships/types [get]
func (h *MembershipHandler) GetMembershipsHandler(c *gin.Context) {
	ctx := c.Request.Context()
	memberships, err := h.services.MembershipServices.GetMembershipTypes(ctx)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	resp := make([]*MembershipResponse, 0, len(memberships))
	for _, m := range memberships {
		resp = append(resp, &MembershipResponse{
			MembershipTypeID: m.MembershipTypeID,
			Name:             m.Name,
			Description:      m.Description,
			DurationDays:     m.DurationDays,
			Price:            m.Price,
			IsActive:         m.IsActive,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}
