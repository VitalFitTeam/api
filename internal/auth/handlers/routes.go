package authhandlers

import "github.com/gin-gonic/gin"

func (r *AuthHandlers) AuthRoutes(rg *gin.RouterGroup) {

	authGroup := rg.Group("/auth")
	{
		// 2. Asocia los métodos del Handler a las rutas
		authGroup.POST("/register", r.RegisterUserClientHandler)
		authGroup.POST("/register-staff", r.RegisterUserStaffHandler)
		authGroup.PUT("/activate", r.ActivateUserHandler)
		authGroup.POST("/login", r.LoginHandler)
		//authGroup.POST("/logout", r.Handler.LogoutHandler)
		// ...
	}
}
