package EndPointHandler

import (
	"Coffee/Services"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	e.POST("/create", Services.CreateCoffee)

}
