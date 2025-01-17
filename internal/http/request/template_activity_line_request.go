package request

type CreateOrUpdateTemplateActivityLineRequest struct {
	TemplateActivityID    string `json:"template_activity_id" validate:"required,uuid"`
	TemplateActivityLines []struct {
		ID                 string `json:"id" validate:"omitempty,uuid"`
		Name               string `json:"name" validate:"required"`
		Description        string `json:"description" validate:"omitempty"`
		Status             string `json:"status" validate:"required,template_activity_line_status_validation"`
		TemplateQuestionID string `json:"template_question_id" validate:"required,uuid"`
		ColorHexCode       string `json:"color_hex_code" validate:"omitempty"`
	} `json:"template_activity_lines" validate:"required,dive"`
	DeletedTemplateActivityLineIDs []string `json:"deleted_template_activity_line_ids" validate:"omitempty,dive,uuid"`
}
