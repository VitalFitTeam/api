package appservices

import (
	authdomain "github.com/vitalfit/api/internal/auth/domain"
	authservices "github.com/vitalfit/api/internal/auth/services"
	logs "github.com/vitalfit/api/internal/shared/errors"
	"github.com/vitalfit/api/internal/store"
	"go.uber.org/zap"
)

type Services struct {
	AuthServices authdomain.AuthServicesInterface
	logs.LogErrors
}

func NewServices(store store.Storage, logger *zap.SugaredLogger) Services {
	return Services{
		AuthServices: authservices.NewAuthServices(store),
		LogErrors:    logs.NewLogErrors(logger),
	}
}
