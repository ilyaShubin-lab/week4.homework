package payment

import (
	"github.com/google/uuid"

	"boilerplates/payment/internal/model"
)

func (s *ServiceSuite) TestPaySucces() {
	transactionUUID, err := s.service.Pay(s.T().Context(), "order-1", "user-1", "PAYMENT_METHOD_CARD")
	s.Require().NoError(err)
	s.Require().NotEmpty(transactionUUID)

	_, parseErr := uuid.Parse(transactionUUID)
	s.Require().NoError(parseErr)
}

func (s *ServiceSuite) TestPayInvalidMethod() {
	transactionUUID, err := s.service.Pay(s.T().Context(), "order-1", "user-1", "")
	s.Require().ErrorIs(err, model.ErrInvalidPaymentMethod)
	s.Require().Empty(transactionUUID)
}
