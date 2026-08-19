package errors

import (
	"fmt"
	"net/http"
)

// Standard Error Codes
const (
	ErrCodeBadRequest          = "ERR_BAD_REQUEST"
	ErrCodeNotFound            = "ERR_NOT_FOUND"
	ErrCodeUnprocessableEntity = "ERR_UNPROCESSABLE_ENTITY"
	ErrCodeInternalServer      = "ERR_INTERNAL_SERVER_ERROR"
	ErrCodeValidationError     = "ERR_VALIDATION_ERROR"
	ErrCodeUnauthorized        = "ERR_UNAUTHORIZED"
	ErrCodeForbidden           = "ERR_FORBIDDEN"
	ErrCodeUnsupportedLanguage = "ERR_UNSUPPORTED_LANGUAGE"
)

// AppError defines a structured application error with HTTP metadata.
type AppError struct {
	Code       string      `json:"code"`
	Message    string      `json:"message"`
	HTTPStatus int         `json:"-"`
	Err        error       `json:"-"`
	Details    interface{} `json:"details,omitempty"`
}

// Error implements the standard Go error interface.
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap implements error unwrapping for Go 1.13+ errors.Unwrap / errors.Is / errors.As.
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a custom AppError.
func NewAppError(code string, message string, httpStatus int, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Err:        err,
	}
}

// NewBadRequest returns a 400 Bad Request error.
func NewBadRequest(message string, err error) *AppError {
	return NewAppError(ErrCodeBadRequest, message, http.StatusBadRequest, err)
}

// NewNotFound returns a 404 Not Found error.
func NewNotFound(message string, err error) *AppError {
	return NewAppError(ErrCodeNotFound, message, http.StatusNotFound, err)
}

// NewUnprocessableEntity returns a 422 Unprocessable Entity error.
func NewUnprocessableEntity(message string, err error) *AppError {
	return NewAppError(ErrCodeUnprocessableEntity, message, http.StatusUnprocessableEntity, err)
}

// NewInternal returns a 500 Internal Server Error.
func NewInternal(message string, err error) *AppError {
	return NewAppError(ErrCodeInternalServer, message, http.StatusInternalServerError, err)
}

// NewValidationError returns a 400 Validation Error with additional detail field metadata.
func NewValidationError(message string, details interface{}) *AppError {
	return &AppError{
		Code:       ErrCodeValidationError,
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
		Details:    details,
	}
}

// Wrap wraps an existing error into an AppError with specified parameters.
func Wrap(err error, code string, message string, httpStatus int) *AppError {
	if err == nil {
		return nil
	}
	return NewAppError(code, message, httpStatus, err)
}
