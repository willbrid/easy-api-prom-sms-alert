package app

import (
	"github.com/willbrid/easy-api-prom-alert-sms/config"
	"github.com/willbrid/easy-api-prom-alert-sms/internal/domain"
	"github.com/willbrid/easy-api-prom-alert-sms/internal/handler"
	"github.com/willbrid/easy-api-prom-alert-sms/internal/microservice"
	"github.com/willbrid/easy-api-prom-alert-sms/internal/usecase"
	"github.com/willbrid/easy-api-prom-alert-sms/pkg/httpserver"
	"github.com/willbrid/easy-api-prom-alert-sms/pkg/logger"

	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfgfile *config.Config, cfgflag *config.ConfigFlag, loggerInstance logger.ILogger) {
	microservices := microservice.NewMicroservices(&cfgfile.EasyAPIPromAlertSMS.Provider, loggerInstance)
	usecases := usecase.NewUsecases(&usecase.Deps{
		Microservices: microservices,
		AlertConfig: &domain.AlertConfig{
			Recipients:           cfgfile.EasyAPIPromAlertSMS.Recipients,
			DefaultRecipientName: cfgfile.EasyAPIPromAlertSMS.Provider.Parameters.To.ParamValue,
			Simulation:           cfgfile.EasyAPIPromAlertSMS.Simulation,
		},
		ILogger: loggerInstance,
	})

	httpServer := httpserver.NewServer(
		fmt.Sprint(":"+fmt.Sprint(cfgflag.ListenPort)),
		cfgflag.EnableHttps,
		cfgflag.CertFile,
		cfgflag.KeyFile,
	)
	handlerInstance := handler.NewHandler(usecases, httpServer.Router, loggerInstance)
	handlerInstance.InitRouter(cfgfile)
	httpServer.Start()

	scheme := map[bool]string{true: "https", false: "http"}[cfgflag.EnableHttps]
	loggerInstance.Info("app server is listening on port %v using %s", cfgflag.ListenPort, scheme)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		loggerInstance.Info("app server - run - signal: %s", s.String())
	case err := <-httpServer.Notify():
		loggerInstance.Error("app server error: %v", err.Error())
	}

	if err := httpServer.Stop(); err != nil {
		loggerInstance.Error("app server - stop - error: %v", err)
	}
}
