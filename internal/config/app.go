package config

type AppConfig struct {
	SecretKey string `env:"SECRET_KEY,required"`
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
