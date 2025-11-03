package membershipsservice

import "github.com/vitalfit/api/internal/store"

type MembershipService struct {
	store store.Storage
}

func NewMembershipService(store store.Storage) *MembershipService {
	return &MembershipService{
		store: store,
	}
}
