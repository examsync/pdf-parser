package errors

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/examsync/pdf-parser/utils/config"
	"github.com/labstack/echo/v5"
)

func TestNewHTTPErrorHandler(t *testing.T) {
	e := echo.New()

	t.Run("AppError handling", func(t *testing.T) {
		cfg := config.ErrorConfig{
			ExposeInternalDetails: false,
			DefaultMessage:        "Default error",
		}
		handler := NewHTTPErrorHandler(cfg)

		req := httptest.NewRequest(http.MethodPost, "/parse", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		appErr := NewBadRequest("Missing file parameter", errors.New("multipart EOF"))
		handler(c, appErr)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected HTTP status 400, got %d", rec.Code)
		}

		var resp ErrorResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal JSON response: %v", err)
		}

		if resp.Success {
			t.Errorf("expected success to be false")
		}
		if resp.Error.Code != ErrCodeBadRequest {
			t.Errorf("expected error code %s, got %s", ErrCodeBadRequest, resp.Error.Code)
		}
		if resp.Error.Message != "Missing file parameter" {
			t.Errorf("expected error message 'Missing file parameter', got %s", resp.Error.Message)
		}
		if resp.Error.InternalError != "" {
			t.Errorf("expected internal_error to be omitted when ExposeInternalDetails=false, got %s", resp.Error.InternalError)
		}
	})

	t.Run("ExposeInternalDetails enabled", func(t *testing.T) {
		cfg := config.ErrorConfig{
			ExposeInternalDetails: true,
			DefaultMessage:        "Default error",
		}
		handler := NewHTTPErrorHandler(cfg)

		req := httptest.NewRequest(http.MethodPost, "/parse", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		appErr := NewInternal("Database timeout", errors.New("context deadline exceeded"))
		handler(c, appErr)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected HTTP status 500, got %d", rec.Code)
		}

		var resp ErrorResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal JSON response: %v", err)
		}

		if resp.Error.InternalError != "context deadline exceeded" {
			t.Errorf("expected internal_error 'context deadline exceeded', got %s", resp.Error.InternalError)
		}
	})

	t.Run("Standard Go error handling", func(t *testing.T) {
		cfg := config.ErrorConfig{
			ExposeInternalDetails: false,
			DefaultMessage:        "An unexpected error occurred",
		}
		handler := NewHTTPErrorHandler(cfg)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		stdErr := errors.New("raw panic or database failure")
		handler(c, stdErr)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected HTTP status 500, got %d", rec.Code)
		}

		var resp ErrorResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal JSON response: %v", err)
		}

		if resp.Error.Message != "An unexpected error occurred" {
			t.Errorf("expected default error message, got %s", resp.Error.Message)
		}
	})

	t.Run("Echo HTTPError handling", func(t *testing.T) {
		cfg := config.ErrorConfig{
			ExposeInternalDetails: false,
			DefaultMessage:        "Default error",
		}
		handler := NewHTTPErrorHandler(cfg)

		req := httptest.NewRequest(http.MethodGet, "/non-existent", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		echoErr := echo.NewHTTPError(http.StatusNotFound, "Route not found")
		handler(c, echoErr)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected HTTP status 404, got %d", rec.Code)
		}

		var resp ErrorResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal JSON response: %v", err)
		}

		if resp.Error.Code != "ERR_HTTP_404" {
			t.Errorf("expected code ERR_HTTP_404, got %s", resp.Error.Code)
		}
	})
}
