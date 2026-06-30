package main

import (
	"net/http"

	"github.com/examsync/pdf-parser/internal/controllers"
	"github.com/examsync/pdf-parser/internal/repositories"
	"github.com/examsync/pdf-parser/internal/services"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// registerHandlers instantiates the MVC layers and configures routing endpoints on Echo
func registerHandlers(e *echo.Echo, db *gorm.DB) {
	// Wire Layers (Tightly Coupled Concrete Structs injection)
	repo := repositories.NewParsedPDFRepository(db)
	service := services.NewParsedPDFService(repo)
	controller := controllers.NewParsedPDFController(service)

	// Register Routes
	e.GET("/pdfs", controller.GetPDFs)
	e.GET("/health", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "healthy",
		})
	})
}
