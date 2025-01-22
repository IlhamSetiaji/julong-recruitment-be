package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectPic struct {
	gorm.Model               `json:"-"`
	ID                       uuid.UUID  `json:"id" gorm:"type:char(36);primaryKey;"`
	ProjectRecruitmentLineID uuid.UUID  `json:"project_recruitment_line_id" gorm:"type:char(36);not null"`
	EmployeeID               *uuid.UUID `json:"employee_id" gorm:"type:char(36);not null"`
	AdministrativeTotal      int        `json:"administrative_total" gorm:"default:null"`

	ProjectRecruitmentLine *ProjectRecruitmentLine `json:"project_recruitment_line" gorm:"foreignKey:ProjectRecruitmentLineID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TestScheduleHeaders    []TestScheduleHeader    `json:"test_schedule_headers" gorm:"foreignKey:ProjectPicID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (pp *ProjectPic) BeforeCreate(tx *gorm.DB) (err error) {
	pp.ID = uuid.New()
	pp.CreatedAt = time.Now()
	pp.UpdatedAt = time.Now()
	return
}

func (pp *ProjectPic) BeforeUpdate(tx *gorm.DB) (err error) {
	pp.UpdatedAt = time.Now()
	return
}

func (ProjectPic) TableName() string {
	return "project_pics"
}
