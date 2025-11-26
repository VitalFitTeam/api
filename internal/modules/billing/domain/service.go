package billingdomain

import (
	"context"

	"github.com/google/uuid"
	"github.com/vitalfit/api/pkg/pagination"
)

type BillingServiceInterface interface {
	//Payment-Methods
	GetPaymentMethodByID(ctx context.Context, methodID uuid.UUID) (*PaymentMethods, error)
	GetPaymentMethods(ctx context.Context) ([]*PaymentMethods, error)
	CreatePaymentMethod(ctx context.Context, paymentMethod *PaymentMethods) error
	UpdatePaymentMethod(ctx context.Context, paymentMethod *PaymentMethods) error
	DeletePaymentMethod(ctx context.Context, methodID uuid.UUID) error
	//Payment-Methods-Branch
	AddPaymentMethodsToBranch(ctx context.Context, branchMethod []*PaymentMethodsBranch) error
	DeletePaymentMethodsFromBranch(ctx context.Context, branchID, methodID uuid.UUID) error
	GetPaymentMethodsFromBranch(ctx context.Context, branchID uuid.UUID) ([]*PaymentMethodsBranch, error)
	UpsertBranchPaymentConfig(ctx context.Context, branchMethod *PaymentMethodsBranch) error
	GetBranchPaymentMethodByID(ctx context.Context, branchID uuid.UUID, methodID uuid.UUID) (*PaymentMethodsBranch, error)
	//Fiscal Documents
	CreateFiscalDocumentType(ctx context.Context, docType *FiscalDocumentType) error
	GetFiscalDocumentTypes(ctx context.Context, fq pagination.PaginatedFeedQuery) ([]*FiscalDocumentType, error)
	GetFiscalDocumentTypesTotal(ctx context.Context, fq pagination.PaginatedFeedQuery) (int64, error)
	GetFiscalDocumentTypeByID(ctx context.Context, docTypeID uuid.UUID) (*FiscalDocumentType, error)
	UpdateFiscalDocumentType(ctx context.Context, docType *FiscalDocumentType) error
	DeleteFiscalDocumentType(ctx context.Context, docTypeID uuid.UUID) error

	// Exchange Rates
	GetLatestRates(ctx context.Context) (map[string]float64, error)
	GetHistoricalRateForCurrency(ctx context.Context, date, currency string) (float64, error)

	//orders
	CreateInvoice(ctx context.Context, invoice *Invoice, items []InvoiceItem) error
}
