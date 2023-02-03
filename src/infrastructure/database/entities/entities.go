package entities

import (
	"time"
)

// Defines values for Fuel.
const (
	DIESEL       Fuel = "DIESEL"
	ELECTRIC     Fuel = "ELECTRIC"
	HYBRIDDIESEL Fuel = "HYBRID_DIESEL"
	HYBRIDPETROL Fuel = "HYBRID_PETROL"
	PETROL       Fuel = "PETROL"
)

// Defines values for Transmission.
const (
	AUTOMATIC Transmission = "AUTOMATIC"
	MANUAL    Transmission = "MANUAL"
)

// Defines values for LockState.
const (
	UNLOCKED LockState = "UNLOCKED"
	LOCKED   LockState = "LOCKED"
)

// Car A specific type of vehicle
type Car struct {
	// Vin A Vehicle Identification Number (VIN) which uniquely identifies a car
	Vin Vin `bson:"_id"`

	// Brand Data that specifies the brand name of the Vehicle manufacturer
	Brand string `bson:"brand"`

	// Model Data that specifies the particular type of Vehicle
	Model string `bson:"model"`

	// ProductionDate Data that specifies the official date the vehicle was declared to have exited production by the manufacturer.
	ProductionDate time.Time `bson:"productionDate"`

	// Color Data on the description of the paint job of a car
	Color string `bson:"technicalSpecification_color"`

	// Consumption Data that specifies the amount of fuel consumed during car operation in units per 100 kilometers
	Consumption Consumption `bson:"technicalSpecification_consumption"`

	// Emissions Data that specifies the CO2 emitted by a car during operation in gram per kilometer
	Emissions Emissions `bson:"technicalSpecification_emissions"`

	// Engine A physical unit that converts fuel into movement
	Engine Engine `bson:"technicalSpecification_engine"`

	// Fuel Data that defines the source of energy that powers the car
	Fuel Fuel `bson:"technicalSpecification_fuel"`

	// FuelCapacity Data that specifies the amount of fuel that can be carried with the car
	FuelCapacity string `bson:"technicalSpecification_fuelCapacity"`

	// NumberOfDoors Data that defines the number of doors that are built into a car
	NumberOfDoors int `bson:"technicalSpecification_numberOfDoors"`

	// NumberOfSeats Data that defines the number of seats that are built into a car
	NumberOfSeats int `bson:"technicalSpecification_numberOfSeats"`

	// Tire A physical unit that serves as the point of contact between a car and the ground
	Tire Tire `bson:"technicalSpecification_tire"`

	// Transmission A physical unit responsible for managing the conversion rate of the engine (can be automated or manually operated)
	Transmission Transmission `bson:"technicalSpecification_transmission"`

	// TrunkVolume Data on the physical volume of the trunk in liters
	TrunkVolume int `bson:"technicalSpecification_trunkVolume"`

	// Weight Data that specifies the total weight of a car when empty in kilograms (kg)
	Weight int `bson:"technicalSpecification_weight"`

	// TrunkLockState Indicates the state of the trunk lock - this is stored in the database to simulate
	// a real car.
	TrunkLockState LockState `bson:"mockData_trunkLockState"`
}

type Consumption struct {
	// City Data that specifies the amount of fuel that is consumed when driving within the city in: kW/100km or l/100km
	City float32 `bson:"city"`

	// Combined Data that specifies the combined amount of fuel that is consumed in: kW / 100 km or l / 100 km
	Combined float32 `bson:"combined"`

	// Overland Data that specifies the amount of fuel that is consumed when driving outside of a city in: kW/100km or l/100km
	Overland float32 `bson:"overland"`
}

type Emissions struct {
	// City Data that specifies the amount of emissions when driving within the city in: g CO2 / km
	City float32 `bson:"city"`

	// Combined Data that specifies the combined amount of emissions in: g CO2 / km. The combination is done by the manufacturer according to an industry-specific standard
	Combined float32 `bson:"combined"`

	// Overland Data that specifies the amount of emissions when driving outside of a city in: g CO2 / km
	Overland float32 `bson:"overland"`
}

type Engine struct {
	// Power Data on the power the engine can provide in kW
	Power int `bson:"power"`

	// Type Data that contains the manufacturer-given type description of the engine
	Type string `bson:"type"`
}

type Tire struct {
	// Manufacturer Data denoting the company responsible for the creation of a physical unit
	Manufacturer string `bson:"manufacturer"`

	// Type Data that contains the manufacturer-given type description of the tire
	Type string `bson:"type"`
}

// Fuel Data that defines the source of energy that powers the car
type Fuel string

// Transmission A physical unit responsible for managing the conversion rate of the engine (can be automated or manually operated)
type Transmission string

// Vin A Vehicle Identification Number (VIN) which uniquely identifies a car
type Vin = string

// LockState Indicates the state of a lock
type LockState = string
