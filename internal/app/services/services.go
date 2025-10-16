package appservices

import (
	"github.com/vitalfit/api/config"
	authdomain "github.com/vitalfit/api/internal/modules/auth/domain"
	authservices "github.com/vitalfit/api/internal/modules/auth/services"
	logs "github.com/vitalfit/api/internal/shared/errors"
	"github.com/vitalfit/api/internal/store"
	"github.com/vitalfit/api/pkg/mailer"
	"go.uber.org/zap"
)

type Services struct {
	AuthServices authdomain.AuthServicesInterface
	UserServices authdomain.UserServicesInterface
	logs.LogErrors
	Logger *zap.SugaredLogger
}

func NewServices(store store.Storage, logger *zap.SugaredLogger, cfg config.Config, auth authdomain.Authenticator, mailer mailer.Client) Services {
	return Services{
		AuthServices: authservices.NewAuthServices(store, cfg, auth, mailer),
		UserServices: authservices.NewUserService(store),
		LogErrors:    logs.NewLogErrors(logger),
		Logger:       logger,
	}
}
