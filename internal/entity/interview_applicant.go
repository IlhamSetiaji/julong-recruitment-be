package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InterviewApplicant struct {
	gorm.Model       `json:"-"`
	ID               uuid.UUID         `json:"id" gorm:"type:char(36);primaryKey;"`
	InterviewID      uuid.UUID         `json:"interview_id" gorm:"type:char(36);not null"`
	ApplicantID      uuid.UUID         `json:"applicant_id" gorm:"type:char(36);default:null"`
	UserProfileID    uuid.UUID         `json:"user_profile_id" gorm:"type:char(36);not null"`
	StartTime        time.Time         `json:"start_time" gorm:"type:time;not null"`
	EndTime          time.Time         `json:"end_time" gorm:"type:time;not null"`
	StartedTime      *time.Time        `json:"started_time" gorm:"type:time;default:null"`
	EndedTime        *time.Time        `json:"ended_time" gorm:"type:time;default:null"`
	AssessmentStatus AssessmentStatus  `json:"assessment_status" gorm:"type:text;default:null"`
	FinalResult      FinalResultStatus `json:"final_result" gorm:"type:text;default:null"`

	Interview   *Interview   `json:"interview" gorm:"foreignKey:InterviewID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserProfile *UserProfile `json:"user_profile" gorm:"foreignKey:UserProfileID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Applicant   *Applicant   `json:"applicant" gorm:"foreignKey:ApplicantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (ia *InterviewApplicant) BeforeCreate(tx *gorm.DB) (err error) {
	ia.ID = uuid.New()
	ia.CreatedAt = time.Now()
	ia.UpdatedAt = time.Now()
	return
}

func (ia *InterviewApplicant) BeforeUpdate(tx *gorm.DB) (err error) {
	ia.UpdatedAt = time.Now()
	return
}

func (InterviewApplicant) TableName() string {
	return "interview_applicants"
}
