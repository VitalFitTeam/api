package marketinghandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	shared_errors "github.com/vitalfit/api/internal/shared/errors"
)

// @Summary		Create a new banner
// @Description	Adds a new promotional banner to the system.
// @Tags			Marketing
// @Accept			json
// @Produce		json
// @Security		ApiKeyAuth
// @Param			banner	body		CreateBannerPayload		true	"Banner creation payload"
// @Success		201		{object}	marketingdomain.Banner	"Banner created successfully"
// @Failure		400		{object}	object{error=string}	"error: Bad Request"
// @Failure		409		{object}	object{error=string}	"error: Conflict"
// @Failure		500		{object}	object{error=string}	"error: Internal Server Error"
// @Router			/marketing/banners [post]
func (h *MarketingHandler) CreateBannerHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var payload CreateBannerPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	banner, err := payload.toBanner()
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	err = h.services.MarketingServices.CreateBanner(ctx, banner)
	if err != nil {
		switch err {
		case shared_errors.ErrConflict:
			h.services.LogErrors.ConflictResponse(c, err)
		default:
			h.services.LogErrors.InternalServerError(c, err)
		}
		return
	}
	c.JSON(http.StatusCreated, banner)
}

// @Summary		Update a banner
// @Description	Updates an existing promotional banner's details.
// @Tags			Marketing
// @Accept			json
// @Produce		json
// @Security		ApiKeyAuth
// @Param			id		path		string					true	"Banner UUID"
// @Param			banner	body		UpdateBannerPayload		true	"Banner update payload"
// @Success		200		{object}	marketingdomain.Banner	"Banner updated successfully"
// @Failure		400		{object}	object{error=string}	"error: Bad Request - Invalid ID or payload"
// @Failure		404		{object}	object{error=string}	"error: Not Found - Banner not found"
// @Failure		500		{object}	object{error=string}	"error: Internal Server Error"
// @Router			/marketing/banners/{id} [put]
func (h *MarketingHandler) UpdateBannerHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var payload UpdateBannerPayload
	bannerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	banner, err := payload.toBanner()
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	banner.BannerID = bannerID
	err = h.services.MarketingServices.UpdateBanner(ctx, banner)
	if err != nil {
		switch err {
		case shared_errors.ErrNotFound:
			h.services.LogErrors.NotFoundResponse(c)
		default:
			h.services.LogErrors.InternalServerError(c, err)
		}
		return
	}
	c.JSON(http.StatusOK, banner)
}

// @Summary		Delete a banner
// @Description	Deletes a promotional banner from the system.
// @Tags			Marketing
// @Produce		json
// @Security		ApiKeyAuth
// @Param			id	path		string					true	"Banner UUID"
// @Success		204	{object}	nil						"No Content"
// @Failure		400	{object}	object{error=string}	"error: Bad Request - Invalid ID"
// @Failure		404	{object}	object{error=string}	"error: Not Found - Banner not found"
// @Failure		500	{object}	object{error=string}	"error: Internal Server Error"
// @Router			/marketing/banners/{id} [delete]
func (h *MarketingHandler) DeleteBannerHandler(c *gin.Context) {
	ctx := c.Request.Context()
	bannerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	err = h.services.MarketingServices.DeleteBanner(ctx, bannerID)
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

// @Summary		Get banner by ID
// @Description	Retrieves a single promotional banner by its UUID.
// @Tags			Marketing
// @Produce		json
// @Security		ApiKeyAuth
// @Param			id	path		string	true	"Banner UUID"
// @Success		200	{object}	object{data=BannerResponse}
// @Failure		400	{object}	object{error=string}	"error: Bad Request - Invalid ID"
// @Failure		404	{object}	object{error=string}	"error: Not Found - Banner not found"
// @Failure		500	{object}	object{error=string}	"error: Internal Server Error"
// @Router			/marketing/banners/{id} [get]
func (h *MarketingHandler) GetBannerByIDHandler(c *gin.Context) {
	ctx := c.Request.Context()
	bannerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	banner, err := h.services.MarketingServices.GetBannerByID(ctx, bannerID)
	if err != nil {
		switch err {
		case shared_errors.ErrNotFound:
			h.services.LogErrors.NotFoundResponse(c)
		default:
			h.services.LogErrors.InternalServerError(c, err)
		}
		return
	}
	resp := &BannerResponse{
		BannerID: banner.BannerID,
		Name:     banner.Name,
		ImageURL: banner.ImageURL,
		LinkURL:  banner.LinkURL,
		IsActive: banner.IsActive,
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// @Summary		List all banners
// @Description	Retrieves a list of all promotional banners.
// @Tags			Marketing
// @Produce		json
// @Security		ApiKeyAuth
// @Success		200	{object}	object{data=[]BannerResponse}
// @Failure		500	{object}	object{error=string}	"error: Internal Server Error"
// @Router			/marketing/banners [get]
func (h *MarketingHandler) GetBannersHandler(c *gin.Context) {
	ctx := c.Request.Context()
	banners, err := h.services.MarketingServices.GetBanners(ctx)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	resp := make([]*BannerResponse, 0, len(banners))
	for _, banner := range banners {
		resp = append(resp, &BannerResponse{

			BannerID: banner.BannerID,
			Name:     banner.Name,
			ImageURL: banner.ImageURL,
			LinkURL:  banner.LinkURL,
			IsActive: banner.IsActive,
		})

	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}
