package authrepository

import (
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

func (s *UserRepositoryDAO) Create() {

}
