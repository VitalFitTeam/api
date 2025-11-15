package billingrepository

import "gorm.io/gorm"

type FiscalDocumentStore struct {
	db *gorm.DB
}

func NewFiscalDocumentStore(db *gorm.DB) *FiscalDocumentStore {
	return &FiscalDocumentStore{db: db}

}
