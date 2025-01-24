package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SavedJob struct {
	gorm.Model    `json:"-"`
	JobPostingID  uuid.UUID `json:"job_posting_id" gorm:"type:char(36);primaryKey"`
	UserProfileID uuid.UUID `json:"user_profile_id" gorm:"type:char(36);primaryKey"`

	JobPosting  *JobPosting  `json:"job_posting" gorm:"foreignKey:JobPostingID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserProfile *UserProfile `json:"user_profile" gorm:"foreignKey:UserProfileID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (sj *SavedJob) BeforeCreate(tx *gorm.DB) (err error) {
	sj.CreatedAt = sj.CreatedAt
	sj.UpdatedAt = sj.UpdatedAt
	return nil
}

func (sj *SavedJob) BeforeUpdate(tx *gorm.DB) (err error) {
	sj.UpdatedAt = sj.UpdatedAt
	return nil
}

func (SavedJob) TableName() string {
	return "saved_jobs"
}
