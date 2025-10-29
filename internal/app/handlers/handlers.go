package apphandlers

import (
	appservices "github.com/vitalfit/api/internal/app/services"
	authhandlers "github.com/vitalfit/api/internal/modules/auth/handlers"
	branchhandlers "github.com/vitalfit/api/internal/modules/branches/handlers"
	instructorhandler "github.com/vitalfit/api/internal/modules/instructor/handler"
)

type Handlers struct {
	AuthHandlers       authhandlers.AuthHandlersInterface
	BranchHandlers     branchhandlers.BranchHandlersInterface
	InstructorHandlers instructorhandler.InstructorHandlersInterface
}

func NewAppHandlers(services appservices.Services) Handlers {
	return Handlers{
		AuthHandlers:       authhandlers.NewAuthHandlers(services),
		BranchHandlers:     branchhandlers.NewBranchHandlers(services),
		InstructorHandlers: instructorhandler.NewInstructorHandlers(services),
	}
}
