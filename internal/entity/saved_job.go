package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SavedJob struct {
	JobPostingID  uuid.UUID `json:"job_posting_id" gorm:"type:char(36);primaryKey"`
	UserProfileID uuid.UUID `json:"user_profile_id" gorm:"type:char(36);primaryKey"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	JobPosting  *JobPosting  `json:"job_posting" gorm:"foreignKey:JobPostingID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserProfile *UserProfile `json:"user_profile" gorm:"foreignKey:UserProfileID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (sj *SavedJob) BeforeCreate(tx *gorm.DB) (err error) {
	sj.CreatedAt = time.Now()
	sj.UpdatedAt = time.Now()
	return nil
}

func (sj *SavedJob) BeforeUpdate(tx *gorm.DB) (err error) {
	sj.UpdatedAt = time.Now()
	return nil
}

func (SavedJob) TableName() string {
	return "saved_jobs"
}
