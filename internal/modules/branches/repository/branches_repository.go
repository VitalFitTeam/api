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

func (s *BranchesStore) GetBranches(ctx context.Context, fq pagination.PaginatedFeedQuery) ([]*branchdomain.Branch, error) {
	var branches []*branchdomain.Branch

	result := s.db.
		Preload("State.Country").
		Preload("Manager").
		Limit(fq.Limit).
		Offset(fq.Offset).Where("name ILIKE ? AND status = ?", "%"+fq.Search+"%", fq.Status).
		Find(&branches)

	if result.Error != nil {
		return nil, result.Error
	}

	return branches, nil

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
