package staffrepository

import (
	"context"

	"github.com/google/uuid"
	staffdomain "github.com/vitalfit/api/internal/modules/staff/domain"
	"github.com/vitalfit/api/pkg/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StaffStore struct {
	db *gorm.DB
}

func NewStaffStore(db *gorm.DB) *StaffStore {
	return &StaffStore{db: db}
}

func (s *StaffStore) AssignStaffToBranch(ctx context.Context, branchID uuid.UUID, staffID []uuid.UUID) error {
	var linksToCreate []staffdomain.BranchStaff
	for _, pid := range staffID {
		linksToCreate = append(linksToCreate, staffdomain.BranchStaff{
			BranchID: branchID,
			UserID:   pid,
		})
	}

	if len(linksToCreate) == 0 {
		return nil
	}

	err := db.WithTX(s.db, func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).
			Clauses(clause.OnConflict{DoNothing: true}).
			Create(&linksToCreate).Error; err != nil {
			return err
		}
		return nil
	})

	return err
}
