package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdministrativeSelectionStatus string

const (
	ADMINISTRATIVE_SELECTION_STATUS_DRAFT       AdministrativeSelectionStatus = "DRAFT"
	ADMINISTRATIVE_SELECTION_STATUS_IN_PROGRESS AdministrativeSelectionStatus = "IN PROGRESS"
	ADMINISTRATIVE_SELECTION_STATUS_COMPLETED   AdministrativeSelectionStatus = "COMPLETED"
)

type AdministrativeSelection struct {
	gorm.Model     `json:"-"`
	ID             uuid.UUID                     `json:"id" gorm:"type:char(36);primaryKey;"`
	JobPostingID   uuid.UUID                     `json:"job_posting_id" gorm:"type:char(36);not null"`
	ProjectPicID   uuid.UUID                     `json:"project_pic_id" gorm:"type:char(36);not null"`
	Status         AdministrativeSelectionStatus `json:"status" gorm:"not null"`
	VerifiedAt     time.Time                     `json:"verified_at" gorm:"type:timestamp;null"`
	VerifiedBy     uuid.UUID                     `json:"verified_by" gorm:"type:char(36);null"`
	DocumentDate   time.Time                     `json:"document_date" gorm:"type:date;not null"`
	DocumentNumber string                        `json:"document_number" gorm:"type:varchar(255);not null"`

	JobPosting            *JobPosting            `json:"job_posting" gorm:"foreignKey:JobPostingID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ProjectPIC            *ProjectPic            `json:"project_pic" gorm:"foreignKey:ProjectPicID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	AdministrativeResults []AdministrativeResult `json:"administrative_results" gorm:"foreignKey:AdministrativeSelectionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	TotalApplicants int `json:"total_applicants" gorm:"-"`
}

func (a *AdministrativeSelection) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.New()
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	return
}

func (a *AdministrativeSelection) BeforeUpdate(tx *gorm.DB) (err error) {
	a.UpdatedAt = time.Now()
	return
}

func (a *AdministrativeSelection) BeforeDelete(tx *gorm.DB) (err error) {
	if a.DeletedAt.Valid {
		return nil
	}

	randomString := uuid.New().String()

	a.DocumentNumber = a.DocumentNumber + "_deleted" + randomString

	if err := tx.Model(&a).Where("id = ?", a.ID).Updates((map[string]interface{}{
		"document_number": a.DocumentNumber,
	})).Error; err != nil {
		return err
	}

	return nil
}

func (AdministrativeSelection) TableName() string {
	return "administrative_selections"
}
