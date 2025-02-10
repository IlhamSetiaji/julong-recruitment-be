package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FgdScheduleStatus string

const (
	FGD_SCHEDULE_STATUS_DRAFT       FgdScheduleStatus = "DRAFT"
	FGD_SCHEDULE_STATUS_IN_PROGRESS FgdScheduleStatus = "IN PROGRESS"
	FGD_SCHEDULE_STATUS_COMPLETED   FgdScheduleStatus = "COMPLETED"
)

type FgdSchedule struct {
	gorm.Model                 `json:"-"`
	ID                         uuid.UUID         `json:"id" gorm:"type:char(36);primaryKey;"`
	JobPostingID               uuid.UUID         `json:"job_posting_id" gorm:"type:char(36);not null"`
	ProjectPicID               uuid.UUID         `json:"project_pic_id" gorm:"type:char(36);not null"`
	ProjectRecruitmentHeaderID uuid.UUID         `json:"project_recruitment_header_id" gorm:"type:char(36);not null"`
	ProjectRecruitmentLineID   uuid.UUID         `json:"project_recruitment_line_id" gorm:"type:char(36);not null"`
	Name                       string            `json:"name" gorm:"type:varchar(255);not null"`
	DocumentNumber             string            `json:"document_number" gorm:"type:varchar(255);not null"`
	ScheduleDate               time.Time         `json:"schedule_date" gorm:"type:date;not null"`
	StartTime                  time.Time         `json:"start_time" gorm:"type:time;not null"`
	EndTime                    time.Time         `json:"end_time" gorm:"type:time;not null"`
	LocationLink               string            `json:"location_link" gorm:"type:text;default:null"`
	Description                string            `json:"description" gorm:"type:text;default:null"`
	RangeDuration              *int              `json:"range_duration" gorm:"type:int;default:0"`
	TotalCandidate             int               `json:"total_candidate" gorm:"type:int;default:0"`
	Status                     FgdScheduleStatus `json:"status" gorm:"type:varchar(255);default:'DRAFT'"`

	JobPosting               *JobPosting               `json:"job_posting" gorm:"foreignKey:JobPostingID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ProjectPic               *ProjectPic               `json:"project_pic" gorm:"foreignKey:ProjectPicID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ProjectRecruitmentHeader *ProjectRecruitmentHeader `json:"project_recruitment_header" gorm:"foreignKey:ProjectRecruitmentHeaderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ProjectRecruitmentLine   *ProjectRecruitmentLine   `json:"project_recruitment_line" gorm:"foreignKey:ProjectRecruitmentLineID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	FgdApplicants            []FgdApplicant            `json:"fgd_applicants" gorm:"foreignKey:FgdScheduleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	FgdAssessors             []FgdAssessor             `json:"fgd_assessors" gorm:"foreignKey:FgdScheduleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (ia *FgdSchedule) BeforeCreate(tx *gorm.DB) (err error) {
	ia.ID = uuid.New()
	ia.CreatedAt = time.Now()
	ia.UpdatedAt = time.Now()
	return
}

func (ia *FgdSchedule) BeforeUpdate(tx *gorm.DB) (err error) {
	ia.UpdatedAt = time.Now()
	return
}

func (tsh *FgdSchedule) BeforeDelete(tx *gorm.DB) (err error) {
	if tsh.DeletedAt.Valid {
		return nil
	}

	randomString := uuid.New().String()

	tsh.DocumentNumber = tsh.DocumentNumber + "_deleted" + randomString

	if err := tx.Model(&FgdApplicant{}).Where("interview_id = ?", tsh.ID).Delete(&FgdApplicant{}).Error; err != nil {
		return err
	}

	if err := tx.Model(&FgdAssessor{}).Where("interview_id = ?", tsh.ID).Delete(&FgdAssessor{}).Error; err != nil {
		return err
	}

	if err := tx.Model(&tsh).Where("id = ?", tsh.ID).Updates((map[string]interface{}{
		"document_number": tsh.DocumentNumber,
	})).Error; err != nil {
		return err
	}

	return nil
}

func (FgdSchedule) TableName() string {
	return "fgd_schedules"
}
