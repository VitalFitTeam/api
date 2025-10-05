package authrepository

import (
	"context"

	authdomain "github.com/vitalfit/api/internal/auth/domain"
	"gorm.io/gorm"
)

type UserRepositoryDAO struct {
	db *gorm.DB
}

func NewUserRepositoryDAO(db *gorm.DB) *UserRepositoryDAO {
	return &UserRepositoryDAO{
		db: db,
	}
}

func (s *UserRepositoryDAO) GetUser() error {

	return nil
}

func (s *UserRepositoryDAO) Create(ctx context.Context, user authdomain.Users) error {
	err := s.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}
