package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuestionOption struct {
	gorm.Model `json:"-"`
	ID         uuid.UUID `json:"id" gorm:"type:char(36);primaryKey;"`
	QuestionID uuid.UUID `json:"question_id" gorm:"type:char(36);not null"`
	OptionText string    `json:"option_text" gorm:"type:text;not null"`

	Question *Question `json:"question" gorm:"foreignKey:QuestionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (qo *QuestionOption) BeforeCreate(tx *gorm.DB) (err error) {
	qo.ID = uuid.New()
	qo.CreatedAt = time.Now()
	qo.UpdatedAt = time.Now()
	return
}

func (qo *QuestionOption) BeforeUpdate(tx *gorm.DB) (err error) {
	qo.UpdatedAt = time.Now()
	return
}

func (QuestionOption) TableName() string {
	return "question_options"
}
