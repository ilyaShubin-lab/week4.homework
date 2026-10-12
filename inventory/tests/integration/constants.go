//go:build integration

package integration

import "time"

const (
	projectName = "inventory-e2e"

	mongoContainerName = "inventory-e2e-mongo"
	appContainerName   = "inventory-e2e-app"

	inventoryDockerfile = "deploy/docker/inventory/Dockerfile"

	grpcPort = "50051"

	startupTimeout = 3 * time.Minute

	envFilePath = "deploy/compose/inventory/.env"
)

const (
	engineUUID   = "550e8400-e29b-41d4-a716-446655440001"
	engineName   = "Двигатель импульсный ИД-500"
	fuelTankUUID = "550e8400-e29b-41d4-a716-446655440002"
	portholeUUID = "550e8400-e29b-41d4-a716-446655440003"

	seededPartsCount = 3
)
