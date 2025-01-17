package request

import (
	"time"
)

type CreateProjectRecruitmentHeader struct {
	TemplateActivityID string    `json:"template_activity_id" validate:"required,uuid"`
	Name               string    `json:"name" validate:"required"`
	Description        string    `json:"description" validate:"omitempty"`
	DocumentDate       time.Time `json:"document_date" validate:"required,datetime=2006-01-02"`
	DocumentNumber     string    `json:"document_number" validate:"required"`
	RecruitmentType    string    `json:"recruitment_type" validate:"required,recruitment_type_validation"`
	StartDate          time.Time `json:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate            time.Time `json:"end_date" validate:"required,datetime=2006-01-02"`
	Status             string    `json:"status" validate:"required,project_recruitment_header_status_validation"`
}

type UpdateProjectRecruitmentHeader struct {
	ID                 string    `json:"id" validate:"required,uuid"`
	TemplateActivityID string    `json:"template_activity_id" validate:"required,uuid"`
	Name               string    `json:"name" validate:"required"`
	Description        string    `json:"description" validate:"omitempty"`
	DocumentDate       time.Time `json:"document_date" validate:"required,datetime=2006-01-02"`
	DocumentNumber     string    `json:"document_number" validate:"required"`
	RecruitmentType    string    `json:"recruitment_type" validate:"required,recruitment_type_validation"`
	StartDate          time.Time `json:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate            time.Time `json:"end_date" validate:"required,datetime=2006-01-02"`
	Status             string    `json:"status" validate:"required,project_recruitment_header_status_validation"`
}
