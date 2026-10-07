package part

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"

	"boilerplates/inventory/internal/model"
	"boilerplates/inventory/internal/repository/converter"
	repoModel "boilerplates/inventory/internal/repository/model"
)

func (r *repository) List(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	query := buildFilter(filter)

	cursor, err := r.collection.Find(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("find parts: %w", err)
	}

	defer func() {
		if closeErr := cursor.Close(ctx); closeErr != nil {
			log.Printf("close mongo cursor: %v", closeErr)
		}
	}()

	var parts []repoModel.Part
	err = cursor.All(ctx, &parts)
	if err != nil {
		return nil, fmt.Errorf("decode parts: %w", err)
	}

	res := make([]model.Part, 0, len(parts))
	for _, p := range parts {
		res = append(res, converter.PartToModel(p))
	}

	return res, nil
}

func buildFilter(f model.PartsFilter) bson.M {
	query := bson.M{}

	if len(f.UUIDs) > 0 {
		query["_id"] = bson.M{"$in": f.UUIDs}
	}
	if len(f.Names) > 0 {
		query["name"] = bson.M{"$in": f.Names}
	}
	if len(f.Categories) > 0 {
		query["category"] = bson.M{"$in": f.Categories}
	}
	if len(f.ManufacturerCountries) > 0 {
		query["manufacturer.country"] = bson.M{"$in": f.ManufacturerCountries}
	}
	if len(f.Tags) > 0 {
		query["tags"] = bson.M{"$in": f.Tags}
	}
	return query
}
