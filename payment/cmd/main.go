package main

import (
	"context"
	"fmt"
	"net"

	paymentv1API "boilerplates/payment/internal/api/payment/v1"
	"boilerplates/payment/internal/config"
	paymentService "boilerplates/payment/internal/service/payment"
	"boilerplates/platform/pkg/logger"
	paymentv1 "boilerplates/shared/pkg/proto/payment/v1"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const configPath = "./deploy/compose/payment/.env"

func main() {
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	err = logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJSON(),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to init logger: %v", err))
	}

	// --- бизнес-логика
	svc := paymentService.NewService()
	// --- транспорт
	apiV1 := paymentv1API.NewAPI(svc)
	// --- сеть
	lis, err := net.Listen("tcp", config.AppConfig().PaymentGRPC.Address())
	if err != nil {
		logger.Fatal(ctx, "failed to listen", zap.Error(err))
	}

	s := grpc.NewServer()

	// вторая проверка соответствия контракту
	paymentv1.RegisterPaymentServiceServer(s, apiV1)

	// Чтобы grpcurl / Postman видели список методов без .proto
	reflection.Register(s)

	logger.Info(ctx, "gRPC server listening", zap.String("address", lis.Addr().String()))

	err = s.Serve(lis)
	if err != nil {
		logger.Fatal(ctx, "failed to serve", zap.Error(err))
	}
}
