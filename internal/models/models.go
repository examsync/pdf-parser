package models

// ParsedPDF represents a GORM database entity for parsed PDFs.
type ParsedPDF struct {
	FileName string `json:"file_name" gorm:"not null"`
}
