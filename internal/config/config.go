package config

import "github.com/caarlos0/env/v11"

type WebServer struct {
	AppConfig
	WebServerConfig
	DatabaseConfig
	AccrualConfig
}

func Load[T any]() (T, error) {
	var cfg T
	err := env.Parse(&cfg)

	return cfg, err
}
