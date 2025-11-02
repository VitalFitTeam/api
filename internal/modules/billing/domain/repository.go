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

	AddPaymentMethodsToBranch(ctx context.Context, branchMethod []*PaymentMethodsBranch) error
	DeletePaymentMethodsFromBranch(ctx context.Context, branchID, methodID uuid.UUID) error
	GetPaymentMethodsFromBranch(ctx context.Context, branchID uuid.UUID) ([]*PaymentMethodsBranch, error)
	UpsertBranchPaymentConfig(ctx context.Context, branchMethod *PaymentMethodsBranch) error
}
