package controllers

import (
	"io"
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

// Parse handles the HTTP multipart form request containing a PDF file to extract and store notification details.
func (c *ExamNotificationController) Parse(ctx *echo.Context) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Missing 'file' field in multipart form",
		})
	}

	src, err := file.Open()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to open the uploaded file: " + err.Error(),
		})
	}
	defer src.Close()

	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to read the uploaded file bytes: " + err.Error(),
		})
	}

	notification, err := c.service.ParsePDF(fileBytes)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to parse PDF: " + err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, notification)
}
