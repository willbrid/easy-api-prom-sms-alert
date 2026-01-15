package microservice

import (
	"github.com/willbrid/easy-api-prom-alert-sms/config"
	"github.com/willbrid/easy-api-prom-alert-sms/internal/microservice/alert"
	"github.com/willbrid/easy-api-prom-alert-sms/pkg/httpclient"
	"github.com/willbrid/easy-api-prom-alert-sms/pkg/httpparam"
	"github.com/willbrid/easy-api-prom-alert-sms/pkg/logger"
)

type IAlertMicroservice interface {
	Consume(url string, body string) error
	GetUrlAndBody(dest string, message string) (string, string, error)
}

type Deps struct {
	Provider      *config.Provider
	IHTTPClient   httpclient.IHTTPClient
	IParamFactory httpparam.IParamFactory
	ILogger       logger.ILogger
}

type Microservices struct {
	IAlertMicroservice IAlertMicroservice
}

func NewMicroservices(deps Deps) *Microservices {
	return &Microservices{
		IAlertMicroservice: alert.NewAlertMicroservice(deps.Provider, deps.IHTTPClient, deps.IParamFactory, deps.ILogger),
	}
}
