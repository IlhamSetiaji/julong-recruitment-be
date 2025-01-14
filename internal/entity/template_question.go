package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TemplateQuestionStatus string

const (
	TEMPLATE_QUESTION_STATUS_ACTIVE   TemplateQuestionStatus = "ACTIVE"
	TEMPLATE_QUESTION_STATUS_INACTIVE TemplateQuestionStatus = "INACTIVE"
)

type TemplateQuestion struct {
	gorm.Model      `json:"-"`
	ID              uuid.UUID              `json:"id" gorm:"type:char(36);primaryKey;"`
	DocumentSetupID uuid.UUID              `json:"document_setup_id" gorm:"type:char(36);not null"`
	FormType        string                 `json:"form_type" gorm:"type:varchar(255);not null"`
	Description     string                 `json:"description" gorm:"type:text;default:null"`
	Duration        int                    `json:"duration" gorm:"not null"`
	Status          TemplateQuestionStatus `json:"status" gorm:"default:'ACTIVE'"`

	Questions []Question `json:"questions" gorm:"foreignKey:TemplateQuestionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (tq *TemplateQuestion) BeforeCreate(tx *gorm.DB) (err error) {
	tq.ID = uuid.New()
	tq.CreatedAt = time.Now()
	tq.UpdatedAt = time.Now()
	return
}

func (tq *TemplateQuestion) BeforeUpdate(tx *gorm.DB) (err error) {
	tq.UpdatedAt = time.Now()
	return
}

func (TemplateQuestion) TableName() string {
	return "template_questions"
}
