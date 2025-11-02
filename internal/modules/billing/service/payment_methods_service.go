package billingservice

import (
	"context"

	"github.com/google/uuid"
	billingdomain "github.com/vitalfit/api/internal/modules/billing/domain"
)

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

func (s *BillingService) CreatePaymentMethod(ctx context.Context, paymentMethod *billingdomain.PaymentMethods) error {
	err := s.store.PaymentMethods.CreatePaymentMethod(ctx, paymentMethod)
	if err != nil {
		return err
	}
	return nil

}

func (s *BillingService) UpdatePaymentMethod(ctx context.Context, paymentMethod *billingdomain.PaymentMethods) error {
	err := s.store.PaymentMethods.UpdatePaymentMethod(ctx, paymentMethod)
	if err != nil {
		return err
	}
	return nil
}

func (s *BillingService) DeletePaymentMethod(ctx context.Context, methodID uuid.UUID) error {
	err := s.store.PaymentMethods.DeletePaymentMethod(ctx, methodID)
	if err != nil {
		return err
	}
	return nil
}
