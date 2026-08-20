package main

import (
	"net/http"

	"github.com/examsync/pdf-parser/internal/controllers"
	"github.com/examsync/pdf-parser/internal/repositories"
	"github.com/examsync/pdf-parser/internal/services"
	"github.com/examsync/pdf-parser/utils/logger"
	"github.com/labstack/echo/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// registerHandlers instantiates the MVC layers and configures routing endpoints on Echo
func registerHandlers(e *echo.Echo, db *gorm.DB) {
	// Wire Layers (Tightly Coupled Concrete Structs injection)
	repo := repositories.NewExamNotificationRepository(db)
	service := services.NewExamNotificationService(repo)
	controller := controllers.NewExamNotificationController(service)

	// Log API Handler Route Registrations
	if logger.Log != nil {
		logger.Log.WithFields(logrus.Fields{
			"POST /parse":                   "ExamNotificationController.Parse",
			"GET /notifications/:filename": "ExamNotificationController.GetByFileName",
			"GET /health":                   "HealthCheckHandler",
		}).Info("Registering API handler routes")
	}

	// Register Routes
	e.POST("/parse", controller.Parse)
	e.GET("/notifications/:filename", controller.GetByFileName)
	e.GET("/health", func(c *echo.Context) error {
		if logger.Log != nil {
			logger.Log.WithFields(logrus.Fields{
				"handler": "HealthCheck",
				"path":    c.Request().URL.Path,
			}).Debug("Health check request received")
		}
		return c.JSON(http.StatusOK, map[string]string{
			"status": "healthy",
		})
	})
}
