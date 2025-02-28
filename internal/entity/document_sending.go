package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DocumentSendingStatus string

const (
	DOCUMENT_SENDING_STATUS_DRAFT     DocumentSendingStatus = "DRAFT"
	DOCUMENT_SENDING_STATUS_PENDING   DocumentSendingStatus = "PENDING"
	DOCUMENT_SENDING_STATUS_SENT      DocumentSendingStatus = "SENT"
	DOCUMENT_SENDING_STATUS_FAILED    DocumentSendingStatus = "FAILED"
	DOCUMENT_SENDING_STATUS_REVISED   DocumentSendingStatus = "REVISED"
	DOCUMENT_SENDING_STATUS_REJECTED  DocumentSendingStatus = "REJECTED"
	DOCUMENT_SENDING_STATUS_APPROVED  DocumentSendingStatus = "APPROVED"
	DOCUMENT_SENDING_STATUS_COMPLETED DocumentSendingStatus = "COMPLETED"
)

type DocumentSending struct {
	gorm.Model               `json:"-"`
	ID                       uuid.UUID              `json:"id" gorm:"type:char(36);primaryKey;"`
	ProjectRecruitmentLineID uuid.UUID              `json:"project_recruitment_line_id" gorm:"type:char(36);not null"`
	ApplicantID              uuid.UUID              `json:"applicant_id" gorm:"type:char(36);not null"`
	DocumentSetupID          uuid.UUID              `json:"document_setup_id" gorm:"type:char(36);not null"`
	DocumentDate             time.Time              `json:"document_date" gorm:"type:date;not null"`
	DocumentNumber           string                 `json:"document_number" gorm:"type:varchar(255);not null;unique"`
	Status                   DocumentSendingStatus  `json:"status" gorm:"default:'PENDING'"`
	RecruitmentType          ProjectRecruitmentType `json:"recruitment_type" gorm:"type:enum('INTERNAL','EXTERNAL');default:'EXTERNAL'"`
	BasicWage                float64                `json:"basic_wage" gorm:"type:decimal(10,2);default:null"`
	PositionalAllowance      float64                `json:"positional_allowance" gorm:"type:decimal(10,2);default:null"`
	OperationalAllowance     float64                `json:"operational_allowance" gorm:"type:decimal(10,2);default:null"`
	MealAllowance            float64                `json:"meal_allowance" gorm:"type:decimal(10,2);default:null"`
	JobLocation              string                 `json:"job_location" gorm:"type:text;default:null"`
	HometripTicket           string                 `json:"hometrip_ticket" gorm:"type:text;default:null"`
	PeriodAgreement          string                 `json:"period_agreement" gorm:"type:text;default:null"`
	HomeLocation             string                 `json:"home_location" gorm:"type:text;default:null"`
	JobLevelID               *uuid.UUID             `json:"job_level_id" gorm:"type:char(36);default:null"`
	JobID                    *uuid.UUID             `json:"job_id" gorm:"type:char(36);default:null"`
	JoinedDate               *time.Time             `json:"joined_date" gorm:"type:date;default:null"`
	ForOrganizationID        *uuid.UUID             `json:"for_organization_id" gorm:"type:char(36);default:null"`
	OrganizationLocationID   *uuid.UUID             `json:"organization_location_id" gorm:"type:char(36);default:null"`
	JobPostingID             uuid.UUID              `json:"job_posting_id" gorm:"type:char(36);default:null"`
	DetailContent            string                 `json:"detail_content" gorm:"type:text;default:null"`
	Path                     string                 `json:"path" gorm:"type:text;default:null"`

	ProjectRecruitmentLine *ProjectRecruitmentLine `json:"project_recruitment_line" gorm:"foreignKey:ProjectRecruitmentLineID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Applicant              *Applicant              `json:"applicant" gorm:"foreignKey:ApplicantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	DocumentSetup          *DocumentSetup          `json:"document_setup" gorm:"foreignKey:DocumentSetupID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	JobPosting             *JobPosting             `json:"job_posting" gorm:"foreignKey:JobPostingID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	JobLevelName        *string `json:"job_level_name" gorm:"-"`
	JobName             *string `json:"job_name" gorm:"-"`
	ForOrganizationName *string `json:"for_organization_name" gorm:"-"`
}

func (ds *DocumentSending) BeforeCreate(tx *gorm.DB) (err error) {
	ds.ID = uuid.New()
	ds.CreatedAt = time.Now()
	ds.UpdatedAt = time.Now()
	return
}

func (ds *DocumentSending) BeforeUpdate(tx *gorm.DB) (err error) {
	ds.UpdatedAt = time.Now()
	return
}

func (ds *DocumentSending) BeforeDelete(tx *gorm.DB) (err error) {
	if ds.DeletedAt.Valid {
		return nil
	}

	randomString := uuid.New().String()

	ds.DocumentNumber = ds.DocumentNumber + "_deleted" + randomString

	if err := tx.Model(&ds).Where("id = ?", ds.ID).Updates((map[string]interface{}{
		"document_number": ds.DocumentNumber,
	})).Error; err != nil {
		return err
	}

	return nil
}

func (DocumentSending) TableName() string {
	return "document_sendings"
}
