package store

import (
	authdomain "github.com/vitalfit/api/internal/auth/domain"
	authrepository "github.com/vitalfit/api/internal/auth/repository"
	"gorm.io/gorm"
)

type Storage struct {
	Users authdomain.UserRepository
	Roles authdomain.RolesRepository
}

func NewStorage(db *gorm.DB) Storage {
	return Storage{
		Users: authrepository.NewUserStore(db),
		Roles: authrepository.NewRoleStore(db),
	}
}
