package models

import "time"

type ExamNotification struct {
	ID                     uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	AdvertisementNumber    string    `json:"advertisement_number" gorm:"not null;unique"` // e.g., "08/2026"
	PostName               string    `json:"post_name" gorm:"not null"`                  // e.g., "Company Commander"
	Department             string    `json:"department"`                                 // e.g., "Home Department (Special Branch), Bihar"
	TotalVacancies         int       `json:"total_vacancies"`                            // e.g., 65
	StartDate              time.Time `json:"start_date"`
	EndDate                time.Time `json:"end_date"`
	EligibilityCutOffDate  time.Time `json:"eligibility_cut_off_date"`
	AgeMin                 int       `json:"age_min"`
	AgeMaxMaleUR           int       `json:"age_max_male_ur"`
	AgeMaxFemaleUR         int       `json:"age_max_female_ur"`
	AgeMaxBC               int       `json:"age_max_bc"`
	AgeMaxSCST             int       `json:"age_max_scst"`
	EducationQualification string    `json:"education_qualification"`
	ApplicationFee         float64   `json:"application_fee"`
	RequiredDocuments      string    `json:"required_documents" gorm:"type:text"`
	EligibilityCriteria    string    `json:"eligibility_criteria" gorm:"type:text"`
}
