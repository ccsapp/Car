// Package api provides primitives to interact with the openapi HTTP API.
package api

import (
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/labstack/echo/v4"

	"DCar/logic/model"
)

// Controller represents all server handlers.
type Controller interface {
	// GetCars Get VINs of all Cars
	// (GET /cars)
	GetCars(ctx echo.Context) error
	// AddCar Add a New Vehicle
	// (POST /cars)
	AddCar(ctx echo.Context) error
	// DeleteCar DeleteOne a Car With All Components
	// (DELETE /cars/{vin})
	DeleteCar(ctx echo.Context, vin model.VinParam) error
	// GetCar Get All Information About a Specific Car
	// (GET /cars/{vin})
	GetCar(ctx echo.Context, vin model.VinParam) error
}

// ControllerWrapper converts echo contexts to parameters.
type ControllerWrapper struct {
	Handler Controller
}

// GetCars converts echo context to params.
func (w *ControllerWrapper) GetCars(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetCars(ctx)
	return err
}

// AddCar converts echo context to params.
func (w *ControllerWrapper) AddCar(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.AddCar(ctx)
	return err
}

// DeleteCar converts echo context to params.
func (w *ControllerWrapper) DeleteCar(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "vin" -------------
	var vin model.VinParam

	err = runtime.BindStyledParameterWithLocation("simple", false, "vin", runtime.ParamLocationPath, ctx.Param("vin"), &vin)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter vin: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.DeleteCar(ctx, vin)
	return err
}

// GetCar converts echo context to params.
func (w *ControllerWrapper) GetCar(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "vin" -------------
	var vin model.VinParam

	err = runtime.BindStyledParameterWithLocation("simple", false, "vin", runtime.ParamLocationPath, ctx.Param("vin"), &vin)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter vin: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetCar(ctx, vin)
	return err
}

// EchoRouter
// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, controller Controller) error {
	return RegisterHandlersWithBaseURL(router, controller, "")
}

// RegisterHandlersWithBaseURL
// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, controller Controller, baseURL string) error {
	wrapper := ControllerWrapper{
		Handler: controller,
	}

	router.GET(baseURL+"/cars", wrapper.GetCars)
	router.POST(baseURL+"/cars", wrapper.AddCar)
	router.DELETE(baseURL+"/cars/:vin", wrapper.DeleteCar)
	router.GET(baseURL+"/cars/:vin", wrapper.GetCar)

	return nil
}
