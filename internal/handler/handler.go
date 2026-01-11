package handler

import (
	"github.com/willbrid/easy-api-prom-alert-sms/config"
	"github.com/willbrid/easy-api-prom-alert-sms/internal/handler/httphandler"
	"github.com/willbrid/easy-api-prom-alert-sms/internal/handler/middleware"
	"github.com/willbrid/easy-api-prom-alert-sms/internal/usecase"
	"github.com/willbrid/easy-api-prom-alert-sms/pkg/httpserver"
	"github.com/willbrid/easy-api-prom-alert-sms/pkg/logger"

	"net/http"
)

type Handler struct {
	Usecases        *usecase.Usecases
	iServer         httpserver.IServer
	iAuthMiddleware middleware.IAuthMiddleware
	iLogger         logger.ILogger
}

func NewHandler(usecases *usecase.Usecases, iServer httpserver.IServer, iAuthMiddleware middleware.IAuthMiddleware, iLogger logger.ILogger) *Handler {
	return &Handler{usecases, iServer, iAuthMiddleware, iLogger}
}

func (h *Handler) InitRouter(cfg *config.Config) {
	router := h.iServer.GetRouter()
	router.Use(func(httpHandler http.Handler) http.Handler {
		return h.iAuthMiddleware.Authenticate(httpHandler, cfg)
	})

	httphandler := httphandler.NewHTTPHandler(h.Usecases, h.iLogger)

	router.HandleFunc("/healthz", httphandler.HandleHealthCheck).Methods("GET")
	router.HandleFunc("/api-alert", httphandler.HandleAlert).Methods("POST")
}
