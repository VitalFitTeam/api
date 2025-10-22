package branchrepository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	branchdomain "github.com/vitalfit/api/internal/modules/branches/domain"
	shared_errors "github.com/vitalfit/api/internal/shared/errors"
	"github.com/vitalfit/api/pkg/db"
	"github.com/vitalfit/api/pkg/pagination"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BranchesStore struct {
	db *gorm.DB
}

func NewBranchesStore(db *gorm.DB) *BranchesStore {
	return &BranchesStore{db: db}
}

func (s *BranchesStore) create(ctx context.Context, tx *gorm.DB, branch *branchdomain.Branch) error {
	err := tx.WithContext(ctx).Omit("PaymentMethodsLinks").Create(&branch).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if pgErr.ConstraintName == "branch_tax_id_key" {
				return shared_errors.ErrConflict
			}
		}
		return err
	}
	return nil

}
func (s *BranchesStore) CreateBranch(ctx context.Context, branch *branchdomain.Branch) (*branchdomain.Branch, error) {
	// transaction
	err := db.WithTX(s.db, func(tx *gorm.DB) error {
		err := s.create(ctx, tx, branch)
		if err != nil {
			return err // rollback
		}
		return nil // commit
	})

	if err != nil {
		return nil, err
	}

	return branch, nil
}

func (s *BranchesStore) GetBranches(ctx context.Context, fq pagination.PaginatedFeedQuery) (*branchdomain.BranchQueryResults, error) {
	baseQuery := s.db.WithContext(ctx).
		Model(&branchdomain.Branch{}).
		Joins("JOIN states ON states.state_id = branch.state_id").
		Joins("JOIN countries ON countries.country_id = states.country_id")

	if fq.Search != "" {
		baseQuery = baseQuery.Where("branch.name ILIKE ?", "%"+fq.Search+"%")
	}

	if fq.Location != "" {
		locationPattern := "%" + fq.Location + "%"
		baseQuery = baseQuery.Where("states.name ILIKE ? OR countries.name ILIKE ?", locationPattern, locationPattern)
	}

	if fq.TaxID != "" {
		taxIDPattern := "%" + fq.TaxID + "%"
		baseQuery = baseQuery.Where("branch.tax_id ILIKE ?", taxIDPattern)
	}

	var counts []branchdomain.StatusCount
	if err := baseQuery.Session(&gorm.Session{NewDB: true}).
		Model(&branchdomain.Branch{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Find(&counts).Error; err != nil {
		return nil, err
	}

	var activeCount, inactiveCount, manteinanceCount int64
	for _, c := range counts {
		switch c.Status {
		case "Active":
			activeCount = c.Count
		case "Inactive":
			inactiveCount = c.Count
		case "Manteinance":
			manteinanceCount = c.Count
		}
	}

	var branches []*branchdomain.Branch

	query := baseQuery.
		Preload("State.Country").
		Preload("Manager").
		Limit(fq.Limit).
		Offset(fq.Offset).
		Order("created_at " + fq.Sort)

	if fq.Status != "" {
		query = query.Where("status = ?", fq.Status)
	}

	if err := query.Find(&branches).Error; err != nil {
		return nil, err
	}

	response := &branchdomain.BranchQueryResults{
		Branches:         branches,
		ActiveCount:      activeCount,
		InactiveCount:    inactiveCount,
		ManteinanceCount: manteinanceCount,
	}

	return response, nil

}

func (s *BranchesStore) AddPaymentMethodsToBranch(ctx context.Context, branchID uuid.UUID, paymentLinks []branchdomain.PaymentMethodsBranch) error {
	if len(paymentLinks) == 0 {
		return nil
	}

	for i := range paymentLinks {
		paymentLinks[i].BranchID = branchID
	}

	return db.WithTX(s.db, func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "branch_id"}, {Name: "method_id"}},
			DoNothing: true,
		}).Create(&paymentLinks).Error
	})
}

// softdelete
func (s *BranchesStore) Delete(ctx context.Context, branchID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, db.QueryTimeoutDuration)
	defer cancel()
	// transaction
	return db.WithTX(s.db, func(tx *gorm.DB) error {
		result := tx.WithContext(ctx).Delete(&branchdomain.Branch{}, branchID)
		if result.Error != nil {
			return result.Error //rollback
		}

		if result.RowsAffected == 0 {
			return shared_errors.ErrNotFound //rollback
		}

		return nil // commit
	})
}

func (s *BranchesStore) GetByID(ctx context.Context, branchID uuid.UUID) (*branchdomain.Branch, error) {
	var branch branchdomain.Branch
	err := s.db.WithContext(ctx).
		Joins("State").Joins("State.Country").Joins("Manager").
		Preload("OperatingHours").
		Preload("PaymentMethodsLinks.Method").
		First(&branch, "branch.branch_id = ?", branchID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared_errors.ErrNotFound
		}
		return nil, err
	}
	return &branch, nil
}

func (s *BranchesStore) Update(ctx context.Context, branch *branchdomain.Branch) error {
	ctx, cancel := context.WithTimeout(ctx, db.QueryTimeoutDuration)
	defer cancel()
	// transaction
	return db.WithTX(s.db, func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Model(&branchdomain.Branch{BranchID: branch.BranchID}).Omit("OperatingHours", "PaymentMethodsLinks").Updates(branch).Error; err != nil {
			return err // rollback
		}
		if len(branch.OperatingHours) > 0 {
			err := tx.WithContext(ctx).Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "branch_id"}, {Name: "day_of_week"}},
				DoUpdates: clause.AssignmentColumns([]string{"open_time", "close_time", "is_closed"}),
			}).Create(&branch.OperatingHours).Error

			if err != nil {
				return err // rollback
			}
		}
		if err := s.AddPaymentMethodsToBranch(ctx, branch.BranchID, branch.PaymentMethodsLinks); err != nil {
			return err // rollback
		}

		return nil
	})
}
