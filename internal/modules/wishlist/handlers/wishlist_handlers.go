package wishlisthandlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	shared_errors "github.com/vitalfit/api/internal/shared/errors"
)

// @Summary		Add service to wishlist
// @Description	Adds a service to the authenticated user's wishlist
// @Tags			Wishlist
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			payload	body		AddToWishlistRequest	true	"Service ID"
// @Success		201		{object}	AddToWishlistResponse	"Service added to wishlist"
// @Failure		400		{object}	map[string]interface{}	"Bad Request"
// @Failure		409		{object}	map[string]interface{}	"Service already in wishlist"
// @Failure		500		{object}	map[string]interface{}	"Internal Server Error"
// @Router			/wishlist [post]
func (h *WishlistHandlers) AddToWishlistHandler(c *gin.Context) {
	ctx := c.Request.Context()

	var payload AddToWishlistRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	user := h.services.UserServices.GetUserFromContext(c)

	wishlistID, err := h.services.WishlistServices.AddToWishlist(ctx, user.UserID, payload.ServiceID)
	if err != nil {
		switch {
		case errors.Is(err, shared_errors.ErrConflict):
			h.services.LogErrors.ConflictResponse(c, errors.New("service already in wishlist"))
		case errors.Is(err, shared_errors.ErrNotFound):
			h.services.LogErrors.NotFoundResponse(c)
		default:
			h.services.LogErrors.InternalServerError(c, err)
		}
		return
	}

	c.JSON(http.StatusCreated, AddToWishlistResponse{
		WishlistID: wishlistID,
		Message:    "Service added to wishlist successfully",
	})
}

// @Summary		Remove service from wishlist
// @Description	Removes a service from the authenticated user's wishlist
// @Tags			Wishlist
// @Security		ApiKeyAuth
// @Produce		json
// @Param			id	path		string						true	"Wishlist Item UUID"
// @Success		200	{object}	RemoveFromWishlistResponse	"Service removed from wishlist"
// @Failure		400	{object}	map[string]interface{}		"Bad Request"
// @Failure		404	{object}	map[string]interface{}		"Not Found"
// @Failure		500	{object}	map[string]interface{}		"Internal Server Error"
// @Router			/wishlist/{id} [delete]
func (h *WishlistHandlers) RemoveFromWishlistHandler(c *gin.Context) {
	ctx := c.Request.Context()

	wishlistID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, errors.New("invalid wishlist id format"))
		return
	}

	if err := h.services.WishlistServices.RemoveFromWishlist(ctx, wishlistID); err != nil {
		if errors.Is(err, shared_errors.ErrNotFound) {
			h.services.LogErrors.NotFoundResponse(c)
			return
		}
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, RemoveFromWishlistResponse{
		Message: "Service removed from wishlist successfully",
	})
}

// @Summary		Get user wishlist
// @Description	Returns all services in the authenticated user's wishlist
// @Tags			Wishlist
// @Security		ApiKeyAuth
// @Produce		json
// @Success		200	{object}	object{data=[]WishlistItemResponse}	"Wishlist items"
// @Failure		500	{object}	map[string]interface{}				"Internal Server Error"
// @Router			/wishlist [get]
func (h *WishlistHandlers) GetUserWishlistHandler(c *gin.Context) {
	ctx := c.Request.Context()

	user := h.services.UserServices.GetUserFromContext(c)

	wishlist, err := h.services.WishlistServices.GetUserWishlist(ctx, user.UserID)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	response := make([]*WishlistItemResponse, 0, len(wishlist))
	for _, item := range wishlist {
		response = append(response, &WishlistItemResponse{
			WishlistID:  item.WishlistID,
			ServiceID:   item.ServiceID,
			ServiceName: item.ServiceName,
			Description: item.Description,
			CreatedAt:   item.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}
