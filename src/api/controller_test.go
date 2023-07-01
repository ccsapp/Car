package api

import (
	"DCar/mocks"
	"context"
	"errors"
	carTypes "github.com/ccsapp/cargotypes"
	openapiTypes "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
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
		TrunkLockState: carTypes.UNLOCKED,
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

func TestController_GetCars_success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	vins := []string{"12345678901234567", "12345678901234568", "12345678901234569"}

	request, _ := http.NewRequestWithContext(ctx, "GET", "https://example.com/cars", nil)

	mockEchoContext := mocks.NewMockContext(ctrl)
	mockCrud := mocks.NewMockICRUD(ctrl)

	mockEchoContext.EXPECT().Request().Return(request)
	mockCrud.
		EXPECT().ReadAllVins(ctx).Return(vins, nil)
	mockEchoContext.EXPECT().JSON(http.StatusOK, vins)

	controller := NewController(mockCrud)
	err := controller.GetCars(mockEchoContext)
	assert.Nil(t, err)
}

func TestController_GetCars_crudError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	request, _ := http.NewRequestWithContext(ctx, "GET", "https://example.com/cars", nil)

	mockEchoContext := mocks.NewMockContext(ctrl)
	mockCrud := mocks.NewMockICRUD(ctrl)

	mockEchoContext.EXPECT().Request().Return(request)

	crudError := errors.New("crud error")
	mockCrud.
		EXPECT().
		ReadAllVins(ctx).Return(nil, crudError)

	controller := NewController(mockCrud)
	err := controller.GetCars(mockEchoContext)
	assert.ErrorIs(t, err, crudError)
}

func TestController_AddCar_success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	request, _ := http.NewRequestWithContext(ctx, "POST", "https://example.com/cars", nil)

	mockEchoContext := mocks.NewMockContext(ctrl)
	mockCrud := mocks.NewMockICRUD(ctrl)
	mockEchoContext.EXPECT().Bind(gomock.Any()).SetArg(0, exampleModelCar).Return(nil)
	mockCrud.
		EXPECT().CreateCar(ctx, &exampleModelCar).Return(exampleModelCar.Vin, nil)
	mockEchoContext.EXPECT().Request().Return(request)
	mockEchoContext.EXPECT().JSON(http.StatusCreated, exampleModelCar.Vin)

	controller := NewController(mockCrud)
	err := controller.AddCar(mockEchoContext)
	assert.Nil(t, err)

}

func TestController_AddCar_duplicate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	request, _ := http.NewRequestWithContext(ctx, "POST", "https://example.com/cars", nil)

	mockEchoContext := mocks.NewMockContext(ctrl)
	mockCrud := mocks.NewMockICRUD(ctrl)
	mockEchoContext.EXPECT().Bind(gomock.Any()).SetArg(0, exampleModelCar).Return(nil)
	mockEchoContext.EXPECT().Request().Return(request)
	mockCrud.
		EXPECT().CreateCar(ctx, &exampleModelCar).Return("",
		mongo.WriteException{WriteErrors: []mongo.WriteError{{Code: 11000}}})

	controller := NewController(mockCrud)
	err := controller.AddCar(mockEchoContext)
	assert.Equal(t, echo.NewHTTPError(http.StatusConflict, "VIN already exists"), err)
}

func TestController_AddCar_unexpectedBindError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEchoContext := mocks.NewMockContext(ctrl)
	mockCrud := mocks.NewMockICRUD(ctrl)

	bindError := errors.New("bind error")
	mockEchoContext.EXPECT().Bind(gomock.Any()).Return(bindError)

	controller := NewController(mockCrud)
	err := controller.AddCar(mockEchoContext)
	assert.ErrorIs(t, err, bindError)
}

func TestController_AddCar_unexpectedCrudError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	request, _ := http.NewRequestWithContext(ctx, "POST", "https://example.com/cars", nil)

	mockEchoContext := mocks.NewMockContext(ctrl)
	mockCrud := mocks.NewMockICRUD(ctrl)
	mockEchoContext.EXPECT().Bind(gomock.Any()).SetArg(0, exampleModelCar).Return(nil)
	mockEchoContext.EXPECT().Request().Return(request)
	crudError := errors.New("crud error")
	mockCrud.
		EXPECT().CreateCar(ctx, &exampleModelCar).Return("", crudError)

	controller := NewController(mockCrud)
	err := controller.AddCar(mockEchoContext)
	assert.ErrorIs(t, err, crudError)
}

func TestController_DeleteCar_success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	vin := "12345678901234569"

	request, _ := http.NewRequestWithContext(ctx, "GET", "https://example.com/cars", nil)

	mockEchoContext := mocks.NewMockContext(ctrl)
	mockCrud := mocks.NewMockICRUD(ctrl)

	mockEchoContext.EXPECT().Request().Return(request)
	mockCrud.
		EXPECT().DeleteCar(ctx, vin).Return(true, nil)
	mockEchoContext.EXPECT().NoContent(http.StatusNoContent)

	controller := NewController(mockCrud)
	err := controller.DeleteCar(mockEchoContext, vin)
	assert.Nil(t, err)
}

func TestController_DeleteCar_notFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	vin := "12345678901234569"

	request, _ := http.NewRequestWithContext(ctx, "GET", "https://example.com/cars", nil)

	mockEchoContext := mocks.NewMockContext(ctrl)
	mockCrud := mocks.NewMockICRUD(ctrl)

	mockEchoContext.EXPECT().Request().Return(request)
	mockCrud.
		EXPECT().DeleteCar(ctx, vin).Return(false, nil)

	controller := NewController(mockCrud)
	err := controller.DeleteCar(mockEchoContext, vin)
	assert.Equal(t, echo.NewHTTPError(http.StatusNotFound, "VIN not found"), err)
}

func TestController_DeleteCar_unexpectedCrudError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	vin := "12345678901234569"

	request, _ := http.NewRequestWithContext(ctx, "GET", "https://example.com/cars", nil)

	mockEchoContext := mocks.NewMockContext(ctrl)
	mockCrud := mocks.NewMockICRUD(ctrl)

	mockEchoContext.EXPECT().Request().Return(request)
	crudError := errors.New("crud error")
	mockCrud.
		EXPECT().DeleteCar(ctx, vin).Return(false, crudError)

	controller := NewController(mockCrud)
	err := controller.DeleteCar(mockEchoContext, vin)
	assert.ErrorIs(t, err, crudError)
}

func TestController_GetCar_success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	vin := "12345678901234569"

	request, _ := http.NewRequestWithContext(ctx, "GET", "https://example.com/cars", nil)

	mockEchoContext := mocks.NewMockContext(ctrl)
	mockCrud := mocks.NewMockICRUD(ctrl)

	mockEchoContext.EXPECT().Request().Return(request)
	mockCrud.
		EXPECT().ReadCar(ctx, vin).Return(exampleModelCar, nil)
	mockEchoContext.EXPECT().JSON(http.StatusOK, exampleModelCar)

	controller := NewController(mockCrud)
	err := controller.GetCar(mockEchoContext, vin)
	assert.Nil(t, err)
}

func TestController_GetCar_notFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	vin := "12345678901234569"

	request, _ := http.NewRequestWithContext(ctx, "GET", "https://example.com/cars", nil)

	mockEchoContext := mocks.NewMockContext(ctrl)
	mockCrud := mocks.NewMockICRUD(ctrl)

	mockEchoContext.EXPECT().Request().Return(request)
	mockCrud.
		EXPECT().
		ReadCar(ctx, vin).Return(carTypes.Car{}, mongo.ErrNoDocuments)

	controller := NewController(mockCrud)
	err := controller.GetCar(mockEchoContext, vin)
	assert.Equal(t, echo.NewHTTPError(http.StatusNotFound), err)
}

func TestController_GetCar_unexpectedCrudError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	vin := "12345678901234569"

	request, _ := http.NewRequestWithContext(ctx, "GET", "https://example.com/cars", nil)

	mockEchoContext := mocks.NewMockContext(ctrl)
	mockCrud := mocks.NewMockICRUD(ctrl)

	mockEchoContext.EXPECT().Request().Return(request)
	crudError := errors.New("crud error")
	mockCrud.
		EXPECT().
		ReadCar(ctx, vin).Return(carTypes.Car{}, crudError)

	controller := NewController(mockCrud)
	err := controller.GetCar(mockEchoContext, vin)
	assert.ErrorIs(t, err, crudError)
}

func TestController_ChangeTrunkLockState_success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	vin := "12345678901234569"

	request, _ := http.NewRequestWithContext(ctx, "GET", "https://example.com/cars", nil)

	mockEchoContext := mocks.NewMockContext(ctrl)
	mockCrud := mocks.NewMockICRUD(ctrl)

	mockEchoContext.EXPECT().Request().Return(request)
	mockEchoContext.EXPECT().Bind(gomock.Any()).SetArg(0, carTypes.UNLOCKED).Return(nil)
	mockCrud.EXPECT().SetTrunkLockState(ctx, vin, carTypes.UNLOCKED).Return(nil)
	mockEchoContext.EXPECT().NoContent(http.StatusNoContent)

	controller := NewController(mockCrud)
	err := controller.ChangeTrunkLockState(mockEchoContext, vin)
	assert.Nil(t, err)
}

func TestController_ChangeTrunkLockState_carNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	vin := "12345678901234569"

	request, _ := http.NewRequestWithContext(ctx, "GET", "https://example.com/cars", nil)

	mockEchoContext := mocks.NewMockContext(ctrl)
	mockCrud := mocks.NewMockICRUD(ctrl)

	mockEchoContext.EXPECT().Request().Return(request)
	mockEchoContext.EXPECT().Bind(gomock.Any()).SetArg(0, carTypes.UNLOCKED).Return(nil)
	mockCrud.EXPECT().SetTrunkLockState(ctx, vin, carTypes.UNLOCKED).Return(mongo.ErrNoDocuments)

	controller := NewController(mockCrud)
	err := controller.ChangeTrunkLockState(mockEchoContext, vin)
	assert.Equal(t, echo.NewHTTPError(http.StatusNotFound, "VIN not found"), err)
}

func TestController_ChangeTrunkLockState_crudError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	vin := "12345678901234569"

	request, _ := http.NewRequestWithContext(ctx, "GET", "https://example.com/cars", nil)

	mockEchoContext := mocks.NewMockContext(ctrl)
	mockCrud := mocks.NewMockICRUD(ctrl)

	crudError := errors.New("crud error")

	mockEchoContext.EXPECT().Request().Return(request)
	mockEchoContext.EXPECT().Bind(gomock.Any()).SetArg(0, carTypes.LOCKED).Return(nil)
	mockCrud.EXPECT().SetTrunkLockState(ctx, vin, carTypes.LOCKED).Return(crudError)

	controller := NewController(mockCrud)
	err := controller.ChangeTrunkLockState(mockEchoContext, vin)
	assert.ErrorIs(t, err, crudError)
}
