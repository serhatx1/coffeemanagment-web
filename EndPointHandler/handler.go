package EndPointHandler

import (
	"Coffee/Services"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	e.POST("/coffee/create", Services.CreateCoffee)
	e.POST("/stock/create", Services.CreateStock)
	e.POST("/stock/check", Services.CheckStock)
	e.POST("/coffee/check", Services.CheckCoffee)
	e.POST("/product/order", Services.OrderCoffee)
	e.POST("/coffee/order/choose", Services.StockControl)
	e.POST("/coffee/order/complete", Services.OrderComplete)

}
