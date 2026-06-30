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

// Create inserts a new ExamNotification record into the database, running auto-migration first.
func (r *ExamNotificationRepository) Create(notification *models.ExamNotification) error {
	if err := r.db.AutoMigrate(&models.ExamNotification{}); err != nil {
		return err
	}
	return r.db.Create(notification).Error
}

// GetAll retrieves all ExamNotification records, performing migration if necessary.
func (r *ExamNotificationRepository) GetAll() ([]models.ExamNotification, error) {
	// Auto migrate schema if not exists
	if err := r.db.AutoMigrate(&models.ExamNotification{}); err != nil {
		return nil, err
	}

	var notifications []models.ExamNotification
	if err := r.db.Find(&notifications).Error; err != nil {
		return nil, err
	}

	return notifications, nil
}
