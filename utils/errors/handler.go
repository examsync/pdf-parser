package errors

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/examsync/pdf-parser/utils/config"
	"github.com/examsync/pdf-parser/utils/logger"
	"github.com/labstack/echo/v5"
	"github.com/sirupsen/logrus"
)

// NewHTTPErrorHandler returns an Echo HTTPErrorHandler configured with ErrorConfig rules.
func NewHTTPErrorHandler(cfg config.ErrorConfig) echo.HTTPErrorHandler {
	defaultMsg := cfg.DefaultMessage
	if defaultMsg == "" {
		defaultMsg = "An unexpected internal server error occurred"
	}

	return func(c *echo.Context, err error) {
		if err == nil {
			return
		}

		var (
			httpStatus    = http.StatusInternalServerError
			errCode       = ErrCodeInternalServer
			userMessage   = defaultMsg
			details       interface{}
			internalCause string
		)

		var appErr *AppError
		var echoErr *echo.HTTPError

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
		} else if errors.As(err, &echoErr) {
			if echoErr.Code != 0 {
				httpStatus = echoErr.Code
			}
			errCode = fmt.Sprintf("ERR_HTTP_%d", httpStatus)
			userMessage = echoErr.Message
			if echoErr.Unwrap() != nil {
				internalCause = echoErr.Unwrap().Error()
			}
		} else if err != nil {
			internalCause = err.Error()
		}

		// Log error with context metadata
		logFields := logrus.Fields{
			"status":    httpStatus,
			"code":      errCode,
			"path":      c.Request().URL.Path,
			"method":    c.Request().Method,
			"remote_ip": c.RealIP(),
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
				Path:      c.Request().URL.Path,
			},
		}

		if cfg.ExposeInternalDetails && internalCause != "" {
			resp.Error.InternalError = internalCause
		}

		_ = c.JSON(httpStatus, resp)
	}
}
