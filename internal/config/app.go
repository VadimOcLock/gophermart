package config

import (
	"strconv"
	"strings"
	"time"
)

type AppConfig struct {
	JWTConfig
	HealthConfig
	PrometheusConfig
}

type HealthConfig struct {
	Port int `env:"HEALTH_PORT" envDefault:"9090"`
}

type PrometheusConfig struct {
	Port int `env:"PROMETHEUS_PORT" envDefault:"9091"`
}

type JWTConfig struct {
	SecretKey       string        `env:"SECRET_KEY"`
	TokenExpiration time.Duration `env:"TOKEN_EXPIRATION" envDefault:"24h"`
}

type WebServerConfig struct {
	SrvAddr string `env:"RUN_ADDRESS"`
}

func (cfg WebServerConfig) Port() int {
	s := strings.Split(cfg.SrvAddr, ":")
	if len(s) != 2 {
		return 0
	}
	port, err := strconv.Atoi(s[1])
	if err != nil {
		return 0
	}

	return port
}

type DatabaseConfig struct {
	DSN string `env:"DATABASE_URI"`
}

type AccrualConfig struct {
	SrvAddr string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}
