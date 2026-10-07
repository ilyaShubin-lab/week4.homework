package converter

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"boilerplates/inventory/internal/model"
	inventoryV1 "boilerplates/shared/pkg/proto/inventory/v1"
)

func PartToProto(part model.Part) *inventoryV1.Part {
	return &inventoryV1.Part{
		Uuid:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      categoryToProto(part.Category),
		Dimensions:    dimensionsToProto(part.Dimensions),
		Manufacturer:  manufacturerToProto(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      metadataToProto(part.Metadata),
		CreatedAt:     timestamppb.New(part.CreatedAt),
		UpdatedAt:     timeToProto(part.UpdatedAt),
	}
}

func categoryToProto(c model.Category) inventoryV1.Category {
	switch c {
	case model.CategoryEngine:
		return inventoryV1.Category_CATEGORY_ENGINE
	case model.CategoryFuel:
		return inventoryV1.Category_CATEGORY_FUEL
	case model.CategoryPorthole:
		return inventoryV1.Category_CATEGORY_PORTHOLE
	case model.CategoryWing:
		return inventoryV1.Category_CATEGORY_WING
	default:
		return inventoryV1.Category_CATEGORY_UNSPECIFIED
	}
}

func timeToProto(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

func dimensionsToProto(data *model.Dimensions) *inventoryV1.Dimensions {
	if data == nil {
		return nil
	}
	return &inventoryV1.Dimensions{
		Length: data.Length,
		Width:  data.Width,
		Height: data.Height,
		Weight: data.Weight,
	}
}

func manufacturerToProto(data *model.Manufacturer) *inventoryV1.Manufacturer {
	if data == nil {
		return nil
	}

	return &inventoryV1.Manufacturer{
		Name:    data.Name,
		Country: data.Country,
		Website: data.Website,
	}
}

func metadataToProto(m map[string]model.Value) map[string]*inventoryV1.Value {
	if m == nil {
		return nil
	}

	res := make(map[string]*inventoryV1.Value, len(m))

	for k, v := range m {
		res[k] = valueToProto(v)
	}

	return res
}

func valueToProto(v model.Value) *inventoryV1.Value {
	switch {
	case v.StringValue != nil:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_StringValue{StringValue: *v.StringValue},
		}
	case v.Int64Value != nil:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_Int64Value{Int64Value: *v.Int64Value},
		}
	case v.DoubleValue != nil:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_DoubleValue{DoubleValue: *v.DoubleValue},
		}
	case v.BoolValue != nil:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_BoolValue{BoolValue: *v.BoolValue},
		}
	default:
		return nil
	}
}

func PartsListToProto(parts []model.Part) []*inventoryV1.Part {
	result := make([]*inventoryV1.Part, 0, len(parts))

	for _, val := range parts {
		result = append(result, PartToProto(val))
	}

	return result
}

func PartsFilterToModel(filter *inventoryV1.PartsFilter) model.PartsFilter {
	if filter == nil {
		return model.PartsFilter{}
	}

	result := model.PartsFilter{
		UUIDs:                 filter.GetUuids(),
		Names:                 filter.GetNames(),
		Categories:            categoriesToModel(filter.GetCategories()),
		ManufacturerCountries: filter.GetManufacturerCountries(),
		Tags:                  filter.GetTags(),
	}
	return result
}

func categoriesToModel(categories []inventoryV1.Category) []model.Category {
	if categories == nil {
		return nil
	}

	result := make([]model.Category, 0, len(categories))
	for _, c := range categories {
		result = append(result, categoryToModel(c))
	}
	return result
}

func categoryToModel(c inventoryV1.Category) model.Category {
	switch c {
	case inventoryV1.Category_CATEGORY_ENGINE:
		return model.CategoryEngine
	case inventoryV1.Category_CATEGORY_FUEL:
		return model.CategoryFuel
	case inventoryV1.Category_CATEGORY_PORTHOLE:
		return model.CategoryPorthole
	case inventoryV1.Category_CATEGORY_WING:
		return model.CategoryWing
	default:
		return model.CategoryUnknown
	}
}
