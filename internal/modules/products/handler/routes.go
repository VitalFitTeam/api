package productshandler

import appservices "github.com/vitalfit/api/internal/app/services"

type ProductsHandlerInterface interface {
}

type ProductsHandler struct {
	services appservices.Services
}

func NewProductsHandler(services appservices.Services) *ProductsHandler {
	return &ProductsHandler{services: services}
}
