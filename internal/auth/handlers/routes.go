package authhandlers

import (
	"github.com/gin-gonic/gin"
	shared_jwt "github.com/vitalfit/api/internal/shared/middleware/auth"
)

func (r *AuthHandlers) AuthRoutes(rg *gin.RouterGroup) {

	authGroup := rg.Group("/auth")
	{ //public routes
		authGroup.POST("/register", r.RegisterUserClientHandler)
		authGroup.PUT("/activate", r.ActivateUserHandler)
		authGroup.POST("/login", r.LoginHandler)
	}
}

func (r *AuthHandlers) UserRoutes(rg *gin.RouterGroup) {
	m := shared_jwt.NewAuthMiddleware(r.services)
	userGroup := rg.Group("/user").Use(m.AuthJwtTokenMiddleware(), m.CheckRoleAccess("super_admin"))
	{ //private routes
		userGroup.GET("/whoami", r.whoami)
		userGroup.POST("/register-staff", r.RegisterUserStaffHandler)
	}
}
