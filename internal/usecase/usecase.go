package usecase

import (
	"github.com/rs/zerolog"

	"github.com/willbrid/easy-api-prom-alert-sms/internal/domain"
	"github.com/willbrid/easy-api-prom-alert-sms/internal/microservice"
	"github.com/willbrid/easy-api-prom-alert-sms/internal/usecase/alert"
)

type IAlertUsecase interface {
	Send(domain.Alert) error
}

type Usecases struct {
	IAlertUsecase IAlertUsecase
}

type Deps struct {
	Microservices *microservice.Microservices
	AlertConfig   *domain.AlertConfig
	Logger        zerolog.Logger
}

func NewUsecases(deps *Deps) *Usecases {
	return &Usecases{
		IAlertUsecase: alert.NewAlertUseCase(deps.Microservices.IAlertMicroservice, deps.AlertConfig, deps.Logger),
	}
}
