package order

import (
	"github.com/stretchr/testify/mock"

	"boilerplates/order/internal/model"
)

func (s *ServiceSuite) TestGetSuccess() {
	expected := model.Order{
		OrderUUID:  "order-1",
		UserUUID:   "user-1",
		PartUUIDs:  []string{"p1"},
		TotalPrice: 100,
		Status:     model.OrderStatusPendingPayment,
	}

	s.orderRepository.EXPECT().
		Get(mock.Anything, "order-1").
		Return(expected, nil)

	order, err := s.service.Get(s.T().Context(), "order-1")

	s.Require().NoError(err)
	s.Require().Equal(expected, order)
}

func (s *ServiceSuite) TestGetNotFound() {
	s.orderRepository.EXPECT().
		Get(mock.Anything, "unknown").
		Return(model.Order{}, model.ErrOrderNotFound)

	_, err := s.service.Get(s.T().Context(), "unknown")

	s.Require().ErrorIs(err, model.ErrOrderNotFound)
}
