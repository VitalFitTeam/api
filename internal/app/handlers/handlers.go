package apphandlers

import (
	appservices "github.com/vitalfit/api/internal/app/services"
	authhandlers "github.com/vitalfit/api/internal/auth/handlers"
)

type Handlers struct {
	AuthHandlers authhandlers.AuthHandlersInterface
}

func NewAppHandlers(services appservices.Services) Handlers {
	return Handlers{
		AuthHandlers: authhandlers.NewAuthHandlers(services),
	}

}
