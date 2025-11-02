package billingservice

import (
	"context"

	"github.com/google/uuid"
	billingdomain "github.com/vitalfit/api/internal/modules/billing/domain"
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

func (s *BillingService) GetPaymentMethodByID(ctx context.Context, methodID uuid.UUID) (*billingdomain.PaymentMethods, error) {
	paymentMethod, err := s.store.PaymentMethods.GetPaymentMethodByID(ctx, methodID)
	if err != nil {
		return nil, err
	}
	return paymentMethod, nil
}

func (s *BillingService) GetPaymentMethods(ctx context.Context) ([]*billingdomain.PaymentMethods, error) {
	paymentMethods, err := s.store.PaymentMethods.GetPaymentMethods(ctx)
	if err != nil {
		return nil, err
	}
	return paymentMethods, nil
}
