package order

import (
	"testing"

	"github.com/stretchr/testify/suite"

	clientMocks "boilerplates/order/internal/client/grpc/mocks"
	repoMocks "boilerplates/order/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	orderRepository *repoMocks.MockOrderRepository
	inventoryClient *clientMocks.MockInventoryClient
	paymentClient   *clientMocks.MockPaymentClient

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.orderRepository = repoMocks.NewMockOrderRepository(s.T())
	s.inventoryClient = clientMocks.NewMockInventoryClient(s.T())
	s.paymentClient = clientMocks.NewMockPaymentClient(s.T())
	s.service = NewService(s.orderRepository, s.inventoryClient, s.paymentClient)
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
