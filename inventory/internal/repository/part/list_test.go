package part

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"

	"boilerplates/inventory/internal/model"
)

func TestBuildFilter(t *testing.T) {
	tests := []struct {
		name   string
		filter model.PartsFilter
		want   bson.M
	}{
		{
			name:   "пустой фильтр — все документы",
			filter: model.PartsFilter{},
			want:   bson.M{},
		},
		{
			filter: model.PartsFilter{UUIDs: []string{}, Tags: []string{}},
			want:   bson.M{},
		},
		{
			name:   "uuids фильтруются по _id",
			filter: model.PartsFilter{UUIDs: []string{"uuid-1", "uuid-2"}},
			want:   bson.M{"_id": bson.M{"$in": []string{"uuid-1", "uuid-2"}}},
		},
		{
			name:   "имена",
			filter: model.PartsFilter{Names: []string{"Топливный бак ТБ-200"}},
			want:   bson.M{"name": bson.M{"$in": []string{"Топливный бак ТБ-200"}}},
		},
		{
			name:   "категории",
			filter: model.PartsFilter{Categories: []model.Category{model.CategoryEngine, model.CategoryPorthole}},
			want:   bson.M{"category": bson.M{"$in": []model.Category{model.CategoryEngine, model.CategoryPorthole}}},
		},
		{
			name:   "страна производителя через manufacturer.country",
			filter: model.PartsFilter{ManufacturerCountries: []string{"Россия"}},
			want:   bson.M{"manufacturer.country": bson.M{"$in": []string{"Россия"}}},
		},
		{
			name:   "теги",
			filter: model.PartsFilter{Tags: []string{"бак", "обзор"}},
			want:   bson.M{"tags": bson.M{"$in": []string{"бак", "обзор"}}},
		},
		{
			name: "несколько условий объединяются",
			filter: model.PartsFilter{
				Names: []string{"Иллюминатор бронированный"},
				Tags:  []string{"обзор"},
			},
			want: bson.M{
				"name": bson.M{"$in": []string{"Иллюминатор бронированный"}},
				"tags": bson.M{"$in": []string{"обзор"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildFilter(tt.filter)

			require.Equal(t, tt.want, got)
		})
	}
}
