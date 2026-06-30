package repositories

import (
	"github.com/examsync/pdf-parser/internal/models"
	"gorm.io/gorm"
)

// ParsedPDFRepository manages database operations for the ParsedPDF model.
type ParsedPDFRepository struct {
	db *gorm.DB
}

// NewParsedPDFRepository creates a new instance of ParsedPDFRepository.
func NewParsedPDFRepository(db *gorm.DB) *ParsedPDFRepository {
	return &ParsedPDFRepository{db: db}
}

// GetAll retrieves all ParsedPDF records, performing migration and seeding if necessary.
func (r *ParsedPDFRepository) GetAll() ([]models.ParsedPDF, error) {
	// Auto migrate schema if not exists
	if err := r.db.AutoMigrate(&models.ParsedPDF{}); err != nil {
		return nil, err
	}

	// Seed dummy records if empty
	var count int64
	if err := r.db.Model(&models.ParsedPDF{}).Count(&count).Error; err != nil {
		return nil, err
	}

	if count == 0 {
		dummyPDFs := []models.ParsedPDF{
			{FileName: "gobyexample.pdf"},
		}
		for _, pdf := range dummyPDFs {
			if err := r.db.Create(&pdf).Error; err != nil {
				return nil, err
			}
		}
	}

	var pdfs []models.ParsedPDF
	if err := r.db.Find(&pdfs).Error; err != nil {
		return nil, err
	}

	return pdfs, nil
}
