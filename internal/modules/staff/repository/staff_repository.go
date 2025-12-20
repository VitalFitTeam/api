package staffrepository

import "gorm.io/gorm"

type StaffStore struct {
	db *gorm.DB
}

func NewStaffStore(db *gorm.DB) *StaffStore {
	return &StaffStore{db: db}
}
