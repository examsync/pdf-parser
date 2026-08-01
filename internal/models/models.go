package models

type ExamNotification struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	FileName string `json:"file_name" gorm:"not null"`
	RawText  string `json:"raw_text" gorm:"type:text"`
}

