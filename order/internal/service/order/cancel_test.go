package order

import (
	"github.com/stretchr/testify/mock"

	"boilerplates/order/internal/model"
)

func (s *ServiceSuite) TestCancelSuccess() {
	s.orderRepository.EXPECT().Get(mock.Anything, "order-1").Return(pendingOrder(), nil)

	s.orderRepository.EXPECT().
		Update(mock.Anything, mock.MatchedBy(func(o model.Order) bool {
			return o.Status == model.OrderStatusCancelled
		})).
		Return(nil)

	err := s.service.Cancel(s.T().Context(), "order-1")

	s.Require().NoError(err)
}

func (s *ServiceSuite) TestCancelOrderNotFound() {
	s.orderRepository.EXPECT().
		Get(mock.Anything, "unknown").
		Return(model.Order{}, model.ErrOrderNotFound)

	err := s.service.Cancel(s.T().Context(), "unknown")

	s.Require().ErrorIs(err, model.ErrOrderNotFound)
}

func (s *ServiceSuite) TestCancelPaidOrder() {
	order := pendingOrder()
	order.Status = model.OrderStatusPaid

	s.orderRepository.EXPECT().Get(mock.Anything, "order-1").Return(order, nil)

	// Update не настроен — оплаченный заказ трогать нельзя.

	err := s.service.Cancel(s.T().Context(), "order-1")

	s.Require().ErrorIs(err, model.ErrOrderAlreadyPaid)
}

func (s *ServiceSuite) TestCancelAlreadyCancelled() {
	order := pendingOrder()
	order.Status = model.OrderStatusCancelled

	s.orderRepository.EXPECT().Get(mock.Anything, "order-1").Return(order, nil)

	err := s.service.Cancel(s.T().Context(), "order-1")

	s.Require().ErrorIs(err, model.ErrOrderCancelled)
}
