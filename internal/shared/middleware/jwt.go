package jwt

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	appservices "github.com/vitalfit/api/internal/app/services"
)

type jwtAuthMiddleware struct {
	services appservices.Services
}

func NewJWTAuthMiddleware(services appservices.Services) *jwtAuthMiddleware {
	return &jwtAuthMiddleware{
		services: services,
	}
}

func (j *jwtAuthMiddleware) AuthJwtTokenMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			j.services.LogErrors.UnauthorizedErrorResponse(c, fmt.Errorf("missing authorization header"))
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			j.services.LogErrors.UnauthorizedErrorResponse(c, fmt.Errorf("authorization header is malformed"))
			return
		}

		token := parts[1]
		jwtToken, err := j.services.AuthServices.ValidateToken(token)
		if err != nil {
			j.services.LogErrors.UnauthorizedErrorResponse(c, err)
			return
		}

		claims, _ := jwtToken.Claims.(jwt.MapClaims)

		userID, err := uuid.Parse(claims["sub"].(string))
		if err != nil {
			j.services.LogErrors.UnauthorizedErrorResponse(c, err)
			return
		}

		ctx := c.Request.Context()

		user, err := j.services.UserServices.GetByID(ctx, userID)
		if err != nil {
			j.services.LogErrors.UnauthorizedErrorResponse(c, err)
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
