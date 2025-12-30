package policiesservices

import "github.com/vitalfit/api/internal/store"

type PoliciesServices struct {
	store store.Storage
}

func NewPoliciesServices(store store.Storage) *PoliciesServices {
	return &PoliciesServices{
		store: store,
	}
}
