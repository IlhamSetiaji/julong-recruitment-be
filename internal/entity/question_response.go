package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuestionResponse struct {
	gorm.Model          `json:"-"`
	ID                  uuid.UUID `json:"id" gorm:"type:char(36);primaryKey;"`
	QuestionID          uuid.UUID `json:"question_id" gorm:"type:char(36);primaryKey;not null"`
	JobPostingID        uuid.UUID `json:"job_posting_id" gorm:"type:char(36);not null"`
	UserProfileID       uuid.UUID `json:"user_profile_id" gorm:"type:char(36);primaryKey;not null"`
	InterviewAssessorID uuid.UUID `json:"interview_assessor_id" gorm:"type:char(36);default:null"`
	FgdAssessorID       uuid.UUID `json:"fgd_assessor_id" gorm:"type:char(36);default:null"`
	Answer              string    `json:"answer" gorm:"type:text;default:null"`
	AnswerFile          string    `json:"answer_file" gorm:"type:text;default:null"`

	Question          *Question          `json:"question" gorm:"foreignKey:QuestionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserProfile       *UserProfile       `json:"user_profile" gorm:"foreignKey:UserProfileID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	JobPosting        *JobPosting        `json:"job_posting" gorm:"foreignKey:JobPostingID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	InterviewAssessor *InterviewAssessor `json:"interview_assessor" gorm:"foreignKey:InterviewAssessorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	FgdAssessor       *FgdAssessor       `json:"fgd_assessor" gorm:"foreignKey:FgdAssessorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (qr *QuestionResponse) BeforeCreate(tx *gorm.DB) (err error) {
	qr.ID = uuid.New()
	qr.CreatedAt = time.Now()
	qr.UpdatedAt = time.Now()
	return
}

func (qr *QuestionResponse) BeforeUpdate(tx *gorm.DB) (err error) {
	qr.UpdatedAt = time.Now()
	return
}

func (QuestionResponse) TableName() string {
	return "question_responses"
}
