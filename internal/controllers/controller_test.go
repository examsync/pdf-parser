package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/examsync/pdf-parser/internal/services"
	appErrors "github.com/examsync/pdf-parser/utils/errors"
	"github.com/labstack/echo/v5"
)

func TestExamNotificationController_Parse_MissingFile(t *testing.T) {
	e := echo.New()
	service := services.NewExamNotificationService(nil)
	ctrl := NewExamNotificationController(service)

	req := httptest.NewRequest(http.MethodPost, "/parse", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := ctrl.Parse(c)
	if err == nil {
		t.Fatalf("expected error when file form field is missing, got nil")
	}

	var appErr *appErrors.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T: %v", err, err)
	}

	if appErr.HTTPStatus != http.StatusBadRequest {
		t.Errorf("expected HTTP status 400, got %d", appErr.HTTPStatus)
	}

	if appErr.Code != appErrors.ErrCodeBadRequest {
		t.Errorf("expected code %s, got %s", appErrors.ErrCodeBadRequest, appErr.Code)
	}
}
