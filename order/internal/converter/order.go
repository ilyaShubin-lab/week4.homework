package converter

import (
	"github.com/google/uuid"

	"boilerplates/order/internal/model"
	orderV1 "boilerplates/shared/pkg/openapi/order/v1"
)

func OrderToDTO(order model.Order) orderV1.OrderDto {
	partUUIDs := make([]uuid.UUID, 0, len(order.PartUUIDs))
	for _, p := range order.PartUUIDs {
		partUUIDs = append(partUUIDs, uuid.MustParse(p))
	}

	dto := orderV1.OrderDto{
		OrderUUID:  uuid.MustParse(order.OrderUUID),
		UserUUID:   uuid.MustParse(order.UserUUID),
		PartUuids:  partUUIDs,
		TotalPrice: order.TotalPrice,
		Status:     orderV1.OrderStatus(order.Status),
	}

	if order.TransactionUUID != nil {
		dto.TransactionUUID = orderV1.NewOptUUID(uuid.MustParse(*order.TransactionUUID))
	}
	if order.PaymentMethod != nil {
		dto.PaymentMethod = orderV1.NewOptPaymentMethod(orderV1.PaymentMethod(*order.PaymentMethod))
	}

	return dto
}

func PaymentMethodToModel(method orderV1.PaymentMethod) model.PaymentMethod {
	switch method {
	case orderV1.PaymentMethodPAYMENTMETHODCARD:
		return model.PaymentMethodCard
	case orderV1.PaymentMethodPAYMENTMETHODSBP:
		return model.PaymentMethodSBP
	case orderV1.PaymentMethodPAYMENTMETHODCREDITCARD:
		return model.PaymentMethodCreditCard
	case orderV1.PaymentMethodPAYMENTMETHODINVESTORMONEY:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnknown
	}
}
