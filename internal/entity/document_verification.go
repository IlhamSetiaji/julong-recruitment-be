package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DocumentVerification struct {
	gorm.Model         `json:"-"`
	ID                 uuid.UUID `json:"id" gorm:"type:char(36);primaryKey;"`
	TemplateQuestionID uuid.UUID `json:"template_question_id" gorm:"type:char(36);not null"`
	Name               string    `json:"name" gorm:"type:text;not null"`
	Format             string    `json:"format" gorm:"type:text;not null"`

	TemplateQuestion *TemplateQuestion `json:"template_question" gorm:"foreignKey:TemplateQuestionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (dv *DocumentVerification) BeforeCreate(tx *gorm.DB) (err error) {
	dv.ID = uuid.New()
	dv.CreatedAt = time.Now()
	dv.UpdatedAt = time.Now()
	return
}

func (dv *DocumentVerification) BeforeUpdate(tx *gorm.DB) (err error) {
	dv.UpdatedAt = time.Now()
	return
}

func (DocumentVerification) TableName() string {
	return "document_verifications"
}
