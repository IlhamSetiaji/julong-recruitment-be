package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Question struct {
	gorm.Model         `json:"-"`
	ID                 uuid.UUID `json:"id" gorm:"type:char(36);primaryKey;"`
	TemplateQuestionID uuid.UUID `json:"template_question_id" gorm:"type:char(36);not null"`
	AnswerTypeID       uuid.UUID `json:"answer_type_id" gorm:"type:char(36);not null"`
	Name               string    `json:"name" gorm:"type:varchar(255);not null"`
	Number             int       `json:"number" gorm:"type:int;not null"`

	TemplateQuestion  *TemplateQuestion  `json:"template_question" gorm:"foreignKey:TemplateQuestionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	AnswerType        *AnswerType        `json:"answer_type" gorm:"foreignKey:AnswerTypeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	QuestionOptions   []QuestionOption   `json:"question_options" gorm:"foreignKey:QuestionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	QuestionResponses []QuestionResponse `json:"question_responses" gorm:"foreignKey:QuestionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (q *Question) BeforeCreate(tx *gorm.DB) (err error) {
	q.ID = uuid.New()
	q.CreatedAt = time.Now()
	q.UpdatedAt = time.Now()
	return
}

func (q *Question) BeforeUpdate(tx *gorm.DB) (err error) {
	q.UpdatedAt = time.Now()
	return
}

func (Question) TableName() string {
	return "questions"
}
