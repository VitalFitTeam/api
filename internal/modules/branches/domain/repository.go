package branchdomain

import (
	"context"

	"github.com/google/uuid"
	"github.com/vitalfit/api/pkg/pagination"
)

type BranchesRepository interface {
	GetByID(ctx context.Context, branchID uuid.UUID) (*Branch, error)
	Update(ctx context.Context, branch *Branch) error
	CreateBranch(ctx context.Context, branch *Branch) (*Branch, error)
	GetBranches(ctx context.Context, fq pagination.PaginatedFeedQuery) (*BranchQueryResults, error)
	AddPaymentMethodsToBranch(ctx context.Context, branchID uuid.UUID, paymentLinks []PaymentMethodsBranch) error
	Delete(ctx context.Context, branchID uuid.UUID) error
	GetBranchStatusCount(ctx context.Context) (*BranchStatusCount, error)
}

type LocationRepository interface {
	FindOrCreateStateByCountry(ctx context.Context, stateName string, countryName string) (*States, error)
}

type PaymentMethodsRepository interface {
	GetPaymentMethods(ctx context.Context) ([]*PaymentMethods, error)
	GetPaymentMethodByID(ctx context.Context, methodID uuid.UUID) (*PaymentMethods, error)
}
