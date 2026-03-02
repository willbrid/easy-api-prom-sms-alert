package app

import (
	"github.com/willbrid/easy-api-prom-alert-sms/config"
	"github.com/willbrid/easy-api-prom-alert-sms/internal/domain"
	"github.com/willbrid/easy-api-prom-alert-sms/internal/handler"
	"github.com/willbrid/easy-api-prom-alert-sms/internal/handler/middleware"
	"github.com/willbrid/easy-api-prom-alert-sms/internal/microservice"
	"github.com/willbrid/easy-api-prom-alert-sms/internal/usecase"
	"github.com/willbrid/easy-api-prom-alert-sms/pkg/httpclient"
	"github.com/willbrid/easy-api-prom-alert-sms/pkg/httpparam"
	"github.com/willbrid/easy-api-prom-alert-sms/pkg/httpserver"

	"github.com/rs/zerolog"

	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfgfile *config.Config, cfgflag *config.ConfigFlag, logger zerolog.Logger) {
	httpClient := httpclient.NewHTTPClient()
	paramFactory := httpparam.NewParamFactory()
	microservices := microservice.NewMicroservices(microservice.Deps{
		Provider:      &cfgfile.EasyAPIPromAlertSMS.Provider,
		IHTTPClient:   httpClient,
		IParamFactory: paramFactory,
		Logger:        logger,
	})
	usecases := usecase.NewUsecases(&usecase.Deps{
		Microservices: microservices,
		AlertConfig: &domain.AlertConfig{
			Recipients:           cfgfile.EasyAPIPromAlertSMS.Recipients,
			DefaultRecipientName: cfgfile.EasyAPIPromAlertSMS.Parameters.To.ParamValue,
			Simulation:           cfgfile.EasyAPIPromAlertSMS.Simulation,
		},
		Logger: logger,
	})

	httpServer := httpserver.NewServer(
		fmt.Sprint(":"+fmt.Sprint(cfgflag.ListenPort)),
		cfgflag.EnableHttps,
		cfgflag.CertFile,
		cfgflag.KeyFile,
	)
	authMiddleware := middleware.NewAuthMiddleware(logger)
	handlerInstance := handler.NewHandler(usecases, httpServer, authMiddleware, logger)
	handlerInstance.InitRouter(cfgfile)
	httpServer.Start()

	scheme := map[bool]string{true: "https", false: "http"}[cfgflag.EnableHttps]
	logger.Info().Str("scheme", scheme).Int("port", cfgflag.ListenPort).Msg("app server starting")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Info().Str("signal", s.String()).Msg("app server stopping")
	case err := <-httpServer.Notify():
		logger.Error().Err(err).Msg("app server stopping")
	}

	if err := httpServer.Stop(); err != nil {
		logger.Error().Err(err).Msg("app server stopping")
	}
}
