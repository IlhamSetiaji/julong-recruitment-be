package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InterviewApplicant struct {
	gorm.Model `json:"-"`
	ID         uuid.UUID `json:"id" gorm:"type:char(36);primaryKey;"`
	
}
