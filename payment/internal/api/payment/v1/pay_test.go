package v1

import (
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"boilerplates/payment/internal/model"
	paymentv1 "boilerplates/shared/pkg/proto/payment/v1"
)

func (a *APISuite) TestPayOrderSuccess() {
	expectedTransactionUUID := "tx-123"
	// ── СЦЕНАРИЙ ──
	a.paymentService.EXPECT().
		Pay(mock.Anything, "order-1", "user-1", "PAYMENT_METHOD_CARD").
		Return(expectedTransactionUUID, nil)

	// 2. ДЕЙСТВИЕ — одна строка, вызов тестируемого
	resp, err := a.api.PayOrder(a.T().Context(), &paymentv1.PayOrderRequest{
		OrderUuid:     "order-1",
		UserUuid:      "user-1",
		PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
	})

	// 3. ПРОВЕРКИ — что получилось
	a.Require().NoError(err)
	a.Require().Equal(expectedTransactionUUID, resp.GetTransactionUuid())
}

func (a *APISuite) TestInvalidPaymentMethod() {
	// ── СЦЕНАРИЙ ──
	a.paymentService.EXPECT().
		Pay(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return("", model.ErrInvalidPaymentMethod)

	// 2. ДЕЙСТВИЕ — одна строка, вызов тестируемого
	resp, err := a.api.PayOrder(a.T().Context(), &paymentv1.PayOrderRequest{})

	a.Require().Equal(codes.InvalidArgument, status.Code(err))
	a.Require().Empty(resp.GetTransactionUuid())
}
