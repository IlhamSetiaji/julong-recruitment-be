package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MPRequestStatus string

const (
	MPR_STATUS_OPEN      MPRequestStatus = "OPEN"
	MPR_STATUS_ON_GOING  MPRequestStatus = "ON_GOING"
	MPR_STATUS_COMPLETED MPRequestStatus = "COMPLETED"
)

type MPRequest struct {
	gorm.Model `json:"-"`
	ID         uuid.UUID       `json:"id" gorm:"type:char(36);primaryKey;"`
	MPRCloneID *uuid.UUID      `json:"mpr_clone_id" gorm:"type:char(36);not null;unique"`
	Status     MPRequestStatus `json:"status" gorm:"default:'OPEN'"`

	JobPosting *JobPosting `json:"job_posting" gorm:"foreignKey:MPRequestID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (mpr *MPRequest) BeforeCreate(tx *gorm.DB) (err error) {
	mpr.ID = uuid.New()
	mpr.CreatedAt = time.Now()
	mpr.UpdatedAt = time.Now()
	return nil
}

func (mpr *MPRequest) BeforeUpdate(tx *gorm.DB) (err error) {
	mpr.UpdatedAt = time.Now()
	return nil
}

func (MPRequest) TableName() string {
	return "mp_requests"
}
