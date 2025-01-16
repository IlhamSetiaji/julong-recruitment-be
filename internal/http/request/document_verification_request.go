package request

type CreateDocumentVerificationRequest struct {
	TemplateQuestionID string `json:"template_question_id" validate:"required,uuid"`
	Name               string `json:"name" validate:"required"`
	Format             string `json:"format" validate:"required"`
}

type UpdateDocumentVerificationRequest struct {
	ID                 string `json:"id" validate:"required,uuid"`
	TemplateQuestionID string `json:"template_question_id" validate:"required,uuid"`
	Name               string `json:"name" validate:"required"`
	Format             string `json:"format" validate:"required"`
}
