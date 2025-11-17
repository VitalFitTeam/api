package billingservice

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/vitalfit/api/config"
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
	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	httpClient := &http.Client{
		Transport: customTransport,
		Timeout:   10 * time.Second,
	}

	return &BillingService{
		store: store,
		cache: cache,
		cfg:   cfg,
		http:  httpClient,
	}
}
