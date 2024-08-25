package Services

import (
	"Coffee/DB"
	"Coffee/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateCoffee(c echo.Context) error {
	var beverage model.Beverage
	if err := c.Bind(&beverage); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data: " + err.Error()})
	}

	if err := DB.DB.Create(&beverage).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, beverage)
}
func CheckCoffee(c echo.Context) error {
	var beverage model.Beverage
	CoffeeId := c.QueryParam("id")
	err := DB.DB.First(&beverage, CoffeeId).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusCreated, beverage)

}
