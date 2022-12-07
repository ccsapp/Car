package mappers

import (
	"DCar/infrastructure/database/entities"
	"DCar/logic/model"
	openapiTypes "github.com/deepmap/oapi-codegen/pkg/types"
)

func MapCarToDb(car *model.Car) entities.Car {
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
		Transmission: entities.Transmission(car.TechnicalSpecification.Transmission),
		TrunkVolume:  car.TechnicalSpecification.TrunkVolume,
		Weight:       car.TechnicalSpecification.Weight,
	}
}

func MapCarFromDb(car *entities.Car) model.Car {
	return model.Car{
		Vin:                    car.Vin,
		Brand:                  car.Brand,
		DynamicData:            model.ExampleDynamicData(), // TODO: do not use example data
		Model:                  car.Model,
		ProductionDate:         openapiTypes.Date{Time: car.ProductionDate},
		TechnicalSpecification: mapTechnicalSpecificationFromDb(car),
	}
}

func mapTechnicalSpecificationFromDb(car *entities.Car) model.TechnicalSpecification {
	return model.TechnicalSpecification{
		Color:         car.Color,
		Consumption:   mapTechnicalSpecificationConsumptionFromDb(&car.Consumption),
		Emissions:     mapTechnicalSpecificationEmissionsFromDb(&car.Emissions),
		Engine:        mapTechnicalSpecificationEngineFromDb(&car.Engine),
		Fuel:          model.TechnicalSpecificationFuel(car.Fuel),
		FuelCapacity:  car.FuelCapacity,
		NumberOfDoors: car.NumberOfDoors,
		NumberOfSeats: car.NumberOfSeats,
		Tire:          mapTechnicalSpecificationTireFromDb(&car.Tire),
		Transmission:  model.TechnicalSpecificationTransmission(car.Transmission),
		TrunkVolume:   car.TrunkVolume,
		Weight:        car.Weight,
	}
}

func mapTechnicalSpecificationConsumptionFromDb(consumption *entities.Consumption) model.TechnicalSpecificationConsumption {
	return model.TechnicalSpecificationConsumption{
		City:     consumption.City,
		Combined: consumption.Combined,
		Overland: consumption.Overland,
	}
}

func mapTechnicalSpecificationEmissionsFromDb(emissions *entities.Emissions) model.TechnicalSpecificationEmissions {
	return model.TechnicalSpecificationEmissions{
		City:     emissions.City,
		Combined: emissions.Combined,
		Overland: emissions.Overland,
	}
}

func mapTechnicalSpecificationEngineFromDb(engine *entities.Engine) model.TechnicalSpecificationEngine {
	return model.TechnicalSpecificationEngine{
		Power: engine.Power,
		Type:  engine.Type,
	}
}

func mapTechnicalSpecificationTireFromDb(tire *entities.Tire) model.TechnicalSpecificationTire {
	return model.TechnicalSpecificationTire{
		Manufacturer: tire.Manufacturer,
		Type:         tire.Type,
	}
}
