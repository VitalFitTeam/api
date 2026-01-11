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

func (m *Manager) UpdateExpiredMembershipsCronjob() {
	m.logger.Info("running update expired memberships cronjob")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	if err := m.appservices.MembershipServices.UpdateExpiredMemberships(ctx); err != nil {
		m.logger.Errorw("failed to update expired memberships", "error", err)
		return
	}
}

func (m *Manager) TestPushNoti() {
	m.logger.Info("running test push noti cronjob")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	token := "fdpe9i-OYersBulmuLxddC:APA91bE-udl2lAwKLein6h6nzjvDJbgzhRv_vzMLoJc8quPCAHU3tT95hP_1gP87aFBWLM6lzx_RNfrCqd5hjm6VjsLZGd_5MY4UdNRUxUx4wiZ_Ll2zhXU"
	m.logger.Infow("token", "token", token)
	data := map[string]string{
		"test_data": "test data",
	}
	m.push.SendPush(ctx, token, "test", "test", data)
}

func (m *Manager) TestChurn() {
	m.logger.Info("running test churn cronjob")
	token := "fdpe9i-OYersBulmuLxddC:APA91bE-udl2lAwKLein6h6nzjvDJbgzhRv_vzMLoJc8quPCAHU3tT95hP_1gP87aFBWLM6lzx_RNfrCqd5hjm6VjsLZGd_5MY4UdNRUxUx4wiZ_Ll2zhXU"
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	xd, err := m.appservices.ReportServices.DetectAndFlagChurnRisk(ctx)
	data, err := m.push.ConvertSliceToDataMap(xd, "churn_risk_data")
	if err != nil {
		m.logger.Errorw("failed to convert struct to data map", "error", err)
		return
	}
	m.logger.Infow("xd", "xd", data)
	m.push.SendPush(ctx, token, "test", "test", data)

}

func (m *Manager) testBroadcast() {
	m.logger.Info("running test broadcast cronjob")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	if err := m.appservices.NotificationServices.SendBroadcast(ctx, "que onda papu", "qlq"); err != nil {
		m.logger.Errorw("failed to send broadcast", "error", err)
		return

	}
}
