package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FinalResultStatus string

const (
	FINAL_RESULT_STATUS_DRAFT       FinalResultStatus = "DRAFT"
	FINAL_RESULT_STATUS_IN_PROGRESS FinalResultStatus = "IN PROGRESS"
	FINAL_RESULT_STATUS_COMPLETED   FinalResultStatus = "COMPLETED"
	FINAL_RESULT_STATUS_ACCEPTED    FinalResultStatus = "ACCEPTED"
	FINAL_RESULT_STATUS_REJECTED    FinalResultStatus = "REJECTED"
)

type AssessmentStatus string

const (
	ASSESSMENT_STATUS_DRAFT       AssessmentStatus = "DRAFT"
	ASSESSMENT_STATUS_IN_PROGRESS AssessmentStatus = "IN_PROGRESS"
	ASSESSMENT_STATUS_COMPLETED   AssessmentStatus = "COMPLETED"
)

type TestApplicant struct {
	gorm.Model           `json:"-"`
	ID                   uuid.UUID         `json:"id" gorm:"type:char(36);primaryKey;"`
	TestScheduleHeaderID uuid.UUID         `json:"test_schedule_header_id" gorm:"type:char(36);not null"`
	ApplicantID          uuid.UUID         `json:"applicant_id" gorm:"type:char(36);default:null"`
	UserProfileID        uuid.UUID         `json:"user_profile_id" gorm:"type:char(36);not null"`
	StartTime            time.Time         `json:"start_time" gorm:"type:time;not null"`
	EndTime              time.Time         `json:"end_time" gorm:"type:time;not null"`
	AssessmentStatus     AssessmentStatus  `json:"assessment_status" gorm:"type:text;default:null"`
	FinalResult          FinalResultStatus `json:"final_result" gorm:"type:text;default:null"`
	StartedTime          *time.Time        `json:"started_time" gorm:"type:datetime;default:null"`
	EndedTime            *time.Time        `json:"ended_time" gorm:"type:datetime;default:null"`

	TestScheduleHeader *TestScheduleHeader `json:"test_schedule_header" gorm:"foreignKey:TestScheduleHeaderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserProfile        *UserProfile        `json:"user_profile" gorm:"foreignKey:UserProfileID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Applicant          *Applicant          `json:"applicant" gorm:"foreignKey:ApplicantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (ta *TestApplicant) BeforeCreate(tx *gorm.DB) (err error) {
	ta.ID = uuid.New()
	ta.CreatedAt = time.Now()
	ta.UpdatedAt = time.Now()
	return
}

func (ta *TestApplicant) BeforeUpdate(tx *gorm.DB) (err error) {
	ta.UpdatedAt = time.Now()
	return
}

func (TestApplicant) TableName() string {
	return "test_applicants"
}
