package staffhandlers

import appservices "github.com/vitalfit/api/internal/app/services"

type StaffHandlersInterface interface {
}

type StaffHandlers struct {
	services appservices.Services
}

func NewStaffHandlers(services appservices.Services) *StaffHandlers {
	return &StaffHandlers{services: services}
}
