package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/examsync/pdf-parser/internal/services"
	"github.com/examsync/pdf-parser/utils/config"
	"github.com/examsync/pdf-parser/utils/errors"
	"github.com/gin-gonic/gin"
)

func TestExamNotificationController_GetByFileName_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(errors.ErrorHandlerMiddleware(config.ErrorConfig{}))

	service := services.NewExamNotificationService(nil)
	controller := NewExamNotificationController(service)

	r.GET("/notifications/:filename", controller.GetByFileName)

	req := httptest.NewRequest(http.MethodGet, "/notifications/non_existent.pdf", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code == http.StatusOK {
		t.Fatalf("Expected non-200 status code for missing PDF file, got %d", rec.Code)
	}
}

func TestExamNotificationController_GetByFileName_MissingParam(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(errors.ErrorHandlerMiddleware(config.ErrorConfig{}))

	service := services.NewExamNotificationService(nil)
	controller := NewExamNotificationController(service)

	r.GET("/notifications", controller.GetByFileName)

	req := httptest.NewRequest(http.MethodGet, "/notifications", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code == http.StatusOK {
		t.Fatalf("Expected non-200 status code when filename parameter is missing, got %d", rec.Code)
	}
}
