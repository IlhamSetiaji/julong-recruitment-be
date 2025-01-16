package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TemplateActivityLineStatus string

const (
	TEMPLATE_ACTIVITY_LINE_STATUS_ACTIVE   TemplateActivityLineStatus = "ACTIVE"
	TEMPLATE_ACTIVITY_LINE_STATUS_INACTIVE TemplateActivityLineStatus = "INACTIVE"
)

type TemplateActivityLine struct {
	gorm.Model         `json:"-"`
	ID                 uuid.UUID                  `json:"id" gorm:"type:char(36);primaryKey;"`
	TemplateActivityID uuid.UUID                  `json:"template_activity_id" gorm:"type:char(36);not null"`
	Description        string                     `json:"description" gorm:"type:text;default:null"`
	Status             TemplateActivityLineStatus `json:"status" gorm:"default:'ACTIVE'"`
	QuestionTemplateID uuid.UUID                  `json:"question_template_id" gorm:"type:char(36);not null"`
	ColorHexCode       string                     `json:"color_hex_code" gorm:"type:varchar(10);not null"`

	TemplateActivity *TemplateActivity `json:"template_activity" gorm:"foreignKey:TemplateActivityID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TemplateQuestion *TemplateQuestion `json:"template_question" gorm:"foreignKey:QuestionTemplateID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (tal *TemplateActivityLine) BeforeCreate(tx *gorm.DB) (err error) {
	tal.ID = uuid.New()
	tal.CreatedAt = time.Now()
	tal.UpdatedAt = time.Now()
	return nil
}

func (tal *TemplateActivityLine) BeforeUpdate(tx *gorm.DB) (err error) {
	tal.UpdatedAt = time.Now()
	return nil
}

func (TemplateActivityLine) TableName() string {
	return "template_activity_lines"
}
