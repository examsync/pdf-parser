package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/examsync/pdf-parser/internal/services"
	"github.com/labstack/echo/v5"
)

func TestExamNotificationController_GetByFileName_NotFound(t *testing.T) {
	e := echo.New()
	service := services.NewExamNotificationService(nil)
	controller := NewExamNotificationController(service)

	e.GET("/pdf/:filename", controller.GetByFileName)

	req := httptest.NewRequest(http.MethodGet, "/pdf/non_existent.pdf", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusOK {
		t.Fatalf("Expected non-200 status code for missing PDF file, got %d", rec.Code)
	}
}

func TestExamNotificationController_GetByFileName_MissingParam(t *testing.T) {
	e := echo.New()
	service := services.NewExamNotificationService(nil)
	controller := NewExamNotificationController(service)

	e.GET("/pdf", controller.GetByFileName)

	req := httptest.NewRequest(http.MethodGet, "/pdf", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusOK {
		t.Fatalf("Expected non-200 status code when filename parameter is missing, got %d", rec.Code)
	}
}
