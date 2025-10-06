package authhandlers

import (
	"github.com/gin-gonic/gin"
	shared_jwt "github.com/vitalfit/api/internal/shared/middleware"
)

func (r *AuthHandlers) AuthRoutes(rg *gin.RouterGroup) {

	authGroup := rg.Group("/auth")
	{ //public routes
		authGroup.POST("/register", r.RegisterUserClientHandler)
		authGroup.POST("/register-staff", r.RegisterUserStaffHandler)
		authGroup.PUT("/activate", r.ActivateUserHandler)
		authGroup.POST("/login", r.LoginHandler)

	}
}

func (r *AuthHandlers) UserRoutes(rg *gin.RouterGroup) {
	m := shared_jwt.NewJWTAuthMiddleware(r.services)
	userGroup := rg.Group("/user").Use(m.AuthJwtTokenMiddleware())
	{
		userGroup.GET("/whoami", r.whoami)
	}
}
