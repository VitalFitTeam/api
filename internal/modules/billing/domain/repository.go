package billingdomain

import (
	"context"

	"github.com/google/uuid"
	"github.com/vitalfit/api/pkg/pagination"
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
	GetBranchPaymentMethodByID(ctx context.Context, branchID uuid.UUID, methodID uuid.UUID) (*PaymentMethodsBranch, error)
}

type FiscalDocumentRepository interface {
	CreateFiscalDocumentType(ctx context.Context, docType *FiscalDocumentType) error
	GetFiscalDocumentTypes(ctx context.Context, fq pagination.PaginatedFeedQuery) ([]*FiscalDocumentType, error)
	GetFiscalDocumentTypesTotal(ctx context.Context, fq pagination.PaginatedFeedQuery) (int64, error)
	GetFiscalDocumentTypeByID(ctx context.Context, docTypeID uuid.UUID) (*FiscalDocumentType, error)
	UpdateFiscalDocumentType(ctx context.Context, docType *FiscalDocumentType) error
	DeleteFiscalDocumentType(ctx context.Context, docTypeID uuid.UUID) error
}
