package main

import (
	"os"

	"github.com/WhatsWithAlex/user-segments-go-service/internal/app"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func preCheckEnv() bool {
	_, err := os.Stat(".env")
	return os.Getenv("ENVIRONMENT") != "production" && !os.IsNotExist(err)
}

func main() {
	if preCheckEnv() {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	err := app.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("main: failed to run app")
	}
	log.Info().Msg("main: app is running")
}
