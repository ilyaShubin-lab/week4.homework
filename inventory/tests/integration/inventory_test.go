//go:build integration

package integration

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	inventoryV1 "boilerplates/shared/pkg/proto/inventory/v1"
)

func (s *InventorySuite) TestHealthCheck() {
	resp, err := grpc_health_v1.NewHealthClient(s.conn).
		Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})

	s.Require().NoError(err)
	s.Equal(grpc_health_v1.HealthCheckResponse_SERVING, resp.GetStatus())
}

func (s *InventorySuite) TestListPartsAll() {
	resp, err := s.client.ListParts(context.Background(), &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{},
	})
	s.Require().NoError(err)
	s.Require().Len(resp.GetParts(), seededPartsCount)

	uuids := make([]string, 0, len(resp.GetParts()))
	for _, p := range resp.GetParts() {
		uuids = append(uuids, p.GetUuid())
	}

	s.ElementsMatch([]string{engineUUID, fuelTankUUID, portholeUUID}, uuids)
}

func (s *InventorySuite) TestListPartsByCategory() {
	resp, err := s.client.ListParts(context.Background(), &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			Categories: []inventoryV1.Category{inventoryV1.Category_CATEGORY_ENGINE},
		},
	})
	s.Require().NoError(err)
	s.Require().Len(resp.GetParts(), 1)
	s.Equal(engineUUID, resp.GetParts()[0].GetUuid())
}

func (s *InventorySuite) TestGetPart() {
	resp, err := s.client.GetPart(context.Background(), &inventoryV1.GetPartRequest{Uuid: engineUUID})
	s.Require().NoError(err)

	part := resp.GetPart()
	s.Equal(engineName, part.GetName())
	s.Equal(inventoryV1.Category_CATEGORY_ENGINE, part.GetCategory())
	s.Equal("Россия", part.GetManufacturer().GetCountry())
	s.True(part.GetMetadata()["certified"].GetBoolValue())
}

func (s *InventorySuite) TestGetPartNotFound() {
	_, err := s.client.GetPart(context.Background(), &inventoryV1.GetPartRequest{
		Uuid: "00000000-0000-0000-0000-000000000000",
	})
	s.Require().Error(err)
	s.Equal(codes.NotFound, status.Code(err))
}

func (s *InventorySuite) TestGetPartEmptyUUID() {
	_, err := s.client.GetPart(context.Background(), &inventoryV1.GetPartRequest{Uuid: ""})
	s.Require().Error(err)
	s.Equal(codes.InvalidArgument, status.Code(err))
}
