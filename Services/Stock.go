package Services

import (
	"Coffee/DB"
	"Coffee/model"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateStock(c echo.Context) error {
	var Stock model.Stock
	var existingStock model.Stock
	if err := c.Bind(&Stock); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data: " + err.Error()})
	}
	if err := DB.DB.Unscoped().Delete(existingStock); err == nil {
		fmt.Println(existingStock)
		if err := DB.DB.Delete(&existingStock).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	DB.DB.First(&existingStock).Where("id=?", 1)
	if err := DB.DB.Model(&existingStock).Updates(&Stock).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, Stock)
}

func CheckStock(c echo.Context) error {
	var Stock model.Stock
	err := DB.DB.First(&Stock, 1).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Stock not found"})
	}
	return c.JSON(http.StatusOK, Stock)
}
