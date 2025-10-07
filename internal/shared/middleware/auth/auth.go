package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	appservices "github.com/vitalfit/api/internal/app/services"
	authdomain "github.com/vitalfit/api/internal/auth/domain"
)

type AuthMiddleware struct {
	services appservices.Services
}

func NewAuthMiddleware(services appservices.Services) *AuthMiddleware {
	return &AuthMiddleware{
		services: services,
	}
}

func (j *AuthMiddleware) AuthJwtTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			j.services.LogErrors.UnauthorizedErrorResponse(c, fmt.Errorf("missing authorization header"))
			c.Abort()
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			j.services.LogErrors.UnauthorizedErrorResponse(c, fmt.Errorf("authorization header is malformed"))
			c.Abort()
			return
		}

		token := parts[1]
		jwtToken, err := j.services.AuthServices.ValidateToken(token)
		if err != nil {
			j.services.LogErrors.UnauthorizedErrorResponse(c, err)
			c.Abort()
			return
		}

		claims, _ := jwtToken.Claims.(jwt.MapClaims)

		userID, err := uuid.Parse(claims["sub"].(string))
		if err != nil {
			j.services.LogErrors.UnauthorizedErrorResponse(c, err)
			c.Abort()
			return
		}

		ctx := c.Request.Context()

		user, err := j.services.UserServices.GetByID(ctx, userID)
		if err != nil {
			j.services.LogErrors.UnauthorizedErrorResponse(c, err)
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	}
}

// checks role access to the endpoint
func (j *AuthMiddleware) CheckRoleAccess(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := j.services.UserServices.GetUserFromContext(c)

		if user == nil {
			j.services.LogErrors.UnauthorizedErrorResponse(c, fmt.Errorf("user or user role not found in context for role check"))
			c.Abort()
			return
		}

		// The rest of the logic is correct
		allowed, err := j.CheckRolePrecedence(c.Request.Context(), user, requiredRole)
		if err != nil {
			j.services.LogErrors.UnauthorizedErrorResponse(c, err)
			c.Abort()
			return
		}

		if !allowed {
			j.services.LogErrors.ForbiddenResponse(c)
			c.Abort()
			return
		}
		c.Next()
	}
}

// compares users level with the level required
func (j *AuthMiddleware) CheckRolePrecedence(ctx context.Context, user *authdomain.Users, roleName string) (bool, error) {
	role, err := j.services.UserServices.GetRoleByName(ctx, roleName)
	if err != nil {
		return false, err
	}
	return user.Role.Level >= role.Level, nil
}
