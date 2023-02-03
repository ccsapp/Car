package main

import (
	"DCar/infrastructure/database"
	"DCar/infrastructure/database/db"
	"DCar/testdata"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
	"time"
)

type ApiTestSuite struct {
	suite.Suite
	dbConnection db.IConnection
	collection   string
	app          *echo.Echo
}

func (suite *ApiTestSuite) SetupSuite() {
	// load the environment variables for the database layer
	config, err := db.LoadConfigFromFile("testdata/testdb.env")
	if err != nil {
		suite.T().Fatal(err.Error())
	}

	// generate a collection name so that concurrent executions do not interfere
	config.CollectionPrefix = fmt.Sprintf("test-%d-", time.Now().Unix())
	suite.collection = config.CollectionPrefix + database.CarsCollectionBaseName

	// create a new database connection
	dbConnection, err := db.NewDbConnection(config)
	if err != nil {
		suite.T().Fatal(err.Error())
	}

	suite.dbConnection = dbConnection

	app, err := newApp(dbConnection, config)
	if err != nil {
		suite.T().Fatal(err.Error())
	}
	suite.app = app
}

func (suite *ApiTestSuite) TearDownSuite() {
	// close the database connection when the program exits
	if err := suite.dbConnection.CleanUpDatabase(); err != nil {
		suite.T().Fatal(err)
	}
}

func (suite *ApiTestSuite) TearDownTest() {
	// clear the collection after each test
	if err := suite.dbConnection.DropCollection(context.Background(), suite.collection); err != nil {
		suite.T().Fatal(err)
	}
}

func TestApiTestSuite(t *testing.T) {
	suite.Run(t, new(ApiTestSuite))
}

func newApiTest(handler http.Handler, name string) *apitest.APITest {
	return apitest.New(name).
		Debug().
		Handler(handler).
		Report(apitest.SequenceDiagram())
}

func (suite *ApiTestSuite) TestVinOverview_empty() {
	newApiTest(suite.app, "Get the Empty VIN Overview").
		Get("/cars").
		Expect(suite.T()).
		Status(http.StatusOK).
		Body("[]").
		End()
}

func (suite *ApiTestSuite) TestGetCar_invalidFormat() {
	newApiTest(suite.app, "Get a Car with Invalid VIN").
		Get("/cars/abc").
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()
}

func (suite *ApiTestSuite) TestAddCar_success() {
	// add the example car to the database
	newApiTest(suite.app, "Add a New Car").
		Post("/cars").
		JSON(testdata.ExampleCar).
		Expect(suite.T()).
		Status(http.StatusCreated).
		Body(testdata.ExampleCarVin).
		End()

	newApiTest(suite.app, "Get the VIN Overview").
		Get("/cars").
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCarVinArray).
		End()

	// validate that the car was added to the database
	newApiTest(suite.app, "Get the New Car").
		Get("/cars/" + testdata.ExampleCarVinString).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCarWithDynamicData).
		End()
}

func (suite *ApiTestSuite) TestAddCar_duplicate() {
	// add the example car to the database
	newApiTest(suite.app, "Add a New Car").
		Post("/cars").
		JSON(testdata.ExampleCar).
		Expect(suite.T()).
		Status(http.StatusCreated).
		End()

	// try to add another car with the same VIN
	newApiTest(suite.app, "Add a Duplicate Car with the same VIN").
		Post("/cars").
		JSON(testdata.ExampleCarDuplicate).
		Expect(suite.T()).
		Status(http.StatusConflict).
		End()

	// validate that the original car did not change
	newApiTest(suite.app, "Get the Original Car").
		Get("/cars/" + testdata.ExampleCarVinString).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCarWithDynamicData).
		End()

	// validate that the VIN does only appear once in the overview
	newApiTest(suite.app, "Get the VIN Overview").
		Get("/cars").
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCarVinArray).
		End()
}

func (suite *ApiTestSuite) TestAddCar_invalidJson() {
	newApiTest(suite.app, "Add Invalid JSON as Car").
		Post("/cars").
		JSON(testdata.ExampleNoJson).
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()

	suite.TestVinOverview_empty()
}

func (suite *ApiTestSuite) TestAddCar_invalidCar() {
	newApiTest(suite.app, "Add Invalid Object as Car").
		Post("/cars").
		JSON(testdata.ExampleNoCar).
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()

	suite.TestVinOverview_empty()
}

func (suite *ApiTestSuite) TestAddCar_invalidCarEnums() {
	newApiTest(suite.app, "Add Car with Invalid Enum Values").
		Post("/cars").
		JSON(testdata.ExampleCarWrongEnum).
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()

	suite.TestVinOverview_empty()
}

func (suite *ApiTestSuite) TestRemoveCar_noSuchCar() {
	newApiTest(suite.app, "Remove Non-existent Car").
		Delete("/cars/" + testdata.ExampleCarVinString).
		Expect(suite.T()).
		Status(http.StatusNotFound).
		End()
}

func (suite *ApiTestSuite) TestRemoveCar_invalidFormat() {
	newApiTest(suite.app, "Remove Non-existent Car").
		Delete("/cars/xyz").
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()
}

func (suite *ApiTestSuite) TestRemoveCar_success() {
	// add two cars
	newApiTest(suite.app, "Add the First Car").
		Post("/cars").
		JSON(testdata.ExampleCar).
		Expect(suite.T()).
		Status(http.StatusCreated).
		End()

	newApiTest(suite.app, "Add the Second Car").
		Post("/cars").
		JSON(testdata.ExampleCar2).
		Expect(suite.T()).
		Status(http.StatusCreated).
		End()

	// verify that both cars exist
	newApiTest(suite.app, "Get the VIN Overview of Both Cars").
		Get("/cars").
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ArrayString(testdata.ExampleCarVin, testdata.ExampleCar2Vin)).
		End()

	// remove one car
	newApiTest(suite.app, "Remove Second Car").
		Delete("/cars/" + testdata.ExampleCar2VinString).
		Expect(suite.T()).
		Status(http.StatusNoContent).
		End()

	// verify that the car was removed
	newApiTest(suite.app, "Request Removed Car").
		Get("/cars/" + testdata.ExampleCar2VinString).
		Expect(suite.T()).
		Status(http.StatusNotFound).
		End()

	// verify that the other car is still there
	newApiTest(suite.app, "Request Remaining Car").
		Get("/cars/" + testdata.ExampleCarVinString).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCarWithDynamicData).
		End()

	newApiTest(suite.app, "Get the VIN Overview").
		Get("/cars").
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCarVinArray).
		End()
}

func (suite *ApiTestSuite) TestChangeTrunkLockState_noSuchCar() {
	newApiTest(suite.app, "Change Trunk Lock State of Non-existent Car").
		Put("/cars/" + testdata.ExampleCarVinString + "/trunkLock").
		JSON(testdata.QuoteString("UNLOCKED")).
		Expect(suite.T()).
		Status(http.StatusNotFound).
		End()
}

func (suite *ApiTestSuite) TestChangeTrunkLockState_invalidVinFormat() {
	newApiTest(suite.app, "Change Trunk Lock State of Car with Invalid VIN").
		Put("/cars/xyz/trunkLock").
		JSON(testdata.QuoteString("UNLOCKED")).
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()
}

func (suite *ApiTestSuite) TestChangeTrunkLockState_invalidLockState() {
	newApiTest(suite.app, "Add the First Car").
		Post("/cars").
		JSON(testdata.ExampleCar).
		Expect(suite.T()).
		Status(http.StatusCreated).
		End()

	newApiTest(suite.app, "Change Trunk Lock State of Car with Invalid Lock State").
		Put("/cars/" + testdata.ExampleCarVinString + "/trunkLock").
		JSON(testdata.QuoteString("xyz")).
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()
}

func (suite *ApiTestSuite) TestChangeTrunkLockState_successUnchanged() {
	newApiTest(suite.app, "Add the First Car").
		Post("/cars").
		JSON(testdata.ExampleCar).
		Expect(suite.T()).
		Status(http.StatusCreated).
		End()

	newApiTest(suite.app, "Set Lock State to LOCKED").
		Put("/cars/" + testdata.ExampleCarVinString + "/trunkLock").
		JSON(testdata.QuoteString("LOCKED")).
		Expect(suite.T()).
		Status(http.StatusNoContent).
		End()

	newApiTest(suite.app, "Get the Car").
		Get("/cars/" + testdata.ExampleCarVinString).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCarWithDynamicData).
		End()
}

func (suite *ApiTestSuite) TestChangeTrunkLockState_successChanged() {
	newApiTest(suite.app, "Add the First Car").
		Post("/cars").
		JSON(testdata.ExampleCar).
		Expect(suite.T()).
		Status(http.StatusCreated).
		End()

	newApiTest(suite.app, "Set Lock State to UNLOCKED").
		Put("/cars/" + testdata.ExampleCarVinString + "/trunkLock").
		JSON(testdata.QuoteString("UNLOCKED")).
		Expect(suite.T()).
		Status(http.StatusNoContent).
		End()

	newApiTest(suite.app, "Get the Car").
		Get("/cars/" + testdata.ExampleCarVinString).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCarUnlockedTrunk).
		End()

	newApiTest(suite.app, "Lock the Trunk Again").
		Put("/cars/" + testdata.ExampleCarVinString + "/trunkLock").
		JSON(testdata.QuoteString("LOCKED")).
		Expect(suite.T()).
		Status(http.StatusNoContent).
		End()

	newApiTest(suite.app, "Get the Car").
		Get("/cars/" + testdata.ExampleCarVinString).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCarWithDynamicData).
		End()
}
