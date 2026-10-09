package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"boilerplates/order/internal/api/health"
	"boilerplates/order/internal/config"
	"boilerplates/platform/pkg/closer"
	"boilerplates/platform/pkg/logger"
	orderV1 "boilerplates/shared/pkg/openapi/order/v1"
)

type App struct {
	diContainer *diContainer
	httpServer  *http.Server
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	return a.runHTTPServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initHTTPServer,
	}

	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJSON(),
	)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	orderServer, err := orderV1.NewServer(a.diContainer.OrderV1API(ctx))
	if err != nil {
		return fmt.Errorf("create ogen server: %w", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", health.Handler)
	mux.Handle("/", orderServer)

	a.httpServer = &http.Server{
		Addr:    config.AppConfig().OrderHTTP.Address(),
		Handler: mux,

		ReadHeaderTimeout: config.AppConfig().OrderHTTP.ReadTimeout(),
	}

	closer.AddNamed("HTTP server", func(ctx context.Context) error {
		if err := a.httpServer.Shutdown(ctx); err != nil {

			if closeErr := a.httpServer.Close(); closeErr != nil {
				logger.Error(ctx, "failed to force close HTTP server", zap.Error(closeErr))
			}
			return fmt.Errorf("http server shutdown: %w", err)
		}
		return nil
	})

	return nil
}

func (a *App) runHTTPServer(ctx context.Context) error {
	logger.Info(ctx, " HTTP server listening", zap.String("address", a.httpServer.Addr))

	err := a.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http serve: %w", err)
	}

	return nil
}
