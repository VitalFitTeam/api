package branchdomain

import (
	"context"

	"github.com/google/uuid"
)

type BranchesServicesInterface interface {
	CreateBranch(ctx context.Context, branch *Branch) error
	GetPaymentMethodByID(ctx context.Context, methodID uuid.UUID) (*PaymentMethods, error)
	GetPaymentMethods(ctx context.Context) ([]*PaymentMethods, error)
}

type LocationsServicesInterface interface {
	FindOrCreateStateByCountry(ctx context.Context, stateName string, countryName string) (*States, error)
}
