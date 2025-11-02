package billingdomain

import (
	"context"

	"github.com/google/uuid"
)

type BillingServiceInterface interface {
	GetPaymentMethodByID(ctx context.Context, methodID uuid.UUID) (*PaymentMethods, error)
	GetPaymentMethods(ctx context.Context) ([]*PaymentMethods, error)
}
