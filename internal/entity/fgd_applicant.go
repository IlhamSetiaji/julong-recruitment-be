package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FgdApplicant struct {
	gorm.Model       `json:"-"`
	ID               uuid.UUID         `json:"id" gorm:"type:char(36);primaryKey;"`
	FgdScheduleID    uuid.UUID         `json:"fgd_schedule_id" gorm:"type:char(36);not null"`
	ApplicantID      uuid.UUID         `json:"applicant_id" gorm:"type:char(36);default:null"`
	UserProfileID    uuid.UUID         `json:"user_profile_id" gorm:"type:char(36);not null"`
	StartTime        time.Time         `json:"start_time" gorm:"type:time;not null"`
	EndTime          time.Time         `json:"end_time" gorm:"type:time;not null"`
	StartedTime      *time.Time        `json:"started_time" gorm:"type:time;default:null"`
	EndedTime        *time.Time        `json:"ended_time" gorm:"type:time;default:null"`
	AssessmentStatus AssessmentStatus  `json:"assessment_status" gorm:"type:text;default:null"`
	FinalResult      FinalResultStatus `json:"final_result" gorm:"type:text;default:null"`

	FgdSchedule *FgdSchedule `json:"fgd_schedule" gorm:"foreignKey:FgdScheduleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserProfile *UserProfile `json:"user_profile" gorm:"foreignKey:UserProfileID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Applicant   *Applicant   `json:"applicant" gorm:"foreignKey:ApplicantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	FgdResults  []FgdResult  `json:"fgd_results" gorm:"foreignKey:FgdApplicantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (fa *FgdApplicant) BeforeCreate(tx *gorm.DB) (err error) {
	fa.ID = uuid.New()
	fa.CreatedAt = time.Now()
	fa.UpdatedAt = time.Now()
	return
}

func (fa *FgdApplicant) BeforeUpdate(tx *gorm.DB) (err error) {
	fa.UpdatedAt = time.Now()
	return
}

func (FgdApplicant) TableName() string {
	return "fgd_applicants"
}
