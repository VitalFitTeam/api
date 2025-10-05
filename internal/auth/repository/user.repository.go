package authrepository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	authdomain "github.com/vitalfit/api/internal/auth/domain"
	shared_errors "github.com/vitalfit/api/internal/shared/errors"
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

func (s *UserRepositoryDAO) Create(ctx context.Context, tx *gorm.DB, user *authdomain.Users) error {
	err := tx.WithContext(ctx).Create(&user).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if pgErr.ConstraintName == "users_email_key" {
				return shared_errors.ErrConflict
			}
		}
		return err
	}
	return nil
}

func (s *UserRepositoryDAO) createUserInvitation(ctx context.Context, tx *gorm.DB, code string, userID uuid.UUID, invitationExp time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, db.QueryTimeoutDuration)
	defer cancel()
	return tx.WithContext(ctx).Create(&authdomain.UserInvitations{
		Token:  code,
		UserID: userID,
		Expiry: time.Now().Add(invitationExp),
	}).Error
}

func (s *UserRepositoryDAO) CreateAndInvitate(ctx context.Context, user *authdomain.Users, token string, invitationExp time.Duration) error {
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

func (s *UserRepositoryDAO) delete(ctx context.Context, tx *gorm.DB, userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, db.QueryTimeoutDuration)
	defer cancel()

	result := tx.WithContext(ctx).Unscoped().Delete(&authdomain.Users{}, userID)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *UserRepositoryDAO) softDelete(ctx context.Context, tx *gorm.DB, userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, db.QueryTimeoutDuration)
	defer cancel()

	result := tx.WithContext(ctx).Delete(&authdomain.Users{}, userID)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Elimina las invitaciones asociadas a ese usuario.
func (s *UserRepositoryDAO) deleteUserInvitations(ctx context.Context, tx *gorm.DB, userID uuid.UUID) error {
	// Elimina todos los registros de invitaciones que tienen este UserID
	result := tx.WithContext(ctx).Where("user_id = ?", userID).Delete(&authdomain.UserInvitations{})

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *UserRepositoryDAO) Delete(ctx context.Context, userID uuid.UUID) error {
	return db.WithTX(s.db, func(tx *gorm.DB) error {
		if err := s.delete(ctx, tx, userID); err != nil {
			return err //rollback
		}
		if err := s.deleteUserInvitations(ctx, tx, userID); err != nil {
			return err //rollback
		}
		return nil //commit
	})
}
