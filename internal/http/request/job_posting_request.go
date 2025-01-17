package request

import "mime/multipart"

type CreateJobPostingRequest struct {
	ProjectRecruitmentHeaderID string                `form:"project_recruitment_header_id" validate:"required,uuid"`
	MPRequestID                string                `form:"mp_request_id" validate:"required,uuid"`
	JobID                      string                `form:"job_id" validate:"required,uuid"`
	ForOrganizationID          string                `form:"for_organization_id" validate:"required,uuid"`
	ForOrganizationLocationID  string                `form:"for_organization_location_id" validate:"required,uuid"`
	DocumentNumber             string                `form:"document_number" validate:"required"`
	DocumentDate               string                `form:"document_date" validate:"required"`
	RecruitmentType            string                `form:"recruitment_type" validate:"required,recruitment_type_validation"`
	StartDate                  string                `form:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate                    string                `form:"end_date" validate:"required,datetime=2006-01-02"`
	Status                     string                `form:"status" validate:"required,job_posting_status_validation"`
	SalaryMin                  string                `form:"salary_min" validate:"required"`
	SalaryMax                  string                `form:"salary_max" validate:"required"`
	ContentDescription         string                `form:"content_description" validate:"omitempty"`
	OrganizationLogo           *multipart.FileHeader `form:"organization_logo" validate:"omitempty"`
	Poster                     *multipart.FileHeader `form:"poster" validate:"omitempty"`
	Link                       string                `form:"link" validate:"omitempty"`
	OrganizationLogoPath       string                `form:"organization_logo_path" validate:"omitempty"`
	PosterPath                 string                `form:"poster_path" validate:"omitempty"`
}

type UpdateJobPostingRequest struct {
	ID                         string `json:"id" validate:"required,uuid"`
	ProjectRecruitmentHeaderID string `json:"project_recruitment_header_id" validate:"required,uuid"`
	MPRequestID                string `json:"mp_request_id" validate:"required,uuid"`
	JobID                      string `json:"job_id" validate:"required,uuid"`
	ForOrganizationID          string `json:"for_organization_id" validate:"required,uuid"`
	ForOrganizationLocationID  string `json:"for_organization_location_id" validate:"required,uuid"`
	DocumentNumber             string `json:"document_number" validate:"required"`
	DocumentDate               string `json:"document_date" validate:"required"`
	RecruitmentType            string `json:"recruitment_type" validate:"required,recruitment_type_validation"`
	StartDate                  string `json:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate                    string `json:"end_date" validate:"required,datetime=2006-01-02"`
	Status                     string `json:"status" validate:"required,job_posting_status_validation"`
	SalaryMin                  string `json:"salary_min" validate:"required"`
	SalaryMax                  string `json:"salary_max" validate:"required"`
	ContentDescription         string `json:"content_description" validate:"omitempty"`
	OrganizationLogo           string `json:"organization_logo" validate:"omitempty"`
	Poster                     string `json:"poster" validate:"omitempty"`
	Link                       string `json:"link" validate:"omitempty"`
}
