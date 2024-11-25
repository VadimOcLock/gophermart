package config

import "time"

type AppConfig struct {
	JWTConfig
}

type JWTConfig struct {
	SecretKey       string        `env:"SECRET_KEY"`
	TokenExpiration time.Duration `env:"TOKEN_EXPIRATION" envDefault:"24h"`
}

type WebServerConfig struct {
	SrvAddr string `env:"RUN_ADDRESS"`
}

type DatabaseConfig struct {
	DSN string `env:"DATABASE_URI"`
}

type AccrualConfig struct {
	SrvAddr string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}
