package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InterviewResultStatus string

const (
	INTERVIEW_RESULT_STATUS_ACCEPTED InterviewResultStatus = "ACCEPTED"
	INTERVIEW_RESULT_STATUS_REJECTED InterviewResultStatus = "REJECTED"
)

type InterviewResult struct {
	gorm.Model           `json:"-"`
	ID                   uuid.UUID             `json:"id" gorm:"type:char(36);primaryKey;"`
	InterviewApplicantID uuid.UUID             `json:"interview_id" gorm:"type:char(36);not null"`
	InterviewAssessorID  uuid.UUID             `json:"interview_assessor_id" gorm:"type:char(36);not null"`
	Status               InterviewResultStatus `json:"status" gorm:"type:varchar(255);default:'REJECTED'"`

	InterviewApplicant *InterviewApplicant `json:"interview_applicant" gorm:"foreignKey:InterviewApplicantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	InterviewAssessor  *InterviewAssessor  `json:"interview_assessor" gorm:"foreignKey:InterviewAssessorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (ir *InterviewResult) BeforeCreate(tx *gorm.DB) (err error) {
	ir.ID = uuid.New()
	ir.CreatedAt = time.Now()
	ir.UpdatedAt = time.Now()
	return
}

func (ir *InterviewResult) BeforeUpdate(tx *gorm.DB) (err error) {
	ir.UpdatedAt = time.Now()
	return
}

func (InterviewResult) TableName() string {
	return "interview_results"
}
