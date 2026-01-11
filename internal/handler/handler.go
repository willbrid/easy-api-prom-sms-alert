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
	Usecases        *usecase.Usecases
	Router          *mux.Router
	iAuthMiddleware middleware.IAuthMiddleware
	iLogger         logger.ILogger
}

func NewHandler(usecases *usecase.Usecases, router *mux.Router, iAuthMiddleware middleware.IAuthMiddleware, iLogger logger.ILogger) *Handler {
	return &Handler{usecases, router, iAuthMiddleware, iLogger}
}

func (h *Handler) InitRouter(cfg *config.Config) {
	h.Router.Use(func(httpHandler http.Handler) http.Handler {
		return h.iAuthMiddleware.Authenticate(httpHandler, cfg)
	})

	httphandler := httphandler.NewHTTPHandler(h.Usecases, h.iLogger)

	h.Router.HandleFunc("/healthz", httphandler.HandleHealthCheck).Methods("GET")
	h.Router.HandleFunc("/api-alert", httphandler.HandleAlert).Methods("POST")
}
