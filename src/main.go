package main

import (
	"DCar/api"
	"DCar/infrastructure/database"
	"DCar/infrastructure/database/db"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

func main() {
	e := echo.New()

	// create a new database connection
	dbConnection, err := db.NewDbConnection()
	if err != nil {
		log.Fatal(err.Error())
	}

	// close the database connection when the program exits
	defer func() {
		if err := dbConnection.CleanUpDatabase(); err != nil {
			log.Fatal(err)
		}
	}()

	// add OpenAPI validation to the echo instance
	err = api.AddOpenApiValidationMiddleware(e)
	if err != nil {
		log.Fatal(err.Error())
	}

	// create a high level CRUD interface for the database and attach it to a controller handling the requests
	err = api.RegisterHandlers(e, api.NewController(database.NewICRUD(dbConnection)))
	if err != nil {
		log.Fatal(err.Error())
	}

	// Use custom error handling that logs any errors that occur but passes any HTTP errors directly to the client.
	// Any other errors are converted to HTTP 500 errors.
	e.Use(func(fun echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := fun(c); err != nil {
				if err, isHttpError := err.(*echo.HTTPError); isHttpError {
					return err
				}
				e.Logger.Error(err.Error())
				return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
			}
			return nil
		}
	})

	// start the server on port 80
	e.Logger.Fatal(e.Start(":80"))
}
