package payment

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"boilerplates/payment/internal/model"
	"boilerplates/platform/pkg/logger"
)

func (s *service) Pay(ctx context.Context, orderUUID, userUUID, paymentMethod string) (string, error) {
	if paymentMethod == "" || paymentMethod == "PAYMENT_METHOD_UNSPECIFIED" {
		return "", model.ErrInvalidPaymentMethod
	}
	transactionUUID := uuid.NewString()

	logger.Info(ctx, "payment processed",
		zap.String("order_uuid", orderUUID),
		zap.String("user_uuid", userUUID),
		zap.String("payment_method", paymentMethod),
		zap.String("transaction_uuid", transactionUUID),
	)

	return transactionUUID, nil
}
