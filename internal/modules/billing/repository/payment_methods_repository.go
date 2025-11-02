package billingrepository

import (
	"context"

	"github.com/google/uuid"
	billingdomain "github.com/vitalfit/api/internal/modules/billing/domain"
	shared_errors "github.com/vitalfit/api/internal/shared/errors"
	"github.com/vitalfit/api/pkg/db"
	"gorm.io/gorm"
)

type PaymentMethodsStore struct {
	db *gorm.DB
}

func NewPaymentMethodsStore(db *gorm.DB) *PaymentMethodsStore {
	return &PaymentMethodsStore{db: db}
}

func (s *PaymentMethodsStore) GetPaymentMethods(ctx context.Context) ([]*billingdomain.PaymentMethods, error) {
	var paymentMethods []*billingdomain.PaymentMethods
	err := s.db.WithContext(ctx).Find(&paymentMethods).Error
	if err != nil {
		return nil, err
	}
	return paymentMethods, nil
}

func (s *PaymentMethodsStore) GetPaymentMethodByID(ctx context.Context, methodID uuid.UUID) (*billingdomain.PaymentMethods, error) {
	ctx, cancel := context.WithTimeout(ctx, db.QueryTimeoutDuration)
	defer cancel()
	var paymentMethod billingdomain.PaymentMethods
	err := s.db.WithContext(ctx).Where("method_id = ?", methodID).First(&paymentMethod).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, shared_errors.ErrNotFound
		}
		return nil, err
	}
	return &paymentMethod, nil
}

// func (s *BranchesStore) AddPaymentMethodsToBranch(ctx context.Context, branchID uuid.UUID, paymentLinks []branchdomain.PaymentMethodsBranch) error {
// 	if len(paymentLinks) == 0 {
// 		return nil
// 	}

// 	for i := range paymentLinks {
// 		paymentLinks[i].BranchID = branchID
// 	}

// 	return db.WithTX(s.db, func(tx *gorm.DB) error {
// 		return tx.WithContext(ctx).Clauses(clause.OnConflict{
// 			Columns:   []clause.Column{{Name: "branch_id"}, {Name: "method_id"}},
// 			DoNothing: true,
// 		}).Create(&paymentLinks).Error
// 	})
// }
