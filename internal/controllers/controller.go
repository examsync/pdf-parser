package controllers

import (
	"net/http"

	"github.com/examsync/pdf-parser/internal/services"
	"github.com/labstack/echo/v5"
)

// ExamNotificationController handles HTTP requests for exam notifications.
type ExamNotificationController struct {
	service *services.ExamNotificationService
}

// NewExamNotificationController creates a new instance of ExamNotificationController.
func NewExamNotificationController(service *services.ExamNotificationService) *ExamNotificationController {
	return &ExamNotificationController{service: service}
}

// GetNotifications handles the HTTP request to get all exam notifications.
func (c *ExamNotificationController) GetNotifications(ctx *echo.Context) error {
	notifications, err := c.service.GetNotifications()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, notifications)
}
