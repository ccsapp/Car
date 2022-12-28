package database

import (
	"DCar/infrastructure/database/db"
	"DCar/infrastructure/database/entities"
	"DCar/infrastructure/database/mappers"
	"DCar/logic/model"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const CarsCollectionBaseName = "cars"

var conversionError = errors.New("invalid type in database")

// IsDuplicateKeyError checks if the error is a duplicate key error. A duplicate key error can occur if you try to
// insert a car with a duplicate VIN.
func IsDuplicateKeyError(err error) bool {
	return mongo.IsDuplicateKeyError(err)
}

// IsNotFoundError checks if the error is a not found error. A not found error can occur if you try to read a
// car with a VIN that does not exist.
func IsNotFoundError(err error) bool {
	return err == mongo.ErrNoDocuments
}

// ICRUD is a high level database interface. It directly maps to the business logic and abstracts away the
// database entities and the database connection.
type ICRUD interface {
	// CreateCar creates a new car in the database and returns the VIN. If the VIN already exists, an error is returned.
	// You can check if the error is such an error with IsDuplicateKeyError. Any other errors are unexpected.
	CreateCar(ctx context.Context, car *model.Car) (model.Vin, error)

	// ReadAllVins returns all VINs in the database. If there are no cars in the database, an empty slice is returned.
	// Any errors are unexpected.
	ReadAllVins(ctx context.Context) ([]model.Vin, error)

	// DeleteCar deletes the car with the given VIN and returns true. If the car does not exist, false is returned.
	// Any errors are unexpected.
	DeleteCar(ctx context.Context, vin model.Vin) (bool, error)

	// ReadCar returns the car with the given VIN. If the car does not exist, an error is returned. You can check if the
	// error is such an error with IsNotFoundError. Any other errors are unexpected.
	ReadCar(ctx context.Context, vin model.Vin) (model.Car, error)
}

type crud struct {
	db         db.IConnection
	collection string
}

func NewICRUD(db db.IConnection, config *db.Config) ICRUD {
	return &crud{
		db:         db,
		collection: config.CollectionPrefix + CarsCollectionBaseName,
	}
}

func (c *crud) CreateCar(ctx context.Context, car *model.Car) (model.Vin, error) {
	res, err := c.db.Insert(ctx, c.collection, mappers.MapCarToDb(car))
	if err != nil {
		return "", err
	}
	vin, ok := res.InsertedID.(model.Vin)
	if !ok {
		return "", conversionError
	}
	return vin, nil
}

func (c *crud) ReadAllVins(ctx context.Context) ([]model.Vin, error) {
	var ids []bson.M
	if err := c.db.GetIDs(ctx, c.collection, &ids); err != nil {
		return nil, err
	}
	vins := make([]model.Vin, len(ids))
	for i, id := range ids {
		vins[i] = id["_id"].(model.Vin)
	}
	return vins, nil
}

func (c *crud) DeleteCar(ctx context.Context, vin model.Vin) (bool, error) {
	res, err := c.db.DeleteOne(ctx, c.collection, bson.D{{"_id", vin}})
	if err != nil {
		return false, err
	}
	return res.DeletedCount > 0, err
}

func (c *crud) ReadCar(ctx context.Context, vin model.Vin) (model.Car, error) {
	res := c.db.FindOne(ctx, c.collection, bson.D{{"_id", vin}})
	var car entities.Car
	err := res.Decode(&car)
	if err != nil {
		return model.Car{}, err
	}
	return mappers.MapCarFromDb(&car), nil
}
