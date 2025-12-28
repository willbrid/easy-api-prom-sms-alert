package handler

import (
	"github.com/willbrid/easy-api-prom-alert-sms/config"
	"github.com/willbrid/easy-api-prom-alert-sms/internal/handler/httphandler"
	"github.com/willbrid/easy-api-prom-alert-sms/internal/handler/middleware"
	"github.com/willbrid/easy-api-prom-alert-sms/internal/usecase"
	"github.com/willbrid/easy-api-prom-alert-sms/pkg/logger"

	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	Usecases *usecase.Usecases
	Router   *mux.Router
	iLogger  logger.ILogger
}

func NewHandler(usecases *usecase.Usecases, router *mux.Router, iLogger logger.ILogger) *Handler {
	return &Handler{usecases, router, iLogger}
}

func (h *Handler) InitRouter(cfg *config.Config) {
	h.Router.Use(func(h http.Handler) http.Handler {
		return middleware.AuthMiddleware(h, cfg)
	})

	httphandler := httphandler.NewHTTPHandler(h.Usecases, h.iLogger)

	h.Router.HandleFunc("/healthz", httphandler.HandleHealthCheck).Methods("GET")
	h.Router.HandleFunc("/api-alert", httphandler.HandleAlert).Methods("POST")
}
