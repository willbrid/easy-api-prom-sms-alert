package usecase

import (
	"github.com/willbrid/easy-api-prom-alert-sms/internal/domain"
	"github.com/willbrid/easy-api-prom-alert-sms/internal/microservice"
	"github.com/willbrid/easy-api-prom-alert-sms/internal/usecase/alert"
	"github.com/willbrid/easy-api-prom-alert-sms/pkg/logger"
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
	ILogger       logger.ILogger
}

func NewUsecases(deps *Deps) *Usecases {
	return &Usecases{
		IAlertUsecase: alert.NewAlertUseCase(deps.Microservices.IAlertMicroservice, deps.AlertConfig, deps.ILogger),
	}
}
