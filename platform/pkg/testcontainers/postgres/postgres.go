package postgres

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	postgresPort   = "5432/tcp"
	startupTimeout = time.Minute

	defaultImage = "postgres:17.0-alpine3.20"
)

type Container struct {
	container testcontainers.Container
	pool      *pgxpool.Pool
	dsn       string
}

func NewContainer(ctx context.Context, opts ...Option) (_ *Container, err error) {
	cfg := &Config{
		ImageName: defaultImage,
		Database:  "test",
		User:      "test",
		Password:  "test",
	}
	for _, opt := range opts {
		opt(cfg)
	}

	req := testcontainers.ContainerRequest{
		Image:        cfg.ImageName,
		ExposedPorts: []string{postgresPort},

		Env: map[string]string{
			"POSTGRES_DB":       cfg.Database,
			"POSTGRES_USER":     cfg.User,
			"POSTGRES_PASSWORD": cfg.Password,
		},

		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(startupTimeout),
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("start postgres container: %w", err)
	}

	defer func() {
		if err == nil {
			return
		}
		if termErr := c.Terminate(ctx); termErr != nil {
			err = errors.Join(err, fmt.Errorf("terminate postgres container: %w", termErr))
		}
	}()

	host, err := c.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("get container host: %w", err)
	}

	port, err := c.MappedPort(ctx, postgresPort)
	if err != nil {
		return nil, fmt.Errorf("get mapped port: %w", err)
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.User, cfg.Password, net.JoinHostPort(host, port.Port()), cfg.Database,
	)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return &Container{container: c, pool: pool, dsn: dsn}, nil
}

func (c *Container) Pool() *pgxpool.Pool { return c.pool }

func (c *Container) DSN() string { return c.dsn }

func (c *Container) Terminate(ctx context.Context) error {
	c.pool.Close()
	return c.container.Terminate(ctx)
}
