package policiesdomain

import (
	"context"
)

type PoliciesRepository interface {
	CreatePolicy(ctx context.Context, policy *CommercialPolicy) error
	GetPolicyByKey(ctx context.Context, key string) (*CommercialPolicy, error)
	UpdatePolicy(ctx context.Context, policy *CommercialPolicy) error
	GetPoliciesList(ctx context.Context, policy *CommercialPolicy) ([]CommercialPolicy, error)
}
