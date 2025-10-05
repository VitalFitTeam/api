package app

import (
	"github.com/vitalfit/api/config"
	apphandlers "github.com/vitalfit/api/internal/app/handlers"
	appservices "github.com/vitalfit/api/internal/app/services"
	"github.com/vitalfit/api/internal/store"
	"go.uber.org/zap"
)

func BuildApplication(cfg *config.Config, store store.Storage) *application {
	//logger
	logger := zap.Must(zap.NewProduction()).Sugar()
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
