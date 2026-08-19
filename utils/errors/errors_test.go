package errors

import (
	"errors"
	"net/http"
	"testing"
)

func TestAppError_ErrorAndUnwrap(t *testing.T) {
	baseErr := errors.New("underlying DB error")
	appErr := NewAppError(ErrCodeInternalServer, "Database connection failed", http.StatusInternalServerError, baseErr)

	expectedStr := "[ERR_INTERNAL_SERVER_ERROR] Database connection failed: underlying DB error"
	if appErr.Error() != expectedStr {
		t.Errorf("expected error string %q, got %q", expectedStr, appErr.Error())
	}

	if !errors.Is(appErr, baseErr) {
		t.Errorf("expected errors.Is(appErr, baseErr) to be true")
	}

	unwrapped := errors.Unwrap(appErr)
	if unwrapped != baseErr {
		t.Errorf("expected unwrapped error to be baseErr, got %v", unwrapped)
	}
}

func TestHelperConstructors(t *testing.T) {
	t.Run("NewBadRequest", func(t *testing.T) {
		err := NewBadRequest("invalid form data", nil)
		if err.HTTPStatus != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, err.HTTPStatus)
		}
		if err.Code != ErrCodeBadRequest {
			t.Errorf("expected code %s, got %s", ErrCodeBadRequest, err.Code)
		}
	})

	t.Run("NewNotFound", func(t *testing.T) {
		err := NewNotFound("resource missing", nil)
		if err.HTTPStatus != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, err.HTTPStatus)
		}
		if err.Code != ErrCodeNotFound {
			t.Errorf("expected code %s, got %s", ErrCodeNotFound, err.Code)
		}
	})

	t.Run("NewUnprocessableEntity", func(t *testing.T) {
		err := NewUnprocessableEntity("corrupt pdf", nil)
		if err.HTTPStatus != http.StatusUnprocessableEntity {
			t.Errorf("expected status %d, got %d", http.StatusUnprocessableEntity, err.HTTPStatus)
		}
		if err.Code != ErrCodeUnprocessableEntity {
			t.Errorf("expected code %s, got %s", ErrCodeUnprocessableEntity, err.Code)
		}
	})

	t.Run("NewValidationError", func(t *testing.T) {
		details := map[string]string{"file": "file is required"}
		err := NewValidationError("validation failed", details)
		if err.HTTPStatus != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, err.HTTPStatus)
		}
		if err.Code != ErrCodeValidationError {
			t.Errorf("expected code %s, got %s", ErrCodeValidationError, err.Code)
		}
		if err.Details == nil {
			t.Errorf("expected details to be populated")
		}
	})

	t.Run("Wrap", func(t *testing.T) {
		if Wrap(nil, "CODE", "msg", 500) != nil {
			t.Errorf("expected Wrap(nil) to return nil")
		}
		base := errors.New("io error")
		wrapped := Wrap(base, ErrCodeInternalServer, "failed file read", http.StatusInternalServerError)
		if wrapped.Err != base {
			t.Errorf("expected wrapped inner error to be base error")
		}
	})
}
