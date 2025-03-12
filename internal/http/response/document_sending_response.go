package response

import (
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
)

type DocumentSendingResponse struct {
	ID                       uuid.UUID                     `json:"id"`
	ProjectRecruitmentLineID uuid.UUID                     `json:"project_recruitment_line_id"`
	ApplicantID              uuid.UUID                     `json:"applicant_id"`
	DocumentSetupID          uuid.UUID                     `json:"document_setup_id"`
	DocumentDate             time.Time                     `json:"document_date"`
	JoinedDate               *time.Time                    `json:"joined_date"`
	DocumentNumber           string                        `json:"document_number"`
	Status                   entity.DocumentSendingStatus  `json:"status"`
	RecruitmentType          entity.ProjectRecruitmentType `json:"recruitment_type"`
	BasicWage                float64                       `json:"basic_wage"`
	PositionalAllowance      float64                       `json:"positional_allowance"`
	OperationalAllowance     float64                       `json:"operational_allowance"`
	MealAllowance            float64                       `json:"meal_allowance"`
	JobLocation              string                        `json:"job_location"`
	HometripTicket           string                        `json:"hometrip_ticket"`
	PeriodAgreement          string                        `json:"period_agreement"`
	HomeLocation             string                        `json:"home_location"`
	JobLevelID               *uuid.UUID                    `json:"job_level_id"`
	JobID                    *uuid.UUID                    `json:"job_id"`
	JobPostingID             uuid.UUID                     `json:"job_posting_id"`
	ForOrganizationID        *uuid.UUID                    `json:"for_organization_id"`
	OrganizationLocationID   *uuid.UUID                    `json:"organization_location_id"`
	DetailContent            string                        `json:"detail_content"`
	Path                     string                        `json:"path"`
	SyncMidsuit              entity.SyncMidsuitEnum        `json:"sync_midsuit"`
	HiredStatus              entity.HiredStatusEnum        `json:"hired_status"`
	CreatedAt                time.Time                     `json:"created_at"`
	UpdatedAt                time.Time                     `json:"updated_at"`

	// ProjectRecruitmentLine *ProjectRecruitmentLineResponse `json:"project_recruitment_line"`
	// Applicant              *ApplicantResponse              `json:"applicant"`
	DocumentSetup          *DocumentSetupResponse          `json:"document_setup"`
	JobPosting             *JobPostingResponse             `json:"job_posting"`
	ProjectRecruitmentLine *ProjectRecruitmentLineResponse `json:"project_recruitment_line"`
	Applicant              *ApplicantResponse              `json:"applicant"`

	ForOrganizationName      *string                              `json:"for_organization_name"`
	OrganizationLocationName *string                              `json:"organization_location_name"`
	JobLevel                 *SendFindJobLevelByIDMessageResponse `json:"job_level"`
	Job                      *SendFindJobByIDMessageResponse      `json:"job"`
}
