package app

import (
	"github.com/vitalfit/api/config"
	apphandlers "github.com/vitalfit/api/internal/app/handlers"
	appservices "github.com/vitalfit/api/internal/app/services"
	authservices "github.com/vitalfit/api/internal/auth/services"
	"github.com/vitalfit/api/internal/store"
	"github.com/vitalfit/api/pkg/mailer"
	rate_mw "github.com/vitalfit/api/pkg/ratelimiter"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func BuildApplication(cfg *config.Config, db *gorm.DB) *application {
	//logger
	//initialize store

	logger := zap.Must(zap.NewProduction()).Sugar()
	mailer, err := mailer.NewResendClient(cfg.Mail.Resend.ApiKey, cfg.Mail.FromEmail)
	if err != nil {
		logger.Errorw("error creating mailer", "error", err.Error())
	}
	auth := authservices.NewJWTAuthenticator(cfg.Auth.Token.Secret, cfg.Auth.Token.Iss, cfg.Auth.Token.Iss)
	rateLimiter := rate_mw.NewFixedWindowLimiter(cfg.RateLimiter.RequestsPerTimeFrame, cfg.RateLimiter.TimeFrame)
	store := store.NewStorage(db, *cfg, mailer, auth)
	services := appservices.NewServices(store, logger)
	handlers := apphandlers.NewAppHandlers(services)
	defer logger.Sync()
	return &application{
		Config:      cfg,
		Logger:      logger,
		Store:       store,
		Services:    services,
		Handlers:    handlers,
		ratelimiter: rateLimiter,
	}
}
