package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DocumentSetup struct {
	gorm.Model      `json:"-"`
	ID              uuid.UUID              `json:"id" gorm:"type:char(36);primaryKey;"`
	DocumentTypeID  uuid.UUID              `json:"document_type_id" gorm:"type:char(36);not null"`
	RecruitmentType ProjectRecruitmentType `json:"recruitment_type" gorm:"not null"`
	Header          string                 `json:"header" gorm:"type:text;not null"`
	Body            string                 `json:"body" gorm:"type:text;not null"`
	Footer          string                 `json:"footer" gorm:"type:text;not null"`

	DocumentType     *DocumentType     `json:"document_type" gorm:"foreignKey:DocumentTypeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	DocumentSendings []DocumentSending `json:"document_sendings" gorm:"foreignKey:DocumentSetupID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (ds *DocumentSetup) BeforeCreate(tx *gorm.DB) (err error) {
	ds.ID = uuid.New()
	ds.CreatedAt = time.Now()
	ds.UpdatedAt = time.Now()
	return
}

func (ds *DocumentSetup) BeforeUpdate(tx *gorm.DB) (err error) {
	ds.UpdatedAt = time.Now()
	return
}

func (DocumentSetup) TableName() string {
	return "document_setups"
}
