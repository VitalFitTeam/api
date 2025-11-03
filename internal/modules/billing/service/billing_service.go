package billingservice

import (
	"github.com/vitalfit/api/internal/store"
)

type BillingService struct {
	store store.Storage
}

func NewBillingService(store store.Storage) *BillingService {
	return &BillingService{
		store: store,
	}
}
