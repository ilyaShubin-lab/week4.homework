package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"boilerplates/payment/internal/config/env"
)

var appConfig *config

type config struct {
	Logger      LoggerConfig
	PaymentGRPC PaymentGRPCConfig
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

	paymentGRPCCfg, err := env.NewPaymentGRPCConfig()
	if err != nil {
		return fmt.Errorf("payment grpc config: %w", err)
	}

	appConfig = &config{
		Logger:      loggerCfg,
		PaymentGRPC: paymentGRPCCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
