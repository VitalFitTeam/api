package store

import (
	"github.com/vitalfit/api/config"
	authdomain "github.com/vitalfit/api/internal/auth/domain"
	authrepository "github.com/vitalfit/api/internal/auth/repository"
	"github.com/vitalfit/api/pkg/mailer"
	"gorm.io/gorm"
)

type Storage struct {
	Users authdomain.UserRepository
	Roles authdomain.RolesRepository
	config.Config
	Mailer mailer.Client
	Auth   authdomain.Authenticator
}

func NewStorage(db *gorm.DB, cfg config.Config, mailer mailer.Client, Auth authdomain.Authenticator) Storage {
	return Storage{
		Users:  authrepository.NewUserRepositoryDAO(db),
		Roles:  authrepository.NewRoleStore(db),
		Config: cfg,
		Mailer: mailer,
		Auth:   Auth,
	}
}
