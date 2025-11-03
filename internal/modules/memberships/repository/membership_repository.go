package membershipsrepository

import "gorm.io/gorm"

type MembershipStore struct {
	db *gorm.DB
}

func NewMembershipStore(db *gorm.DB) *MembershipStore {
	return &MembershipStore{
		db: db,
	}
}
