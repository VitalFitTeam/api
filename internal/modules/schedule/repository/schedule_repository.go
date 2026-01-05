package schedulerepository

import (
	"context"
	"time"

	"github.com/google/uuid"
	bookingdomain "github.com/vitalfit/api/internal/modules/booking/domain"
	scheduledomain "github.com/vitalfit/api/internal/modules/schedule/domain"
	"github.com/vitalfit/api/pkg/db"
	"gorm.io/gorm"
)

type ScheduleStore struct {
	db *gorm.DB
}

func NewScheduleStore(db *gorm.DB) *ScheduleStore {
	return &ScheduleStore{db: db}
}

// ----------------------------------------
// CreateClass
// ----------------------------------------

func (s *ScheduleStore) CreateClass(ctx context.Context, class *scheduledomain.Class) error {
	err := db.WithTX(s.db, func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Create(class).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// ----------------------------------------
// CreateClasses (Batch)
// ----------------------------------------

func (s *ScheduleStore) CreateClasses(ctx context.Context, classes []scheduledomain.Class) error {
	return db.WithTX(s.db, func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Create(&classes).Error; err != nil {
			return err
		}
		return nil
	})
}

// ----------------------------------------
// GetClassesByBranch
// ----------------------------------------

func (s *ScheduleStore) GetClassesByBranch(ctx context.Context, branchID uuid.UUID) ([]scheduledomain.Class, error) {
	var classes []scheduledomain.Class

	err := db.WithTX(s.db, func(tx *gorm.DB) error {

		if err := tx.WithContext(ctx).
			Preload("Service").
			Preload("Instructor").
			Preload("Branch").
			Where("branch_id = ?", branchID).
			Find(&classes).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return classes, nil
}

// ----------------------------------------
// GetUpcomingClassesByBranch
// ----------------------------------------

func (s *ScheduleStore) GetUpcomingClassesByBranch(ctx context.Context, branchID uuid.UUID) ([]scheduledomain.Class, error) {
	var classes []scheduledomain.Class

	err := db.WithTX(s.db, func(tx *gorm.DB) error {

		if err := tx.WithContext(ctx).
			Preload("Service").
			Preload("Instructor").
			Preload("Branch").
			Where("branch_id = ?", branchID).
			Where("ends_at > ?", time.Now()).
			Find(&classes).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return classes, nil
}

// ----------------------------------------
// GetClassByID
// ----------------------------------------

func (s *ScheduleStore) GetClassByID(ctx context.Context, classID uuid.UUID) (*scheduledomain.Class, error) {
	var class scheduledomain.Class

	err := db.WithTX(s.db, func(tx *gorm.DB) error {

		if err := tx.WithContext(ctx).
			Preload("Service").
			Preload("Instructor").
			Preload("Branch").
			First(&class, "class_id = ?", classID).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, gorm.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &class, nil
}

// ----------------------------------------
// UpdateClass
// ----------------------------------------

func (s *ScheduleStore) UpdateClass(ctx context.Context, class *scheduledomain.Class) error {

	result := s.db.WithContext(ctx).
		Model(&scheduledomain.Class{}).
		Where("class_id = ?", class.ClassID).
		Updates(class)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *ScheduleStore) GetAvailableClassesForBranch(ctx context.Context, branchID uuid.UUID, startTime, endTime time.Time) ([]scheduledomain.Class, error) {
	var classes []scheduledomain.Class
	bookingCountSubquery := s.db.Model(&bookingdomain.Booking{}).
		Select("COALESCE(COUNT(booking_id), 0)").
		Where("class_id = classes.class_id").
		Where("status = 'Confirmed'").
		Where("deleted_at IS NULL")

	err := s.db.WithContext(ctx).
		Model(&scheduledomain.Class{}).
		Preload("Service").
		Where("branch_id = ?", branchID).
		Where("is_visible = ?", true).
		Where("starts_at BETWEEN ? AND ?", startTime, endTime).
		Where("max_capacity > (?)", bookingCountSubquery).
		Order("starts_at ASC").
		Find(&classes).Error

	if err != nil {
		return nil, err
	}

	return classes, nil
}

// ----------------------------------------
// DeleteClass (soft delete)
// ----------------------------------------

func (s *ScheduleStore) DeleteClass(ctx context.Context, classID uuid.UUID) error {

	result := s.db.WithContext(ctx).
		Where("class_id = ?", classID).
		Delete(&scheduledomain.Class{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
