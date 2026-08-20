package main

import (
	"net/http"

	"github.com/examsync/pdf-parser/internal/controllers"
	"github.com/examsync/pdf-parser/internal/repositories"
	"github.com/examsync/pdf-parser/internal/services"
	"github.com/examsync/pdf-parser/utils/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// registerHandlers instantiates the MVC layers and configures routing endpoints on Gin
func registerHandlers(r *gin.Engine, db *gorm.DB) {
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
		}).Info("Registering API handler routes on Gin engine")
	}

	// Register Routes
	r.POST("/parse", controller.Parse)
	r.GET("/notifications/:filename", controller.GetByFileName)
	r.GET("/health", func(c *gin.Context) {
		if logger.Log != nil {
			logger.Log.WithFields(logrus.Fields{
				"handler": "HealthCheck",
				"path":    c.Request.URL.Path,
			}).Debug("Health check request received")
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})
}
