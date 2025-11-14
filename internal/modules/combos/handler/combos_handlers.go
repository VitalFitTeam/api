package comboshandler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	shared_errors "github.com/vitalfit/api/internal/shared/errors"
	"github.com/vitalfit/api/pkg/pagination"
)

// @Summary		Create a new package
// @Description	Adds a new package with its items to the system.
// @Tags			Packages
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			package	body		CreatePackageRequest	true	"Package creation payload"
// @Success		201		{object}	map[string]interface{}	"Package created successfully"
// @Failure		400		{object}	map[string]interface{}	"Bad Request: Invalid payload"
// @Failure		500		{object}	map[string]interface{}	"Internal Server Error"
// @Router			/packages [post]
func (h *CombosHandler) CreatePackageHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var payload CreatePackageRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	pkg, err := payload.ToPackage()
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	err = h.services.CombosServices.CreatePackage(ctx, pkg)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(201, gin.H{"message": "Package created successfully"})

}

// @Summary		List all packages
// @Description	Retrieves a paginated list of all available packages, with optional searching by name.
// @Tags			Packages
// @Security		ApiKeyAuth
// @Produce		json
// @Param			limit	query		int										false	"Number of results per page"	default(10)
// @Param			page	query		int										false	"Page number for pagination"	default(1)
// @Param			sort	query		string									false	"Sort order (asc/desc)"			enums(asc, desc)	default(desc)
// @Param			search	query		string									false	"Search term for package name"
// @Success		200		{object}	pagination.PaginatedResponseTotal[PackageResponse]	"A paginated list of packages"
// @Failure		400		{object}	map[string]interface{}					"Bad Request: Invalid query parameters"
// @Failure		500		{object}	map[string]interface{}					"Internal server error"
// @Router			/packages [get]
func (h *CombosHandler) GetPackageHandler(c *gin.Context) {
	ctx := c.Request.Context()
	fq := pagination.PaginatedFeedQuery{
		Page:   1,
		Limit:  10,
		Search: "",
	}
	fq, err := fq.Parse(c.Request)
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	nextURL := fmt.Sprintf("/packages?limit=%d&page=%d&sort=%s", fq.Limit, fq.Page+1, fq.Sort)
	previousPage := fq.Page - 1
	if previousPage <= 0 {
		previousPage = 1
	}
	previousURL := fmt.Sprintf("/packages?limit=%d&page=%d&sort=%s", fq.Limit, previousPage, fq.Sort)

	packages, err := h.services.CombosServices.GetPackages(ctx, fq)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	total, err := h.services.CombosServices.GetPackagesTotal(ctx, fq)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	data := make([]PackageResponse, len(packages))
	for i, pkg := range packages {
		data[i] = PackageResponse{
			PackageID:   pkg.PackageID,
			Name:        pkg.Name,
			Description: pkg.Description,
			Price:       pkg.Price,
			IsActive:    pkg.IsActive,
			StartAt:     pkg.StartAt,
			EndAt:       pkg.EndAt,
			CreatedAt:   pkg.CreatedAt,
			UpdatedAt:   pkg.UpdatedAt,
		}
	}

	resp := pagination.PaginatedResponseTotal[PackageResponse]{
		Data:     data,
		Next:     nextURL,
		Previous: previousURL,
		Count:    int64(len(packages)),
		Total:    total,
	}
	c.JSON(200, gin.H{"data": resp})
}

// @Summary		Get package by ID
// @Description	Retrieves detailed information about a specific package by its UUID, including its items.
// @Tags			Packages
// @Security		ApiKeyAuth
// @Produce		json
// @Param			id	path		string					true	"Package UUID"
// @Success		200	{object}	object{data=PackageResponseByID}	"Package details"
// @Failure		400	{object}	map[string]interface{}	"Bad Request: Invalid UUID format"
// @Failure		404	{object}	map[string]interface{}	"Not Found: Package not found"
// @Failure		500	{object}	map[string]interface{}	"Internal Server Error"
// @Router			/packages/{id} [get]
func (h *CombosHandler) GetPackageByIDHandler(c *gin.Context) {
	ctx := c.Request.Context()
	packageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	pkg, err := h.services.CombosServices.GetPackageByID(ctx, packageID)
	if err != nil {
		switch err {
		case shared_errors.ErrNotFound:
			h.services.LogErrors.NotFoundResponse(c)
		default:
			h.services.LogErrors.InternalServerError(c, err)
		}
		return
	}

	packageItems := make([]PackageItemResponse, len(pkg.PackageItems))
	for i, item := range pkg.PackageItems {
		packageItems[i] = PackageItemResponse{
			ServiceID:        item.ServiceID,
			SessionsIncluded: item.SessionsIncluded,
		}
	}

	resp := PackageResponseByID{
		PackageResponse: PackageResponse{
			PackageID:   pkg.PackageID,
			Name:        pkg.Name,
			Description: pkg.Description,
			Price:       pkg.Price,
			IsActive:    pkg.IsActive,
			StartAt:     pkg.StartAt,
			EndAt:       pkg.EndAt,
			CreatedAt:   pkg.CreatedAt,
			UpdatedAt:   pkg.UpdatedAt,
		},
		PackageItems: packageItems,
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// @Summary		Update a package
// @Description	Updates an existing package's information, including its items.
// @Tags			Packages
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			id		path		string					true	"Package UUID"
// @Param			package	body		CreatePackageRequest	true	"Payload with fields to update"
// @Success		204		{object}	nil						"Package updated successfully"
// @Failure		400		{object}	map[string]interface{}	"Bad Request: Invalid UUID or payload"
// @Failure		404		{object}	map[string]interface{}	"Not Found: Package not found"
// @Failure		500		{object}	map[string]interface{}	"Internal Server Error"
// @Router			/packages/{id} [put]
func (h *CombosHandler) UpdatePackageHandler(c *gin.Context) {
	ctx := c.Request.Context()
	packageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	var payload CreatePackageRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	pkg, err := payload.ToPackage()
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	pkg.PackageID = packageID

	err = h.services.CombosServices.UpdatePackage(ctx, pkg)
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

// @Summary		Delete a package
// @Description	Deletes a specific package by its UUID.
// @Tags			Packages
// @Security		ApiKeyAuth
// @Produce		json
// @Param			id	path		string					true	"Package UUID"
// @Success		204	{object}	nil						"Package deleted successfully"
// @Failure		400	{object}	map[string]interface{}	"Bad Request: Invalid UUID format"
// @Failure		404	{object}	map[string]interface{}	"Not Found: Package not found"
// @Failure		500	{object}	map[string]interface{}	"Internal Server Error"
// @Router			/packages/{id} [delete]
func (h *CombosHandler) DeletePackageHandler(c *gin.Context) {
	ctx := c.Request.Context()
	packageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	err = h.services.CombosServices.DeletePackage(ctx, packageID)
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
