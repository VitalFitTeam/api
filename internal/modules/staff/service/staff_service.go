package staffservice

import "github.com/vitalfit/api/internal/store"

type StaffService struct {
	store store.Storage
}

func NewStaffService(store store.Storage) *StaffService {
	return &StaffService{store: store}
}
