package authrepository

import (
	"context"
	"time"

	"github.com/google/uuid"
	authdomain "github.com/vitalfit/api/internal/modules/auth/domain"
	"gorm.io/gorm"
)

type SessionStore struct {
	db *gorm.DB
}

func NewSessionStore(db *gorm.DB) *SessionStore {
	return &SessionStore{db: db}
}

func (s *SessionStore) Create(ctx context.Context, session *authdomain.Session) error {
	return s.db.WithContext(ctx).Create(session).Error

}

func (s *SessionStore) GetByRefreshToken(refreshToken string) (*authdomain.Session, error) {
	var session authdomain.Session
	if err := s.db.Where("refresh_token = ? AND is_blocked = ? AND expires_at > ?",
		refreshToken, false, time.Now()).First(&session).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, nil
		default:
			return nil, err
		}
	}
	return &session, nil

}

func (s *SessionStore) Revoke(sessionID uuid.UUID) error {
	return s.db.Model(&authdomain.Session{}).
		Where("id = ?", sessionID).
		Update("is_blocked", true).Error
}

func (s *SessionStore) RevokeAllForUser(userID uuid.UUID) error {
	return s.db.Model(&authdomain.Session{}).
		Where("user_id = ? AND is_blocked = ?", userID, false).
		Update("is_blocked", true).Error
}
