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

type Interview struct {
	gorm.Model                 `json:"-"`
	ID                         uuid.UUID       `json:"id" gorm:"type:char(36);primaryKey;"`
	JobPostingID               uuid.UUID       `json:"job_posting_id" gorm:"type:char(36);not null"`
	ProjectPicID               uuid.UUID       `json:"project_pic_id" gorm:"type:char(36);not null"`
	ProjectRecruitmentHeaderID uuid.UUID       `json:"project_recruitment_header_id" gorm:"type:char(36);not null"`
	ProjectRecruitmentLineID   uuid.UUID       `json:"project_recruitment_line_id" gorm:"type:char(36);not null"`
	Name                       string          `json:"name" gorm:"type:varchar(255);not null"`
	DocumentNumber             string          `json:"document_number" gorm:"type:varchar(255);not null"`
	ScheduleDate               time.Time       `json:"schedule_date" gorm:"type:date;not null"`
	StartTime                  time.Time       `json:"start_time" gorm:"type:time;not null"`
	EndTime                    time.Time       `json:"end_time" gorm:"type:time;not null"`
	LocationLink               string          `json:"location_link" gorm:"type:text;default:null"`
	Description                string          `json:"description" gorm:"type:text;default:null"`
	RangeDuration              *int            `json:"range_duration" gorm:"type:int;default:0"`
	TotalCandidate             int             `json:"total_candidate" gorm:"type:int;default:0"`
	Status                     InterviewStatus `json:"status" gorm:"type:varchar(255);default:'DRAFT'"`
	MeetingLink                string          `json:"meeting_link" gorm:"type:text;default:null"`

	JobPosting               *JobPosting               `json:"job_posting" gorm:"foreignKey:JobPostingID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ProjectPic               *ProjectPic               `json:"project_pic" gorm:"foreignKey:ProjectPicID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ProjectRecruitmentHeader *ProjectRecruitmentHeader `json:"project_recruitment_header" gorm:"foreignKey:ProjectRecruitmentHeaderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ProjectRecruitmentLine   *ProjectRecruitmentLine   `json:"project_recruitment_line" gorm:"foreignKey:ProjectRecruitmentLineID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	InterviewApplicants      []InterviewApplicant      `json:"interview_applicants" gorm:"foreignKey:InterviewID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	InterviewAssessors       []InterviewAssessor       `json:"interview_assessors" gorm:"foreignKey:InterviewID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (ia *Interview) BeforeCreate(tx *gorm.DB) (err error) {
	ia.ID = uuid.New()
	ia.CreatedAt = time.Now()
	ia.UpdatedAt = time.Now()
	return
}

func (ia *Interview) BeforeUpdate(tx *gorm.DB) (err error) {
	ia.UpdatedAt = time.Now()
	return
}

func (tsh *Interview) BeforeDelete(tx *gorm.DB) (err error) {
	if tsh.DeletedAt.Valid {
		return nil
	}

	randomString := uuid.New().String()

	tsh.DocumentNumber = tsh.DocumentNumber + "_deleted" + randomString

	if err := tx.Model(&InterviewApplicant{}).Where("interview_id = ?", tsh.ID).Delete(&InterviewApplicant{}).Error; err != nil {
		return err
	}

	if err := tx.Model(&InterviewAssessor{}).Where("interview_id = ?", tsh.ID).Delete(&InterviewAssessor{}).Error; err != nil {
		return err
	}

	if err := tx.Model(&tsh).Where("id = ?", tsh.ID).Updates((map[string]interface{}{
		"document_number": tsh.DocumentNumber,
	})).Error; err != nil {
		return err
	}

	return nil
}

func (Interview) TableName() string {
	return "interviews"
}
