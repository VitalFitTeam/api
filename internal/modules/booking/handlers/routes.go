package bookinghandlers

import (
	"github.com/gin-gonic/gin"
	appservices "github.com/vitalfit/api/internal/app/services"
	"github.com/vitalfit/api/internal/shared/middleware/auth"
)

type BookingHandlersInterface interface {
	BookingRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware)

	GetClientScheduleHandler(c *gin.Context)
	CreateBookingHandler(c *gin.Context)
	CancelBookingHandler(c *gin.Context)
}

type BookingHandlers struct {
	services appservices.Services
}

func NewBookingHandlers(services appservices.Services) *BookingHandlers {
	return &BookingHandlers{services: services}
}

func (h *BookingHandlers) BookingRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware) {

	scheduleRoutes := rg.Group("/schedule")
	{
		scheduleRoutes.Use(m.AuthJwtTokenMiddleware())
		scheduleRoutes.GET("", m.RBACPermission("booking:schedule:list"), h.GetClientScheduleHandler)
		scheduleRoutes.POST("/:classId/book", m.RBACPermission("booking:create"), h.CreateBookingHandler)
	}

	bookingRoutes := rg.Group("/bookings")
	{
		bookingRoutes.Use(m.AuthJwtTokenMiddleware())
		bookingRoutes.DELETE("/:bookingId", m.RBACPermission("booking:cancel"), h.CancelBookingHandler)
	}
}
