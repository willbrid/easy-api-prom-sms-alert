package httphandler

import (
	"github.com/willbrid/easy-api-prom-alert-sms/internal/usecase"
	"github.com/willbrid/easy-api-prom-alert-sms/pkg/logger"
)

type HTTPHandler struct {
	Usecases *usecase.Usecases
	iLogger  logger.ILogger
}

func NewHTTPHandler(usecases *usecase.Usecases, iLogger logger.ILogger) *HTTPHandler {
	return &HTTPHandler{usecases, iLogger}
}
