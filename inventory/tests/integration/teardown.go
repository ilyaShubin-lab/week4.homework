//go:build integration

package integration

import (
	"context"

	"go.uber.org/zap"

	"boilerplates/platform/pkg/logger"
)

func teardownTestEnvironment(ctx context.Context, env *TestEnvironment) {
	if env == nil {
		return
	}

	if env.App != nil {
		if err := env.App.Terminate(ctx); err != nil {
			logger.Error(ctx, "failed to terminate inventory container", zap.Error(err))
		}
	}

	if env.Mongo != nil {
		if err := env.Mongo.Terminate(ctx); err != nil {
			logger.Error(ctx, "failed to terminate mongo container", zap.Error(err))
		}
	}

	if env.Network != nil {
		if err := env.Network.Remove(ctx); err != nil {
			logger.Error(ctx, "failed to remove network", zap.Error(err))
		}
	}
}
