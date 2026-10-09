package app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1API "boilerplates/order/internal/api/order/v1"
	grpcClient "boilerplates/order/internal/client/grpc"
	inventoryClient "boilerplates/order/internal/client/grpc/inventory/v1"
	paymentClient "boilerplates/order/internal/client/grpc/payment/v1"
	"boilerplates/order/internal/config"
	"boilerplates/order/internal/repository"
	orderRepository "boilerplates/order/internal/repository/order"
	"boilerplates/order/internal/service"
	orderService "boilerplates/order/internal/service/order"
	"boilerplates/platform/pkg/closer"
	pgMigrator "boilerplates/platform/pkg/migrator/pg"
	orderV1 "boilerplates/shared/pkg/openapi/order/v1"
	inventoryV1 "boilerplates/shared/pkg/proto/inventory/v1"
	paymentV1 "boilerplates/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	orderV1API orderV1.Handler

	orderService    service.OrderService
	orderRepository repository.OrderRepository

	inventoryClient grpcClient.InventoryClient
	paymentClient   grpcClient.PaymentClient

	postgresPool *pgxpool.Pool
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

// ---------- верхний уровень: API ----------

func (d *diContainer) OrderV1API(ctx context.Context) orderV1.Handler {
	if d.orderV1API == nil {
		d.orderV1API = orderV1API.NewAPI(d.OrderService(ctx))
	}

	return d.orderV1API
}

// ---------- бизнес-логика ----------
func (d *diContainer) OrderService(ctx context.Context) service.OrderService {
	if d.orderService == nil {
		d.orderService = orderService.NewService(
			d.OrderRepository(ctx),
			d.InventoryClient(),
			d.PaymentClient(),
		)
	}

	return d.orderService
}

// ---------- хранилище ----------

func (d *diContainer) OrderRepository(ctx context.Context) repository.OrderRepository {
	if d.orderRepository == nil {
		d.orderRepository = orderRepository.NewRepository(d.PostgresPool(ctx))
	}

	return d.orderRepository
}

func (d *diContainer) PostgresPool(ctx context.Context) *pgxpool.Pool {
	if d.postgresPool == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().Postgres.URI())
		if err != nil {
			panic(fmt.Sprintf("failed to create postgres pool: %v", err))
		}

		closer.AddNamed("PostgreSQL pool", func(_ context.Context) error {
			pool.Close()
			return nil
		})
		if err = pool.Ping(ctx); err != nil {
			panic(fmt.Sprintf("failed to ping postgres: %v", err))
		}

		d.runMigrations(ctx, pool)

		d.postgresPool = pool
	}

	return d.postgresPool
}

func (d *diContainer) runMigrations(ctx context.Context, pool *pgxpool.Pool) {

	db := stdlib.OpenDBFromPool(pool)

	m := pgMigrator.NewMigrator(db, config.AppConfig().Postgres.MigrationsDir())
	if err := m.Up(ctx); err != nil {
		panic(fmt.Sprintf("failed to apply migrations: %v", err))
	}
	if err := db.Close(); err != nil {
		panic(fmt.Sprintf("failed to close migrations db: %v", err))
	}
}

// ---------- внешние gRPC-клиенты ----------

func (d *diContainer) InventoryClient() grpcClient.InventoryClient {
	if d.inventoryClient == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().InventoryGRPC.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {

			panic(fmt.Sprintf("failed to create inventory grpc client: %v", err))
		}

		closer.AddNamed("Inventory gRPC connection", func(_ context.Context) error {
			return conn.Close()
		})

		d.inventoryClient = inventoryClient.NewClient(inventoryV1.NewInventoryServiceClient(conn))
	}

	return d.inventoryClient
}

func (d *diContainer) PaymentClient() grpcClient.PaymentClient {
	if d.paymentClient == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().PaymentGRPC.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create payment grpc client: %v", err))
		}

		closer.AddNamed("Payment gRPC connection", func(_ context.Context) error {
			return conn.Close()
		})

		d.paymentClient = paymentClient.NewClient(paymentV1.NewPaymentServiceClient(conn))
	}

	return d.paymentClient
}
