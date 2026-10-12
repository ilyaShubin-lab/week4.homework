//go:build integration

package integration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/testcontainers/testcontainers-go/wait"

	"boilerplates/platform/pkg/logger"
	"boilerplates/platform/pkg/testcontainers"
	"boilerplates/platform/pkg/testcontainers/app"
	"boilerplates/platform/pkg/testcontainers/mongo"
	"boilerplates/platform/pkg/testcontainers/network"
	"boilerplates/platform/pkg/testcontainers/path"
)

func setupTestEnvironment(ctx context.Context) (env *TestEnvironment, err error) {
	env = &TestEnvironment{}

	defer func() {
		if err != nil {
			teardownTestEnvironment(ctx, env)
		}
	}()

	root := path.GetProjectRoot()

	if err = godotenv.Load(filepath.Join(root, envFilePath)); err != nil {
		return env, fmt.Errorf("load %s (выполни task env:generate): %w", envFilePath, err)
	}

	env.Network, err = network.NewNetwork(ctx, projectName)
	if err != nil {
		return env, fmt.Errorf("create network: %w", err)
	}

	env.Mongo, err = mongo.NewContainer(ctx,
		mongo.WithNetworkName(env.Network.Name()),
		mongo.WithContainerName(mongoContainerName),
		mongo.WithImageName(os.Getenv(testcontainers.MongoImageNameKey)),
		mongo.WithDatabase(os.Getenv(testcontainers.MongoDatabaseKey)),
		mongo.WithAuth(
			os.Getenv(testcontainers.MongoUsernameKey),
			os.Getenv(testcontainers.MongoPasswordKey),
		),
		mongo.WithAuthDB(os.Getenv(testcontainers.MongoAuthDBKey)),
		mongo.WithLogger(logger.Logger()),
	)
	if err != nil {
		return env, fmt.Errorf("start mongo: %w", err)
	}

	appEnv := map[string]string{
		"GRPC_HOST":      "0.0.0.0",
		"GRPC_PORT":      grpcPort,
		"LOGGER_LEVEL":   "debug",
		"LOGGER_AS_JSON": "false",

		testcontainers.MongoHostKey: mongoContainerName,
		testcontainers.MongoPortKey: testcontainers.MongoPort,

		testcontainers.MongoDatabaseKey: os.Getenv(testcontainers.MongoDatabaseKey),
		testcontainers.MongoAuthDBKey:   os.Getenv(testcontainers.MongoAuthDBKey),
		testcontainers.MongoUsernameKey: os.Getenv(testcontainers.MongoUsernameKey),
		testcontainers.MongoPasswordKey: os.Getenv(testcontainers.MongoPasswordKey),
	}

	env.App, err = app.NewContainer(ctx,
		app.WithName(appContainerName),
		app.WithDockerfile(root, inventoryDockerfile),
		app.WithPort(grpcPort),
		app.WithNetwork(env.Network.Name()),
		app.WithEnv(appEnv),
		app.WithLogOutput(os.Stdout),

		app.WithStartupWait(
			wait.ForLog("gRPC server listening").WithStartupTimeout(startupTimeout),
		),
		app.WithLogger(logger.Logger()),
	)
	if err != nil {
		return env, fmt.Errorf("start inventory app: %w", err)
	}

	return env, nil
}
