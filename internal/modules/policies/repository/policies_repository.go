package policiesrepository

import "gorm.io/gorm"

type PoliciesStore struct {
	db *gorm.DB
}

func NewPoliciesStore(db *gorm.DB) *PoliciesStore {
	return &PoliciesStore{
		db: db,
	}

}
