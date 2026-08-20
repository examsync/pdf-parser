package services

import (
	"testing"

	"github.com/examsync/pdf-parser/utils/errors"
)

func TestExamNotificationService_GetByFileName_EmptyFileName(t *testing.T) {
	service := NewExamNotificationService(nil)
	_, err := service.GetByFileName("")
	if err == nil {
		t.Fatal("Expected error for empty filename, got nil")
	}

	appErr, ok := err.(*errors.AppError)
	if !ok {
		t.Fatalf("Expected *errors.AppError, got %T", err)
	}

	if appErr.Code != errors.ErrCodeBadRequest {
		t.Errorf("Expected error code %s, got %s", errors.ErrCodeBadRequest, appErr.Code)
	}
}

func TestExamNotificationService_GetByFileName_NotFound(t *testing.T) {
	service := NewExamNotificationService(nil)
	_, err := service.GetByFileName("non_existent_file_12345.pdf")
	if err == nil {
		t.Fatal("Expected error for non-existent file, got nil")
	}

	appErr, ok := err.(*errors.AppError)
	if !ok {
		t.Fatalf("Expected *errors.AppError, got %T", err)
	}

	if appErr.Code != errors.ErrCodeNotFound {
		t.Errorf("Expected error code %s, got %s", errors.ErrCodeNotFound, appErr.Code)
	}
}
