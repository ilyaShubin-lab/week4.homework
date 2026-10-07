package part

import (
	"github.com/stretchr/testify/mock"

	"boilerplates/inventory/internal/model"
)

func (s *ServiceSuite) TestGetSuccess() {
	partUUID := "part-1"
	expected := model.Part{
		UUID:  partUUID,
		Name:  "Супер пупер двигатель",
		Price: 22233,
	}

	s.partRepository.EXPECT().
		Get(mock.Anything, partUUID).
		Return(expected, nil)

	part, err := s.service.Get(s.T().Context(), partUUID)

	s.Require().NoError(err)
	s.Require().Equal(expected, part)
}
