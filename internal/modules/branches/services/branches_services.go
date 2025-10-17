package branchservices

import (
	"github.com/vitalfit/api/config"
	"github.com/vitalfit/api/internal/store"
)

type BranchServices struct {
	store  store.Storage
	config config.Config
}

func NewBranchServices(store store.Storage, cfg config.Config) *BranchServices {
	return &BranchServices{
		store:  store,
		config: cfg,
	}

}
