package models

type ExamNotification struct {
	ID                  uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	FileName            string `json:"file_name" gorm:"not null"`
	ImportantDates      string `json:"important_dates" gorm:"type:text"`
	EligibilityCriteria string `json:"eligibility_criteria" gorm:"type:text"`
	RequiredDocuments   string `json:"required_documents" gorm:"type:text"`
	Fee                 string `json:"fee" gorm:"type:text"`
}
