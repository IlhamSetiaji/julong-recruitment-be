package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DocumentAgreementStatus string

const (
	DOCUMENT_AGREEMENT_STATUS_SUBMITTED DocumentAgreementStatus = "SUBMITTED"
	DOCUMENT_AGREEMENT_STATUS_APPROVED  DocumentAgreementStatus = "APPROVED"
	DOCUMENT_AGREEMENT_STATUS_REJECTED  DocumentAgreementStatus = "REJECTED"
	DOCUMENT_AGREEMENT_STATUS_REVISED   DocumentAgreementStatus = "REVISED"
	DOCUMENT_AGREEMENT_STATUS_COMPLETED DocumentAgreementStatus = "COMPLETED"
)

type DocumentAgreement struct {
	gorm.Model
	ID                uuid.UUID               `json:"id" gorm:"type:char(36);primaryKey;"`
	DocumentSendingID uuid.UUID               `json:"document_sending_id" gorm:"type:char(36);not null"`
	ApplicantID       uuid.UUID               `json:"applicant_id" gorm:"type:char(36);not null"`
	Status            DocumentAgreementStatus `json:"status" gorm:"default:'SUBMITTED'"`
	Path              string                  `json:"path" gorm:"type:text;not null"`

	DocumentSending *DocumentSending `json:"document_sending" gorm:"foreignKey:DocumentSendingID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Applicant       *Applicant       `json:"applicant" gorm:"foreignKey:ApplicantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (da *DocumentAgreement) BeforeCreate(tx *gorm.DB) (err error) {
	da.ID = uuid.New()
	da.CreatedAt = time.Now()
	da.UpdatedAt = time.Now()
	return
}

func (da *DocumentAgreement) BeforeUpdate(tx *gorm.DB) (err error) {
	da.UpdatedAt = time.Now()
	return
}

func (DocumentAgreement) TableName() string {
	return "document_agreements"
}
