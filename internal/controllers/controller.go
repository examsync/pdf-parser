package controllers

import (
	"io"
	"net/http"

	"github.com/examsync/pdf-parser/internal/services"
	"github.com/examsync/pdf-parser/utils/errors"
	"github.com/examsync/pdf-parser/utils/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ExamNotificationController handles HTTP requests for exam notifications using Gin Gonic.
type ExamNotificationController struct {
	service *services.ExamNotificationService
}

// NewExamNotificationController creates a new instance of ExamNotificationController.
func NewExamNotificationController(service *services.ExamNotificationService) *ExamNotificationController {
	return &ExamNotificationController{service: service}
}

// Parse handles the HTTP multipart form request containing a PDF file to extract and store notification details.
func (c *ExamNotificationController) Parse(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		if logger.Log != nil {
			logger.Log.WithFields(logrus.Fields{
				"handler": "Parse",
				"status":  http.StatusBadRequest,
				"error":   err.Error(),
			}).Warn("--- [HANDLER: Parse] Failed to extract multipart form file from request ---")
		}
		_ = ctx.Error(errors.NewBadRequest("Missing 'file' field in multipart form", err))
		return
	}

	if logger.Log != nil {
		logger.Log.WithFields(logrus.Fields{
			"handler":  "Parse",
			"filename": file.Filename,
			"size":     file.Size,
		}).Info("--- [HANDLER: Parse] Receiving PDF file upload request for parsing ---")
	}

	src, err := file.Open()
	if err != nil {
		if logger.Log != nil {
			logger.Log.WithFields(logrus.Fields{
				"handler":  "Parse",
				"filename": file.Filename,
				"error":    err.Error(),
			}).Error("--- [HANDLER: Parse] Failed to open uploaded file stream ---")
		}
		_ = ctx.Error(errors.NewInternal("Failed to open the uploaded file", err))
		return
	}
	defer src.Close()

	fileBytes, err := io.ReadAll(src)
	if err != nil {
		if logger.Log != nil {
			logger.Log.WithFields(logrus.Fields{
				"handler":  "Parse",
				"filename": file.Filename,
				"error":    err.Error(),
			}).Error("--- [HANDLER: Parse] Failed to read uploaded file bytes ---")
		}
		_ = ctx.Error(errors.NewInternal("Failed to read the uploaded file bytes", err))
		return
	}

	notification, err := c.service.ParsePDF(file.Filename, fileBytes)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	if logger.Log != nil {
		logger.Log.WithFields(logrus.Fields{
			"handler":         "Parse",
			"notification_id": notification.ID,
			"filename":        notification.FileName,
			"language":        notification.Language,
			"text_length":     len(notification.RawText),
		}).Info("---- [HANDLER: Parse] Successfully parsed and processed PDF notification ----")
	}

	ctx.JSON(http.StatusCreated, notification)
}

// GetByFileName handles HTTP GET requests to read/retrieve a PDF notification by file name.
func (c *ExamNotificationController) GetByFileName(ctx *gin.Context) {
	fileName := ctx.Param("filename")
	if fileName == "" {
		fileName = ctx.Query("filename")
	}

	if logger.Log != nil {
		logger.Log.WithFields(logrus.Fields{
			"handler":  "GetByFileName",
			"filename": fileName,
		}).Info("--- [HANDLER: GetByFileName] Receiving request to read notification by file name ---")
	}

	if fileName == "" {
		if logger.Log != nil {
			logger.Log.WithFields(logrus.Fields{
				"handler": "GetByFileName",
				"status":  http.StatusBadRequest,
			}).Warn("--- [HANDLER: GetByFileName] Missing required filename parameter in request ---")
		}
		_ = ctx.Error(errors.NewBadRequest("Missing required 'filename' parameter", nil))
		return
	}

	notification, err := c.service.GetByFileName(fileName)
	if err != nil {
		if logger.Log != nil {
			logger.Log.WithFields(logrus.Fields{
				"handler":  "GetByFileName",
				"filename": fileName,
				"error":    err.Error(),
			}).Warn("---- [HANDLER: GetByFileName] Failed to retrieve PDF notification ----")
		}
		_ = ctx.Error(err)
		return
	}

	if logger.Log != nil {
		logger.Log.WithFields(logrus.Fields{
			"handler":         "GetByFileName",
			"notification_id": notification.ID,
			"filename":        notification.FileName,
			"language":        notification.Language,
		}).Info("---- [HANDLER: GetByFileName] Successfully retrieved PDF notification ----")
	}

	ctx.JSON(http.StatusOK, notification)
}
