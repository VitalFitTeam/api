package accessservice

import "github.com/vitalfit/api/internal/store"

type AccessService struct {
	store store.Storage
}

func NewAccessServices(store store.Storage) *AccessService {
	return &AccessService{
		store: store,
	}
}
