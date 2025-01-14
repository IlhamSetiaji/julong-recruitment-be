package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuestionResponse struct {
	gorm.Model `json:"-"`
	ID         uuid.UUID  `json:"id" gorm:"type:char(36);primaryKey;"`
	QuestionID uuid.UUID  `json:"question_id" gorm:"type:char(36);primaryKey;not null"`
	UserID     *uuid.UUID `json:"user_id" gorm:"type:char(36);primaryKey;not null"`
	Answer     string     `json:"answer" gorm:"type:text;default"`

	Question *Question `json:"question" gorm:"foreignKey:QuestionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (qr *QuestionResponse) BeforeCreate(tx *gorm.DB) (err error) {
	qr.ID = uuid.New()
	qr.CreatedAt = time.Now()
	qr.UpdatedAt = time.Now()
	return
}

func (qr *QuestionResponse) BeforeUpdate(tx *gorm.DB) (err error) {
	qr.UpdatedAt = time.Now()
	return
}

func (QuestionResponse) TableName() string {
	return "question_responses"
}
