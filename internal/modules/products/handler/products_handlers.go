package productshandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *ProductsHandler) ListServiceCategoriesHandler(c *gin.Context) {
	ctx := c.Request.Context()
	serviceCategories, err := h.services.ProductsServices.ListServiceCategories(ctx)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": serviceCategories})

}

func (h *ProductsHandler) CreateServiceHandler(c *gin.Context) {

}

func (h *ProductsHandler) GetServicesHandler(c *gin.Context) {

}

func (h *ProductsHandler) DeleteServiceHandler(c *gin.Context) {

}

func (h *ProductsHandler) GetServiceByIDHandler(c *gin.Context) {

}

func (h *ProductsHandler) UpdatServicetHandler(c *gin.Context) {

}
