package authhandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/vitalfit/api/internal/shared/middleware/auth"
)

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
	}
}
