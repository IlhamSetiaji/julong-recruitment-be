package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Skill struct {
	gorm.Model  `json:"-"`
	ID          uuid.UUID `json:"id" gorm:"type:char(36);primaryKey;"`
	ApplicantID uuid.UUID `json:"applicant_id" gorm:"type:char(36);not null"`
	Name        string    `json:"name" gorm:"type:varchar(255);not null"`
	Description string    `json:"description" gorm:"type:text;not null"`
	Certificate string    `json:"certificate" gorm:"type:text;not null"`

	Applicant *Applicant `json:"applicant" gorm:"foreignKey:ApplicantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (s *Skill) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.New()
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	return nil
}

func (s *Skill) BeforeUpdate(tx *gorm.DB) (err error) {
	s.UpdatedAt = time.Now()
	return nil
}

func (Skill) TableName() string {
	return "skills"
}
