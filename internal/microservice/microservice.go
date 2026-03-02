package microservice

import (
	"github.com/rs/zerolog"

	"github.com/willbrid/easy-api-prom-alert-sms/config"
	"github.com/willbrid/easy-api-prom-alert-sms/internal/microservice/alert"
	"github.com/willbrid/easy-api-prom-alert-sms/pkg/httpclient"
	"github.com/willbrid/easy-api-prom-alert-sms/pkg/httpparam"
)

type IAlertMicroservice interface {
	Consume(url string, body string) error
	GetUrlAndBody(dest string, message string) (string, string, error)
}

type Deps struct {
	Provider      *config.Provider
	IHTTPClient   httpclient.IHTTPClient
	IParamFactory httpparam.IParamFactory
	Logger        zerolog.Logger
}

type Microservices struct {
	IAlertMicroservice IAlertMicroservice
}

func NewMicroservices(deps Deps) *Microservices {
	return &Microservices{
		IAlertMicroservice: alert.NewAlertMicroservice(deps.Provider, deps.IHTTPClient, deps.IParamFactory, deps.Logger),
	}
}
