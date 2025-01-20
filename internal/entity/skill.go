package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Skill struct {
	gorm.Model  `json:"-"`
	ID          uuid.UUID `json:"id" gorm:"type:char(36);primaryKey;"`
	Name        string    `json:"name" gorm:"type:varchar(255);not null"`
	Description string    `json:"description" gorm:"type:text;not null"`
	Certificate string    `json:"certificate" gorm:"type:text;not null"`
}
