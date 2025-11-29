package productsrepository

import (
	"context"

	"github.com/google/uuid"
	productsdomain "github.com/vitalfit/api/internal/modules/products/domain"
	"github.com/vitalfit/api/pkg/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *ProductsStore) AssignBranchService(ctx context.Context, branchServices []*productsdomain.ServiceBranchDetail) error {
	if len(branchServices) == 0 {
		return nil
	}
	err := db.WithTX(s.db, func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).
			Clauses(clause.OnConflict{DoNothing: true}).
			Create(&branchServices).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *ProductsStore) GetBranchService(ctx context.Context, branchID uuid.UUID) ([]*productsdomain.ServiceBranchDetail, error) {
	var branchServices []*productsdomain.ServiceBranchDetail

	err := db.WithTX(s.db, func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).
			Preload("Service").
			Preload("Branch").
			Find(&branchServices, "branch_id = ?", branchID).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return branchServices, nil
}

func (s *ProductsStore) UpdateBranchService(ctx context.Context, branchService *productsdomain.ServiceBranchDetail) error {
	result := s.db.WithContext(ctx).
		Model(&productsdomain.ServiceBranchDetail{}).
		Where("branch_id = ? AND service_id = ?", branchService.BranchID, branchService.ServiceID).
		Updates(branchService)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *ProductsStore) DeleteBranchService(ctx context.Context, branchID uuid.UUID, serviceID uuid.UUID) error {
	association := productsdomain.ServiceBranchDetail{
		BranchID:  branchID,
		ServiceID: serviceID,
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

func (s *ProductsStore) GetBranchServiceByID(ctx context.Context, branchID uuid.UUID, serviceID uuid.UUID) (*productsdomain.ServiceBranchDetail, error) {
	var branchService productsdomain.ServiceBranchDetail
	err := db.WithTX(s.db, func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).
			Preload("Service").
			Preload("Branch").
			First(&branchService, "branch_id = ? AND service_id = ?", branchID, serviceID).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &branchService, nil
}

func (s *ProductsStore) GetBranchServicesByIDs(ctx context.Context, branchID uuid.UUID, serviceIDs []uuid.UUID) (map[uuid.UUID]*productsdomain.ServiceBranchDetail, error) {
	if len(serviceIDs) == 0 {
		return make(map[uuid.UUID]*productsdomain.ServiceBranchDetail), nil
	}

	var branchServices []*productsdomain.ServiceBranchDetail
	err := s.db.WithContext(ctx).
		Where("branch_id = ?", branchID).
		Where("service_id IN ?", serviceIDs).
		Find(&branchServices).Error

	if err != nil {
		return nil, err
	}

	branchServicesMap := make(map[uuid.UUID]*productsdomain.ServiceBranchDetail, len(branchServices))
	for _, bs := range branchServices {
		branchServicesMap[bs.ServiceID] = bs
	}

	return branchServicesMap, nil
}
