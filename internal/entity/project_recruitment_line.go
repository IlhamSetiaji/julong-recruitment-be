package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectRecruitmentLine struct {
	gorm.Model                 `json:"-"`
	ID                         uuid.UUID `json:"id" gorm:"type:char(36);primaryKey;"`
	ProjectRecruitmentHeaderID uuid.UUID `json:"project_recruitment_header_id" gorm:"type:char(36);not null"`
	TemplateActivityID         uuid.UUID `json:"template_activity_id" gorm:"type:char(36);not null"`
	StartDate                  time.Time `json:"start_date" gorm:"type:date;not null"`
	EndDate                    time.Time `json:"end_date" gorm:"type:date;not null"`

	ProjectRecruitmentHeader *ProjectRecruitmentHeader `json:"project_recruitment_header" gorm:"foreignKey:ProjectRecruitmentHeaderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TemplateActivity         *TemplateActivity         `json:"template_activity" gorm:"foreignKey:TemplateActivityID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (prl *ProjectRecruitmentLine) BeforeCreate(tx *gorm.DB) (err error) {
	prl.ID = uuid.New()
	prl.CreatedAt = time.Now()
	prl.UpdatedAt = time.Now()
	return nil
}

func (prl *ProjectRecruitmentLine) BeforeUpdate(tx *gorm.DB) (err error) {
	prl.UpdatedAt = time.Now()
	return nil
}

func (ProjectRecruitmentLine) TableName() string {
	return "project_recruitment_lines"
}
