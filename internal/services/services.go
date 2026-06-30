package services

import (
	"github.com/examsync/pdf-parser/internal/models"
	"github.com/examsync/pdf-parser/internal/repositories"
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
	text, err := pdf.ExtractText(fileBytes)
	if err != nil {
		return nil, err
	}

	notification := pdf.ParseNotification(fileName, text)

	if err := s.repo.Create(notification); err != nil {
		return nil, err
	}

	return notification, nil
}
