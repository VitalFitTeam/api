package authhandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/vitalfit/api/internal/shared/middleware/auth"
)

func (r *AuthHandlers) AuthRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware) {

	authGroup := rg.Group("/auth")
	{ //public routes
		authGroup.POST("/register", r.registerUserClientHandler)
		authGroup.PUT("/activate", r.activateUserHandler)
		authGroup.POST("/login", r.loginHandler)

		passwordGroup := authGroup.Group("/password")
		{
			passwordGroup.POST("/forgot", r.forgotPasswordHandler)
			// passwordGroup.POST("/reset", r.resetPasswordHandler)
		}

		protectedGroup := authGroup.Group("/").Use(m.AuthJwtTokenMiddleware(), m.CheckRoleAccess("branch_admin"))
		{
			protectedGroup.POST("/register-staff", r.registerUserStaffHandler)
		}

	}

}

func (r *AuthHandlers) UserRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware) {
	userGroup := rg.Group("/user").Use(m.AuthJwtTokenMiddleware())
	{ //private routes
		userGroup.GET("/whoami", r.whoami)
	}
}
