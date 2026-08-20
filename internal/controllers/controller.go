package controllers

import (
	"io"
	"net/http"

	"github.com/examsync/pdf-parser/internal/services"
	"github.com/examsync/pdf-parser/utils/errors"
	"github.com/examsync/pdf-parser/utils/logger"
	"github.com/labstack/echo/v5"
	"github.com/sirupsen/logrus"
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
		if logger.Log != nil {
			logger.Log.WithFields(logrus.Fields{
				"handler": "Parse",
				"status":  http.StatusBadRequest,
				"error":   err.Error(),
			}).Warn("Failed to extract multipart form file from request")
		}
		return errors.NewBadRequest("Missing 'file' field in multipart form", err)
	}

	if logger.Log != nil {
		logger.Log.WithFields(logrus.Fields{
			"handler":  "Parse",
			"filename": file.Filename,
			"size":     file.Size,
		}).Info("Receiving PDF file upload request for parsing")
	}

	src, err := file.Open()
	if err != nil {
		if logger.Log != nil {
			logger.Log.WithFields(logrus.Fields{
				"handler":  "Parse",
				"filename": file.Filename,
				"error":    err.Error(),
			}).Error("Failed to open uploaded file stream")
		}
		return errors.NewInternal("Failed to open the uploaded file", err)
	}
	defer src.Close()

	fileBytes, err := io.ReadAll(src)
	if err != nil {
		if logger.Log != nil {
			logger.Log.WithFields(logrus.Fields{
				"handler":  "Parse",
				"filename": file.Filename,
				"error":    err.Error(),
			}).Error("Failed to read uploaded file bytes")
		}
		return errors.NewInternal("Failed to read the uploaded file bytes", err)
	}

	notification, err := c.service.ParsePDF(file.Filename, fileBytes)
	if err != nil {
		return err
	}

	if logger.Log != nil {
		logger.Log.WithFields(logrus.Fields{
			"handler":         "Parse",
			"notification_id": notification.ID,
			"filename":        notification.FileName,
			"language":        notification.Language,
			"text_length":     len(notification.RawText),
		}).Info("Successfully parsed and processed PDF notification")
	}

	return ctx.JSON(http.StatusCreated, notification)
}

// GetByFileName handles HTTP GET requests to read/retrieve a PDF notification by file name.
func (c *ExamNotificationController) GetByFileName(ctx *echo.Context) error {
	fileName := ctx.Param("filename")
	if fileName == "" {
		fileName = ctx.QueryParam("filename")
	}

	if logger.Log != nil {
		logger.Log.WithFields(logrus.Fields{
			"handler":  "GetByFileName",
			"filename": fileName,
		}).Info("Receiving HTTP GET request to read notification by file name")
	}

	if fileName == "" {
		if logger.Log != nil {
			logger.Log.WithFields(logrus.Fields{
				"handler": "GetByFileName",
				"status":  http.StatusBadRequest,
			}).Warn("Missing required filename parameter in request")
		}
		return errors.NewBadRequest("Missing required 'filename' parameter", nil)
	}

	notification, err := c.service.GetByFileName(fileName)
	if err != nil {
		if logger.Log != nil {
			logger.Log.WithFields(logrus.Fields{
				"handler":  "GetByFileName",
				"filename": fileName,
				"error":    err.Error(),
			}).Warn("Failed to retrieve PDF notification by file name")
		}
		return err
	}

	if logger.Log != nil {
		logger.Log.WithFields(logrus.Fields{
			"handler":         "GetByFileName",
			"notification_id": notification.ID,
			"filename":        notification.FileName,
			"language":        notification.Language,
		}).Info("Successfully retrieved PDF notification by file name")
	}

	return ctx.JSON(http.StatusOK, notification)
}
