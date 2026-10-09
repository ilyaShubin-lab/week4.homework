//go:build integration

package order

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/suite"

	"boilerplates/order/internal/model"
	pgMigrator "boilerplates/platform/pkg/migrator/pg"
	"boilerplates/platform/pkg/testcontainers/path"
	"boilerplates/platform/pkg/testcontainers/postgres"
)

type RepositorySuite struct {
	suite.Suite

	pg   *postgres.Container
	repo *repository
}

func TestRepositoryIntegration(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}

func (s *RepositorySuite) SetupSuite() {
	ctx := context.Background()

	pg, err := postgres.NewContainer(ctx,
		postgres.WithDatabase("order_test"),
		postgres.WithAuth("order_test", "order_test"),
	)
	s.Require().NoError(err)
	s.pg = pg

	migrationsDir := filepath.Join(path.GetProjectRoot(), "order", "migrations")
	db := stdlib.OpenDBFromPool(pg.Pool())
	s.Require().NoError(pgMigrator.NewMigrator(db, migrationsDir).Up(ctx))
	s.Require().NoError(db.Close())

	s.repo = NewRepository(pg.Pool())
}

func (s *RepositorySuite) TearDownSuite() {
	if s.pg != nil {
		s.Require().NoError(s.pg.Terminate(context.Background()))
	}
}

func (s *RepositorySuite) SetupTest() {
	_, err := s.pg.Pool().Exec(context.Background(), "TRUNCATE TABLE orders")
	s.Require().NoError(err)
}

func newOrder() model.Order {
	return model.Order{
		OrderUUID:  uuid.NewString(),
		UserUUID:   uuid.NewString(),
		PartUUIDs:  []string{uuid.NewString(), uuid.NewString()},
		TotalPrice: 1500.5,
		Status:     model.OrderStatusPendingPayment,
	}
}

func (s *RepositorySuite) TestCreateAndGet() {
	ctx := context.Background()
	order := newOrder()

	s.Require().NoError(s.repo.Create(ctx, order))

	got, err := s.repo.Get(ctx, order.OrderUUID)
	s.Require().NoError(err)
	s.Equal(order, got)
}

func (s *RepositorySuite) TestGetNotFound() {
	_, err := s.repo.Get(context.Background(), uuid.NewString())
	s.ErrorIs(err, model.ErrOrderNotFound)
}

func (s *RepositorySuite) TestUpdatePaid() {
	ctx := context.Background()
	order := newOrder()
	s.Require().NoError(s.repo.Create(ctx, order))

	txUUID := uuid.NewString()
	method := model.PaymentMethodCard
	order.Status = model.OrderStatusPaid
	order.TransactionUUID = &txUUID
	order.PaymentMethod = &method

	s.Require().NoError(s.repo.Update(ctx, order))

	got, err := s.repo.Get(ctx, order.OrderUUID)
	s.Require().NoError(err)
	s.Equal(order, got)
}

func (s *RepositorySuite) TestUpdateNotFound() {
	err := s.repo.Update(context.Background(), newOrder())
	s.ErrorIs(err, model.ErrOrderNotFound)
}

func (s *RepositorySuite) TestCreateInvalidStatus() {
	order := newOrder()
	order.Status = "SHIPPED_TO_MARS"

	err := s.repo.Create(context.Background(), order)
	s.Error(err)
}
