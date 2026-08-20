package services

import (
	"os"
	"path/filepath"

	"github.com/examsync/pdf-parser/internal/models"
	"github.com/examsync/pdf-parser/internal/repositories"
	"github.com/examsync/pdf-parser/utils/errors"
	"github.com/examsync/pdf-parser/utils/pdf"
)

// ExamNotificationService handles business logic operations for exam notifications.
type ExamNotificationService struct {
	repo *repositories.ExamNotificationRepository
}

// NewExamNotificationService creates a new instance of ExamNotificationService.
func NewExamNotificationService(repo *repositories.ExamNotificationRepository) *ExamNotificationService {
	return &ExamNotificationService{repo: repo}
}

// ParsePDF parses notification data from raw PDF bytes, saves it to the database, and returns it.
func (s *ExamNotificationService) ParsePDF(fileName string, fileBytes []byte) (*models.ExamNotification, error) {
	if len(fileBytes) == 0 {
		return nil, errors.NewBadRequest("Uploaded file is empty", nil)
	}

	text, err := pdf.ExtractText(fileBytes)
	if err != nil {
		return nil, errors.NewUnprocessableEntity("Failed to extract text from PDF document", err)
	}

	lang, err := pdf.DetectLanguage(text)
	if err != nil {
		return nil, errors.NewAppError(errors.ErrCodeUnsupportedLanguage, "Unsupported document language. Only English and Hindi PDFs are supported", 422, err)
	}

	notification := pdf.ParseNotification(fileName, text, lang)

	if s.repo != nil {
		if err := s.repo.Create(notification); err != nil {
			return nil, errors.NewInternal("Failed to save notification to database", err)
		}
	}

	return notification, nil
}

// GetByFileName retrieves an existing notification record by file name from the repository,
// or attempts to locate and parse a PDF file with that name from local storage.
func (s *ExamNotificationService) GetByFileName(fileName string) (*models.ExamNotification, error) {
	if fileName == "" {
		return nil, errors.NewBadRequest("File name parameter cannot be empty", nil)
	}

	if s.repo != nil {
		notification, err := s.repo.GetByFileName(fileName)
		if err == nil && notification != nil {
			return notification, nil
		}
	}

	// Fallback: check if physical PDF file exists in workspace or uploads folder
	searchPaths := []string{
		fileName,
		filepath.Join("uploads", fileName),
	}

	for _, path := range searchPaths {
		fileBytes, err := os.ReadFile(path)
		if err == nil && len(fileBytes) > 0 {
			return s.ParsePDF(fileName, fileBytes)
		}
	}

	return nil, errors.NewNotFound("PDF file or notification record not found for file name: "+fileName, nil)
}

