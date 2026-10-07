package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	inventoryV1API "boilerplates/inventory/internal/api/inventory/v1"
	partRepository "boilerplates/inventory/internal/repository/part"
	partService "boilerplates/inventory/internal/service/part"
	inventoryV1 "boilerplates/shared/pkg/proto/inventory/v1"
)

const grpcPort = "50051"

func main() {
	if err := run(); err != nil {
		log.Fatalf("inventory service: %v", err)
	}
}

func run() error {
	ctx := context.Background()
	mongoURI := getEnv("MONGO_URI",
		"mongodb://inventory-service-user:inventory-service-password@localhost:27017/?authSource=admin")
	mongoDB := getEnv("MONGO_DATABASE", "inventory-service")

	// --- подключение к Mongo
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return fmt.Errorf("mongo connect: %w", err)
	}
	defer func() {
		closeErr := client.Disconnect(ctx)
		if closeErr != nil {
			log.Printf("mongo disconnect: %v", closeErr)
		}
	}()

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err = client.Ping(pingCtx, nil)
	if err != nil {
		return fmt.Errorf("mongo ping: %w", err)
	}
	log.Println("connected to MongoDB")

	// --- хранилище
	repo := partRepository.NewRepository(client.Database(mongoDB))
	err = repo.InitParts(ctx)
	if err != nil {
		return fmt.Errorf("init parts: %w", err)
	}
	// --- бизнес-логика
	svc := partService.NewService(repo)
	// --- транспорт
	apiV1 := inventoryV1API.NewAPI(svc)
	// --- сеть
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	s := grpc.NewServer()

	// вторая проверка соответствия контракту
	inventoryV1.RegisterInventoryServiceServer(s, apiV1)

	// Чтобы grpcurl / Postman видели список методов без .proto
	reflection.Register(s)

	log.Printf("gRPC server listening on :%s", grpcPort)
	// Serve блокируется, пока сервер работает. Если упадёт — ошибка уйдёт в main.
	return s.Serve(lis)
}

func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v != "" {
		return v
	}
	return fallback
}
