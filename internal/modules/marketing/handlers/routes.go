package marketinghandlers

import appservices "github.com/vitalfit/api/internal/app/services"

type MarketingHandlerInterface interface {
}

type MarketingHandler struct {
	services appservices.Services
}

func NewMarketingHandler(services appservices.Services) *MarketingHandler {
	return &MarketingHandler{services: services}
}
