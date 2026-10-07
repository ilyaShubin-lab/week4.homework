package payment

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"boilerplates/platform/pkg/logger"
)

type ServiceSuite struct {
	suite.Suite
	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.service = NewService()
}

func (s *ServiceSuite) SetupSuite() {
	logger.SetNopLogger()
}

func TestServiceSuite(t *testing.T) { suite.Run(t, new(ServiceSuite)) }
