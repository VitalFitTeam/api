package authrepository

import (
	"context"

	authdomain "github.com/vitalfit/api/internal/auth/domain"
	shared_errors "github.com/vitalfit/api/internal/shared/errors"
	"gorm.io/gorm"
)

type RoleStore struct {
	db *gorm.DB
}

func NewRoleStore(db *gorm.DB) *RoleStore {
	return &RoleStore{db: db}
}

func (s *RoleStore) GetByName(ctx context.Context, name string) (*authdomain.Roles, error) {
	var role authdomain.Roles
	err := s.db.WithContext(ctx).Where("name = ?", name).First(&role).Error
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, shared_errors.ErrNotFound
		case gorm.ErrDuplicatedKey:
			return nil, shared_errors.ErrConflict
		default:
			return nil, err
		}

	}
	return &role, nil
}
