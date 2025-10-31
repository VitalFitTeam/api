package instructorrepository

import (
	"context"

	"github.com/google/uuid"
	instructordomain "github.com/vitalfit/api/internal/modules/instructor/domain"
	"github.com/vitalfit/api/pkg/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *InstructorStore) AssignInstructorsToBranch(ctx context.Context, branchID uuid.UUID, instructorID []uuid.UUID) error {
	var linksToCreate []instructordomain.BranchInstructor
	for _, pid := range instructorID {
		linksToCreate = append(linksToCreate, instructordomain.BranchInstructor{
			BranchID:     branchID,
			InstructorID: pid,
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
func (s *InstructorStore) ListBranchInstructors(ctx context.Context, branchID uuid.UUID) ([]*instructordomain.Instructor, error) {
	var instructors []*instructordomain.Instructor

	err := s.db.WithContext(ctx).
		Preload("User").
		Joins("JOIN branch_instructors ON branch_instructors.instructor_id = instructors.instructor_id").
		Where("branch_instructors.branch_id = ?", branchID).
		Find(&instructors).Error

	if err != nil {
		return nil, err
	}

	return instructors, nil
}

func (s *InstructorStore) RemoveInstructorFromBranch(ctx context.Context, branchID uuid.UUID, instructorID uuid.UUID) error {

	association := instructordomain.BranchInstructor{
		BranchID:     branchID,
		InstructorID: instructorID,
	}

	result := s.db.WithContext(ctx).Delete(&association)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
