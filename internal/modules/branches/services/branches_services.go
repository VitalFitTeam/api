package branchservices

import (
	"context"

	"github.com/google/uuid"
	"github.com/vitalfit/api/config"
	branchdomain "github.com/vitalfit/api/internal/modules/branches/domain"
	"github.com/vitalfit/api/internal/store"
)

type BranchServices struct {
	store  store.Storage
	config config.Config
}

func NewBranchServices(store store.Storage, cfg config.Config) *BranchServices {
	return &BranchServices{
		store:  store,
		config: cfg,
	}

}

func (s *BranchServices) CreateBranch(ctx context.Context, branch *branchdomain.Branch) error {
	branch, err := s.store.Branches.CreateBranch(ctx, branch)
	if err != nil {
		return err
	}
	if err := s.store.Branches.AddPaymentMethodsToBranch(ctx, branch.BranchID, branch.PaymentMethodsLinks); err != nil {
		return err
	}
	return nil
}

func (s *BranchServices) GetPaymentMethodByID(ctx context.Context, methodID uuid.UUID) (*branchdomain.PaymentMethods, error) {
	paymentMethod, err := s.store.PaymentMethods.GetPaymentMethodByID(ctx, methodID)
	if err != nil {
		return nil, err
	}
	return paymentMethod, nil
}

func (s *BranchServices) GetPaymentMethods(ctx context.Context) ([]*branchdomain.PaymentMethods, error) {
	paymentMethods, err := s.store.PaymentMethods.GetPaymentMethods(ctx)
	if err != nil {
		return nil, err
	}
	return paymentMethods, nil
}
