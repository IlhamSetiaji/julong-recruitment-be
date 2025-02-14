package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DocumentVerificationLine struct {
	gorm.Model                   `json:"-"`
	ID                           uuid.UUID `json:"id" gorm:"type:char(36);primaryKey;"`
	DocumentVerificationHeaderID uuid.UUID `json:"document_verification_header_id" gorm:"type:char(36);not null"`
	DocumentVerificationID       uuid.UUID `json:"document_verification_id" gorm:"type:char(36);not null"`
	Path                         string    `json:"path" gorm:"type:text;not null"`

	DocumentVerificationHeader *DocumentVerificationHeader `json:"document_verification_header" gorm:"foreignKey:DocumentVerificationHeaderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	DocumentVerification       *DocumentVerification       `json:"document_verification" gorm:"foreignKey:DocumentVerificationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (dvl *DocumentVerificationLine) BeforeCreate(tx *gorm.DB) (err error) {
	dvl.ID = uuid.New()
	dvl.CreatedAt = time.Now()
	dvl.UpdatedAt = time.Now()
	return
}

func (dvl *DocumentVerificationLine) BeforeUpdate(tx *gorm.DB) (err error) {
	dvl.UpdatedAt = time.Now()
	return
}

func (DocumentVerificationLine) TableName() string {
	return "document_verification_lines"
}
