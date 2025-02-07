package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InterviewAssessor struct {
	gorm.Model  `json:"-"`
	ID          uuid.UUID  `json:"id" gorm:"type:char(36);primaryKey;"`
	InterviewID uuid.UUID  `json:"interview_id" gorm:"type:char(36);not null"`
	EmployeeID  *uuid.UUID `json:"employee_id" gorm:"type:char(36);not null"`

	Interview         *Interview         `json:"interview" gorm:"foreignKey:InterviewID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	QuestionResponses []QuestionResponse `json:"question_responses" gorm:"foreignKey:InterviewAssessorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	InterviewResults  []InterviewResult  `json:"interview_results" gorm:"foreignKey:InterviewAssessorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (ia *InterviewAssessor) BeforeCreate(tx *gorm.DB) (err error) {
	ia.ID = uuid.New()
	ia.CreatedAt = time.Now()
	ia.UpdatedAt = time.Now()
	return
}

func (ia *InterviewAssessor) BeforeUpdate(tx *gorm.DB) (err error) {
	ia.UpdatedAt = time.Now()
	return
}

func (InterviewAssessor) TableName() string {
	return "interview_assessors"
}
