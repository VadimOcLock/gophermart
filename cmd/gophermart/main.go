package main

import (
	"context"
	"os"
	"time"

	"github.com/VadimOcLock/gophermart/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()

	// Config.
	cfg, err := config.Load[config.WebServer]()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration.")
	}

	// Flags.
	cfg, err = config.ParseFlags(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse configuration flags.")
	}

	// Logger.
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).
		With().Timestamp().Logger()

	// Connection.

	// Migrations.
	// Store.
	// Service.
	// UseCase.
	// Handler.
	// Server.
	// Run.
}
