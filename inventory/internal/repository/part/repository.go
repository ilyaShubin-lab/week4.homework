package part

// имя папки

import (
	"go.mongodb.org/mongo-driver/mongo"

	def "boilerplates/inventory/internal/repository"
)

// Компилятор проверит, что все методы интерфейса на месте
var _ def.PartRepository = (*repository)(nil)

const partsCollection = "parts"

type repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *repository {
	return &repository{
		collection: db.Collection(partsCollection),
	}
}
