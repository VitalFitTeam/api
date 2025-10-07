package authrepository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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

func (s *UserRepositoryDAO) Create(ctx context.Context, tx *gorm.DB, user *authdomain.Users) error {
	err := tx.WithContext(ctx).Create(&user).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if pgErr.ConstraintName == "users_email_key" {
				return shared_errors.ErrConflict
			}
			if pgErr.ConstraintName == "users_identity_document_key" {
				return shared_errors.ErrConflict
			}
		}
		return err
	}
	return nil
}

func (s *UserRepositoryDAO) GetByID(ctx context.Context, userID uuid.UUID) (*authdomain.Users, error) {
	var user authdomain.Users

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result := s.db.WithContext(ctx).
		Preload("Role").
		Where("user_id = ?", userID).
		First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, shared_errors.ErrNotFound
		}
		return nil, result.Error
	}
	return &user, nil
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

func (s *UserRepositoryDAO) Activate(ctx context.Context, code string) error {
	return db.WithTX(s.db, func(tx *gorm.DB) error {
		user, err := s.getUserFromInvitation(ctx, tx, code)
		if err != nil {
			return err
		}
		user.IsValidated = true
		if err := s.Update(ctx, user); err != nil {
			return err
		}

		if err := s.deleteUserInvitations(ctx, tx, user.UserID); err != nil {
			return err
		}

		return nil
	})

}

func (s *UserRepositoryDAO) GetByEmail(ctx context.Context, email string) (*authdomain.Users, error) {
	var user authdomain.Users
	ctx, cancel := context.WithTimeout(ctx, db.QueryTimeoutDuration)
	defer cancel()
	err := s.db.WithContext(ctx).Where("email = ?", email).Where("is_validated = ?", true).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared_errors.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserRepositoryDAO) Update(ctx context.Context, user *authdomain.Users) error {
	err := s.db.WithContext(ctx).Save(user).Error
	if err != nil {
		return err
	}
	return nil
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

// func (s *UserRepositoryDAO) softDelete(ctx context.Context, tx *gorm.DB, userID uuid.UUID) error {
// 	ctx, cancel := context.WithTimeout(ctx, db.QueryTimeoutDuration)
// 	defer cancel()

// 	result := tx.WithContext(ctx).Delete(&authdomain.Users{}, userID)

// 	if result.Error != nil {
// 		return result.Error
// 	}
// 	return nil
// }

// Elimina las invitaciones asociadas a ese usuario.
func (s *UserRepositoryDAO) deleteUserInvitations(ctx context.Context, tx *gorm.DB, userID uuid.UUID) error {
	// Elimina todos los registros de invitaciones que tienen este UserID
	result := tx.WithContext(ctx).Where("user_id = ?", userID).Delete(&authdomain.UserInvitations{})

	if result.Error != nil {
		return result.Error
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

func (s *UserRepositoryDAO) getUserFromInvitation(ctx context.Context, tx *gorm.DB, code string) (*authdomain.Users, error) {

	var invitation authdomain.UserInvitations

	hash := sha256.Sum256([]byte(code))
	hashCode := hex.EncodeToString(hash[:])

	ctx, cancel := context.WithTimeout(ctx, db.QueryTimeoutDuration)
	defer cancel()

	result := tx.WithContext(ctx).
		Preload("Users").
		Where("token = ? AND expiry > ?", hashCode, time.Now()).
		First(&invitation)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, shared_errors.ErrNotFound
		}
		return nil, result.Error
	}

	return &invitation.Users, nil
}

func (s *UserRepositoryDAO) CreatePasswordResetToken(ctx context.Context, userID uuid.UUID, key string, tokenExp time.Duration) error {
	return db.WithTX(s.db, func(tx *gorm.DB) error {
		if err := s.userResetToken(ctx, tx, userID, key, tokenExp); err != nil {
			return err //rollback
		}
		return nil //commit

	})
}

func (s *UserRepositoryDAO) DeleteResetToken(ctx context.Context, userID uuid.UUID) error {
	return db.WithTX(s.db, func(tx *gorm.DB) error {
		if err := s.deleteUserReset(ctx, tx, userID); err != nil {
			return err //rollback
		}
		return nil //commit
	})
}

func (s *UserRepositoryDAO) userResetToken(ctx context.Context, tx *gorm.DB, userID uuid.UUID, key string, tokenExp time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, db.QueryTimeoutDuration)
	defer cancel()

	err := tx.WithContext(ctx).Create(&authdomain.PasswordResetToken{
		Token:  key,
		UserID: userID,
		Expiry: time.Now().Add(tokenExp),
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *UserRepositoryDAO) deleteUserReset(ctx context.Context, tx *gorm.DB, userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, db.QueryTimeoutDuration)
	defer cancel()
	err := s.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&authdomain.PasswordResetToken{}).Error
	if err != nil {
		return err
	}
	return nil
}
