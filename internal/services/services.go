package services

import (
	"github.com/examsync/pdf-parser/internal/models"
	"github.com/examsync/pdf-parser/internal/repositories"
)

// ParsedPDFService handles business logic operations for parsed PDFs.
type ParsedPDFService struct {
	repo *repositories.ParsedPDFRepository
}

// NewParsedPDFService creates a new instance of ParsedPDFService.
func NewParsedPDFService(repo *repositories.ParsedPDFRepository) *ParsedPDFService {
	return &ParsedPDFService{repo: repo}
}

// GetPDFs retrieves parsed PDFs from the data repository.
func (s *ParsedPDFService) GetPDFs() ([]models.ParsedPDF, error) {
	return s.repo.GetAll()
}
