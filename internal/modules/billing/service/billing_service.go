package billingservice

import (
	"context"
	"net/http"
	"time"

	"github.com/vitalfit/api/config"
	billingdomain "github.com/vitalfit/api/internal/modules/billing/domain"
	"github.com/vitalfit/api/internal/store"
	"github.com/vitalfit/api/internal/store/cache"
)

type BillingService struct {
	store store.Storage
	cache cache.Storage
	cfg   config.Config
	http  *http.Client
}

func NewBillingService(store store.Storage, cache cache.Storage, cfg config.Config) *BillingService {
	return &BillingService{
		store: store,
		cache: cache,
		cfg:   cfg,
		http:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *BillingService) Checkout(ctx context.Context, invoice *billingdomain.Invoice) error {
	return nil
}

func (s *BillingService) CreatePayment(ctx context.Context, payment *billingdomain.Payment) error {
	return nil
}
