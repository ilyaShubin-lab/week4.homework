package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"boilerplates/order/internal/config/env"
)

var appConfig *config

type config struct {
	Logger        LoggerConfig
	OrderHTTP     OrderHTTPConfig
	InventoryGRPC InventoryGRPCConfig
	PaymentGRPC   PaymentGRPCConfig
	Postgres      PostgresConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("load .env file: %w", err)
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return fmt.Errorf("logger config: %w", err)
	}
	orderHTTPCfg, err := env.NewOrderHTTPConfig()
	if err != nil {
		return fmt.Errorf("order http config: %w", err)
	}
	inventoryGRPCCfg, err := env.NewInventoryGRPCConfig()
	if err != nil {
		return fmt.Errorf("inventory grpc config: %w", err)
	}
	paymentGRPCCfg, err := env.NewPaymentGRPCConfig()
	if err != nil {
		return fmt.Errorf("payment grpc config: %w", err)
	}

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return fmt.Errorf("postgres config: %w", err)
	}

	appConfig = &config{
		Logger:        loggerCfg,
		OrderHTTP:     orderHTTPCfg,
		InventoryGRPC: inventoryGRPCCfg,
		PaymentGRPC:   paymentGRPCCfg,
		Postgres:      postgresCfg,
	}
	return nil
}

func AppConfig() *config {
	return appConfig
}
