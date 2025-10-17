package branchhandlers

import appservices "github.com/vitalfit/api/internal/app/services"

type BranchHandlersInterface interface {
}

type BranchHandlers struct {
	services appservices.Services
}

func NewBranchHandlers(services appservices.Services) *BranchHandlers {
	return &BranchHandlers{services: services}
}
