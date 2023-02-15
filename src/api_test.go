package main

import (
	"DCar/infrastructure/database"
	"DCar/infrastructure/database/db"
	"DCar/testdata"
	"DCar/testhelpers"
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
	dbConnection       db.IConnection
	collection         string
	app                *echo.Echo
	recordingFormatter *testhelpers.RecordingFormatter
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

func (suite *ApiTestSuite) SetupTest() {
	suite.recordingFormatter = testhelpers.NewRecordingFormatter()
}

func (suite *ApiTestSuite) TearDownSuite() {
	// close the database connection when the program exits
	if err := suite.dbConnection.CleanUpDatabase(); err != nil {
		suite.T().Fatal(err)
	}
}

func (suite *ApiTestSuite) TearDownTest() {
	// generate the sequence diagram for the test
	suite.recordingFormatter.SetOutFileName(suite.T().Name())
	suite.recordingFormatter.SetTitle(suite.T().Name())

	diagramFormatter := apitest.SequenceDiagram()
	diagramFormatter.Format(suite.recordingFormatter.GetRecorder())

	// clear the collection after each test
	if err := suite.dbConnection.DropCollection(context.Background(), suite.collection); err != nil {
		suite.T().Fatal(err)
	}
}

func TestApiTestSuite(t *testing.T) {
	suite.Run(t, new(ApiTestSuite))
}

func (suite *ApiTestSuite) newApiTest() *apitest.APITest {
	return apitest.New().
		Debug().
		Handler(suite.app).
		Report(suite.recordingFormatter)
}

func (suite *ApiTestSuite) TestVinOverview_empty() {
	suite.newApiTest().
		Get("/cars").
		Expect(suite.T()).
		Status(http.StatusOK).
		Body("[]").
		End()
}

func (suite *ApiTestSuite) TestGetCar_invalidFormat() {
	suite.newApiTest().
		Get("/cars/abc").
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()
}

func (suite *ApiTestSuite) TestAddCar_success() {
	// add the example car to the database
	suite.newApiTest().
		Post("/cars").
		JSON(testdata.ExampleCar).
		Expect(suite.T()).
		Status(http.StatusCreated).
		Body(testdata.ExampleCarVin).
		End()

	suite.newApiTest().
		Get("/cars").
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCarVinArray).
		End()

	// validate that the car was added to the database
	suite.newApiTest().
		Get("/cars/" + testdata.ExampleCarVinString).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCarWithDynamicData).
		End()
}

func (suite *ApiTestSuite) TestAddCar_duplicate() {
	// add the example car to the database
	suite.newApiTest().
		Post("/cars").
		JSON(testdata.ExampleCar).
		Expect(suite.T()).
		Status(http.StatusCreated).
		End()

	// try to add another car with the same VIN
	suite.newApiTest().
		Post("/cars").
		JSON(testdata.ExampleCarDuplicate).
		Expect(suite.T()).
		Status(http.StatusConflict).
		End()

	// validate that the original car did not change
	suite.newApiTest().
		Get("/cars/" + testdata.ExampleCarVinString).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCarWithDynamicData).
		End()

	// validate that the VIN does only appear once in the overview
	suite.newApiTest().
		Get("/cars").
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCarVinArray).
		End()
}

func (suite *ApiTestSuite) TestAddCar_invalidJson() {
	suite.newApiTest().
		Post("/cars").
		JSON(testdata.ExampleNoJson).
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()

	suite.TestVinOverview_empty()
}

func (suite *ApiTestSuite) TestAddCar_invalidCar() {
	suite.newApiTest().
		Post("/cars").
		JSON(testdata.ExampleNoCar).
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()

	suite.TestVinOverview_empty()
}

func (suite *ApiTestSuite) TestAddCar_invalidCarEnums() {
	suite.newApiTest().
		Post("/cars").
		JSON(testdata.ExampleCarWrongEnum).
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()

	suite.TestVinOverview_empty()
}

func (suite *ApiTestSuite) TestRemoveCar_noSuchCar() {
	suite.newApiTest().
		Delete("/cars/" + testdata.ExampleCarVinString).
		Expect(suite.T()).
		Status(http.StatusNotFound).
		End()
}

func (suite *ApiTestSuite) TestRemoveCar_invalidFormat() {
	suite.newApiTest().
		Delete("/cars/xyz").
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()
}

func (suite *ApiTestSuite) TestRemoveCar_success() {
	// add two cars
	suite.newApiTest().
		Post("/cars").
		JSON(testdata.ExampleCar).
		Expect(suite.T()).
		Status(http.StatusCreated).
		End()

	suite.newApiTest().
		Post("/cars").
		JSON(testdata.ExampleCar2).
		Expect(suite.T()).
		Status(http.StatusCreated).
		End()

	// verify that both cars exist
	suite.newApiTest().
		Get("/cars").
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ArrayString(testdata.ExampleCarVin, testdata.ExampleCar2Vin)).
		End()

	// remove one car
	suite.newApiTest().
		Delete("/cars/" + testdata.ExampleCar2VinString).
		Expect(suite.T()).
		Status(http.StatusNoContent).
		End()

	// verify that the car was removed
	suite.newApiTest().
		Get("/cars/" + testdata.ExampleCar2VinString).
		Expect(suite.T()).
		Status(http.StatusNotFound).
		End()

	// verify that the other car is still there
	suite.newApiTest().
		Get("/cars/" + testdata.ExampleCarVinString).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCarWithDynamicData).
		End()

	suite.newApiTest().
		Get("/cars").
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCarVinArray).
		End()
}

func (suite *ApiTestSuite) TestChangeTrunkLockState_noSuchCar() {
	suite.newApiTest().
		Put("/cars/" + testdata.ExampleCarVinString + "/trunkLock").
		JSON(testdata.QuoteString("UNLOCKED")).
		Expect(suite.T()).
		Status(http.StatusNotFound).
		End()
}

func (suite *ApiTestSuite) TestChangeTrunkLockState_invalidVinFormat() {
	suite.newApiTest().
		Put("/cars/xyz/trunkLock").
		JSON(testdata.QuoteString("UNLOCKED")).
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()
}

func (suite *ApiTestSuite) TestChangeTrunkLockState_invalidLockState() {
	suite.newApiTest().
		Post("/cars").
		JSON(testdata.ExampleCar).
		Expect(suite.T()).
		Status(http.StatusCreated).
		End()

	suite.newApiTest().
		Put("/cars/" + testdata.ExampleCarVinString + "/trunkLock").
		JSON(testdata.QuoteString("xyz")).
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()
}

func (suite *ApiTestSuite) TestChangeTrunkLockState_successUnchanged() {
	suite.newApiTest().
		Post("/cars").
		JSON(testdata.ExampleCar).
		Expect(suite.T()).
		Status(http.StatusCreated).
		End()

	suite.newApiTest().
		Put("/cars/" + testdata.ExampleCarVinString + "/trunkLock").
		JSON(testdata.QuoteString("LOCKED")).
		Expect(suite.T()).
		Status(http.StatusNoContent).
		End()

	suite.newApiTest().
		Get("/cars/" + testdata.ExampleCarVinString).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCarWithDynamicData).
		End()
}

func (suite *ApiTestSuite) TestChangeTrunkLockState_successChanged() {
	suite.newApiTest().
		Post("/cars").
		JSON(testdata.ExampleCar).
		Expect(suite.T()).
		Status(http.StatusCreated).
		End()

	suite.newApiTest().
		Put("/cars/" + testdata.ExampleCarVinString + "/trunkLock").
		JSON(testdata.QuoteString("UNLOCKED")).
		Expect(suite.T()).
		Status(http.StatusNoContent).
		End()

	suite.newApiTest().
		Get("/cars/" + testdata.ExampleCarVinString).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCarUnlockedTrunk).
		End()

	suite.newApiTest().
		Put("/cars/" + testdata.ExampleCarVinString + "/trunkLock").
		JSON(testdata.QuoteString("LOCKED")).
		Expect(suite.T()).
		Status(http.StatusNoContent).
		End()

	suite.newApiTest().
		Get("/cars/" + testdata.ExampleCarVinString).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(testdata.ExampleCarWithDynamicData).
		End()
}
