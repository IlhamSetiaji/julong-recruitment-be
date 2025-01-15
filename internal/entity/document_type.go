package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DocumentType struct {
	gorm.Model `json:"-"`
	ID         uuid.UUID `json:"id" gorm:"type:char(36);primaryKey;"`
	Name       string    `json:"name" gorm:"type:varchar(255);not null"`

	MailTemplates  []MailTemplate  `json:"mail_templates" gorm:"foreignKey:DocumentTypeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	DocumentSetups []DocumentSetup `json:"document_setups" gorm:"foreignKey:DocumentTypeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (dt *DocumentType) BeforeCreate(tx *gorm.DB) (err error) {
	dt.ID = uuid.New()
	dt.CreatedAt = time.Now()
	dt.UpdatedAt = time.Now()
	return
}

func (dt *DocumentType) BeforeUpdate(tx *gorm.DB) (err error) {
	dt.UpdatedAt = time.Now()
	return
}

func (DocumentType) TableName() string {
	return "document_types"
}
