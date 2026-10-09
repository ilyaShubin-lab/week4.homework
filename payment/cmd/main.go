package main

import (
	"context"
	"fmt"
	"syscall"
	"time"

	"go.uber.org/zap"

	"boilerplates/payment/internal/app"
	"boilerplates/payment/internal/config"
	"boilerplates/platform/pkg/closer"
	"boilerplates/platform/pkg/logger"
)

const (
	configPath      = "./deploy/compose/payment/.env"
	shutdownTimeout = 5 * time.Second
)

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	ctx := context.Background()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	defer gracefulShutdown()

	a, err := app.New(ctx)
	if err != nil {
		logger.Error(ctx, "failed to create app", zap.Error(err))
		return
	}

	if err = a.Run(ctx); err != nil {
		logger.Error(ctx, "app run failed", zap.Error(err))
	}
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "shutdown error", zap.Error(err))
	}
}
