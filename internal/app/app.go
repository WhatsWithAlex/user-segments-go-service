package app

import (
	"errors"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var EnvLoadError = errors.New("app: error loading environment variables")

func Run() (err error) {
	env, err := LoadEnv()
	if err != nil {
		return errors.Join(EnvLoadError, err)
	}
	log.Info().Msg("app: environment variables loaded successfully")

	if env.Environment == "development" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		log.Warn().Msg("app: running in development mode")
	}

	return
}
