package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"boilerplates/order/internal/model"
	"boilerplates/order/internal/repository/converter"
	repoModel "boilerplates/order/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, orderUUID string) (model.Order, error) {
	const query = `
		SELECT order_uuid, user_uuid, part_uuids, total_price,
		       transaction_uuid, payment_method, status
		FROM orders
		WHERE order_uuid = $1`

	var o repoModel.Order

	err := r.pool.QueryRow(ctx, query, orderUUID).Scan(
		&o.OrderUUID,
		&o.UserUUID,
		&o.PartUUIDs,
		&o.TotalPrice,
		&o.TransactionUUID,
		&o.PaymentMethod,
		&o.Status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Order{}, model.ErrOrderNotFound
		}
		return model.Order{}, fmt.Errorf("select order: %w", err)
	}
	return converter.OrderToModel(o), nil
}
