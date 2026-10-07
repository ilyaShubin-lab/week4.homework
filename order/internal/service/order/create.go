package order

import (
	"context"

	"github.com/google/uuid"

	"boilerplates/order/internal/model"
)

func (s *service) Create(ctx context.Context, userUUID string, partUUIDs []string) (model.Order, error) {
	parts, err := s.inventoryClient.ListParts(ctx, model.PartsFilter{UUIDs: partUUIDs})
	if err != nil {
		return model.Order{}, err
	}

	// заказ создаётся, только если нашлись ВСЕ запрошенные детали.
	if len(parts) != len(partUUIDs) {
		return model.Order{}, model.ErrPartsNotFound
	}

	var total float64
	for _, p := range parts {
		total += p.Price
	}

	order := model.Order{
		OrderUUID:  uuid.NewString(),
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: total,
		Status:     model.OrderStatusPendingPayment,
	}

	err = s.orderRepository.Create(ctx, order)
	if err != nil {
		return model.Order{}, err
	}
	return order, nil
}
