package authhandlers

import "github.com/gin-gonic/gin"

// RegisterRoutes define todas las rutas específicas del módulo Auth.
func (r *AuthHandlers) AuthRoutes(rg *gin.RouterGroup) {
	// 1. Define el grupo de rutas para /v1/auth
	authGroup := rg.Group("/auth")
	{
		// 2. Asocia los métodos del Handler a las rutas
		authGroup.POST("/register", r.RegisterUserClientHandler)
		authGroup.POST("/register-staff", r.RegisterUserStaffHandler)
		//authGroup.POST("/login", r.Handler.LoginHandler)
		//authGroup.POST("/logout", r.Handler.LogoutHandler)
		// ...
	}
}
