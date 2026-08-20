package controllers

import (
	"io"
	"net/http"

	"github.com/examsync/pdf-parser/internal/services"
	"github.com/examsync/pdf-parser/utils/errors"
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
		return errors.NewBadRequest("Missing 'file' field in multipart form", err)
	}

	src, err := file.Open()
	if err != nil {
		return errors.NewInternal("Failed to open the uploaded file", err)
	}
	defer src.Close()

	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return errors.NewInternal("Failed to read the uploaded file bytes", err)
	}

	notification, err := c.service.ParsePDF(file.Filename, fileBytes)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, notification)
}

// GetByFileName handles HTTP GET requests to read/retrieve a PDF notification by file name.
func (c *ExamNotificationController) GetByFileName(ctx *echo.Context) error {
	fileName := ctx.Param("filename")
	if fileName == "" {
		fileName = ctx.QueryParam("filename")
	}

	if fileName == "" {
		return errors.NewBadRequest("Missing required 'filename' parameter", nil)
	}

	notification, err := c.service.GetByFileName(fileName)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, notification)
}

