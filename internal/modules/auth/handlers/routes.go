package authhandlers

import (
	"github.com/gin-gonic/gin"
	appservices "github.com/vitalfit/api/internal/app/services"
	"github.com/vitalfit/api/internal/shared/middleware/auth"
)

type AuthHandlersInterface interface {
	AuthRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware)
	UserRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware)
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

		protectedGroup := authGroup.Group("/").Use(m.AuthJwtTokenMiddleware(), m.CheckRoleAccess("branch_admin"))
		{
			protectedGroup.POST("/register-staff", r.RegisterUserStaffHandler)
		}

	}

}

func (r *AuthHandlers) UserRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware) {
	userGroup := rg.Group("/user").Use(m.AuthJwtTokenMiddleware())
	{ //private routes
		userGroup.GET("/whoami", r.WhoAmI)
		userGroup.GET("/branch-admins", r.GetBranchAdminsHandler).Use(m.CheckRoleAccess("super_admin"))
	}
}
