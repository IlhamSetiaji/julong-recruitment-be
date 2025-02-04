package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InterviewStatus string

const (
	INTERVIEW_STATUS_DRAFT       InterviewStatus = "DRAFT"
	INTERVIEW_STATUS_IN_PROGRESS InterviewStatus = "IN PROGRESS"
	INTERVIEW_STATUS_COMPLETED   InterviewStatus = "COMPLETED"
)

type InterviewApplicant struct {
	gorm.Model     `json:"-"`
	ID             uuid.UUID       `json:"id" gorm:"type:char(36);primaryKey;"`
	JobPostingID   uuid.UUID       `json:"job_posting_id" gorm:"type:char(36);not null"`
	ProjectPicID   uuid.UUID       `json:"project_pic_id" gorm:"type:char(36);not null"`
	Name           string          `json:"name" gorm:"type:varchar(255);not null"`
	DocumentNumber string          `json:"document_number" gorm:"type:varchar(255);not null"`
	ScheduleDate   time.Time       `json:"schedule_date" gorm:"type:date;not null"`
	StartTime      time.Time       `json:"start_time" gorm:"type:datetime;not null"`
	EndTime        time.Time       `json:"end_time" gorm:"type:datetime;not null"`
	LocationLink   string          `json:"location_link" gorm:"type:text;default:null"`
	Description    string          `json:"description" gorm:"type:text;default:null"`
	RangeDuration  *int            `json:"range_duration" gorm:"type:int;default:0"`
	TotalCandidate int             `json:"total_candidate" gorm:"type:int;default:0"`
	Status         InterviewStatus `json:"status" gorm:"type:varchar(255);default:'DRAFT'"`
}

func (ia *InterviewApplicant) BeforeCreate(tx *gorm.DB) (err error) {
	ia.ID = uuid.New()
	ia.CreatedAt = time.Now()
	ia.UpdatedAt = time.Now()
	return
}

func (ia *InterviewApplicant) BeforeUpdate(tx *gorm.DB) (err error) {
	ia.UpdatedAt = time.Now()
	return
}

func (tsh *InterviewApplicant) BeforeDelete(tx *gorm.DB) (err error) {
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

func (InterviewApplicant) TableName() string {
	return "interview_applicants"
}
