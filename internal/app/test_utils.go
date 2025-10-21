package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vitalfit/api/config"
	apphandlers "github.com/vitalfit/api/internal/app/handlers"
	appservices "github.com/vitalfit/api/internal/app/services"
	authmocks "github.com/vitalfit/api/internal/modules/auth/mocks"
	"github.com/vitalfit/api/internal/store"
	mailermocks "github.com/vitalfit/api/pkg/mailer/mocks"
	"github.com/vitalfit/api/pkg/ratelimiter"
	"go.uber.org/zap"
)

func NewTestApplication(t *testing.T, cfg *config.Config) *application {
	t.Helper()

	logger := zap.NewNop().Sugar()
	// Uncomment to enable logs
	//logger := zap.Must(zap.NewProduction()).Sugar()
	testAuth := &authmocks.TestAuthenticator{}
	mailer := &mailermocks.MockMailer{}
	mockStore := store.NewMockStore()

	// Rate limiter
	rateLimiter := ratelimiter.NewFixedWindowLimiter(
		cfg.RateLimiter.RequestsPerTimeFrame,
		cfg.RateLimiter.TimeFrame,
	)
	mockServices := appservices.NewServices(mockStore, logger, *cfg, testAuth, mailer)
	mockHandlers := apphandlers.NewAppHandlers(mockServices)

	return &application{
		Logger:      logger,
		ratelimiter: rateLimiter,
		Store:       mockStore,
		Services:    mockServices,
		Handlers:    mockHandlers,
		Config:      cfg,
	}
}

func ExecuteRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	return rr
}

func CheckResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d", expected, actual)
	}
}
