package api

import (
	"DCar/infrastructure/database"
	"DCar/logic/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

type controller struct {
	crud database.ICRUD
}

// NewController creates a new controller instance and takes a high level CRUD interface as a parameter.
func NewController(crud database.ICRUD) Controller {
	return controller{
		crud,
	}
}

func (c controller) GetCars(ctx echo.Context) error {
	allVins, err := c.crud.ReadAllVins(ctx.Request().Context())

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, allVins)
}

func (c controller) AddCar(ctx echo.Context) error {
	// get request body
	var car model.Car

	// bind errors are unexpected since we validated the request body
	err := ctx.Bind(&car)
	if err != nil {
		return err
	}

	vin, err := c.crud.CreateCar(ctx.Request().Context(), &car)
	if err != nil {
		if database.IsDuplicateKeyError(err) {
			return echo.NewHTTPError(http.StatusConflict, "VIN already exists")
		}
		return err
	}
	return ctx.JSON(http.StatusCreated, vin)
}

func (c controller) DeleteCar(ctx echo.Context, vin model.VinParam) error {
	deleted, err := c.crud.DeleteCar(ctx.Request().Context(), vin)
	if err != nil {
		return err
	}
	if deleted {
		return ctx.NoContent(http.StatusNoContent)
	}
	return echo.NewHTTPError(http.StatusNotFound, "VIN not found")
}

func (c controller) GetCar(ctx echo.Context, vin model.VinParam) error {
	car, err := c.crud.ReadCar(ctx.Request().Context(), vin)
	if err != nil {
		if database.IsNotFoundError(err) {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return err
	}
	return ctx.JSON(http.StatusOK, car)
}