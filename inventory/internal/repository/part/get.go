package part

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"boilerplates/inventory/internal/model"
	"boilerplates/inventory/internal/repository/converter"
	repoModel "boilerplates/inventory/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, uuid string) (model.Part, error) {
	var part repoModel.Part
	filter := bson.M{"_id": uuid}

	err := r.collection.FindOne(ctx, filter).Decode(&part)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Part{}, model.ErrPartNotFound
		}
		return model.Part{}, fmt.Errorf("find part %s: %w", uuid, err)
	}

	return converter.PartToModel(part), nil
}
