package response

import (
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
)

type JobPostingResponse struct {
	ID                         uuid.UUID                     `json:"id"`
	ProjectRecruitmentHeaderID uuid.UUID                     `json:"project_recruitment_header_id"`
	MPRequestID                *uuid.UUID                    `json:"mp_request_id"`
	JobID                      *uuid.UUID                    `json:"job_id"`
	ForOrganizationID          *uuid.UUID                    `json:"for_organization_id"`
	ForOrganizationLocationID  *uuid.UUID                    `json:"for_organization_location_id"`
	DocumentNumber             string                        `json:"document_number"`
	DocumentDate               string                        `json:"document_date"`
	RecruitmentType            entity.ProjectRecruitmentType `json:"recruitment_type"`
	StartDate                  string                        `json:"start_date"`
	EndDate                    string                        `json:"end_date"`
	Status                     entity.JobPostingStatus       `json:"status"`
	SalaryMin                  string                        `json:"salary_min"`
	SalaryMax                  string                        `json:"salary_max"`
	ContentDescription         string                        `json:"content_description"`
	OrganizationLogo           *string                       `json:"organization_logo"`
	Poster                     *string                       `json:"poster"`
	Link                       string                        `json:"link"`
	IsApplied                  bool                          `json:"is_applied"`
	IsSaved                    bool                          `json:"is_saved"`
	AppliedDate                time.Time                     `json:"applied_date"`
	ApplicantStatus            entity.ApplicantStatus        `json:"apply_status"`
	ApplicantProcessStatus     entity.ApplicantProcessStatus `json:"apply_process_status"`
	MinimumWorkExperience      string                        `json:"minimum_work_experience"`
	Name                       string                        `json:"name"`
	IsShow                     string                        `json:"is_show"`

	ForOrganizationName     string `json:"for_organization_name"`
	ForOrganizationLocation string `json:"for_organization_location"`
	JobName                 string `json:"job_name"`
	TotalApplicant          int    `json:"total_applicant"`

	ProjectRecruitmentHeader *ProjectRecruitmentHeaderResponse `json:"project_recruitment_header"`
	MPRequest                *MPRequestHeaderResponse          `json:"mp_request"`
}
