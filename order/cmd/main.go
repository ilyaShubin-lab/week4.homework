package main

import (
	"context"
	"database/sql"
	"fmt" // ← добавился: fmt.Errorf
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1API "boilerplates/order/internal/api/order/v1"
	inventoryClient "boilerplates/order/internal/client/grpc/inventory/v1"
	paymentClient "boilerplates/order/internal/client/grpc/payment/v1"
	"boilerplates/order/internal/migrator"
	orderRepository "boilerplates/order/internal/repository/order"
	orderService "boilerplates/order/internal/service/order"
	orderV1 "boilerplates/shared/pkg/openapi/order/v1"
	inventoryv1 "boilerplates/shared/pkg/proto/inventory/v1"
	paymentv1 "boilerplates/shared/pkg/proto/payment/v1"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("order service: %v", err)
	}
}

func run() error {
	ctx := context.Background()

	dsn := getEnv("POSTGRES_DSN",
		"postgres://order-service-user:order-service-password@localhost:5432/order-service?sslmode=disable",
	)
	migrationsDir := getEnv("MIGRATIONS_DIR", "order/migrations")

	if err := runMigrations(dsn, migrationsDir); err != nil {
		return err
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return fmt.Errorf("create pgx pool: %w", err)
	}
	defer pool.Close()

	if err = pool.Ping(ctx); err != nil {
		return fmt.Errorf("ping pgx pool: %w", err)
	}

	inventoryConn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("connect to inventory: %w", err)
	}
	defer func() {
		if closeErr := inventoryConn.Close(); closeErr != nil {
			log.Printf("close inventory connection: %v", closeErr)
		}
	}()

	paymentConn, err := grpc.NewClient(
		"localhost:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("connect to payment: %w", err)
	}
	defer func() {
		if closeErr := paymentConn.Close(); closeErr != nil {
			log.Printf("close payment connection: %v", closeErr)
		}
	}()

	invClient := inventoryClient.NewClient(inventoryv1.NewInventoryServiceClient(inventoryConn))
	payClient := paymentClient.NewClient(paymentv1.NewPaymentServiceClient(paymentConn))

	repo := orderRepository.NewRepository(pool)
	service := orderService.NewService(repo, invClient, payClient)
	api := orderV1API.NewAPI(service)

	orderServer, err := orderV1.NewServer(api)
	if err != nil {
		return fmt.Errorf("create openapi server: %w", err)
	}

	httpServer := &http.Server{
		Addr:              ":8080",
		Handler:           orderServer,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("HTTP server listening on :8080")

	return httpServer.ListenAndServe()
}

func runMigrations(dsn, dir string) error {
	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("open db for migrations: %w", err)
	}
	defer func() {
		if closeErr := sqlDB.Close(); closeErr != nil {
			log.Printf("close migrations db: %v", closeErr)
		}
	}()

	if err = sqlDB.Ping(); err != nil {
		return fmt.Errorf("ping db: %w", err)
	}

	if err = migrator.NewMigrator(sqlDB, dir).Up(); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}

	log.Println("migrations applied")
	return nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
