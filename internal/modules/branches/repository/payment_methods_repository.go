package branchrepository

import (
	"context"

	branchdomain "github.com/vitalfit/api/internal/modules/branches/domain"
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

func (s *PaymentMethodsStore) GetPaymentMethods(ctx context.Context) ([]*branchdomain.PaymentMethods, error) {
	var paymentMethods []*branchdomain.PaymentMethods
	err := s.db.WithContext(ctx).Find(&paymentMethods).Error
	if err != nil {
		return nil, err
	}
	return paymentMethods, nil
}

func (s *PaymentMethodsStore) GetPaymentMethodByName(ctx context.Context, name string) (*branchdomain.PaymentMethods, error) {
	ctx, cancel := context.WithTimeout(ctx, db.QueryTimeoutDuration)
	defer cancel()
	var paymentMethod branchdomain.PaymentMethods
	err := s.db.WithContext(ctx).Where("name = ?", name).First(&paymentMethod).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, shared_errors.ErrNotFound
		}
		return nil, err
	}
	return &paymentMethod, nil
}
