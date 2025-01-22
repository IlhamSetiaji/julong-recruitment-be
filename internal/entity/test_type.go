package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TestTypeStatus string

const (
	TEST_TYPE_STATUS_ACTIVE   TestTypeStatus = "ACTIVE"
	TEST_TYPE_STATUS_INACTIVE TestTypeStatus = "INACTIVE"
)

type TestType struct {
	gorm.Model      `json:"-"`
	ID              uuid.UUID              `json:"id" gorm:"type:char(36);primaryKey;"`
	Name            string                 `json:"name" gorm:"type:text;not null"`
	RecruitmentType ProjectRecruitmentType `json:"recruitment_type" gorm:"type:varchar(255);not null"`
	Status          TestTypeStatus         `json:"status" gorm:"type:varchar(255);default:'ACTIVE'"`
}

func (tt *TestType) BeforeCreate(tx *gorm.DB) (err error) {
	tt.ID = uuid.New()
	tt.CreatedAt = time.Now()
	tt.UpdatedAt = time.Now()
	return
}

func (tt *TestType) BeforeUpdate(tx *gorm.DB) (err error) {
	tt.UpdatedAt = time.Now()
	return
}

func (TestType) TableName() string {
	return "test_types"
}
