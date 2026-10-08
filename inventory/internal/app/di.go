package app

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	inventoryV1API "boilerplates/inventory/internal/api/inventory/v1"
	"boilerplates/inventory/internal/config"
	"boilerplates/inventory/internal/repository"
	partRepo "boilerplates/inventory/internal/repository/part"
	"boilerplates/inventory/internal/service"
	partSvc "boilerplates/inventory/internal/service/part"
	"boilerplates/platform/pkg/closer"
	inventoryV1 "boilerplates/shared/pkg/proto/inventory/v1"
)

type diContainer struct {
	inventoryV1API inventoryV1.InventoryServiceServer
	partService    service.PartService
	partRepository repository.PartRepository

	mongoDBClient *mongo.Client
	mongoDBHandle *mongo.Database
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) InventoryV1API(ctx context.Context) inventoryV1.InventoryServiceServer {
	if d.inventoryV1API == nil {
		d.inventoryV1API = inventoryV1API.NewAPI(d.PartService(ctx))
	}
	return d.inventoryV1API
}

func (d *diContainer) PartService(ctx context.Context) service.PartService {
	if d.partService == nil {
		d.partService = partSvc.NewService(d.PartRepository(ctx))
	}
	return d.partService
}

func (d *diContainer) PartRepository(ctx context.Context) repository.PartRepository {
	if d.partRepository == nil {
		repo := partRepo.NewRepository(d.MongoDBHandle(ctx))

		if err := repo.InitParts(ctx); err != nil {
			panic(fmt.Sprintf("failed to init parts: %v", err))
		}
		d.partRepository = repo
	}
	return d.partRepository
}

func (d *diContainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if d.mongoDBHandle == nil {
		d.mongoDBHandle = d.MongoDBClient(ctx).Database(config.AppConfig().Mongo.DatabaseName())
	}
	return d.mongoDBHandle
}

func (d *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if d.mongoDBClient == nil {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
		if err != nil {
			panic(fmt.Sprintf("failed to connect MongoDB: %v", err))
		}

		if err = client.Ping(ctx, readpref.Primary()); err != nil {
			panic(fmt.Sprintf("failed to ping MongoDB: %v", err))
		}

		closer.AddNamed("MongoDB client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})
		d.mongoDBClient = client

	}
	return d.mongoDBClient
}
