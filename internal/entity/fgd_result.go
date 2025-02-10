package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FgdResultStatus string

const (
	FGD_RESULT_STATUS_ACCEPTED FgdResultStatus = "ACCEPTED"
	FGD_RESULT_STATUS_REJECTED FgdResultStatus = "REJECTED"
)

type FgdResult struct {
	gorm.Model     `json:"-"`
	ID             uuid.UUID       `json:"id" gorm:"type:char(36);primaryKey;"`
	FgdApplicantID uuid.UUID       `json:"fgd_id" gorm:"type:char(36);not null"`
	FgdAssessorID  uuid.UUID       `json:"fgd_assessor_id" gorm:"type:char(36);not null"`
	Status         FgdResultStatus `json:"status" gorm:"type:varchar(255);default:'REJECTED'"`

	FgdApplicant *FgdApplicant `json:"fgd_applicant" gorm:"foreignKey:FgdApplicantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	FgdAssessor  *FgdAssessor  `json:"fgd_assessor" gorm:"foreignKey:FgdAssessorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (fr *FgdResult) BeforeCreate(tx *gorm.DB) (err error) {
	fr.ID = uuid.New()
	fr.CreatedAt = time.Now()
	fr.UpdatedAt = time.Now()
	return
}

func (fr *FgdResult) BeforeUpdate(tx *gorm.DB) (err error) {
	fr.UpdatedAt = time.Now()
	return
}

func (FgdResult) TableName() string {
	return "fgd_results"
}
