package combosservices

import "github.com/vitalfit/api/internal/store"

type CombosServices struct {
	store store.Storage
}

func NewCombosServices(store store.Storage) *CombosServices {
	return &CombosServices{store: store}
}
