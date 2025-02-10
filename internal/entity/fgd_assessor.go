package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FgdAssessor struct {
	gorm.Model    `json:"-"`
	ID            uuid.UUID  `json:"id" gorm:"type:char(36);primaryKey;"`
	FgdScheduleID uuid.UUID  `json:"fgd_schedule_id" gorm:"type:char(36);not null"`
	EmployeeID    *uuid.UUID `json:"employee_id" gorm:"type:char(36);not null"`

	FgdSchedule       *FgdSchedule       `json:"fgd_schedule" gorm:"foreignKey:FgdScheduleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	QuestionResponses []QuestionResponse `json:"question_responses" gorm:"foreignKey:FgdAssessorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	FgdResults        []FgdResult        `json:"fgd_results" gorm:"foreignKey:FgdAssessorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (fa *FgdAssessor) BeforeCreate(tx *gorm.DB) (err error) {
	fa.ID = uuid.New()
	fa.CreatedAt = time.Now()
	fa.UpdatedAt = time.Now()
	return
}

func (fa *FgdAssessor) BeforeUpdate(tx *gorm.DB) (err error) {
	fa.UpdatedAt = fa.UpdatedAt
	return
}

func (FgdAssessor) TableName() string {
	return "fgd_assessors"
}
