package order

import (
	"errors"

	"github.com/stretchr/testify/mock"

	"boilerplates/order/internal/model"
)

func pendingOrder() model.Order {
	return model.Order{
		OrderUUID:  "order-1",
		UserUUID:   "user-1",
		PartUUIDs:  []string{"p1"},
		TotalPrice: 100,
		Status:     model.OrderStatusPendingPayment,
	}
}

func (s *ServiceSuite) TestPaySuccess() {
	order := pendingOrder()
	txUUID := "tx-123"

	s.orderRepository.EXPECT().Get(mock.Anything, "order-1").Return(order, nil)

	s.paymentClient.EXPECT().
		PayOrder(mock.Anything, "order-1", "user-1", model.PaymentMethodCard).
		Return(txUUID, nil)

	// Ключевая проверка: в хранилище уезжает заказ с новым статусом И транзакцией.
	s.orderRepository.EXPECT().
		Update(mock.Anything, mock.MatchedBy(func(o model.Order) bool {
			return o.Status == model.OrderStatusPaid &&
				o.TransactionUUID != nil && *o.TransactionUUID == txUUID &&
				o.PaymentMethod != nil && *o.PaymentMethod == model.PaymentMethodCard
		})).
		Return(nil)

	res, err := s.service.Pay(s.T().Context(), "order-1", model.PaymentMethodCard)

	s.Require().NoError(err)
	s.Require().Equal(txUUID, res)
}

func (s *ServiceSuite) TestPayOrderNotFound() {
	s.orderRepository.EXPECT().
		Get(mock.Anything, "unknown").
		Return(model.Order{}, model.ErrOrderNotFound)

	// Ни payment, ни Update не настроены — до них дойти не должно.

	_, err := s.service.Pay(s.T().Context(), "unknown", model.PaymentMethodCard)

	s.Require().ErrorIs(err, model.ErrOrderNotFound)
}

func (s *ServiceSuite) TestPayAlreadyPaid() {
	order := pendingOrder()
	order.Status = model.OrderStatusPaid

	s.orderRepository.EXPECT().Get(mock.Anything, "order-1").Return(order, nil)

	// Самое важное в этом тесте: paymentClient НЕ настроен.
	// Если сервис всё-таки спишет деньги второй раз — падение.

	_, err := s.service.Pay(s.T().Context(), "order-1", model.PaymentMethodCard)

	s.Require().ErrorIs(err, model.ErrOrderAlreadyPaid)
}

func (s *ServiceSuite) TestPayCancelledOrder() {
	order := pendingOrder()
	order.Status = model.OrderStatusCancelled

	s.orderRepository.EXPECT().Get(mock.Anything, "order-1").Return(order, nil)

	_, err := s.service.Pay(s.T().Context(), "order-1", model.PaymentMethodCard)

	s.Require().ErrorIs(err, model.ErrOrderCancelled)
}

func (s *ServiceSuite) TestPayPaymentUnavailable() {
	s.orderRepository.EXPECT().Get(mock.Anything, "order-1").Return(pendingOrder(), nil)

	s.paymentClient.EXPECT().
		PayOrder(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return("", model.ErrPaymentUnavailable)

	// Update не настроен: платёж не прошёл — статус менять нельзя.

	_, err := s.service.Pay(s.T().Context(), "order-1", model.PaymentMethodCard)

	s.Require().ErrorIs(err, model.ErrPaymentUnavailable)
}

func (s *ServiceSuite) TestPayUpdateError() {
	repoErr := errors.New("storage failed")

	s.orderRepository.EXPECT().Get(mock.Anything, "order-1").Return(pendingOrder(), nil)
	s.paymentClient.EXPECT().PayOrder(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return("tx-123", nil)
	s.orderRepository.EXPECT().Update(mock.Anything, mock.Anything).Return(repoErr)

	_, err := s.service.Pay(s.T().Context(), "order-1", model.PaymentMethodCard)

	s.Require().ErrorIs(err, repoErr)
}
