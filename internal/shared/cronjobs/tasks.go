package cronjobs

import (
	"context"
	"time"
)

func (m *Manager) GetRatesCronjob() {
	m.logger.Info("running get currency rates cronjob")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rates, err := m.appservices.BillingServices.GetLatestRates(ctx)
	if err != nil {
		m.logger.Errorw("failed to get rates", "error", err)
		return
	}
	if len(rates) == 0 {
		m.logger.Warn("rates map is empty")
		return
	}
}
