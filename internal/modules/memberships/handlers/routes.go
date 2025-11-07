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
	services appservices.Services
}

func NewMembershipHandler(services appservices.Services) *MembershipHandler {
	return &MembershipHandler{
		services: services,
	}
}

func (r *MembershipHandler) MembershipRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware) {
	membershipPlansGroup := rg.Group("/membership-plans")
	{
		membershipPlansGroup.Use(m.AuthJwtTokenMiddleware())

		membershipPlansGroup.POST("", m.RBACPermission("memberships:create"), r.CreateMembershipHandler)
		membershipPlansGroup.GET("", m.RBACPermission("memberships:list"), r.GetMembershipsHandler)
		membershipPlansGroup.GET("/summary", m.RBACPermission("memberships:list"), r.GetSummaryMembershipsHandler)
		membershipPlansGroup.GET("/:id", m.RBACPermission("memberships:get"), r.GetMembershipByIDHandler)
		membershipPlansGroup.PUT("/:id", m.RBACPermission("memberships:update"), r.UpdateMembershipHandler)
		membershipPlansGroup.DELETE("/:id", m.RBACPermission("memberships:delete"), r.DeleteMembershipHandler)
	}
}
