package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MailTemplate struct {
	gorm.Model     `json:"-"`
	ID             uuid.UUID                `json:"id" gorm:"type:char(36);primaryKey;"`
	Name           string                   `json:"name" gorm:"type:varchar(255);not null"`
	DocumentTypeID uuid.UUID                `json:"document_type_id" gorm:"type:char(36);default:null"`
	FormType       TemplateQuestionFormType `json:"form_type" gorm:"type:varchar(255);not null"`
	Subject        string                   `json:"subject" gorm:"type:text;not null"`
	Body           string                   `json:"body" gorm:"type:text;not null"`

	DocumentType *DocumentType `json:"document_type" gorm:"foreignKey:DocumentTypeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (mt *MailTemplate) BeforeCreate(tx *gorm.DB) (err error) {
	mt.ID = uuid.New()
	mt.CreatedAt = time.Now()
	mt.UpdatedAt = time.Now()
	return
}

func (mt *MailTemplate) BeforeUpdate(tx *gorm.DB) (err error) {
	mt.UpdatedAt = time.Now()
	return
}

func (MailTemplate) TableName() string {
	return "mail_templates"
}
