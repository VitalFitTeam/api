package billingservice

import "github.com/vitalfit/api/internal/store"

type fiscalDocumentService struct {
	store store.Storage
}

func NewFiscalDocumentService(store store.Storage) *fiscalDocumentService {
	return &fiscalDocumentService{store: store}
}
