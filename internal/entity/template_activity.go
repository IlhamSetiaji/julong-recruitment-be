package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TemplateActivityStatus string

const (
	TEMPLATE_ACTIVITY_STATUS_ACTIVE   TemplateActivityStatus = "ACTIVE"
	TEMPLATE_ACTIVITY_STATUS_INACTIVE TemplateActivityStatus = "INACTIVE"
)

type TemplateActivity struct {
	gorm.Model      `json:"-"`
	ID              uuid.UUID              `json:"id" gorm:"type:char(36);primaryKey;"`
	Name            string                 `json:"name" gorm:"type:varchar(255);not null"`
	Description     string                 `json:"description" gorm:"type:text;default:null"`
	RecruitmentType ProjectRecruitmentType `json:"recruitment_type" gorm:"not null"`
	Status          TemplateActivityStatus `json:"status" gorm:"default:'ACTIVE'"`

	TemplateActivityLines []TemplateActivityLine `json:"template_activity_lines" gorm:"foreignKey:TemplateActivityID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (ta *TemplateActivity) BeforeCreate(tx *gorm.DB) (err error) {
	ta.ID = uuid.New()
	ta.CreatedAt = time.Now()
	ta.UpdatedAt = time.Now()
	return nil
}

func (ta *TemplateActivity) BeforeUpdate(tx *gorm.DB) (err error) {
	ta.UpdatedAt = time.Now()
	return nil
}

func (TemplateActivity) TableName() string {
	return "template_activities"
}
