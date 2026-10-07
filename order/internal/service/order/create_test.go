package order

import (
	"github.com/stretchr/testify/mock"

	"boilerplates/order/internal/model"
)

func (s *ServiceSuite) TestCreateSuccess() {
	partUUIDs := []string{"p1", "p2"}

	s.inventoryClient.EXPECT().
		ListParts(mock.Anything, model.PartsFilter{UUIDs: partUUIDs}).
		Return([]model.Part{{UUID: "p1", Price: 100}, {UUID: "p2", Price: 50}}, nil)

	s.orderRepository.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(nil)

	order, err := s.service.Create(s.T().Context(), "user-1", partUUIDs)

	s.Require().NoError(err)
	s.Require().Equal(float64(150), order.TotalPrice)
	s.Require().Equal(model.OrderStatusPendingPayment, order.Status)
}

func (s *ServiceSuite) TestCreatePartsNotFound() {
	// Вернули одну деталь вместо двух — заказ создаваться не должен.
	s.inventoryClient.EXPECT().ListParts(mock.Anything, mock.Anything).
		Return([]model.Part{{UUID: "p1", Price: 100}}, nil)

	_, err := s.service.Create(s.T().Context(), "user-1", []string{"p1", "p2"})

	s.Require().ErrorIs(err, model.ErrPartsNotFound)
}
