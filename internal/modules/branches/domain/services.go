package branchdomain

import (
	"context"
)

type BranchesServicesInterface interface {
	CreateBranch(ctx context.Context, branch *Branch) error
	GetPaymentMethodByName(ctx context.Context, name string) (*PaymentMethods, error)
	GetPaymentMethods(ctx context.Context) ([]*PaymentMethods, error)
}

type LocationsServicesInterface interface {
	FindOrCreateStateByCountry(ctx context.Context, stateName string, countryName string) (*States, error)
}
