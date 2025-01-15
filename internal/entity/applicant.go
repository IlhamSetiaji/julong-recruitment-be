package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApplicantStatus string

const (
	APPLICANT_STATUS_APPLIED   ApplicantStatus = "APPLIED"
	APPLICANT_STATUS_SHORTLIST ApplicantStatus = "SHORTLIST"
	APPLICANT_STATUS_REJECTED  ApplicantStatus = "REJECTED"
	APPLICANT_STATUS_HIRED     ApplicantStatus = "HIRED"
)

type Applicant struct {
	gorm.Model   `json:"-"`
	ID           uuid.UUID       `json:"id" gorm:"type:char(36);primaryKey;"`
	UserID       *uuid.UUID      `json:"user_id" gorm:"type:char(36);not null"`
	JobPostingID uuid.UUID       `json:"job_posting_id" gorm:"type:char(36);not null"`
	AppliedDate  time.Time       `json:"applied_date" gorm:"type:date;not null"`
	Status       ApplicantStatus `json:"status" gorm:"not null"`

	DocumentSendings []DocumentSending `json:"document_sendings" gorm:"foreignKey:ApplicantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (a *Applicant) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.New()
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	return
}

func (a *Applicant) BeforeUpdate(tx *gorm.DB) (err error) {
	a.UpdatedAt = time.Now()
	return
}

func (Applicant) TableName() string {
	return "applicants"
}
