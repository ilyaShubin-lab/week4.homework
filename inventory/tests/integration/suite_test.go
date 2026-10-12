//go:build integration

package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"boilerplates/platform/pkg/logger"
	inventoryV1 "boilerplates/shared/pkg/proto/inventory/v1"
)

type InventorySuite struct {
	suite.Suite

	env    *TestEnvironment
	conn   *grpc.ClientConn
	client inventoryV1.InventoryServiceClient
}

func TestInventoryE2E(t *testing.T) {
	suite.Run(t, new(InventorySuite))
}

func (s *InventorySuite) SetupSuite() {
	s.Require().NoError(logger.Init("info", false))

	ctx := context.Background()

	env, err := setupTestEnvironment(ctx)
	s.Require().NoError(err)
	s.env = env

	s.conn, err = grpc.NewClient(env.App.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	s.Require().NoError(err)

	s.client = inventoryV1.NewInventoryServiceClient(s.conn)
}

func (s *InventorySuite) TearDownSuite() {
	if s.conn != nil {
		_ = s.conn.Close()
	}
	teardownTestEnvironment(context.Background(), s.env)
}
