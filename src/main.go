package main

import (
	"DCar/api"
	"DCar/environment"
	"DCar/infrastructure/database"
	"DCar/infrastructure/database/db"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

// newApp allows production as well as testing to create a new Echo instance for the API.
// Configuration values are read from the environment.
func newApp(dbConnection db.IConnection) (*echo.Echo, error) {
	app := echo.New()

	// add OpenAPI validation to the echo instance
	err := api.AddOpenApiValidationMiddleware(app)
	if err != nil {
		return nil, err
	}

	// create a high level CRUD interface for the database and attach it to a controller handling the requests
	err = api.RegisterHandlers(app, api.NewController(database.NewICRUD(dbConnection, environment.GetEnvironment())))
	if err != nil {
		return nil, err
	}

	// Use custom error handling that logs any errors that occur but passes any HTTP errors directly to the client.
	// Any other errors are converted to HTTP 500 errors.
	app.Use(func(fun echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := fun(c); err != nil {
				if err, isHttpError := err.(*echo.HTTPError); isHttpError {
					return err
				}
				app.Logger.Error(err.Error())
				return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
			}
			return nil
		}
	})

	return app, nil
}

func main() {
	// create a new database connection
	dbConnection, err := db.NewDbConnection(environment.GetEnvironment())
	if err != nil {
		log.Fatal(err.Error())
	}

	// close the database connection when the program exits
	defer func() {
		if err := dbConnection.CleanUpDatabase(); err != nil {
			log.Fatal(err)
		}
	}()

	app, err := newApp(dbConnection)
	if err != nil {
		log.Fatal(err.Error())
	}

	// start the server on the configured port
	app.Logger.Fatal(app.Start(fmt.Sprintf(":%d", environment.GetEnvironment().GetAppExposePort())))
}
