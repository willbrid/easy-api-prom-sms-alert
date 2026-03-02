package alert

import (
	"github.com/rs/zerolog"
	"github.com/willbrid/easy-api-prom-alert-sms/internal/domain"
	"github.com/willbrid/easy-api-prom-alert-sms/internal/microservice"
	"github.com/willbrid/easy-api-prom-alert-sms/internal/pkg/alerthelper"
	"github.com/willbrid/easy-api-prom-alert-sms/internal/pkg/recipienthelper"
)

type AlertUseCase struct {
	iMsc        microservice.IAlertMicroservice
	alertConfig *domain.AlertConfig
	logger      zerolog.Logger
}

func NewAlertUseCase(iMsc microservice.IAlertMicroservice, alertConfig *domain.AlertConfig, logger zerolog.Logger) *AlertUseCase {
	return &AlertUseCase{iMsc, alertConfig, logger}
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
				auc.logger.Error().Err(err).Msg("failed to get url and body request")
				return err
			}

			if auc.alertConfig.Simulation {
				auc.logger.Info().Str("simulation", "true").Msgf("send request with url %s and body %s", url, body)
			} else {
				if err := auc.iMsc.Consume(url, body); err != nil {
					auc.logger.Error().Err(err).Str("url", url).Str("body", body).Msg("failed to send request")
					return err
				}
			}
		}
	}

	return nil
}
