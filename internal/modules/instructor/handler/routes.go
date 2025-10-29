package instructorhandler

import appservices "github.com/vitalfit/api/internal/app/services"

type InstructorHandlersInterface interface {
}

type InstructorHandlers struct {
	services appservices.Services
}

func NewInstructorHandlers(services appservices.Services) *InstructorHandlers {
	return &InstructorHandlers{services: services}
}
