package app

import (
	"github.com/vitalfit/api/config"
	apphandlers "github.com/vitalfit/api/internal/app/handlers"
	appservices "github.com/vitalfit/api/internal/app/services"
	"github.com/vitalfit/api/internal/store"
	"github.com/vitalfit/api/pkg/mailer"
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
	store := store.NewStorage(db, *cfg, mailer)
	services := appservices.NewServices(store, logger)
	handlers := apphandlers.NewAppHandlers(services)
	defer logger.Sync()
	return &application{
		Config:   cfg,
		Logger:   logger,
		store:    store,
		services: services,
		handlers: handlers,
	}
}
