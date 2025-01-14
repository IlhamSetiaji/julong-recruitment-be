package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnswerType struct {
	gorm.Model `json:"-"`
	ID         uuid.UUID `json:"id" gorm:"type:char(36);primaryKey;"`
	Name       string    `json:"name" gorm:"type:varchar(255);not null"`

	Questions []Question `json:"questions" gorm:"foreignKey:AnswerTypeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (at *AnswerType) BeforeCreate(tx *gorm.DB) (err error) {
	at.ID = uuid.New()
	at.CreatedAt = time.Now()
	at.UpdatedAt = time.Now()
	return
}

func (at *AnswerType) BeforeUpdate(tx *gorm.DB) (err error) {
	at.UpdatedAt = time.Now()
	return
}

func (AnswerType) TableName() string {
	return "answer_types"
}
