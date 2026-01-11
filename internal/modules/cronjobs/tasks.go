package cronjobs

import (
	"context"
	"fmt"
	"time"

	notidomain "github.com/vitalfit/api/internal/modules/notifications/domain"
	"github.com/vitalfit/api/pkg/mailer"
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

func (m *Manager) TestBroadcast() {
	now := time.Now()
	m.logger.Infof("running test broadcast cronjob %v", now)

}

func (m *Manager) NotifyClassReminder() {
	m.logger.Info("running notify class reminder cronjob")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	BookingRemind, err := m.appservices.BookingServices.GetUpcomingClassReminders(ctx)
	if err != nil {
		m.logger.Errorw("failed to get upcoming class reminders", "error", err)
		return
	}
	if len(BookingRemind) == 0 {
		m.logger.Warn("BookingRemind is empty")
		return
	}
	var notifications []notidomain.Notification
	for _, booking := range BookingRemind {
		var notification notidomain.Notification
		notification.UserID = booking.UserID
		notification.Title = fmt.Sprintf("Class Reminder %s", booking.ServiceName)
		notification.Message = fmt.Sprintf("Your %s class starts in %v", booking.ServiceName, booking.TimeUntilStart)
		notification.Type = string(notidomain.ClassReminder)
		notification.IsRead = false
		metadata, err := m.push.ConvertStructToDataMap(booking)
		if err != nil {
			m.logger.Errorw("error converting metadata")
		}
		metaInterface := make(map[string]interface{}, len(metadata))
		for k, v := range metadata {
			metaInterface[k] = v
		}
		notification.Metadata = metaInterface
		session, err := m.store.Session.GetUserSessions(ctx, booking.UserID)
		if err != nil {
			m.logger.Errorw("failed to get user sessions")
		}
		for _, ses := range session {
			if ses.DeviceToken != "" {
				err := m.push.SendPush(ctx, ses.DeviceToken, fmt.Sprintf("Class Reminder %s", booking.ServiceName), fmt.Sprintf("Your %s class starts in %v", booking.ServiceName, booking.TimeUntilStart), metadata)
				if err != nil {
					m.logger.Errorw("failed to send push notification", "error", err)
				}
			}
		}
		notifications = append(notifications, notification)
	}
	if err := m.appservices.NotificationServices.CreateBatchNotifications(ctx, notifications); err != nil {
		m.logger.Errorw("failed to create batch notifications", "error", err)
		return
	}
	m.logger.Infow("notifications", "notifications", notifications)

}

func (m *Manager) MembershipExpiringNotification() {
	m.logger.Info("running membership expiring cronjob")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	memberships, err := m.appservices.MembershipServices.GetExpiringMemberships(ctx, 3)
	if err != nil {
		m.logger.Errorw("failed to get expiring memberships", "error", err)
		return
	}
	if len(memberships) == 0 {
		m.logger.Warn("memberships is empty")
		return
	}
	var notifications []notidomain.Notification
	for _, membership := range memberships {
		var notification notidomain.Notification
		notification.UserID = membership.UserID
		notification.Title = fmt.Sprintf("%s Expiring ", membership.MembershipName)
		notification.Message = fmt.Sprintf("Your memberships %s is expiring in %d days", membership.MembershipName, membership.DaysRemaining)
		notification.Type = string(notidomain.MembershipExpiring)
		notification.IsRead = false
		metadata, err := m.push.ConvertStructToDataMap(membership)
		if err != nil {
			m.logger.Errorw("error converting metadata")
		}
		metaInterface := make(map[string]interface{}, len(metadata))
		for k, v := range metadata {
			metaInterface[k] = v
		}
		notification.Metadata = metaInterface

		mem := membership
		go func() {
			renewalURL := fmt.Sprintf("%s/en/memberships", m.config.FrontURLE)
			isProdEnv := m.config.Env == "production"

			data := struct {
				UserName       string
				MembershipName string
				DaysRemaining  int
				RenewalURL     string
			}{
				UserName:       mem.UserName,
				MembershipName: mem.MembershipName,
				DaysRemaining:  mem.DaysRemaining,
				RenewalURL:     renewalURL,
			}

			if _, err := m.Mailer.Send(mailer.MembershipExpiring, mem.UserName, mem.UserEmail, data, !isProdEnv); err != nil {
				m.logger.Errorw("failed to send membership expiring email", "error", err)
			}
		}()
		notifications = append(notifications, notification)
	}
	if len(notifications) > 0 {
		if err := m.appservices.NotificationServices.CreateBatchNotifications(ctx, notifications); err != nil {
			m.logger.Errorw("failed to create batch notifications", "error", err)
		}
	}
	m.logger.Infow("memberships", "memberships", memberships)
}
