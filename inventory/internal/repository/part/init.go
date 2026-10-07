package part

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	repoModel "boilerplates/inventory/internal/repository/model"
)

func ptr[T any](v T) *T { return &v }

func (r *repository) InitParts(ctx context.Context) error {
	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("count parts: %w", err)
	}

	if count > 0 {
		return nil
	}

	now := time.Now()
	parts := []repoModel.Part{
		{
			UUID:          "550e8400-e29b-41d4-a716-446655440001",
			Name:          "Двигатель импульсный ИД-500",
			Description:   "Основной маршевый двигатель",
			Price:         1500000.50,
			StockQuantity: 3,
			Category:      1, // ENGINE
			Dimensions: &repoModel.Dimensions{
				Length: 400, Width: 200, Height: 200, Weight: 5000,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "Космофлот",
				Country: "Россия",
				Website: "https://kosmoflot.ru",
			},
			Tags: []string{"двигатель", "маршевый"},
			Metadata: map[string]repoModel.Value{
				"warranty_years": {Int64Value: ptr(int64(5))},
				"certified":      {BoolValue: ptr(true)},
			},
			CreatedAt: now,
			UpdatedAt: ptr(now),
		},
		{
			UUID:          "550e8400-e29b-41d4-a716-446655440002",
			Name:          "Топливный бак ТБ-200",
			Description:   "Основной топливный бак",
			Price:         250000,
			StockQuantity: 10,
			Category:      2, // FUEL
			Tags:          []string{"топливо", "бак"},
			CreatedAt:     now,
			UpdatedAt:     ptr(now),
		},
		{
			UUID:          "550e8400-e29b-41d4-a716-446655440003",
			Name:          "Иллюминатор бронированный",
			Description:   "Обзорный иллюминатор",
			Price:         45000,
			StockQuantity: 25,
			Category:      3, // PORTHOLE
			Tags:          []string{"обзор"},
			CreatedAt:     now,
			UpdatedAt:     ptr(now),
		},
	}

	docs := make([]any, 0, len(parts))
	for _, p := range parts {
		docs = append(docs, p)
	}
	_, err = r.collection.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("insert parts: %w", err)
	}
	return nil
}
