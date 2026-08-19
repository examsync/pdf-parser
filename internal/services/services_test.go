package services

import (
	"errors"
	"testing"

	appErrors "github.com/examsync/pdf-parser/utils/errors"
)

func TestExamNotificationService_ParsePDF_EmptyFile(t *testing.T) {
	service := NewExamNotificationService(nil)
	_, err := service.ParsePDF("empty.pdf", []byte{})

	if err == nil {
		t.Fatalf("expected error for empty file, got nil")
	}

	var appErr *appErrors.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T: %v", err, err)
	}

	if appErr.Code != appErrors.ErrCodeBadRequest {
		t.Errorf("expected code %s, got %s", appErrors.ErrCodeBadRequest, appErr.Code)
	}
}

func TestExamNotificationService_ParsePDF_CorruptPDF(t *testing.T) {
	service := NewExamNotificationService(nil)
	_, err := service.ParsePDF("corrupt.pdf", []byte("invalid pdf content"))

	if err == nil {
		t.Fatalf("expected error for corrupt PDF bytes, got nil")
	}

	var appErr *appErrors.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T: %v", err, err)
	}

	if appErr.Code != appErrors.ErrCodeUnprocessableEntity {
		t.Errorf("expected code %s, got %s", appErrors.ErrCodeUnprocessableEntity, appErr.Code)
	}
}
