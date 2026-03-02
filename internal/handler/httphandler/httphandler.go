package httphandler

import (
	"github.com/rs/zerolog"

	"github.com/willbrid/easy-api-prom-alert-sms/internal/usecase"
)

type HTTPHandler struct {
	Usecases *usecase.Usecases
	logger   zerolog.Logger
}

func NewHTTPHandler(usecases *usecase.Usecases, logger zerolog.Logger) *HTTPHandler {
	return &HTTPHandler{usecases, logger}
}
