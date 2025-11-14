package combosrepository

import "gorm.io/gorm"

type CombosStore struct {
	db *gorm.DB
}

func NewCombosStore(db *gorm.DB) *CombosStore {
	return &CombosStore{db: db}
}
