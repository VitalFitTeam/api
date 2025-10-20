package branchdomain

import (
	"context"
)

type BranchesRepository interface {
	//GetByID(ctx context.Context, branchID uuid.UUID) (*Branch, error)
	//Update(ctx context.Context, branch *Branch) error
	CreateBranch(ctx context.Context, branch *Branch) error
}

type LocationRepository interface {
	FindOrCreateStateByCountry(ctx context.Context, stateName string, countryName string) (*States, error)
}

type PaymentMethodsRepository interface {
	GetPaymentMethods(ctx context.Context) ([]*PaymentMethods, error)
	GetPaymentMethodByName(ctx context.Context, name string) (*PaymentMethods, error)
}
