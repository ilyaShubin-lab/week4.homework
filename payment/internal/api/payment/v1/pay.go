package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"boilerplates/payment/internal/model"
	paymentv1 "boilerplates/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	transactionUUID, err := a.paymentService.Pay(ctx, req.GetOrderUuid(),
		req.GetUserUuid(), req.GetPaymentMethod().String(),
	)
	if err != nil {
		if errors.Is(err, model.ErrInvalidPaymentMethod) {
			return nil, status.Error(codes.InvalidArgument, "invelid payment method")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &paymentv1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}
