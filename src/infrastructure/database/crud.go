package database

//go:generate mockgen -source=./crud.go -package=mocks -destination=../../mocks/mock_crud.go

import (
	"DCar/infrastructure/database/db"
	"DCar/infrastructure/database/entities"
	"DCar/infrastructure/database/mappers"
	"context"
	"errors"
	carTypes "github.com/ccsapp/cargotypes"
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

type CrudConfig interface {
	GetAppCollectionPrefix() string
}

// ICRUD is a high level database interface. It directly maps to the business logic and abstracts away the
// database entities and the database connection.
type ICRUD interface {
	// CreateCar creates a new car in the database and returns the VIN. If the VIN already exists, an error is returned.
	// You can check if the error is such an error with IsDuplicateKeyError. Any other errors are unexpected.
	CreateCar(ctx context.Context, car *carTypes.Car) (carTypes.Vin, error)

	// ReadAllVins returns all VINs in the database. If there are no cars in the database, an empty slice is returned.
	// Any errors are unexpected.
	ReadAllVins(ctx context.Context) ([]carTypes.Vin, error)

	// DeleteCar deletes the car with the given VIN and returns true. If the car does not exist, false is returned.
	// Any errors are unexpected.
	DeleteCar(ctx context.Context, vin carTypes.Vin) (bool, error)

	// ReadCar returns the car with the given VIN. If the car does not exist, an error is returned. You can check if the
	// error is such an error with IsNotFoundError. Any other errors are unexpected.
	ReadCar(ctx context.Context, vin carTypes.Vin) (carTypes.Car, error)

	// SetTrunkLockState sets the trunk lock state of the car with the given VIN. If the car does not exist,
	// an error is returned. You can check if the error is such an error with IsNotFoundError. Any other errors are
	// unexpected.
	SetTrunkLockState(ctx context.Context, vin carTypes.Vin, state carTypes.DynamicDataLockState) error
}

type crud struct {
	db         db.IConnection
	collection string
}

func NewICRUD(db db.IConnection, config CrudConfig) ICRUD {
	return &crud{
		db:         db,
		collection: config.GetAppCollectionPrefix() + CarsCollectionBaseName,
	}
}

func (c *crud) CreateCar(ctx context.Context, car *carTypes.Car) (carTypes.Vin, error) {
	res, err := c.db.Insert(ctx, c.collection, mappers.MapCarToDb(car))
	if err != nil {
		return "", err
	}
	vin, ok := res.InsertedID.(carTypes.Vin)
	if !ok {
		return "", conversionError
	}
	return vin, nil
}

func (c *crud) ReadAllVins(ctx context.Context) ([]carTypes.Vin, error) {
	var ids []bson.M
	if err := c.db.GetIDs(ctx, c.collection, &ids); err != nil {
		return nil, err
	}
	vins := make([]carTypes.Vin, len(ids))
	for i, id := range ids {
		vins[i] = id["_id"].(carTypes.Vin)
	}
	return vins, nil
}

func (c *crud) DeleteCar(ctx context.Context, vin carTypes.Vin) (bool, error) {
	res, err := c.db.DeleteOne(ctx, c.collection, bson.D{{"_id", vin}})
	if err != nil {
		return false, err
	}
	return res.DeletedCount > 0, nil
}

func (c *crud) ReadCar(ctx context.Context, vin carTypes.Vin) (carTypes.Car, error) {
	res := c.db.FindOne(ctx, c.collection, bson.D{{"_id", vin}})
	var car entities.Car
	err := res.Decode(&car)
	if err != nil {
		return carTypes.Car{}, err
	}
	return mappers.MapCarFromDb(&car), nil
}

func (c *crud) SetTrunkLockState(ctx context.Context, vin carTypes.Vin, state carTypes.DynamicDataLockState) error {
	res, err := c.db.UpdateOne(ctx, c.collection, bson.D{{"_id", vin}}, bson.D{{"mockData_trunkLockState",
		state}})

	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
