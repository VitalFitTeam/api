package membershipshandlers

import (
	"github.com/gin-gonic/gin"
	appservices "github.com/vitalfit/api/internal/app/services"
	"github.com/vitalfit/api/internal/shared/middleware/auth"
)

type MembershipsHandlerInterface interface {
	MembershipRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware)
}

type MembershipHandler struct {
	service appservices.Services
}

func NewMembershipHandler(service appservices.Services) *MembershipHandler {
	return &MembershipHandler{
		service: service,
	}
}

func (r *MembershipHandler) MembershipRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware) {

}
