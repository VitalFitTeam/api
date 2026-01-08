package authservices

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	authdomain "github.com/vitalfit/api/internal/modules/auth/domain"
	"github.com/vitalfit/api/pkg/otp"
)

func (h *AuthService) CreateUserSession(ctx context.Context, user *authdomain.Users, userAgent, string, clientIP string) (string, string, error) {
	claims := jwt.MapClaims{
		"sub": user.UserID,
		"exp": time.Now().Add(h.config.Auth.Token.AccessExp).Unix(), // Usar config de Access Token
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": h.config.Auth.Token.Iss,
		"aud": h.config.Auth.Token.Aud,
	}

	accessToken, err := h.auth.GenerateToken(claims)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := otp.GenerateRandomString()
	if err != nil {
		return "", "", err
	}

	session := &authdomain.Session{
		UserID:       user.UserID,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientIP:     clientIP,
		IsBlocked:    false,
		ExpiresAt:    time.Now().Add(h.config.Auth.Token.RefreshExp),
	}

	if err := h.store.Session.Create(ctx, session); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
