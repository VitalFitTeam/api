package accessrepository

import "gorm.io/gorm"

type AccessStore struct {
	db *gorm.DB
}

func NewAccessStore(db *gorm.DB) *AccessStore {
	return &AccessStore{
		db: db,
	}
}
