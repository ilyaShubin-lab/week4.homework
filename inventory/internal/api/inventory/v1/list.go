package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"boilerplates/inventory/internal/converter"
	inventoryv1 "boilerplates/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(
	ctx context.Context,
	req *inventoryv1.ListPartsRequest,
) (*inventoryv1.ListPartsResponse, error) {
	parts, err := a.partService.List(ctx, converter.PartsFilterToModel(req.GetFilter()))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &inventoryv1.ListPartsResponse{
		Parts: converter.PartsListToProto(parts),
	}, nil
}
