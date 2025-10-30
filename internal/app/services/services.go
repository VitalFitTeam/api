package appservices

import (
	"github.com/vitalfit/api/config"
	authdomain "github.com/vitalfit/api/internal/modules/auth/domain"
	authservices "github.com/vitalfit/api/internal/modules/auth/services"
	branchdomain "github.com/vitalfit/api/internal/modules/branches/domain"
	branchservices "github.com/vitalfit/api/internal/modules/branches/services"
	instructordomain "github.com/vitalfit/api/internal/modules/instructor/domain"
	instructorservices "github.com/vitalfit/api/internal/modules/instructor/services"
	logs "github.com/vitalfit/api/internal/shared/errors"
	"github.com/vitalfit/api/internal/store"
	"github.com/vitalfit/api/pkg/mailer"
	"go.uber.org/zap"
)

type Services struct {
	AuthServices       authdomain.AuthServicesInterface
	UserServices       authdomain.UserServicesInterface
	BranchesServices   branchdomain.BranchesServicesInterface
	LocationsServices  branchdomain.LocationsServicesInterface
	InstructorServices instructordomain.InstructorServiceInterface
	logs.LogErrors
	Logger *zap.SugaredLogger
}

func NewServices(store store.Storage, logger *zap.SugaredLogger, cfg config.Config, auth authdomain.Authenticator, mailer mailer.Client) Services {
	return Services{
		AuthServices:       authservices.NewAuthServices(store, cfg, auth, mailer),
		UserServices:       authservices.NewUserService(store),
		BranchesServices:   branchservices.NewBranchServices(store, cfg),
		LocationsServices:  branchservices.NewLocationsServices(store),
		InstructorServices: instructorservices.NewInstructorServices(store, cfg),
		LogErrors:          logs.NewLogErrors(logger),
		Logger:             logger,
	}
}
