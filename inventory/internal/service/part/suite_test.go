package part

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"boilerplates/inventory/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	partRepository *mocks.MockPartRepository

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.partRepository = mocks.NewMockPartRepository(s.T())

	s.service = NewService(
		s.partRepository,
	)
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
