package cronjobs

import (
	"github.com/robfig/cron/v3"
	appservices "github.com/vitalfit/api/internal/app/services"
	billingdomain "github.com/vitalfit/api/internal/modules/billing/domain"
	"github.com/vitalfit/api/internal/store"
	"github.com/vitalfit/api/internal/store/cache"
	"go.uber.org/zap"
)

type Manager struct {
	cron           *cron.Cron
	store          store.Storage
	appservices    appservices.Services
	cache          cache.Storage
	logger         *zap.SugaredLogger
	billingService billingdomain.BillingServiceInterface
}

func NewManager(store store.Storage, appservices appservices.Services, logger *zap.SugaredLogger, cache cache.Storage) *Manager {
	return &Manager{
		cron:        cron.New(),
		store:       store,
		appservices: appservices,
		cache:       cache,
		logger:      logger,
	}
}

func (m *Manager) Start() {
	m.logger.Info("Initializing Cronjobs...")

	m.registerRoutes()

	m.cron.Start()

	m.logger.Info("Cronjobs running in background")
}

func (m *Manager) Stop() {
	m.cron.Stop()
	m.logger.Info("Cronjobs stopped")
}

func (m *Manager) registerRoutes() {
	var err error
	_, err = m.cron.AddFunc("@every 2h", m.GetRatesCronjob)
	if err != nil {
		m.logger.Errorw("Error getting rates cronjob", "error", err)
	}

}
