package combosrepository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	combosdomain "github.com/vitalfit/api/internal/modules/combos/domain"
	shared_errors "github.com/vitalfit/api/internal/shared/errors"
	"github.com/vitalfit/api/pkg/db"
	"github.com/vitalfit/api/pkg/pagination"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CombosStore struct {
	db *gorm.DB
}

func NewCombosStore(db *gorm.DB) *CombosStore {
	return &CombosStore{db: db}
}

func (s *CombosStore) CreatePackage(ctx context.Context, pkg *combosdomain.Package) error {
	return db.WithTX(s.db, func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Create(pkg).Error
	})
}

func (s *CombosStore) GetPackages(ctx context.Context, fq pagination.PaginatedFeedQuery) ([]*combosdomain.Package, error) {
	var packages []*combosdomain.Package
	query := s.db.WithContext(ctx)
	if fq.Search != "" {
		searchQuery := "%" + fq.Search + "%"
		query = query.Where("name ILIKE ?", searchQuery)
	}

	err := query.Limit(fq.Limit).Offset(fq.Page*fq.Limit - fq.Limit).Order("created_at " + fq.Sort).Find(&packages).Error
	return packages, err
}

func (s *CombosStore) GetPackagesTotal(ctx context.Context, fq pagination.PaginatedFeedQuery) (int64, error) {
	var count int64
	query := s.db.WithContext(ctx).Model(&combosdomain.Package{})
	if fq.Search != "" {
		searchQuery := "%" + fq.Search + "%"
		query = query.Where("name ILIKE ?", searchQuery)
	}
	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *CombosStore) GetPackageByID(ctx context.Context, packageID uuid.UUID) (*combosdomain.Package, error) {
	var pkg combosdomain.Package
	err := s.db.WithContext(ctx).Preload("PackageItems.Service").Preload("PackageItems").First(&pkg, "package_id = ?", packageID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared_errors.ErrNotFound
		}
		return nil, err
	}
	return &pkg, nil
}

func (s *CombosStore) UpdatePackage(ctx context.Context, pkg *combosdomain.Package) error {
	return db.WithTX(s.db, func(tx *gorm.DB) error {
		var existingPackage combosdomain.Package
		if err := tx.WithContext(ctx).First(&existingPackage, "package_id = ?", pkg.PackageID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return shared_errors.ErrNotFound
			}
			return err
		}

		if err := tx.WithContext(ctx).Model(&existingPackage).Updates(pkg).Error; err != nil {
			return err
		}

		if err := tx.WithContext(ctx).Where("package_id = ?", pkg.PackageID).Delete(&combosdomain.PackageItem{}).Error; err != nil {
			return err
		}

		if len(pkg.PackageItems) > 0 {
			for i := range pkg.PackageItems {
				pkg.PackageItems[i].PackageID = pkg.PackageID
			}
			if err := tx.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&pkg.PackageItems).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *CombosStore) DeletePackage(ctx context.Context, packageID uuid.UUID) error {
	result := s.db.WithContext(ctx).Delete(&combosdomain.Package{}, packageID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return shared_errors.ErrNotFound
	}
	return nil
}
