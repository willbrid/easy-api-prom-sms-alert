package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

func InitLogger() zerolog.Logger {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = time.RFC3339

	logger := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Caller().
		Logger()

	return logger
}
