package v1

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"boilerplates/order/internal/converter"
	"boilerplates/order/internal/model"
	orderV1 "boilerplates/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	txUUID, err := a.orderService.Pay(ctx, params.OrderUUID.String(),
		converter.PaymentMethodToModel(req.PaymentMethod))
	if err != nil {
		switch {
		case errors.Is(err, model.ErrOrderNotFound):
			return &orderV1.NotFoundError{Code: 404, Message: "order not found"}, nil
		case errors.Is(err, model.ErrOrderAlreadyPaid), errors.Is(err, model.ErrOrderCancelled):
			return &orderV1.ConflictError{Code: 409, Message: err.Error()}, nil
		default:
			return &orderV1.InternalServerError{Code: 500, Message: "internal error"}, nil
		}
	}
	return &orderV1.PayOrderResponse{TransactionUUID: uuid.MustParse(txUUID)}, nil
}
