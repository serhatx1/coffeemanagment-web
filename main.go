package main

import (
	"Coffee/DB"
	"Coffee/EndPointHandler"
	"Coffee/model"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Echo!")
	})

	DB.Init()

	if err := DB.DB.AutoMigrate(&model.Beverage{}, &model.Stock{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	EndPointHandler.SetupRoutes(e)
	e.Logger.Fatal(e.Start(":8080"))

}
