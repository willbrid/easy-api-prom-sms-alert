package alert

import (
	"github.com/willbrid/easy-api-prom-alert-sms/internal/domain"
	"github.com/willbrid/easy-api-prom-alert-sms/internal/microservice"
	"github.com/willbrid/easy-api-prom-alert-sms/internal/pkg/alerthelper"
	"github.com/willbrid/easy-api-prom-alert-sms/internal/pkg/recipienthelper"
	"github.com/willbrid/easy-api-prom-alert-sms/pkg/logger"
)

type AlertUseCase struct {
	iMsc        microservice.IAlertMicroservice
	alertConfig *domain.AlertConfig
	iLogger     logger.ILogger
}

func NewAlertUseCase(iMsc microservice.IAlertMicroservice, alertConfig *domain.AlertConfig, iLogger logger.ILogger) *AlertUseCase {
	return &AlertUseCase{iMsc, alertConfig, iLogger}
}

func (auc *AlertUseCase) Send(alertData domain.Alert) error {
	recipients := auc.alertConfig.Recipients
	defaultRecipientName := auc.alertConfig.DefaultRecipientName

	for _, alert := range alertData.Data.Alerts {
		alertMsg := alerthelper.GetMsgFromAlert(alert)
		recipientName := alerthelper.GetRecipientFromAlert(alert, defaultRecipientName)
		members := recipienthelper.GetRecipientMembers(recipients, recipientName)

		for _, member := range members {
			url, body, err := auc.iMsc.GetUrlAndBody(member, alertMsg)

			if err != nil {
				return err
			}

			if auc.alertConfig.Simulation {
				auc.iLogger.Info("send request with url %s and body %s", url, body)
			} else {
				if err := auc.iMsc.Consume(url, body); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
