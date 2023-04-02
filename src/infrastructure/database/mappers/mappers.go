package mappers

import (
	"DCar/infrastructure/database/entities"
	"DCar/logic/model"
	carTypes "github.com/ccsapp/cargotypes"
	openapiTypes "github.com/deepmap/oapi-codegen/pkg/types"
)

// MapCarToDb maps a car from the domain to a car in the database.
// The trunk lock state is not mapped and set to LOCKED.
func MapCarToDb(car *carTypes.Car) entities.Car {
	return entities.Car{
		Vin:            car.Vin,
		Brand:          car.Brand,
		Model:          car.Model,
		ProductionDate: car.ProductionDate.Time,
		Color:          car.TechnicalSpecification.Color,
		Consumption: entities.Consumption{
			City:     car.TechnicalSpecification.Consumption.City,
			Combined: car.TechnicalSpecification.Consumption.Combined,
			Overland: car.TechnicalSpecification.Consumption.Overland,
		},
		Emissions: entities.Emissions{
			City:     car.TechnicalSpecification.Emissions.City,
			Combined: car.TechnicalSpecification.Emissions.Combined,
			Overland: car.TechnicalSpecification.Emissions.Overland,
		},
		Engine: entities.Engine{
			Power: car.TechnicalSpecification.Engine.Power,
			Type:  car.TechnicalSpecification.Engine.Type,
		},
		Fuel:          entities.Fuel(car.TechnicalSpecification.Fuel),
		FuelCapacity:  car.TechnicalSpecification.FuelCapacity,
		NumberOfDoors: car.TechnicalSpecification.NumberOfDoors,
		NumberOfSeats: car.TechnicalSpecification.NumberOfSeats,
		Tire: entities.Tire{
			Manufacturer: car.TechnicalSpecification.Tire.Manufacturer,
			Type:         car.TechnicalSpecification.Tire.Type,
		},
		Transmission:   entities.Transmission(car.TechnicalSpecification.Transmission),
		TrunkVolume:    car.TechnicalSpecification.TrunkVolume,
		Weight:         car.TechnicalSpecification.Weight,
		TrunkLockState: entities.LOCKED,
	}
}

func MapCarFromDb(car *entities.Car) carTypes.Car {
	return carTypes.Car{
		Vin:                    car.Vin,
		Brand:                  car.Brand,
		DynamicData:            model.ExampleDynamicData(car.TrunkLockState), // TODO: do not use example data
		Model:                  car.Model,
		ProductionDate:         openapiTypes.Date{Time: car.ProductionDate},
		TechnicalSpecification: mapTechnicalSpecificationFromDb(car),
	}
}

func mapTechnicalSpecificationFromDb(car *entities.Car) carTypes.TechnicalSpecification {
	return carTypes.TechnicalSpecification{
		Color:         car.Color,
		Consumption:   mapTechnicalSpecificationConsumptionFromDb(&car.Consumption),
		Emissions:     mapTechnicalSpecificationEmissionsFromDb(&car.Emissions),
		Engine:        mapTechnicalSpecificationEngineFromDb(&car.Engine),
		Fuel:          carTypes.TechnicalSpecificationFuel(car.Fuel),
		FuelCapacity:  car.FuelCapacity,
		NumberOfDoors: car.NumberOfDoors,
		NumberOfSeats: car.NumberOfSeats,
		Tire:          mapTechnicalSpecificationTireFromDb(&car.Tire),
		Transmission:  carTypes.TechnicalSpecificationTransmission(car.Transmission),
		TrunkVolume:   car.TrunkVolume,
		Weight:        car.Weight,
	}
}

func mapTechnicalSpecificationConsumptionFromDb(consumption *entities.Consumption) carTypes.TechnicalSpecificationConsumption {
	return carTypes.TechnicalSpecificationConsumption{
		City:     consumption.City,
		Combined: consumption.Combined,
		Overland: consumption.Overland,
	}
}

func mapTechnicalSpecificationEmissionsFromDb(emissions *entities.Emissions) carTypes.TechnicalSpecificationEmissions {
	return carTypes.TechnicalSpecificationEmissions{
		City:     emissions.City,
		Combined: emissions.Combined,
		Overland: emissions.Overland,
	}
}

func mapTechnicalSpecificationEngineFromDb(engine *entities.Engine) carTypes.TechnicalSpecificationEngine {
	return carTypes.TechnicalSpecificationEngine{
		Power: engine.Power,
		Type:  engine.Type,
	}
}

func mapTechnicalSpecificationTireFromDb(tire *entities.Tire) carTypes.TechnicalSpecificationTire {
	return carTypes.TechnicalSpecificationTire{
		Manufacturer: tire.Manufacturer,
		Type:         tire.Type,
	}
}
