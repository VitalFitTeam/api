package branchdomain

import (
	"context"
)

type BranchesServicesInterface interface {
	CreateBranch(ctx context.Context, branch *Branch) error
}

type LocationsServicesInterface interface {
	FindOrCreateStateByCountry(ctx context.Context, stateName string, countryName string) (*States, error)
}
