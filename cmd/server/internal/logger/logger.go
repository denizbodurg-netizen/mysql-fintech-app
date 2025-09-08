package logger

import (
	"os"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func New(env string) zerolog.logger{
	if env == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	return log.Logger
}