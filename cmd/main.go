package main

import (
	"github.com/willbrid/easy-api-prom-alert-sms/config"
	"github.com/willbrid/easy-api-prom-alert-sms/internal/app"
	"github.com/willbrid/easy-api-prom-alert-sms/pkg/logging"

	"github.com/go-playground/validator/v10"
)

func main() {
	validate := validator.New(validator.WithRequiredStructEnabled())
	logger := logging.InitLogger()

	configFlag, err := config.LoadConfigFlag(validate, logger)
	if err != nil {
		logger.Error().Err(err).Msg("failed to load configuration flags")
	}

	viperInstance, err := config.ReadConfigFile(configFlag.ConfigFile, logger)
	if err != nil {
		logger.Error().Err(err).Msg("failed to read configuration file")
	}

	configLoaded, err := config.LoadConfig(viperInstance, validate, logger)
	if err != nil {
		logger.Error().Err(err).Msg("failed to load configuration file")
	}

	logger.Info().Str("config_file", configFlag.ConfigFile).Msg("configuration file was loaded successfully")
	app.Run(configLoaded, configFlag, logger)
}
