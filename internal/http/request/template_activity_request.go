package request

type CreateTemplateActivityRequest struct {
	Name            string `json:"name" validate:"required"`
	Description     string `json:"description" validate:"omitempty"`
	RecruitmentType string `json:"recruitment_type" validate:"required,recruitment_type_validation"`
	Status          string `json:"status" validate:"required,template_activity_status_validation"`
}

type UpdateTemplateActivityRequest struct {
	ID              string `json:"id" validate:"required,uuid"`
	Name            string `json:"name" validate:"required"`
	Description     string `json:"description" validate:"omitempty"`
	RecruitmentType string `json:"recruitment_type" validate:"required,recruitment_type_validation"`
	Status          string `json:"status" validate:"required,template_activity_status_validation"`
}
