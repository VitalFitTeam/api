package billingdomain

import (
	"context"

	"github.com/google/uuid"
)

type BillingRepository interface {
}

type PaymentMethodsRepository interface {
	GetPaymentMethods(ctx context.Context) ([]*PaymentMethods, error)
	GetPaymentMethodByID(ctx context.Context, methodID uuid.UUID) (*PaymentMethods, error)

	CreatePaymentMethod(ctx context.Context, paymentMethod *PaymentMethods) error
	UpdatePaymentMethod(ctx context.Context, paymentMethod *PaymentMethods) error
	DeletePaymentMethod(ctx context.Context, methodID uuid.UUID) error
}
