package repositories

import (
	"time"

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

// GetAll retrieves all ExamNotification records, performing migration and seeding if necessary.
func (r *ExamNotificationRepository) GetAll() ([]models.ExamNotification, error) {
	// Auto migrate schema if not exists
	if err := r.db.AutoMigrate(&models.ExamNotification{}); err != nil {
		return nil, err
	}

	// Seed dummy records if empty
	var count int64
	if err := r.db.Model(&models.ExamNotification{}).Count(&count).Error; err != nil {
		return nil, err
	}

	if count == 0 {
		startDate := time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2026, 7, 30, 0, 0, 0, 0, time.UTC)
		cutOffDate := time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)

		bppscNotification := models.ExamNotification{
			AdvertisementNumber:    "08/2026",
			PostName:               "Company Commander (gqYe lekns’kd)",
			Department:             "Home Department (Special Branch), Bihar (x`g foHkkx fof'k" + `es"k 'kk[kk, fcgkj)`,
			TotalVacancies:         65,
			StartDate:              startDate,
			EndDate:                endDate,
			EligibilityCutOffDate:  cutOffDate,
			AgeMin:                 20,
			AgeMaxMaleUR:           37,
			AgeMaxFemaleUR:         40,
			AgeMaxBC:               40,
			AgeMaxSCST:             42,
			EducationQualification: "Graduation or equivalent from a recognized university",
			ApplicationFee:         100.0,
			RequiredDocuments:      "Active Mobile Number; Active Email ID; Class 10 (Matric) Board Certificate (for DOB, Name, Parentage verification); Graduation Degree/Marksheet; Photograph (Color, white background, 15-25 KB, .jpg/.jpeg/.gif format); Signature (Hindi & English separately in black/blue ink, white background, 15-25 KB, .jpg/.jpeg/.gif format); Caste/Category Certificate (if applicable); Domicile/Residence Certificate (if applicable); EWS Income & Asset Certificate (if applicable); Ex-servicemen Discharge Book/Defense Service Certificate (if applicable); Regular Govt Service Certificate (if applicable)",
			EligibilityCriteria:    "Nationality: Indian; Gender: Male, Female, and Third Gender; Physical Standards: Height (Male UR/BC: min 165 cm, Male EBC/SC/ST: min 160 cm, Female all categories: min 155 cm); Chest (Male UR/BC/EBC: min 81 cm unexpanded / 86 cm expanded, Male SC/ST: min 79 cm unexpanded / 84 cm expanded, min 5cm expansion required); Weight (Female all categories: min 48 kg); Third Gender Physical Standards: Same as BC category",
		}

		if err := r.db.Create(&bppscNotification).Error; err != nil {
			return nil, err
		}
	}

	var notifications []models.ExamNotification
	if err := r.db.Find(&notifications).Error; err != nil {
		return nil, err
	}

	return notifications, nil
}
