package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnswerType struct {
	gorm.Model `json:"-"`
	ID         uuid.UUID `json:"id" gorm:"type:char(36);primaryKey;"`
	Name       string    `json:"name" gorm:"type:varchar(255);not null"`

	Questions []Question `json:"questions" gorm:"foreignKey:AnswerTypeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
