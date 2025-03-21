package request

type CreateTemplateQuestion struct {
	DocumentSetupID string `json:"document_setup_id" validate:"omitempty,uuid"`
	Name            string `json:"name" validate:"required"`
	FormType        string `json:"form_type" validate:"required,form_type_validation"`
	Description     string `json:"description" validate:"omitempty"`
	Duration        *int   `json:"duration" validate:"omitempty,gte=0"`
	Status          string `json:"status" validate:"required,template_question_status_validation"`
}

type UpdateTemplateQuestion struct {
	ID              string `json:"id" validate:"required,uuid"`
	DocumentSetupID string `json:"document_setup_id" validate:"omitempty,uuid"`
	Name            string `json:"name" validate:"required"`
	FormType        string `json:"form_type" validate:"required,form_type_validation"`
	Description     string `json:"description" validate:"omitempty"`
	Duration        *int   `json:"duration" validate:"omitempty,gte=0"`
	Status          string `json:"status" validate:"required,template_question_status_validation"`
}
