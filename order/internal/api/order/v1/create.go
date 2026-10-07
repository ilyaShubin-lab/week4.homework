package v1

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"boilerplates/order/internal/model"
	orderV1 "boilerplates/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	partUUIDs := make([]string, 0, len(req.PartUuids))
	for _, p := range req.PartUuids {
		partUUIDs = append(partUUIDs, p.String())
	}

	order, err := a.orderService.Create(ctx, req.UserUUID.String(), partUUIDs)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrPartsNotFound):
			return &orderV1.BadRequestError{Code: 400, Message: "some parts not found"}, nil
		case errors.Is(err, model.ErrInventoryUnavailable):
			return &orderV1.BadGatewayError{Code: 502, Message: "inventory unavailable"}, nil
		default:
			return &orderV1.InternalServerError{Code: 500, Message: "internal error"}, nil
		}
	}

	return &orderV1.CreateOrderResponse{
		UUID:       uuid.MustParse(order.OrderUUID),
		TotalPrice: float32(order.TotalPrice),
	}, nil
}
