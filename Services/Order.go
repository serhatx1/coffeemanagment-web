package Services

import (
	"Coffee/DB"
	"Coffee/model"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"reflect"
)

func OrderCoffee(c echo.Context) error {
	var Products model.Order
	if err := c.Bind(&Products); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}
	var price float32
	for i := 0; i < len(Products.Beverages); i++ {
		var B model.Beverage
		err := DB.DB.Where("id=?", Products.Beverages[i].ID).First(&B).Error
		if err != nil {
			fmt.Println("hi")
			return c.JSON(http.StatusUnprocessableEntity, err)
		}
		errs := isStockEnough(c, B)
		if errs == false {
			fmt.Println("hi")
			return c.JSON(http.StatusUnprocessableEntity, err)
		}
		price += B.Price
		Products.Beverages[i] = B
	}

	Products.Price = price
	err := DB.DB.Create(&Products).Error
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}
	return c.JSON(http.StatusOK, Products)
}
func isStockEnough(c echo.Context, b model.Beverage) bool {
	var stock model.Stock
	err := DB.DB.Where("id = ?", 1).First(&stock).Error
	if err != nil {
		return false
	}
	var recipe map[string]int
	err = json.Unmarshal([]byte(b.Recipe), &recipe)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return false
	}
	for ingredient, requiredQuantity := range recipe {

		value := reflect.ValueOf(&stock).Elem()
		field := value.FieldByName(ingredient)
		if !(field.IsValid() && field.Kind() == reflect.Int) {
			return false
		}
		if field.Int() < int64(requiredQuantity) {
			fmt.Println("hi")
			return false
		}
	}
	return true
}
func StockControl(c echo.Context) error {
	var bm *model.Beverage
	var b *model.Beverage
	if err := c.Bind(&bm); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}
	fmt.Println("id", bm.ID)
	err := DB.DB.Where("id=?", bm.ID).First(&b).Error
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}
	var recipe map[string]int
	err = json.Unmarshal([]byte(b.Recipe), &recipe)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return c.JSON(http.StatusUnprocessableEntity, err)
	}
	var stock model.Stock
	err = DB.DB.Where("id=?", 1).First(&stock).Error
	for ingredient, requiredQuantity := range recipe {
		value := reflect.ValueOf(&stock).Elem()
		fmt.Println("ingredient", ingredient, value)
		field := value.FieldByName(ingredient)
		if !(field.IsValid() && field.Kind() == reflect.Int) {
			return c.JSON(http.StatusUnprocessableEntity, err)
		}
		if field.Int() < int64(requiredQuantity) {
			return c.JSON(http.StatusUnprocessableEntity, err)
		}
	}
	return c.JSON(http.StatusOK, stock)

}
func OrderComplete(c echo.Context) error {
	var req struct {
		ID uint `json:"id"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "Invalid request payload"})
	}

	var order model.Order
	err := DB.DB.Preload("Beverages").First(&order, req.ID).Error
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Order not found"})
	}

	if order.Completed {
		return c.JSON(http.StatusOK, order)
	}

	order.Completed = true

	for _, beverage := range order.Beverages {
		err := ReduceStock(c, beverage)
		if err != nil {
			order.Completed = false
			DB.DB.Save(&order)
			return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "Failed to reduce stock"})
		}
	}

	err = DB.DB.Save(&order).Error
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "Failed to update order"})
	}

	return c.JSON(http.StatusOK, order)
}

func ReduceStock(c echo.Context, b model.Beverage) error {
	var recipe map[string]int
	err := json.Unmarshal([]byte(b.Recipe), &recipe)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "Error unmarshalling JSON"})
	}

	var stock model.Stock
	err = DB.DB.Where("id=?", 1).First(&stock).Error
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Stock not found"})
	}

	value := reflect.ValueOf(&stock).Elem()
	for ingredient, requiredQuantity := range recipe {
		field := value.FieldByName(ingredient)
		if !(field.IsValid() && field.Kind() == reflect.Int) {
			return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": fmt.Sprintf("Invalid field or type: %s", ingredient)})
		}

		availableQuantity := int(field.Int())
		if availableQuantity < requiredQuantity {
			return c.JSON(http.StatusConflict, map[string]string{"error": fmt.Sprintf("Insufficient stock for %s", ingredient)})
		}

		field.SetInt(int64(availableQuantity - requiredQuantity))
	}

	err = DB.DB.Save(&stock).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update stock"})
	}

	return c.JSON(http.StatusOK, stock)
}
