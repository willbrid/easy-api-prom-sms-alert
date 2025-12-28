package main

import (
	"github.com/willbrid/easy-api-prom-alert-sms/config"
	"github.com/willbrid/easy-api-prom-alert-sms/internal/app"
	"github.com/willbrid/easy-api-prom-alert-sms/pkg/logger"

	"github.com/go-playground/validator/v10"
)

func main() {
	var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())
	loggerInstance := logger.NewLogger()

	configFlag, err := config.LoadConfigFlag(validate)
	if err != nil {
		loggerInstance.Fatal("failed to load configuration flags: %v", err.Error())
	}

	viperInstance, err := config.ReadConfigFile(configFlag.ConfigFile)
	if err != nil {
		loggerInstance.Fatal("failed to read configuration file: %v", err.Error())
	}

	configLoaded, err := config.LoadConfig(viperInstance, validate)
	if err != nil {
		loggerInstance.Fatal("failed to load configuration file: %v", err.Error())
	}

	loggerInstance.Info("configuration file '%s' was loaded successfully", configFlag.ConfigFile)

	app.Run(configLoaded, configFlag, loggerInstance)
}
