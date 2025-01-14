package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JobPostingStatus string

const (
	JOB_POSTING_STATUS_DRAFT       JobPostingStatus = "DRAFT"
	JOB_POSTING_STATUS_SUBMITTED   JobPostingStatus = "SUBMITTED"
	JOB_POSTING_STATUS_APPROVED    JobPostingStatus = "APPROVED"
	JOB_POSTING_STATUS_REJECTED    JobPostingStatus = "REJECTED"
	JOB_POSTING_STATUS_CLOSE       JobPostingStatus = "CLOSE"
	JOB_POSTING_STATUS_IN_PROGRESS JobPostingStatus = "IN PROGRESS"
	JOB_POSTING_STATUS_COMPLETED   JobPostingStatus = "COMPLETED"
)

type JobPosting struct {
	gorm.Model                 `json:"-"`
	ID                         uuid.UUID              `json:"id" gorm:"type:char(36);primaryKey;"`
	ProjectRecruitmentHeaderID uuid.UUID              `json:"project_recruitment_header_id" gorm:"type:char(36);not null"`
	MPRequestID                *uuid.UUID             `json:"mp_request_id" gorm:"type:char(36);not null'unique"`
	JobID                      *uuid.UUID             `json:"job_id" gorm:"type:char(36);not null"`
	ForOrganizationID          *uuid.UUID             `json:"for_organization_id" gorm:"type:char(36);not null"`
	ForOrganizationLocationID  *uuid.UUID             `json:"for_organization_location_id" gorm:"type:char(36);not null"`
	DocumentNumber             string                 `json:"document_number" gorm:"type:varchar(255);not null;unique"`
	DocumentDate               time.Time              `json:"document_date" gorm:"type:date;not null"`
	RecruitmentType            ProjectRecruitmentType `json:"recruitment_type" gorm:"not null"`
	StartDate                  time.Time              `json:"start_date" gorm:"type:date;not null"`
	EndDate                    time.Time              `json:"end_date" gorm:"type:date;not null"`
	Status                     JobPostingStatus       `json:"status" gorm:"default:'DRAFT'"`
	SalaryMin                  string                 `json:"salary_min" gorm:"type:varchar(255);not null"`
	SalaryMax                  string                 `json:"salary_max" gorm:"type:varchar(255);not null"`
	ContentDescription         string                 `json:"content_description" gorm:"type:text;default:null"`
	OrganizationLogo           string                 `json:"organization_logo" gorm:"type:text;default:null"`
	Poster                     string                 `json:"poster" gorm:"type:text;default:null"`
	Link                       string                 `json:"link" gorm:"type:text;default:null"`

	ProjectRecruitmentHeader ProjectRecruitmentHeader `json:"project_recruitment_header" gorm:"foreignKey:ProjectRecruitmentHeaderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	MPRequest                *MPRequest               `json:"mp_request" gorm:"foreignKey:MPRequestID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	ForOrganizationName      string `json:"for_organization_name" gorm:"-"`
	ForOrganizationLocation  string `json:"for_organization_location" gorm:"-"`
	ForOrganizationStructure string `json:"for_organization_structure" gorm:"-"`
	JobName                  string `json:"job_name" gorm:"-"`
}

func (jp *JobPosting) BeforeCreate(tx *gorm.DB) (err error) {
	jp.ID = uuid.New()
	jp.CreatedAt = time.Now()
	jp.UpdatedAt = time.Now()
	return nil
}

func (jp *JobPosting) BeforeUpdate(tx *gorm.DB) (err error) {
	jp.UpdatedAt = time.Now()
	return nil
}

func (JobPosting) TableName() string {
	return "job_postings"
}
