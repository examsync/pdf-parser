package repositories

import (
	"github.com/examsync/pdf-parser/internal/models"
	"gorm.io/gorm"
)

// ExamNotificationRepository manages database operations for the ExamNotification model.
type ExamNotificationRepository struct {
	db *gorm.DB
}

// NewExamNotificationRepository creates a new instance of ExamNotificationRepository.
func NewExamNotificationRepository(db *gorm.DB) *ExamNotificationRepository {
	return &ExamNotificationRepository{db: db}
}

// Create inserts a new ExamNotification record into the database.
func (r *ExamNotificationRepository) Create(notification *models.ExamNotification) error {
	return r.db.Create(notification).Error
}

// GetAll retrieves all ExamNotification records from the database.
func (r *ExamNotificationRepository) GetAll() ([]models.ExamNotification, error) {
	var notifications []models.ExamNotification
	if err := r.db.Find(&notifications).Error; err != nil {
		return nil, err
	}

	return notifications, nil
}
