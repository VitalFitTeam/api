package store

import (
	authdomain "github.com/vitalfit/api/internal/modules/auth/domain"
	authrepository "github.com/vitalfit/api/internal/modules/auth/repository"
	branchdomain "github.com/vitalfit/api/internal/modules/branches/domain"
	branchrepository "github.com/vitalfit/api/internal/modules/branches/repository"
	"gorm.io/gorm"
)

type Storage struct {
	Users    authdomain.UserRepository
	Roles    authdomain.RolesRepository
	Branches branchdomain.BranchesRepository
}

func NewStorage(db *gorm.DB) Storage {
	return Storage{
		Users:    authrepository.NewUserStore(db),
		Roles:    authrepository.NewRoleStore(db),
		Branches: branchrepository.NewBranchesStore(db),
	}
}
