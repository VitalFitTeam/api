package billingrepository

import "gorm.io/gorm"

type Billingstore struct {
	db *gorm.DB
}

func NewBillingStore(db *gorm.DB) *Billingstore {
	return &Billingstore{
		db: db,
	}
}
