package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Question struct {
	gorm.Model         `json:"-"`
	ID                 uuid.UUID `json:"id" gorm:"type:char(36);primaryKey;"`
	TemplateQuestionID uuid.UUID `json:"template_question_id" gorm:"type:char(36);not null"`
	AnswerTypeID       uuid.UUID `json:"answer_type_id" gorm:"type:char(36);not null"`
	Name               string    `json:"name" gorm:"type:varchar(255);not null"`

	TemplateQuestion *TemplateQuestion `json:"template_question" gorm:"foreignKey:TemplateQuestionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	AnswerType       *AnswerType       `json:"answer_type" gorm:"foreignKey:AnswerTypeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
