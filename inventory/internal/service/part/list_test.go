package part

import (
	"errors"

	"github.com/stretchr/testify/mock"

	"boilerplates/inventory/internal/model"
)

func (s *ServiceSuite) TestListSucces() {
	filter := model.PartsFilter{
		UUIDs:      []string{"part-1", "part-2"},
		Categories: []model.Category{model.CategoryEngine},
	}

	expected := []model.Part{
		{UUID: "part-1", Name: "Двигатель", Price: 100500},
		{UUID: "part-2", Name: "Иллюминатор", Price: 700},
	}

	s.partRepository.EXPECT().
		List(s.T().Context(), filter).
		Return(expected, nil)

	parts, err := s.service.List(s.T().Context(), filter)

	s.Require().NoError(err)
	s.Require().Equal(expected, parts)
}

func (s *ServiceSuite) TestListEmpty() {
	filter := model.PartsFilter{UUIDs: []string{"unknown"}}

	expected := []model.Part{}
	s.partRepository.EXPECT().
		List(mock.Anything, filter).
		Return(expected, nil)
	parts, err := s.service.List(s.T().Context(), filter)

	s.Require().NoError(err)
	s.Require().Empty(parts)
}

func (s *ServiceSuite) TestListRepoError() {
	filter := model.PartsFilter{}
	repoError := errors.New("storage is down")

	s.partRepository.EXPECT().
		List(mock.Anything, filter).
		Return(nil, repoError)

	parts, err := s.service.List(s.T().Context(), filter)

	s.Require().Error(err, repoError)
	s.Require().Nil(parts)
}
