package policiesrepository

import (
	"context"

	policiesdomain "github.com/vitalfit/api/internal/modules/policies/domain"
	"gorm.io/gorm"
)

type PoliciesStore struct {
	db *gorm.DB
}

func NewPoliciesStore(db *gorm.DB) *PoliciesStore {
	return &PoliciesStore{
		db: db,
	}

}

func (s *PoliciesStore) CreatePolicyTx(ctx context.Context, tx *gorm.DB, policy *policiesdomain.CommercialPolicy) error {

	return nil
}

func (s *PoliciesStore) CreatePolicy(ctx context.Context, policy *policiesdomain.CommercialPolicy) error {

	return nil
}
