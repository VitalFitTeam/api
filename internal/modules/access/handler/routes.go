package accesshandler

import appservices "github.com/vitalfit/api/internal/app/services"

type AccessHandlerInterface interface {
}

type AccessHandler struct {
	services appservices.Services
}

func NewAccessHandler(services appservices.Services) *AccessHandler {
	return &AccessHandler{
		services: services,
	}
}
