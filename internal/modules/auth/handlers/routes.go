package authhandlers

import (
	"github.com/gin-gonic/gin"
	appservices "github.com/vitalfit/api/internal/app/services"
	"github.com/vitalfit/api/internal/shared/middleware/auth"
)

type AuthHandlersInterface interface {
	AuthRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware)
	UserRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware)
	AdminRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware)
	//Users
	RegisterUserStaffHandler(c *gin.Context)
	RegisterUserClientHandler(c *gin.Context)
	ActivateUserHandler(c *gin.Context)
	LoginHandler(c *gin.Context)
	WhoAmI(c *gin.Context)
	ForgotPasswordHandler(c *gin.Context)
	ResetPasswordHandler(c *gin.Context)
	GetBranchAdminsHandler(c *gin.Context)
	//Roles
	GetRolesHandler(c *gin.Context)
	CreateRoleHandler(c *gin.Context)
	GetRoleByIDHandler(c *gin.Context)
	UpdateRoleHandler(c *gin.Context)
	DeleteRoleHandler(c *gin.Context)
	GetPermissionsHandler(c *gin.Context)
	AssignRolePermissionHandler(c *gin.Context)
	DeleteRolePermissionHandler(c *gin.Context)
}

type AuthHandlers struct {
	services appservices.Services
}

func NewAuthHandlers(services appservices.Services) *AuthHandlers {
	return &AuthHandlers{services: services}
}

func (r *AuthHandlers) AuthRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware) {

	authGroup := rg.Group("/auth")
	{ //public routes
		authGroup.POST("/register", r.RegisterUserClientHandler)
		authGroup.PUT("/activate", r.ActivateUserHandler)
		authGroup.POST("/login", r.LoginHandler)

		passwordGroup := authGroup.Group("/password")
		{
			passwordGroup.POST("/forgot", r.ForgotPasswordHandler)
			passwordGroup.POST("/reset", r.ResetPasswordHandler)
		}

		protectedGroup := authGroup.Group("/").Use(
			m.AuthJwtTokenMiddleware(),
			m.RBACPermission("users:create"),
		)
		{
			protectedGroup.POST("/register-staff", r.RegisterUserStaffHandler)
		}
	}

}

func (r *AuthHandlers) UserRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware) {
	userGroup := rg.Group("/user").Use(m.AuthJwtTokenMiddleware())
	{ //private routes
		userGroup.GET("/whoami", r.WhoAmI)

		userGroup.GET("/branch-admins",
			m.RBACPermission("users:list"),
			r.GetBranchAdminsHandler,
		)
	}
}

func (r *AuthHandlers) AdminRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware) {

	adminGroup := rg.Group("/admin")

	adminGroup.Use(m.AuthJwtTokenMiddleware())
	{

		adminGroup.GET("/permissions", m.RBACPermission("permissions:list"), r.GetPermissionsHandler)

		roleGroup := adminGroup.Group("/roles")
		{
			// GET /admin/roles
			roleGroup.GET("", m.RBACPermission("roles:list"), r.GetRolesHandler)

			// POST /admin/roles
			roleGroup.POST("", m.RBACPermission("roles:create"), r.CreateRoleHandler)

			// GET /admin/roles/{id}
			roleGroup.GET("/:id", m.RBACPermission("roles:list"), r.GetRoleByIDHandler)

			// PUT /admin/roles/{id}
			roleGroup.PUT("/:id", m.RBACPermission("roles:update"), r.UpdateRoleHandler)

			// DELETE /admin/roles/{id}
			roleGroup.DELETE("/:id", m.RBACPermission("roles:delete"), r.DeleteRoleHandler)

			// POST /admin/roles/{id}/permissions
			roleGroup.POST("/:id/permissions", m.RBACPermission("roles:assign_permissions"), r.AssignRolePermissionHandler)

			// DELETE /admin/roles/{id}/permissions
			roleGroup.DELETE("/:id/permissions", m.RBACPermission("roles:assign_permissions"), r.DeleteRolePermissionHandler)
		}
	}
}
