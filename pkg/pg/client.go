package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/matchsystems/werr"
)

func New(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, werr.Wrap(err)
	}

	config.ConnConfig, err = NewConnConfig(cfg)
	if err != nil {
		return nil, werr.Wrap(err)
	}

	connConfig, err := NewConnConfig(cfg)
	if err != nil {
		return nil, werr.Wrap(err)
	}

	conn, err := pgx.ConnectConfig(ctx, connConfig)
	if err != nil {
		return nil, werr.Wrap(err)
	}
	defer func() {
		_ = conn.Close(ctx)
	}()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, werr.Wrap(err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, werr.Wrap(err)
	}

	return pool, nil
}

func NewConn(ctx context.Context, cfg Config) (*pgx.Conn, error) {
	connConfig, err := NewConnConfig(cfg)
	if err != nil {
		return nil, werr.Wrap(err)
	}

	conn, err := pgx.ConnectConfig(ctx, connConfig)
	if err != nil {
		return nil, werr.Wrap(err)
	}

	if err = conn.Ping(ctx); err != nil {
		return nil, werr.Wrap(err)
	}

	return conn, nil
}

func NewConnConfig(cfg Config) (*pgx.ConnConfig, error) {
	config, err := pgx.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, werr.Wrap(err)
	}

	return config, nil
}
