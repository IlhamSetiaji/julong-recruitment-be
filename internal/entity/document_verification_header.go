package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DocumentVerificationHeaderStatus string

const (
	DOCUMENT_VERIFICATION_HEADER_STATUS_PENDING  DocumentVerificationHeaderStatus = "PENDING"
	DOCUMENT_VERIFICATION_HEADER_STATUS_VERIFIED DocumentVerificationHeaderStatus = "ACCEPTED"
	DOCUMENT_VERIFICATION_HEADER_STATUS_REJECTED DocumentVerificationHeaderStatus = "REJECTED"
)

type DocumentVerificationHeader struct {
	gorm.Model               `json:"-"`
	ID                       uuid.UUID                        `json:"id" gorm:"type:char(36);primaryKey;"`
	ProjectRecruitmentLineID uuid.UUID                        `json:"project_recruitment_line_id" gorm:"type:char(36);not null"`
	ApplicantID              uuid.UUID                        `json:"applicant_id" gorm:"type:char(36);not null"`
	JobPostingID             uuid.UUID                        `json:"job_posting_id" gorm:"type:char(36);not null"`
	VerifiedBy               uuid.UUID                        `json:"verified_by" gorm:"type:char(36);defaut:null"`
	Status                   DocumentVerificationHeaderStatus `json:"status" gorm:"type:varchar(255);default:'PENDING'"`

	ProjectRecruitmentLine    *ProjectRecruitmentLine    `json:"project_recruitment_line" gorm:"foreignKey:ProjectRecruitmentLineID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Applicant                 *Applicant                 `json:"applicant" gorm:"foreignKey:ApplicantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	JobPosting                *JobPosting                `json:"job_posting" gorm:"foreignKey:JobPostingID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	DocumentVerificationLines []DocumentVerificationLine `json:"document_verification_lines" gorm:"foreignKey:DocumentVerificationHeaderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (dv *DocumentVerificationHeader) BeforeCreate(tx *gorm.DB) (err error) {
	dv.ID = uuid.New()
	dv.CreatedAt = time.Now()
	dv.UpdatedAt = time.Now()
	return
}

func (dv *DocumentVerificationHeader) BeforeUpdate(tx *gorm.DB) (err error) {
	dv.UpdatedAt = time.Now()
	return
}

func (DocumentVerificationHeader) TableName() string {
	return "document_verification_headers"
}
