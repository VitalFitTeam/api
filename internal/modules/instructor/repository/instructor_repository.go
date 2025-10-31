package instructorrepository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	authdomain "github.com/vitalfit/api/internal/modules/auth/domain"
	instructordomain "github.com/vitalfit/api/internal/modules/instructor/domain"
	shared_errors "github.com/vitalfit/api/internal/shared/errors"
	"github.com/vitalfit/api/pkg/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InstructorStore struct {
	db *gorm.DB
}

func NewInstructorStore(db *gorm.DB) *InstructorStore {
	return &InstructorStore{db: db}
}

func (s *InstructorStore) Create(ctx context.Context, tx *gorm.DB, instructor *instructordomain.Instructor) error {

	if err := tx.WithContext(ctx).Create(&instructor.User).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return shared_errors.ErrConflict
		}
		return err
	}

	instructor.UserID = instructor.User.UserID

	if err := tx.WithContext(ctx).Omit("User").Create(instructor).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return shared_errors.ErrConflict
		}
		return err
	}
	return nil
}

func (s *InstructorStore) GetInstructors(ctx context.Context) ([]*instructordomain.Instructor, error) {
	ctx, cancel := context.WithTimeout(ctx, db.QueryTimeoutDuration)
	defer cancel()
	var instructors []*instructordomain.Instructor
	err := s.db.WithContext(ctx).
		Joins("JOIN users ON users.user_id = instructors.user_id").
		Where("users.is_validated = ?", true).
		Preload("User").Preload("Specialties").
		Find(&instructors).Error
	if err != nil {
		return nil, err
	}
	return instructors, nil
}

func (s *InstructorStore) Delete(ctx context.Context, instructorID uuid.UUID) error {
	return db.WithTX(s.db, func(tx *gorm.DB) error {
		var instructor instructordomain.Instructor
		if err := tx.WithContext(ctx).Where("instructor_id = ?", instructorID).First(&instructor).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return shared_errors.ErrNotFound
			}
			return err
		}

		if err := tx.WithContext(ctx).Delete(&instructordomain.Instructor{}, instructorID).Error; err != nil {
			return err
		}

		if err := tx.WithContext(ctx).Delete(&authdomain.Users{}, instructor.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return err
		}

		return nil
	})
}

func (s *InstructorStore) GetByID(ctx context.Context, instructorID uuid.UUID) (*instructordomain.Instructor, error) {
	ctx, cancel := context.WithTimeout(ctx, db.QueryTimeoutDuration)
	defer cancel()
	var instructor *instructordomain.Instructor
	err := s.db.WithContext(ctx).Preload("User").Preload("Specialties").Where("instructor_id = ?", instructorID).First(&instructor).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared_errors.ErrNotFound
		}
		return nil, err
	}
	return instructor, nil
}

func (s *InstructorStore) Update(ctx context.Context, instructor *instructordomain.Instructor) error {
	return db.WithTX(s.db, func(tx *gorm.DB) error {
		var existingInstructor instructordomain.Instructor
		if err := tx.WithContext(ctx).Where("instructor_id = ?", instructor.InstructorID).First(&existingInstructor).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return shared_errors.ErrNotFound
			}
			return err
		}
		userUpdates := make(map[string]interface{})
		if instructor.User.FirstName != "" {
			userUpdates["first_name"] = instructor.User.FirstName
		}
		if instructor.User.LastName != "" {
			userUpdates["last_name"] = instructor.User.LastName
		}
		if instructor.User.Email != "" {
			userUpdates["email"] = instructor.User.Email
		}
		if instructor.User.Phone != "" {
			userUpdates["phone"] = instructor.User.Phone
		}
		if instructor.User.Gender != "" {
			userUpdates["gender"] = instructor.User.Gender
		}
		if !instructor.User.BirthDate.IsZero() {
			userUpdates["birth_date"] = instructor.User.BirthDate
		}
		if instructor.User.ProfilePictureURL != "" {
			userUpdates["profile_picture_url"] = instructor.User.ProfilePictureURL
		}

		if err := tx.WithContext(ctx).Model(instructor).Omit("User").Updates(instructor).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *InstructorStore) CreateAndInvitate(ctx context.Context, instructor *instructordomain.Instructor, token string, invitationExp time.Duration) error {
	//transacction
	return db.WithTX(s.db, func(tx *gorm.DB) error {

		if err := s.Create(ctx, tx, instructor); err != nil {
			return err //rollback
		}

		if err := s.createUserInvitation(ctx, tx, token, instructor.UserID, invitationExp); err != nil {
			return err //rollback
		}

		return nil //commit
	})
}

func (s *InstructorStore) createUserInvitation(ctx context.Context, tx *gorm.DB, code string, userID uuid.UUID, invitationExp time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, db.QueryTimeoutDuration)
	defer cancel()
	return tx.WithContext(ctx).Create(&authdomain.UserInvitations{
		Token:  code,
		UserID: userID,
		Expiry: time.Now().Add(invitationExp),
	}).Error
}

func (s *InstructorStore) AssignInstructorSpecialtyTx(ctx context.Context, tx *gorm.DB, instructorID uuid.UUID, specialityID []uuid.UUID) error {
	var linksToCreate []instructordomain.InstructorSpecialty
	for _, pid := range specialityID {
		linksToCreate = append(linksToCreate, instructordomain.InstructorSpecialty{
			InstructorID: instructorID,
			CategoryID:   pid,
		})
	}
	if len(linksToCreate) == 0 {
		return nil
	}
	err := tx.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&linksToCreate).Error
	if err != nil {
		switch err {
		case shared_errors.ErrConflict:
			return shared_errors.ErrConflict
		case shared_errors.ErrNotFound:
			return shared_errors.ErrNotFound
		default:
			return err
		}
	}
	return nil
}

func (s *InstructorStore) AssignInstructorSpecialty(ctx context.Context, instructorID uuid.UUID, specialityID []uuid.UUID) error {
	return db.WithTX(s.db, func(tx *gorm.DB) error {
		return s.AssignInstructorSpecialtyTx(ctx, tx, instructorID, specialityID)
	})
}
func (s *InstructorStore) DeleteInstructorSpecialty(ctx context.Context, instructorID uuid.UUID, specialtyID uuid.UUID) error {
	return db.WithTX(s.db, func(tx *gorm.DB) error {
		result := tx.WithContext(ctx).
			Where("instructor_id = ? AND category_id = ?", instructorID, specialtyID).
			Delete(&instructordomain.InstructorSpecialty{})

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return shared_errors.ErrNotFound
		}

		return nil
	})
}
