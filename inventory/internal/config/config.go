package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"boilerplates/inventory/internal/config/env"
)

var appConfig *config

type config struct {
	Logger        LoggerConfig
	InventoryGRPC InventoryGRPCConfig
	Mongo         MongoConfig
}

func Load(path ...string) error {
	// файл → переменные окружения
	err := godotenv.Load(path...)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("load .env file: %w", err)
	}

	// переменные окружения → типизированные конфиги, по одному на блок
	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return fmt.Errorf("logger config: %w", err)
	}

	inventoryGRPCCfg, err := env.NewInventoryGRPCConfig()
	if err != nil {
		return fmt.Errorf("inventory grpc config: %w", err)
	}

	// Mongo
	mongoCfg, err := env.NewMongoConfig()
	if err != nil {
		return fmt.Errorf("mongo config: %w", err)
	}

	// всё собрали
	appConfig = &config{
		Logger:        loggerCfg,
		InventoryGRPC: inventoryGRPCCfg,
		Mongo:         mongoCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
