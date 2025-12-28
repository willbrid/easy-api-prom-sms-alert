package microservice

import (
	"github.com/willbrid/easy-api-prom-alert-sms/config"
	"github.com/willbrid/easy-api-prom-alert-sms/internal/microservice/alert"
	"github.com/willbrid/easy-api-prom-alert-sms/pkg/logger"
)

type IAlertMicroservice interface {
	Consume(url string, body string) error
	GetUrlAndBody(dest string, message string) (string, string, error)
}

type Microservices struct {
	IAlertMicroservice IAlertMicroservice
}

func NewMicroservices(provider *config.Provider, iLogger logger.ILogger) *Microservices {
	return &Microservices{
		IAlertMicroservice: alert.NewAlertMicroservice(provider, iLogger),
	}
}
