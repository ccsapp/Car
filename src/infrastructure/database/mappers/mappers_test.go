package mappers

import (
	"DCar/infrastructure/database/entities"
	carTypes "github.com/ccsapp/cargotypes"
	openapiTypes "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

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

var exampleDatabaseCar = entities.Car{
	Vin:            "12345678901234567",
	Brand:          "Volkswagen",
	Model:          "Golf",
	ProductionDate: time.Date(2022, 12, 01, 0, 0, 0, 0, time.UTC),
	Color:          "black",
	Consumption: entities.Consumption{
		City:     6.4,
		Combined: 5.2,
		Overland: 4.6,
	},
	Emissions: entities.Emissions{
		City:     120,
		Combined: 100,
		Overland: 90,
	},
	Engine: entities.Engine{
		Power: 110,
		Type:  "someType",
	},
	Fuel:          entities.ELECTRIC,
	FuelCapacity:  "54.0L;85.2kWh",
	NumberOfDoors: 5,
	NumberOfSeats: 5,
	Tire: entities.Tire{
		Manufacturer: "GOODYEAR",
		Type:         "185/65R15",
	},
	Transmission:   entities.MANUAL,
	TrunkVolume:    435,
	Weight:         1320,
	TrunkLockState: entities.LOCKED,
}

func TestMapCarToDb(t *testing.T) {
	assert.Equal(t, exampleDatabaseCar, MapCarToDb(&exampleModelCar))
}

func TestMapCarFromDb(t *testing.T) {
	assert.Equal(t, exampleModelCar, MapCarFromDb(&exampleDatabaseCar))
}
