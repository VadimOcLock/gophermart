package config

import (
	"flag"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
	"strings"

	"github.com/VadimOcLock/gophermart/internal/errorz"
)

const (
	defaultWebSrvAddr     = "localhost:8080"
	defaultAccrualSrvAddr = "localhost:8081"
	defaultDatabaseDSN    = ""
	defaultSecretKey      = "1234567890"
)

type netAddress struct {
	Host string
	Port int
}

func (n *netAddress) String() string {
	return fmt.Sprintf("%s:%d", n.Host, n.Port)
}

func (n *netAddress) Set(value string) error {
	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		return errorz.ErrInvalidAddressFormat
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid port: %w", err)
	}
	n.Host = parts[0]
	n.Port = port

	return nil
}

func ParseFlags(cfg *WebServer) error {
	var (
		flagWebSrvAddr     string
		flagAccrualSrvAddr string
		flagDatabaseDSN    string
		flagSecretKey      string
	)

	flag.StringVar(&flagWebSrvAddr, "a", defaultWebSrvAddr, "server addr host and port")
	flag.StringVar(&flagAccrualSrvAddr, "r", defaultAccrualSrvAddr, "accrual server addr host and port")
	flag.StringVar(&flagDatabaseDSN, "d", defaultDatabaseDSN, "database dsn")
	flag.StringVar(&flagSecretKey, "k", defaultSecretKey, "database dsn")

	flag.Parse()

	var webSrvAddr netAddress
	if err := webSrvAddr.Set(flagWebSrvAddr); err != nil {
		return fmt.Errorf("error parsing server address: %w", err)
	}
	var accrualSrvAddr netAddress
	if err := accrualSrvAddr.Set(flagAccrualSrvAddr); err != nil {
		return fmt.Errorf("error parsing accrual server address: %w", err)
	}

	if envVal := os.Getenv("RUN_ADDRESS"); envVal == "" {
		cfg.WebServerConfig.SrvAddr = webSrvAddr.String()
	}
	if envVal := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envVal == "" {
		cfg.AccrualConfig.SrvAddr = accrualSrvAddr.String()
	}
	if envVal := os.Getenv("DATABASE_URI"); envVal == "" {
		cfg.DatabaseConfig.DSN = flagDatabaseDSN
	}
	if envVal := os.Getenv("SECRET_KEY"); envVal == "" {
		cfg.AppConfig.SecretKey = flagSecretKey
	}

	log.Info().Msgf("%s", cfg)

	return nil
}
