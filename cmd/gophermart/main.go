package main

import (
	"context"
	"os"
	"syscall"
	"time"

	"github.com/VadimOcLock/gophermart/internal/server"
	"github.com/VadimOcLock/gophermart/pkg/lifecycle"
	"github.com/VadimOcLock/gophermart/pkg/migrations"
	"github.com/safeblock-dev/wr/taskgroup"

	"github.com/VadimOcLock/gophermart/pkg/pg"

	"github.com/VadimOcLock/gophermart/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const migrationsFolderPath = "file://schema/migrations"

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
	pgClient, err := pg.New(ctx, pg.Config{
		DSN: cfg.DatabaseConfig.DSN,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database.")
	}
	defer pgClient.Close()
	log.Info().Msg("Connected to database.")

	// Migrations.
	if err = migrations.Run(cfg.DatabaseConfig.DSN, migrationsFolderPath); err != nil {
		log.Fatal().Err(err).Msg("Failed to run database migration.")
	}
	log.Info().Msg("migrations applied successfully")

	// Server.
	srv := server.New(pgClient, cfg)

	// Run.
	tasks := taskgroup.New()
	tasks.Add(taskgroup.SignalHandler(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM))
	tasks.Add(lifecycle.HTTPServer(srv))
	_ = tasks.Run()
}
