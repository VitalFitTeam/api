package authrepository

import (
	"context"
	"time"

	"github.com/google/uuid"
	authdomain "github.com/vitalfit/api/internal/auth/domain"
	"github.com/vitalfit/api/pkg/db"
	"gorm.io/gorm"
)

type UserRepositoryDAO struct {
	db *gorm.DB
}

func NewUserRepositoryDAO(db *gorm.DB) *UserRepositoryDAO {
	return &UserRepositoryDAO{
		db: db,
	}
}

func (s *UserRepositoryDAO) GetUser() error {
	return nil
}

func (s *UserRepositoryDAO) Create(ctx context.Context, tx *gorm.DB, user authdomain.Users) error {
	return tx.WithContext(ctx).Create(&user).Error
}

func (s *UserRepositoryDAO) createUserInvitation(ctx context.Context, tx *gorm.DB, code string, userID uuid.UUID, invitationExp time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, authdomain.QueryTimeoutDuration)
	defer cancel()
	return tx.WithContext(ctx).Create(&authdomain.UserInvitations{
		Token:  code,
		UserID: userID,
		Expiry: time.Now().Add(invitationExp),
	}).Error
}

func (s *UserRepositoryDAO) CreateAndInvitate(ctx context.Context, user authdomain.Users, token string, invitationExp time.Duration) error {
	//transacction
	return db.WithTX(s.db, func(tx *gorm.DB) error {

		if err := s.Create(ctx, tx, user); err != nil {
			return err //rollback
		}

		if err := s.createUserInvitation(ctx, tx, token, user.UserID, invitationExp); err != nil {
			return err //rollback
		}

		return nil //commit
	})
}
