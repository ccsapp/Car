// Package model provides primitives to interact with the openapi HTTP API.
package model

import (
	openapiTypes "github.com/deepmap/oapi-codegen/pkg/types"
)

// Defines values for DynamicDataEngineState.
const (
	OFF DynamicDataEngineState = "OFF"
	ON  DynamicDataEngineState = "ON"
)

// Defines values for TechnicalSpecificationFuel.
const (
	DIESEL       TechnicalSpecificationFuel = "DIESEL"
	ELECTRIC     TechnicalSpecificationFuel = "ELECTRIC"
	HYBRIDDIESEL TechnicalSpecificationFuel = "HYBRID_DIESEL"
	HYBRIDPETROL TechnicalSpecificationFuel = "HYBRID_PETROL"
	PETROL       TechnicalSpecificationFuel = "PETROL"
)

// Defines values for TechnicalSpecificationTransmission.
const (
	AUTOMATIC TechnicalSpecificationTransmission = "AUTOMATIC"
	MANUAL    TechnicalSpecificationTransmission = "MANUAL"
)

// Car A specific type of vehicle
type Car struct {
	// Brand Data that specifies the brand name of the Vehicle manufacturer
	Brand *string `json:"brand,omitempty"`

	// DynamicData Data that changes during a car's operation
	DynamicData *DynamicData `json:"dynamicData,omitempty"`

	// Model Data that specifies the particular type of a Vehicle
	Model *string `json:"model,omitempty"`

	// ProductionDate Data that specifies the official date the vehicle was declared to have exited production by the manufacturer.
	ProductionDate         *openapiTypes.Date      `json:"productionDate,omitempty"`
	TechnicalSpecification *TechnicalSpecification `json:"technicalSpecification,omitempty"`

	// Vin A Vehicle Identification Number (VIN) which uniquely identifies a car
	Vin *Vin `json:"vin,omitempty"`
}

// DynamicData Data that changes during a car's operation
type DynamicData struct {
	// DoorsLockState Data that specifies whether an object is locked or unlocked
	DoorsLockState *LockState              `json:"doorsLockState,omitempty"`
	EngineState    *DynamicDataEngineState `json:"engineState,omitempty"`

	// FuelLevelPercentage Data that specifies the relation of remaining fuelCapacity to the maximum fuelCapacity in percentage
	FuelLevelPercentage *int `json:"fuelLevelPercentage,omitempty"`

	// Position Data that specifies the GeoCoordinate of a car
	Position *struct {
		// Latitude Data that specifies the distance from the equator
		Latitude *float32 `json:"latitude,omitempty"`

		// Longitude Data that specifies the distance east or west from a line (meridian) passing through Greenwich
		Longitude *float32 `json:"longitude,omitempty"`
	} `json:"position,omitempty"`

	// TrunkLockState Data that specifies whether an object is locked or unlocked
	TrunkLockState *LockState `json:"trunkLockState,omitempty"`
}

// DynamicDataEngineState defines model for DynamicData.EngineState.
type DynamicDataEngineState string

// LockState Data that specifies whether an object is locked or unlocked
type LockState = interface{}

// TechnicalSpecification defines model for technicalSpecification.
type TechnicalSpecification struct {
	// Color Data on the description of the paint job of a car
	Color *string `json:"color,omitempty"`

	// Consumption Data that specifies the amount of fuel consumed during car operation in units per 100 kilometers
	Consumption *struct {
		// City Data that specifies the amount of fuel that is consumed when driving within the city in: kW/100km or l/100km
		City *float32 `json:"city,omitempty"`

		// Combined Data that specifies the combined amount of fuel that is consumed in: kW / 100 km or l / 100 km
		Combined *float32 `json:"combined,omitempty"`

		// Overland Data that specifies the amount of fuel that is consumed when driving outside of a city in: kW/100km or l/100km
		Overland *float32 `json:"overland,omitempty"`
	} `json:"consumption,omitempty"`

	// Emissions Data that specifies the CO2 emitted by a car during operation in gram per kilometer
	Emissions *struct {
		// City Data that specifies the amount of emissions when driving within the city in: g CO2 / km
		City *float32 `json:"city,omitempty"`

		// Combined Data that specifies the combined amount of emissions in: g CO2 / km. The combination is done by the manufacturer according to an industry-specific standard
		Combined *float32 `json:"combined,omitempty"`

		// Overland Data that specifies the amount of emissions when driving outside of a city in: g CO2 / km
		Overland *float32 `json:"overland,omitempty"`
	} `json:"emissions,omitempty"`

	// Engine A physical unit that converts fuel into movement
	Engine *struct {
		// Power Data on the power the engine can provide in kW
		Power *int `json:"power,omitempty"`

		// Type Data that contains the manufacturer-given type description of the engine
		Type *string `json:"type,omitempty"`
	} `json:"engine,omitempty"`

	// Fuel Data that defines the source of energy that powers the car
	Fuel *TechnicalSpecificationFuel `json:"fuel,omitempty"`

	// FuelCapacity Data that specifies the amount of fuel that can be carried with the car
	FuelCapacity *string `json:"fuelCapacity,omitempty"`

	// NumberOfDoors Data that defines the number of doors that are built into a car
	NumberOfDoors *int `json:"numberOfDoors,omitempty"`

	// NumberOfSeats Data that defines the number of seats that are built into a car
	NumberOfSeats *int `json:"numberOfSeats,omitempty"`

	// Tire A physical unit that serves as the point of contact between a car and the ground
	Tire *struct {
		// Manufacturer Data denoting the company responsible for the creation of a physical unit
		Manufacturer *string `json:"manufacturer,omitempty"`

		// Type Data that contains the manufacturer-given type description of the tire
		Type *string `json:"type,omitempty"`
	} `json:"tire,omitempty"`

	// Transmission A physical unit responsible for managing the conversion rate of the engine (can be automated or manually operated)
	Transmission *TechnicalSpecificationTransmission `json:"transmission,omitempty"`

	// TrunkVolume Data on the physical volume of the trunk in liters
	TrunkVolume *int `json:"trunkVolume,omitempty"`

	// Weight Data that specifies the total weight of a car when empty in kilograms (kg)
	Weight *int `json:"weight,omitempty"`
}

// TechnicalSpecificationFuel Data that defines the source of energy that powers the car
type TechnicalSpecificationFuel string

// TechnicalSpecificationTransmission A physical unit responsible for managing the conversion rate of the engine (can be automated or manually operated)
type TechnicalSpecificationTransmission string

// Vin A Vehicle Identification Number (VIN) which uniquely identifies a car
type Vin = string

// VinParam A Vehicle Identification Number (VIN) which uniquely identifies a car
type VinParam = Vin

// AddVehicleJSONRequestBody defines body for AddVehicle for application/json ContentType.
type AddVehicleJSONRequestBody = Car
