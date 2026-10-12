//go:build integration

package integration

import (
	"boilerplates/platform/pkg/testcontainers/app"
	"boilerplates/platform/pkg/testcontainers/mongo"
	"boilerplates/platform/pkg/testcontainers/network"
)

type TestEnvironment struct {
	Network *network.Network
	Mongo   *mongo.Container
	App     *app.Container
}
