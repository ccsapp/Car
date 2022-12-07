package model

func ExampleDynamicData() DynamicData {
	return DynamicData{
		DoorsLockState:      UNLOCKED,
		EngineState:         OFF,
		FuelLevelPercentage: 23,
		Position: DynamicDataPosition{
			Latitude:  49.0069,
			Longitude: 8.4037,
		},
		TrunkLockState: UNLOCKED,
	}
}
