package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdministrativeResultStatus string

const (
	ADMINISTRATIVE_RESULT_STATUS_ACCEPTED    AdministrativeResultStatus = "ACCEPTED"
	ADMINISTRATIVE_RESULT_STATUS_REJECTED    AdministrativeResultStatus = "REJECTED"
	ADMINISTRATIVE_RESULT_STATUS_SHORTLISTED AdministrativeResultStatus = "SHORTLISTED"
)

type AdministrativeResult struct {
	gorm.Model                `json:"-"`
	ID                        uuid.UUID                  `json:"id" gorm:"type:char(36);primaryKey;"`
	AdministrativeSelectionID uuid.UUID                  `json:"administrative_selection_id" gorm:"type:char(36);not null"`
	UserProfileID             uuid.UUID                  `json:"user_profile_id" gorm:"type:char(36);not null"`
	Status                    AdministrativeResultStatus `json:"status" gorm:"not null"`

	AdministrativeSelection *AdministrativeSelection `json:"administrative_selection" gorm:"foreignKey:AdministrativeSelectionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserProfile             *UserProfile             `json:"user_profile" gorm:"foreignKey:UserProfileID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (a *AdministrativeResult) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.New()
	a.CreatedAt = a.CreatedAt.Local()
	a.UpdatedAt = a.UpdatedAt.Local()
	return
}

func (a *AdministrativeResult) BeforeUpdate(tx *gorm.DB) (err error) {
	a.UpdatedAt = a.UpdatedAt.Local()
	return
}

func (AdministrativeResult) TableName() string {
	return "administrative_results"
}
