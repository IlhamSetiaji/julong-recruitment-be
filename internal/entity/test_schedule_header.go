package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TestScheduleStatus string

const (
	TEST_SCHEDULE_STATUS_DRAFT       TestScheduleStatus = "DRAFT"
	TEST_SCHEDULE_STATUS_IN_PROGRESS TestScheduleStatus = "IN PROGRESS"
	TEST_SCHEDULE_STATUS_COMPLETED   TestScheduleStatus = "COMPLETED"
)

type TestScheduleHeader struct {
	gorm.Model     `json:"-"`
	ID             uuid.UUID          `json:"id" gorm:"type:char(36);primaryKey;"`
	JobPostingID   uuid.UUID          `json:"job_posting_id" gorm:"type:char(36);not null"`
	TestTypeID     uuid.UUID          `json:"test_type_id" gorm:"type:char(36);not null"`
	ProjectPicID   uuid.UUID          `json:"project_pic_id" gorm:"type:char(36);not null"`
	JobID          *uuid.UUID         `json:"job_id" gorm:"type:char(36);not null"`
	Name           string             `json:"name" gorm:"type:text;not null"`
	DocumentNumber string             `json:"document_number" gorm:"type:text;not null"`
	StartDate      time.Time          `json:"start_date" gorm:"type:date;not null"`
	EndDate        time.Time          `json:"end_date" gorm:"type:date;not null"`
	StartTime      time.Time          `json:"start_time" gorm:"type:time;not null"`
	EndTime        time.Time          `json:"end_time" gorm:"type:time;not null"`
	Link           string             `json:"link" gorm:"type:text;not null"`
	Location       string             `json:"location" gorm:"type:text;not null"`
	Description    string             `json:"description" gorm:"type:text;not null"`
	TotalCandidate int                `json:"total_candidate" gorm:"type:int;not null"`
	Status         TestScheduleStatus `json:"status" gorm:"type:varchar(255);default:'DRAFT'"`

	JobPosting *JobPosting `json:"job_posting" gorm:"foreignKey:JobPostingID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TestType   *TestType   `json:"test_type" gorm:"foreignKey:TestTypeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ProjectPic *ProjectPic `json:"project_pic" gorm:"foreignKey:ProjectPicID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	JobName string `json:"job_name" gorm:"-"`
}

func (tsh *TestScheduleHeader) BeforeCreate(tx *gorm.DB) (err error) {
	tsh.ID = uuid.New()
	tsh.CreatedAt = time.Now()
	tsh.UpdatedAt = time.Now()
	return
}

func (tsh *TestScheduleHeader) BeforeUpdate(tx *gorm.DB) (err error) {
	tsh.UpdatedAt = time.Now()
	return
}

func (tsh *TestScheduleHeader) BeforeDelete(tx *gorm.DB) (err error) {
	if tsh.DeletedAt.Valid {
		return nil
	}

	randomString := uuid.New().String()

	tsh.DocumentNumber = tsh.DocumentNumber + "_deleted" + randomString

	if err := tx.Model(&tsh).Where("id = ?", tsh.ID).Updates((map[string]interface{}{
		"document_number": tsh.DocumentNumber,
	})).Error; err != nil {
		return err
	}

	return nil
}

func (TestScheduleHeader) TableName() string {
	return "test_schedule_headers"
}
