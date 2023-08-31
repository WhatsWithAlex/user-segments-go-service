package app

import (
	"errors"
	"os"

	"github.com/WhatsWithAlex/user-segments-go-service/internal/api"
	"github.com/WhatsWithAlex/user-segments-go-service/internal/env"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var EnvLoadError = errors.New("app: error loading environment variables")
var DBSetupError = errors.New("app: error setupping database")
var ServerSetupError = errors.New("app: error setupping server")

func Run() (err error) {
	env, err := env.LoadEnv()
	if err != nil {
		return errors.Join(EnvLoadError, err)
	}
	log.Info().Msg("app: environment variables loaded successfully")

	if env.Environment == "development" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		log.Warn().Msg("app: running in development mode")
	} else {
		gin.SetMode("release")
	}

	db, err := setupDB(&env.DB)
	if err != nil {
		return errors.Join(DBSetupError, err)
	}

	ginServer := gin.Default()
	api.SetupRouter(&env, db, ginServer)
	err = ginServer.Run(":" + env.App.Port)
	if err != nil {
		return errors.Join(ServerSetupError, err)
	}
	return
}
