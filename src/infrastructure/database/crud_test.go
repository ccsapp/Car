package database

import (
	"DCar/infrastructure/database/mappers"
	"DCar/mocks"
	"context"
	"errors"
	"testing"
	"time"

	carTypes "git.scc.kit.edu/cm-tm/cm-team/projectwork/pse/domain/d-cargotypes.git"
	openapiTypes "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionName = "cars"

type TestCrudConfig struct{}

func (c *TestCrudConfig) GetAppCollectionPrefix() string {
	return ""
}

var config = &TestCrudConfig{}

var exampleModelCar = carTypes.Car{
	Brand: "Volkswagen",
	DynamicData: carTypes.DynamicData{
		DoorsLockState:      carTypes.UNLOCKED,
		EngineState:         carTypes.OFF,
		FuelLevelPercentage: 23,
		Position: carTypes.DynamicDataPosition{
			Latitude:  49.0069,
			Longitude: 8.4037,
		},
		TrunkLockState: carTypes.LOCKED,
	},
	Model: "Golf",
	ProductionDate: openapiTypes.Date{
		Time: time.Date(2022, 12, 01, 0, 0, 0, 0, time.UTC),
	},
	TechnicalSpecification: carTypes.TechnicalSpecification{
		Color: "black",
		Consumption: carTypes.TechnicalSpecificationConsumption{
			City:     6.4,
			Combined: 5.2,
			Overland: 4.6,
		},
		Emissions: carTypes.TechnicalSpecificationEmissions{
			City:     120,
			Combined: 100,
			Overland: 90,
		},
		Engine: carTypes.TechnicalSpecificationEngine{

			Power: 110,
			Type:  "someType",
		},
		Fuel:          carTypes.ELECTRIC,
		FuelCapacity:  "54.0L;85.2kWh",
		NumberOfDoors: 5,
		NumberOfSeats: 5,
		Tire: carTypes.TechnicalSpecificationTire{
			Manufacturer: "GOODYEAR",
			Type:         "185/65R15",
		},
		Transmission: carTypes.MANUAL,
		TrunkVolume:  435,
		Weight:       1320,
	},
	Vin: "12345678901234567",
}

func TestCrud_CreateCar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockConnection := mocks.NewMockIConnection(ctrl)
	mockConnection.
		EXPECT().
		Insert(ctx, collectionName, mappers.MapCarToDb(&exampleModelCar)).
		Return(&mongo.InsertOneResult{
			InsertedID: exampleModelCar.Vin,
		}, nil)

	crud := NewICRUD(mockConnection, config)
	vin, err := crud.CreateCar(ctx, &exampleModelCar)

	assert.Nil(t, err)
	assert.Equal(t, exampleModelCar.Vin, vin)
}

func TestCrud_CreateCar_dbError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockConnection := mocks.NewMockIConnection(ctrl)
	dbError := errors.New("db error")
	mockConnection.
		EXPECT().
		Insert(ctx, collectionName, mappers.MapCarToDb(&exampleModelCar)).
		Return(nil, dbError)

	crud := NewICRUD(mockConnection, config)
	vin, err := crud.CreateCar(ctx, &exampleModelCar)

	assert.ErrorIs(t, err, dbError)
	assert.Equal(t, "", vin)
}

func TestCrud_CreateCar_conversionError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockConnection := mocks.NewMockIConnection(ctrl)

	invalidReturn := mongo.InsertOneResult{
		InsertedID: nil, // nil cannot be converted to string
	}

	mockConnection.
		EXPECT().
		Insert(ctx, collectionName, gomock.Any()).
		Return(&invalidReturn, nil)

	crud := NewICRUD(mockConnection, config)
	vin, err := crud.CreateCar(ctx, &exampleModelCar)

	assert.Equal(t, err, conversionError)
	assert.Equal(t, "", vin)
}

func TestCrud_ReadAllVins_NoErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	returnArray := []bson.M{
		{"_id": "JH4DA1840KS004941"},
		{"_id": "WV2YB0257EH008533"},
	}

	expectedVins := []carTypes.Vin{
		"JH4DA1840KS004941",
		"WV2YB0257EH008533",
	}

	mockConnection := mocks.NewMockIConnection(ctrl)

	mockGetIds := func(ctx context.Context, collection string, resultIds *[]bson.M) error {
		// increase slice length by 2 to provide space for the result
		*resultIds = append(*resultIds, nil, nil)
		copy(*resultIds, returnArray[:])
		return nil
	}

	mockConnection.
		EXPECT().
		GetIDs(ctx, collectionName, gomock.Any()).
		DoAndReturn(mockGetIds)

	crud := NewICRUD(mockConnection, config)
	vins, err := crud.ReadAllVins(ctx)

	assert.Nil(t, err)
	assert.Equal(t, expectedVins, vins)
}

func TestCrud_ReadAllVins_dbError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	dbError := errors.New("db error")

	mockConnection := mocks.NewMockIConnection(ctrl)

	mockConnection.
		EXPECT().
		GetIDs(ctx, collectionName, gomock.Any()).
		Return(dbError)

	crud := NewICRUD(mockConnection, config)
	vins, err := crud.ReadAllVins(ctx)

	assert.ErrorIs(t, err, dbError)
	assert.Nil(t, vins)
}

func TestCrud_DeleteCarSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockConnection := mocks.NewMockIConnection(ctrl)
	mockConnection.
		EXPECT().
		DeleteOne(ctx, collectionName, bson.D{{"_id", "12345678901234567"}}).
		Return(&mongo.DeleteResult{
			DeletedCount: 1,
		}, nil)

	crud := NewICRUD(mockConnection, config)
	success, err := crud.DeleteCar(ctx, "12345678901234567")

	assert.Nil(t, err)
	assert.True(t, success)
}

func TestCrud_DeleteCar_dbError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	dbError := errors.New("db error")

	mockConnection := mocks.NewMockIConnection(ctrl)
	mockConnection.
		EXPECT().
		DeleteOne(ctx, collectionName, bson.D{{"_id", "12345678901234567"}}).
		Return(nil, dbError)

	crud := NewICRUD(mockConnection, config)
	success, err := crud.DeleteCar(ctx, "12345678901234567")

	assert.ErrorIs(t, err, dbError)
	assert.False(t, success)
}

func TestCrud_DeleteCar_notFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockConnection := mocks.NewMockIConnection(ctrl)
	mockConnection.
		EXPECT().
		DeleteOne(ctx, collectionName, bson.D{{"_id", "12345678901234567"}}).
		Return(&mongo.DeleteResult{
			DeletedCount: 0,
		}, nil)

	crud := NewICRUD(mockConnection, config)
	success, err := crud.DeleteCar(ctx, "12345678901234567")

	assert.Nil(t, err)
	assert.False(t, success)
}

func TestCrud_ReadCar_success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockConnection := mocks.NewMockIConnection(ctrl)
	mockConnection.
		EXPECT().
		FindOne(ctx, collectionName, bson.D{{"_id", "12345678901234567"}}).
		Return(mongo.NewSingleResultFromDocument(mappers.MapCarToDb(&exampleModelCar), nil, nil))

	crud := NewICRUD(mockConnection, config)
	car, err := crud.ReadCar(ctx, "12345678901234567")

	assert.Nil(t, err)
	assert.Equal(t, exampleModelCar, car)
}

func TestCrud_ReadCar_decodeError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockConnection := mocks.NewMockIConnection(ctrl)
	mockConnection.
		EXPECT().
		FindOne(ctx, collectionName, bson.D{{"_id", "12345678901234567"}}).
		Return(mongo.NewSingleResultFromDocument(nil, errors.New("error"), nil))

	crud := NewICRUD(mockConnection, config)
	car, err := crud.ReadCar(ctx, "12345678901234567")

	assert.NotNil(t, err)
	assert.Equal(t, carTypes.Car{}, car)
}

func TestCrud_SetTrunkLockState_successChanged(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockConnection := mocks.NewMockIConnection(ctrl)
	mockConnection.
		EXPECT().
		UpdateOne(ctx, collectionName, bson.D{{"_id", "12345678901234567"}},
			bson.D{{"mockData_trunkLockState", carTypes.UNLOCKED}}).
		Return(&mongo.UpdateResult{
			MatchedCount: 1,
		}, nil)

	crud := NewICRUD(mockConnection, config)
	err := crud.SetTrunkLockState(ctx, "12345678901234567", carTypes.UNLOCKED)

	assert.Nil(t, err)
}

func TestCrud_SetTrunkLockState_errorCarNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockConnection := mocks.NewMockIConnection(ctrl)
	mockConnection.
		EXPECT().
		UpdateOne(ctx, collectionName, bson.D{{"_id", "12345678901234567"}},
			bson.D{{"mockData_trunkLockState", carTypes.UNLOCKED}}).
		Return(&mongo.UpdateResult{
			MatchedCount: 0,
		}, nil)

	crud := NewICRUD(mockConnection, config)
	err := crud.SetTrunkLockState(ctx, "12345678901234567", carTypes.UNLOCKED)

	assert.True(t, IsNotFoundError(err))
}

func TestCrud_SetTrunkLockState_databaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	databaseError := errors.New("database error")

	mockConnection := mocks.NewMockIConnection(ctrl)
	mockConnection.
		EXPECT().
		UpdateOne(ctx, collectionName, bson.D{{"_id", "12345678901234567"}},
			bson.D{{"mockData_trunkLockState", carTypes.UNLOCKED}}).
		Return(nil, databaseError)

	crud := NewICRUD(mockConnection, config)
	err := crud.SetTrunkLockState(ctx, "12345678901234567", carTypes.UNLOCKED)

	assert.ErrorIs(t, err, databaseError)
}
