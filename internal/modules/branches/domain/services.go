package branchdomain

import (
	"context"

	"github.com/google/uuid"
	"github.com/vitalfit/api/pkg/pagination"
)

type BranchesServicesInterface interface {
	CreateBranch(ctx context.Context, branch *Branch) error
	GetPaymentMethodByID(ctx context.Context, methodID uuid.UUID) (*PaymentMethods, error)
	GetPaymentMethods(ctx context.Context) ([]*PaymentMethods, error)
	GetBranches(ctx context.Context, fq pagination.PaginatedFeedQuery) (*BranchQueryResults, error)
	DeleteBranch(ctx context.Context, branchID uuid.UUID) error
}

type LocationsServicesInterface interface {
	FindOrCreateStateByCountry(ctx context.Context, stateName string, countryName string) (*States, error)
}
