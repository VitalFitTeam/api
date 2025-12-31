package policiesdomain

import (
	"context"

	"gorm.io/gorm"
)

type PoliciesRepository interface {
	CreatePolicyTx(ctx context.Context, tx *gorm.DB, policy *CommercialPolicy) error
	CreatePolicy(ctx context.Context, policy *CommercialPolicy) error
}
