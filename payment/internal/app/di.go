package app

import (
	paymentV1API "boilerplates/payment/internal/api/payment/v1"
	"boilerplates/payment/internal/service"
	paymentService "boilerplates/payment/internal/service/payment"
	paymentV1 "boilerplates/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentv1API   paymentV1.PaymentServiceServer
	paymentService service.PaymentService
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) PaymentV1API() paymentV1.PaymentServiceServer {
	if d.paymentv1API == nil {
		d.paymentv1API = paymentV1API.NewAPI(d.PaymentService())
	}
	return d.paymentv1API
}

func (d *diContainer) PaymentService() service.PaymentService {
	if d.paymentService == nil {
		d.paymentService = paymentService.NewService()
	}
	return d.paymentService
}
