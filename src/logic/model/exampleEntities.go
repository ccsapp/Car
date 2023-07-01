package model

import (
	"DCar/infrastructure/database/entities"
	carTypes "github.com/ccsapp/cargotypes"
)

func ExampleDynamicData(trunkLockState entities.LockState) carTypes.DynamicData {
	return carTypes.DynamicData{
		DoorsLockState:      carTypes.UNLOCKED,
		EngineState:         carTypes.OFF,
		FuelLevelPercentage: 23,
		Position: carTypes.DynamicDataPosition{
			Latitude:  49.0069,
			Longitude: 8.4037,
		},
		TrunkLockState: carTypes.DynamicDataLockState(trunkLockState),
	}
}
