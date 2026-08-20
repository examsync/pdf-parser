package errors

import (
	"errors"
	"net/http"
	"time"

	"github.com/examsync/pdf-parser/utils/config"
	"github.com/examsync/pdf-parser/utils/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// HandleError processes an error and writes the standardized ErrorResponse to gin.Context.
func HandleError(c *gin.Context, cfg config.ErrorConfig, err error) {
	if err == nil {
		return
	}

	defaultMsg := cfg.DefaultMessage
	if defaultMsg == "" {
		defaultMsg = "An unexpected internal server error occurred"
	}

	var (
		httpStatus    = http.StatusInternalServerError
		errCode       = ErrCodeInternalServer
		userMessage   = defaultMsg
		details       interface{}
		internalCause string
	)

	var appErr *AppError

	if errors.As(err, &appErr) {
		if appErr.HTTPStatus != 0 {
			httpStatus = appErr.HTTPStatus
		}
		if appErr.Code != "" {
			errCode = appErr.Code
		}
		if appErr.Message != "" {
			userMessage = appErr.Message
		}
		details = appErr.Details
		if appErr.Err != nil {
			internalCause = appErr.Err.Error()
		}
	} else if err != nil {
		internalCause = err.Error()
	}

	// Log error with context metadata
	logFields := logrus.Fields{
		"status":    httpStatus,
		"code":      errCode,
		"path":      c.Request.URL.Path,
		"method":    c.Request.Method,
		"remote_ip": c.ClientIP(),
	}
	if internalCause != "" {
		logFields["internal_error"] = internalCause
	}

	if logger.Log != nil {
		if httpStatus >= 500 {
			logger.Log.WithFields(logFields).Errorf("HTTP Handler Error: %s", userMessage)
		} else {
			logger.Log.WithFields(logFields).Warnf("HTTP Handler Warning: %s", userMessage)
		}
	}

	// Construct JSON Error Response
	resp := ErrorResponse{
		Success: false,
		Error: ErrorDetails{
			Code:      errCode,
			Message:   userMessage,
			Details:   details,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Path:      c.Request.URL.Path,
		},
	}

	if cfg.ExposeInternalDetails && internalCause != "" {
		resp.Error.InternalError = internalCause
	}

	c.JSON(httpStatus, resp)
}

// ErrorHandlerMiddleware returns a Gin middleware that handles errors attached via c.Error(err).
func ErrorHandlerMiddleware(cfg config.ErrorConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			lastErr := c.Errors.Last().Err
			HandleError(c, cfg, lastErr)
		}
	}
}
