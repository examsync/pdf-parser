package services

import (
	"github.com/examsync/pdf-parser/internal/models"
	"github.com/examsync/pdf-parser/internal/repositories"
)

// ExamNotificationService handles business logic operations for exam notifications.
type ExamNotificationService struct {
	repo *repositories.ExamNotificationRepository
}

// NewExamNotificationService creates a new instance of ExamNotificationService.
func NewExamNotificationService(repo *repositories.ExamNotificationRepository) *ExamNotificationService {
	return &ExamNotificationService{repo: repo}
}

// GetNotifications retrieves exam notifications from the data repository.
func (s *ExamNotificationService) GetNotifications() ([]models.ExamNotification, error) {
	return s.repo.GetAll()
}
