package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectRecruitmentType string

const (
	PROJECT_RECRUITMENT_TYPE_MT ProjectRecruitmentType = "MT_Management Trainee"
	PROJECT_RECRUITMENT_TYPE_PH ProjectRecruitmentType = "PH_Professional Hire"
	PROJECT_RECRUITMENT_TYPE_NS ProjectRecruitmentType = "NS_Non Staff to Staff"
)

type ProjectRecruitmentHeaderStatus string

const (
	PROJECT_RECRUITMENT_HEADER_STATUS_DRAFT       ProjectRecruitmentHeaderStatus = "DRAFT"
	PROJECT_RECRUITMENT_HEADER_STATUS_SUBMITTED   ProjectRecruitmentHeaderStatus = "SUBMITTED"
	PROJECT_RECRUITMENT_HEADER_STATUS_APPROVED    ProjectRecruitmentHeaderStatus = "APPROVED"
	PROJECT_RECRUITMENT_HEADER_STATUS_REJECTED    ProjectRecruitmentHeaderStatus = "REJECTED"
	PROJECT_RECRUITMENT_HEADER_STATUS_CLOSE       ProjectRecruitmentHeaderStatus = "CLOSE"
	PROJECT_RECRUITMENT_HEADER_STATUS_IN_PROGRESS ProjectRecruitmentHeaderStatus = "IN PROGRESS"
	PROJECT_RECRUITMENT_HEADER_STATUS_COMPLETED   ProjectRecruitmentHeaderStatus = "COMPLETED"
)

type ProjectRecruitmentHeader struct {
	gorm.Model         `json:"-"`
	ID                 uuid.UUID                      `json:"id" gorm:"type:char(36);primaryKey;"`
	TemplateActivityID uuid.UUID                      `json:"template_activity_id" gorm:"type:char(36);not null"`
	Name               string                         `json:"name" gorm:"type:varchar(255);not null"`
	Description        string                         `json:"description" gorm:"type:text;default:null"`
	DocumentDate       time.Time                      `json:"document_date" gorm:"type:date;not null"`
	DocumentNumber     string                         `json:"document_number" gorm:"type:varchar(255);not null"`
	RecruitmentType    ProjectRecruitmentType         `json:"recruitment_type" gorm:"not null"`
	StartDate          time.Time                      `json:"start_date" gorm:"type:date;not null"`
	EndDate            time.Time                      `json:"end_date" gorm:"type:date;not null"`
	Status             ProjectRecruitmentHeaderStatus `json:"status" gorm:"default:'DRAFT'"`
	ProjectPicID       *uuid.UUID                     `json:"project_pic_id" gorm:"type:char(36);not null"`

	JobPostings             []JobPosting             `json:"job_postings" gorm:"foreignKey:ProjectRecruitmentHeaderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ProjectRecruitmentLines []ProjectRecruitmentLine `json:"project_recruitment_lines" gorm:"foreignKey:ProjectRecruitmentHeaderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TemplateActivity        *TemplateActivity        `json:"template_activity" gorm:"foreignKey:TemplateActivityID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TestScheduleHeaders     []TestScheduleHeader     `json:"test_schedule_headers" gorm:"foreignKey:ProjectRecruitmentHeaderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (prh *ProjectRecruitmentHeader) BeforeCreate(tx *gorm.DB) (err error) {
	prh.ID = uuid.New()
	prh.CreatedAt = time.Now()
	prh.UpdatedAt = time.Now()
	return nil
}

func (prh *ProjectRecruitmentHeader) BeforeUpdate(tx *gorm.DB) (err error) {
	prh.UpdatedAt = time.Now()
	return nil
}

func (prh *ProjectRecruitmentHeader) BeforeDelete(tx *gorm.DB) (err error) {
	if prh.DeletedAt.Valid {
		return nil
	}

	randomString := uuid.New().String()

	prh.DocumentNumber = prh.DocumentNumber + "_deleted" + randomString

	if err := tx.Model(&prh).Where("id = ?", prh.ID).Updates((map[string]interface{}{
		"document_number": prh.DocumentNumber,
	})).Error; err != nil {
		return err
	}

	return nil
}

func (ProjectRecruitmentHeader) TableName() string {
	return "project_recruitment_headers"
}
