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
	MinimumWorkExperience      string                `form:"minimum_work_experience" validate:"omitempty"`
	OrganizationLogoPath       string                `form:"organization_logo_path" validate:"omitempty"`
	PosterPath                 string                `form:"poster_path" validate:"omitempty"`
}

type UpdateJobPostingRequest struct {
	ID                         string                `form:"id" validate:"required,uuid"`
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
	MinimumWorkExperience      string                `form:"minimum_work_experience" validate:"omitempty"`
	OrganizationLogoPath       string                `form:"organization_logo_path" validate:"omitempty"`
	PosterPath                 string                `form:"poster_path" validate:"omitempty"`
	DeletedOrganizationLogo    string                `form:"deleted_organization_logo" validate:"omitempty"`
	DeletedPoster              string                `form:"deleted_poster" validate:"omitempty"`
}
