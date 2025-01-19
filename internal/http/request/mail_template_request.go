package request

type CreateMailTemplateRequest struct {
	Name           string `json:"name" validate:"required"`
	DocumentTypeID string `json:"document_type_id" validate:"required,uuid"`
	Subject        string `json:"subject" validate:"required"`
	Body           string `json:"body" validate:"required"`
}

type UpdateMailTemplateRequest struct {
	ID             string `json:"id" validate:"required,uuid"`
	Name           string `json:"name" validate:"required"`
	DocumentTypeID string `json:"document_type_id" validate:"required,uuid"`
	Subject        string `json:"subject" validate:"required"`
	Body           string `json:"body" validate:"required"`
}
