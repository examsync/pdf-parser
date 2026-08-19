package services

import (
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
