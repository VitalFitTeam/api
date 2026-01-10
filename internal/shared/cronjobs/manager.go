package cronjobs

import (
	"github.com/robfig/cron/v3"
	appservices "github.com/vitalfit/api/internal/app/services"
	noti "github.com/vitalfit/api/internal/shared/notifications"
	"github.com/vitalfit/api/internal/store"
	"github.com/vitalfit/api/internal/store/cache"
	"go.uber.org/zap"
)

type Manager struct {
	cron        *cron.Cron
	store       store.Storage
	appservices appservices.Services
	cache       cache.Storage
	push        noti.PushService
	logger      *zap.SugaredLogger
}

func NewManager(store store.Storage, appservices appservices.Services, logger *zap.SugaredLogger, cache cache.Storage, noti noti.PushService) *Manager {
	return &Manager{
		cron:        cron.New(),
		store:       store,
		appservices: appservices,
		cache:       cache,
		push:        noti,
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
	_, err = m.cron.AddFunc("@every 2h2s", m.GetRatesCronjob)
	if err != nil {
		m.logger.Errorw("Error getting rates cronjob", "error", err)
	}

	_, err = m.cron.AddFunc("@daily", m.UpdateExpiredMembershipsCronjob)
	if err != nil {
		m.logger.Errorw("Error registering update expired memberships cronjob", "error", err)
	}

	_, err = m.cron.AddFunc("@daily", m.TestPushNoti)
	if err != nil {
		m.logger.Errorw("Error testing push noti cronjobs", "error", err)
	}

}
