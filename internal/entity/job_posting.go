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
	MinimumWorkExperience      string                 `json:"minimum_work_experience" gorm:"type:varchar(255);default:null"`

	ProjectRecruitmentHeader *ProjectRecruitmentHeader `json:"project_recruitment_header" gorm:"foreignKey:ProjectRecruitmentHeaderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	MPRequest                *MPRequest                `json:"mp_request" gorm:"foreignKey:MPRequestID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Applicants               []Applicant               `json:"applicants" gorm:"foreignKey:JobPostingID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TestScheduleHeaders      []TestScheduleHeader      `json:"test_schedule_headers" gorm:"foreignKey:JobPostingID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserProfiles             []UserProfile             `json:"user_profiles" gorm:"many2many:saved_jobs;foreignKey:ID;joinForeignKey:JobPostingID;References:ID;JoinReferences:UserProfileID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	AdministrativeSelections []AdministrativeSelection `json:"administrative_selections" gorm:"foreignKey:JobPostingID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	ForOrganizationName      string                 `json:"for_organization_name" gorm:"-"`
	ForOrganizationLocation  string                 `json:"for_organization_location" gorm:"-"`
	ForOrganizationStructure string                 `json:"for_organization_structure" gorm:"-"`
	JobName                  string                 `json:"job_name" gorm:"-"`
	IsApplied                bool                   `json:"is_applied" gorm:"-"`
	IsSaved                  bool                   `json:"is_saved" gorm:"-"`
	AppliedDate              time.Time              `json:"applied_date" gorm:"-"`
	ApplicantStatus          ApplicantStatus        `json:"apply_status" gorm:"-"`
	ApplicantProcessStatus   ApplicantProcessStatus `json:"apply_process_status" gorm:"-"`
	TotalApplicant           int                    `json:"total_applicant" gorm:"-"`
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

func (jp *JobPosting) BeforeDelete(tx *gorm.DB) (err error) {
	if jp.DeletedAt.Valid {
		return nil
	}

	randomString := uuid.New().String()

	jp.DocumentNumber = jp.DocumentNumber + "_deleted" + randomString

	if err := tx.Model(&jp).Where("id = ?", jp.ID).Updates((map[string]interface{}{
		"document_number": jp.DocumentNumber,
	})).Error; err != nil {
		return err
	}

	return nil
}

func (JobPosting) TableName() string {
	return "job_postings"
}
