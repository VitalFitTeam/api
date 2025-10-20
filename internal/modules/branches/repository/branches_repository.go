package branchrepository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	branchdomain "github.com/vitalfit/api/internal/modules/branches/domain"
	shared_errors "github.com/vitalfit/api/internal/shared/errors"
	"github.com/vitalfit/api/pkg/db"
	"gorm.io/gorm"
)

type BranchesStore struct {
	db *gorm.DB
}

func NewBranchesStore(db *gorm.DB) *BranchesStore {
	return &BranchesStore{db: db}
}

func (s *BranchesStore) create(ctx context.Context, tx *gorm.DB, branch *branchdomain.Branch) error {
	err := tx.WithContext(ctx).Create(&branch).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if pgErr.ConstraintName == "users_tax_id_key" {
				return shared_errors.ErrConflict
			}
		}
		return err
	}
	return nil

}

func (s *BranchesStore) CreateBranch(ctx context.Context, branch *branchdomain.Branch) error {
	//transaction
	return db.WithTX(s.db, func(tx *gorm.DB) error {
		if err := s.create(ctx, tx, branch); err != nil {
			return err //rollback
		}
		return nil
	})
}
