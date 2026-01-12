package app

import (
	"github.com/go-redis/redis/v8"
	"github.com/vitalfit/api/config"
	apphandlers "github.com/vitalfit/api/internal/app/handlers"
	appservices "github.com/vitalfit/api/internal/app/services"
	authservices "github.com/vitalfit/api/internal/modules/auth/services"
	"github.com/vitalfit/api/internal/modules/cronjobs"
	"github.com/vitalfit/api/internal/shared/notifications"
	"github.com/vitalfit/api/internal/store/cache"

	"github.com/vitalfit/api/internal/store"
	"github.com/vitalfit/api/pkg/mailer"
	rate_mw "github.com/vitalfit/api/pkg/ratelimiter"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func BuildApplication(cfg *config.Config, db *gorm.DB, rdb *redis.Client) *application {
	//logger
	//initialize store

	logger := zap.Must(zap.NewProduction()).Sugar()
	mailer, err := mailer.NewResendClient(cfg.Mail.Resend.ApiKey, cfg.Mail.FromEmail)
	if err != nil {
		logger.Errorw("error creating mailer", "error", err.Error())
	}
	auth := authservices.NewJWTAuthenticator(cfg.Auth.Token.Secret, cfg.Auth.Token.Aud, cfg.Auth.Token.Iss)
	rateLimiter := rate_mw.NewFixedWindowLimiter(cfg.RateLimiter.RequestsPerTimeFrame, cfg.RateLimiter.TimeFrame)
	notifications, err := notifications.NewPushService()
	if err != nil {
		logger.Errorw("error creating push service", "error", err.Error())
	}
	store := store.NewStorage(db)

	cache := cache.NewRedisStorage(rdb)
	services := appservices.NewServices(store, logger, *cfg, auth, mailer, cache, *notifications)
	handlers := apphandlers.NewAppHandlers(services)

	cronjob := cronjobs.NewManager(store, services, mailer, logger, cache, *notifications, *cfg)
	defer logger.Sync()
	return &application{
		Config:      cfg,
		Logger:      logger,
		Store:       store,
		Cache:       cache,
		Services:    services,
		Handlers:    handlers,
		ratelimiter: rateLimiter,
		Cronjob:     cronjob,
	}
}
