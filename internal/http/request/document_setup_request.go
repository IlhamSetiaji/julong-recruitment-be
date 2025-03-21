package request

type CreateDocumentSetupRequest struct {
	Title           string `json:"title" validate:"required"`
	DocumentTypeID  string `json:"document_type_id" validate:"required,uuid"`
	RecruitmentType string `json:"recruitment_type" validate:"required,recruitment_type_validation"`
	Header          string `json:"header" validate:"required"`
	Body            string `json:"body" validate:"required"`
	Footer          string `json:"footer" validate:"required"`
}

type UpdateDocumentSetupRequest struct {
	ID              string `json:"id" validate:"required,uuid"`
	Title           string `json:"title" validate:"required"`
	DocumentTypeID  string `json:"document_type_id" validate:"required,uuid"`
	RecruitmentType string `json:"recruitment_type" validate:"required,recruitment_type_validation"`
	Header          string `json:"header" validate:"required"`
	Body            string `json:"body" validate:"required"`
	Footer          string `json:"footer" validate:"required"`
}
