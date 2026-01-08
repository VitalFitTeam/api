package authrepository

import (
	"context"
	"errors"
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

func (s *SessionStore) GetByRefreshToken(ctx context.Context, refreshToken string) (*authdomain.Session, error) {
	var session authdomain.Session
	if err := s.db.WithContext(ctx).Where("refresh_token = ? AND is_blocked = ? AND expires_at > ?",
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

func (s *SessionStore) GetUserSessions(ctx context.Context, userID uuid.UUID) ([]*authdomain.Session, error) {
	var sessions []*authdomain.Session
	if err := s.db.WithContext(ctx).Where("user_id = ? AND is_blocked = ?", userID, false).
		Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

func (s *SessionStore) RotateSession(ctx context.Context, sessionID uuid.UUID, oldToken, newToken string) error {
	result := s.db.WithContext(ctx).Model(&authdomain.Session{}).
		Where("id = ? AND refresh_token = ?", sessionID, oldToken).
		Update("refresh_token", newToken)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		// s.db.Model(&authdomain.Session{}).Where("id = ?", sessionID).Update("is_blocked", true)
		return errors.New("refresh token mismatch: reuse detection")
	}

	return nil
}

func (s *SessionStore) Revoke(ctx context.Context, sessionID uuid.UUID) error {
	return s.db.WithContext(ctx).Model(&authdomain.Session{}).
		Where("id = ?", sessionID).
		Update("is_blocked", true).Error
}

func (s *SessionStore) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	return s.db.WithContext(ctx).Model(&authdomain.Session{}).
		Where("user_id = ? AND is_blocked = ?", userID, false).
		Update("is_blocked", true).Error
}
